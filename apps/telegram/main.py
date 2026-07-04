"""Telegram → Praxis NATS adapter.

Polls the Telegram Bot API and publishes each inbound text message to NATS
as an InputMessage JSON payload consumed by the Praxis kernel worker.

Configuration (environment variables):
  TELEGRAM_BOT_TOKEN    — required; Telegram Bot API token
  NATS_URL              — NATS server; default: nats://localhost:4222
  PRAXIS_INPUT_SUBJECT  — JetStream subject; default: praxis.kernel.input

Wire contract: the published JSON matches internal/transport/nats.InputMessage:
  {"id", "source", "text", "timestamp", "metadata"}
"""

from __future__ import annotations

import asyncio
import json
import logging
import os
import signal
import sys
from datetime import timezone

import nats
from telegram import Update
from telegram.ext import Application, ContextTypes, MessageHandler, filters

logger = logging.getLogger(__name__)


def load_config() -> tuple[str, str, str]:
    """Read and validate required configuration from environment variables.

    Returns (bot_token, nats_url, input_subject).
    Exits with a clear error if TELEGRAM_BOT_TOKEN is not set.
    """
    token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
    if not token:
        logger.error("TELEGRAM_BOT_TOKEN is required but not set")
        sys.exit(1)

    nats_url = os.getenv("NATS_URL", "nats://localhost:4222")
    input_subject = os.getenv("PRAXIS_INPUT_SUBJECT", "praxis.kernel.input")
    return token, nats_url, input_subject


def telegram_update_to_payload(update: Update) -> dict:
    """Map a Telegram Update to the Praxis InputMessage wire format.

    Owns all Telegram-specific field extraction. The returned dict matches
    the JSON schema of internal/transport/nats.InputMessage consumed by the
    Go worker.

    Precondition: update.message and update.message.text are not None.
    """
    msg = update.message
    chat_id = str(msg.chat_id)
    message_id = str(msg.message_id)
    user = msg.from_user

    return {
        "id": f"tg-{chat_id}-{message_id}",
        "source": "telegram",
        "text": msg.text,
        "timestamp": msg.date.astimezone(timezone.utc).isoformat(),
        "metadata": {
            "chat_id": chat_id,
            "message_id": message_id,
            "username": (user.username or "") if user else "",
            "first_name": (user.first_name or "") if user else "",
        },
    }


async def run(token: str, nats_url: str, input_subject: str) -> None:
    """Connect to NATS, start Telegram polling, run until SIGINT/SIGTERM."""
    try:
        nc = await nats.connect(nats_url)
    except Exception as exc:
        logger.error("failed to connect to NATS at %s: %s", nats_url, exc)
        sys.exit(1)

    logger.info("connected to NATS at %s", nats_url)
    js = nc.jetstream()

    async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        if not update.message or not update.message.text:
            return

        payload = telegram_update_to_payload(update)
        data = json.dumps(payload).encode()

        try:
            await js.publish(
                input_subject,
                data,
                headers={"Nats-Msg-Id": payload["id"]},
            )
            logger.info("published id=%s to %s", payload["id"], input_subject)
        except Exception as exc:
            logger.error(
                "publish failed id=%s subject=%s: %s",
                payload["id"],
                input_subject,
                exc,
            )

    app = Application.builder().token(token).build()
    app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_message))

    stop = asyncio.Event()
    loop = asyncio.get_running_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, stop.set)

    async with app:
        await app.start()
        await app.updater.start_polling()
        logger.info("telegram polling started subject=%s", input_subject)

        await stop.wait()
        logger.info("shutdown signal received")

        await app.updater.stop()
        await app.stop()

    await nc.drain()
    logger.info("shutdown complete")


def main() -> None:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(name)s %(message)s",
    )

    token, nats_url, input_subject = load_config()
    logger.info(
        "telegram adapter starting nats_url=%s subject=%s",
        nats_url,
        input_subject,
    )

    try:
        asyncio.run(run(token, nats_url, input_subject))
    except KeyboardInterrupt:
        pass


if __name__ == "__main__":
    main()

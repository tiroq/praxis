"""Telegram → Praxis NATS adapter.

Polls the Telegram Bot API and publishes each inbound text message to NATS
as an InputMessage JSON payload consumed by the Praxis kernel worker.

Configuration (environment variables):
  TELEGRAM_BOT_TOKEN    — required; Telegram Bot API token
  NATS_URL              — NATS server; default: nats://localhost:4222
  PRAXIS_INPUT_SUBJECT  — JetStream subject; default: praxis.kernel.input
  PRAXIS_OUTPUT_SUBJECT — Output subject for worker replies; default: praxis.kernel.output
  PRAXIS_DLQ_SUBJECT    — Dead-letter subject; default: praxis.kernel.dlq
  TELEGRAM_PUBLISH_MAX_RETRIES — Retry attempts per publish; default: 3
  TELEGRAM_RETRY_BASE_MS — Retry base delay in ms; default: 200
  NATS_MAX_RECONNECT_ATTEMPTS — NATS reconnect attempts; default: 60
  NATS_RECONNECT_WAIT_SECONDS — NATS reconnect wait seconds; default: 2
  TELEGRAM_HEALTH_HOST  — Health endpoint host; default: 0.0.0.0
  TELEGRAM_HEALTH_PORT  — Health/metrics endpoint port; default: 8090

Wire contract: the published JSON matches internal/transport/nats.InputMessage:
  {"id", "correlation_id", "source", "text", "timestamp", "metadata"}
"""

from __future__ import annotations

import asyncio
import json
import logging
import os
import signal
import sys
import threading
import time
from dataclasses import dataclass
from datetime import timezone
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any, Callable

import nats
from telegram import Update
from telegram.ext import Application, ContextTypes, MessageHandler, filters

logger = logging.getLogger(__name__)


@dataclass
class Config:
    token: str
    nats_url: str
    input_subject: str
    output_subject: str
    dlq_subject: str
    publish_max_retries: int
    retry_base_seconds: float
    nats_max_reconnect_attempts: int
    nats_reconnect_wait_seconds: float
    health_host: str
    health_port: int


class RuntimeState:
    def __init__(self) -> None:
        self._lock = threading.Lock()
        self.nats_connected = False
        self.polling_active = False
        self.last_error = ""
        self.last_message_unix = 0.0
        self.last_reply_unix = 0.0
        self.metrics = {
            "messages_received_total": 0,
            "publish_success_total": 0,
            "publish_fail_total": 0,
            "publish_retry_total": 0,
            "dlq_publish_total": 0,
            "dlq_publish_fail_total": 0,
            "output_messages_total": 0,
            "replies_sent_total": 0,
            "replies_failed_total": 0,
            "nats_reconnects_total": 0,
        }

    def inc(self, name: str, delta: int = 1) -> None:
        with self._lock:
            self.metrics[name] = self.metrics.get(name, 0) + delta

    def set_nats_connected(self, connected: bool) -> None:
        with self._lock:
            self.nats_connected = connected

    def set_polling_active(self, active: bool) -> None:
        with self._lock:
            self.polling_active = active

    def set_last_error(self, message: str) -> None:
        with self._lock:
            self.last_error = message

    def mark_message(self) -> None:
        with self._lock:
            self.last_message_unix = time.time()

    def mark_reply(self) -> None:
        with self._lock:
            self.last_reply_unix = time.time()

    def snapshot(self) -> dict[str, Any]:
        with self._lock:
            metrics_copy = dict(self.metrics)
            connected = self.nats_connected
            polling = self.polling_active
            error = self.last_error
            last_message = self.last_message_unix
            last_reply = self.last_reply_unix

        status = "ok" if connected and polling else "degraded"
        return {
            "status": status,
            "service": "praxis-telegram",
            "nats_connected": connected,
            "polling_active": polling,
            "last_error": error,
            "last_message_unix": last_message,
            "last_reply_unix": last_reply,
            "metrics": metrics_copy,
        }


def _parse_int(value: str, default: int, minimum: int) -> int:
    try:
        parsed = int(value)
    except ValueError:
        return default
    return parsed if parsed >= minimum else default


def _make_http_handler(state: RuntimeState) -> type[BaseHTTPRequestHandler]:
    class Handler(BaseHTTPRequestHandler):
        def do_GET(self) -> None:  # noqa: N802
            if self.path == "/health":
                snapshot = state.snapshot()
                payload = {
                    "status": snapshot["status"],
                    "service": snapshot["service"],
                    "nats_connected": snapshot["nats_connected"],
                    "polling_active": snapshot["polling_active"],
                    "last_error": snapshot["last_error"],
                }
                body = json.dumps(payload).encode("utf-8")
                self.send_response(HTTPStatus.OK)
                self.send_header("Content-Type", "application/json")
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)
                return

            if self.path == "/metrics":
                snapshot = state.snapshot()
                metrics = snapshot["metrics"]
                lines = [
                    "# TYPE praxis_telegram_messages_received_total counter",
                    f"praxis_telegram_messages_received_total {metrics['messages_received_total']}",
                    "# TYPE praxis_telegram_publish_success_total counter",
                    f"praxis_telegram_publish_success_total {metrics['publish_success_total']}",
                    "# TYPE praxis_telegram_publish_fail_total counter",
                    f"praxis_telegram_publish_fail_total {metrics['publish_fail_total']}",
                    "# TYPE praxis_telegram_publish_retry_total counter",
                    f"praxis_telegram_publish_retry_total {metrics['publish_retry_total']}",
                    "# TYPE praxis_telegram_dlq_publish_total counter",
                    f"praxis_telegram_dlq_publish_total {metrics['dlq_publish_total']}",
                    "# TYPE praxis_telegram_dlq_publish_fail_total counter",
                    f"praxis_telegram_dlq_publish_fail_total {metrics['dlq_publish_fail_total']}",
                    "# TYPE praxis_telegram_output_messages_total counter",
                    f"praxis_telegram_output_messages_total {metrics['output_messages_total']}",
                    "# TYPE praxis_telegram_replies_sent_total counter",
                    f"praxis_telegram_replies_sent_total {metrics['replies_sent_total']}",
                    "# TYPE praxis_telegram_replies_failed_total counter",
                    f"praxis_telegram_replies_failed_total {metrics['replies_failed_total']}",
                    "# TYPE praxis_telegram_nats_reconnects_total counter",
                    f"praxis_telegram_nats_reconnects_total {metrics['nats_reconnects_total']}",
                    "# TYPE praxis_telegram_nats_connected gauge",
                    f"praxis_telegram_nats_connected {1 if snapshot['nats_connected'] else 0}",
                    "# TYPE praxis_telegram_polling_active gauge",
                    f"praxis_telegram_polling_active {1 if snapshot['polling_active'] else 0}",
                ]
                body = ("\n".join(lines) + "\n").encode("utf-8")
                self.send_response(HTTPStatus.OK)
                self.send_header("Content-Type", "text/plain; version=0.0.4")
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)
                return

            self.send_response(HTTPStatus.NOT_FOUND)
            self.end_headers()

        def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
            logger.debug("health-server: " + format, *args)

    return Handler


def start_health_server(state: RuntimeState, host: str, port: int) -> ThreadingHTTPServer:
    handler = _make_http_handler(state)
    server = ThreadingHTTPServer((host, port), handler)
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    logger.info("health server started on %s:%d", host, port)
    return server


def load_config() -> Config:
    """Read and validate required configuration from environment variables.

    Returns full runtime Config.
    Exits with a clear error if TELEGRAM_BOT_TOKEN is not set.
    """
    token = os.getenv("TELEGRAM_BOT_TOKEN", "").strip()
    if not token:
        logger.error("TELEGRAM_BOT_TOKEN is required but not set")
        sys.exit(1)

    nats_url = os.getenv("NATS_URL", "nats://localhost:4222").strip()
    input_subject = os.getenv("PRAXIS_INPUT_SUBJECT", "praxis.kernel.input").strip()
    output_subject = os.getenv("PRAXIS_OUTPUT_SUBJECT", "praxis.kernel.output").strip()
    dlq_subject = os.getenv("PRAXIS_DLQ_SUBJECT", "praxis.kernel.dlq").strip()

    publish_max_retries = _parse_int(os.getenv("TELEGRAM_PUBLISH_MAX_RETRIES", "3"), 3, 1)
    retry_base_ms = _parse_int(os.getenv("TELEGRAM_RETRY_BASE_MS", "200"), 200, 1)
    max_reconnect_attempts = _parse_int(os.getenv("NATS_MAX_RECONNECT_ATTEMPTS", "60"), 60, 1)
    reconnect_wait_seconds = float(_parse_int(os.getenv("NATS_RECONNECT_WAIT_SECONDS", "2"), 2, 1))
    health_host = os.getenv("TELEGRAM_HEALTH_HOST", "0.0.0.0").strip() or "0.0.0.0"
    health_port = _parse_int(os.getenv("TELEGRAM_HEALTH_PORT", "8090"), 8090, 1)

    return Config(
        token=token,
        nats_url=nats_url,
        input_subject=input_subject,
        output_subject=output_subject,
        dlq_subject=dlq_subject,
        publish_max_retries=publish_max_retries,
        retry_base_seconds=retry_base_ms / 1000.0,
        nats_max_reconnect_attempts=max_reconnect_attempts,
        nats_reconnect_wait_seconds=reconnect_wait_seconds,
        health_host=health_host,
        health_port=health_port,
    )


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
        "correlation_id": f"telegram-chat-{chat_id}",
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


async def publish_with_retry(
    js: Any,
    subject: str,
    data: bytes,
    msg_id: str,
    max_retries: int,
    retry_base_seconds: float,
    state: RuntimeState,
) -> tuple[bool, Exception | None]:
    attempts = max(1, max_retries)
    for attempt in range(1, attempts + 1):
        try:
            await js.publish(subject, data, headers={"Nats-Msg-Id": msg_id})
            return True, None
        except Exception as exc:  # pragma: no cover - network path
            state.inc("publish_fail_total")
            state.set_last_error(str(exc))
            if attempt >= attempts:
                return False, exc
            state.inc("publish_retry_total")
            await asyncio.sleep(retry_base_seconds * (2 ** (attempt - 1)))
    return False, RuntimeError("publish retry loop exited unexpectedly")


def render_reply_text(output: dict[str, Any]) -> str:
    if output.get("status") == "error":
        error = output.get("error") or "pipeline failed"
        return f"Praxis error: {error}"

    assistant_reply = output.get("assistant_reply")
    if isinstance(assistant_reply, str) and assistant_reply.strip():
        return assistant_reply.strip()

    result = output.get("result")
    if isinstance(result, dict):
        decision = result.get("decision")
        actions = result.get("actions", [])
        if isinstance(decision, dict):
            outcome = decision.get("outcome", "unknown")
            action_count = len(actions) if isinstance(actions, list) else 0
            return f"Decision: {outcome}\nActions: {action_count}"

    return "Decision: unknown\nActions: 0"


async def handle_output_message(
    app: Application,
    state: RuntimeState,
    pending_chat_by_event_id: dict[str, str],
    output: dict[str, Any],
) -> None:
    state.inc("output_messages_total")
    input_event_id = str(output.get("input_event_id", ""))

    metadata = output.get("metadata")
    chat_id = ""
    if isinstance(metadata, dict):
        chat_id = str(metadata.get("chat_id", ""))
    if not chat_id:
        chat_id = pending_chat_by_event_id.get(input_event_id, "")

    if not chat_id:
        return

    try:
        await app.bot.send_message(chat_id=int(chat_id), text=render_reply_text(output))
        state.inc("replies_sent_total")
        state.mark_reply()
        pending_chat_by_event_id.pop(input_event_id, None)
    except Exception as exc:  # pragma: no cover - network path
        state.inc("replies_failed_total")
        state.set_last_error(str(exc))
        logger.error("reply failed input_event_id=%s chat_id=%s err=%s", input_event_id, chat_id, exc)


async def run(config: Config, state: RuntimeState) -> None:
    """Connect to NATS, start Telegram polling and output replies, run until SIGINT/SIGTERM."""
    pending_chat_by_event_id: dict[str, str] = {}

    def on_disconnected() -> None:
        state.set_nats_connected(False)
        logger.warning("nats disconnected")

    def on_reconnected() -> None:
        state.set_nats_connected(True)
        state.inc("nats_reconnects_total")
        logger.info("nats reconnected")

    def on_closed() -> None:
        state.set_nats_connected(False)
        logger.warning("nats connection closed")

    def on_async_error(exc: Exception) -> None:
        state.set_last_error(str(exc))
        logger.error("nats async error: %s", exc)

    try:
        nc = await nats.connect(
            servers=[config.nats_url],
            max_reconnect_attempts=config.nats_max_reconnect_attempts,
            reconnect_time_wait=config.nats_reconnect_wait_seconds,
            disconnected_cb=on_disconnected,
            reconnected_cb=on_reconnected,
            closed_cb=on_closed,
            error_cb=on_async_error,
        )
    except Exception as exc:
        state.set_last_error(str(exc))
        logger.error("failed to connect to NATS at %s: %s", config.nats_url, exc)
        sys.exit(1)

    state.set_nats_connected(True)
    logger.info("connected to NATS at %s", config.nats_url)
    js = nc.jetstream()

    app = Application.builder().token(config.token).build()

    async def handle_output(msg: Any) -> None:
        try:
            output = json.loads(msg.data.decode("utf-8"))
        except Exception as exc:  # pragma: no cover - network path
            state.set_last_error(str(exc))
            logger.error("invalid output json: %s", exc)
            return

        await handle_output_message(app, state, pending_chat_by_event_id, output)

    await nc.subscribe(config.output_subject, cb=handle_output)
    logger.info("subscribed to output subject=%s", config.output_subject)

    async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
        del context
        if not update.message or not update.message.text:
            return

        state.inc("messages_received_total")
        state.mark_message()

        payload = telegram_update_to_payload(update)
        data = json.dumps(payload).encode()

        ok, err = await publish_with_retry(
            js,
            config.input_subject,
            data,
            payload["id"],
            config.publish_max_retries,
            config.retry_base_seconds,
            state,
        )
        if ok:
            state.inc("publish_success_total")
            chat_id = payload.get("metadata", {}).get("chat_id", "")
            if chat_id:
                pending_chat_by_event_id[payload["id"]] = chat_id
            logger.info(
                "published id=%s correlation_id=%s subject=%s",
                payload["id"],
                payload.get("correlation_id", ""),
                config.input_subject,
            )
            return

        logger.error(
            "publish failed id=%s subject=%s after retries: %s",
            payload["id"],
            config.input_subject,
            err,
        )

        if not config.dlq_subject:
            return

        dlq_payload = {
            "input": payload,
            "error": str(err) if err else "publish failed",
            "failed_at": time.time(),
        }
        try:
            await js.publish(config.dlq_subject, json.dumps(dlq_payload).encode("utf-8"))
            state.inc("dlq_publish_total")
            logger.warning("published to dlq id=%s dlq_subject=%s", payload["id"], config.dlq_subject)
        except Exception as dlq_exc:  # pragma: no cover - network path
            state.inc("dlq_publish_fail_total")
            state.set_last_error(str(dlq_exc))
            logger.error(
                "dlq publish failed id=%s subject=%s: %s",
                payload["id"],
                config.dlq_subject,
                dlq_exc,
            )

    app.add_handler(MessageHandler(filters.TEXT & ~filters.COMMAND, handle_message))

    stop = asyncio.Event()
    loop = asyncio.get_running_loop()
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, stop.set)

    async with app:
        await app.start()
        await app.updater.start_polling()
        state.set_polling_active(True)
        logger.info(
            "telegram polling started input_subject=%s output_subject=%s",
            config.input_subject,
            config.output_subject,
        )

        await stop.wait()
        logger.info("shutdown signal received")

        state.set_polling_active(False)
        await app.updater.stop()
        await app.stop()

    state.set_nats_connected(False)
    await nc.drain()
    logger.info("shutdown complete")


def main() -> None:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)s %(name)s %(message)s",
    )

    config = load_config()
    state = RuntimeState()

    health_server = start_health_server(state, config.health_host, config.health_port)

    logger.info(
        "telegram adapter starting nats_url=%s input_subject=%s output_subject=%s dlq_subject=%s",
        config.nats_url,
        config.input_subject,
        config.output_subject,
        config.dlq_subject,
    )

    try:
        asyncio.run(run(config, state))
    except KeyboardInterrupt:
        pass
    finally:
        health_server.shutdown()
        health_server.server_close()


if __name__ == "__main__":
    main()

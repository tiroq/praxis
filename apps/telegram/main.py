"""Telegram bot entrypoint scaffold."""

from __future__ import annotations

import os


def main() -> None:
    token = os.getenv("TELEGRAM_BOT_TOKEN", "replace-me")
    print("Praxis Telegram scaffold")
    print("Status: scaffold only")
    print(f"Token configured: {'yes' if token != 'replace-me' else 'no'}")
    print("TODO: implement Telegram polling or webhook handling")
    print("TODO: map incoming messages to Praxis capture flows")


if __name__ == "__main__":
    main()

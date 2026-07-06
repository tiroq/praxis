"""Minimal LLM router service for assistant reply generation.

This service intentionally implements a single concrete vertical slice:
POST /v1/reply -> Ollama chat -> assistant_reply JSON.
"""

from __future__ import annotations

import json
import os
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any
from urllib import error as urlerror
from urllib import request as urlrequest


def _load_model_from_config() -> str:
    cfg_path = Path("configs/llm-routing.yaml")
    if not cfg_path.exists():
        return "gemma4:e2b"

    lines = cfg_path.read_text(encoding="utf-8").splitlines()
    in_summarizer = False
    for raw in lines:
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("summarizer:"):
            in_summarizer = True
            continue
        if in_summarizer and not raw.startswith("  "):
            break
        if in_summarizer and line.startswith("model:"):
            _, value = line.split(":", 1)
            model = value.strip()
            if model:
                return model
            break

    return "gemma4:e2b"


def _router_config() -> dict[str, Any]:
    return {
        "host": os.getenv("LLM_ROUTER_HOST", "0.0.0.0"),
        "port": int(os.getenv("LLM_ROUTER_PORT", "8081")),
        "ollama_url": os.getenv("OLLAMA_URL", "http://localhost:11434"),
        "model": os.getenv("LLM_ROUTER_MODEL", _load_model_from_config()),
        "timeout_seconds": float(os.getenv("LLM_ROUTER_TIMEOUT_SECONDS", "8")),
    }


def _generate_reply(input_text: str, cfg: dict[str, Any]) -> str:
    payload = {
        "model": cfg["model"],
        "system": "You are Praxis assistant. Reply concisely and helpfully.",
        "prompt": input_text,
        "stream": False,
        "options": {
            "num_predict": 256,
            "temperature": 0.2,
        },
    }

    req = urlrequest.Request(
        f"{cfg['ollama_url'].rstrip('/')}/api/generate",
        data=json.dumps(payload).encode("utf-8"),
        headers={"Content-Type": "application/json"},
        method="POST",
    )

    with urlrequest.urlopen(req, timeout=cfg["timeout_seconds"]) as resp:
        response_json = json.loads(resp.read().decode("utf-8"))

    content = response_json.get("response")
    if not isinstance(content, str) or not content.strip():
        raise ValueError("empty response content from provider")
    return content.strip()


class LLMRouterHandler(BaseHTTPRequestHandler):
    config = _router_config()

    def do_POST(self) -> None:  # noqa: N802
        if self.path != "/v1/reply":
            self._write_json(HTTPStatus.NOT_FOUND, {"error": "not found"})
            return

        try:
            length = int(self.headers.get("Content-Length", "0"))
        except ValueError:
            self._write_json(HTTPStatus.BAD_REQUEST, {"error": "invalid content-length"})
            return

        body = self.rfile.read(length)
        try:
            payload = json.loads(body.decode("utf-8"))
        except json.JSONDecodeError:
            self._write_json(HTTPStatus.BAD_REQUEST, {"error": "invalid json"})
            return

        input_text = payload.get("input_text", "")
        if not isinstance(input_text, str) or not input_text.strip():
            self._write_json(HTTPStatus.BAD_REQUEST, {"error": "input_text is required"})
            return

        try:
            assistant_reply = _generate_reply(input_text=input_text.strip(), cfg=self.config)
        except (urlerror.URLError, TimeoutError, ValueError) as exc:
            self._write_json(HTTPStatus.BAD_GATEWAY, {"error": str(exc)})
            return

        self._write_json(
            HTTPStatus.OK,
            {
                "assistant_reply": assistant_reply,
                "provider": "ollama",
                "model": self.config["model"],
            },
        )

    def do_GET(self) -> None:  # noqa: N802
        if self.path != "/health":
            self._write_json(HTTPStatus.NOT_FOUND, {"error": "not found"})
            return
        self._write_json(
            HTTPStatus.OK,
            {
                "status": "ok",
                "model": self.config["model"],
            },
        )

    def log_message(self, format: str, *args: Any) -> None:  # noqa: A003
        print(f"llm-router: {format % args}")

    def _write_json(self, status: HTTPStatus, payload: dict[str, Any]) -> None:
        data = json.dumps(payload).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        try:
            self.wfile.write(data)
        except BrokenPipeError:
            # Client canceled before reading response; safe to ignore.
            pass


def main() -> None:
    cfg = _router_config()
    LLMRouterHandler.config = cfg
    server = ThreadingHTTPServer((cfg["host"], cfg["port"]), LLMRouterHandler)
    print(f"Praxis LLM router listening on {cfg['host']}:{cfg['port']} model={cfg['model']}")
    server.serve_forever()


if __name__ == "__main__":
    main()

"""Minimal LLM router service for assistant reply generation and extraction.

This service intentionally implements a single concrete vertical slice:
POST /v1/reply -> Ollama chat -> reply_text JSON.
POST /v1/extract-facts -> Ollama structured extraction -> facts JSON.
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


def _generate_with_system(system: str, prompt: str, cfg: dict[str, Any], num_predict: int) -> str:
    payload = {
        "model": cfg["model"],
        "system": system,
        "prompt": prompt,
        "stream": False,
        "options": {
            "num_predict": num_predict,
            "temperature": 0.0,
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


def _build_input_text(payload: dict[str, Any]) -> str:
    user_message = payload.get("user_message")
    if not isinstance(user_message, str) or not user_message.strip():
        # Legacy compatibility with older workers.
        user_message = payload.get("input_text", "")
    if not isinstance(user_message, str) or not user_message.strip():
        return ""

    conversation = payload.get("conversation")
    if not isinstance(conversation, list) or len(conversation) == 0:
        return user_message.strip()

    lines: list[str] = []
    for item in conversation:
        if not isinstance(item, dict):
            continue
        role = item.get("role")
        text = item.get("text")
        if not isinstance(role, str) or not isinstance(text, str):
            continue
        role = role.strip()
        text = text.strip()
        if role and text:
            lines.append(f"{role}: {text}")

    lines.append(f"user: {user_message.strip()}")
    return "\n".join(lines)


def _build_fact_extraction_prompt(payload: dict[str, Any]) -> str:
    latest = payload.get("latest_user_message")
    if not isinstance(latest, str) or not latest.strip():
        return ""

    conversation = payload.get("conversation")
    lines: list[str] = []
    if isinstance(conversation, list):
        for item in conversation:
            if not isinstance(item, dict):
                continue
            role = item.get("role")
            text = item.get("text")
            if not isinstance(role, str) or not isinstance(text, str):
                continue
            role = role.strip()
            text = text.strip()
            if role and text:
                lines.append(f"{role}: {text}")

    lines.append(f"latest_user_message: {latest.strip()}")
    transcript = "\n".join(lines)
    return (
        "Extract only explicit, durable user facts from the latest user message.\n"
        "Do not infer facts from assistant messages. Do not summarize. Do not retrieve memory.\n"
        "Return JSON only in this exact shape: "
        '{"facts":[{"type":"location|preference|identity|relationship|profile|other",'
        '"value":"string","confidence":0.0}]}.\n'
        "If there are no explicit durable user facts, return {\"facts\":[]}.\n\n"
        f"Conversation:\n{transcript}"
    )


def _extract_json_object(text: str) -> dict[str, Any]:
    stripped = text.strip()
    start = stripped.find("{")
    end = stripped.rfind("}")
    if start == -1 or end == -1 or end < start:
        raise ValueError("provider did not return a JSON object")
    parsed = json.loads(stripped[start : end + 1])
    if not isinstance(parsed, dict):
        raise ValueError("provider JSON was not an object")
    return parsed


def _extract_facts(payload: dict[str, Any], cfg: dict[str, Any]) -> list[dict[str, Any]]:
    prompt = _build_fact_extraction_prompt(payload)
    if not prompt:
        raise ValueError("latest_user_message is required")

    raw = _generate_with_system(
        system="You extract structured candidate user facts. Return JSON only.",
        prompt=prompt,
        cfg=cfg,
        num_predict=256,
    )
    parsed = _extract_json_object(raw)
    facts = parsed.get("facts")
    if not isinstance(facts, list):
        raise ValueError("extraction response missing facts array")

    normalized: list[dict[str, Any]] = []
    for item in facts:
        if not isinstance(item, dict):
            continue
        fact_type = item.get("type")
        value = item.get("value")
        confidence = item.get("confidence")
        if not isinstance(fact_type, str) or not isinstance(value, str):
            continue
        if not isinstance(confidence, int | float):
            continue
        confidence_value = float(confidence)
        if confidence_value < 0 or confidence_value > 1:
            continue
        fact_type = fact_type.strip()
        value = value.strip()
        if not fact_type or not value:
            continue
        normalized.append({
            "type": fact_type,
            "value": value,
            "confidence": confidence_value,
        })

    return normalized


class LLMRouterHandler(BaseHTTPRequestHandler):
    config = _router_config()

    def do_POST(self) -> None:  # noqa: N802
        if self.path not in {"/v1/reply", "/v1/extract-facts"}:
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

        if self.path == "/v1/reply":
            self._handle_reply(payload)
            return

        self._handle_extract_facts(payload)

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

    def _handle_reply(self, payload: dict[str, Any]) -> None:
        input_text = _build_input_text(payload)
        if not input_text:
            self._write_json(HTTPStatus.BAD_REQUEST, {"error": "user_message is required"})
            return

        try:
            assistant_reply = _generate_reply(input_text=input_text.strip(), cfg=self.config)
        except (urlerror.URLError, TimeoutError, ValueError) as exc:
            self._write_json(HTTPStatus.BAD_GATEWAY, {"error": str(exc)})
            return

        self._write_json(
            HTTPStatus.OK,
            {
                "reply_text": assistant_reply,
                "assistant_reply": assistant_reply,
                "provider": "ollama",
                "model": self.config["model"],
            },
        )

    def _handle_extract_facts(self, payload: dict[str, Any]) -> None:
        try:
            facts = _extract_facts(payload, self.config)
        except ValueError as exc:
            self._write_json(HTTPStatus.BAD_GATEWAY, {"error": str(exc)})
            return
        except (urlerror.URLError, TimeoutError) as exc:
            self._write_json(HTTPStatus.BAD_GATEWAY, {"error": str(exc)})
            return

        self._write_json(
            HTTPStatus.OK,
            {
                "facts": facts,
                "provider": "ollama",
                "model": self.config["model"],
            },
        )


def main() -> None:
    cfg = _router_config()
    LLMRouterHandler.config = cfg
    server = ThreadingHTTPServer((cfg["host"], cfg["port"]), LLMRouterHandler)
    print(f"Praxis LLM router listening on {cfg['host']}:{cfg['port']} model={cfg['model']}")
    server.serve_forever()


if __name__ == "__main__":
    main()

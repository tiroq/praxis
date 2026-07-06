# Praxis LLM Router

Minimal Sprint 3 implementation for worker reply generation.

## Endpoints

- POST `/v1/reply`
	- Request JSON: `{"input_text": "..."}`
	- Response JSON: `{"assistant_reply": "...", "provider": "ollama", "model": "..."}`
- GET `/health`

## Runtime Config

- `LLM_ROUTER_HOST` (default `0.0.0.0`)
- `LLM_ROUTER_PORT` (default `8081`)
- `OLLAMA_URL` (default `http://localhost:11434`)
- `LLM_ROUTER_MODEL` (default from `configs/llm-routing.yaml` summarizer model, fallback `gemma4:e2b`)
- `LLM_ROUTER_TIMEOUT_SECONDS` (default `8`)

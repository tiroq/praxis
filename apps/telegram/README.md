# Telegram App

Telegram-facing NATS adapter for Praxis. Polls Telegram Bot API, publishes messages to the Kernel
via NATS JetStream, and sends results back to Telegram chats.

## Architecture

```text
Telegram Bot (user sends message)
  ↓
Telegram polling (long polling via Bot API)
  ↓
Message handler
  ↓
Mapper: telegram_update_to_payload() → InputMessage
  ↓
NATS JetStream Publish (praxis.kernel.input)
  ↓
[NATS Worker processes through Kernel]
  ↓
Output subscriber (praxis.kernel.output)
  ↓
Correlation via chat_id metadata
  ↓
Reply render: render_reply_text() → simple status message
  ↓
Telegram Bot send_message()
  ↓
Chat returns reply
```

## Configuration (environment variables)

| Variable | Default | Description |
|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | (required) | Telegram Bot API token from @BotFather |
| `NATS_URL` | `nats://localhost:4222` | NATS server address |
| `PRAXIS_INPUT_SUBJECT` | `praxis.kernel.input` | JetStream input subject |
| `PRAXIS_OUTPUT_SUBJECT` | `praxis.kernel.output` | JetStream output subject |
| `PRAXIS_DLQ_SUBJECT` | `praxis.kernel.dlq` | Dead-letter subject for failed publishes |
| `TELEGRAM_PUBLISH_MAX_RETRIES` | `3` | Retry attempts per input publish |
| `TELEGRAM_RETRY_BASE_MS` | `200` | Exponential backoff base (ms) |
| `NATS_MAX_RECONNECT_ATTEMPTS` | `60` | NATS reconnect attempts |
| `NATS_RECONNECT_WAIT_SECONDS` | `2` | NATS reconnect wait (seconds) |
| `TELEGRAM_HEALTH_HOST` | `0.0.0.0` | Health/metrics endpoint host |
| `TELEGRAM_HEALTH_PORT` | `8090` | Health/metrics endpoint port |

## Message Contract

**Input** (published to `praxis.kernel.input`):

```json
{
  "id": "tg-<chat_id>-<message_id>",
  "correlation_id": "telegram-chat-<chat_id>",
  "source": "telegram",
  "text": "<user message text>",
  "timestamp": "2026-07-01T12:34:56Z",
  "metadata": {
    "chat_id": "<numeric chat id>",
    "message_id": "<numeric message id>",
    "username": "<telegram username or empty>",
    "first_name": "<telegram first name or empty>"
  }
}
```

**Output** (received from `praxis.kernel.output`):

```json
{
  "input_event_id": "tg-<chat_id>-<message_id>",
  "correlation_id": "telegram-chat-<chat_id>",
  "status": "ok",
  "result": {
    "decision": {
      "id": "...",
      "outcome": "approve|reject|hold|escalate|defer|needs-revision"
    },
    "actions": [
      { "id": "...", ... }
    ]
  },
  "metadata": {
    "chat_id": "..."
  },
  "processed_at": "2026-07-01T12:34:57Z"
}
```

**Reply sent to chat:**

- Success: `Decision: <outcome>\nActions: <count>`
- Error: `Praxis error: <error>`

## Running

### Prerequisites

- Python 3.11+
- `nats-py>=2.3.0`
- `python-telegram-bot>=20.0`
- Active NATS server (see worker/README.md for setup)
- Telegram Bot API token (get from @BotFather)

### Start the Telegram Adapter

```bash
# Install dependencies
pip install -r apps/telegram/requirements.txt

# Set required env vars
export TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN_HERE"
export NATS_URL="nats://localhost:4222"

# Run
python3 apps/telegram/main.py
```

Health check:

```bash
curl http://localhost:8090/health
curl http://localhost:8090/metrics
```

## Testing

### Unit Tests

Tests mapper and output rendering (no live Telegram or NATS):

```bash
python3 -m unittest apps.telegram.test_main -v
```

Tests verify:
- `telegram_update_to_payload()` creates correct wire format
- `render_reply_text()` formats decisions and actions correctly
- `handle_output_message()` correlates outputs back to chats
- Metadata extraction and fallback chain

### End-to-End Verification

Requires running NATS worker. See `task smoke:nats` in Praxis Taskfile.

# Praxis Worker

The worker is the NATS JetStream transport adapter for the Core Kernel.

It subscribes to an input subject, drives each message through the kernel pipeline
(`Event → Review → Decision → Action`), and publishes the result to an output subject.
The kernel itself remains transport-agnostic — all NATS concerns live in
`internal/transport/nats/`.

## Architecture

```text
NATS JetStream (input subject)
  ↓
transport adapter (internal/transport/nats)
  ↓  decode JSON → validate → convert to kernel.Event
Core Kernel (internal/core/kernel)
  ↓  Review → Decision → Action
transport adapter
  ↓  encode OutputMessage → publish
NATS JetStream (output subject)
```

## Configuration (env vars)

| Variable               | Default                   | Description                        |
|------------------------|---------------------------|------------------------------------|
| `NATS_URL`             | `nats://localhost:4222`   | NATS server address                |
| `NATS_STREAM`          | `PRAXIS`                  | JetStream stream name              |
| `NATS_INPUT_SUBJECT`   | `praxis.kernel.input`     | Subject the worker subscribes to   |
| `NATS_OUTPUT_SUBJECT`  | `praxis.kernel.output`    | Subject results are published to   |
| `NATS_DURABLE`         | `praxis-worker`           | Durable consumer name              |
| `NATS_ACK_WAIT_SECONDS`| `30`                      | Ack wait in seconds                |
| `NATS_MAX_DELIVER`     | `3`                       | Max delivery attempts              |
| `PRAXIS_STORAGE_BACKEND` | `sqlite`                | Storage backend: `memory` or `sqlite` |
| `PRAXIS_SQLITE_PATH`   | `build/praxis.db`         | Path to SQLite database file       |

### Storage Configuration

The worker optionally records kernel pipeline execution events to a persistent EventStore.
Event recording is configured via environment variables and is **non-fatal** — if storage
fails to open, the worker logs a warning and continues without event recording.

The worker also optionally persists conversation history as a **SQLite projection** used for
reconstructing chat threads. This history is derived from immutable events and is rebuildable.
The EventStore remains the canonical event source of truth.

- **`PRAXIS_STORAGE_BACKEND`**: Storage backend to use (`memory` or `sqlite`).
  - `memory` — In-memory EventStore; events are not persisted to disk.
  - `sqlite` — SQLite-backed EventStore; events are persisted to the configured database file.
  - Default: `sqlite`

- **`PRAXIS_SQLITE_PATH`**: Path to the SQLite database file (only used when `PRAXIS_STORAGE_BACKEND=sqlite`).
  - Default: `build/praxis.db`

When storage is successfully opened, the worker injects an `EventRecorder` into the kernel via
`kernel.WithEventRecorder()`. The kernel then records a `kernel.pipeline.completed` event after
each successful pipeline execution.

When SQLite storage is enabled, the worker also opens a conversation projection store at
`PRAXIS_SQLITE_PATH` and the subscriber appends both user and assistant messages for each
processed input. Conversation persistence is non-fatal: failures are logged and the reply loop
continues.

**Example:**

```sh
# Use SQLite storage (default)
PRAXIS_STORAGE_BACKEND=sqlite PRAXIS_SQLITE_PATH=build/praxis.db go run ./services/worker

# Use in-memory storage (no persistence)
PRAXIS_STORAGE_BACKEND=memory go run ./services/worker

# Disable storage (not recommended; for testing only)
PRAXIS_STORAGE_BACKEND=invalid go run ./services/worker  # Logs error, continues without storage
```

## Message Contracts

**Input** (`praxis.kernel.input`):

```json
{
  "id": "evt_...",
  "source": "manual",
  "text": "нужно купить билеты в Шанхай",
  "timestamp": "2026-07-01T00:00:00Z",
  "metadata": {}
}
```

**Output** (`praxis.kernel.output`):

```json
{
  "input_event_id": "evt_...",
  "status": "ok",
  "result": { ... },
  "error": null,
  "processed_at": "2026-07-01T00:00:01Z"
}
```

## Ack / Nak Contract

- **Ack** — sent only after a successful publish to the output subject.
- **Nak** — sent if publish fails; message is redelivered up to `MaxDeliver` times.
- **Term** — sent for poison messages (invalid JSON or missing required fields) to
  prevent endless redelivery.

## Running

```sh
# With task runner (no NATS server needed for unit tests):
task build:worker    # compile to build/worker
task run:worker      # go run ./services/worker

# Real end-to-end smoke verification over local JetStream:
task smoke:nats

# With a real NATS server:
NATS_URL=nats://localhost:4222 go run ./services/worker
```

## Unit Tests

No real NATS server is required. All adapter behaviour is tested with in-process
fakes in `internal/transport/nats/nats_test.go`.

```sh
go test ./internal/transport/nats/...
```

`task test` does not require a live NATS server. Real NATS is exercised only by
`task smoke:nats`, which boots local JetStream, starts the worker, publishes a
real input event, and validates the real output message.

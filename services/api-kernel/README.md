# api-kernel

Minimal HTTP wrapper around the Core Kernel pipeline for local testing.

Exposes the `Event → Review → Decision → Action` pipeline over HTTP with no external dependencies.

## Run

```sh
go run ./services/api-kernel
```

Server listens on `:8080`.

## Endpoint

### `POST /v1/kernel/run`

**Request**

```json
{
  "text": "нужно купить билеты в Шанхай",
  "source": "manual"
}
```

- `text` — required, non-empty human-readable input
- `source` — optional, defaults to `"api"`

**Response (200)**

`PipelineResult` as JSON — the complete auditable trace of the run.

**Error (400)**

```json
{ "error": "text must not be empty" }
```

**Error (500)**

```json
{ "error": "kernel error: ..." }
```

## Examples

```sh
# Happy path
curl -X POST localhost:8080/v1/kernel/run \
  -H 'Content-Type: application/json' \
  -d '{"text":"нужно купить билеты в Шанхай","source":"manual"}'

# Missing text → 400
curl -X POST localhost:8080/v1/kernel/run \
  -H 'Content-Type: application/json' \
  -d '{"text":""}'
```

## Tests

```sh
go test ./services/api-kernel/...
```

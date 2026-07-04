# Praxis Configuration Reference

**Complete listing of all `PRAXIS_*` environment variables, organized by subsystem.**

This document is the authoritative reference for Praxis runtime configuration. See [ADR-012](docs/adr/ADR-012-runtime-configuration-naming.md) for rationale and design.

---

## Overview

All Praxis environment variables use the `PRAXIS_<SUBSYSTEM>_<SETTING>` convention.

**Quick lookup:**
- [NATS (messaging)](#nats-subsystem-praxis_nats_)
- [Storage (persistence)](#storage-subsystem-praxis_storage_)
- [Integrations](#integrations)
- [Defaults table](#all-variables-defaults-and-types)

---

## NATS Subsystem (`PRAXIS_NATS_*`)

NATS JetStream configuration for the kernel transport layer.

| Variable | Default | Type | Mutable | Description |
|----------|---------|------|---------|-------------|
| `PRAXIS_NATS_URL` | `nats://localhost:4222` | string | No | NATS server address (scheme + host + port) |
| `PRAXIS_NATS_STREAM` | `PRAXIS` | string | No | JetStream stream name |
| `PRAXIS_NATS_INPUT_SUBJECT` | `praxis.kernel.input` | string | No | Subject for inbound messages to kernel |
| `PRAXIS_NATS_OUTPUT_SUBJECT` | `praxis.kernel.output` | string | No | Subject for outbound kernel results |
| `PRAXIS_NATS_DURABLE` | `praxis-worker` | string | No | Durable consumer name (for pull-subscribe) |
| `PRAXIS_NATS_ACK_WAIT_SECONDS` | `30` | integer | No | Acknowledgement timeout in seconds |
| `PRAXIS_NATS_MAX_DELIVER` | `3` | integer | No | Maximum redelivery attempts before dead-letter |

### Examples

```bash
# Default (local development)
export PRAXIS_NATS_URL=nats://localhost:4222
export PRAXIS_NATS_STREAM=PRAXIS
export PRAXIS_NATS_INPUT_SUBJECT=praxis.kernel.input
export PRAXIS_NATS_OUTPUT_SUBJECT=praxis.kernel.output

# Production (remote NATS cluster)
export PRAXIS_NATS_URL=nats://nats-1:4222,nats-2:4222,nats-3:4222
export PRAXIS_NATS_STREAM=PRAXIS_PROD
export PRAXIS_NATS_INPUT_SUBJECT=prod.kernel.input
export PRAXIS_NATS_OUTPUT_SUBJECT=prod.kernel.output
export PRAXIS_NATS_ACK_WAIT_SECONDS=60
```

### Affected Services

- `services/worker` — subscriber (pull consumer)
- `services/api-kernel` — publisher (Go)
- `apps/telegram` — publisher (Python)
- `cmd/praxis` (CLI) — publisher
- `cmd/nats-smoke` — end-to-end verification

---

## Storage Subsystem (`PRAXIS_STORAGE_*`)

Configuration for the event store (persistence layer).

| Variable | Default | Type | Mutable | Description |
|----------|---------|------|---------|-------------|
| `PRAXIS_STORAGE_BACKEND` | `sqlite` | string | No | Backend type: `sqlite`, `memory` |
| `PRAXIS_STORAGE_SQLITE_PATH` | `build/praxis.db` | string | No | File path for SQLite database (only used if backend=sqlite) |

### Allowed Values

**`PRAXIS_STORAGE_BACKEND`:**
- `sqlite` — persistent SQLite database
- `memory` — ephemeral in-memory store (loses data on restart)
- Any other value → treated as "disabled"; service logs warning and continues without storage

### Examples

```bash
# Production (SQLite, persistent)
export PRAXIS_STORAGE_BACKEND=sqlite
export PRAXIS_STORAGE_SQLITE_PATH=/data/praxis.db

# Development (in-memory, fast, no persistence)
export PRAXIS_STORAGE_BACKEND=memory

# Disable storage (not recommended)
export PRAXIS_STORAGE_BACKEND=disabled
# → service logs: "storage backend 'disabled' not recognized; continuing without storage"
```

### Affected Services

- `services/worker` — records pipeline execution events (optional)

---

## Integrations

### Telegram Adapter (`PRAXIS_TELEGRAM_*`)

Configuration for the Telegram Bot API adapter.

| Variable | Default | Type | Mutable | Secret | Description |
|----------|---------|------|---------|--------|-------------|
| `PRAXIS_TELEGRAM_BOT_TOKEN` | (required) | string | No | **Yes** | Telegram Bot API token (https://t.me/BotFather) |

**Important:** This is a secret. Store in:
- Docker Secrets or Environment in production
- `.env` file locally (git-ignored)
- Kubernetes Secret in K8s deployments
- CI/CD secret vault

```bash
# Local development (.env, git-ignored)
PRAXIS_TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklmnoPQRstUVwxyz

# Docker Compose (pass from host)
# docker-compose run telegram
# (set via .env or -e on command line)

# Kubernetes
apiVersion: v1
kind: Secret
metadata:
  name: praxis-secrets
type: Opaque
data:
  PRAXIS_TELEGRAM_BOT_TOKEN: MTIzNDU2Nzg5OkFCQ2RlZkdISWprbG1ub1BRUnN0VVZ3eHl6
```

### Affected Services

- `apps/telegram` — Telegram polling adapter

### Future Integrations (Not Yet Implemented)

#### OpenAI Integration (`PRAXIS_OPENAI_*`)

Planned for LLM routing subsystem.

| Variable | Default | Type | Secret | Description |
|----------|---------|------|--------|-------------|
| `PRAXIS_OPENAI_API_KEY` | (required if enabled) | string | **Yes** | OpenAI API key |
| `PRAXIS_OPENAI_ORG_ID` | (optional) | string | No | OpenAI organization ID |

#### Anthropic Integration (`PRAXIS_ANTHROPIC_*`)

Planned for LLM routing subsystem.

| Variable | Default | Type | Secret | Description |
|----------|---------|------|--------|-------------|
| `PRAXIS_ANTHROPIC_API_KEY` | (required if enabled) | string | **Yes** | Anthropic Claude API key |

---

## Database Subsystem (`PRAXIS_DB_*`) — *Future*

Planned for direct Postgres integration (currently using migrations in `infra/postgres/init.sql`).

| Variable | Default | Type | Mutable | Secret | Description |
|----------|---------|------|---------|--------|-------------|
| `PRAXIS_DB_URL` | `postgres://localhost/praxis` | string | No | **Yes** | Postgres connection string (scheme://user:password@host:port/db) |
| `PRAXIS_DB_USERNAME` | `praxis` | string | No | **Yes** | Postgres user (alternative to embedding in URL) |
| `PRAXIS_DB_PASSWORD` | (required) | string | No | **Yes** | Postgres password (alternative to embedding in URL) |
| `PRAXIS_DB_POOL_MIN` | `2` | integer | No | No | Minimum connection pool size |
| `PRAXIS_DB_POOL_MAX` | `10` | integer | No | No | Maximum connection pool size |

---

## Logging Subsystem (`PRAXIS_LOG_*`) — *Future*

Planned observability configuration.

| Variable | Default | Type | Description |
|----------|---------|------|-------------|
| `PRAXIS_LOG_LEVEL` | `info` | string | Log level: `debug`, `info`, `warn`, `error` |
| `PRAXIS_LOG_FORMAT` | `text` | string | Log format: `text` or `json` |

---

## HTTP/API Subsystem (`PRAXIS_HTTP_*`) — *Future*

Planned for API server binding.

| Variable | Default | Type | Description |
|----------|---------|------|-------------|
| `PRAXIS_HTTP_HOST` | `0.0.0.0` | string | Bind address |
| `PRAXIS_HTTP_PORT` | `8080` | integer | Bind port |
| `PRAXIS_HTTP_TLS_CERT` | (optional) | string | Path to TLS certificate file |
| `PRAXIS_HTTP_TLS_KEY` | (optional) | string | Path to TLS key file |

---

## LLM Routing Subsystem (`PRAXIS_LLM_*`) — *Future*

Planned for model routing and provider selection.

| Variable | Default | Type | Description |
|----------|---------|------|-------------|
| `PRAXIS_LLM_DEFAULT_PROVIDER` | `anthropic` | string | Default provider when no routing applies |
| `PRAXIS_LLM_TIMEOUT_SECONDS` | `30` | integer | API call timeout |
| `PRAXIS_LLM_RETRY_ATTEMPTS` | `3` | integer | Number of retry attempts on transient failure |

---

## All Variables: Defaults and Types

### Quick Reference Table

| Variable | Default | Type | Secret | Service(s) |
|----------|---------|------|--------|-----------|
| `PRAXIS_NATS_URL` | `nats://localhost:4222` | string | No | worker, api-kernel, CLI, telegram |
| `PRAXIS_NATS_STREAM` | `PRAXIS` | string | No | worker, api-kernel, CLI, smoke |
| `PRAXIS_NATS_INPUT_SUBJECT` | `praxis.kernel.input` | string | No | worker, api-kernel, CLI, telegram |
| `PRAXIS_NATS_OUTPUT_SUBJECT` | `praxis.kernel.output` | string | No | worker, api-kernel, CLI |
| `PRAXIS_NATS_DURABLE` | `praxis-worker` | string | No | worker, api-kernel |
| `PRAXIS_NATS_ACK_WAIT_SECONDS` | `30` | integer | No | worker, api-kernel |
| `PRAXIS_NATS_MAX_DELIVER` | `3` | integer | No | worker, api-kernel |
| `PRAXIS_STORAGE_BACKEND` | `sqlite` | string | No | worker |
| `PRAXIS_STORAGE_SQLITE_PATH` | `build/praxis.db` | string | No | worker |
| `PRAXIS_TELEGRAM_BOT_TOKEN` | (required) | string | **Yes** | telegram |

---

## Configuration Precedence (ADR-008)

Environment variables are read in this order of precedence:

```
code defaults → YAML config files → environment variables → CLI flags
                                    ↑
                           (this layer)
```

When a service starts:
1. Loads **code defaults** (hardcoded in the service).
2. Merges in **YAML config** from `configs/` (if applicable).
3. **Overrides** with **environment variables** (`PRAXIS_*` vars set in this run).
4. **Overrides** with **CLI flags** (if service supports them).

Higher layers completely override lower layers *per key*; missing keys fall through to lower layers.

---

## Best Practices

### Secrets Management

- **Never commit secrets to Git.** Use `.env` files (git-ignored) locally.
- **Use secret stores in production:** Kubernetes Secrets, Docker Secrets, CI/CD vault, AWS Secrets Manager, etc.
- **Mark secrets clearly:** All variables ending in `_TOKEN` or `_PASSWORD` are secrets.
- **Rotate regularly:** Especially API keys like `PRAXIS_TELEGRAM_BOT_TOKEN` and `PRAXIS_OPENAI_API_KEY`.

### Development vs. Production

**Development (local Docker Compose):**
```bash
# Use defaults; minimal overrides needed
export PRAXIS_NATS_URL=nats://nats:4222
export PRAXIS_STORAGE_BACKEND=memory
```

**Staging (Docker Compose on server):**
```bash
# More conservative; persistent storage
export PRAXIS_NATS_URL=nats://nats-server:4222
export PRAXIS_STORAGE_BACKEND=sqlite
export PRAXIS_STORAGE_SQLITE_PATH=/mnt/data/praxis.db
```

**Production (Kubernetes):**
```yaml
# Use ConfigMaps for non-secrets, Secrets for secrets
apiVersion: v1
kind: ConfigMap
metadata:
  name: praxis-config
data:
  PRAXIS_NATS_STREAM: PRAXIS_PROD
  PRAXIS_NATS_INPUT_SUBJECT: prod.kernel.input
  PRAXIS_STORAGE_BACKEND: sqlite
---
apiVersion: v1
kind: Secret
metadata:
  name: praxis-secrets
type: Opaque
data:
  PRAXIS_NATS_URL: (base64-encoded)
  PRAXIS_TELEGRAM_BOT_TOKEN: (base64-encoded)
```

### Debugging Configuration

To see all `PRAXIS_*` variables in use:

```bash
env | grep ^PRAXIS_
```

To see a service's effective configuration:

```bash
# Go services usually log config at startup
go run ./services/worker 2>&1 | grep -i config

# Python services
PRAXIS_LOG_LEVEL=debug python apps/telegram/main.py 2>&1 | grep config
```

---

## Migration

See [MIGRATION_GUIDE-ADR-012.md](docs/MIGRATION_GUIDE-ADR-012.md) for guidance on migrating from old variable names (e.g., `NATS_URL`, `TELEGRAM_BOT_TOKEN`) to new `PRAXIS_*` names.

---

## Related Documentation

- [ADR-012: Runtime Configuration Naming Convention](docs/adr/ADR-012-runtime-configuration-naming.md)
- [ADR-008: Configuration Strategy](docs/adr/ADR-008-configuration-strategy.md)
- [RFC-030 §13: Cross-Cutting Concerns](rfcs/030-system-architecture.md)
- Service READMEs:
  - [services/worker/README.md](services/worker/README.md)
  - [services/api-kernel/README.md](services/api-kernel/README.md)
  - [apps/telegram/README.md](apps/telegram/README.md)

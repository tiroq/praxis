# Migration Guide: ADR-012 Runtime Configuration Naming Convention

This guide helps you migrate from current inconsistent environment variable naming to the standardized `PRAXIS_*` convention.

---

## Quick Reference: Old → New Variable Mapping

| Old Name | New Name | Subsystem | Service(s) |
|----------|----------|-----------|-----------|
| `NATS_URL` | `PRAXIS_NATS_URL` | NATS | worker, api-kernel, CLI, smoke tests |
| `NATS_STREAM` | `PRAXIS_NATS_STREAM` | NATS | worker, api-kernel, CLI, smoke tests |
| `NATS_INPUT_SUBJECT` | `PRAXIS_NATS_INPUT_SUBJECT` | NATS | worker, api-kernel, CLI |
| `NATS_OUTPUT_SUBJECT` | `PRAXIS_NATS_OUTPUT_SUBJECT` | NATS | worker, api-kernel, CLI |
| `NATS_DURABLE` | `PRAXIS_NATS_DURABLE` | NATS | worker, api-kernel |
| `NATS_ACK_WAIT_SECONDS` | `PRAXIS_NATS_ACK_WAIT_SECONDS` | NATS | worker, api-kernel |
| `NATS_MAX_DELIVER` | `PRAXIS_NATS_MAX_DELIVER` | NATS | worker, api-kernel |
| `PRAXIS_STORAGE_BACKEND` | `PRAXIS_STORAGE_BACKEND` | Storage | worker *(already correct)* |
| `PRAXIS_SQLITE_PATH` | `PRAXIS_STORAGE_SQLITE_PATH` | Storage | worker *(rename for consistency)* |
| `TELEGRAM_BOT_TOKEN` | `PRAXIS_TELEGRAM_BOT_TOKEN` | Telegram | Telegram adapter |
| `PRAXIS_INPUT_SUBJECT` | `PRAXIS_NATS_INPUT_SUBJECT` | NATS | Telegram adapter *(consolidate with NATS subsystem)* |

---

## Migration Timeline

### Phase 1: Parallel Support (v0.2.0)
- Services accept **both** old and new variable names.
- New names take precedence; old names trigger deprecation warnings.
- No breaking changes; existing deployments continue to work.
- **Action:** Update your `docker-compose.yml`, Kubernetes manifests, or CI/CD to use new names at your pace.

### Phase 2: Hard Cutover (v0.3.0+)
- Old variable names are no longer supported; services error on startup if detected.
- **Action:** All deployments must use new names.

---

## Migration Steps

### For Docker Compose Users

**Before (v0.1.x):**
```yaml
services:
  worker:
    environment:
      NATS_URL: nats://nats:4222
      NATS_STREAM: PRAXIS
      NATS_INPUT_SUBJECT: praxis.kernel.input
      NATS_OUTPUT_SUBJECT: praxis.kernel.output
      PRAXIS_STORAGE_BACKEND: sqlite
      PRAXIS_SQLITE_PATH: build/praxis.db

  telegram:
    environment:
      TELEGRAM_BOT_TOKEN: "${TELEGRAM_BOT_TOKEN}"
      NATS_URL: nats://nats:4222
      PRAXIS_INPUT_SUBJECT: praxis.kernel.input
```

**After (v0.2.0+, recommended):**
```yaml
services:
  worker:
    environment:
      PRAXIS_NATS_URL: nats://nats:4222
      PRAXIS_NATS_STREAM: PRAXIS
      PRAXIS_NATS_INPUT_SUBJECT: praxis.kernel.input
      PRAXIS_NATS_OUTPUT_SUBJECT: praxis.kernel.output
      PRAXIS_STORAGE_BACKEND: sqlite
      PRAXIS_STORAGE_SQLITE_PATH: build/praxis.db

  telegram:
    environment:
      PRAXIS_TELEGRAM_BOT_TOKEN: "${TELEGRAM_BOT_TOKEN}"
      PRAXIS_NATS_URL: nats://nats:4222
      PRAXIS_NATS_INPUT_SUBJECT: praxis.kernel.input
```

### For Kubernetes / Helm Users

**Before (v0.1.x):**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: praxis-config
data:
  NATS_STREAM: "PRAXIS"
  NATS_INPUT_SUBJECT: "praxis.kernel.input"
  NATS_OUTPUT_SUBJECT: "praxis.kernel.output"
---
apiVersion: v1
kind: Secret
metadata:
  name: praxis-secrets
type: Opaque
data:
  NATS_URL: bmF0czovL25hdHM6NDIyMg==  # base64(nats://nats:4222)
  TELEGRAM_BOT_TOKEN: MTIzNDU2Nzg5MA==  # base64(bot token)
```

**After (v0.2.0+, recommended):**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: praxis-config
data:
  PRAXIS_NATS_STREAM: "PRAXIS"
  PRAXIS_NATS_INPUT_SUBJECT: "praxis.kernel.input"
  PRAXIS_NATS_OUTPUT_SUBJECT: "praxis.kernel.output"
---
apiVersion: v1
kind: Secret
metadata:
  name: praxis-secrets
type: Opaque
data:
  PRAXIS_NATS_URL: bmF0czovL25hdHM6NDIyMg==
  PRAXIS_TELEGRAM_BOT_TOKEN: MTIzNDU2Nzg5MA==
```

### For CI/CD Pipelines

**Before (v0.1.x):**
```bash
#!/bin/bash
export NATS_URL="nats://ci-nats:4222"
export NATS_STREAM="PRAXIS"
export NATS_INPUT_SUBJECT="praxis.kernel.input"
export PRAXIS_STORAGE_BACKEND="memory"
go test ./...
```

**After (v0.2.0+, recommended):**
```bash
#!/bin/bash
export PRAXIS_NATS_URL="nats://ci-nats:4222"
export PRAXIS_NATS_STREAM="PRAXIS"
export PRAXIS_NATS_INPUT_SUBJECT="praxis.kernel.input"
export PRAXIS_STORAGE_BACKEND="memory"
go test ./...
```

---

## Troubleshooting

### Deprecation Warnings in Logs (Phase 1)

If you see:
```
WARN: env var NATS_URL is deprecated; use PRAXIS_NATS_URL instead
```

**Action:** Update your configuration to use the new name. The old name still works during v0.2.0, but support will be removed in v0.3.0.

### Services Not Starting (Phase 2+)

If you see:
```
ERROR: env var NATS_URL is no longer supported. Use PRAXIS_NATS_URL instead.
       See migration guide: docs/MIGRATION_GUIDE-ADR-012.md
```

**Action:** Your deployment is using Praxis v0.3.0+ which no longer accepts old variable names. Update your environment to use new `PRAXIS_*` names (see table above).

### Mixed Old/New Names

**Phase 1 behavior:** If both old and new names are set, the new name takes precedence.

Example:
```bash
export NATS_URL=nats://old:4222
export PRAXIS_NATS_URL=nats://new:4222
# → PRAXIS_NATS_URL is used; NATS_URL is ignored
```

**Phase 2:** This situation cannot occur; only new names are accepted.

---

## Validation Checklist

After migrating to the new naming convention, verify:

- [ ] All `docker-compose.yml` files use `PRAXIS_*` names
- [ ] All Kubernetes manifests (ConfigMaps, Secrets) use `PRAXIS_*` names
- [ ] All CI/CD pipeline scripts use `PRAXIS_*` names
- [ ] Helm charts (if maintained) use `PRAXIS_*` names
- [ ] No services log deprecation warnings at startup
- [ ] All services start cleanly with only new variable names set
- [ ] Documentation has been updated to reflect new names

---

## Questions?

Refer to:
- **ADR-012:** Full architectural rationale and design
- **docs/configuration-reference.md:** Complete list of all `PRAXIS_*` variables
- Service READMEs (e.g., `services/worker/README.md`): service-specific configuration details

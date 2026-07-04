# ADR-012 — Runtime Configuration Naming Convention

**Status:** PROPOSED
**Date:** 2026-07-04
**Authors:** Architecture
**Supersedes:** —
**Superseded by:** —
**Related:** ADR-008 (Configuration Strategy), ADR-007 (Storage Architecture), RFC-030 §13 (Cross-Cutting Concerns)

---

## Purpose

This ADR defines the **runtime configuration naming language for the entire Praxis platform**. It establishes a uniform convention for all environment variables across all services, languages, and subsystems, both current and future.

**Scope:** Runtime configuration names only. This ADR does NOT govern service contracts, NATS subjects, event schemas, or APIs — those remain under RFC authority.

---

## Decision Authority

This ADR governs:

- **Environment variable naming**: All `PRAXIS_*` variables across all services
- **Runtime configuration naming**: The naming language for deployment-time configuration
- **Subsystem naming**: How subsystems are identified in the configuration namespace
- **Configuration ownership**: Which subsystem owns the interpretation of each variable

This ADR does **not** govern:

- **Configuration values**: Specific defaults, ranges, or constraints (owned by service implementations)
- **Loading order or precedence**: How configuration is merged across sources (governed by ADR-008)
- **Service contracts**: NATS subjects, event schemas, APIs, or protocol definitions (governed by RFCs)
- **Runtime behavior**: What services do with configuration; only how they name it
- **Implementation techniques**: Programming languages, frameworks, or loading mechanisms remain local to each service

**Authority hierarchy:**
- RFC-030 defines subsystem concepts; this ADR names them
- ADR-008 defines configuration precedence; this ADR names the env var layer
- This ADR is final authority on naming conventions
- Service implementations are free to evolve provided they respect the naming convention

---

## Context

### The Problem

Current Praxis services use inconsistent environment variable naming conventions:

| Service | Variables | Namespace |
|---------|-----------|-----------|
| Go (worker, CLI, api-kernel) | `NATS_URL`, `NATS_INPUT_SUBJECT`, `NATS_OUTPUT_SUBJECT`, `NATS_STREAM`, `NATS_DURABLE`, `NATS_ACK_WAIT_SECONDS`, `NATS_MAX_DELIVER` | `NATS_*` |
| Go (storage) | `PRAXIS_STORAGE_BACKEND`, `PRAXIS_SQLITE_PATH` | `PRAXIS_*` |
| Python (Telegram adapter) | `TELEGRAM_BOT_TOKEN`, `NATS_URL`, `PRAXIS_INPUT_SUBJECT` | Mixed: `TELEGRAM_*`, `NATS_*`, `PRAXIS_*` |

### Consequences of Inconsistency

1. **Deployment friction:** operators must remember different prefixes for different services.
2. **Kubernetes manifests:** ConfigMaps and Secrets lack consistent naming, complicating templating.
3. **Helm charts:** mixed prefixes require custom logic instead of systematic templating.
4. **Documentation:** each service requires separate env var documentation; no unified reference.
5. **CI/CD pipelines:** environment variable lookup requires service-specific hardcoding.
6. **Future services:** each new language/framework starts without a clear convention.
7. **Namespace pollution:** external vars like `NATS_URL` can collide with other NATS-aware tools or testing frameworks.

## Explicit Non-Goals

This ADR explicitly does **NOT** introduce:

1. **Configuration library or package:** Each service owns its configuration loading. No shared `config.py` module or Go constants package.
2. **Configuration framework:** No `ConfigManager`, `ConfigLoader`, `ConfigProvider`, `ConfigRegistry`, or abstraction layer. Services use standard `os.Getenv()` (Go) or `os.getenv()` (Python).
3. **Dependency injection container:** No central container for configuration objects. Each service constructs its own Config struct/dict.
4. **Global configuration singleton:** No shared state. Each service reads and holds its own configuration instance.
5. **Reflection-based loading:** No automatic field mapping or auto-wiring of configuration variables.
6. **Configuration validation framework:** Validation remains in service startup code. No external schema validator.
7. **Configuration reload or hot-swap:** Configuration is read at startup; changing env vars requires service restart.
8. **Centralized configuration service or API:** Each service reads from its environment; no runtime configuration server.
9. **Mandatory configuration file format:** Services may use YAML, TOML, JSON, or plain env vars; this ADR only standardizes env var *names*, not the configuration substrate (per ADR-008).

**Rationale for each:** Shared configuration code would create service dependencies and violate the Extract/Don't Invent principle. Naming is language-agnostic and scalable; frameworks are not.

---

## ADR Compliance Rules

Every new environment variable introduced into Praxis **must** satisfy all of the following:

1. **Naming convention:** Uses the pattern `PRAXIS_<SUBSYSTEM>_<SETTING_NAME>` with no exceptions
2. **Subsystem assignment:** Belongs to exactly one existing subsystem or justifies a new subsystem bounded by an RFC or ADR
3. **Ownership:** Has exactly one documented architectural owner (from the subsystem's defining RFC/ADR)
4. **Documentation:** Documented in the service README AND in `docs/configuration-reference.md`
5. **Default or required:** Has either a safe default value OR explicit "(required)" designation with explanation
6. **Migration coverage:** If a variable is renamed, a deprecation phase and migration guide are provided
7. **No ADR violations:** Does not contradict this ADR or any superseding ADR

Variables that fail any of these checks **must be rejected in code review**.

---

## Decision

**All externally visible runtime configuration across Praxis services MUST use the `PRAXIS_` namespace prefix, organized by subsystem.**

Environment variables exposed by Praxis services follow a strict naming convention:

```
PRAXIS_<SUBSYSTEM>_<SETTING_NAME>
```

This naming convention applies uniformly across all languages (Go, Python, Rust, Java, etc.), all subsystems, and all future services. No exceptions.

### Naming Convention

#### Subsystems (non-exhaustive; extensible)

New subsystems are added as the platform evolves. All must follow the `PRAXIS_<SUBSYSTEM>_*` pattern.

| Subsystem | Prefix | Examples | Owner |
|-----------|--------|----------|----------|
| NATS (JetStream, broker) | `PRAXIS_NATS_` | `PRAXIS_NATS_URL`, `PRAXIS_NATS_STREAM`, `PRAXIS_NATS_INPUT_SUBJECT` | Transport/Messaging subsystem |
| Storage (persistence) | `PRAXIS_STORAGE_` | `PRAXIS_STORAGE_BACKEND`, `PRAXIS_STORAGE_SQLITE_PATH` | Storage subsystem |
| Telegram adapter | `PRAXIS_TELEGRAM_` | `PRAXIS_TELEGRAM_BOT_TOKEN` | Telegram adapter |
| Worker service | `PRAXIS_WORKER_` | `PRAXIS_WORKER_CONCURRENCY`, `PRAXIS_WORKER_TIMEOUT` | Worker service |
| API service | `PRAXIS_API_` | `PRAXIS_API_HOST`, `PRAXIS_API_PORT` | API service |
| LLM routing | `PRAXIS_LLM_` | `PRAXIS_LLM_DEFAULT_PROVIDER`, `PRAXIS_LLM_TIMEOUT` | LLM routing subsystem |
| HTTP/Web | `PRAXIS_HTTP_` | `PRAXIS_HTTP_HOST`, `PRAXIS_HTTP_PORT` | HTTP transport |
| Logging/Observability | `PRAXIS_LOG_` | `PRAXIS_LOG_LEVEL`, `PRAXIS_LOG_FORMAT` | Observability subsystem |
| Tracing | `PRAXIS_TRACE_` | `PRAXIS_TRACE_ENABLED`, `PRAXIS_TRACE_SAMPLE_RATE` | Observability subsystem |
| Database (Postgres) | `PRAXIS_DB_` | `PRAXIS_DB_URL`, `PRAXIS_DB_USERNAME`, `PRAXIS_DB_PASSWORD` | Storage subsystem |
| Authentication | `PRAXIS_AUTH_` | `PRAXIS_AUTH_ENABLED`, `PRAXIS_AUTH_SECRET_KEY` | Auth subsystem |
| OpenAI integration | `PRAXIS_OPENAI_` | `PRAXIS_OPENAI_API_KEY`, `PRAXIS_OPENAI_ORG_ID` | OpenAI adapter |
| Anthropic integration | `PRAXIS_ANTHROPIC_` | `PRAXIS_ANTHROPIC_API_KEY` | Anthropic adapter |

#### Setting Names

- Use `SNAKE_CASE_UPPERCASE`.
- Be specific: `PRAXIS_NATS_INPUT_SUBJECT` not `PRAXIS_NATS_INPUT`.
- Suffix boolean flags with `_ENABLED` or `_DISABLED` if ambiguous: `PRAXIS_STORAGE_ENABLED`.
- Use `_URL` for endpoints (not `_HOST` alone): `PRAXIS_NATS_URL`.
- Use full words: `PRAXIS_NATS_ACK_WAIT_SECONDS` not `PRAXIS_NATS_ACK_WAIT_SECS`.

#### Namespace Growth Rule

Subsystems should remain stable and focused.

Do **not** create namespace prefixes like:

```
PRAXIS_MISC_*         (hides ownership)
PRAXIS_COMMON_*       (no single owner)
PRAXIS_SHARED_*       (anti-pattern)
PRAXIS_UTIL_*         (violates Extract, don't invent)
PRAXIS_INTERNAL_*     (internal to what?)
PRAXIS_TEMP_*         (temporary configuration is a smell)
PRAXIS_CONFIG_*       (redundant; all env vars are config)
```

These prefixes hide architectural ownership and violate the principle that naming should clarify, not obscure.

Create a new subsystem only when:
1. Architecture clearly defines a new bounded context (in RFC or ADR)
2. The subsystem owns one or more configuration variables
3. Ownership is obvious to operators and developers

**Example of justified new subsystem:** `PRAXIS_LLM_*` because RFC-041 defines LLM routing as a distinct subsystem.

**Example of unjustified namespace:** `PRAXIS_UTIL_CACHE_*` because "util" hides ownership; use `PRAXIS_CACHE_*` instead.

---

**Before (current inconsistency):**
```bash
export NATS_URL=nats://localhost:4222
export NATS_STREAM=PRAXIS
export NATS_INPUT_SUBJECT=praxis.kernel.input
export PRAXIS_STORAGE_BACKEND=sqlite
export PRAXIS_TELEGRAM_BOT_TOKEN=123456:ABC-DEF
```

**After (standardized):**
```bash
export PRAXIS_NATS_URL=nats://localhost:4222
export PRAXIS_NATS_STREAM=PRAXIS
export PRAXIS_NATS_INPUT_SUBJECT=praxis.kernel.input
export PRAXIS_STORAGE_BACKEND=sqlite
export PRAXIS_TELEGRAM_BOT_TOKEN=123456:ABC-DEF
```

---

## Ownership Model

Every runtime configuration variable has exactly one **architectural owner**. Ownership is architectural, not organizational.

### Examples

```
PRAXIS_STORAGE_BACKEND
Owner: Storage subsystem (RFC-033 §6 Configuration Store)
Meaning: Selects storage backend implementation (sqlite, memory, etc.)
Loading: Each application reads and interprets locally

PRAXIS_NATS_URL
Owner: Transport/NATS subsystem
Meaning: Specifies NATS server address and port
Loading: Worker reads via internal/transport/nats/config.go; Telegram reads via apps/telegram/main.py

PRAXIS_TELEGRAM_BOT_TOKEN
Owner: Telegram adapter
Meaning: Authenticates with Telegram Bot API
Loading: Telegram app reads via apps/telegram/main.py

PRAXIS_LOG_LEVEL
Owner: Observability subsystem
Meaning: Controls logging verbosity (debug, info, warn, error)
Loading: Each service reads and applies locally
```

### Ownership Principles

1. **One owner per variable:** no shared governance of a single variable.
2. **Architectural ownership:** ownership is determined by the subsystem/RFC that governs the functionality, not by which team manages the code.
3. **Local loading:** ownership does not require or justify centralized loading. The owner defines *what* the variable means; the application defines *how* to load it.
4. **Semantic clarity:** the variable name must make ownership unmistakable (e.g., `PRAXIS_STORAGE_*` owned by Storage, not Application).

#### Ownership Rule

**Ownership is architectural. Ownership is NOT determined by:**

- **Repository:** A variable owned by Storage subsystem belongs to the storage subsystem regardless of whether code lives in `internal/storage/`, `packages/storage/`, or `services/storage/`.
- **Team:** Ownership is not about which team maintains the code; it's about which subsystem owns the configuration's interpretation.
- **Deployment:** A variable does not belong to "production team" or "staging team"; it belongs to the subsystem that interprets it.
- **Implementation language:** A variable owned by NATS subsystem is owned by NATS whether the reader is Go, Python, Rust, or Java.
- **Loading mechanism:** The owner does not determine whether configuration is loaded via environment variables, YAML files, CLI flags, or shared code; the owner only determines the meaning.

**Ownership is determined by:**
- The RFC or ADR that defines the subsystem responsible for interpreting the variable
- The architectural principle that the variable implements
- The question: "Which subsystem's behavior changes when this variable changes?"

**Example:** `PRAXIS_NATS_URL` is owned by the Transport/NATS subsystem (RFC-030 §13), not by the Worker service that reads it. If the Worker were replaced with a Scheduler, ownership remains with Transport/NATS.

---

## Runtime Configuration vs. Service Contracts

**Critical invariant:** Runtime configuration is separate from service contracts.

### Configuration Controls Only

Runtime variables control **where** to connect, **what** backend to use, **how verbose** to log, and similar operational concerns.

**Configuration does NOT control:**
- NATS subject names (e.g., `praxis.kernel.input`) — governed by RFC-032 (Data Flow)
- Event schemas and JSON contracts — governed by RFC-013 (Event Model)
- API endpoints and request/response formats — governed by RFC-031 (Service Contracts)
- Message payload structures — governed by relevant RFCs

### Examples

| Type | Example | Governed by | Immutable? | Notes |
|------|---------|-------------|-----------|-------|
| Configuration | `PRAXIS_NATS_URL=nats://prod:4222` | ADR-012 | No | Can change per deployment |
| Contract | Subject `praxis.kernel.input` | RFC-032 | Yes (protocol stability) | Cannot change without RFC-032 update |
| Configuration | `PRAXIS_STORAGE_BACKEND=sqlite` | ADR-012 | No | Can change per deployment |
| Contract | Event schema `{"id", "source", "text", …}` | RFC-013 | Yes (append-only) | Cannot change; must extend |
| Configuration | `PRAXIS_LOG_LEVEL=debug` | ADR-012 | No | Can change per deployment |
| Contract | HTTP API response `{"status", "result"}` | RFC-031 | Yes (contract stability) | Cannot break without version bump |

### Why This Matters

**Contracts require coordination and stability; configuration does not.**

- Changing `PRAXIS_NATS_URL` from one deployment to another is safe and expected.
- Changing the `praxis.kernel.input` subject name breaks all producers and consumers; requires RFC amendment.
- Changing `PRAXIS_STORAGE_BACKEND` affects only the running instance; changing Event Store schema affects all historical replay.

**Configuration belongs to operations; contracts belong to architecture.**

---

## Rationale

### Why `PRAXIS_` Namespace

1. **Ownership clarity:** `PRAXIS_*` explicitly signals "this configuration belongs to Praxis"; external users know not to set it unless running Praxis.
2. **Collision avoidance:** prevents accidental shadowing of external tool variables (`NATS_*` is used by NATS testing frameworks, `TELEGRAM_*` by other bots).
3. **Consistency across languages:** Go, Python, Rust, Java — all use the same prefix, eliminating language-specific conventions.
4. **Kubernetes/Helm:** ConfigMaps and Secrets can be templated uniformly; no special case per language.
5. **CI/CD:** scripts can use systematic env var detection: `env | grep ^PRAXIS_` instead of service-specific lists.
6. **Documentation:** single canonical reference: "All Praxis env vars start with `PRAXIS_`".

### Why Subsystem Suffix

1. **Logical grouping:** related settings cluster together (`PRAXIS_NATS_*` groups all NATS config).
2. **Prevents abbreviation:** full subsystem name (`NATS` not `N`) makes intent unmistakable.
3. **Scales:** adding new subsystems (e.g., `PRAXIS_REDIS_*`) does not require renaming existing vars.
4. **Documentation clarity:** env var reference can be organized by subsystem.

### The Praxis Configuration Pattern

Praxis follows a specific pattern for configuration management:

```
Every application owns its configuration loading.

worker:
    ConfigFromEnv()  (internal/transport/nats/config.go)

telegram:
    load_config()  (apps/telegram/main.py)

scheduler:
    load_config()  (services/scheduler/main.py)

api:
    load_config()  (services/api/main.py)
```

**Only the naming convention is centralized. Loading and interpretation are local.**

This pattern ensures:
- Services remain independent
- No hidden configuration dependencies
- Easy to audit where configuration is loaded
- Easy to add new services without framework changes
- Type safety within each service (Go Config struct, Python dict, etc.)

### Why NOT a Shared Configuration Framework

**This ADR explicitly does NOT introduce:**

- `internal/config` package
- `config.Manager`, `config.Loader`, `config.Provider`, `config.Registry`
- Configuration dependency injection container
- Global configuration singleton
- Reflection-based configuration loading
- Configuration file format standardization (Go: struct, Python: dict, future services: free to choose)
- Centralized configuration service or API

**Rationale:**

1. **Extract, don't invent:** Only one Python service exists; no duplication yet to justify shared code.
2. **Ownership remains local:** Each service owns its defaults, loading logic, and validation.
3. **Type safety:** Go uses Config structs; Python uses dicts; future services choose what fits their language.
4. **Future flexibility:** If/when multiple Python or Rust services emerge, extraction is trivial; forcing it now increases upfront burden.
5. **Simplicity:** Naming is a language-agnostic convention; loading is language-specific.
6. **Service independence:** Shared configuration code creates dependencies between services; local loading keeps them independent.

**Standardizing names does NOT justify centralizing implementation.**

---

## Architectural Invariants

This ADR establishes the following non-negotiable architecture laws:

### Invariant 1: Configuration Ownership Separation

```
Configuration belongs to applications.
Naming belongs to architecture.
Loading belongs to applications.
Interpretation belongs to the owning subsystem.
```

Any new configuration variable must be owned by exactly one subsystem (defined by RFC or ADR). Ownership determines *what* the variable means, not *who* implements loading.

### Invariant 2: Standardizing Names Does Not Justify Centralizing Implementation

```
Uniform naming DOES NOT imply shared loading code.
Uniform naming DOES NOT imply shared constants modules.
Uniform naming DOES NOT imply shared frameworks.
```

If a future version of Praxis has 5 Python microservices, they may all follow `PRAXIS_*` naming without sharing configuration code. Extraction happens only when **proven duplication** justifies the coupling cost.

### Invariant 3: Service Independence Through Local Configuration

```
Runtime configuration MUST NOT create dependencies between services.

Services may share naming conventions.
Services must not share configuration implementations.
```

Example of **allowed** (independent):
```python
# telegram/main.py
token = os.getenv("PRAXIS_TELEGRAM_BOT_TOKEN", "")

# scheduler/main.py
backend = os.getenv("PRAXIS_STORAGE_BACKEND", "sqlite")
```

Example of **forbidden** (creates dependency):
```python
# Both telegram and scheduler importing from a shared config module
from praxis.config import load_telegram_config, load_storage_config
# → Creates coupling between independent services
```

---

---

## Consequences

### Positive

- **Consistency:** all future services follow the same convention; no surprises.
- **Deployment simplicity:** operators set `PRAXIS_*` vars; no need to learn service-specific prefixes.
- **Kubernetes-friendly:** ConfigMaps and Secrets named `praxis-config` and `praxis-secrets` work uniformly.
- **CI/CD clarity:** automated scripts can find all Praxis configuration in one env var scan.
- **Documentation:** single unified reference page: "Praxis Configuration Reference" organized by subsystem.
- **Third-party tools:** monitoring/debugging tools (e.g., `env | grep PRAXIS_`) work immediately.

### Negative

- **Backward compatibility break:** existing deployments using `NATS_*`, `TELEGRAM_*` must migrate (mitigated by 2-phase rollout below).
- **Slightly longer var names:** `PRAXIS_NATS_URL` > `NATS_URL` (negligible; readability > brevity).

### Trade-offs

- Accepts a small breaking change (1-2 release cycle migration window) for long-term consistency and clarity.
- Defers shared configuration abstractions until duplication justifies extraction (Extract, don't invent).

---

## Migration Strategy

### Phase 1: Parallel Support (Compatibility Mode)

**Behavior:** Services read BOTH old and new variable names, with deterministic precedence.

**Precedence rule:**
```
if PRAXIS_<SUBSYSTEM>_<SETTING> is set:
    use it
elif legacy variable (e.g., NATS_SETTING) is set:
    use it (log deprecation warning)
else:
    use code default
```

**Deprecation messaging:**
- Log a **WARN** message at startup: `"env var NATS_URL is deprecated; use PRAXIS_NATS_URL instead"`
- Include link to migration guide in warning
- Services continue to function; migration is not forced

**Example (Go):**
```go
func ConfigFromEnv() Config {
    cfg := DefaultConfig()
    
    // New naming (priority)
    if v := os.Getenv("PRAXIS_NATS_URL"); v != "" {
        cfg.URL = v
    } else if v := os.Getenv("NATS_URL"); v != "" {
        // Deprecated
        log.Warn("NATS_URL is deprecated; use PRAXIS_NATS_URL instead")
        cfg.URL = v
    }
    // ... repeat for other fields
    return cfg
}
```

**Example (Python):**
```python
def load_config():
    nats_url = (
        os.getenv("PRAXIS_NATS_URL") or
        os.getenv("NATS_URL") or
        "nats://localhost:4222"
    )
    if os.getenv("NATS_URL") and not os.getenv("PRAXIS_NATS_URL"):
        logger.warning("NATS_URL is deprecated; use PRAXIS_NATS_URL")
    return nats_url
```

**Documentation:** docker-compose.yml and Helm charts use new names; include commented-out legacy names with deprecation notices.

### Phase 2: Hard Cutover (Breaking Change)

**Behavior:** Legacy variables are no longer accepted; services error on startup if detected.

```
if legacy variable (e.g., NATS_URL) is set:
    raise error: "NATS_URL is no longer supported. Use PRAXIS_NATS_URL instead. Migration guide: [link]"
```

**Communication:**
- Release notes: "Breaking Change: legacy environment variable names no longer supported"
- Link to migration guide in error message
- Update all documentation and examples
- Recommend: allow 1–2 release cycles of Phase 1 before Phase 2

---

## Migration Phases: Behavioral Definition

| Phase | Behavior | Old vars | New vars | Status |
|-------|----------|----------|----------|--------|
| **Phase 1** | Both work; deprecation warnings | ✅ Supported (warn) | ✅ Supported (preferred) | Accepting both names |
| **Phase 2** | Only new names accepted | ❌ Error | ✅ Supported (required) | Names fully migrated |

---

## Compatibility Rule

This ADR establishes how changes to configuration variables are classified:

**Changing the meaning of an existing variable is a breaking architectural change.**

- Example: `PRAXIS_NATS_URL` previously meant "local NATS" but now means "remote NATS" → Breaking change
- Impact: Services expecting one behavior get another; potential data loss or errors
- Mitigation: Create a new variable (`PRAXIS_NATS_REMOTE_URL`) instead; use deprecation phase for old one

**Changing only the default value is not necessarily a breaking architectural change.**

- Example: `PRAXIS_LOG_LEVEL` default changes from "info" to "debug"
- Impact: New instances log more verbosely; existing instances with explicit settings are unaffected
- Mitigation: Document in release notes; use cautiously

**Renaming a variable requires migration.**

- Example: `PRAXIS_NATS_URL` becomes `PRAXIS_NATS_ENDPOINT_URL`
- Impact: Existing deployments using old name must be updated
- Mitigation: Two-phase migration (Phase 1: both work, Phase 2: old name removed)

**Removing a variable requires a deprecation phase.**

- Example: `PRAXIS_NATS_DEPRECATED_AUTH` is no longer used
- Impact: Deployments setting this variable should not (silent no-op) or cannot (error)
- Mitigation: Phase 1 ignores it (no-op with warning), Phase 2 errors on it

**Governance:** Changes to variable meaning or removal must be reviewed by the architecture team and documented in an ADR amendment or new ADR.

---

## Review Checklist for Configuration Changes

When reviewing any change that introduces or modifies environment variables, the code reviewer must verify:

**Naming and Ownership:**
- [ ] Does this belong to an existing subsystem (e.g., NATS, Storage, Telegram)?
- [ ] If a new subsystem is needed, is it justified by an RFC or ADR that defines its boundaries?
- [ ] Is the owner clear from the variable name? Could operators immediately identify which subsystem owns it?

**Duplication and Coherence:**
- [ ] Does this duplicate another existing variable? If so, why is a new one needed?
- [ ] Is the variable externally visible (operators can/should set it), or internal (applications should not expose it)?
- [ ] Does the variable belong to a single subsystem, or does it straddle multiple owners (sign of poor design)?

**Compatibility and Migration:**
- [ ] Is this a new variable (no migration needed) or a rename (migration required)?
- [ ] If renaming, is a deprecation phase planned? Which release cycles?
- [ ] If changing meaning, does the commit message explain why this is not breaking?

**Governance Compliance:**
- [ ] Does it follow the pattern `PRAXIS_<SUBSYSTEM>_<SETTING_NAME>`?
- [ ] Does it have a documented owner (from an RFC or ADR)?
- [ ] Is it documented in the service README and in `docs/configuration-reference.md`?
- [ ] Does it have either a safe default or explicit "(required)" designation?
- [ ] Does it violate ADR-012, ADR-008, or any other architecture decision?

**Reject if** any item is unclear or fails.

---

## Implementation Order

### Step 1: ADR Acceptance
- [ ] This ADR is reviewed and accepted by the architecture team.
- [ ] ADR status changed to ACCEPTED.
- [ ] Link to this ADR is added to the README under "Configuration".

### Step 2: Phase 1 Implementation (Parallel Support)
- [ ] Update `internal/transport/nats/config.go` to support both `NATS_*` (deprecated) and `PRAXIS_NATS_*` (new) with deprecation warnings.
- [ ] Update Python services to support both old and new names with warnings.
- [ ] Create `MIGRATION_GUIDE.md` documenting old → new variable mapping.
- [ ] Update `docker-compose.yml` to use new names (comment out old names with notes).
- [ ] Update service READMEs (worker, api, cli, etc.) with env var tables using new names.
- [ ] Create `docs/configuration-reference.md` listing all `PRAXIS_*` variables by subsystem.
- [ ] Add CHANGELOG entry: "Phase 1: Parallel support for old and new environment variable names. Deprecation warnings added."

### Step 3: Monitoring and Feedback (Phase 1)
- [ ] Deploy Phase 1 code to staging/production.
- [ ] Monitor deprecation warning logs.
- [ ] Solicit feedback from operators about migration path.
- [ ] Document any edge cases or concerns.

### Step 4: Phase 2 Implementation (Breaking Change)
- [ ] Remove deprecation-mode code; raise errors on legacy vars.
- [ ] Update all documentation to remove legacy names.
- [ ] CHANGELOG: "Breaking Change: legacy environment variable names no longer supported. All `PRAXIS_*` names required."

### Step 5: Ongoing Governance (for all future services)
- [ ] Every new service must use `PRAXIS_*` naming from day 1; no compatibility mode.
- [ ] Code review checklist item: "Environment variable names follow `PRAXIS_<SUBSYSTEM>_*` convention."
- [ ] New subsystems must be documented in `docs/configuration-reference.md` before release.

---

## Implementation Checklist

### Phase 1: Parallel Support (All items required before Phase 2)

**Code:**
- [ ] `internal/transport/nats/config.go`: read `PRAXIS_NATS_*` first, then `NATS_*` (warn on legacy)
- [ ] `cmd/nats-smoke/main.go`: update config structs and logging
- [ ] `services/worker/main.go`: update config reading and logging
- [ ] `services/api-kernel/main.go`: update config reading (if applicable)
- [ ] `internal/cli/praxiscli/app.go`: update config reading (if applicable)
- [ ] `apps/telegram/main.py`: support both old and new names; warn on legacy
- [ ] All services: add deprecation log calls with link to migration guide

**Configuration Files:**
- [ ] Update `docker-compose.yml` to use new names; comment out old names
- [ ] Update `infra/*/` configuration files if they reference env vars
- [ ] Update Kubernetes manifests (if maintained) to use new names

**Documentation:**
- [ ] Update `services/worker/README.md`: env var table with new names
- [ ] Update `services/api-kernel/README.md`: env var table
- [ ] Create/update `MIGRATION_GUIDE.md` with old→new mapping
- [ ] Create/update `docs/configuration-reference.md` with all `PRAXIS_*` variables
- [ ] Add section to main `README.md`: link to configuration reference
- [ ] Update any other service-specific README files

**Tests:**
- [ ] Add test cases: both old and new vars work, new takes precedence
- [ ] Add test assertions for deprecation warnings
- [ ] Verify startup logs contain migration link

**CI/CD:**
- [ ] Update GitHub Actions workflows to use new variable names
- [ ] Verify no hardcoded old variable names in scripts

### Phase 2: Breaking Change (After Phase 1 monitoring period)

**Code:**
- [ ] Remove legacy variable support code paths
- [ ] Add error handling: raise if legacy vars detected at startup
- [ ] Remove all deprecation warnings
- [ ] Update all logging to reference only new variable names

**Configuration Files:**
- [ ] Remove old names from all configs
- [ ] Update all templates (docker-compose, K8s, Helm) to new names only

**Documentation:**
- [ ] Remove all references to old variable names from READMEs
- [ ] Update `MIGRATION_GUIDE.md` with Phase 2 breaking change notice
- [ ] Update `docs/configuration-reference.md` if needed
- [ ] Add CHANGELOG entry with breaking change and link to migration guide

**Tests:**
- [ ] Add/update test cases: legacy vars raise errors
- [ ] Verify error messages include helpful guidance

---

## Configuration Reference (Quick Lookup)

This section provides a quick reference for all `PRAXIS_*` variables. Detailed defaults and descriptions are in service READMEs.

### NATS Subsystem (`PRAXIS_NATS_*`)

| Variable | Default | Type | Purpose | Secret |
|----------|---------|------|---------|--------|
| `PRAXIS_NATS_URL` | `nats://localhost:4222` | string | NATS server address | No |
| `PRAXIS_NATS_STREAM` | `PRAXIS` | string | JetStream stream name | No |
| `PRAXIS_NATS_INPUT_SUBJECT` | `praxis.kernel.input` | string | Input topic for kernel | No |
| `PRAXIS_NATS_OUTPUT_SUBJECT` | `praxis.kernel.output` | string | Output topic for kernel | No |
| `PRAXIS_NATS_DURABLE` | `praxis-worker` | string | Durable consumer name | No |
| `PRAXIS_NATS_ACK_WAIT_SECONDS` | `30` | integer | Ack timeout in seconds | No |
| `PRAXIS_NATS_MAX_DELIVER` | `3` | integer | Max redelivery attempts | No |

### Storage Subsystem (`PRAXIS_STORAGE_*`)

| Variable | Default | Type | Purpose | Secret |
|----------|---------|------|---------|--------|
| `PRAXIS_STORAGE_BACKEND` | `sqlite` | string | Backend: `memory` or `sqlite` | No |
| `PRAXIS_STORAGE_SQLITE_PATH` | `build/praxis.db` | string | SQLite database file path | No |

### Telegram Integration (`PRAXIS_TELEGRAM_*`)

| Variable | Default | Type | Purpose | Secret |
|----------|---------|------|---------|--------|
| `PRAXIS_TELEGRAM_BOT_TOKEN` | (required) | string | Telegram Bot API token | **Yes** |

### Database Subsystem (`PRAXIS_DB_*`) — *Future*

| Variable | Default | Type | Purpose | Secret |
|----------|---------|------|---------|--------|
| `PRAXIS_DB_URL` | `postgres://localhost/praxis` | string | Postgres connection string | **Yes** |
| `PRAXIS_DB_USERNAME` | `praxis` | string | DB user | **Yes** |
| `PRAXIS_DB_PASSWORD` | (required) | string | DB password | **Yes** |

---

## Non-Goals

This ADR explicitly does **NOT**:

1. **Introduce a configuration library:** Each service owns its configuration loading. No shared `config.py` or constants package.
2. **Implement configuration framework:** No `ConfigManager`, `ConfigLoader`, or abstraction layer. Services use standard `os.Getenv()` (Go) or `os.getenv()` (Python).
3. **Validate configuration schema at deployment:** Validation remains in service startup code. No external schema validator.
4. **Support configuration reload/hot-swap:** Configuration is read at startup; changing env vars requires service restart.
5. **Create a centralized configuration service:** Each service reads from its environment; no runtime configuration API.
6. **Mandate configuration file format:** Services may use YAML, TOML, JSON, or plain env vars; this ADR only standardizes env var *names*, not the configuration substrate (per ADR-008).

---

## Future Extensibility

### Adding New Subsystems

When a new service or subsystem is introduced:

1. Choose a **subsystem abbreviation** (2–6 chars, no underscores): e.g., `LLM`, `REDIS`, `S3`.
2. Define variables as `PRAXIS_<SUBSYSTEM>_<SETTING>`.
3. Document in `docs/configuration-reference.md` before release.
4. Reference this ADR in the service README.

**Example:** if Redis caching is added:
```
PRAXIS_REDIS_URL = redis://localhost:6379
PRAXIS_REDIS_DB = 0
PRAXIS_REDIS_PASSWORD = (secret)
```

### Coordinating with ADR-008 (Configuration Strategy)

This ADR standardizes the **names** of environment variables (the "e" in the `env` layer of ADR-008's precedence model). ADR-008 defines the **layering precedence**: code defaults < YAML files < env vars < CLI flags. This ADR complements ADR-008 by ensuring the env var layer is consistently named across services.

### Future-Proofing Rule

**ADR-012 intentionally avoids describing implementation techniques.**

Configuration loading implementations may evolve. Future Praxis codebases may use:

- **Go:** Structs, dependency injection, singletons, or local helpers
- **Python:** Classes, dataclasses, Pydantic models, or plain dictionaries
- **Rust:** Traits, builder patterns, or procedural macros
- **Java:** Spring beans, Guice, or manual construction
- **Other languages:** TypeScript, Ruby, WASM, or unforeseen technologies

**ADR-012 governs only the naming, not the implementation.**

This is intentional. The naming convention is language-independent and timeless. The implementation approach should evolve as the platform adopts new languages and frameworks, without requiring ADR changes.

**Constraints:**
- The naming convention `PRAXIS_<SUBSYSTEM>_<SETTING>` must remain stable across all languages and implementations.
- Services must own their configuration loading; no language gets a special shared framework.
- Ownership is architectural (subsystem-based), not implementation-based (language-specific).

**Implication:** If Praxis adds a Rust service in the future, that service reads `PRAXIS_*` variables using whatever pattern fits Rust idioms—not Go patterns, not Python patterns. The name is the contract; the loading is free.

---

## When This ADR Must Be Updated

ADR-012 must be reviewed and potentially amended whenever:

1. **A new subsystem is introduced:** A new bounded context defined by RFC or ADR requires a new namespace (e.g., `PRAXIS_LLM_*`, `PRAXIS_REDIS_*`). Update the Subsystems table and ownership documentation.

2. **A runtime configuration namespace is proposed that does not fit existing subsystems:** Architecture must decide: create a new subsystem or consolidate with an existing one. This ADR must document the decision.

3. **A new deployment environment requires additional configuration:** (e.g., separate staging/production secrets, multi-tenant configuration, canary deployments). This ADR should address whether new variables are needed or whether existing precedence model (ADR-008) suffices.

4. **Configuration ownership changes:** If an RFC or subsystem definition changes ownership of a variable (rare), update the Ownership Model section and subsystem table.

5. **The naming convention itself is challenged:** If a project-wide proposal suggests deviating from `PRAXIS_<SUBSYSTEM>_*`, that requires ADR amendment or a new ADR superseding this one.

### When ADR-012 Does NOT Need to Change

- **Adding a new environment variable that follows the convention:** Just document it in the service README and `docs/configuration-reference.md`. No ADR change needed.
- **Changing the default value of an existing variable:** Document in release notes and service README. No ADR change needed (unless the meaning changes; see Compatibility Rule).
- **Deprecating a variable:** Use Phase 1 (both names work) then Phase 2 (old name errors). Update `MIGRATION_GUIDE.md`. No ADR change needed unless the pattern changes.
- **Adding a new service or application:** If it uses only existing subsystem variables, no ADR change needed.

---

## Acceptance Criteria

All of the following MUST be true before this ADR transitions from PROPOSED to ACCEPTED:

### Architectural Alignment

- [ ] All environment variable names follow the pattern `PRAXIS_<SUBSYSTEM>_<SETTING>` with no exceptions.
- [ ] Every variable has exactly one documented architectural owner (from an RFC or subsystem definition).
- [ ] Configuration ownership is separated from loading implementation (configuration is architectural, loading is application-specific).

### Non-Goal Compliance

- [ ] No shared configuration framework, library, or manager is introduced.
- [ ] No reflection-based configuration loading or dependency injection container is created.
- [ ] No global configuration singleton or shared state is introduced.
- [ ] No centralized configuration service or API is created.
- [ ] Services own their configuration loading; no shared implementations across services.

### Invariant Enforcement

- [ ] **Invariant 1 (Ownership Separation):** Configuration belongs to applications; naming belongs to architecture; loading belongs to applications. All three are addressed separately in implementation.
- [ ] **Invariant 2 (No Coupling from Naming):** The naming convention does not imply or require shared code. Services may follow identical naming with completely independent loading logic.
- [ ] **Invariant 3 (Service Independence):** No runtime configuration creates dependencies between services. Verification: `grep -r "import.*config" services/ packages/ apps/ | grep -v self` returns no cross-service configuration imports.

### Documentation Completeness

- [ ] This ADR documents the full naming convention with examples.
- [ ] All current Praxis configuration variables are listed in `docs/configuration-reference.md` organized by subsystem.
- [ ] `MIGRATION_GUIDE.md` exists and includes old→new variable mapping with deployment examples (docker-compose, Kubernetes).
- [ ] Each subsystem section in `docs/configuration-reference.md` identifies the owner and purpose of each variable.

### Phase 1 Implementation Ready

- [ ] Go configuration loading code supports both `NATS_*` (deprecated) and `PRAXIS_NATS_*` (new) with deprecation warnings.
- [ ] Python services support both old and new names with deprecation warnings.
- [ ] Deprecation warnings include link to `MIGRATION_GUIDE.md`.
- [ ] Code review checklist item added: "Environment variable names follow `PRAXIS_<SUBSYSTEM>_*` convention."
- [ ] Phase 1 code is tested with both old and new variables; new variables take precedence.

### Deployment and Operations

- [ ] `docker-compose.yml` uses new variable names in all examples.
- [ ] Kubernetes manifests (if maintained) use new variable names.
- [ ] All service READMEs include env var tables using new names.
- [ ] Build pipeline checks for use of old variable names and warns during CI.

---

## Related RFCs and ADRs

- **RFC-030:** Cross-Cutting Concerns (defines subsystem concepts and ownership model)
- **ADR-008:** Configuration Strategy (defines precedence: code defaults < YAML < env vars < CLI flags)
- **ADR-007:** Storage Architecture (defines storage subsystem and PRAXIS_STORAGE_* variables)
- **RFC-040:** Agent Architecture (defines agent subsystem boundaries)

---

## Questions and Answers

**Q: Why not use environment variable module/registry?**
A: Extract, don't invent. Only one Python service exists; duplication doesn't justify shared code. When multiple services emerge, extraction is easy. Introducing a registry now adds coupling without proven benefit.

**Q: Can I add new variables without following this convention?**
A: No. All environment variables must use `PRAXIS_<SUBSYSTEM>_<SETTING>` from day 1. This is an architectural law (Invariant 1).

**Q: What about third-party services we're integrating (e.g., Upwork API)?**
A: Use `PRAXIS_UPWORK_*` if the integration is part of Praxis subsystems. If you're building a standalone connector, you may use different naming, but the integration's configuration within Praxis services must follow `PRAXIS_*`.

**Q: Can multiple services share configuration loading logic?**
A: Only if they are explicitly coupled by an RFC or ADR that defines them as a subsystem unit. Standard Praxis services remain independent and own their configuration loading. Shared code is allowed AFTER proven duplication; never before.

**Q: Who defines the subsystem for a new feature?**
A: The RFC or ADR that introduces the feature. If you're adding a new subsystem, update the relevant RFC first, then document the variables in this ADR or a follow-up ADR.

**Q: When should we move from Phase 1 to Phase 2?**
A: After at least 1–2 release cycles in Phase 1 and all operational deployments have migrated to new names. Timeline is determined by the project's release cadence and adoption metrics (deprecation warning logs).

---

## Final Governance Statement

**ADR-012 standardizes the language of runtime configuration, not its implementation.**

This distinction is fundamental:

| Governed by ADR-012 | Owned by Applications |
|---------|---------|
| Names: `PRAXIS_<SUBSYSTEM>_<SETTING>` | Loading: how each service reads its config |
| Ownership: which subsystem interprets a variable | Implementation: Go structs, Python dicts, Rust traits, etc. |
| Meaning: what the variable controls | Precedence: which sources to read first (ADR-008) |
| Compatibility: deprecation and migration rules | Format: YAML, TOML, JSON, env vars, or hybrid |

**Architecture owns names.**
Names are platform-wide, language-independent, and stable. The naming convention is decided here and governs all services.

**Applications own loading.**
Each service decides how to read its configuration, apply precedence, validate inputs, and handle errors. No shared framework, no centralized manager, no dependency injection. Each service is autonomous.

**Subsystems own interpretation.**
Each subsystem (defined by RFC or ADR) owns what its variables mean. The Transport/NATS subsystem interprets `PRAXIS_NATS_URL`. The Storage subsystem interprets `PRAXIS_STORAGE_BACKEND`. Ownership is architectural, not organizational or implementation-based.

**Implementations remain local.**
Go services use Go patterns. Python services use Python patterns. Future languages use their idioms. The naming is the contract; the code is free.

---

## References

- **ADR-008:** Configuration Strategy (layering and precedence)
- **ADR-007:** Storage Architecture (storage configuration)
- **RFC-030 §13:** Cross-Cutting Concerns (configuration and secrets)
- **RFC-033 §6:** Storage Model (Configuration Store)
- **docs/adr/** — all other architectural decisions

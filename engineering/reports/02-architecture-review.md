# 02 - Architecture Review

Date: 2026-07-05

## Stage 1 - Existing Knowledge Search

### Existing Components

- Telegram ingress adapter: `apps/telegram/main.py`
- Worker transport subscriber: `internal/transport/natsworker/subscriber.go`
- NATS transport config/client/publisher: `internal/transport/nats/*.go`
- Worker composition root: `services/worker/main.go`

### Reference Implementations

- Golden mapper: `apps/telegram/main.py::telegram_update_to_payload()`
- Health endpoint pattern: `services/api/routes/health.py`

### Golden Examples

- Mapper invariants: `docs/architecture/principles/GOLDEN_MAPPER.md`

### Similar Packages

- NATS adapter package split: `internal/transport/nats` + `internal/transport/natsworker`

### Existing Patterns

- Transport Mapper pattern: `docs/architecture/reference/PATTERN_CATALOG.md`
- Reference registry/search rule: `docs/architecture/reference/REFERENCE_REGISTRY.md`

## Stage 2 - Architecture Review

### Relevant RFCs

- RFC-013: canonical event envelope, correlation/causation/trace, immutable/replay-safe events.
- RFC-030: layer boundaries, edge-sync/internal-async model, observability service responsibility.
- RFC-031: command/event/query contracts, idempotency and observability contracts.
- RFC-032: pipeline data flow, dead-letter handling, correlation model.
- RFC-060: retry/health/observability test obligations.

### Relevant ADRs

- ADR-003 Internal Event Bus (NATS + JetStream direction; currently PROPOSED)
- ADR-006 Transport Boundaries (Kernel transport-agnostic; adapters own translation)
- ADR-009 Observability (health/metrics/traces expectations; currently PROPOSED)
- ADR-011 Build and Development Workflow (Taskfile as canonical verification surface)

### Architecture Summary

Sprint 1 work is infrastructure/adapter hardening around an existing pipeline; Kernel domain logic is not the primary change target. Transport boundaries must remain strict: mapper stays pure, subscriber handles delivery concerns, and composition roots wire dependencies.

### Affected Layers

- Experience/Adapter: Telegram adapter ingress and reply behavior.
- Infrastructure: NATS delivery, retry, DLQ, reconnect, metrics plumbing.
- Composition Root: worker/telegram startup wiring.

### Ownership

- Mapper owns deterministic representation translation only.
- Subscriber/runtime loop owns retry/ack/nak/term flow and publish outcomes.
- Composition root owns config and dependency wiring.

### Invariants

1. Event envelope remains immutable and transport-independent.
2. Correlation and causation remain distinct and preserved end-to-end.
3. No business logic migrates into mapper/transport adapters.
4. Retry and DLQ logic do not violate idempotency requirements.

## Stage 3 - Architecture Diff (pre-code)

### Packages

NEW
- None planned.

MODIFIED
- `apps/telegram` (reply/reconnect/health/metrics concerns)
- `internal/transport/natsworker` (retry/DLQ behavior and metadata preservation)
- `internal/transport/nats` (config surface for delivery tuning)
- `services/worker` (composition/config wiring)

REMOVED
- None.

### Dependencies

Before
- Telegram adapter publishes to NATS subject.
- Worker consumes input subject and publishes output.
- Retry/DLQ/reconnect/health/metrics are partial/incomplete.

After (target)
- Same dependency direction.
- Additional resilience and observability behavior at adapter/infrastructure layer.

### Interfaces

- Reuse existing `InputMessage` / `OutputMessage` contracts.
- Avoid new exported interfaces unless strictly required by duplication.

### Repositories / Storage

- No new storage backend introduced.
- DLQ should remain transport-level routing concern first, not new persistence abstraction.

### Adapters / Services / Composition Root

- Adapter behavior extended, not re-architected.
- Composition roots updated only to wire new config knobs.

## Stage 4 - RFC Impact Matrix

### RFC-013 Event Model

Affected: YES

Reason
- Correlation/causation/trace propagation and event envelope integrity.

Changes
- Ensure ingress and downstream output preserve required metadata.

Risk
- Medium if metadata is dropped on retry/reconnect paths.

### RFC-030 System Architecture

Affected: YES

Reason
- Transport boundary and observability responsibilities.

Changes
- Keep transport resilience in infra/adapter layer only.

Risk
- Medium if business logic leaks into adapter.

### RFC-031 Service Contracts

Affected: YES

Reason
- Idempotency/retry and observability contract obligations.

Changes
- Validate delivery/retry behavior against command/event rules.

Risk
- Medium if retries duplicate side effects.

### RFC-032 Data Flow

Affected: YES

Reason
- DLQ and correlation model are explicit in this RFC.

Changes
- Introduce/verify dead-letter handling path for failed processing.

Risk
- Medium if DLQ routing semantics are underspecified.

### RFC-060 Testing Strategy

Affected: YES

Reason
- Retry/provider health/observability behavior must be tested.

Changes
- Add tests for retry storms, health handling, metrics emission.

Risk
- High if sprint merges without these tests.

Conflict Check
- No direct RFC conflicts identified in reviewed scope.

## Stage 5 - Reference Comparison

Reference Used
- `apps/telegram/main.py::telegram_update_to_payload()` as golden mapper.

Required Comparison Result
- Any resilience logic (retry, reconnect, metrics, health) must not be placed inside mapper function.

Planned Deviations
- None in mapper logic.
- Resilience features belong in runtime loop and composition/config.

Deviation Justification
- N/A (no justified mapper deviation required).

## Stage 6 - Complexity Budget

### Complexity Budget

+ Files: <= 10 modified in first slice.
+ Packages: 3-4 touched (telegram adapter, natsworker, worker main, tests).
+ Interfaces: 0 new unless duplication appears.
+ DTOs: 0 new canonical DTOs.
+ Repositories: 0 new repositories.
+ Adapters: 0 new adapter categories; extend existing.
+ Services: 0 new services.
+ Public APIs: 0 new public API surfaces in first slice.

Budget Rule
- If any slice exceeds this budget without measurable reliability gain, stop and redesign.

## Architecture Gates

- Gate 1 (Architecture Approved): PASS for minimal-slice implementation strategy.
- Gate 2 (Reference Comparison Passed): PASS (golden mapper constraints explicit).
- Gate 3 (Complexity Budget Approved): PASS.

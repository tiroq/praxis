# 01 - Intake

Date: 2026-07-05
Task Source: `engineering/tasks/sprint-1.md`
Normalized Task: `engineering/tasks/current-task.md`

## Stage 0 - Task Intake

```yaml
Task: Sprint 1 - Production-ready Capture Pipeline
Goal: Deliver reliable and observable Telegram->NATS->Kernel capture path with traceable metadata.
Constraints:
  - Follow all RFCs, ADRs, and policy/instruction documents.
  - Search and compare reference implementations first.
  - No speculative abstractions.
  - Use minimal vertical slices.
Unknowns:
  - Exact acceptance thresholds for retry policy (attempt counts, backoff caps).
  - DLQ ownership and retention policy details.
  - Health endpoint contract scope for Telegram adapter process.
  - Metrics backend/export format required in this sprint.
Risks:
  - Scope coupling across Python Telegram adapter and Go worker transport.
  - Introducing retry logic in wrong layer (mapper/adapter boundary violations).
  - Reconnect behavior changes causing duplicate deliveries without idempotency checks.
```

## Assumptions

1. Sprint 1 is executed as multiple small slices, not one large merge.
2. Existing NATS + JetStream stack remains the transport baseline.
3. Baseline verification commands remain canonical: `go test ./...`, `golangci-lint run`, `task dev`.

## Expected Outcome

Capture pipeline supports:

1. End-to-end correlation metadata propagation.
2. Reliable processing under transient transport failures.
3. Explicit failure isolation (DLQ) for non-recoverable deliveries.
4. Observable health and metrics for runtime operation.

## Affected Packages / Services (initial)

- `apps/telegram`
- `internal/transport/nats`
- `internal/transport/natsworker`
- `services/worker`
- `services/api` (reference pattern for health surface)

## Affected RFCs (initial)

- RFC-013 Event Model
- RFC-030 System Architecture
- RFC-031 Service Contracts
- RFC-032 Data Flow
- RFC-060 Testing Strategy

## Gate 0

Status: PASS (with bounded assumptions recorded)

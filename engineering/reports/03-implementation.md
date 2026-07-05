# 03 - Implementation

Date: 2026-07-05

## Stage 7 - Implementation Plan

Implementation strategy is a sequence of minimal vertical slices.

### Slice S1 - Correlation ID Completion (first)

Owner
- Transport/worker maintainers.

Responsibility
- Ensure inbound Telegram payloads and worker outputs preserve correlation metadata required by RFC-013/RFC-032.

Justification
- Existing sprint goal begins with correlation IDs; this is prerequisite observability context for retry/DLQ/metrics.

Scope
1. Confirm ingress mapping and downstream output correlation handling.
2. Add focused tests for missing/propagated correlation IDs.
3. Keep mapper pure; no resilience logic in mapper.

### Slice S2 - Retry Strategy at Delivery Boundary

Scope
1. Implement bounded retry policy in subscriber/runtime boundary (not mapper).
2. Confirm idempotent behavior and non-duplication expectations.
3. Add tests covering retry paths and terminal failure handoff.

### Slice S3 - DLQ Handling

Scope
1. Route terminally failed messages to DLQ subject/flow using existing transport capabilities.
2. Add tests for Term/Nak behavior and DLQ routing semantics.

### Slice S4 - Health Endpoints

Scope
1. Add health surface for telegram adapter runtime (polling + NATS connectivity state).
2. Align shape with existing health pattern.

### Slice S5 - Metrics

Scope
1. Emit counters for received/published/failed/retried/DLQ deliveries.
2. Ensure metadata allows correlation-aware tracing.

### Slice S6 - Reconnect Logic Improvements

Scope
1. Add bounded reconnect strategy and backoff.
2. Add tests for disconnect/reconnect scenarios and stale client handling.

### Slice S7 - Telegram Reply Pipeline

Scope
1. Implement output->Telegram reply path as adapter concern.
2. Keep decision logic outside adapter.

## Stage 8 - Implementation

Status: EXECUTED.

Implemented slices:

1. S1 Correlation ID completion
	- Added `correlation_id` to Telegram ingress payload mapping.
	- Propagated correlation metadata into worker output wire contract.
	- Added correlation fallback (`event ID`) when inbound correlation is missing.

2. S2 Retry strategy (delivery boundary)
	- Added bounded publish retry loop in Telegram adapter for input subject publishing.
	- Existing worker redelivery strategy remains bounded via `MaxDeliver` and `Nak` flow.

3. S3 DLQ
	- Added DLQ subject to NATS transport config (`NATS_DLQ_SUBJECT`, default `praxis.kernel.dlq`).
	- Added worker DLQ publish path when output publish continues failing at terminal delivery attempt.
	- Added Telegram adapter DLQ fallback for unrecoverable ingress publish failure.

4. S4 Health endpoints
	- Added HTTP `/health` endpoint in Telegram adapter runtime.

5. S5 Metrics
	- Added in-process counters and `/metrics` endpoint in Prometheus text format for Telegram adapter runtime.

6. S6 Reconnect logic improvements
	- Added explicit NATS reconnect callbacks and reconnect tuning options in Telegram adapter.

7. S7 Telegram reply pipeline
	- Added output-subject subscription in Telegram adapter.
	- Added output-to-user reply delivery using `chat_id` from output metadata with local fallback map.

Files modified for implementation:

- `apps/telegram/main.py`
- `internal/transport/nats/config.go`
- `internal/transport/nats/messages.go`
- `internal/transport/natsworker/subscriber.go`
- `internal/transport/nats/nats_test.go`
- `internal/transport/natsworker/subscriber_test.go`

## Gate 4

Gate 4 (Implementation Complete): PASS.

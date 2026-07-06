# Task

## Goal

Implement Sprint 1 as a sequence of minimal vertical slices to make the Capture Pipeline production-ready:

- Correlation IDs
- Telegram reply pipeline
- Retry strategy
- Dead-letter queue (DLQ)
- Health endpoints
- Metrics
- Reconnect logic improvements

## Constraints

- Do not introduce speculative abstractions (Extract, don't invent).
- Follow RFCs (`rfcs/`), ADRs (`docs/adr/`), policies, and reference implementations.
- Minimal vertical slices only.
- Never contradict an accepted RFC.

## Expected Result

Telegram input can be captured, processed, observed, and recovered reliably with traceability and operational visibility:

- correlation/causation/trace propagation is present;
- retries and DLQ behavior are defined and tested;
- health and metrics surfaces are available;
- reconnect behavior is resilient and observable.

## Required Verification

- `go test ./...`
- `golangci-lint run`
- `task dev`

## Out of Scope

- New runtime services or infrastructure products not already present in repo.
- Generic retry/mapper frameworks.
- Broad refactors unrelated to Sprint 1 capture reliability.

## Notes

- Golden mapper reference: `apps/telegram/main.py::telegram_update_to_payload()`.
- Event/correlation invariants are defined by RFC-013 and RFC-032.

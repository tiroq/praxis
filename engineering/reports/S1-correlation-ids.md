# S1 - Correlation IDs

Date: 2026-07-05

## Scope

Fix the worker correlation-ID path so the normalized correlation ID is used consistently from InputMessage to Kernel Event to OutputMessage.

## Change Summary

- Normalized correlation ID now comes from the worker's computed Kernel Event and is reused when building OutputMessage.
- Added regression coverage for an InputMessage missing correlation_id.
- Verified the fallback path preserves the input event ID as the effective correlation ID.

## Verification

Executed after the fix:

- `go test ./...` - PASS
- `golangci-lint run` - PASS
- `task dev` - PASS
- `task smoke:nats` - PASS
- `task smoke:praxis` - PASS

## Result

PASS for S1 only.

## Notes

- S2-S7 were not modified in this step.
- The fix is limited to `internal/transport/natsworker/subscriber.go` and its regression test.

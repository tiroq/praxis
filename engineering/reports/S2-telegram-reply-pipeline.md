# S2 - Telegram Reply Pipeline

Date: 2026-07-05

## Scope

Implement the Telegram reply pipeline so worker output messages can be observed and relayed back to the originating Telegram chat.

## Change Summary

- Extracted the reply-handling logic into a testable helper: [handle_output_message](apps/telegram/main.py#L302).
- Kept the adapter boundary intact: Telegram message mapping remains in [telegram_update_to_payload](apps/telegram/main.py#L249), while reply delivery stays in the runtime adapter.
- Added reply pipeline regression coverage in [apps/telegram/test_main.py](apps/telegram/test_main.py).

## Behavior Verified

- Worker output with `metadata.chat_id` sends a reply to that chat.
- Worker output without metadata falls back to the pending chat-id map.
- Reply text is rendered from pipeline status/result payloads.
- Reply messages are removed from the pending map after successful delivery.

## Verification

Executed after the S2-only change:

- `python -m unittest discover -s apps/telegram -p 'test_*.py'` - PASS
- `go test ./...` - PASS
- `golangci-lint run` - PASS
- `task dev` - PASS
- `task smoke:nats` - PASS
- `task smoke:praxis` - PASS

## Result

PASS for S2 only.

## Notes

- S1 correlation-id behavior was left unchanged.
- S3-S7 were not modified in this step.
- The reply pipeline remains an adapter/runtime concern, not a mapper concern.

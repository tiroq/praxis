# Telegram Adapter Runtime Failure Audit

**Date**: 2026-07-04  
**Subject**: `apps/telegram/main.py`  
**Scope**: Failure classification and phased remediation strategy

---

## Corrected Failure Classification

### Startup Failures (Fail-Fast)

**TELEGRAM_BOT_TOKEN missing**
- Behavior: Logs error message, exits via `sys.exit(1)`
- Classification: **Designed fail-fast (correct)**
- Recovery: None (intentional); requires human remediation

**NATS unavailable at startup**
- Behavior: Logs error with URL and exception, exits via `sys.exit(1)`  
- Classification: **Designed fail-fast (correct)**
- Recovery: None (intentional); requires human restart once NATS is available

---

### Runtime Failures (Per-Message)

**NATS publish fails**
- Behavior: Exception caught at `js.publish()` call; logged with message ID, subject, and exception detail
- Classification: **Visible failure without recovery**
  - ✓ Failure is logged (visible)
  - ✗ No retry, buffer, or recovery mechanism
  - Message is discarded after logging
- Current: Polling continues; next message handled independently
- Risk: Message loss with full auditability (ID + exception logged)

**NATS disconnects during runtime**
- Behavior: Same as publish failure—next `js.publish()` attempt raises exception (stale object reference)
- Classification: **Visible failure without recovery** (cascades to publish failures)
  - Subsequent messages trigger the publish exception handler
  - No reconnection or circuit-breaker logic
- Current: Polling continues; adapter becomes non-functional
- Risk: Silent functional death after initial log entry

---

### Inbound Channel Failures

**Telegram polling hangs / API disconnects**
- Behavior: `app.updater.start_polling()` becomes unresponsive; no explicit error handling in `run()`
- Classification: **Unobserved failure**
  - No exception is caught or logged
  - Process remains running but polling thread/coroutine is stuck
  - New Telegram messages are not polled; zero indication in logs
- Current: Adapter appears running but silent inbound failure
- Risk: Complete inbound channel loss with no observability

---

## Phase 1 Acceptance Decision

**Decision: ACCEPT Phase 1 as-is (no code changes required)**

Phase 1 requirements already met:
1. ✓ Fail-fast on missing `TELEGRAM_BOT_TOKEN`
2. ✓ Fail-fast when NATS unavailable at startup
3. ✓ Log publish failures with message ID and exception
4. ✓ Keep polling after per-message failures

**Current implementation achieves minimal correct behavior** within constraints:
- Configuration validation is explicit
- Startup safety is ensured
- Per-message failures are logged (not silent)
- Polling resilience allows transient outages to recover autonomously

**No implementation changes needed for Phase 1.**

---

## Phase 2 Backlog

### Candidates (in priority order)

1. **NATS reconnect strategy**
   - Implements exponential backoff and reconnection attempt loop
   - Detects stale `js` object after disconnect
   - Requirement: Conditional—only if NATS availability SLA requires recovery

2. **Telegram polling health detection**
   - Wraps `app.updater.start_polling()` with watchdog or exception handler
   - Logs/exits if polling becomes unresponsive
   - Requirement: High—prevents silent inbound channel loss

3. **Dead-letter handling**
   - Routes failed publishes to a separate subject or local log
   - Enables post-hoc recovery or manual intervention
   - Constraint: No persistent buffering allowed (out of scope)

4. **Metrics & observability**
   - Counter: messages polled, published, failed (with classification)
   - Gauge: NATS connection state, polling state
   - Requirement: Informational—enables monitoring

5. **Health endpoint**
   - HTTP endpoint reporting NATS connectivity, polling status, recent failure count
   - Requirement: Integration—needed if Praxis has a control plane

6. **Output-to-Telegram error replies**
   - On publish failure, sends error message back to Telegram user
   - Provides user feedback loop
   - Constraint: Praxis adapter itself must not have reply logic; separate service needed

---

## Summary

| Aspect | Status |
|--------|--------|
| Phase 1 terminology corrected | ✓ |
| Phase 1 acceptance decision | Accept as-is |
| Implementation changes required now | No |
| Code modifications needed | None |
| High-priority Phase 2 | Telegram polling health + NATS reconnect |


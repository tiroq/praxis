# Praxis RFC Hygiene Report

- Generated: 2026-06-29T15:11:46Z
- Status: **PASS** (0 error(s), 0 warning(s))
- RFC files scanned: 32

> Read-only Phase 0 tooling. It does not modify RFC files and does not resolve
> any architecture decision. See [docs/ARCHITECTURE_DECISION_QUEUE.md](../../../docs/ARCHITECTURE_DECISION_QUEUE.md).

## Check Summary

| Check | Errors | Warnings | Status |
| --- | --- | --- | --- |
| duplicate_bodies | 0 | 0 | PASS |
| titles | 0 | 0 | PASS |
| references | 0 | 0 | PASS |

## Structural Errors (cause non-zero exit)

_None._

## Warnings (non-blocking)

_None._

## Architecture Decision Advisories (open — not resolved)

These ambiguities are tracked in [docs/ARCHITECTURE_DECISION_QUEUE.md](../../../docs/ARCHITECTURE_DECISION_QUEUE.md) and are reported as advisories only. They never fail the build.

| ID | Status | Summary |
| --- | --- | --- |
| ADQ-001 | open | Canonical Object vs Artifact model is unresolved. |
| ADQ-002 | open | Memory / Knowledge service ownership is undefined in RFC-030/031. |
| ADQ-003 | open | Workflow model has no RFC; deferred for Phases 1-2. |
| ADQ-004 | open | Action state machine differs between RFC-022 and RFC-023. |
| ADQ-005 | open | Invariant ID registry referenced by RFC-061 is undefined in RFC-060. |
| ADQ-006 | open | Some Space RFC storage mappings contradict RFC-033. |

## Artifacts

- `verify/rfc/out/rfc_graph.json` — 32 nodes, 196 edges
- `verify/rfc/out/invariants.json` — 536 draft items (non-authoritative; see ADQ-005)


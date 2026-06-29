---
description: "Use when implementing, modifying, or reviewing Praxis code in services/, packages/, apps/, scripts/, or infra/. Enforces RFC-driven development, the canonical implementation order, and non-negotiable architectural invariants (immutable events, auditable decisions, policy-bound memory, bounded spaces)."
name: "Praxis Implementation Discipline"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Implementation Discipline

RFCs in `./rfcs` are the **architectural source of truth**. Implementation follows
architecture, never the reverse. Never write code that contradicts an accepted RFC.
If code and an RFC disagree, the RFC wins — stop and flag the conflict.

## Before Changing Code

Do this every time, before editing:

1. **Identify relevant RFCs.** Name the specific RFC(s) the change traces back to.
2. **List affected invariants.** State which hard rules (below) the change touches.
3. **Propose a minimal implementation slice.** Smallest change that satisfies the RFC.
4. **Add or update verification tests.** No behavior change ships without a test that
   pins the invariant or contract it relies on.

If you cannot map a change to an RFC, stop and ask — do not invent behavior.

## Implementation Order

Build foundations before dependents. Respect this sequence:

1. RFC-031 — service envelopes (`rfcs/031-service-contracts.md`)
2. RFC-013 — event model (`rfcs/013-event-model.md`)
3. RFC-014 — identity / representation model (`rfcs/014-identity-representation-model.md`)
4. RFC-015 — object lifecycle model (`rfcs/015-object-lifecycle-model.md`)
5. RFC-020 — review system (`rfcs/020-review-system.md`)
6. RFC-021 — decision model (`rfcs/021-decision-model.md`)
7. RFC-023 — action model (`rfcs/023-action-model.md`)

Do not implement a later layer before its prerequisites exist and are tested.

## Hard Rules (Invariants)

These are non-negotiable. Any code that violates one is wrong by definition.

- **Events are immutable.** Never mutate or delete an emitted event; correct via new events.
- **Decisions are explicit and auditable.** Every Decision records who/what, why, and inputs.
- **Reviews never commit Decisions.** A Review produces findings; it cannot enact a Decision.
- **Agents never mutate canonical state directly.** They propose; the system commits.
- **Agents never call LLM providers directly.** All model access goes through the LLM router.
- **Prompt versions are immutable after release.** Released prompts are frozen; ship a new version.
- **Memory is policy-bound.** Reads/writes to memory honor the governing policy; no ad-hoc access.
- **Spaces are bounded contexts.** Keep models, data, and logic within their space.
- **Cross-space communication is explicit.** Use defined contracts/events; no hidden coupling.
- **Derived stores are rebuildable.** Never treat a projection/cache as a source of truth.

## When Unsure

Prefer reading the RFC over guessing. If an RFC is ambiguous or missing, surface the gap
rather than filling it with assumed behavior.

---
description: "Use when implementing, modifying, or reviewing Praxis code in services/, packages/, apps/, scripts/, or infra/. Enforces plan-first RFC-driven development, reuse-before-build, vertical slices, the canonical implementation order, and non-negotiable architectural invariants (immutable events, auditable decisions, policy-bound memory, bounded spaces)."
name: "Praxis Implementation Discipline"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Implementation Discipline

RFCs in `./rfcs` are the **architectural source of truth**. Implementation follows
architecture, never the reverse. Never write code that contradicts an accepted RFC.
If code and an RFC disagree, the RFC wins — stop and flag the conflict.

## 1. Plan Before Coding — Never Start Immediately

Do not write code first. Produce an **Implementation Plan** and get approval (or an
explicit "just do it") before editing. The plan must include:

- **Relevant RFCs** — the specific RFC(s) the work traces back to.
- **Impacted invariants** — which hard rules (below) the change touches.
- **Existing components to reuse** — services, events, commands, queries, aggregates,
  projections, agents, storage already defined.
- **New components to create** — only what reuse cannot cover, each justified by an RFC.
- **Minimal implementation slice** — the smallest end-to-end change that delivers value.
- **Verification strategy** — the tests / verification scripts that will pin behavior.
- **Risks** — what could break, ambiguous RFC areas, and second-source-of-truth dangers.

If you cannot map the work to an RFC, stop and ask — do not invent behavior.

## 2. Prefer Extending Existing Architecture

Never introduce a new abstraction if an RFC or the codebase already defines one. Before
creating anything, search for an existing: **service, event, command, query, aggregate,
projection, agent, storage**. Extend or compose what exists; only create new when reuse
is genuinely impossible, and say why.

## 3. Think in Vertical Slices

Never build entire subsystems at once. Implement the smallest end-to-end slice that
provides value, typically:

```
Command → Event → Projection → Query → Verification
```

Ship one working path through the layers instead of many half-built services.

## 4. RFC Compliance Review (before implementing)

Before writing the slice, explicitly answer each. If any answer is **YES**, stop and explain:

- Does this violate any RFC?
- Does this duplicate an existing concept?
- Does this introduce a second source of truth?
- Does this bypass Review → Decision → Action?
- Does this introduce mutable canonical state?
- Does this weaken auditability?

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

## Avoid Speculative Implementation

Do not implement future features, generic frameworks, unused abstractions, placeholder
services, or speculative configuration. Everything must be justified by an existing RFC
or an accepted implementation slice.

## Minimize Architectural Debt

- Prefer deleting code over adding abstractions.
- Prefer composition over inheritance.
- Prefer explicit behavior over magic.
- Prefer deterministic behavior over convenience.

## Challenge Assumptions

Do not blindly implement the requested design. If a simpler solution satisfies the RFCs,
explain it and recommend it before implementing. If the request contradicts long-term
architecture, explain why instead of complying silently.

## Leave the Repository Healthier

Every change should improve at least one of: documentation, verification, naming,
comments, tests, architecture consistency, or dead-code removal.

## After Implementation — Architecture Review

When a slice is complete, produce a short review:

- **RFC Compliance** — every RFC the change satisfies.
- **Invariants** — every invariant preserved.
- **Future RFC impact** — RFCs that should be updated because implementation revealed
  missing or ambiguous detail.

## After Implementation — Summary

End every implementation with:

- Files changed
- RFCs implemented
- New Commands
- New Events
- New Queries
- New Services
- Tests added
- Verification scripts added
- Remaining work
- Recommended next implementation slice

## When Unsure

Prefer reading the RFC over guessing. If an RFC is ambiguous or missing, surface the gap
rather than filling it with assumed behavior.

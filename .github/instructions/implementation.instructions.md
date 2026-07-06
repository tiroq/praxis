---
name: "Implementation Discipline — Reference"
description: "Reference documentation. For active Copilot instructions, see .github/copilot-instructions.md"
---

# Implementation Discipline

RFCs in `./rfcs` are **architectural source of truth**. Implementation follows architecture, never reverse.

---

## Plan Before Coding

Do not write code first. Produce **Implementation Plan** before editing.

Plan must include:

- **Relevant RFCs** — specific RFC(s) this traces back to
- **Impacted invariants** — which hard rules this change touches (see engineering-laws.instructions.md)
- **Existing components to reuse** — services, events, commands, queries, aggregates, projections, adapters, storage already defined
- **New components to create** — only what reuse cannot cover; each justified by RFC or architecture review
- **Minimal implementation slice** — smallest end-to-end change delivering value
- **Verification strategy** — tests/commands that will validate behavior
- **Risks** — what could break, ambiguous RFC areas, second-source-of-truth dangers

**Cannot map work to RFC?** Stop and ask. Do not invent behavior.

---

## Reuse Before Build

Never introduce new abstraction if RFC or codebase already defines one.

Before creating anything, search for existing:

- service
- event
- command
- query
- aggregate
- projection
- agent
- storage
- adapter

Extend or compose what exists. Only create new when reuse genuinely impossible.

Categories to search for reference implementations:

- mapper
- repository
- worker
- adapter
- storage backend
- transport
- bootstrap
- configuration loader
- service
- reducer
- projection
- event
- command
- query
- aggregate

**If one exists → prefer evolving it** over building new.

---

## Think in Vertical Slices

Never build entire subsystems at once. Implement smallest end-to-end slice providing value:

```
Command → Event → Projection → Query → Verification
```

Ship one working path through layers instead of many half-built services.

---

## Implementation Order

Build foundations before dependents. Respect this sequence:

1. RFC-031 — service envelopes (`rfcs/031-service-contracts.md`)
2. RFC-013 — event model (`rfcs/013-event-model.md`)
3. RFC-014 — identity/representation model (`rfcs/014-identity-representation-model.md`)
4. RFC-015 — object lifecycle model (`rfcs/015-object-lifecycle-model.md`)
5. RFC-020 — review system (`rfcs/020-review-system.md`)
6. RFC-021 — decision model (`rfcs/021-decision-model.md`)
7. RFC-023 — action model (`rfcs/023-action-model.md`)

Do not implement later layer before prerequisites exist and are tested.

---

## Implementation Rules

- Reuse existing code
- Reuse existing packages
- Reuse existing architecture
- Prefer deletion over addition
- Prefer concrete implementations
- Avoid speculative abstractions

---

## RFC Compliance Review (Before Implementing)

Before writing slice, explicitly answer each. If any is **YES**, stop and explain:

- Does this violate any RFC?
- Does this duplicate existing concept?
- Does this introduce second source of truth?
- Does this bypass Review → Decision → Action?
- Does this introduce mutable canonical state?
- Does this weaken auditability?

---

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


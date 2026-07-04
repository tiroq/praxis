---
description: "Use when implementing, modifying, or reviewing Praxis code in services/, packages/, apps/, scripts/, or infra/. Enforces plan-first RFC-driven development, reuse-before-build, vertical slices, the canonical implementation order, and non-negotiable architectural invariants (immutable events, auditable decisions, policy-bound memory, bounded spaces)."
name: "Praxis Implementation Discipline"
applyTo: ["services/**", "packages/**", "apps/**", "scripts/**", "infra/**"]
---

# Praxis Implementation Discipline

RFCs in `./rfcs` are **architectural source of truth**. Implementation follows architecture, never reverse. Never write code contradicting accepted RFC. If code and RFC disagree, RFC wins — stop and flag conflict.

## 1. Architecture Review Gate — Stop Before New Abstractions

**Before implementing any:**

- new package
- new repository
- new engine
- new storage mechanism
- new runtime
- new adapter
- new projection
- new public interface
- new composition root

**You MUST:**

1. Read and follow `.github/instructions/praxis-architecture-guardian.instructions.md`
2. Complete full 10-phase architecture review
3. Produce all deliverables (RFC summary, ownership analysis, dependency graph, etc.)
4. **STOP and wait for explicit approval**
5. After approval: proceed to implementation

Never skip architecture review.

Never implement first and justify later.

Never assume approval.

---

## 2. Plan Before Coding — Never Start Immediately

Do not write code first. Produce **Implementation Plan** and get approval (or explicit "just do it") before editing. Plan must include:

- **Relevant RFCs** — specific RFC(s) work traces back to.
- **Impacted invariants** — which hard rules (below) change touches.
- **Existing components to reuse** — services, events, commands, queries, aggregates, projections, agents, storage already defined.
- **New components to create** — only what reuse cannot cover, each justified by RFC.
- **Minimal implementation slice** — smallest end-to-end change that delivers value.
- **Verification strategy** — tests / verification scripts that will pin behavior.
- **Risks** — what could break, ambiguous RFC areas, second-source-of-truth dangers.

Cannot map work to RFC? Stop and ask — do not invent behavior.

## 3. Prefer Extending Existing Architecture

Never introduce new abstraction if RFC or codebase already defines one. Before creating anything, search for existing: **service, event, command, query, aggregate, projection, agent, storage**. Extend or compose what exists; only create new when reuse genuinely impossible, and say why.

## 4. Think in Vertical Slices

Never build entire subsystems at once. Implement smallest end-to-end slice providing value, typically:

```
Command → Event → Projection → Query → Verification
```

Ship one working path through layers instead of many half-built services.

## 5. RFC Compliance Review (before implementing)

Before writing slice, explicitly answer each. If any answer is **YES**, stop and explain:

- Does this violate any RFC?
- Does this duplicate existing concept?
- Does this introduce second source of truth?
- Does this bypass Review → Decision → Action?
- Does this introduce mutable canonical state?
- Does this weaken auditability?

## 6. Implementation Order

Build foundations before dependents. Respect this sequence:

1. RFC-031 — service envelopes (`rfcs/031-service-contracts.md`)
2. RFC-013 — event model (`rfcs/013-event-model.md`)
3. RFC-014 — identity / representation model (`rfcs/014-identity-representation-model.md`)
4. RFC-015 — object lifecycle model (`rfcs/015-object-lifecycle-model.md`)
5. RFC-020 — review system (`rfcs/020-review-system.md`)
6. RFC-021 — decision model (`rfcs/021-decision-model.md`)
7. RFC-023 — action model (`rfcs/023-action-model.md`)

Do not implement later layer before prerequisites exist and are tested.

## 7. Hard Rules (Invariants)

Non-negotiable. Any code violating one is wrong by definition.

- **Events are immutable.** Never mutate or delete emitted event; correct via new events.
- **Decisions are explicit and auditable.** Every Decision records who/what, why, inputs.
- **Reviews never commit Decisions.** Review produces findings; cannot enact Decision.
- **Agents never mutate canonical state directly.** They propose; system commits.
- **Agents never call LLM providers directly.** All model access goes through LLM router.
- **Prompt versions are immutable after release.** Released prompts frozen; ship new version.
- **Memory is policy-bound.** Reads/writes to memory honor governing policy; no ad-hoc access.
- **Spaces are bounded contexts.** Keep models, data, logic within their space.
- **Cross-space communication is explicit.** Use defined contracts/events; no hidden coupling.
- **Derived stores are rebuildable.** Never treat projection/cache as source of truth.

## 8. Avoid Speculative Implementation

Do not implement future features, generic frameworks, unused abstractions, placeholder services, or speculative config. Everything must be justified by existing RFC or accepted implementation slice.

## 9. Minimize Architectural Debt

- Prefer deleting code over adding abstractions.
- Prefer composition over inheritance.
- Prefer explicit behavior over magic.
- Prefer deterministic behavior over convenience.

## 10. Challenge Assumptions

Do not blindly implement requested design. If simpler solution satisfies RFCs, explain and recommend it before implementing. If request contradicts long-term architecture, explain why instead of complying silently.

## 11. Leave the Repository Healthier

Every change should improve at least one of: documentation, verification, naming, comments, tests, architecture consistency, or dead-code removal.

## 12. After Implementation — Architecture Review

When slice complete, produce short review:

- **RFC Compliance** — every RFC change satisfies.
- **Invariants** — every invariant preserved.
- **Future RFC impact** — RFCs to update because implementation revealed missing or ambiguous detail.

## 13. After Implementation — Summary

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

## 14. When Unsure

Prefer reading RFC over guessing. If RFC ambiguous or missing, surface gap rather than filling with assumed behavior.

## 15. Build and Verification Operations

Rules:
- Prefer Taskfile commands over raw go/python commands.
- Use `task test` for normal validation.
- Use `task build` before claiming binaries compile.
- Use `task verify:rfc` after changing RFCs or docs under `rfcs/`.
- Use `task graph:rebuild` after changing graph-relevant docs or code.
- Never commit `build/` or `graphify-out/`.
- Do not create ad-hoc build scripts unless Taskfile cannot express operation.
- No command canonical unless in `Taskfile.yml`.

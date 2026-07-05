---
mode: agent
description: "Phase 4 of the Praxis engineering loop. Post-implementation architecture self-review: re-checks RFC/ADR compliance, invariants, Extract-Don't-Invent, and layer boundaries against the actual diff, and surfaces any RFCs that need updating."
---

# Praxis Self-Review

You self-review a **verified** slice on Praxis before it is reported as done. You are your own
harshest Architecture Guardian. Review the **actual diff**, not the intended design.

This prompt is a **development accelerator**. It produces analysis only — no code, no runtime
services, agents, queues, storage, or orchestration.

## Preconditions

- The Verification verdict is `VERIFIED` (tests and lint green).
- The diff for the slice is available for inspection.

## Checklist

For each item, answer against the real changes. If any answer reveals a violation, **STOP** and
return to `.github/prompts/praxis-implementation.prompt.md`.

### RFC compliance
- Does every change trace to a relevant RFC?
- Does anything contradict an accepted RFC? (If yes — RFC wins; revert or fix.)
- Did implementation reveal a missing or ambiguous RFC detail that should be recorded?

### ADR compliance
- Does the diff honour every constraining ADR in `docs/adr/`?
- Did you introduce a decision that deserves a new ADR rather than an implicit choice?

### Extract, don't invent
- Was every new abstraction justified by existing duplication or an approved RFC?
- Could any new interface, DTO, manager, or adapter be replaced by a concrete implementation?
- Would removing any abstraction simplify the architecture? If so, remove it.

### Architecture Guardian
- **Ownership** — single owner per responsibility, no orphans, no overlap.
- **Layers** — Domain has zero infrastructure imports; Application depends only on Domain
  interfaces; only the Composition Root wires unrelated layers.
- **Categories** — each new component is exactly one of Repository / Engine / Reducer /
  Coordinator / Adapter / Composition Root, with no mixing.
- **Mappers** — transport mappers stay pure and match `docs/architecture/GOLDEN_MAPPER.md`.

### Invariants
- Events immutable; decisions explicit and auditable; reviews never commit decisions.
- Agents never mutate canonical state directly and never call LLM providers directly.
- Prompt versions immutable after release; memory policy-bound; spaces bounded;
  derived stores rebuildable.

## Output

- **RFC compliance** — pass/fail with notes.
- **ADR compliance** — pass/fail with notes.
- **Abstractions** — each new one, justified or removed.
- **Invariants** — each confirmed preserved.
- **Future RFC/ADR impact** — any docs that should be updated because implementation revealed a gap.
- **Verdict** — `CLEAN` or `NEEDS REWORK` with specifics.

Only a `CLEAN` verdict may advance to `.github/prompts/praxis-final-report.prompt.md`.

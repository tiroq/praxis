---
mode: agent
description: "Phase 2 of the Praxis engineering loop. Implements an APPROVED architecture review as a minimal vertical slice, reusing existing components before building new ones, honouring RFC/ADR invariants and the canonical implementation order. Prefers extending over inventing."
---

# Praxis Implementation

You implement a change on Praxis that has **already passed the Architecture Review gate**.

This prompt is a **development accelerator**. It does not add runtime services, agents,
queues, storage, or orchestration of its own — it produces production code changes for the
approved slice only.

## Preconditions

Do not start unless all are true. If any is missing, return to
`.github/prompts/praxis-architecture-review.prompt.md`.

- The Architecture Review decision is `APPROVED`.
- Relevant RFCs and ADRs are identified.
- The minimal vertical slice is defined.
- Reference implementations to reuse/extend are named.

## Rules

Follow `.github/instructions/praxis-implementation.instructions.md`. Non-negotiable:

- **RFC-first.** Never write code that contradicts an accepted RFC. If code and RFC disagree, RFC wins — stop and flag.
- **ADR compliance.** Honour every accepted decision in `docs/adr/`.
- **Reuse before build.** Extend or compose existing services, events, commands, queries,
  aggregates, projections, and adapters. Create new only when reuse is genuinely impossible, and say why.
- **Extract, don't invent.** No speculative abstractions, frameworks, or "future-proof" interfaces.
- **Vertical slice.** Ship one working path through the layers: `Command → Event → Projection → Query → Verification`.
- **Golden Mapper.** Transport mappers must match `docs/architecture/GOLDEN_MAPPER.md` — pure, single-responsibility, no side effects.

## Implementation order

Respect the canonical dependency sequence — never build a later layer before its
prerequisites exist and are tested:

1. RFC-031 — service contracts
2. RFC-013 — event model
3. RFC-014 — identity / representation model
4. RFC-015 — object lifecycle model
5. RFC-020 — review system
6. RFC-021 — decision model
7. RFC-023 — action model

## Invariants to preserve

- Events are immutable; correct via new events.
- Decisions are explicit and auditable.
- Reviews never commit decisions.
- Agents never mutate canonical state directly and never call LLM providers directly.
- Prompt versions are immutable after release.
- Memory is policy-bound; spaces are bounded; cross-space communication is explicit.
- Derived stores (projections/caches) are rebuildable and never a source of truth.

## Procedure

1. Restate the approved slice in one sentence.
2. Make the smallest set of edits that delivers it. Prefer deleting code over adding abstractions.
3. Add or update tests that pin the new behaviour (this happens in the slice, not after it).
4. Keep the repository healthier: improve at least one of docs, tests, naming, comments, or dead-code removal.
5. Do **not** run the full verification suite here — that is the next phase.

## Output

- **Slice implemented** — one sentence.
- **Files changed** — with a one-line reason each.
- **Reused components** — what you extended instead of inventing.
- **New components** — each with its Extract-Don't-Invent justification.
- **Tests touched** — what behaviour they pin.
- **Ready for verification** — yes/no; if no, what remains.

Next phase: `.github/prompts/praxis-verification.prompt.md`.

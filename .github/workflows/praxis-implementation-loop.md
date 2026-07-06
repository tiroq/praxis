# Praxis Implementation Loop

A lightweight, prompt-based engineering workflow for building **Praxis itself**.

> **Scope guard.** This is a development accelerator, not part of Praxis runtime architecture.
> It introduces **no services, agents, queues, storage, orchestration, or product code**. Every
> artefact in this workflow is documentation/prompt only. Despite living under `.github/workflows/`,
> this file is a Markdown guide — it is **not** a GitHub Actions pipeline.

## Purpose

Turn the Praxis architectural discipline (RFC-first, ADR compliance, Extract-Don't-Invent,
Architecture Guardian, tests + lint before done) into a repeatable five-phase loop that any
contributor or coding agent can follow for a single change.

## The loop

```
┌─────────────────────────────────────────────────────────────────┐
│  1. Architecture Review   → gate: RFC-first, ADR, reuse, extract  │
│  2. Implementation        → minimal vertical slice, reuse first   │
│  3. Verification          → tests + lint (Taskfile), must pass    │
│  4. Self-Review           → re-check invariants against the diff   │
│  5. Final Report          → files, commands, risks, remaining debt │
└─────────────────────────────────────────────────────────────────┘
        │                                                     ▲
        └──────────── any failure sends you back ─────────────┘
```

Each phase is a prompt file. Run them in order:

| Phase | Prompt | Gate |
|---|---|---|
| 1. Architecture Review | [`.github/prompts/praxis-architecture-review.prompt.md`](../prompts/praxis-architecture-review.prompt.md) | `APPROVED` to proceed |
| 2. Implementation | [`.github/prompts/praxis-implementation.prompt.md`](../prompts/praxis-implementation.prompt.md) | slice implemented |
| 3. Verification | [`.github/prompts/praxis-verification.prompt.md`](../prompts/praxis-verification.prompt.md) | `VERIFIED` (tests + lint) |
| 4. Self-Review | [`.github/prompts/praxis-self-review.prompt.md`](../prompts/praxis-self-review.prompt.md) | `CLEAN` |
| 5. Final Report | [`.github/prompts/praxis-final-report.prompt.md`](../prompts/praxis-final-report.prompt.md) | report emitted |

## Enforced gates

The loop refuses to advance unless each rule holds:

1. **RFC-first review.** Every change traces to an RFC in `rfcs/`. RFCs are the architectural
   source of truth; if code and RFC disagree, the RFC wins. Ambiguous RFC → stop and flag.
2. **ADR compliance.** Every change honours accepted decisions in `docs/adr/`. A contradiction
   requires a new ADR, never a silent deviation.
3. **Reference Implementation search.** Before building anything, find the closest existing
   component to reuse or extend (services, events, commands, queries, aggregates, projections,
   adapters; transport mappers follow `docs/architecture/principles/GOLDEN_MAPPER.md`).
4. **Extract, don't invent.** New abstractions are valid only when real duplication exists or an
   approved RFC defines them. "We may need it later" is never sufficient.
5. **Architecture Guardian checks.** Ownership, layer boundaries, single-category components, and
   the Praxis invariants (immutable events, auditable decisions, reviews never commit decisions,
   agents never mutate canonical state or call LLM providers directly, policy-bound memory,
   bounded spaces, rebuildable derived stores).
6. **Tests before completion.** `task test` (and `task test:storage`, `task build` as relevant)
   must pass.
7. **Lint before completion.** `go vet ./...`, `gofmt -l .` (must be empty), and `task verify:rfc`
   must pass. No `--no-verify` shortcuts.
8. **Final report.** Changed files, commands run with results, risks, and remaining debt.

## Feedback rules

- **Architecture Review → BLOCKED:** resolve the RFC/ADR/abstraction issue before any code.
- **Verification → BLOCKED:** fix the failing command; return to Implementation. A change is not
  complete until tests **and** lint pass.
- **Self-Review → NEEDS REWORK:** return to Implementation with the specific violation.
- Only after a `CLEAN` self-review is the Final Report produced.

## Reference material

- Architecture Guardian rules — [`.github/instructions/praxis-architecture-guardian.instructions.md`](../instructions/praxis-architecture-guardian.instructions.md)
- Implementation discipline — [`.github/instructions/praxis-implementation.instructions.md`](../instructions/praxis-implementation.instructions.md)
- RFCs — `rfcs/`
- ADRs — `docs/adr/`
- Canonical commands — `Taskfile.yml`

## What this workflow must never do

- Add runtime services, agents, queues, storage, or orchestration engines.
- Introduce an agents framework or product code.
- Modify Praxis product behaviour.
- Weaken or bypass any architectural gate to force a change through.

---
mode: agent
description: "Phase 5 of the Praxis engineering loop. Produces the final report for a completed slice: changed files, commands run with results, RFCs/ADRs touched, risks, and remaining architectural debt. Closes the loop."
---

# Praxis Final Report

You produce the **final report** for a slice that has passed implementation, verification, and
self-review. This is the closing artefact of the Praxis engineering loop.

This prompt is a **development accelerator**. It emits a report only — no code, no runtime
services, agents, queues, storage, or orchestration.

## Preconditions

- Verification verdict: `VERIFIED`.
- Self-review verdict: `CLEAN`.

## Report format

Output exactly the following sections. Be concrete; cite file paths and RFC/ADR numbers.

### 1. Summary
One or two sentences: what the slice delivered and which RFC(s) it advances.

### 2. Changed files
A list of every changed/added/removed file with a one-line reason each.

### 3. Commands run
Every command executed during implementation and verification, each with its result:

| Command | Result |
|---|---|
| `task test` | PASS/FAIL |
| `go vet ./...` | PASS/FAIL |
| `gofmt -l .` | empty/`<files>` |
| `task verify:rfc` | PASS/FAIL |
| … | … |

### 4. RFCs & ADRs
- RFCs implemented or touched (by number).
- ADRs honoured; any new ADR proposed.

### 5. Architecture notes
- New Commands / Events / Queries / Services (if any).
- Reused components (Extract-Don't-Invent outcome).
- Invariants preserved.

### 6. Risks
Known risks introduced by this change: fragile areas, ambiguous RFC regions, second-source-of-truth
dangers, anything a reviewer should watch.

### 7. Remaining debt
- Deferred work and why it was deferred.
- Follow-up RFC/ADR updates needed.
- Recommended next implementation slice.

## Rules

- Report only what actually happened; do not claim commands or tests that were not run.
- If any precondition is unmet, say so and stop — do not fabricate a clean report.
- Never recommend committing `build/` or `graphify-out/`.

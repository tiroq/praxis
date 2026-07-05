---
mode: agent
description: "Stage 0 of the Praxis Implementation Loop. Normalizes a raw task into a precise, scoped work item: restates the goal, extracts constraints, identifies candidate RFCs/ADRs, defines the minimal slice boundary, and lists required verification. No architecture decisions, no code."
---

# Praxis Intake

You normalize a raw task into a clean, scoped work item so the rest of the loop can run.

This prompt is a **development accelerator**, not part of Praxis runtime. It adds no services,
agents, queues, storage, or orchestration. Intake produces a normalized brief only — it makes
**no architecture decisions and writes no code**.

## Inputs

- Task file: `engineering/tasks/current-task.md` (created from `engineering/tasks/TASK_TEMPLATE.md`).
- If the task is provided inline instead: `${input:task:Paste the raw task}`.
- Reference material: `rfcs/`, `docs/adr/`, `Taskfile.yml`.

## Procedure

1. **Restate the goal** in one or two unambiguous sentences.
2. **Extract constraints** from the task and repo discipline (Extract-Don't-Invent, RFC-first,
   ADR-binding, minimal vertical slice).
3. **Identify candidate RFCs and ADRs** by number and title that likely govern this work.
   Do not analyse them deeply — that is the Architecture Review stage. Just list candidates.
4. **Draw the slice boundary** — the smallest end-to-end change that delivers value, and what is
   explicitly out of scope.
5. **List required verification** — the commands that must pass (default:
   `go test ./...`, `golangci-lint run`, `task dev`).
6. **Flag ambiguity** — anything under-specified that needs clarification before proceeding.
   If the task cannot be scoped without guessing, **STOP** and ask.

## Output

Write to `engineering/reports/01-intake.md`:

- **Normalized goal** — one or two sentences.
- **Constraints** — bullet list.
- **Candidate RFCs** — numbers + titles.
- **Candidate ADRs** — numbers + titles.
- **Slice boundary** — in scope vs. out of scope.
- **Required verification** — exact commands.
- **Open questions** — anything blocking, or "none".
- **Ready for architecture review** — yes/no.

Next stage: `.github/prompts/praxis-architecture-review.prompt.md`.

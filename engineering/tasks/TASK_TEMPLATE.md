# Task

> Copy this template to `engineering/tasks/current-task.md` and fill it in before
> starting the Praxis Implementation Loop. One task = one minimal vertical slice.

## Goal

Implement:

<!-- One or two sentences. What behaviour should exist after this slice? -->

## Constraints

- Do not introduce speculative abstractions (Extract, don't invent).
- Follow RFCs (`rfcs/`), ADRs (`docs/adr/`), policies, and reference implementations.
- Minimal vertical slice only (`Command → Event → Projection → Query → Verification`).
- Never contradict an accepted RFC; if code and RFC disagree, the RFC wins.

## Expected Result

<!-- Observable outcome. What can be demonstrated or queried once this is done? -->

## Required Verification

- `go test ./...`
- `golangci-lint run`
- `task dev`

## Out of Scope

<!-- What this task explicitly does NOT do. Prevents scope creep. -->

## Notes

<!-- Relevant RFC/ADR numbers, reference implementations to reuse, known risks. -->

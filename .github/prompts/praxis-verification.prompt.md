---
mode: agent
description: "Phase 3 of the Praxis engineering loop. Runs tests and lint before a change can be called complete, using the canonical Taskfile commands. Verification failures block completion. Records every command and its result."
---

# Praxis Verification

You verify an implemented slice on Praxis. **A change is not complete until tests and lint pass.**

This prompt is a **development accelerator**. It runs verification commands only; it adds no
runtime services, agents, queues, storage, or orchestration.

## Canonical commands

Prefer Taskfile commands over raw `go`/`python` invocations. No command is canonical unless it
lives in `Taskfile.yml`.

### Tests (required)
- `task test` — full Go test suite (normal validation).
- `task test:storage` — storage layer tests, when storage code changed.
- `task test:coverage` — coverage profile into `build/coverage.out`, when coverage matters.
- `task build` — required before claiming binaries compile.

### Lint / hygiene (required)
- `go vet ./...` — static analysis.
- `gofmt -l .` — formatting check; output **must be empty**.
- `task verify:rfc` — RFC hygiene, required after changing anything under `rfcs/` or docs.

### Optional, when relevant
- `task graph:rebuild` — after changing graph-relevant docs or code.
- `task smoke:nats` — real NATS JetStream smoke test, only when transport changed and a broker is available.

## Rules

- Run tests **and** lint. Passing one is not enough.
- If any command fails, verification **fails**. Do not mark the change complete.
  Return to `.github/prompts/praxis-implementation.prompt.md` with the failure details.
- Never commit `build/` or `graphify-out/`.
- Do not weaken, skip, or `--no-verify` any check to force a pass.
- Do not invent ad-hoc build scripts; if the Taskfile cannot express it, say so.

## Procedure

1. Determine which command set applies from the files changed.
2. Run tests, then lint. Capture the exact command and its pass/fail result.
3. On failure, stop and report the first failing command with its output.

## Output

- **Commands run** — exact command + `PASS`/`FAIL` for each.
- **Test result** — overall pass/fail; note any skipped suites and why.
- **Lint result** — `go vet`, `gofmt -l .`, `task verify:rfc` outcomes.
- **Verdict** — `VERIFIED` (all green) or `BLOCKED` with the failing command.

Only a `VERIFIED` verdict may advance to `.github/prompts/praxis-self-review.prompt.md`.

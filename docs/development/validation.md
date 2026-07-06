# Validation & Verification Discipline

**Execute every applicable verification command. Only skip commands that genuinely do not exist.**

Never claim commands passed unless executed.

---

## Verification Operations

Run every applicable verification. Examples:

```bash
go test ./...
golangci-lint run
task test
task build
task verify:rfc
task graph:rebuild
task smoke:*
task report
```

---

## Build and Verification Rules

- **Prefer Taskfile commands** over raw go/python commands
- **Use `task test`** for normal validation
- **Use `task build`** before claiming binaries compile
- **Use `task verify:rfc`** after changing RFCs or docs under `rfcs/`
- **Use `task graph:rebuild`** after changing graph-relevant docs or code
- **Never commit** `build/` or `graphify-out/`
- **Do not create ad-hoc build scripts** unless Taskfile cannot express operation
- **No command is canonical** unless in `Taskfile.yml`

---

## Continuous Validation During Implementation

**During implementation,** before committing code, validate:

- RFC compliance (see [implementation-guide.md](implementation-guide.md))
- Architecture compliance (no invariant violations)
- No unnecessary abstractions
- No unnecessary interfaces
- No unnecessary DTOs
- No unnecessary services
- No hidden business logic

---

## Verification Phase — After Implementation

Run all tests and validation commands. Expected output:

- All unit tests pass
- All integration tests pass
- Linting passes
- Build succeeds
- Smoke tests pass (if applicable)
- Verification report passes

---

## Repair Loop — If Verification Fails

If verification fails:

1. Analyze the failure
2. Repair the issue
3. Verify again
4. Repeat

**Maximum repair iterations: 5.**

If issue cannot be fixed in 5 iterations, surface the blocker and document.

---

## After Verification — Final Review

Verify that after your changes:

- Architecture remains simpler (fewer files, fewer abstractions, fewer layers)
- Ownership remains correct (single owner per responsibility)
- Boundaries remain intact (no cross-layer violations)
- No duplication introduced
- No unnecessary abstractions added
- Code follows project conventions

---

## When Unsure

Prefer reading RFC over guessing. If RFC ambiguous or missing, surface gap rather than filling with assumed behavior.

Do not implement when blocked by ambiguity.

Ask for clarification instead of inventing.

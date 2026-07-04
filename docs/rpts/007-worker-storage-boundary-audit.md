# Worker Storage Integration Boundary Audit

**Date**: 2026-07-04  
**Auditor**: GitHub Copilot  
**Scope**: Worker storage integration architecture boundary compliance

## Executive Summary

**VERDICT: ✅ PASS**

The worker storage integration now fully complies with architectural boundaries. The
previous violation (storage importing kernel) has been corrected by moving the adapter
to the composition root (services/worker).

## Architecture Boundary Requirements

The system enforces strict layering:

1. **internal/core/kernel** (domain logic) must not import:
   - internal/storage
   - database/sql (except sql/driver for interfaces)
   
2. **internal/storage** (persistence) must not import:
   - internal/core/kernel
   
3. **services/worker** (composition root) is allowed to import both:
   - internal/core/kernel
   - internal/storage

## Audit Findings

### Initial State (FAIL)

Boundary violation detected:

```bash
$ go list -deps ./internal/storage | rg 'internal/core/kernel'
github.com/tiroq/praxis/internal/core/kernel  # ❌ VIOLATION
```

**Root Cause**: `internal/storage/adapter.go` contained `EventRecorderAdapter` which
imported `kernel.EventRecord` type. This violated the principle that storage (infrastructure)
must not depend on kernel (domain).

### Corrective Action

**Fix Applied**: Move adapter to composition root

1. **Created**: `services/worker/storage_recorder.go`
   - Contains `eventRecorderAdapter` (unexported, worker-local)
   - Imports both `kernel.EventRecord` and `eventstore.EventRecord`
   - Lives at composition boundary where both dependencies are valid

2. **Updated**: Worker service files
   - `services/worker/main.go`: changed `storage.NewEventRecorderAdapter` → `newEventRecorderAdapter`
   - `services/worker/main_test.go`: same update

3. **Removed**: Boundary-violating code
   - Deleted `internal/storage/adapter.go`
   - Removed adapter tests from `internal/storage/storage_test.go`
   - Removed `kernel` import from `internal/storage/storage_test.go`

### Final State (PASS)

All boundaries verified:

```bash
# ✅ Storage does NOT import kernel
$ go list -deps ./internal/storage | rg 'internal/core/kernel'
(empty output - no violation)

# ✅ Kernel does NOT import storage (only sql/driver for interfaces)
$ go list -deps ./internal/core/kernel | rg 'internal/storage|database/sql|sqlite'
database/sql/driver

# ✅ Worker imports BOTH (composition root)
$ go list -deps ./services/worker | rg 'internal/storage|internal/core/kernel'
github.com/tiroq/praxis/internal/core/kernel
github.com/tiroq/praxis/internal/storage/eventstore
github.com/tiroq/praxis/internal/storage/sqlite
github.com/tiroq/praxis/internal/storage
```

## Verification Tests

### Unit Tests

```bash
$ go test ./...
ok      github.com/tiroq/praxis/internal/storage        1.402s
ok      github.com/tiroq/praxis/services/worker         1.965s
# All tests pass (14 packages)
```

### Lint Check

```bash
$ golangci-lint run
0 issues.
```

### Build Check

```bash
$ task dev
# All binaries built successfully
```

### Integration Tests (smoke:praxis)

```json
{
  "worker_flow_ok": true,
  "message": "validated praxis publish/watch flow over NATS JetStream"
}
```

**Storage Verification**:
```json
{
  "event_count": 1,
  "event_type": "kernel.pipeline.completed",
  "verification": "passed"
}
```

✅ Kernel pipeline execution events are correctly persisted to storage via the adapter.

### Coverage Report

Generated via `task report`:
- All smoke tests passed
- Coverage reports updated in `build/`

## Compliance Matrix

| Requirement | Status | Evidence |
|-------------|--------|----------|
| Kernel does not import storage | ✅ PASS | `go list -deps` shows no storage imports |
| Storage does not import kernel | ✅ PASS | `go list -deps` shows no kernel imports |
| Worker composes both | ✅ PASS | Worker imports both kernel and storage |
| Persisted events verified | ✅ PASS | smoke:praxis verified `kernel.pipeline.completed` in DB |
| Tests pass | ✅ PASS | `go test ./...` all green |
| Linter clean | ✅ PASS | `golangci-lint run` zero issues |
| Build succeeds | ✅ PASS | `task dev` completed successfully |

## Architectural Notes

### Why This Matters

The original boundary violation (storage → kernel) created a circular conceptual dependency:

- **Kernel** defines business logic and event contracts
- **Storage** persists data using those contracts
- If storage imports kernel, storage becomes tightly coupled to domain

This violates:
- **Dependency Inversion Principle**: Infrastructure should depend on abstractions, not domain
- **Clean Architecture**: Domain should be the innermost layer with no outward dependencies
- **Praxis RFC-030**: Kernel must remain transport-agnostic and infrastructure-agnostic

### The Fix: Adapter at Composition Root

The adapter pattern is correct, but placement matters:

- ❌ **Bad**: Adapter in `internal/storage` (infrastructure layer imports domain)
- ✅ **Good**: Adapter in `services/worker` (composition root knows about both)

The composition root (worker service) is the only place where it's architecturally sound
to have dependencies on both kernel and storage. This is its purpose: to wire together
independent layers.

### Alternative Considered

Could have made `kernel.EventRecord` an interface instead of a struct, allowing storage
to implement it without importing kernel. **Rejected** because:

1. EventRecord is data, not behavior—interfaces for data are antipatterns
2. Would require kernel to define an interface in a separate package (awkward)
3. Current solution is simpler and more idiomatic

## Recommendations

1. ✅ **Keep adapter in services/worker** — do not move back to internal/storage
2. ✅ **Add architecture test** — consider adding a test that fails if storage imports kernel
3. ✅ **Document pattern** — this is the canonical way to adapt between layers at composition boundary

## Conclusion

The worker storage integration now satisfies all architectural invariants:

- ✅ Kernel remains infrastructure-agnostic
- ✅ Storage remains domain-agnostic  
- ✅ Worker composes both at the boundary
- ✅ Event recording works end-to-end
- ✅ All tests pass
- ✅ No linting issues

**Final Verdict**: **PASS** — All boundary requirements met.

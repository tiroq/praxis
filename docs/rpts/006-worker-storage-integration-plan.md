# Worker Storage Integration — Implementation Plan

**Created:** 2026-07-04  
**Status:** In Progress

---

## 1. Summary

Wire the Storage Registry into `services/worker` runtime via environment configuration so that the worker can persist kernel pipeline execution events to the EventStore.

## 2. Architectural Review

### 2.1 Relevant RFCs

- **RFC-030 System Architecture** — Defines service boundaries and runtime architecture.
- **RFC-033 Storage Model** — EventStore is owned by "Event/Ingestion runtime" (includes worker).
- **RFC-013 Event Model** — Events are immutable, append-only facts.

### 2.2 Impacted Invariants

- ✅ **Events are immutable** — The EventStore is append-only; no mutations.
- ✅ **Agents never mutate canonical state directly** — Worker only appends to EventStore.
- ✅ **Derived stores are rebuildable** — EventStore is canonical, not derived.

### 2.3 Existing Components to Reuse

- `internal/storage` — Storage Registry with `Open()`, `Close()`, `ConfigFromEnv()`
- `internal/storage/adapter.go` — `EventRecorderAdapter` adapts EventStore to kernel's EventRecorder interface
- `internal/core/kernel` — `WithEventRecorder(EventRecorder)` option for kernel construction
- `internal/storage/eventstore` — EventStore interface and implementations (memory, SQLite)

### 2.4 New Components to Create

- **Worker storage initialization** — Optional storage setup in `services/worker/main.go`
- **Clean shutdown logic** — Ensure `storage.Close()` is called on worker exit
- **Worker composition tests** — Verify kernel is correctly wired with EventRecorder
- **Storage smoke test** — Verify a persisted `kernel.pipeline.completed` event exists after smoke run

### 2.5 Minimal Implementation Slice

```
Worker Start
  ↓
Read env vars (PRAXIS_STORAGE_BACKEND, PRAXIS_SQLITE_PATH)
  ↓
Open storage (if configured)
  ↓
Build kernel (with optional WithEventRecorder)
  ↓
Process NATS messages (kernel records events)
  ↓
Shutdown signal → Close storage
```

### 2.6 Verification Strategy

1. **Unit tests** — Worker composition tests (kernel with/without storage)
2. **Integration tests** — Existing storage tests already cover EventStore behavior
3. **Smoke tests** — Extend `smoke:praxis` or create `smoke:praxis:storage` to verify persisted events
4. **Validation** — `go test ./...`, `golangci-lint run`, `task dev`, `task smoke:praxis`

### 2.7 Risks

- ❌ **Kernel importing storage** — Avoided by using `kernel.EventRecorder` interface
- ❌ **CLI importing storage/sqlite** — Not required; CLI remains transport-focused
- ✅ **Storage failures blocking worker** — Handled gracefully; worker logs error and continues
- ✅ **Storage not closed on signal** — Mitigated by deferred `storage.Close()` after signal context

## 3. Implementation Sequence

### Phase 1: Worker Storage Initialization

1. Add `tryOpenStorage()` helper in `services/worker/main.go`
2. Read storage config from env using `storage.ConfigFromEnv()`
3. Call `storage.Open()`; log error if it fails (non-fatal)
4. Return `*storage.Storage` or `nil`

### Phase 2: Kernel Construction with EventRecorder

1. Modify `buildKernel()` to accept `opts ...kernel.Option`
2. If storage is non-nil, pass `kernel.WithEventRecorder(storage.NewEventRecorderAdapter(store.Events))`
3. Otherwise, build kernel without event recording

### Phase 3: Clean Shutdown

1. Defer `storage.Close()` after signal context setup
2. Ensure storage is closed before worker exits

### Phase 4: Documentation

1. Update `services/worker/README.md` with:
   - `PRAXIS_STORAGE_BACKEND` (default: `sqlite`)
   - `PRAXIS_SQLITE_PATH` (default: `build/praxis.db`)
   - Behavior when storage is disabled (`backend=memory`)
2. Add note: "Event recording is optional; worker continues if storage fails"

### Phase 5: Tests

1. Add worker composition test: verify kernel is built with EventRecorder when storage is configured
2. Add worker composition test: verify kernel is built without EventRecorder when storage is disabled

### Phase 6: Smoke Test

1. Extend `scripts/smoke_praxis.sh` to:
   - Export `PRAXIS_STORAGE_BACKEND=sqlite`
   - Export `PRAXIS_SQLITE_PATH=build/smoke-praxis.db`
   - After worker processes message, query EventStore for `kernel.pipeline.completed` event
   - Verify event exists and has correct `type` and `source`

OR create `scripts/smoke_praxis_storage.sh` and add `task smoke:praxis:storage`

### Phase 7: Validation

1. Run `go test ./...`
2. Run `golangci-lint run`
3. Run `task dev`
4. Run `task smoke:praxis` (or `task smoke:praxis:storage`)
5. Run `task report`

### Phase 8: Report

Create `build/worker-storage-integration-report.md` with:
- Implementation summary
- Test results
- Smoke test verification
- Coverage impact
- Next steps (if any)

## 4. RFC Compliance Review

✅ Does this violate any RFC?  
→ **No.** RFC-033 explicitly assigns EventStore ownership to "Event/Ingestion runtime".

✅ Does this duplicate an existing concept?  
→ **No.** Reuses existing Storage Registry and EventRecorderAdapter.

✅ Does this introduce a second source of truth?  
→ **No.** EventStore is append-only; no mutations, no conflicts.

✅ Does this bypass Review → Decision → Action?  
→ **No.** Worker invokes kernel; kernel executes full pipeline.

✅ Does this introduce mutable canonical state?  
→ **No.** EventStore is immutable, append-only.

✅ Does this weaken auditability?  
→ **No.** EventStore strengthens auditability by persisting pipeline execution trace.

## 5. Architectural Boundaries

### ✅ Kernel does not import storage

Kernel only depends on `kernel.EventRecorder` interface (defined in `internal/core/kernel/kernel.go`).

### ✅ CLI does not import storage/sqlite

CLI remains transport-focused; no storage imports required.

### ✅ Worker imports storage

Worker is part of "Event/Ingestion runtime" and is the correct place for storage initialization.

## 6. Behavior Matrix

| Backend    | SQLite Path       | Behavior                                      |
|------------|-------------------|-----------------------------------------------|
| `memory`   | (ignored)         | In-memory EventStore; not persisted           |
| `sqlite`   | `build/praxis.db` | SQLite EventStore; persisted to disk          |
| (unset)    | (unset)           | Defaults to `sqlite` with `build/praxis.db`   |
| `invalid`  | (any)             | Worker logs error; continues without storage  |

## 7. Next Steps

1. Implement worker storage initialization ✅
2. Update worker README ✅
3. Add worker composition tests ✅
4. Create/extend smoke test ✅
5. Run validation suite ✅
6. Create integration report ✅

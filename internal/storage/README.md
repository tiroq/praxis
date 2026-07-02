# Storage Layer

This directory contains storage implementations for Praxis persistent data.

## Architecture

Storage in Praxis follows strict architectural boundaries defined in RFC-033:

- **Storage is service-owned** — every service owns its persistence
- **No shared database integration** — services integrate through Commands, Events, and Queries
- **Domain-facing interfaces** — storage contracts are technology-neutral
- **Adapters live outside core** — concrete implementations (SQLite, Postgres, etc.) are adapters

## EventStore

The EventStore is the append-only persistent store for immutable Event Records.

### Interface

```go
type EventStore interface {
    Append(ctx context.Context, event EventRecord) error
    Get(ctx context.Context, id string) (EventRecord, error)
    List(ctx context.Context, filter ListFilter) ([]EventRecord, error)
    Close() error
}
```

### Implementations

1. **Memory** (`eventstore/memory.go`)
   - In-memory implementation
   - Thread-safe with sync.RWMutex
   - Used for tests and local development
   - All events lost on process exit

2. **SQLite** (`sqlite/eventstore.go`)
   - Durable SQLite-backed implementation
   - Uses modernc.org/sqlite (pure Go)
   - Supports both file-based and `:memory:` databases
   - WAL mode enabled for better concurrency

### Design Principles

1. **Append-only** — events are never updated or deleted
2. **Immutable** — EventRecord fields cannot be changed after storage
3. **Ordered** — List returns events deterministically (occurred_at ASC, id ASC)
4. **Validated** — required fields and JSON payload are validated on Append
5. **Isolated** — no direct coupling to Kernel or other domain logic

## Boundary Rules

### What Kernel CAN do:
- Define domain events (internal/core/kernel/event.go)
- Emit events during business logic
- Accept optional EventStore via dependency injection

### What Kernel MUST NOT do:
- Import `database/sql`
- Import `internal/storage/sqlite`
- Import `modernc.org/sqlite`
- Know about storage implementation details
- Couple to any specific persistence technology

### What storage adapters CAN do:
- Import domain interfaces from `internal/storage/eventstore`
- Map EventRecord to/from storage-specific formats
- Implement EventStore interface
- Handle storage-specific concerns (transactions, indexes, migrations)

### What storage adapters MUST NOT do:
- Import `internal/core/kernel`
- Modify Kernel behavior
- Expose storage-specific types to domain code

## Future Work

The following storage categories are **deferred** until their architectural decisions are resolved:

- **Artifact Store** — deferred pending RFC clarification
- **Canonical Object Store** — deferred until ADR-001/ADQ-001 resolved
- **Postgres adapter** — deferred until multi-database strategy defined
- **Migration framework** — deferred until schema versioning strategy accepted

## Testing

### Run storage tests

```bash
task test:storage
```

### Coverage requirements

- `internal/storage/eventstore` ≥ 90%
- `internal/storage/sqlite` ≥ 80%

### Test strategy

Contract tests in `eventstore/eventstore_test.go` run against both Memory and SQLite implementations to ensure behavioral equivalence.

## Import Verification

To verify boundary enforcement:

```bash
# Kernel must not import storage implementations
go list -f '{{.Imports}}' ./internal/core/kernel | grep -E 'database/sql|sqlite'

# CLI must not import SQLite directly
go list -f '{{.Imports}}' ./internal/cli/praxiscli | grep sqlite
```

Both commands should return nothing (empty output = compliance).

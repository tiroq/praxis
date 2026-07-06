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
It is the canonical event source used for replay and reconstruction of derived state.

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

## Conversation History Projection

Conversation history persistence in `internal/storage/conversationstore` is a concrete
SQLite-backed projection for read/reconstruction convenience. It is derived from immutable
events and can be rebuilt. It is not canonical truth.

## Storage Registry

The Storage Registry is the composition point for all Praxis persistence. It owns the construction of backend-specific stores and provides a single point of configuration and lifecycle management.

### Interface

```go
type Storage struct {
    Events eventstore.EventStore
}

func Open(ctx context.Context, cfg Config) (*Storage, error)
func (s *Storage) Close() error
```

### Configuration

Storage backend is selected via `Config`:

```go
type Config struct {
    Backend    string // "memory" or "sqlite"
    SQLitePath string // path to SQLite database file
}
```

Environment variables:
- `PRAXIS_STORAGE_BACKEND` — backend selection (default: `sqlite`)
- `PRAXIS_SQLITE_PATH` — SQLite database path (default: `build/praxis.db`)

### Usage

```go
// From environment
cfg := storage.ConfigFromEnv()
store, err := storage.Open(ctx, cfg)
if err != nil {
    return err
}
defer store.Close()

// Use the event store
err = store.Events.Append(ctx, eventRecord)
```
` ≥ 85% (includes registry)
- `internal/storage/eventstore` ≥ 90%
- `internal/storage/sqlite` ≥ 80%

### Test strategy

Contract tests in `eventstore/eventstore_test.go` run against both Memory and SQLite implementations to ensure behavioral equivalence.

Storage Registry tests in `storage_test.go` verify:
- Config parsing (defaults, environment variables)
- Backend selection (memory, sqlite, unsupported)
- Lifecycle management (open, close, persistence)
- Kernel DI boundary (EventRecorder adapter)
Unknown backends return `ErrUnsupportedBackend`.

### Kernel Dependency Injection

The Kernel can optionally record pipeline execution events via an EventRecorder interface:

```go
// Create storage
store, _ := storage.Open(ctx, storage.DefaultConfig())
defer store.Close()

// Create adapter
adapter := storage.NewEventRecorderAdapter(store.Events)

// Create kernel with event recording
kernel := kernel.New(reviewer, decisionMaker, planner, 
    kernel.WithEventRecorder(adapter))

// Run pipeline — events are recorded automatically
result, err := kernel.Run(ctx, event)
```

**Important**: The Kernel only depends on the `EventRecorder` interface defined in `internal/core/kernel/kernel.go`. It does NOT import any storage adapter packages. The storage package provides the adapter that bridges `EventRecorder` to `eventstore.EventStore`.

### Why This Boundary Matters

The Storage Registry allows:
- Kernel to optionally persist events WITHOUT coupling to SQLite
- CLI and services to select backends via configuration
- Tests to use memory storage without file I/O
- Future backends (Postgres, etc.) to be added without touching Kernel

The Kernel's `EventRecorder` interface uses only simple types (string, time, json.RawMessage) that the kernel package already depends on. This avoids circular dependencies and preserves the kernel's adapter-free guarantee.

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

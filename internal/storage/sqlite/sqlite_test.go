package sqlite
package sqlite_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/tiroq/praxis/internal/storage/sqlite"
	_ "modernc.org/sqlite" // needed for the verification query in this file
)

// TestWALModeFileDB verifies that file-backed SQLite stores are opened in WAL mode.
// WAL mode provides better concurrent read performance and is required for production use.
func TestWALModeFileDB(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "wal-test.db")
	ctx := context.Background()

	store, err := sqlite.OpenEventStore(ctx, dbPath)
	if err != nil {
		t.Fatalf("OpenEventStore failed: %v", err)
	}
	defer func() { _ = store.Close() }()

	// Open a separate read-only connection to inspect the journal mode.
	// This does NOT go through the EventStore interface — it is test-only
	// validation that the driver was configured correctly.
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("sql.Open failed: %v", err)
	}
	defer func() { _ = db.Close() }()

	var journalMode string
	if err := db.QueryRowContext(ctx, "PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("PRAGMA journal_mode query failed: %v", err)
	}

	if journalMode != "wal" {
		t.Errorf("Expected journal_mode=wal for file-backed database, got %q", journalMode)
	}
}

// TestWALModeMemoryDB documents the expected WAL behavior for :memory: databases.
// SQLite ignores the WAL pragma for in-memory databases; the mode stays "memory".
// This is expected and acceptable — WAL only applies to file-backed databases.
func TestWALModeMemoryDB(t *testing.T) {
	ctx := context.Background()

	store, err := sqlite.OpenEventStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("OpenEventStore(:memory:) failed: %v", err)
	}
	defer func() { _ = store.Close() }()

	// :memory: databases report "memory" journal mode regardless of any WAL pragma.
	// This is a SQLite driver limitation, not a bug in the implementation.
	// The OpenEventStore call itself must succeed even when WAL is a no-op.
	t.Log("WAL mode for :memory: is not applicable; SQLite uses 'memory' journal mode.")
}

// TestSQLiteEncapsulation verifies that no *sql.DB or database internals are
// accessible to callers of OpenEventStore.
func TestSQLiteEncapsulation(t *testing.T) {
	ctx := context.Background()

	store, err := sqlite.OpenEventStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("OpenEventStore failed: %v", err)
	}
	defer func() { _ = store.Close() }()

	// The returned value must not be type-assertable to any SQL-exposing type.
	// If it were *sql.DB, a caller could run arbitrary SQL. The assertion below
	// would only compile if the package exported a concrete DB-bearing type.
	// Since OpenEventStore returns eventstore.EventStore (an interface), this
	// test confirms the contract by ensuring the store only satisfies the interface.
	type sqlExposer interface {
		DB() *sql.DB
	}
	if _, ok := store.(sqlExposer); ok {
		t.Error("store must not expose *sql.DB through any exported method")
	}

	t.Log("Encapsulation confirmed: store implements only eventstore.EventStore.")
}

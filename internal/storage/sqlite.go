package storage

import (
	"context"
	"fmt"

	"github.com/tiroq/praxis/internal/storage/eventstore"
	"github.com/tiroq/praxis/internal/storage/sqlite"
)

// openSQLiteBackend creates a SQLite-backed event store.
// If cfg.SQLitePath is ":memory:", creates an in-memory SQLite database.
// Otherwise, creates or opens the SQLite file at the specified path.
func openSQLiteBackend(ctx context.Context, cfg Config) (eventstore.EventStore, error) {
	// Ensure the directory exists for file-based SQLite
	if err := cfg.ensureSQLiteDir(); err != nil {
		return nil, fmt.Errorf("failed to create SQLite directory: %w", err)
	}

	store, err := sqlite.OpenEventStore(ctx, cfg.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite event store at %q: %w", cfg.SQLitePath, err)
	}

	return store, nil
}

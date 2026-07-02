package storage

import (
	"context"

	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// openMemoryBackend creates an in-memory event store.
// This backend is suitable for testing and local development where
// persistence across restarts is not required.
func openMemoryBackend(ctx context.Context, cfg Config) (eventstore.EventStore, error) {
	return eventstore.NewMemoryEventStore(), nil
}

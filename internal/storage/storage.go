package storage

import (
	"context"
	"fmt"

	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// Storage is the composition point for all Praxis persistence.
// It owns the construction of backend-specific stores and provides
// a single point of configuration and lifecycle management.
//
// Storage does not implement any business logic; it is purely a
// registry that delegates to specialized stores (EventStore, etc.).
type Storage struct {
	Events eventstore.EventStore
	config Config
}

// Open creates a new Storage instance using the provided configuration.
// It constructs the appropriate backend implementations based on Config.Backend.
//
// Supported backends:
//   - "memory": in-memory event store (suitable for testing and development)
//   - "sqlite": SQLite-backed event store (persistent storage)
//
// The caller is responsible for calling Close() to release resources.
func Open(ctx context.Context, cfg Config) (*Storage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("storage.Open: invalid config: %w", err)
	}

	var events eventstore.EventStore
	var err error

	switch cfg.Backend {
	case BackendMemory:
		events, err = openMemoryBackend(ctx, cfg)
	case BackendSQLite:
		events, err = openSQLiteBackend(ctx, cfg)
	default:
		return nil, ErrUnsupportedBackend{Backend: cfg.Backend}
	}

	if err != nil {
		return nil, fmt.Errorf("storage.Open: failed to open %s backend: %w", cfg.Backend, err)
	}

	return &Storage{
		Events: events,
		config: cfg,
	}, nil
}

// MustOpen is like Open but panics on error.
// This is useful in tests and startup code where misconfiguration should be fatal.
func MustOpen(ctx context.Context, cfg Config) *Storage {
	s, err := Open(ctx, cfg)
	if err != nil {
		panic(fmt.Sprintf("storage.MustOpen: %v", err))
	}
	return s
}

// Close releases all resources held by the Storage instance.
// It closes the underlying event store and any other active backends.
// After Close is called, the Storage instance should not be used.
func (s *Storage) Close() error {
	if s.Events != nil {
		if err := s.Events.Close(); err != nil {
			return fmt.Errorf("storage.Close: failed to close event store: %w", err)
		}
	}
	return nil
}

// Config returns the configuration that was used to open this Storage instance.
// The returned Config is a copy and cannot be used to reconfigure the Storage.
func (s *Storage) Config() Config {
	return s.config
}

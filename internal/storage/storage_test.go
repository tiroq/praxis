package storage_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tiroq/praxis/internal/storage"
	"github.com/tiroq/praxis/internal/storage/eventstore"
)

func TestDefaultConfig(t *testing.T) {
	cfg := storage.DefaultConfig()

	if cfg.Backend != storage.BackendSQLite {
		t.Errorf("expected backend %q, got %q", storage.BackendSQLite, cfg.Backend)
	}

	if cfg.SQLitePath != "build/praxis.db" {
		t.Errorf("expected SQLitePath %q, got %q", "build/praxis.db", cfg.SQLitePath)
	}
}

func TestConfigFromEnv(t *testing.T) {
	tests := []struct {
		name        string
		envBackend  string
		envPath     string
		wantBackend string
		wantPath    string
	}{
		{
			name:        "defaults when env not set",
			envBackend:  "",
			envPath:     "",
			wantBackend: storage.BackendSQLite,
			wantPath:    "build/praxis.db",
		},
		{
			name:        "memory backend from env",
			envBackend:  storage.BackendMemory,
			envPath:     "",
			wantBackend: storage.BackendMemory,
			wantPath:    "build/praxis.db",
		},
		{
			name:        "custom sqlite path from env",
			envBackend:  storage.BackendSQLite,
			envPath:     "/tmp/test.db",
			wantBackend: storage.BackendSQLite,
			wantPath:    "/tmp/test.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envBackend != "" {
				os.Setenv("PRAXIS_STORAGE_BACKEND", tt.envBackend)
				defer os.Unsetenv("PRAXIS_STORAGE_BACKEND")
			}
			if tt.envPath != "" {
				os.Setenv("PRAXIS_SQLITE_PATH", tt.envPath)
				defer os.Unsetenv("PRAXIS_SQLITE_PATH")
			}

			cfg := storage.ConfigFromEnv()

			if cfg.Backend != tt.wantBackend {
				t.Errorf("expected backend %q, got %q", tt.wantBackend, cfg.Backend)
			}
			if cfg.SQLitePath != tt.wantPath {
				t.Errorf("expected SQLitePath %q, got %q", tt.wantPath, cfg.SQLitePath)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       storage.Config
		wantError bool
	}{
		{
			name: "valid memory config",
			cfg: storage.Config{
				Backend: storage.BackendMemory,
			},
			wantError: false,
		},
		{
			name: "valid sqlite config",
			cfg: storage.Config{
				Backend:    storage.BackendSQLite,
				SQLitePath: "test.db",
			},
			wantError: false,
		},
		{
			name: "unsupported backend",
			cfg: storage.Config{
				Backend: "postgres",
			},
			wantError: true,
		},
		{
			name: "sqlite with empty path",
			cfg: storage.Config{
				Backend:    storage.BackendSQLite,
				SQLitePath: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestOpenMemoryBackend(t *testing.T) {
	ctx := context.Background()
	cfg := storage.Config{
		Backend: storage.BackendMemory,
	}

	s, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer s.Close()

	if s.Events == nil {
		t.Fatal("expected Events to be non-nil")
	}

	// Verify we can append an event
	event := newTestEvent("test-1")
	if err := s.Events.Append(ctx, event); err != nil {
		t.Errorf("Append() error = %v", err)
	}

	// Verify we can retrieve the event
	retrieved, err := s.Events.Get(ctx, "test-1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if retrieved.ID != "test-1" {
		t.Errorf("expected ID %q, got %q", "test-1", retrieved.ID)
	}
}

func TestOpenSQLiteBackend(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := storage.Config{
		Backend:    storage.BackendSQLite,
		SQLitePath: dbPath,
	}

	s, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer s.Close()

	if s.Events == nil {
		t.Fatal("expected Events to be non-nil")
	}

	// Verify we can append an event
	event := newTestEvent("test-1")
	if err := s.Events.Append(ctx, event); err != nil {
		t.Errorf("Append() error = %v", err)
	}

	// Verify we can retrieve the event
	retrieved, err := s.Events.Get(ctx, "test-1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if retrieved.ID != "test-1" {
		t.Errorf("expected ID %q, got %q", "test-1", retrieved.ID)
	}
}

func TestOpenUnsupportedBackend(t *testing.T) {
	ctx := context.Background()
	cfg := storage.Config{
		Backend: "postgres",
	}

	_, err := storage.Open(ctx, cfg)
	if err == nil {
		t.Fatal("expected error for unsupported backend, got nil")
	}

	var unsupportedErr storage.ErrUnsupportedBackend
	if !errors.As(err, &unsupportedErr) {
		t.Errorf("expected ErrUnsupportedBackend in error chain, got %T: %v", err, err)
	}
}

func TestStorageClose(t *testing.T) {
	ctx := context.Background()
	cfg := storage.Config{
		Backend: storage.BackendMemory,
	}

	s, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := s.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestSQLitePersistence(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persist.db")

	cfg := storage.Config{
		Backend:    storage.BackendSQLite,
		SQLitePath: dbPath,
	}

	// Open storage, write event, close
	s1, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	event := newTestEvent("persist-1")
	if err := s1.Events.Append(ctx, event); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	if err := s1.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Reopen storage and verify event persisted
	s2, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() after reopen error = %v", err)
	}
	defer s2.Close()

	retrieved, err := s2.Events.Get(ctx, "persist-1")
	if err != nil {
		t.Fatalf("Get() after reopen error = %v", err)
	}
	if retrieved.ID != "persist-1" {
		t.Errorf("expected ID %q, got %q", "persist-1", retrieved.ID)
	}
	if retrieved.Type != "test.event" {
		t.Errorf("expected Type %q, got %q", "test.event", retrieved.Type)
	}
}

func TestMustOpen(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		cfg := storage.Config{
			Backend: storage.BackendMemory,
		}

		s := storage.MustOpen(ctx, cfg)
		defer s.Close()

		if s.Events == nil {
			t.Fatal("expected Events to be non-nil")
		}
	})

	t.Run("panic on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got nil")
			}
		}()

		cfg := storage.Config{
			Backend: "invalid",
		}
		storage.MustOpen(ctx, cfg)
	})
}

func TestStorageConfig(t *testing.T) {
	ctx := context.Background()
	cfg := storage.Config{
		Backend:    storage.BackendMemory,
		SQLitePath: "test.db",
	}

	s, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer s.Close()

	retrieved := s.Config()
	if retrieved.Backend != cfg.Backend {
		t.Errorf("expected Backend %q, got %q", cfg.Backend, retrieved.Backend)
	}
	if retrieved.SQLitePath != cfg.SQLitePath {
		t.Errorf("expected SQLitePath %q, got %q", cfg.SQLitePath, retrieved.SQLitePath)
	}
}

// Helper functions

func newTestEvent(id string) eventstore.EventRecord {
	payload := map[string]string{"message": "test"}
	payloadJSON, _ := json.Marshal(payload)

	occurredAt, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")

	return eventstore.EventRecord{
		ID:         id,
		Type:       "test.event",
		Source:     "test",
		SubjectID:  "test-subject",
		OccurredAt: occurredAt,
		Payload:    payloadJSON,
		Metadata:   map[string]string{"test": "true"},
	}
}

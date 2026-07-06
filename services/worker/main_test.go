package main

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage"
	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// TestBuildKernel_WithoutOptions verifies that buildKernel creates a functional
// kernel when no options are provided.
func TestBuildKernel_WithoutOptions(t *testing.T) {
	k := buildKernel()
	if k == nil {
		t.Fatal("buildKernel() returned nil")
	}

	// Verify that the kernel can execute the pipeline without event recording
	ctx := context.Background()
	event := kernel.Event{
		ID:     "test-event-1",
		Source: "test",
		Text:   "urgent review required",
	}

	result, err := k.Run(ctx, event)
	if err != nil {
		t.Fatalf("kernel.Run() failed: %v", err)
	}

	if result.EventID != event.ID {
		t.Errorf("expected EventID=%s, got %s", event.ID, result.EventID)
	}
}

// TestBuildKernel_WithEventRecorder verifies that buildKernel accepts
// WithEventRecorder option and creates a kernel that records events.
func TestBuildKernel_WithEventRecorder(t *testing.T) {
	ctx := context.Background()

	// Create an in-memory storage for testing
	cfg := storage.Config{
		Backend: storage.BackendMemory,
	}
	store, err := storage.Open(ctx, cfg)
	if err != nil {
		t.Fatalf("storage.Open() failed: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("store.Close() failed: %v", err)
		}
	}()

	// Build kernel with EventRecorder
	recorder := newEventRecorderAdapter(store.Events)
	k := buildKernel(kernel.WithEventRecorder(recorder))
	if k == nil {
		t.Fatal("buildKernel() returned nil")
	}

	// Execute the pipeline
	event := kernel.Event{
		ID:     "test-event-2",
		Source: "test",
		Text:   "купить билеты в Шанхай",
	}

	result, err := k.Run(ctx, event)
	if err != nil {
		t.Fatalf("kernel.Run() failed: %v", err)
	}

	if result.EventID != event.ID {
		t.Errorf("expected EventID=%s, got %s", event.ID, result.EventID)
	}

	// Verify that an event was recorded
	events, err := store.Events.List(ctx, eventstore.ListFilter{})
	if err != nil {
		t.Fatalf("store.Events.List() failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("expected 1 recorded event, got %d", len(events))
	}

	recordedEvent := events[0]
	if recordedEvent.Type != "kernel.pipeline.completed" {
		t.Errorf("expected event type 'kernel.pipeline.completed', got '%s'", recordedEvent.Type)
	}
	if recordedEvent.Source != "kernel" {
		t.Errorf("expected event source 'kernel', got '%s'", recordedEvent.Source)
	}
	if recordedEvent.CausationID != event.ID {
		t.Errorf("expected causation_id=%s, got %s", event.ID, recordedEvent.CausationID)
	}
}

// TestTryOpenStorage_Success verifies that tryOpenStorage successfully opens
// an in-memory storage backend.
func TestTryOpenStorage_Success(t *testing.T) {
	// Save original env vars and restore after test
	origBackend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	defer func() {
		if origBackend != "" {
			_ = os.Setenv("PRAXIS_STORAGE_BACKEND", origBackend)
		} else {
			_ = os.Unsetenv("PRAXIS_STORAGE_BACKEND")
		}
	}()

	_ = os.Setenv("PRAXIS_STORAGE_BACKEND", "memory")

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	store := tryOpenStorage(ctx, logger)
	if store == nil {
		t.Fatal("tryOpenStorage() returned nil for valid memory backend")
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("store.Close() failed: %v", err)
		}
	}()

	if store.Events == nil {
		t.Error("store.Events is nil")
	}
}

// TestTryOpenStorage_InvalidBackend verifies that tryOpenStorage returns nil
// (non-fatal) when an invalid backend is configured.
func TestTryOpenStorage_InvalidBackend(t *testing.T) {
	// Save original env vars and restore after test
	origBackend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	defer func() {
		if origBackend != "" {
			_ = os.Setenv("PRAXIS_STORAGE_BACKEND", origBackend)
		} else {
			_ = os.Unsetenv("PRAXIS_STORAGE_BACKEND")
		}
	}()

	_ = os.Setenv("PRAXIS_STORAGE_BACKEND", "invalid-backend")

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	store := tryOpenStorage(ctx, logger)
	if store != nil {
		defer func() {
			if err := store.Close(); err != nil {
				t.Errorf("store.Close() failed: %v", err)
			}
		}()
		t.Error("tryOpenStorage() should return nil for invalid backend")
	}
}

// TestTryOpenStorage_SQLite verifies that tryOpenStorage successfully opens
// a SQLite backend with a temporary database file.
func TestTryOpenStorage_SQLite(t *testing.T) {
	// Save original env vars and restore after test
	origBackend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	origPath := os.Getenv("PRAXIS_SQLITE_PATH")
	defer func() {
		if origBackend != "" {
			_ = os.Setenv("PRAXIS_STORAGE_BACKEND", origBackend)
		} else {
			_ = os.Unsetenv("PRAXIS_STORAGE_BACKEND")
		}
		if origPath != "" {
			_ = os.Setenv("PRAXIS_SQLITE_PATH", origPath)
		} else {
			_ = os.Unsetenv("PRAXIS_SQLITE_PATH")
		}
	}()

	// Create a temporary directory for the test database
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test-worker.db"

	_ = os.Setenv("PRAXIS_STORAGE_BACKEND", "sqlite")
	_ = os.Setenv("PRAXIS_SQLITE_PATH", dbPath)

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	store := tryOpenStorage(ctx, logger)
	if store == nil {
		t.Fatal("tryOpenStorage() returned nil for valid SQLite backend")
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("store.Close() failed: %v", err)
		}
	}()

	if store.Events == nil {
		t.Error("store.Events is nil")
	}
}

// TestTryOpenConversationStore_SQLite verifies that the worker can open
// the SQLite conversation projection store when SQLite backend is configured.
func TestTryOpenConversationStore_SQLite(t *testing.T) {
	origBackend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	origPath := os.Getenv("PRAXIS_SQLITE_PATH")
	defer func() {
		if origBackend != "" {
			_ = os.Setenv("PRAXIS_STORAGE_BACKEND", origBackend)
		} else {
			_ = os.Unsetenv("PRAXIS_STORAGE_BACKEND")
		}
		if origPath != "" {
			_ = os.Setenv("PRAXIS_SQLITE_PATH", origPath)
		} else {
			_ = os.Unsetenv("PRAXIS_SQLITE_PATH")
		}
	}()

	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test-conversation.db"

	_ = os.Setenv("PRAXIS_STORAGE_BACKEND", "sqlite")
	_ = os.Setenv("PRAXIS_SQLITE_PATH", dbPath)

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	store := tryOpenConversationStore(ctx, logger)
	if store == nil {
		t.Fatal("tryOpenConversationStore() returned nil for valid SQLite backend")
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("store.Close() failed: %v", err)
		}
	}()
}

// TestTryOpenConversationStore_NonSQLite verifies that conversation projection
// initialization is skipped for non-SQLite backends.
func TestTryOpenConversationStore_NonSQLite(t *testing.T) {
	origBackend := os.Getenv("PRAXIS_STORAGE_BACKEND")
	defer func() {
		if origBackend != "" {
			_ = os.Setenv("PRAXIS_STORAGE_BACKEND", origBackend)
		} else {
			_ = os.Unsetenv("PRAXIS_STORAGE_BACKEND")
		}
	}()

	_ = os.Setenv("PRAXIS_STORAGE_BACKEND", "memory")

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	store := tryOpenConversationStore(ctx, logger)
	if store != nil {
		defer func() {
			if err := store.Close(); err != nil {
				t.Errorf("store.Close() failed: %v", err)
			}
		}()
		t.Fatal("tryOpenConversationStore() should return nil for non-SQLite backend")
	}
}

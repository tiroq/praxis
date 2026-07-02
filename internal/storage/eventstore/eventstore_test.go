package eventstore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/tiroq/praxis/internal/storage/eventstore"
	"github.com/tiroq/praxis/internal/storage/sqlite"
)

// storeFactory creates a new EventStore for testing.
type storeFactory func(t *testing.T) eventstore.EventStore

// makeMemoryStore creates an in-memory event store.
func makeMemoryStore(t *testing.T) eventstore.EventStore {
	store := eventstore.NewMemoryEventStore()
	t.Cleanup(func() {
		store.Close()
	})
	return store
}

// makeSQLiteStore creates a SQLite event store using a temporary file.
func makeSQLiteStore(t *testing.T) eventstore.EventStore {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := sqlite.OpenEventStore(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("failed to open SQLite store: %v", err)
	}

	t.Cleanup(func() {
		store.Close()
	})

	return store
}

// makeSQLiteMemoryStore creates a SQLite event store using an in-memory database.
func makeSQLiteMemoryStore(t *testing.T) eventstore.EventStore {
	store, err := sqlite.OpenEventStore(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("failed to open SQLite memory store: %v", err)
	}

	t.Cleanup(func() {
		store.Close()
	})

	return store
}

// testEventRecord creates a valid test event record.
func testEventRecord(id string) eventstore.EventRecord {
	payload, _ := json.Marshal(map[string]interface{}{
		"message": "test message",
		"value":   42,
	})

	return eventstore.EventRecord{
		ID:            id,
		Type:          "test.event",
		Source:        "test-source",
		SubjectID:     "subject-1",
		CorrelationID: "corr-1",
		CausationID:   "cause-1",
		TraceID:       "trace-1",
		OccurredAt:    time.Now().UTC(),
		Payload:       payload,
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}
}

// runContractTests runs the contract tests against a given store factory.
func runContractTests(t *testing.T, factory storeFactory) {
	t.Run("AppendAndGet", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event := testEventRecord("event-1")

		// Append the event
		if err := store.Append(ctx, event); err != nil {
			t.Fatalf("Append failed: %v", err)
		}

		// Get the event
		retrieved, err := store.Get(ctx, "event-1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		// Verify fields
		if retrieved.ID != event.ID {
			t.Errorf("ID mismatch: got %q, want %q", retrieved.ID, event.ID)
		}
		if retrieved.Type != event.Type {
			t.Errorf("Type mismatch: got %q, want %q", retrieved.Type, event.Type)
		}
		if retrieved.Source != event.Source {
			t.Errorf("Source mismatch: got %q, want %q", retrieved.Source, event.Source)
		}
		if retrieved.SubjectID != event.SubjectID {
			t.Errorf("SubjectID mismatch: got %q, want %q", retrieved.SubjectID, event.SubjectID)
		}
		if !retrieved.OccurredAt.Equal(event.OccurredAt) {
			t.Errorf("OccurredAt mismatch: got %v, want %v", retrieved.OccurredAt, event.OccurredAt)
		}
		if string(retrieved.Payload) != string(event.Payload) {
			t.Errorf("Payload mismatch: got %s, want %s", retrieved.Payload, event.Payload)
		}
	})

	t.Run("DuplicateAppendRejected", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event := testEventRecord("event-dup")

		// First append should succeed
		if err := store.Append(ctx, event); err != nil {
			t.Fatalf("First Append failed: %v", err)
		}

		// Second append with same ID should fail
		err := store.Append(ctx, event)
		if err == nil {
			t.Fatal("Expected duplicate event error, got nil")
		}

		dupErr, ok := err.(eventstore.ErrDuplicateEvent)
		if !ok {
			t.Fatalf("Expected ErrDuplicateEvent, got %T: %v", err, err)
		}
		if dupErr.ID != "event-dup" {
			t.Errorf("Expected ID %q in error, got %q", "event-dup", dupErr.ID)
		}
	})

	t.Run("GetMissingReturnsNotFound", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		_, err := store.Get(ctx, "nonexistent")
		if err == nil {
			t.Fatal("Expected not found error, got nil")
		}

		notFoundErr, ok := err.(eventstore.ErrEventNotFound)
		if !ok {
			t.Fatalf("Expected ErrEventNotFound, got %T: %v", err, err)
		}
		if notFoundErr.ID != "nonexistent" {
			t.Errorf("Expected ID %q in error, got %q", "nonexistent", notFoundErr.ID)
		}
	})

	t.Run("PayloadValidation", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event := testEventRecord("event-invalid-payload")
		event.Payload = []byte("not valid json")

		err := store.Append(ctx, event)
		if err == nil {
			t.Fatal("Expected invalid JSON error, got nil")
		}
	})

	t.Run("RequiredFieldValidation", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		tests := []struct {
			name     string
			mutate   func(*eventstore.EventRecord)
			expected string
		}{
			{
				name:     "MissingID",
				mutate:   func(e *eventstore.EventRecord) { e.ID = "" },
				expected: "id",
			},
			{
				name:     "MissingType",
				mutate:   func(e *eventstore.EventRecord) { e.Type = "" },
				expected: "type",
			},
			{
				name:     "MissingSource",
				mutate:   func(e *eventstore.EventRecord) { e.Source = "" },
				expected: "source",
			},
			{
				name:     "MissingSubjectID",
				mutate:   func(e *eventstore.EventRecord) { e.SubjectID = "" },
				expected: "subject_id",
			},
			{
				name:     "MissingOccurredAt",
				mutate:   func(e *eventstore.EventRecord) { e.OccurredAt = time.Time{} },
				expected: "occurred_at",
			},
			{
				name:     "MissingPayload",
				mutate:   func(e *eventstore.EventRecord) { e.Payload = nil },
				expected: "payload",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				event := testEventRecord(fmt.Sprintf("event-%s", tt.name))
				tt.mutate(&event)

				err := store.Append(ctx, event)
				if err == nil {
					t.Fatal("Expected validation error, got nil")
				}

				missingErr, ok := err.(eventstore.ErrMissingField)
				if !ok {
					t.Fatalf("Expected ErrMissingField, got %T: %v", err, err)
				}
				if string(missingErr) != tt.expected {
					t.Errorf("Expected field %q in error, got %q", tt.expected, string(missingErr))
				}
			})
		}
	})

	t.Run("NilMetadataRoundTrips", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event := testEventRecord("event-nil-metadata")
		event.Metadata = nil

		if err := store.Append(ctx, event); err != nil {
			t.Fatalf("Append failed: %v", err)
		}

		retrieved, err := store.Get(ctx, "event-nil-metadata")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.Metadata == nil {
			t.Error("Expected non-nil Metadata, got nil")
		}
		if len(retrieved.Metadata) != 0 {
			t.Errorf("Expected empty Metadata, got %d entries", len(retrieved.Metadata))
		}
	})

	t.Run("ListByType", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		// Append events with different types
		event1 := testEventRecord("list-type-1")
		event1.Type = "type.a"
		event1.OccurredAt = time.Now().UTC().Add(-2 * time.Minute)

		event2 := testEventRecord("list-type-2")
		event2.Type = "type.b"
		event2.OccurredAt = time.Now().UTC().Add(-1 * time.Minute)

		event3 := testEventRecord("list-type-3")
		event3.Type = "type.a"
		event3.OccurredAt = time.Now().UTC()

		for _, e := range []eventstore.EventRecord{event1, event2, event3} {
			if err := store.Append(ctx, e); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		// List events of type.a
		events, err := store.List(ctx, eventstore.ListFilter{Type: "type.a"})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 2 {
			t.Fatalf("Expected 2 events, got %d", len(events))
		}

		if events[0].ID != "list-type-1" || events[1].ID != "list-type-3" {
			t.Errorf("Expected IDs [list-type-1, list-type-3], got [%s, %s]", events[0].ID, events[1].ID)
		}
	})

	t.Run("ListBySource", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event1 := testEventRecord("list-source-1")
		event1.Source = "source-x"
		event1.OccurredAt = time.Now().UTC()

		event2 := testEventRecord("list-source-2")
		event2.Source = "source-y"
		event2.OccurredAt = time.Now().UTC()

		for _, e := range []eventstore.EventRecord{event1, event2} {
			if err := store.Append(ctx, e); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		events, err := store.List(ctx, eventstore.ListFilter{Source: "source-x"})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}
		if events[0].ID != "list-source-1" {
			t.Errorf("Expected ID list-source-1, got %s", events[0].ID)
		}
	})

	t.Run("ListBySubjectID", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event1 := testEventRecord("list-subject-1")
		event1.SubjectID = "subj-1"
		event1.OccurredAt = time.Now().UTC()

		event2 := testEventRecord("list-subject-2")
		event2.SubjectID = "subj-2"
		event2.OccurredAt = time.Now().UTC()

		for _, e := range []eventstore.EventRecord{event1, event2} {
			if err := store.Append(ctx, e); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		events, err := store.List(ctx, eventstore.ListFilter{SubjectID: "subj-1"})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}
		if events[0].ID != "list-subject-1" {
			t.Errorf("Expected ID list-subject-1, got %s", events[0].ID)
		}
	})

	t.Run("ListByCorrelationID", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		event1 := testEventRecord("list-corr-1")
		event1.CorrelationID = "corr-x"
		event1.OccurredAt = time.Now().UTC()

		event2 := testEventRecord("list-corr-2")
		event2.CorrelationID = "corr-y"
		event2.OccurredAt = time.Now().UTC()

		for _, e := range []eventstore.EventRecord{event1, event2} {
			if err := store.Append(ctx, e); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		events, err := store.List(ctx, eventstore.ListFilter{CorrelationID: "corr-x"})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(events))
		}
		if events[0].ID != "list-corr-1" {
			t.Errorf("Expected ID list-corr-1, got %s", events[0].ID)
		}
	})

	t.Run("ListOrderDeterministic", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		baseTime := time.Now().UTC().Truncate(time.Second)

		// Create events with same occurred_at but different IDs
		event1 := testEventRecord("event-c")
		event1.OccurredAt = baseTime

		event2 := testEventRecord("event-a")
		event2.OccurredAt = baseTime

		event3 := testEventRecord("event-b")
		event3.OccurredAt = baseTime

		// Append in non-alphabetical order
		for _, e := range []eventstore.EventRecord{event1, event2, event3} {
			if err := store.Append(ctx, e); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		// List should return in ID order when occurred_at is equal
		events, err := store.List(ctx, eventstore.ListFilter{})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		if len(events) < 3 {
			t.Fatalf("Expected at least 3 events, got %d", len(events))
		}

		// Find our test events in the results
		var testEvents []eventstore.EventRecord
		for _, e := range events {
			if e.ID == "event-a" || e.ID == "event-b" || e.ID == "event-c" {
				testEvents = append(testEvents, e)
			}
		}

		if len(testEvents) != 3 {
			t.Fatalf("Expected 3 test events, got %d", len(testEvents))
		}

		// Should be in alphabetical order by ID
		if testEvents[0].ID != "event-a" || testEvents[1].ID != "event-b" || testEvents[2].ID != "event-c" {
			t.Errorf("Expected IDs [event-a, event-b, event-c], got [%s, %s, %s]",
				testEvents[0].ID, testEvents[1].ID, testEvents[2].ID)
		}
	})

	t.Run("ListLimitAndOffset", func(t *testing.T) {
		store := factory(t)
		ctx := context.Background()

		baseTime := time.Now().UTC().Truncate(time.Second)

		// Append 5 events
		for i := 0; i < 5; i++ {
			event := testEventRecord(fmt.Sprintf("event-limit-%d", i))
			event.OccurredAt = baseTime.Add(time.Duration(i) * time.Second)
			if err := store.Append(ctx, event); err != nil {
				t.Fatalf("Append failed: %v", err)
			}
		}

		// Test limit
		events, err := store.List(ctx, eventstore.ListFilter{Limit: 2})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}
		if len(events) < 2 {
			t.Errorf("Expected at least 2 events with limit, got %d", len(events))
		}

		// Test offset
		events, err = store.List(ctx, eventstore.ListFilter{Limit: 10, Offset: 2})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		// Count our test events
		count := 0
		for _, e := range events {
			if e.Type == "test.event" && e.Source == "test-source" {
				count++
			}
		}
		if count < 3 {
			t.Errorf("Expected at least 3 events with offset, got %d", count)
		}
	})
}

// TestMemoryEventStore runs contract tests against the memory implementation.
func TestMemoryEventStore(t *testing.T) {
	runContractTests(t, makeMemoryStore)

	t.Run("ConcurrentAppends", func(t *testing.T) {
		store := makeMemoryStore(t)
		ctx := context.Background()

		const numGoroutines = 10
		const eventsPerGoroutine = 10

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for g := 0; g < numGoroutines; g++ {
			go func(goroutineID int) {
				defer wg.Done()

				for i := 0; i < eventsPerGoroutine; i++ {
					event := testEventRecord(fmt.Sprintf("concurrent-%d-%d", goroutineID, i))
					if err := store.Append(ctx, event); err != nil {
						t.Errorf("Concurrent Append failed: %v", err)
					}
				}
			}(g)
		}

		wg.Wait()

		// Verify all events were stored
		events, err := store.List(ctx, eventstore.ListFilter{})
		if err != nil {
			t.Fatalf("List failed: %v", err)
		}

		expectedCount := numGoroutines * eventsPerGoroutine
		if len(events) < expectedCount {
			t.Errorf("Expected at least %d events, got %d", expectedCount, len(events))
		}
	})
}

// TestSQLiteEventStore runs contract tests against the SQLite implementation.
func TestSQLiteEventStore(t *testing.T) {
	t.Run("MemoryDB", func(t *testing.T) {
		runContractTests(t, makeSQLiteMemoryStore)
	})

	t.Run("FileDB", func(t *testing.T) {
		runContractTests(t, makeSQLiteStore)
	})

	t.Run("PersistenceAcrossReopen", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "persist.db")
		ctx := context.Background()

		// First session: create and append
		store1, err := sqlite.OpenEventStore(ctx, dbPath)
		if err != nil {
			t.Fatalf("Failed to open store: %v", err)
		}

		event := testEventRecord("persist-1")
		if err := store1.Append(ctx, event); err != nil {
			t.Fatalf("Append failed: %v", err)
		}

		if err := store1.Close(); err != nil {
			t.Fatalf("Close failed: %v", err)
		}

		// Second session: reopen and verify
		store2, err := sqlite.OpenEventStore(ctx, dbPath)
		if err != nil {
			t.Fatalf("Failed to reopen store: %v", err)
		}
		defer store2.Close()

		retrieved, err := store2.Get(ctx, "persist-1")
		if err != nil {
			t.Fatalf("Get failed after reopen: %v", err)
		}

		if retrieved.ID != "persist-1" {
			t.Errorf("Expected ID persist-1, got %s", retrieved.ID)
		}

		// Verify we can't append duplicate after reopen
		err = store2.Append(ctx, event)
		if err == nil {
			t.Fatal("Expected duplicate error after reopen, got nil")
		}
	})

	t.Run("MultipleStoresIndependent", func(t *testing.T) {
		tmpDir := t.TempDir()
		ctx := context.Background()

		store1, err := sqlite.OpenEventStore(ctx, filepath.Join(tmpDir, "store1.db"))
		if err != nil {
			t.Fatalf("Failed to open store1: %v", err)
		}
		defer store1.Close()

		store2, err := sqlite.OpenEventStore(ctx, filepath.Join(tmpDir, "store2.db"))
		if err != nil {
			t.Fatalf("Failed to open store2: %v", err)
		}
		defer store2.Close()

		event := testEventRecord("independent-1")
		if err := store1.Append(ctx, event); err != nil {
			t.Fatalf("Append to store1 failed: %v", err)
		}

		// Event should not exist in store2
		_, err = store2.Get(ctx, "independent-1")
		if err == nil {
			t.Fatal("Expected not found error in store2, got nil")
		}

		if _, ok := err.(eventstore.ErrEventNotFound); !ok {
			t.Fatalf("Expected ErrEventNotFound, got %T: %v", err, err)
		}
	})
}

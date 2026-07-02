package eventstore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/tiroq/praxis/internal/storage/eventstore"
	"github.com/tiroq/praxis/internal/storage/sqlite"
)

// benchEventRecord returns a realistic EventRecord for benchmarks.
// The payload intentionally mirrors a real-world event shape.
func benchEventRecord(id string) eventstore.EventRecord {
	payload, _ := json.Marshal(map[string]interface{}{
		"user_id":    "user-123",
		"session_id": "sess-456abc",
		"action":     "login",
		"ip":         "192.168.1.42",
		"user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"timestamp":  time.Now().Unix(),
		"success":    true,
	})
	return eventstore.EventRecord{
		ID:            id,
		Type:          "user.login",
		Source:        "auth-service",
		SubjectID:     "user-123",
		CorrelationID: "corr-bench-001",
		CausationID:   "cause-bench-001",
		TraceID:       "trace-bench-001",
		OccurredAt:    time.Now().UTC(),
		Payload:       payload,
		Metadata: map[string]string{
			"region":  "us-east-1",
			"version": "1.0.0",
		},
	}
}

// ── Memory benchmarks ────────────────────────────────────────────────────────

func BenchmarkAppendMemory(b *testing.B) {
	store := eventstore.NewMemoryEventStore()
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-mem-append-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Append failed: %v", err)
		}
	}
}

func BenchmarkGetMemory(b *testing.B) {
	store := eventstore.NewMemoryEventStore()
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	const preload = 100
	for i := 0; i < preload; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-mem-get-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Setup Append failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("bench-mem-get-%d", i%preload)
		if _, err := store.Get(ctx, id); err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkListMemory(b *testing.B) {
	store := eventstore.NewMemoryEventStore()
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-mem-list-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Setup Append failed: %v", err)
		}
	}

	filter := eventstore.ListFilter{Type: "user.login", Limit: 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.List(ctx, filter); err != nil {
			b.Fatalf("List failed: %v", err)
		}
	}
}

// ── SQLite benchmarks ────────────────────────────────────────────────────────
// All SQLite benchmarks use ":memory:" to avoid disk I/O noise.

func BenchmarkAppendSQLite(b *testing.B) {
	store, err := sqlite.OpenEventStore(context.Background(), ":memory:")
	if err != nil {
		b.Fatalf("OpenEventStore failed: %v", err)
	}
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-sqlite-append-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Append failed: %v", err)
		}
	}
}

func BenchmarkGetSQLite(b *testing.B) {
	store, err := sqlite.OpenEventStore(context.Background(), ":memory:")
	if err != nil {
		b.Fatalf("OpenEventStore failed: %v", err)
	}
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	const preload = 100
	for i := 0; i < preload; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-sqlite-get-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Setup Append failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := fmt.Sprintf("bench-sqlite-get-%d", i%preload)
		if _, err := store.Get(ctx, id); err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkListSQLite(b *testing.B) {
	store, err := sqlite.OpenEventStore(context.Background(), ":memory:")
	if err != nil {
		b.Fatalf("OpenEventStore failed: %v", err)
	}
	defer func() { _ = store.Close() }()
	ctx := context.Background()

	for i := 0; i < 100; i++ {
		event := benchEventRecord(fmt.Sprintf("bench-sqlite-list-%d", i))
		if err := store.Append(ctx, event); err != nil {
			b.Fatalf("Setup Append failed: %v", err)
		}
	}

	filter := eventstore.ListFilter{Type: "user.login", Limit: 20}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := store.List(ctx, filter); err != nil {
			b.Fatalf("List failed: %v", err)
		}
	}
}

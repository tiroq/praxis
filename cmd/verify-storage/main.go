// verify-storage is a smoke test helper that verifies the EventStore
// contains expected kernel pipeline events after a smoke test run.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tiroq/praxis/internal/storage"
	"github.com/tiroq/praxis/internal/storage/eventstore"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <db_path>\n", os.Args[0])
		os.Exit(1)
	}

	dbPath := os.Args[1]

	ctx := context.Background()
	cfg := storage.Config{
		Backend:    storage.BackendSQLite,
		SQLitePath: dbPath,
	}

	store, err := storage.Open(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open storage: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := store.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close storage: %v\n", err)
		}
	}()

	// Query for kernel.pipeline.completed events
	filter := eventstore.ListFilter{
		Type: "kernel.pipeline.completed",
	}

	events, err := store.Events.List(ctx, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list events: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Fprintf(os.Stderr, "no kernel.pipeline.completed events found\n")
		os.Exit(1)
	}

	// Report the results
	result := map[string]any{
		"event_count":  len(events),
		"event_type":   "kernel.pipeline.completed",
		"events":       make([]map[string]any, 0, len(events)),
		"verification": "passed",
	}

	for _, event := range events {
		eventInfo := map[string]any{
			"id":             event.ID,
			"type":           event.Type,
			"source":         event.Source,
			"subject_id":     event.SubjectID,
			"correlation_id": event.CorrelationID,
			"causation_id":   event.CausationID,
			"occurred_at":    event.OccurredAt.Format("2006-01-02T15:04:05Z07:00"),
			"metadata":       event.Metadata,
		}
		result["events"] = append(result["events"].([]map[string]any), eventInfo)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal result: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

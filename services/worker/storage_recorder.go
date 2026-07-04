package main

import (
	"context"

	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// eventRecorderAdapter adapts an eventstore.EventStore to the kernel's EventRecorder interface.
// This allows the kernel to record events without depending on storage implementation details.
//
// This adapter lives in the worker service (composition root) rather than internal/storage
// to preserve architectural boundaries: storage must not import kernel.
type eventRecorderAdapter struct {
	store eventstore.EventStore
}

// newEventRecorderAdapter creates a new adapter that wraps an EventStore.
func newEventRecorderAdapter(store eventstore.EventStore) *eventRecorderAdapter {
	return &eventRecorderAdapter{store: store}
}

// Append converts a kernel event record to a storage event record and appends it.
func (a *eventRecorderAdapter) Append(ctx context.Context, event kernel.EventRecord) error {
	// Convert kernel event record to storage event record
	storageEvent := eventstore.EventRecord{
		ID:            event.ID,
		CorrelationID: event.CorrelationID,
		CausationID:   event.CausationID,
		TraceID:       event.TraceID,
		Type:          event.Type,
		Source:        event.Source,
		SubjectID:     event.SubjectID,
		OccurredAt:    event.OccurredAt,
		Payload:       event.Payload,
		Metadata:      event.Metadata,
		CreatedAt:     event.CreatedAt,
	}

	return a.store.Append(ctx, storageEvent)
}

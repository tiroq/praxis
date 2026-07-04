package storage

import (
	"context"

	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage/eventstore"
)

// EventRecorderAdapter adapts an eventstore.EventStore to the kernel's EventRecorder interface.
// This allows the kernel to record events without depending on storage implementation details.
type EventRecorderAdapter struct {
	store eventstore.EventStore
}

// NewEventRecorderAdapter creates a new adapter that wraps an EventStore.
func NewEventRecorderAdapter(store eventstore.EventStore) *EventRecorderAdapter {
	return &EventRecorderAdapter{store: store}
}

// Append converts a kernel event record to a storage event record and appends it.
func (a *EventRecorderAdapter) Append(ctx context.Context, event kernel.EventRecord) error {
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

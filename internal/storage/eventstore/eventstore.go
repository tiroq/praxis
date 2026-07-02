package eventstore

import "context"

// ListFilter specifies criteria for filtering events.
type ListFilter struct {
	Type          string // filter by event type
	Source        string // filter by event source
	SubjectID     string // filter by subject ID
	CorrelationID string // filter by correlation ID
	Limit         int    // maximum number of events to return (0 = default)
	Offset        int    // number of events to skip
}

// EventStore defines the interface for persistent event storage.
// All implementations must be append-only and never modify or delete events.
type EventStore interface {
	// Append stores a new event.
	// Returns ErrDuplicateEvent if an event with the same ID already exists.
	// Returns ErrMissingField if required fields are missing.
	// Returns ErrInvalidJSON if the payload is not valid JSON.
	Append(ctx context.Context, event EventRecord) error

	// Get retrieves an event by ID.
	// Returns ErrEventNotFound if the event does not exist.
	Get(ctx context.Context, id string) (EventRecord, error)

	// List retrieves events matching the filter criteria.
	// Results are ordered deterministically by occurred_at ascending, then ID ascending.
	// If filter.Limit is 0, a safe default is used.
	List(ctx context.Context, filter ListFilter) ([]EventRecord, error)

	// Close releases any resources held by the store.
	Close() error
}

package eventstore

import (
	"context"
	"sort"
	"sync"
	"time"
)

// MemoryEventStore is an in-memory implementation of EventStore.
// It is thread-safe and suitable for testing and local development.
// Events are stored in an append-only fashion.
type MemoryEventStore struct {
	mu     sync.RWMutex
	events map[string]EventRecord // indexed by ID
	order  []string               // ordered list of event IDs
}

// NewMemoryEventStore creates a new in-memory event store.
func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		events: make(map[string]EventRecord),
		order:  make([]string, 0),
	}
}

// Append stores a new event in memory.
func (m *MemoryEventStore) Append(ctx context.Context, event EventRecord) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for duplicate ID
	if _, exists := m.events[event.ID]; exists {
		return ErrDuplicateEvent{ID: event.ID}
	}

	// Set CreatedAt if not already set
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	// Ensure metadata is not nil for consistent round-tripping
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}

	// Store a clone to ensure immutability
	clone := event.Clone()
	m.events[event.ID] = clone
	m.order = append(m.order, event.ID)

	return nil
}

// Get retrieves an event by ID.
func (m *MemoryEventStore) Get(ctx context.Context, id string) (EventRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	event, exists := m.events[id]
	if !exists {
		return EventRecord{}, ErrEventNotFound{ID: id}
	}

	// Return a clone to ensure immutability
	return event.Clone(), nil
}

// List retrieves events matching the filter criteria.
// Results are ordered by occurred_at ascending, then ID ascending.
func (m *MemoryEventStore) List(ctx context.Context, filter ListFilter) ([]EventRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Collect matching events
	var matches []EventRecord
	for _, id := range m.order {
		event := m.events[id]

		// Apply filters
		if filter.Type != "" && event.Type != filter.Type {
			continue
		}
		if filter.Source != "" && event.Source != filter.Source {
			continue
		}
		if filter.SubjectID != "" && event.SubjectID != filter.SubjectID {
			continue
		}
		if filter.CorrelationID != "" && event.CorrelationID != filter.CorrelationID {
			continue
		}

		matches = append(matches, event.Clone())
	}

	// Sort by occurred_at, then by ID for deterministic ordering
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].OccurredAt.Equal(matches[j].OccurredAt) {
			return matches[i].ID < matches[j].ID
		}
		return matches[i].OccurredAt.Before(matches[j].OccurredAt)
	})

	// Apply limit and offset
	limit := filter.Limit
	if limit == 0 {
		limit = 100 // safe default
	}

	start := filter.Offset
	if start > len(matches) {
		return []EventRecord{}, nil
	}

	end := start + limit
	if end > len(matches) {
		end = len(matches)
	}

	return matches[start:end], nil
}

// Close releases resources. For memory store, this is a no-op.
func (m *MemoryEventStore) Close() error {
	return nil
}

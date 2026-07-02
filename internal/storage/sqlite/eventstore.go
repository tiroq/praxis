package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/tiroq/praxis/internal/storage/eventstore"
	_ "modernc.org/sqlite" // SQLite driver
)

// EventStore is a SQLite-backed implementation of eventstore.EventStore.
type EventStore struct {
	db *sql.DB
}

// OpenEventStore opens or creates a SQLite event store at the specified path.
// Use ":memory:" for an in-memory database (useful for tests).
func OpenEventStore(ctx context.Context, path string) (*EventStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}

	// Create schema
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &EventStore{db: db}, nil
}

// Append stores a new event in SQLite.
func (s *EventStore) Append(ctx context.Context, event eventstore.EventRecord) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		return err
	}

	// Set CreatedAt if not already set
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	// Ensure metadata is not nil for consistent round-tripping
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	// Insert the event
	query := `
		INSERT INTO events (
			id, type, source, subject_id, correlation_id, causation_id, trace_id,
			occurred_at, payload, metadata, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		event.ID,
		event.Type,
		event.Source,
		event.SubjectID,
		event.CorrelationID,
		event.CausationID,
		event.TraceID,
		event.OccurredAt.Format(time.RFC3339Nano),
		string(event.Payload),
		string(metadataJSON),
		event.CreatedAt.Format(time.RFC3339Nano),
	)

	if err != nil {
		// Check for duplicate ID (UNIQUE constraint violation)
		if isSQLiteConstraintError(err) {
			return eventstore.ErrDuplicateEvent{ID: event.ID}
		}
		return err
	}

	return nil
}

// Get retrieves an event by ID.
func (s *EventStore) Get(ctx context.Context, id string) (eventstore.EventRecord, error) {
	query := `
		SELECT id, type, source, subject_id, correlation_id, causation_id, trace_id,
			occurred_at, payload, metadata, created_at
		FROM events
		WHERE id = ?
	`

	var event eventstore.EventRecord
	var occurredAtStr, createdAtStr string
	var payloadStr, metadataJSON string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.Type,
		&event.Source,
		&event.SubjectID,
		&event.CorrelationID,
		&event.CausationID,
		&event.TraceID,
		&occurredAtStr,
		&payloadStr,
		&metadataJSON,
		&createdAtStr,
	)

	if err == sql.ErrNoRows {
		return eventstore.EventRecord{}, eventstore.ErrEventNotFound{ID: id}
	}
	if err != nil {
		return eventstore.EventRecord{}, err
	}

	// Parse timestamps
	event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAtStr)
	if err != nil {
		return eventstore.EventRecord{}, err
	}

	event.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
	if err != nil {
		return eventstore.EventRecord{}, err
	}

	// Convert payload string to json.RawMessage
	event.Payload = json.RawMessage(payloadStr)

	// Parse metadata
	if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
		return eventstore.EventRecord{}, err
	}

	return event, nil
}

// List retrieves events matching the filter criteria.
func (s *EventStore) List(ctx context.Context, filter eventstore.ListFilter) ([]eventstore.EventRecord, error) {
	// Build query with filters
	query := `
		SELECT id, type, source, subject_id, correlation_id, causation_id, trace_id,
			occurred_at, payload, metadata, created_at
		FROM events
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, filter.Type)
	}
	if filter.Source != "" {
		query += " AND source = ?"
		args = append(args, filter.Source)
	}
	if filter.SubjectID != "" {
		query += " AND subject_id = ?"
		args = append(args, filter.SubjectID)
	}
	if filter.CorrelationID != "" {
		query += " AND correlation_id = ?"
		args = append(args, filter.CorrelationID)
	}

	// Order deterministically
	query += " ORDER BY occurred_at ASC, id ASC"

	// Apply limit
	limit := filter.Limit
	if limit == 0 {
		limit = 100 // safe default
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, filter.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var events []eventstore.EventRecord
	for rows.Next() {
		var event eventstore.EventRecord
		var occurredAtStr, createdAtStr string
		var payloadStr, metadataJSON string

		err := rows.Scan(
			&event.ID,
			&event.Type,
			&event.Source,
			&event.SubjectID,
			&event.CorrelationID,
			&event.CausationID,
			&event.TraceID,
			&occurredAtStr,
			&payloadStr,
			&metadataJSON,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAtStr)
		if err != nil {
			return nil, err
		}

		event.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return nil, err
		}

		// Convert payload string to json.RawMessage
		event.Payload = json.RawMessage(payloadStr)

		// Parse metadata
		if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Close closes the database connection.
func (s *EventStore) Close() error {
	return s.db.Close()
}

// isSQLiteConstraintError checks if the error is a SQLite constraint violation.
func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// SQLite UNIQUE constraint violation contains "UNIQUE constraint failed"
	return contains(err.Error(), "UNIQUE constraint failed") || contains(err.Error(), "constraint failed")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexAny(s, substr) >= 0)
}

func indexAny(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

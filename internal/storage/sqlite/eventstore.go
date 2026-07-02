package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/tiroq/praxis/internal/storage/eventstore"
	_ "modernc.org/sqlite" // register SQLite driver
)

// store is the unexported SQLite-backed implementation of eventstore.EventStore.
// Callers receive only the eventstore.EventStore interface from OpenEventStore;
// no *sql.DB is accessible outside this package.
type store struct {
	db     *sql.DB
	closed atomic.Bool
}

// OpenEventStore opens or creates a SQLite event store at the specified path.
// Use ":memory:" for an in-memory database (useful for tests).
// The return type is eventstore.EventStore; no concrete SQLite type is exposed.
func OpenEventStore(ctx context.Context, path string) (eventstore.EventStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency on file-backed databases.
	// For :memory: the WAL pragma is accepted but has no effect.
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}

	// Apply schema (idempotent — uses CREATE TABLE IF NOT EXISTS).
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &store{db: db}, nil
}

// Append stores a new event in SQLite.
func (s *store) Append(ctx context.Context, event eventstore.EventRecord) error {
	// Check context cancellation and closed state before any DB work.
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return eventstore.ErrStoreClosed
	}

	// Validate the event before touching the database.
	if err := event.Validate(); err != nil {
		return err
	}

	// Set CreatedAt if not already set.
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	// Ensure metadata is not nil for consistent round-tripping.
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}

	metadataJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

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
		if isSQLiteConstraintError(err) {
			return eventstore.ErrDuplicateEvent{ID: event.ID}
		}
		return err
	}

	return nil
}

// Get retrieves an event by ID.
func (s *store) Get(ctx context.Context, id string) (eventstore.EventRecord, error) {
	if err := ctx.Err(); err != nil {
		return eventstore.EventRecord{}, err
	}
	if s.closed.Load() {
		return eventstore.EventRecord{}, eventstore.ErrStoreClosed
	}

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

	event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAtStr)
	if err != nil {
		return eventstore.EventRecord{}, err
	}

	event.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
	if err != nil {
		return eventstore.EventRecord{}, err
	}

	event.Payload = json.RawMessage(payloadStr)

	if err := json.Unmarshal([]byte(metadataJSON), &event.Metadata); err != nil {
		return eventstore.EventRecord{}, err
	}

	return event, nil
}

// List retrieves events matching the filter criteria.
func (s *store) List(ctx context.Context, filter eventstore.ListFilter) ([]eventstore.EventRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.closed.Load() {
		return nil, eventstore.ErrStoreClosed
	}

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

	query += " ORDER BY occurred_at ASC, id ASC"

	limit := filter.Limit
	if limit == 0 {
		limit = 100
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

		event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAtStr)
		if err != nil {
			return nil, err
		}

		event.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			return nil, err
		}

		event.Payload = json.RawMessage(payloadStr)

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

// Close closes the underlying database connection.
// Close is idempotent: subsequent calls return nil without closing again.
// After Close returns, all Append/Get/List calls return ErrStoreClosed.
func (s *store) Close() error {
	if s.closed.Swap(true) {
		// Already closed — idempotent.
		return nil
	}
	return s.db.Close()
}

// isSQLiteConstraintError reports whether err is a SQLite UNIQUE constraint violation.
func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return contains(msg, "UNIQUE constraint failed") || contains(msg, "constraint failed")
}

func contains(s, substr string) bool {
	return indexAny(s, substr) >= 0 || s == substr
}

func indexAny(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

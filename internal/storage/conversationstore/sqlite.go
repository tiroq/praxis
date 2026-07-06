package conversationstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite" // register SQLite driver
)

// SQLiteStore is the unexported SQLite-backed implementation of ConversationStore.
// Callers receive only the ConversationStore interface; no *sql.DB is accessible.
type SQLiteStore struct {
	db     *sql.DB
	closed atomic.Bool
}

// OpenStore opens or creates a SQLite conversation store at the specified path.
// Use ":memory:" for an in-memory database (useful for tests).
// Applies schema idempotently using CREATE TABLE IF NOT EXISTS.
func OpenStore(ctx context.Context, path string) (ConversationStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode for better concurrency.
	// For :memory: the WAL pragma is accepted but has no effect.
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}

	// Apply schema (idempotent).
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

// CreateConversation persists a new Conversation.
func (s *SQLiteStore) CreateConversation(ctx context.Context, conv *Conversation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}
	if conv == nil {
		return ErrMissingField("conversation")
	}

	query := `
		INSERT INTO conversations (id, correlation_id, lifecycle, created_at, updated_at, last_message_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		conv.ID,
		conv.CorrelationID,
		conv.Lifecycle,
		conv.CreatedAt,
		conv.UpdatedAt,
		conv.LastMessageAt,
	)
	if err != nil {
		if isSQLiteConstraintError(err) {
			return ErrDuplicateConversation{ID: conv.ID}
		}
		return err
	}
	return nil
}

// GetConversation retrieves a Conversation by ID.
func (s *SQLiteStore) GetConversation(ctx context.Context, id string) (*Conversation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.closed.Load() {
		return nil, ErrStoreClosed
	}

	query := `
		SELECT id, correlation_id, lifecycle, created_at, updated_at, last_message_at
		FROM conversations
		WHERE id = ?
	`

	var conv Conversation
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&conv.ID,
		&conv.CorrelationID,
		&conv.Lifecycle,
		&conv.CreatedAt,
		&conv.UpdatedAt,
		&conv.LastMessageAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrConversationNotFound{ID: id}
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// GetConversationByCorrelationID retrieves or creates a Conversation by correlation_id.
func (s *SQLiteStore) GetConversationByCorrelationID(ctx context.Context, correlationID string) (*Conversation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.closed.Load() {
		return nil, ErrStoreClosed
	}

	query := `
		SELECT id, correlation_id, lifecycle, created_at, updated_at, last_message_at
		FROM conversations
		WHERE correlation_id = ?
		LIMIT 1
	`

	var conv Conversation
	err := s.db.QueryRowContext(ctx, query, correlationID).Scan(
		&conv.ID,
		&conv.CorrelationID,
		&conv.Lifecycle,
		&conv.CreatedAt,
		&conv.UpdatedAt,
		&conv.LastMessageAt,
	)
	if err == sql.ErrNoRows {
		// Create a new conversation
		conv = *NewConversation(correlationID)
		if err := s.CreateConversation(ctx, &conv); err != nil {
			return nil, err
		}
		return &conv, nil
	}
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// AppendMessage persists a new Message in append-only fashion.
func (s *SQLiteStore) AppendMessage(ctx context.Context, msg *Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}
	if msg == nil {
		return ErrMissingField("message")
	}

	// Verify conversation exists
	if _, err := s.GetConversation(ctx, msg.ConversationID); err != nil {
		return fmt.Errorf("cannot append message: %w", err)
	}

	// Marshal metadata
	metadataJSON, err := json.Marshal(msg.Metadata)
	if err != nil {
		return ErrInvalidJSON{Field: "metadata", Err: err}
	}

	query := `
		INSERT INTO messages (id, conversation_id, event_id, role, content, timestamp, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		msg.ID,
		msg.ConversationID,
		msg.EventID,
		msg.Role,
		msg.Content,
		msg.Timestamp,
		string(metadataJSON),
		msg.CreatedAt,
	)
	if err != nil {
		if isSQLiteConstraintError(err) {
			return ErrDuplicateMessage{ID: msg.ID}
		}
		return err
	}

	// Update conversation metadata
	if err := s.UpdateConversationMetadata(ctx, msg.ConversationID); err != nil {
		return fmt.Errorf("failed to update conversation metadata after appending message: %w", err)
	}

	return nil
}

// ListMessages retrieves all messages for a conversation, ordered by timestamp.
func (s *SQLiteStore) ListMessages(ctx context.Context, conversationID string, filter ListFilter) ([]*Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.closed.Load() {
		return nil, ErrStoreClosed
	}

	// Default limit if not specified
	limit := filter.Limit
	if limit <= 0 {
		limit = 1000 // safe default
	}

	query := `
		SELECT id, conversation_id, event_id, role, content, timestamp, metadata, created_at
		FROM messages
		WHERE conversation_id = ?
		ORDER BY timestamp ASC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, conversationID, limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		var metadataJSON string

		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.EventID,
			&msg.Role,
			&msg.Content,
			&msg.Timestamp,
			&metadataJSON,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if metadataJSON != "" {
			if err := json.Unmarshal([]byte(metadataJSON), &msg.Metadata); err != nil {
				return nil, ErrInvalidJSON{Field: "metadata", Err: err}
			}
		} else {
			msg.Metadata = make(map[string]string)
		}

		messages = append(messages, &msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// UpdateConversationMetadata updates last_message_at and updated_at.
func (s *SQLiteStore) UpdateConversationMetadata(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	// Get the timestamp of the last message
	query := `SELECT timestamp FROM messages WHERE conversation_id = ? ORDER BY timestamp DESC LIMIT 1`
	var lastMessageAt *string
	err := s.db.QueryRowContext(ctx, query, id).Scan(&lastMessageAt)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	updateQuery := `
		UPDATE conversations
		SET updated_at = ?, last_message_at = ?
		WHERE id = ?
	`
	_, err = s.db.ExecContext(ctx, updateQuery, now, lastMessageAt, id)
	return err
}

// ArchiveConversation marks a conversation as archived.
func (s *SQLiteStore) ArchiveConversation(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	query := `UPDATE conversations SET lifecycle = ?, updated_at = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, "archived", now, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrConversationNotFound{ID: id}
	}
	return nil
}

// Close releases database resources.
func (s *SQLiteStore) Close() error {
	if !s.closed.CompareAndSwap(false, true) {
		return ErrStoreClosed
	}
	return s.db.Close()
}

// isSQLiteConstraintError checks if an error is a SQLite UNIQUE constraint error.
func isSQLiteConstraintError(err error) bool {
	var sqlErr interface {
		Code() string
	}
	return errors.As(err, &sqlErr) && sqlErr.Code() == "SQLITE_CONSTRAINT"
}

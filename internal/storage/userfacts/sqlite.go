package userfacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	_ "modernc.org/sqlite" // register SQLite driver
)

// SQLiteStore is a concrete SQLite-backed store for extracted user fact candidates.
type SQLiteStore struct {
	db     *sql.DB
	closed atomic.Bool
}

// OpenStore opens or creates a SQLite user facts store.
func OpenStore(ctx context.Context, path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := ensureSchemaMigrations(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

// Append persists a candidate fact without revising or merging prior facts.
func (s *SQLiteStore) Append(ctx context.Context, fact *CandidateFact) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}
	if fact == nil {
		return ErrMissingField("fact")
	}
	if err := validateFact(fact); err != nil {
		return err
	}

	query := `
		INSERT INTO candidate_user_facts
			(id, user_id, correlation_id, type, value, confidence, source_event_id, source_message_id, validation_state, validation_updated_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query,
		fact.ID,
		fact.UserID,
		fact.CorrelationID,
		fact.Type,
		fact.Value,
		fact.Confidence,
		fact.SourceEventID,
		fact.SourceMessageID,
		fact.ValidationState,
		fact.ValidationUpdatedAt,
		fact.CreatedAt,
	)
	if err != nil {
		if isSQLiteConstraintError(err) {
			return ErrDuplicateFact{ID: fact.ID}
		}
		return err
	}
	return nil
}

// ListBySourceEvent returns candidate facts extracted from one source event.
func (s *SQLiteStore) ListBySourceEvent(ctx context.Context, sourceEventID string) ([]*CandidateFact, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.closed.Load() {
		return nil, ErrStoreClosed
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, correlation_id, type, value, confidence, source_event_id, source_message_id, validation_state, validation_updated_at, created_at
		FROM candidate_user_facts
		WHERE source_event_id = ?
		ORDER BY created_at ASC, id ASC
	`, sourceEventID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var facts []*CandidateFact
	for rows.Next() {
		var fact CandidateFact
		if err := rows.Scan(
			&fact.ID,
			&fact.UserID,
			&fact.CorrelationID,
			&fact.Type,
			&fact.Value,
			&fact.Confidence,
			&fact.SourceEventID,
			&fact.SourceMessageID,
			&fact.ValidationState,
			&fact.ValidationUpdatedAt,
			&fact.CreatedAt,
		); err != nil {
			return nil, err
		}
		facts = append(facts, &fact)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return facts, nil
}

// TransitionValidationState moves a fact to the next validation stage and records the transition.
func (s *SQLiteStore) TransitionValidationState(
	ctx context.Context,
	factID string,
	toState ValidationState,
	actor string,
	reason string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.closed.Load() {
		return ErrStoreClosed
	}
	factID = strings.TrimSpace(factID)
	actor = strings.TrimSpace(actor)
	reason = strings.TrimSpace(reason)
	if factID == "" {
		return ErrMissingField("fact_id")
	}
	if actor == "" {
		return ErrMissingField("actor")
	}
	if reason == "" {
		return ErrMissingField("reason")
	}
	if !IsValidValidationState(toState) {
		return ErrInvalidValidationState{State: string(toState)}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var fromStateRaw string
	if err = tx.QueryRowContext(ctx, `
		SELECT validation_state
		FROM candidate_user_facts
		WHERE id = ?
	`, factID).Scan(&fromStateRaw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrFactNotFound{ID: factID}
		}
		return err
	}

	fromState := ValidationState(fromStateRaw)
	if !CanTransitionValidationState(fromState, toState) {
		return ErrInvalidValidationTransition{
			From: fromState,
			To:   toState,
		}
	}

	transitionedAt := time.Now().UTC().Format(time.RFC3339Nano)
	result, err := tx.ExecContext(ctx, `
		UPDATE candidate_user_facts
		SET validation_state = ?, validation_updated_at = ?
		WHERE id = ?
	`, toState, transitionedAt, factID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrFactNotFound{ID: factID}
	}

	if _, err = tx.ExecContext(ctx, `
		INSERT INTO candidate_user_fact_validation_transitions
			(fact_id, from_state, to_state, actor, reason, transitioned_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, factID, fromState, toState, actor, reason, transitionedAt); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
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

func validateFact(fact *CandidateFact) error {
	required := map[string]string{
		"id":                    fact.ID,
		"user_id":               fact.UserID,
		"correlation_id":        fact.CorrelationID,
		"type":                  fact.Type,
		"value":                 fact.Value,
		"source_event_id":       fact.SourceEventID,
		"source_message_id":     fact.SourceMessageID,
		"validation_updated_at": fact.ValidationUpdatedAt,
		"created_at":            fact.CreatedAt,
	}
	for field, value := range required {
		if strings.TrimSpace(value) == "" {
			return ErrMissingField(field)
		}
	}
	if fact.Confidence < 0 || fact.Confidence > 1 {
		return ErrInvalidConfidence(fact.Confidence)
	}
	if !IsValidValidationState(fact.ValidationState) {
		return ErrInvalidValidationState{State: string(fact.ValidationState)}
	}
	return nil
}

func isSQLiteConstraintError(err error) bool {
	var codeInt interface {
		Code() int
	}
	if errors.As(err, &codeInt) {
		code := codeInt.Code()
		if code == 1555 || code == 2067 || code == 275 {
			return true
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "sqlite_constraint") || strings.Contains(msg, "constraint failed")
}

func ensureSchemaMigrations(ctx context.Context, db *sql.DB) error {
	columns, err := tableColumns(ctx, db, "candidate_user_facts")
	if err != nil {
		return err
	}
	if _, ok := columns["validation_state"]; !ok {
		if _, err := db.ExecContext(ctx, `
			ALTER TABLE candidate_user_facts
			ADD COLUMN validation_state TEXT NOT NULL DEFAULT 'extracted'
		`); err != nil {
			return err
		}
	}
	if _, ok := columns["validation_updated_at"]; !ok {
		if _, err := db.ExecContext(ctx, `
			ALTER TABLE candidate_user_facts
			ADD COLUMN validation_updated_at TEXT NOT NULL DEFAULT ''
		`); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(ctx, `
		UPDATE candidate_user_facts
		SET validation_updated_at = created_at
		WHERE validation_updated_at = ''
	`); err != nil {
		return err
	}

	return nil
}

func tableColumns(ctx context.Context, db *sql.DB, table string) (map[string]struct{}, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	columns := make(map[string]struct{})
	for rows.Next() {
		var (
			cid        int
			name       string
			typ        string
			notNull    int
			defaultV   interface{}
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultV, &primaryKey); err != nil {
			return nil, err
		}
		columns[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

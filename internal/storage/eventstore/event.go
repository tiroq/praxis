package eventstore

import (
	"encoding/json"
	"time"
)

// EventRecord is the stored representation of an event.
// It aligns with RFC-013 event envelope and is designed for append-only persistence.
type EventRecord struct {
	// Identity
	ID            string `json:"id"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`
	TraceID       string `json:"trace_id"`

	// Provenance
	Type       string    `json:"type"`
	Source     string    `json:"source"`
	SubjectID  string    `json:"subject_id"`
	OccurredAt time.Time `json:"occurred_at"`

	// Content
	Payload  json.RawMessage   `json:"payload"`
	Metadata map[string]string `json:"metadata"`

	// System fields
	CreatedAt time.Time `json:"created_at"` // when the record was stored
}

// Validate checks that an EventRecord has all required fields and valid data.
func (e *EventRecord) Validate() error {
	if e.ID == "" {
		return ErrMissingField("id")
	}
	if e.Type == "" {
		return ErrMissingField("type")
	}
	if e.Source == "" {
		return ErrMissingField("source")
	}
	if e.SubjectID == "" {
		return ErrMissingField("subject_id")
	}
	if e.OccurredAt.IsZero() {
		return ErrMissingField("occurred_at")
	}
	if len(e.Payload) == 0 {
		return ErrMissingField("payload")
	}

	// Validate that Payload is valid JSON
	var tmp interface{}
	if err := json.Unmarshal(e.Payload, &tmp); err != nil {
		return ErrInvalidJSON{Field: "payload", Err: err}
	}

	return nil
}

// Clone creates a deep copy of the EventRecord.
// This is used internally to ensure immutability.
func (e *EventRecord) Clone() EventRecord {
	clone := EventRecord{
		ID:            e.ID,
		CorrelationID: e.CorrelationID,
		CausationID:   e.CausationID,
		TraceID:       e.TraceID,
		Type:          e.Type,
		Source:        e.Source,
		SubjectID:     e.SubjectID,
		OccurredAt:    e.OccurredAt,
		CreatedAt:     e.CreatedAt,
	}

	// Deep copy Payload
	if e.Payload != nil {
		clone.Payload = make(json.RawMessage, len(e.Payload))
		copy(clone.Payload, e.Payload)
	}

	// Deep copy Metadata
	if e.Metadata != nil {
		clone.Metadata = make(map[string]string, len(e.Metadata))
		for k, v := range e.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

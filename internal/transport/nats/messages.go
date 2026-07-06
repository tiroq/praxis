package nats

import (
	"encoding/json"
	"time"
)

// InputMessage is the wire format for events arriving on the input subject.
// Fields map directly to kernel.Event fields; only the minimum required set
// is exposed over the wire.
type InputMessage struct {
	// ID is the globally unique event identifier (e.g. "evt_...").
	ID string `json:"id"`

	// CorrelationID groups related messages into one logical activity.
	CorrelationID string `json:"correlation_id,omitempty"`

	// Source identifies the originating system or user (e.g. "manual", "telegram").
	Source string `json:"source"`

	// Text is the primary human-readable content to be processed by the kernel.
	Text string `json:"text"`

	// Timestamp is the time the event occurred or was observed.
	Timestamp time.Time `json:"timestamp"`

	// Metadata contains optional enrichment fields.
	Metadata map[string]string `json:"metadata"`
}

// Validate checks that the InputMessage satisfies the minimum wire-level requirements.
func (m InputMessage) Validate() error {
	if m.ID == "" {
		return ErrEmptyMessageID
	}
	if m.Text == "" {
		return ErrEmptyMessageText
	}
	return nil
}

// OutputMessage is the wire format published to the output subject after the
// kernel pipeline completes.
type OutputMessage struct {
	// InputEventID is the ID of the InputMessage that triggered this result.
	InputEventID string `json:"input_event_id"`

	// CorrelationID groups related messages into one logical activity.
	CorrelationID string `json:"correlation_id,omitempty"`

	// Status is "ok" on success or "error" on failure.
	Status string `json:"status"`

	// Result holds the full pipeline result JSON on success.
	Result json.RawMessage `json:"result,omitempty"`

	// ReplyText holds the LLM-generated user-facing reply text when available.
	ReplyText string `json:"reply_text,omitempty"`

	// Error holds the error string on failure.
	Error *string `json:"error,omitempty"`

	// Metadata carries transport-level context propagated from input.
	Metadata map[string]string `json:"metadata,omitempty"`

	// ProcessedAt is when the pipeline completed.
	ProcessedAt time.Time `json:"processed_at"`
}

package nats

import (
	"time"

	"github.com/tiroq/praxis/internal/core/kernel"
)

// InputMessage is the wire format for events arriving on the input subject.
// Fields map directly to kernel.Event fields; only the minimum required set
// is exposed over the wire.
type InputMessage struct {
	// ID is the globally unique event identifier (e.g. "evt_...").
	ID string `json:"id"`

	// Source identifies the originating system or user (e.g. "manual", "telegram").
	Source string `json:"source"`

	// Text is the primary human-readable content to be processed by the kernel.
	Text string `json:"text"`

	// Timestamp is the time the event occurred or was observed.
	Timestamp time.Time `json:"timestamp"`

	// Metadata contains optional enrichment fields.
	Metadata map[string]string `json:"metadata"`
}

// validate checks that the InputMessage satisfies the minimum requirements
// for a valid kernel.Event.
func (m InputMessage) validate() error {
	if m.ID == "" {
		return ErrEmptyMessageID
	}
	if m.Text == "" {
		return ErrEmptyMessageText
	}
	return nil
}

// toKernelEvent converts a validated InputMessage to a kernel.Event.
// OccurredAt and ObservedAt are both set from m.Timestamp; if Timestamp is
// zero the current UTC time is used. Confidence defaults to 1.0 for
// externally supplied events.
func (m InputMessage) toKernelEvent() kernel.Event {
	ts := m.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	meta := m.Metadata
	if meta == nil {
		meta = map[string]string{}
	}

	return kernel.Event{
		ID:               m.ID,
		Source:           m.Source,
		Text:             m.Text,
		OccurredAt:       ts,
		ObservedAt:       time.Now().UTC(),
		Type:             "external.text",
		ContentType:      "text/plain",
		Payload:          map[string]any{},
		Metadata:         meta,
		Confidence:       1.0,
		TrustLevel:       kernel.TrustLevelMedium,
		ValidationStatus: kernel.ValidationStatusPending,
		Origin:           kernel.EventOriginExternal,
	}
}

// OutputMessage is the wire format published to the output subject after the
// kernel pipeline completes.
type OutputMessage struct {
	// InputEventID is the ID of the InputMessage that triggered this result.
	InputEventID string `json:"input_event_id"`

	// Status is "ok" on success or "error" on failure.
	Status string `json:"status"`

	// Result holds the full kernel.PipelineResult on success.
	Result *kernel.PipelineResult `json:"result,omitempty"`

	// Error holds the error string on failure.
	Error *string `json:"error,omitempty"`

	// ProcessedAt is when the pipeline completed.
	ProcessedAt time.Time `json:"processed_at"`
}

// newOutputOK constructs a success OutputMessage from a PipelineResult.
func newOutputOK(inputEventID string, result kernel.PipelineResult) OutputMessage {
	return OutputMessage{
		InputEventID: inputEventID,
		Status:       "ok",
		Result:       &result,
		ProcessedAt:  time.Now().UTC(),
	}
}

// newOutputError constructs an error OutputMessage from an error value.
func newOutputError(inputEventID string, err error) OutputMessage {
	s := err.Error()
	return OutputMessage{
		InputEventID: inputEventID,
		Status:       "error",
		Error:        &s,
		ProcessedAt:  time.Now().UTC(),
	}
}

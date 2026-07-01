package kernel

import "time"

// TrustLevel represents the degree of confidence in an event's authenticity.
type TrustLevel string

const (
	TrustLevelHigh    TrustLevel = "high"
	TrustLevelMedium  TrustLevel = "medium"
	TrustLevelLow     TrustLevel = "low"
	TrustLevelUnknown TrustLevel = "unknown"
)

// ValidationStatus represents the result of event validation checks.
type ValidationStatus string

const (
	ValidationStatusValid   ValidationStatus = "valid"
	ValidationStatusInvalid ValidationStatus = "invalid"
	ValidationStatusPending ValidationStatus = "pending"
)

// EventOrigin classifies where an event originated.
type EventOrigin string

const (
	EventOriginExternal    EventOrigin = "external"
	EventOriginInternal    EventOrigin = "internal"
	EventOriginUser        EventOrigin = "user"
	EventOriginAI          EventOrigin = "ai"
	EventOriginScheduled   EventOrigin = "scheduled"
	EventOriginIntegration EventOrigin = "integration"
)

// Event is an immutable fact that something happened, was observed, or was received.
// It is the primary entry point into the Praxis execution spine.
// Per RFC-013: events are immutable once recorded; the Payload and all identity
// fields must never be mutated.
type Event struct {
	// Identity
	ID            string // globally unique; must be non-empty
	CorrelationID string // groups related events into one logical activity
	CausationID   string // references the event that directly caused this event
	TraceID       string // end-to-end distributed tracing

	// Provenance
	Source     string    // originating system, user, or process
	Actor      string    // who or what performed the action
	OccurredAt time.Time // when the event occurred or was observed
	ObservedAt time.Time // when the event was received by Praxis

	// Content
	Type        string            // classification, e.g. "email.received"
	ContentType string            // MIME type or format of the payload
	Text        string            // human-readable summary or raw text (primary ingestion field)
	Payload     map[string]any    // immutable raw event data
	Metadata    map[string]string // enrichable context (processing info, tags, etc.)

	// Quality
	Confidence       float64          // 0.0–1.0 certainty in the event's accuracy
	TrustLevel       TrustLevel       // degree of confidence in authenticity
	ValidationStatus ValidationStatus // result of validation checks
	SchemaVersion    string           // version of the event schema
	Origin           EventOrigin      // classification by origin
}

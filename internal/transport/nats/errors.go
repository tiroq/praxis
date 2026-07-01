package nats

import "errors"

// Sentinel errors for the NATS transport adapter.
// Callers should use errors.Is for matching.
var (
	// ErrEmptyMessageID is returned when an inbound message has no id field.
	ErrEmptyMessageID = errors.New("nats: message id must not be empty")

	// ErrEmptyMessageText is returned when an inbound message has no text field.
	ErrEmptyMessageText = errors.New("nats: message text must not be empty")

	// ErrPublishFailed is returned when the result event cannot be published.
	ErrPublishFailed = errors.New("nats: failed to publish result event")

	// ErrInvalidJSON is returned when the inbound message payload is not valid JSON.
	ErrInvalidJSON = errors.New("nats: invalid JSON payload")
)

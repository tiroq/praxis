package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// jsPublish is the minimal JetStream interface needed by Publisher.
// Defined as a narrow interface to allow fakes in unit tests.
type jsPublish interface {
	Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error)
}

// Publisher publishes OutputMessages to a NATS JetStream subject.
type Publisher struct {
	js      jsPublish
	subject string
}

// NewPublisher constructs a Publisher that writes to the given subject.
func NewPublisher(js nats.JetStreamContext, subject string) *Publisher {
	return &Publisher{js: js, subject: subject}
}

// Publish serialises msg to JSON and publishes it to the output subject.
// Returns ErrPublishFailed (wrapping the underlying cause) on any error.
func (p *Publisher) Publish(msg OutputMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%w: marshal: %w", ErrPublishFailed, err)
	}

	if _, err := p.js.Publish(p.subject, data); err != nil {
		return fmt.Errorf("%w: %w", ErrPublishFailed, err)
	}
	return nil
}

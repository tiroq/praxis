package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// Client holds a live NATS connection and a JetStream context.
// It is the single point of contact with the NATS server for the adapter.
type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

// NewClient dials the NATS server at cfg.URL, obtains a JetStream context,
// and ensures the configured stream exists with the required subjects.
// The caller must call Close when the client is no longer needed.
func NewClient(cfg Config) (*Client, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats: connect to %s: %w", cfg.URL, err)
	}

	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats: obtain JetStream context: %w", err)
	}

	if err := ensureStream(js, cfg); err != nil {
		nc.Close()
		return nil, err
	}

	return &Client{conn: nc, js: js}, nil
}

// JetStream returns the JetStream context for use by publisher and subscriber.
func (c *Client) JetStream() nats.JetStreamContext { return c.js }

// Close drains and closes the underlying NATS connection gracefully.
func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Drain()
	}
}

// ensureStream creates the stream if it does not already exist.
// If the stream exists, the configuration is left untouched.
func ensureStream(js nats.JetStreamContext, cfg Config) error {
	_, err := js.StreamInfo(cfg.StreamName)
	if err == nil {
		// Stream already exists.
		return nil
	}
	if err != nats.ErrStreamNotFound {
		return fmt.Errorf("nats: check stream %q: %w", cfg.StreamName, err)
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     cfg.StreamName,
		Subjects: []string{cfg.InputSubject, cfg.OutputSubject},
	})
	if err != nil {
		return fmt.Errorf("nats: create stream %q: %w", cfg.StreamName, err)
	}
	return nil
}

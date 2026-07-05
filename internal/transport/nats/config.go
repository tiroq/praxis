// Package nats provides a JetStream transport adapter that wraps the Core Kernel.
// It is purely infrastructure: it connects to NATS, decodes inbound messages,
// drives the kernel pipeline, and publishes result events. The kernel itself is
// never modified; all NATS concerns live exclusively in this package.
package nats

import (
	"os"
	"strconv"
	"time"
)

// Config holds the NATS JetStream configuration for the worker adapter.
// All fields have safe defaults that are applied by DefaultConfig.
type Config struct {
	// URL is the NATS server address.
	URL string

	// StreamName is the JetStream stream name that must exist or be created.
	StreamName string

	// InputSubject is the subject the worker subscribes to for inbound events.
	InputSubject string

	// OutputSubject is the subject the worker publishes processed results to.
	OutputSubject string

	// DLQSubject is the subject terminally failed messages are published to.
	DLQSubject string

	// Durable is the durable consumer name for at-least-once delivery.
	Durable string

	// AckWait is the duration the server waits for an ack before redelivering.
	AckWait time.Duration

	// MaxDeliver is the maximum number of delivery attempts before the message
	// is moved to the dead-letter or simply dropped.
	MaxDeliver int
}

// DefaultConfig returns a Config pre-populated with safe defaults.
// Override individual fields before passing to NewClient.
func DefaultConfig() Config {
	return Config{
		URL:           "nats://localhost:4222",
		StreamName:    "PRAXIS",
		InputSubject:  "praxis.kernel.input",
		OutputSubject: "praxis.kernel.output",
		DLQSubject:    "praxis.kernel.dlq",
		Durable:       "praxis-worker",
		AckWait:       30 * time.Second,
		MaxDeliver:    3,
	}
}

// ConfigFromEnv reads the Config from environment variables, falling back to
// DefaultConfig values for any unset variable.
//
// Environment variables:
//
//	NATS_URL               – NATS server URL
//	NATS_STREAM            – JetStream stream name
//	NATS_INPUT_SUBJECT     – inbound subject
//	NATS_OUTPUT_SUBJECT    – outbound subject
//	NATS_DLQ_SUBJECT       – dead-letter subject
//	NATS_DURABLE           – durable consumer name
//	NATS_ACK_WAIT_SECONDS  – ack wait in seconds (integer)
//	NATS_MAX_DELIVER        – max delivery attempts (integer)
func ConfigFromEnv() Config {
	cfg := DefaultConfig()

	if v := os.Getenv("NATS_URL"); v != "" {
		cfg.URL = v
	}
	if v := os.Getenv("NATS_STREAM"); v != "" {
		cfg.StreamName = v
	}
	if v := os.Getenv("NATS_INPUT_SUBJECT"); v != "" {
		cfg.InputSubject = v
	}
	if v := os.Getenv("NATS_OUTPUT_SUBJECT"); v != "" {
		cfg.OutputSubject = v
	}
	if v := os.Getenv("NATS_DLQ_SUBJECT"); v != "" {
		cfg.DLQSubject = v
	}
	if v := os.Getenv("NATS_DURABLE"); v != "" {
		cfg.Durable = v
	}
	if v := os.Getenv("NATS_ACK_WAIT_SECONDS"); v != "" {
		if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
			cfg.AckWait = time.Duration(secs) * time.Second
		}
	}
	if v := os.Getenv("NATS_MAX_DELIVER"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.MaxDeliver = n
		}
	}

	return cfg
}

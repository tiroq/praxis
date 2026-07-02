package praxiscli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nuid"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

const (
	defaultSource          = "praxis-cli"
	metadataCorrelationKey = "correlation_id"
)

// PublishedMessage contains identifiers of the message sent to the input subject.
type PublishedMessage struct {
	MessageID     string
	CorrelationID string
}

// PublishRequest describes a single inbound message publish operation.
type PublishRequest struct {
	Text          string
	Source        string
	MessageID     string
	CorrelationID string
	Metadata      map[string]string
}

// WatchRequest controls watch command behavior.
type WatchRequest struct {
	Writer      io.Writer
	MaxMessages int
	PollTimeout time.Duration
}

// App orchestrates CLI transport operations over NATS.
type App struct {
	cfg       natstransport.Config
	logger    *slog.Logger
	factory   TransportFactory
	now       func() time.Time
	idFactory func(prefix string) string
}

// NewApp constructs the CLI application service.
func NewApp(cfg natstransport.Config, logger *slog.Logger, factory TransportFactory) *App {
	return &App{
		cfg:       cfg,
		logger:    logger,
		factory:   factory,
		now:       func() time.Time { return time.Now().UTC() },
		idFactory: newID,
	}
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, nuid.Next())
}

// Publish validates input, builds an InputMessage, and publishes it to the configured input subject.
func (a *App) Publish(_ context.Context, req PublishRequest) (PublishedMessage, error) {
	text := strings.TrimSpace(req.Text)
	if text == "" {
		return PublishedMessage{}, errors.New("publish text must not be empty")
	}

	id := strings.TrimSpace(req.MessageID)
	if id == "" {
		id = a.idFactory("evt")
	}

	corr := strings.TrimSpace(req.CorrelationID)
	if corr == "" {
		corr = a.idFactory("corr")
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = defaultSource
	}

	meta := cloneMetadata(req.Metadata)
	meta[metadataCorrelationKey] = corr

	msg := natstransport.InputMessage{
		ID:            id,
		CorrelationID: corr,
		Source:        source,
		Text:          text,
		Timestamp:     a.now(),
		Metadata:      meta,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return PublishedMessage{}, fmt.Errorf("marshal input message: %w", err)
	}

	transport, err := a.factory(a.cfg)
	if err != nil {
		return PublishedMessage{}, fmt.Errorf("connect transport: %w", err)
	}
	defer transport.Close()

	if err := transport.PublishInput(a.cfg.InputSubject, payload, id); err != nil {
		return PublishedMessage{}, fmt.Errorf("publish input message: %w", err)
	}

	a.logger.Info("published input message", "subject", a.cfg.InputSubject, "message_id", id, "correlation_id", corr)
	return PublishedMessage{MessageID: id, CorrelationID: corr}, nil
}

// Watch subscribes to output messages and pretty-prints each output payload.
func (a *App) Watch(ctx context.Context, req WatchRequest) error {
	if req.Writer == nil {
		return errors.New("watch writer must not be nil")
	}

	pollTimeout := req.PollTimeout
	if pollTimeout <= 0 {
		pollTimeout = time.Second
	}

	transport, err := a.factory(a.cfg)
	if err != nil {
		return fmt.Errorf("connect transport: %w", err)
	}
	defer transport.Close()

	sub, err := transport.SubscribeOutput(a.cfg.OutputSubject, a.cfg.StreamName)
	if err != nil {
		return fmt.Errorf("subscribe output: %w", err)
	}
	defer func() { _ = sub.Close() }()

	processed := 0
	for {
		if req.MaxMessages > 0 && processed >= req.MaxMessages {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}

		delivery, err := sub.Next(pollTimeout)
		if err != nil {
			if errors.Is(err, ErrPollTimeout) {
				continue
			}
			if errors.Is(err, ErrConnectionClosed) {
				return nil
			}
			return fmt.Errorf("fetch output message: %w", err)
		}

		if err := a.processDelivery(req.Writer, delivery); err != nil {
			return err
		}
		processed++
	}
}

func (a *App) processDelivery(w io.Writer, delivery Delivery) error {
	var out natstransport.OutputMessage
	if err := json.Unmarshal(delivery.Data(), &out); err != nil {
		_ = delivery.Nak()
		return fmt.Errorf("decode output message: %w", err)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		_ = delivery.Nak()
		return fmt.Errorf("write output message: %w", err)
	}

	if err := delivery.Ack(); err != nil {
		return fmt.Errorf("ack output message: %w", err)
	}

	return nil
}

func cloneMetadata(src map[string]string) map[string]string {
	out := make(map[string]string, len(src)+1)
	for k, v := range src {
		out[k] = v
	}
	return out
}

// Delivery is one message delivery returned by a Subscription.
type Delivery interface {
	Data() []byte
	Ack() error
	Nak() error
}

// Subscription consumes output deliveries from the transport.
type Subscription interface {
	Next(timeout time.Duration) (Delivery, error)
	Close() error
}

// Transport describes the transport client behavior used by the CLI logic.
type Transport interface {
	PublishInput(subject string, payload []byte, msgID string) error
	SubscribeOutput(subject, stream string) (Subscription, error)
	Close()
}

// TransportFactory builds a transport client from config.
type TransportFactory func(cfg natstransport.Config) (Transport, error)

var (
	// ErrPollTimeout indicates no message was received in the configured poll interval.
	ErrPollTimeout = errors.New("poll timeout")
	// ErrConnectionClosed indicates the transport connection has closed.
	ErrConnectionClosed = errors.New("connection closed")
)

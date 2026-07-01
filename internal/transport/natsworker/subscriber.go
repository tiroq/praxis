package natsworker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/core/kernel"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

// KernelRunner is the interface the Subscriber uses to execute the pipeline.
// It is satisfied by *kernel.Kernel and defined here for fake-based tests.
type KernelRunner interface {
	Run(ctx context.Context, event kernel.Event) (kernel.PipelineResult, error)
}

// MessagePublisher is the interface the Subscriber uses to publish results.
type MessagePublisher interface {
	Publish(msg natstransport.OutputMessage) error
}

// Subscriber wraps a JetStream pull-subscribe loop and routes each message
// through the kernel pipeline before publishing the result.
type Subscriber struct {
	js        nats.JetStreamContext
	cfg       natstransport.Config
	kernel    KernelRunner
	publisher MessagePublisher
	logger    *slog.Logger
}

// NewSubscriber constructs a Subscriber. All parameters are required.
func NewSubscriber(
	js nats.JetStreamContext,
	cfg natstransport.Config,
	k KernelRunner,
	pub MessagePublisher,
	logger *slog.Logger,
) *Subscriber {
	return &Subscriber{
		js:        js,
		cfg:       cfg,
		kernel:    k,
		publisher: pub,
		logger:    logger,
	}
}

type pullSubscription interface {
	Fetch(n int, opts ...nats.PullOpt) ([]*nats.Msg, error)
	Unsubscribe() error
}

// Run subscribes to the input subject and processes messages until ctx is cancelled.
func (s *Subscriber) Run(ctx context.Context) error {
	sub, err := s.js.PullSubscribe(
		s.cfg.InputSubject,
		s.cfg.Durable,
		nats.AckWait(s.cfg.AckWait),
		nats.MaxDeliver(s.cfg.MaxDeliver),
	)
	if err != nil {
		return fmt.Errorf("nats: pull subscribe %q: %w", s.cfg.InputSubject, err)
	}
	defer sub.Unsubscribe() //nolint:errcheck

	s.logger.Info("nats subscriber started",
		"subject", s.cfg.InputSubject,
		"durable", s.cfg.Durable,
	)

	return s.runLoop(ctx, sub)
}

func (s *Subscriber) runLoop(ctx context.Context, sub pullSubscription) error {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("nats subscriber stopping", "reason", ctx.Err())
			return nil
		default:
		}

		msgs, err := sub.Fetch(1, nats.MaxWait(500*time.Millisecond))
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			if errors.Is(err, nats.ErrConnectionClosed) || errors.Is(err, context.Canceled) {
				return nil
			}
			s.logger.Warn("nats fetch error", "err", err)
			continue
		}

		for _, msg := range msgs {
			s.handleMessage(ctx, &natsMsgWrapper{m: msg})
		}
	}
}

type natsMsg interface {
	Ack() error
	Nak() error
	Term() error
	GetData() []byte
	GetSubject() string
}

type natsMsgWrapper struct{ m *nats.Msg }

func (w *natsMsgWrapper) Ack() error         { return w.m.Ack() }
func (w *natsMsgWrapper) Nak() error         { return w.m.Nak() }
func (w *natsMsgWrapper) Term() error        { return w.m.Term() }
func (w *natsMsgWrapper) GetData() []byte    { return w.m.Data }
func (w *natsMsgWrapper) GetSubject() string { return w.m.Subject }

func (s *Subscriber) handleMessage(ctx context.Context, msg natsMsg) {
	var input natstransport.InputMessage
	if err := json.Unmarshal(msg.GetData(), &input); err != nil {
		s.logger.Error("nats: invalid JSON - terminating message",
			"err", err,
			"subject", msg.GetSubject(),
		)
		_ = msg.Term()
		return
	}

	if err := input.Validate(); err != nil {
		s.logger.Error("nats: invalid message - terminating",
			"id", input.ID,
			"err", err,
		)
		_ = msg.Term()
		return
	}

	event := toKernelEvent(input)
	result, kernelErr := s.kernel.Run(ctx, event)

	var out natstransport.OutputMessage
	if kernelErr != nil {
		s.logger.Error("nats: kernel pipeline failed",
			"event_id", input.ID,
			"err", kernelErr,
		)
		out = newOutputError(input.ID, kernelErr)
	} else {
		out = newOutputOK(input.ID, result)
	}

	if err := s.publisher.Publish(out); err != nil {
		s.logger.Error("nats: publish failed - naking message",
			"event_id", input.ID,
			"err", err,
		)
		_ = msg.Nak()
		return
	}

	if err := msg.Ack(); err != nil {
		s.logger.Warn("nats: ack failed", "event_id", input.ID, "err", err)
	}

	s.logger.Info("nats: message processed",
		"event_id", input.ID,
		"status", out.Status,
	)
}

func toKernelEvent(m natstransport.InputMessage) kernel.Event {
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
		CorrelationID:    m.CorrelationID,
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

func newOutputOK(inputEventID string, result kernel.PipelineResult) natstransport.OutputMessage {
	payload, err := json.Marshal(result)
	if err != nil {
		err = fmt.Errorf("marshal pipeline result: %w", err)
		return newOutputError(inputEventID, err)
	}

	return natstransport.OutputMessage{
		InputEventID: inputEventID,
		Status:       "ok",
		Result:       payload,
		ProcessedAt:  time.Now().UTC(),
	}
}

func newOutputError(inputEventID string, err error) natstransport.OutputMessage {
	s := err.Error()
	return natstransport.OutputMessage{
		InputEventID: inputEventID,
		Status:       "error",
		Error:        &s,
		ProcessedAt:  time.Now().UTC(),
	}
}

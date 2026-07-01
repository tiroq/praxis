package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/core/kernel"
)

// KernelRunner is the interface the Subscriber uses to execute the pipeline.
// It is satisfied by *kernel.Kernel and is defined here to allow easy mocking
// in unit tests without importing the kernel in test helpers.
type KernelRunner interface {
	Run(ctx context.Context, event kernel.Event) (kernel.PipelineResult, error)
}

// MessagePublisher is the interface the Subscriber uses to publish results.
// Satisfied by *Publisher; separated for testability.
type MessagePublisher interface {
	Publish(msg OutputMessage) error
}

// Subscriber wraps a JetStream push-subscribe loop and routes each message
// through the kernel pipeline before publishing the result.
//
// Ack / Nak contract (enforced strictly):
//   - Ack is sent only after a successful Publish.
//   - If Publish fails, the message is Nak'd so it is redelivered.
//   - If the payload is invalid JSON or fails validation, the message is
//     terminated (TermMsg) to prevent endless redelivery of poison messages.
type Subscriber struct {
	js        nats.JetStreamContext
	cfg       Config
	kernel    KernelRunner
	publisher MessagePublisher
	logger    *slog.Logger
}

// NewSubscriber constructs a Subscriber. All parameters are required.
func NewSubscriber(
	js nats.JetStreamContext,
	cfg Config,
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

// pullSubscription is the minimal interface for a JetStream pull subscription.
// Defined to allow fake implementations in unit tests.
type pullSubscription interface {
	Fetch(n int, opts ...nats.PullOpt) ([]*nats.Msg, error)
	Unsubscribe() error
}

// Run subscribes to the input subject and processes messages until ctx is
// cancelled.  It blocks until the context is done.
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

// runLoop is the core fetch-and-process loop. It is separated from Run so
// that it can be tested with a fake pullSubscription without a real server.
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

// natsMsg is the minimal interface of *nats.Msg used by handleMessage.
// Defined here so tests can pass a fake without a live NATS connection.
type natsMsg interface {
	Ack() error
	Nak() error
	Term() error
	GetData() []byte
	GetSubject() string
}

// natsMsgWrapper adapts *nats.Msg to the natsMsg interface.
type natsMsgWrapper struct{ m *nats.Msg }

func (w *natsMsgWrapper) Ack() error         { return w.m.Ack() }
func (w *natsMsgWrapper) Nak() error         { return w.m.Nak() }
func (w *natsMsgWrapper) Term() error        { return w.m.Term() }
func (w *natsMsgWrapper) GetData() []byte    { return w.m.Data }
func (w *natsMsgWrapper) GetSubject() string { return w.m.Subject }

// handleMessage decodes, validates, runs the pipeline, and publishes the result.
// Ack/Nak/Term semantics are enforced strictly inside this method.
func (s *Subscriber) handleMessage(ctx context.Context, msg natsMsg) {
	var input InputMessage
	if err := json.Unmarshal(msg.GetData(), &input); err != nil {
		s.logger.Error("nats: invalid JSON — terminating message",
			"err", err,
			"subject", msg.GetSubject(),
		)
		_ = msg.Term()
		return
	}

	if err := input.validate(); err != nil {
		s.logger.Error("nats: invalid message — terminating",
			"id", input.ID,
			"err", err,
		)
		_ = msg.Term()
		return
	}

	event := input.toKernelEvent()
	result, kernelErr := s.kernel.Run(ctx, event)

	var out OutputMessage
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
		s.logger.Error("nats: publish failed — naking message",
			"event_id", input.ID,
			"err", err,
		)
		_ = msg.Nak()
		return
	}

	// Ack only after a successful publish.
	if err := msg.Ack(); err != nil {
		s.logger.Warn("nats: ack failed", "event_id", input.ID, "err", err)
	}

	s.logger.Info("nats: message processed",
		"event_id", input.ID,
		"status", out.Status,
	)
}

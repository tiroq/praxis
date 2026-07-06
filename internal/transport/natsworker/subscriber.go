package natsworker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage/conversationstore"
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

// ReplyGenerator builds a user-facing assistant reply from input and kernel output.
type ReplyGenerator func(ctx context.Context, input natstransport.InputMessage, result kernel.PipelineResult) (string, error)

// Subscriber wraps a JetStream pull-subscribe loop and routes each message
// through the kernel pipeline before publishing the result.
// Optionally persists conversation history to a rebuildable SQLite projection store.
type Subscriber struct {
	js                nats.JetStreamContext
	cfg               natstransport.Config
	kernel            KernelRunner
	publisher         MessagePublisher
	replyGenerator    ReplyGenerator
	replyTimeout      time.Duration
	logger            *slog.Logger
	conversationStore *conversationstore.SQLiteStore // optional; nil if not provided
}

const defaultReplyTimeout = 5 * time.Second

// NewSubscriber constructs a Subscriber. Parameters are required except conversationStore which is optional.
func NewSubscriber(
	js nats.JetStreamContext,
	cfg natstransport.Config,
	k KernelRunner,
	pub MessagePublisher,
	logger *slog.Logger,
) *Subscriber {
	return &Subscriber{
		js:                js,
		cfg:               cfg,
		kernel:            k,
		publisher:         pub,
		replyTimeout:      defaultReplyTimeout,
		logger:            logger,
		conversationStore: nil,
	}
}

// WithConversationStore sets the concrete SQLite conversation projection store.
// Returns the Subscriber for method chaining.
func (s *Subscriber) WithConversationStore(store *conversationstore.SQLiteStore) *Subscriber {
	s.conversationStore = store
	return s
}

// WithReplyGenerator configures optional LLM-powered assistant reply generation.
// If timeout is <= 0, a safe default timeout is used.
func (s *Subscriber) WithReplyGenerator(generator ReplyGenerator, timeout time.Duration) *Subscriber {
	s.replyGenerator = generator
	if timeout > 0 {
		s.replyTimeout = timeout
	} else {
		s.replyTimeout = defaultReplyTimeout
	}
	return s
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
	GetNumDelivered() uint64
}

type natsMsgWrapper struct{ m *nats.Msg }

func (w *natsMsgWrapper) Ack() error         { return w.m.Ack() }
func (w *natsMsgWrapper) Nak() error         { return w.m.Nak() }
func (w *natsMsgWrapper) Term() error        { return w.m.Term() }
func (w *natsMsgWrapper) GetData() []byte    { return w.m.Data }
func (w *natsMsgWrapper) GetSubject() string { return w.m.Subject }
func (w *natsMsgWrapper) GetNumDelivered() uint64 {
	meta, err := w.m.Metadata()
	if err != nil || meta == nil || meta.NumDelivered == 0 {
		return 1
	}
	return meta.NumDelivered
}

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
	correlationID := event.CorrelationID

	// Optionally persist input message to conversation store (RFC-033: Message Projection)
	var conversationID string
	if s.conversationStore != nil {
		conv, err := s.conversationStore.GetConversationByCorrelationID(ctx, correlationID)
		if err != nil {
			s.logger.Warn("nats: failed to get/create conversation",
				"correlation_id", correlationID,
				"err", err,
			)
			// Non-fatal; continue without persistence
		} else {
			conversationID = conv.ID
			inputMsg := conversationstore.NewMessage(
				conv.ID,
				input.ID,
				"user",
				input.Text,
				input.Timestamp.Format(time.RFC3339Nano),
				input.Metadata,
			)
			if err := s.conversationStore.AppendMessage(ctx, inputMsg); err != nil {
				s.logger.Warn("nats: failed to persist input message",
					"event_id", input.ID,
					"conversation_id", conversationID,
					"err", err,
				)
				// Non-fatal; continue without persistence
			}
		}
	}

	result, kernelErr := s.kernel.Run(ctx, event)

	var out natstransport.OutputMessage
	if kernelErr != nil {
		s.logger.Error("nats: kernel pipeline failed",
			"event_id", input.ID,
			"err", kernelErr,
		)
		out = newOutputError(input.ID, correlationID, input.Metadata, kernelErr)
	} else {
		assistantReply := ""
		if s.replyGenerator != nil {
			replyCtx := ctx
			cancel := func() {}
			if s.replyTimeout > 0 {
				replyCtx, cancel = context.WithTimeout(ctx, s.replyTimeout)
			}
			reply, replyErr := s.replyGenerator(replyCtx, input, result)
			cancel()

			if replyErr != nil {
				s.logger.Warn("nats: llm reply generation failed; using deterministic fallback",
					"event_id", input.ID,
					"err", replyErr,
				)
				assistantReply = deterministicFallbackReply(input.Text)
			} else {
				assistantReply = strings.TrimSpace(reply)
				if assistantReply == "" {
					s.logger.Warn("nats: llm reply generation returned empty reply; using deterministic fallback",
						"event_id", input.ID,
					)
					assistantReply = deterministicFallbackReply(input.Text)
				}
			}
		}

		out = newOutputOK(input.ID, correlationID, input.Metadata, result, assistantReply)
	}

	if err := s.publisher.Publish(out); err != nil {
		delivered := msg.GetNumDelivered()
		s.logger.Error("nats: publish failed - naking message",
			"event_id", input.ID,
			"num_delivered", delivered,
			"err", err,
		)

		if s.shouldSendToDLQ(delivered) {
			if dlqErr := s.publishDLQ(input, out, err, delivered); dlqErr != nil {
				s.logger.Error("nats: dlq publish failed",
					"event_id", input.ID,
					"num_delivered", delivered,
					"err", dlqErr,
				)
			} else {
				if ackErr := msg.Ack(); ackErr != nil {
					s.logger.Warn("nats: ack after dlq failed", "event_id", input.ID, "err", ackErr)
				}
				return
			}
		}

		_ = msg.Nak()
		return
	}

	// Optionally persist output message to conversation store (RFC-033: Message Projection)
	if s.conversationStore != nil && conversationID != "" && out.Status == "ok" {
		assistantText := strings.TrimSpace(out.AssistantReply)
		if assistantText == "" {
			assistantText = string(out.Result)
		}

		outputMsg := conversationstore.NewMessage(
			conversationID,
			out.InputEventID,
			"assistant",
			assistantText,
			out.ProcessedAt.Format(time.RFC3339Nano),
			out.Metadata,
		)
		if err := s.conversationStore.AppendMessage(ctx, outputMsg); err != nil {
			s.logger.Warn("nats: failed to persist output message",
				"event_id", out.InputEventID,
				"conversation_id", conversationID,
				"err", err,
			)
			// Non-fatal; continue
		}
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

	corr := strings.TrimSpace(m.CorrelationID)
	if corr == "" {
		corr = m.ID
	}

	return kernel.Event{
		ID:               m.ID,
		CorrelationID:    corr,
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

func newOutputOK(inputEventID string, correlationID string, metadata map[string]string, result kernel.PipelineResult, assistantReply string) natstransport.OutputMessage {
	payload, err := json.Marshal(result)
	if err != nil {
		err = fmt.Errorf("marshal pipeline result: %w", err)
		return newOutputError(inputEventID, correlationID, metadata, err)
	}

	return natstransport.OutputMessage{
		InputEventID:   inputEventID,
		CorrelationID:  correlationID,
		Status:         "ok",
		Result:         payload,
		AssistantReply: strings.TrimSpace(assistantReply),
		Metadata:       cloneMetadata(metadata),
		ProcessedAt:    time.Now().UTC(),
	}
}

func deterministicFallbackReply(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "I could not generate a reply right now."
	}
	if len(trimmed) > 160 {
		trimmed = trimmed[:160] + "..."
	}
	return "I received your message: " + trimmed
}

func newOutputError(inputEventID string, correlationID string, metadata map[string]string, err error) natstransport.OutputMessage {
	s := err.Error()
	return natstransport.OutputMessage{
		InputEventID:  inputEventID,
		CorrelationID: correlationID,
		Status:        "error",
		Error:         &s,
		Metadata:      cloneMetadata(metadata),
		ProcessedAt:   time.Now().UTC(),
	}
}

func cloneMetadata(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func (s *Subscriber) shouldSendToDLQ(numDelivered uint64) bool {
	if s.cfg.DLQSubject == "" {
		return false
	}
	if s.cfg.MaxDeliver <= 0 {
		return false
	}
	return numDelivered >= uint64(s.cfg.MaxDeliver)
}

func (s *Subscriber) publishDLQ(input natstransport.InputMessage, out natstransport.OutputMessage, cause error, numDelivered uint64) error {
	payload := map[string]any{
		"input":          input,
		"output":         out,
		"error":          cause.Error(),
		"num_delivered":  numDelivered,
		"failed_at":      time.Now().UTC(),
		"failed_subject": s.cfg.OutputSubject,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal dlq payload: %w", err)
	}

	if _, err := s.js.Publish(s.cfg.DLQSubject, data); err != nil {
		return fmt.Errorf("publish dlq message: %w", err)
	}

	s.logger.Warn("nats: message moved to dlq",
		"event_id", input.ID,
		"correlation_id", input.CorrelationID,
		"dlq_subject", s.cfg.DLQSubject,
		"num_delivered", numDelivered,
	)

	return nil
}

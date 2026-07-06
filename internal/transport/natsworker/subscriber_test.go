package natsworker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	natspkg "github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/core/kernel"
	"github.com/tiroq/praxis/internal/storage/conversationstore"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

type fakeKernel struct {
	result kernel.PipelineResult
	err    error
}

func (f *fakeKernel) Run(_ context.Context, _ kernel.Event) (kernel.PipelineResult, error) {
	return f.result, f.err
}

type fakePublisher struct {
	published []natstransport.OutputMessage
	err       error
}

func (f *fakePublisher) Publish(msg natstransport.OutputMessage) error {
	if f.err != nil {
		return f.err
	}
	f.published = append(f.published, msg)
	return nil
}

type fakeMsg struct {
	data    []byte
	acked   bool
	naked   bool
	termed  bool
	subject string
	num     uint64
}

func (m *fakeMsg) Ack() error         { m.acked = true; return nil }
func (m *fakeMsg) Nak() error         { m.naked = true; return nil }
func (m *fakeMsg) Term() error        { m.termed = true; return nil }
func (m *fakeMsg) GetData() []byte    { return m.data }
func (m *fakeMsg) GetSubject() string { return m.subject }
func (m *fakeMsg) GetNumDelivered() uint64 {
	if m.num == 0 {
		return 1
	}
	return m.num
}

func handlerHarness(k KernelRunner, pub MessagePublisher) *Subscriber {
	return &Subscriber{
		kernel:    k,
		publisher: pub,
		logger:    noopLogger(),
		cfg:       natstransport.DefaultConfig(),
	}
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 100}))
}

func TestToKernelEvent(t *testing.T) {
	ts := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	event := toKernelEvent(natstransport.InputMessage{
		ID:            "evt_123",
		CorrelationID: "corr_123",
		Source:        "manual",
		Text:          "buy tickets",
		Timestamp:     ts,
		Metadata:      map[string]string{"key": "val"},
	})

	if event.ID != "evt_123" || event.CorrelationID != "corr_123" {
		t.Fatalf("event = %#v", event)
	}
	if event.Source != "manual" || event.Text != "buy tickets" {
		t.Fatalf("event = %#v", event)
	}
	if !event.OccurredAt.Equal(ts) || event.Metadata["key"] != "val" {
		t.Fatalf("event = %#v", event)
	}
}

func TestOutputBuilders(t *testing.T) {
	ok := newOutputOK("evt_1", "corr_1", map[string]string{"chat_id": "5"}, kernel.PipelineResult{EventID: "evt_1"}, "hello from llm")
	if ok.Status != "ok" || len(ok.Result) == 0 || ok.Error != nil {
		t.Fatalf("ok = %#v", ok)
	}
	if ok.CorrelationID != "corr_1" || ok.Metadata["chat_id"] != "5" {
		t.Fatalf("ok correlation/metadata = %#v", ok)
	}
	if ok.AssistantReply != "hello from llm" {
		t.Fatalf("ok assistant_reply = %#v", ok)
	}

	errOut := newOutputError("evt_2", "corr_2", map[string]string{"chat_id": "7"}, errors.New("boom"))
	if errOut.Status != "error" || errOut.Error == nil || *errOut.Error != "boom" {
		t.Fatalf("errOut = %#v", errOut)
	}
	if errOut.CorrelationID != "corr_2" || errOut.Metadata["chat_id"] != "7" {
		t.Fatalf("errOut correlation/metadata = %#v", errOut)
	}
}

func TestHandleMessageContracts(t *testing.T) {
	t.Run("invalid json terms", func(t *testing.T) {
		s := handlerHarness(&fakeKernel{}, &fakePublisher{})
		msg := &fakeMsg{data: []byte("{bad json")}
		s.handleMessage(context.Background(), msg)
		if !msg.termed || msg.acked || msg.naked {
			t.Fatalf("msg = %#v", msg)
		}
	})

	t.Run("validation terms", func(t *testing.T) {
		s := handlerHarness(&fakeKernel{}, &fakePublisher{})
		data, _ := json.Marshal(natstransport.InputMessage{ID: "evt_1"})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.termed {
			t.Fatalf("msg = %#v", msg)
		}
	})

	t.Run("publish failure naks", func(t *testing.T) {
		pub := &fakePublisher{err: errors.New("publish error")}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_ok"}}
		s := handlerHarness(k, pub)
		data, _ := json.Marshal(natstransport.InputMessage{ID: "evt_ok", Text: "buy tickets"})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.naked || msg.acked {
			t.Fatalf("msg = %#v", msg)
		}
	})

	t.Run("success acks", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_ok"}}
		s := handlerHarness(k, pub)
		data, _ := json.Marshal(natstransport.InputMessage{
			ID:            "evt_ok",
			CorrelationID: "corr_ok",
			Text:          "hello",
			Timestamp:     time.Now().UTC(),
			Metadata:      map[string]string{"chat_id": "99"},
		})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.acked || msg.naked || msg.termed || len(pub.published) != 1 {
			t.Fatalf("msg = %#v published=%d", msg, len(pub.published))
		}
		if pub.published[0].CorrelationID != "corr_ok" || pub.published[0].Metadata["chat_id"] != "99" {
			t.Fatalf("published = %#v", pub.published[0])
		}
	})

	t.Run("llm reply is invoked and published", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_llm_ok"}}
		s := handlerHarness(k, pub)

		called := false
		s.WithReplyGenerator(func(_ context.Context, input natstransport.InputMessage, _ kernel.PipelineResult) (string, error) {
			called = true
			if input.ID != "evt_llm_ok" {
				t.Fatalf("expected input id evt_llm_ok, got %s", input.ID)
			}
			return "LLM says hi", nil
		}, 2*time.Second)

		data, _ := json.Marshal(natstransport.InputMessage{
			ID:            "evt_llm_ok",
			CorrelationID: "corr_llm_ok",
			Text:          "hello",
			Timestamp:     time.Now().UTC(),
		})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)

		if !called {
			t.Fatal("expected llm reply generator to be called")
		}
		if len(pub.published) != 1 {
			t.Fatalf("expected one published output, got %d", len(pub.published))
		}
		if pub.published[0].AssistantReply != "LLM says hi" {
			t.Fatalf("expected assistant reply to be published, got %#v", pub.published[0])
		}
	})

	t.Run("llm provider error falls back deterministically", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_llm_err"}}
		s := handlerHarness(k, pub)

		s.WithReplyGenerator(func(_ context.Context, _ natstransport.InputMessage, _ kernel.PipelineResult) (string, error) {
			return "", errors.New("provider unavailable")
		}, 2*time.Second)

		data, _ := json.Marshal(natstransport.InputMessage{
			ID:            "evt_llm_err",
			CorrelationID: "corr_llm_err",
			Text:          "please summarize this",
			Timestamp:     time.Now().UTC(),
		})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)

		if len(pub.published) != 1 {
			t.Fatalf("expected one published output, got %d", len(pub.published))
		}
		if pub.published[0].AssistantReply == "" {
			t.Fatalf("expected fallback assistant reply, got %#v", pub.published[0])
		}
		if pub.published[0].InputEventID != "evt_llm_err" || pub.published[0].CorrelationID != "corr_llm_err" {
			t.Fatalf("expected ids preserved, got %#v", pub.published[0])
		}
	})

	t.Run("llm timeout falls back deterministically", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_llm_timeout"}}
		s := handlerHarness(k, pub)

		s.WithReplyGenerator(func(ctx context.Context, _ natstransport.InputMessage, _ kernel.PipelineResult) (string, error) {
			<-ctx.Done()
			return "", ctx.Err()
		}, 10*time.Millisecond)

		data, _ := json.Marshal(natstransport.InputMessage{
			ID:            "evt_llm_timeout",
			CorrelationID: "corr_llm_timeout",
			Text:          "timeout please",
			Timestamp:     time.Now().UTC(),
		})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)

		if len(pub.published) != 1 {
			t.Fatalf("expected one published output, got %d", len(pub.published))
		}
		if pub.published[0].AssistantReply == "" {
			t.Fatalf("expected timeout fallback assistant reply, got %#v", pub.published[0])
		}
		if pub.published[0].InputEventID != "evt_llm_timeout" || pub.published[0].CorrelationID != "corr_llm_timeout" {
			t.Fatalf("expected ids preserved, got %#v", pub.published[0])
		}
	})

	t.Run("missing correlation defaults through to output", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_missing_corr"}}
		s := handlerHarness(k, pub)
		data, _ := json.Marshal(natstransport.InputMessage{
			ID:       "evt_missing_corr",
			Text:     "hello",
			Metadata: map[string]string{"chat_id": "123"},
		})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.acked || len(pub.published) != 1 {
			t.Fatalf("msg = %#v published=%d", msg, len(pub.published))
		}
		if pub.published[0].CorrelationID != "evt_missing_corr" {
			t.Fatalf("expected fallback correlation_id, got %#v", pub.published[0])
		}
	})

	t.Run("kernel error publishes error output", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{err: errors.New("kernel blew up")}
		s := handlerHarness(k, pub)
		data, _ := json.Marshal(natstransport.InputMessage{ID: "evt_err", CorrelationID: "corr_err", Text: "something"})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.acked || len(pub.published) != 1 || pub.published[0].Status != "error" {
			t.Fatalf("msg = %#v published=%#v", msg, pub.published)
		}
		if pub.published[0].CorrelationID != "corr_err" {
			t.Fatalf("published = %#v", pub.published[0])
		}
	})
}

func TestToKernelEventDefaultsCorrelationID(t *testing.T) {
	event := toKernelEvent(natstransport.InputMessage{ID: "evt_no_corr", Text: "hello"})
	if event.CorrelationID != "evt_no_corr" {
		t.Fatalf("expected fallback correlation id, got %#v", event)
	}
}

func TestHandleMessage_PersistsConversationProjection(t *testing.T) {
	ctx := context.Background()

	store, err := conversationstore.OpenStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("OpenStore(:memory:) failed: %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("store.Close() failed: %v", closeErr)
		}
	}()

	pub := &fakePublisher{}
	k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_persist_1"}}
	s := handlerHarness(k, pub).WithConversationStore(store)
	s.WithReplyGenerator(func(_ context.Context, _ natstransport.InputMessage, _ kernel.PipelineResult) (string, error) {
		return "hello from llm reply", nil
	}, time.Second)

	data, _ := json.Marshal(natstransport.InputMessage{
		ID:            "evt_persist_1",
		CorrelationID: "telegram-chat-555",
		Source:        "telegram",
		Text:          "hello from telegram",
		Timestamp:     time.Date(2026, 7, 6, 10, 0, 0, 0, time.UTC),
		Metadata:      map[string]string{"chat_id": "555", "username": "alice"},
	})

	msg := &fakeMsg{data: data}
	s.handleMessage(ctx, msg)

	if !msg.acked || msg.naked || msg.termed {
		t.Fatalf("unexpected ack state: %#v", msg)
	}

	conv, err := store.GetConversationByCorrelationID(ctx, "telegram-chat-555")
	if err != nil {
		t.Fatalf("GetConversationByCorrelationID() failed: %v", err)
	}

	history, err := store.ListMessages(ctx, conv.ID, conversationstore.ListFilter{})
	if err != nil {
		t.Fatalf("ListMessages() failed: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 persisted projection messages (user+assistant), got %d", len(history))
	}

	roles := map[string]int{}
	for _, m := range history {
		roles[m.Role]++
		if m.EventID != "evt_persist_1" {
			t.Fatalf("expected projection message to preserve event id, got %q", m.EventID)
		}
		if m.Role == "assistant" && m.Content != "hello from llm reply" {
			t.Fatalf("expected assistant message content from llm reply, got %q", m.Content)
		}
	}
	if roles["user"] != 1 || roles["assistant"] != 1 {
		t.Fatalf("expected one user and one assistant message, got roles=%#v", roles)
	}
	if conv.CorrelationID != "telegram-chat-555" {
		t.Fatalf("expected preserved correlation id, got %q", conv.CorrelationID)
	}
}

var _ natsMsg = (*fakeMsg)(nil)

type fakePullSub struct {
	calls   []fetchCall
	callIdx int
}

type fetchCall struct {
	msgs []*natspkg.Msg
	err  error
}

func (f *fakePullSub) Fetch(_ int, _ ...natspkg.PullOpt) ([]*natspkg.Msg, error) {
	if f.callIdx >= len(f.calls) {
		return nil, natspkg.ErrTimeout
	}
	c := f.calls[f.callIdx]
	f.callIdx++
	return c.msgs, c.err
}

func (f *fakePullSub) Unsubscribe() error { return nil }

func TestRunLoop(t *testing.T) {
	t.Run("context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := handlerHarness(&fakeKernel{}, &fakePublisher{}).runLoop(ctx, &fakePullSub{}); err != nil {
			t.Fatalf("runLoop error = %v", err)
		}
	})

	t.Run("connection closed exits", func(t *testing.T) {
		fake := &fakePullSub{calls: []fetchCall{{err: natspkg.ErrConnectionClosed}}}
		if err := handlerHarness(&fakeKernel{}, &fakePublisher{}).runLoop(context.Background(), fake); err != nil {
			t.Fatalf("runLoop error = %v", err)
		}
	})

	t.Run("processes message", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_loop"}}
		s := handlerHarness(k, pub)
		input := natstransport.InputMessage{ID: "evt_loop", Text: "buy tickets", Timestamp: time.Now().UTC()}
		data, _ := json.Marshal(input)
		msg := &natspkg.Msg{Data: data, Subject: "praxis.kernel.input"}
		ctx, cancel := context.WithCancel(context.Background())
		fake := &fakePullSub{calls: []fetchCall{{msgs: []*natspkg.Msg{msg}}, {err: natspkg.ErrTimeout}}}
		go func() {
			time.Sleep(600 * time.Millisecond)
			cancel()
		}()
		if err := s.runLoop(ctx, fake); err != nil {
			t.Fatalf("runLoop error = %v", err)
		}
		if len(pub.published) != 1 {
			t.Fatalf("published = %d", len(pub.published))
		}
	})
}

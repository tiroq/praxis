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
}

func (m *fakeMsg) Ack() error         { m.acked = true; return nil }
func (m *fakeMsg) Nak() error         { m.naked = true; return nil }
func (m *fakeMsg) Term() error        { m.termed = true; return nil }
func (m *fakeMsg) GetData() []byte    { return m.data }
func (m *fakeMsg) GetSubject() string { return m.subject }

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
	ok := newOutputOK("evt_1", kernel.PipelineResult{EventID: "evt_1"})
	if ok.Status != "ok" || len(ok.Result) == 0 || ok.Error != nil {
		t.Fatalf("ok = %#v", ok)
	}

	errOut := newOutputError("evt_2", errors.New("boom"))
	if errOut.Status != "error" || errOut.Error == nil || *errOut.Error != "boom" {
		t.Fatalf("errOut = %#v", errOut)
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
		data, _ := json.Marshal(natstransport.InputMessage{ID: "evt_ok", Text: "hello", Timestamp: time.Now().UTC()})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.acked || msg.naked || msg.termed || len(pub.published) != 1 {
			t.Fatalf("msg = %#v published=%d", msg, len(pub.published))
		}
	})

	t.Run("kernel error publishes error output", func(t *testing.T) {
		pub := &fakePublisher{}
		k := &fakeKernel{err: errors.New("kernel blew up")}
		s := handlerHarness(k, pub)
		data, _ := json.Marshal(natstransport.InputMessage{ID: "evt_err", Text: "something"})
		msg := &fakeMsg{data: data}
		s.handleMessage(context.Background(), msg)
		if !msg.acked || len(pub.published) != 1 || pub.published[0].Status != "error" {
			t.Fatalf("msg = %#v published=%#v", msg, pub.published)
		}
	})
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

package praxiscli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

func TestPublish_GeneratesIDsAndPublishesInputMessage(t *testing.T) {
	fake := &fakeTransport{}
	app := newTestApp(fake)
	app.idFactory = func(prefix string) string { return prefix + "_test" }
	now := time.Date(2026, 7, 2, 0, 0, 0, 0, time.UTC)
	app.now = func() time.Time { return now }

	got, err := app.Publish(context.Background(), PublishRequest{Text: "Buy tickets to Shanghai"})
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if got.MessageID != "evt_test" || got.CorrelationID != "corr_test" {
		t.Fatalf("PublishedMessage = %#v", got)
	}

	var msg natstransport.InputMessage
	if err := json.Unmarshal(fake.publishedPayload, &msg); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if msg.ID != "evt_test" || msg.CorrelationID != "corr_test" {
		t.Fatalf("payload IDs = id:%q corr:%q", msg.ID, msg.CorrelationID)
	}
	if msg.Metadata[metadataCorrelationKey] != "corr_test" {
		t.Fatalf("metadata correlation_id = %q", msg.Metadata[metadataCorrelationKey])
	}
	if msg.Source != defaultSource {
		t.Fatalf("source = %q", msg.Source)
	}
	if !msg.Timestamp.Equal(now) {
		t.Fatalf("timestamp = %v", msg.Timestamp)
	}
}

func TestPublish_UsesProvidedValuesAndErrors(t *testing.T) {
	fake := &fakeTransport{}
	app := newTestApp(fake)

	_, err := app.Publish(context.Background(), PublishRequest{
		Text:          "test",
		MessageID:     "evt_given",
		CorrelationID: "corr_given",
		Source:        "manual",
		Metadata:      map[string]string{"x": "y"},
	})
	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}

	var msg natstransport.InputMessage
	if err := json.Unmarshal(fake.publishedPayload, &msg); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if msg.Source != "manual" || msg.Metadata["x"] != "y" {
		t.Fatalf("payload source/metadata invalid: %#v", msg)
	}

	if _, err := app.Publish(context.Background(), PublishRequest{Text: "   "}); err == nil || !strings.Contains(err.Error(), "publish text") {
		t.Fatalf("expected empty-text error, got %v", err)
	}

	app.factory = func(natstransport.Config) (Transport, error) { return nil, errors.New("dial") }
	if _, err := app.Publish(context.Background(), PublishRequest{Text: "ok"}); err == nil || !strings.Contains(err.Error(), "connect transport") {
		t.Fatalf("expected factory error, got %v", err)
	}

	app = newTestApp(&fakeTransport{publishErr: errors.New("boom")})
	if _, err := app.Publish(context.Background(), PublishRequest{Text: "ok"}); err == nil || !strings.Contains(err.Error(), "publish input message") {
		t.Fatalf("expected publish error, got %v", err)
	}
}

func TestWatch_SuccessAndErrorBranches(t *testing.T) {
	t.Run("success ack", func(t *testing.T) {
		d := &fakeDelivery{data: mustJSON(t, natstransport.OutputMessage{InputEventID: "evt_1", Status: "ok"})}
		app := newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{msg: d}}}})

		var out bytes.Buffer
		err := app.Watch(context.Background(), WatchRequest{Writer: &out, MaxMessages: 1, PollTimeout: time.Millisecond})
		if err != nil {
			t.Fatalf("Watch() error = %v", err)
		}
		if d.ackCount != 1 || d.nakCount != 0 {
			t.Fatalf("ack=%d nak=%d", d.ackCount, d.nakCount)
		}
		if !strings.Contains(out.String(), "\"input_event_id\": \"evt_1\"") {
			t.Fatalf("output = %q", out.String())
		}
	})

	t.Run("decode/write errors nak", func(t *testing.T) {
		dec := &fakeDelivery{data: []byte("not-json")}
		app := newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{msg: dec}}}})
		if err := app.Watch(context.Background(), WatchRequest{Writer: io.Discard, MaxMessages: 1, PollTimeout: time.Millisecond}); err == nil || !strings.Contains(err.Error(), "decode output message") {
			t.Fatalf("decode branch err = %v", err)
		}
		if dec.nakCount != 1 {
			t.Fatalf("decode branch nakCount = %d", dec.nakCount)
		}

		wr := &fakeDelivery{data: mustJSON(t, natstransport.OutputMessage{Status: "ok"})}
		app = newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{msg: wr}}}})
		if err := app.Watch(context.Background(), WatchRequest{Writer: errWriter{}, MaxMessages: 1, PollTimeout: time.Millisecond}); err == nil || !strings.Contains(err.Error(), "write output message") {
			t.Fatalf("write branch err = %v", err)
		}
		if wr.nakCount != 1 {
			t.Fatalf("write branch nakCount = %d", wr.nakCount)
		}
	})

	t.Run("subscribe/fetch errors", func(t *testing.T) {
		app := newTestApp(&fakeTransport{})
		if err := app.Watch(context.Background(), WatchRequest{Writer: io.Discard, MaxMessages: 1, PollTimeout: time.Millisecond}); err == nil || !strings.Contains(err.Error(), "subscribe output") {
			t.Fatalf("subscribe branch err = %v", err)
		}

		app = newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{err: errors.New("fetch boom")}}}})
		if err := app.Watch(context.Background(), WatchRequest{Writer: io.Discard, PollTimeout: time.Millisecond}); err == nil || !strings.Contains(err.Error(), "fetch output message") {
			t.Fatalf("fetch branch err = %v", err)
		}
	})

	t.Run("timeout connection context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		app := newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{err: ErrPollTimeout}, {err: ErrPollTimeout}}, afterEach: cancel}})
		if err := app.Watch(ctx, WatchRequest{Writer: io.Discard, PollTimeout: time.Millisecond}); err != nil {
			t.Fatalf("timeout branch err = %v", err)
		}

		app = newTestApp(&fakeTransport{sub: &fakeSubscription{next: []nextResult{{err: ErrConnectionClosed}}}})
		if err := app.Watch(context.Background(), WatchRequest{Writer: io.Discard, PollTimeout: time.Millisecond}); err != nil {
			t.Fatalf("connection branch err = %v", err)
		}
	})

	t.Run("writer required", func(t *testing.T) {
		app := newTestApp(&fakeTransport{})
		if err := app.Watch(context.Background(), WatchRequest{}); err == nil || !strings.Contains(err.Error(), "watch writer") {
			t.Fatalf("writer validation err = %v", err)
		}
	})
}

func TestHelpers(t *testing.T) {
	orig := map[string]string{"a": "1"}
	cloned := cloneMetadata(orig)
	cloned["a"] = "2"
	if orig["a"] != "1" {
		t.Fatalf("cloneMetadata mutated original")
	}

	id := newID("evt")
	if !strings.HasPrefix(id, "evt_") || len(id) <= len("evt_") {
		t.Fatalf("newID invalid: %q", id)
	}

	app := newTestApp(&fakeTransport{})
	err := app.processDelivery(io.Discard, &fakeDeliveryWithAckError{fakeDelivery{data: mustJSON(t, natstransport.OutputMessage{Status: "ok"})}})
	if err == nil || !strings.Contains(err.Error(), "ack output message") {
		t.Fatalf("processDelivery ack err = %v", err)
	}
}

func newTestApp(fake *fakeTransport) *App {
	app := NewApp(
		natstransport.DefaultConfig(),
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		func(natstransport.Config) (Transport, error) { return fake, nil },
	)
	app.idFactory = func(prefix string) string { return prefix + "_id" }
	app.now = func() time.Time { return time.Unix(0, 0).UTC() }
	return app
}

type fakeTransport struct {
	publishSubject   string
	publishedPayload []byte
	publishErr       error
	sub              Subscription
}

func (f *fakeTransport) PublishInput(subject string, payload []byte, _ string) error {
	f.publishSubject = subject
	f.publishedPayload = append([]byte(nil), payload...)
	return f.publishErr
}

func (f *fakeTransport) SubscribeOutput(_, _ string) (Subscription, error) {
	if f.sub == nil {
		return nil, errors.New("no subscription")
	}
	return f.sub, nil
}

func (f *fakeTransport) Close() {}

type nextResult struct {
	msg Delivery
	err error
}

type fakeSubscription struct {
	next      []nextResult
	idx       int
	afterEach func()
}

func (s *fakeSubscription) Next(_ time.Duration) (Delivery, error) {
	if s.idx >= len(s.next) {
		return nil, ErrPollTimeout
	}
	item := s.next[s.idx]
	s.idx++
	if s.afterEach != nil {
		s.afterEach()
	}
	return item.msg, item.err
}

func (s *fakeSubscription) Close() error { return nil }

type fakeDelivery struct {
	data     []byte
	ackCount int
	nakCount int
}

func (d *fakeDelivery) Data() []byte { return d.data }
func (d *fakeDelivery) Ack() error {
	d.ackCount++
	return nil
}
func (d *fakeDelivery) Nak() error {
	d.nakCount++
	return nil
}

type fakeDeliveryWithAckError struct{ fakeDelivery }

func (d *fakeDeliveryWithAckError) Ack() error {
	d.ackCount++
	return fmt.Errorf("ack boom")
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

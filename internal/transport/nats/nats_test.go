package nats

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	natspkg "github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/core/kernel"
)

// ---------------------------------------------------------------------------
// Fakes / stubs — no real NATS server required.
// ---------------------------------------------------------------------------

// fakeKernel is a KernelRunner stub for tests.
type fakeKernel struct {
	result kernel.PipelineResult
	err    error
}

func (f *fakeKernel) Run(_ context.Context, _ kernel.Event) (kernel.PipelineResult, error) {
	return f.result, f.err
}

// fakePublisher is a MessagePublisher stub that records calls.
type fakePublisher struct {
	published []OutputMessage
	err       error
}

func (f *fakePublisher) Publish(msg OutputMessage) error {
	if f.err != nil {
		return f.err
	}
	f.published = append(f.published, msg)
	return nil
}

// fakeMsg simulates a nats.Msg for ack/nak/term tracking.
// It implements the natsMsg interface defined in subscriber.go.
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

// ---------------------------------------------------------------------------
// Config tests
// ---------------------------------------------------------------------------

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.URL != "nats://localhost:4222" {
		t.Errorf("want URL nats://localhost:4222, got %q", cfg.URL)
	}
	if cfg.StreamName != "PRAXIS" {
		t.Errorf("want StreamName PRAXIS, got %q", cfg.StreamName)
	}
	if cfg.InputSubject != "praxis.kernel.input" {
		t.Errorf("want InputSubject praxis.kernel.input, got %q", cfg.InputSubject)
	}
	if cfg.OutputSubject != "praxis.kernel.output" {
		t.Errorf("want OutputSubject praxis.kernel.output, got %q", cfg.OutputSubject)
	}
	if cfg.Durable != "praxis-worker" {
		t.Errorf("want Durable praxis-worker, got %q", cfg.Durable)
	}
	if cfg.AckWait != 30*time.Second {
		t.Errorf("want AckWait 30s, got %v", cfg.AckWait)
	}
	if cfg.MaxDeliver != 3 {
		t.Errorf("want MaxDeliver 3, got %d", cfg.MaxDeliver)
	}
}

func TestConfigFromEnv_Defaults(t *testing.T) {
	// Ensure none of the env vars are set.
	for _, k := range []string{
		"NATS_URL", "NATS_STREAM", "NATS_INPUT_SUBJECT",
		"NATS_OUTPUT_SUBJECT", "NATS_DURABLE",
		"NATS_ACK_WAIT_SECONDS", "NATS_MAX_DELIVER",
	} {
		t.Setenv(k, "")
	}

	cfg := ConfigFromEnv()
	def := DefaultConfig()

	if cfg.URL != def.URL {
		t.Errorf("URL: want %q, got %q", def.URL, cfg.URL)
	}
	if cfg.StreamName != def.StreamName {
		t.Errorf("StreamName: want %q, got %q", def.StreamName, cfg.StreamName)
	}
	if cfg.AckWait != def.AckWait {
		t.Errorf("AckWait: want %v, got %v", def.AckWait, cfg.AckWait)
	}
	if cfg.MaxDeliver != def.MaxDeliver {
		t.Errorf("MaxDeliver: want %d, got %d", def.MaxDeliver, cfg.MaxDeliver)
	}
}

func TestConfigFromEnv_Overrides(t *testing.T) {
	os.Setenv("NATS_URL", "nats://remote:4222")
	os.Setenv("NATS_STREAM", "MYSTREAM")
	os.Setenv("NATS_INPUT_SUBJECT", "my.input")
	os.Setenv("NATS_OUTPUT_SUBJECT", "my.output")
	os.Setenv("NATS_DURABLE", "my-consumer")
	os.Setenv("NATS_ACK_WAIT_SECONDS", "60")
	os.Setenv("NATS_MAX_DELIVER", "5")
	t.Cleanup(func() {
		for _, k := range []string{
			"NATS_URL", "NATS_STREAM", "NATS_INPUT_SUBJECT",
			"NATS_OUTPUT_SUBJECT", "NATS_DURABLE",
			"NATS_ACK_WAIT_SECONDS", "NATS_MAX_DELIVER",
		} {
			os.Unsetenv(k)
		}
	})

	cfg := ConfigFromEnv()

	if cfg.URL != "nats://remote:4222" {
		t.Errorf("URL: got %q", cfg.URL)
	}
	if cfg.StreamName != "MYSTREAM" {
		t.Errorf("StreamName: got %q", cfg.StreamName)
	}
	if cfg.InputSubject != "my.input" {
		t.Errorf("InputSubject: got %q", cfg.InputSubject)
	}
	if cfg.OutputSubject != "my.output" {
		t.Errorf("OutputSubject: got %q", cfg.OutputSubject)
	}
	if cfg.Durable != "my-consumer" {
		t.Errorf("Durable: got %q", cfg.Durable)
	}
	if cfg.AckWait != 60*time.Second {
		t.Errorf("AckWait: got %v", cfg.AckWait)
	}
	if cfg.MaxDeliver != 5 {
		t.Errorf("MaxDeliver: got %d", cfg.MaxDeliver)
	}
}

// ---------------------------------------------------------------------------
// InputMessage decode / validate tests
// ---------------------------------------------------------------------------

func TestInputMessage_Decode_Valid(t *testing.T) {
	raw := `{
		"id": "evt_abc",
		"source": "manual",
		"text": "нужно купить билеты в Шанхай",
		"timestamp": "2026-07-01T00:00:00Z",
		"metadata": {}
	}`
	var msg InputMessage
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("unexpected decode error: %v", err)
	}
	if msg.ID != "evt_abc" {
		t.Errorf("ID: got %q", msg.ID)
	}
	if msg.Text != "нужно купить билеты в Шанхай" {
		t.Errorf("Text: got %q", msg.Text)
	}
}

func TestInputMessage_Decode_InvalidJSON(t *testing.T) {
	var msg InputMessage
	err := json.Unmarshal([]byte("{not json}"), &msg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInputMessage_Validate_MissingID(t *testing.T) {
	msg := InputMessage{Text: "some text"}
	if err := msg.validate(); !errors.Is(err, ErrEmptyMessageID) {
		t.Errorf("want ErrEmptyMessageID, got %v", err)
	}
}

func TestInputMessage_Validate_MissingText(t *testing.T) {
	msg := InputMessage{ID: "evt_1"}
	if err := msg.validate(); !errors.Is(err, ErrEmptyMessageText) {
		t.Errorf("want ErrEmptyMessageText, got %v", err)
	}
}

func TestInputMessage_Validate_Valid(t *testing.T) {
	msg := InputMessage{ID: "evt_1", Text: "buy tickets"}
	if err := msg.validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// toKernelEvent mapping tests
// ---------------------------------------------------------------------------

func TestInputMessage_toKernelEvent(t *testing.T) {
	ts := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	msg := InputMessage{
		ID:        "evt_123",
		Source:    "manual",
		Text:      "buy tickets",
		Timestamp: ts,
		Metadata:  map[string]string{"key": "val"},
	}
	event := msg.toKernelEvent()

	if event.ID != "evt_123" {
		t.Errorf("ID: got %q", event.ID)
	}
	if event.Text != "buy tickets" {
		t.Errorf("Text: got %q", event.Text)
	}
	if event.Source != "manual" {
		t.Errorf("Source: got %q", event.Source)
	}
	if !event.OccurredAt.Equal(ts) {
		t.Errorf("OccurredAt: got %v", event.OccurredAt)
	}
	if event.Confidence != 1.0 {
		t.Errorf("Confidence: got %f", event.Confidence)
	}
	if event.Metadata["key"] != "val" {
		t.Errorf("Metadata: got %v", event.Metadata)
	}
}

func TestInputMessage_toKernelEvent_ZeroTimestamp(t *testing.T) {
	msg := InputMessage{ID: "evt_x", Text: "hi"}
	before := time.Now().UTC()
	event := msg.toKernelEvent()
	after := time.Now().UTC()

	if event.OccurredAt.Before(before) || event.OccurredAt.After(after) {
		t.Errorf("OccurredAt not in expected range: %v", event.OccurredAt)
	}
}

// ---------------------------------------------------------------------------
// OutputMessage encoding tests
// ---------------------------------------------------------------------------

func TestNewOutputOK_Encoding(t *testing.T) {
	result := kernel.PipelineResult{EventID: "evt_1"}
	out := newOutputOK("evt_1", result)

	if out.Status != "ok" {
		t.Errorf("Status: got %q", out.Status)
	}
	if out.InputEventID != "evt_1" {
		t.Errorf("InputEventID: got %q", out.InputEventID)
	}
	if out.Result == nil {
		t.Error("Result must not be nil on ok")
	}
	if out.Error != nil {
		t.Errorf("Error must be nil on ok, got %v", *out.Error)
	}

	data, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var decoded OutputMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if decoded.Status != "ok" {
		t.Errorf("round-trip Status: got %q", decoded.Status)
	}
}

func TestNewOutputError_Encoding(t *testing.T) {
	out := newOutputError("evt_2", errors.New("something failed"))

	if out.Status != "error" {
		t.Errorf("Status: got %q", out.Status)
	}
	if out.Error == nil {
		t.Fatal("Error must not be nil on error output")
	}
	if *out.Error != "something failed" {
		t.Errorf("Error: got %q", *out.Error)
	}
	if out.Result != nil {
		t.Error("Result must be nil on error output")
	}
}

// ---------------------------------------------------------------------------
// Handler (Subscriber.handleMessage) ack/nak/term contract tests
// ---------------------------------------------------------------------------

// handlerHarness wires a Subscriber without a real JetStream context so that
// handleMessage can be tested in isolation.
func handlerHarness(k KernelRunner, pub MessagePublisher) *Subscriber {
	return &Subscriber{
		kernel:    k,
		publisher: pub,
		logger:    noopLogger(),
		cfg:       DefaultConfig(),
	}
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 100}))
}

func TestHandleMessage_InvalidJSON_Terms(t *testing.T) {
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	msg := &fakeMsg{data: []byte("{bad json")}
	s.handleMessage(context.Background(), msg)

	if !msg.termed {
		t.Error("expected Term() on invalid JSON")
	}
	if msg.acked || msg.naked {
		t.Error("must not ack or nak invalid JSON")
	}
}

func TestHandleMessage_MissingText_Terms(t *testing.T) {
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	input := InputMessage{ID: "evt_1"} // no Text
	data, _ := json.Marshal(input)
	msg := &fakeMsg{data: data}
	s.handleMessage(context.Background(), msg)

	if !msg.termed {
		t.Error("expected Term() on missing text")
	}
}

func TestHandleMessage_MissingID_Terms(t *testing.T) {
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	input := InputMessage{Text: "hello"}
	data, _ := json.Marshal(input)
	msg := &fakeMsg{data: data}
	s.handleMessage(context.Background(), msg)

	if !msg.termed {
		t.Error("expected Term() on missing id")
	}
}

func TestHandleMessage_PublishFails_Naks(t *testing.T) {
	pub := &fakePublisher{err: errors.New("publish error")}
	k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_ok"}}
	s := handlerHarness(k, pub)

	input := InputMessage{ID: "evt_ok", Text: "buy tickets"}
	data, _ := json.Marshal(input)
	msg := &fakeMsg{data: data}
	s.handleMessage(context.Background(), msg)

	if !msg.naked {
		t.Error("expected Nak() when publish fails")
	}
	if msg.acked {
		t.Error("must not Ack when publish fails")
	}
}

func TestHandleMessage_Success_AcksAfterPublish(t *testing.T) {
	pub := &fakePublisher{}
	k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_ok"}}
	s := handlerHarness(k, pub)

	input := InputMessage{ID: "evt_ok", Text: "нужно купить билеты в Шанхай", Timestamp: time.Now().UTC()}
	data, _ := json.Marshal(input)
	msg := &fakeMsg{data: data}
	s.handleMessage(context.Background(), msg)

	if !msg.acked {
		t.Error("expected Ack() after successful publish")
	}
	if msg.naked || msg.termed {
		t.Error("must not Nak or Term on success")
	}
	if len(pub.published) != 1 {
		t.Errorf("expected 1 published message, got %d", len(pub.published))
	}
	if pub.published[0].Status != "ok" {
		t.Errorf("expected status ok, got %q", pub.published[0].Status)
	}
}

func TestHandleMessage_KernelError_PublishesErrorOutput(t *testing.T) {
	pub := &fakePublisher{}
	k := &fakeKernel{err: errors.New("kernel blew up")}
	s := handlerHarness(k, pub)

	input := InputMessage{ID: "evt_err", Text: "something"}
	data, _ := json.Marshal(input)
	msg := &fakeMsg{data: data}
	s.handleMessage(context.Background(), msg)

	if !msg.acked {
		t.Error("expected Ack even when kernel fails (publish succeeded)")
	}
	if len(pub.published) != 1 {
		t.Fatalf("expected 1 published message, got %d", len(pub.published))
	}
	if pub.published[0].Status != "error" {
		t.Errorf("expected status error, got %q", pub.published[0].Status)
	}
}

// ---------------------------------------------------------------------------
// fakeMsg implements the natsMsg interface (defined in subscriber.go).
// Verified by compile: fakeMsg must satisfy natsMsg for the tests to build.
// ---------------------------------------------------------------------------
var _ natsMsg = (*fakeMsg)(nil)

// ---------------------------------------------------------------------------
// Publisher tests (using fake jsPublish — no real NATS required)
// ---------------------------------------------------------------------------

// fakeJSPublish is a minimal jsPublish stub.
type fakeJSPublish struct {
	published []struct {
		subj string
		data []byte
	}
	err error
}

func (f *fakeJSPublish) Publish(subj string, data []byte, _ ...natspkg.PubOpt) (*natspkg.PubAck, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.published = append(f.published, struct {
		subj string
		data []byte
	}{subj, data})
	return &natspkg.PubAck{}, nil
}

func TestNewPublisher_ReturnsNonNil(t *testing.T) {
	p := newPublisherFromInterface(&fakeJSPublish{}, "out.subject")
	if p == nil {
		t.Fatal("expected non-nil Publisher")
	}
}

func TestPublisher_Publish_Success(t *testing.T) {
	fake := &fakeJSPublish{}
	p := newPublisherFromInterface(fake, "praxis.kernel.output")
	out := newOutputOK("evt_1", kernel.PipelineResult{EventID: "evt_1"})
	if err := p.Publish(out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.published) != 1 {
		t.Fatalf("expected 1 publish call, got %d", len(fake.published))
	}
	if fake.published[0].subj != "praxis.kernel.output" {
		t.Errorf("subject: got %q", fake.published[0].subj)
	}
	var decoded OutputMessage
	if err := json.Unmarshal(fake.published[0].data, &decoded); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if decoded.Status != "ok" {
		t.Errorf("status: got %q", decoded.Status)
	}
}

func TestPublisher_Publish_Failure_ReturnsErrPublishFailed(t *testing.T) {
	fake := &fakeJSPublish{err: errors.New("broker down")}
	p := newPublisherFromInterface(fake, "out")
	out := newOutputError("evt_x", errors.New("boom"))
	err := p.Publish(out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrPublishFailed) {
		t.Errorf("want ErrPublishFailed, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// NewSubscriber constructor test
// ---------------------------------------------------------------------------

func TestNewSubscriber_ReturnsNonNil(t *testing.T) {
	s := NewSubscriber(nil, DefaultConfig(), &fakeKernel{}, &fakePublisher{}, noopLogger())
	if s == nil {
		t.Fatal("expected non-nil Subscriber")
	}
}

// ---------------------------------------------------------------------------
// natsMsgWrapper field passthrough tests
// ---------------------------------------------------------------------------

func TestNatsMsgWrapper_GetDataAndSubject(t *testing.T) {
	raw := []byte(`{"hello":"world"}`)
	m := &natspkg.Msg{Data: raw, Subject: "test.subject"}
	w := &natsMsgWrapper{m: m}

	if string(w.GetData()) != string(raw) {
		t.Errorf("GetData: got %q", w.GetData())
	}
	if w.GetSubject() != "test.subject" {
		t.Errorf("GetSubject: got %q", w.GetSubject())
	}
}

func TestNatsMsgWrapper_AckNakTermDoNotPanic(t *testing.T) {
	// Without a real subscription bound, Ack/Nak/Term return errors but must
	// not panic. This exercises the wrapper pass-through methods.
	m := &natspkg.Msg{}
	w := &natsMsgWrapper{m: m}
	_ = w.Ack()
	_ = w.Nak()
	_ = w.Term()
}

// ---------------------------------------------------------------------------
// ensureStream tests (using a fake streamManager — no real NATS required)
// ---------------------------------------------------------------------------

// fakeStreamManager implements streamManager for testing ensureStream.
type fakeStreamManager struct {
	streamExists  bool
	streamInfoErr error
	addStreamErr  error
	added         []*natspkg.StreamConfig
}

func (f *fakeStreamManager) StreamInfo(_ string, _ ...natspkg.JSOpt) (*natspkg.StreamInfo, error) {
	if f.streamExists {
		return &natspkg.StreamInfo{}, nil
	}
	if f.streamInfoErr != nil {
		return nil, f.streamInfoErr
	}
	return nil, natspkg.ErrStreamNotFound
}

func (f *fakeStreamManager) AddStream(cfg *natspkg.StreamConfig, _ ...natspkg.JSOpt) (*natspkg.StreamInfo, error) {
	if f.addStreamErr != nil {
		return nil, f.addStreamErr
	}
	f.added = append(f.added, cfg)
	return &natspkg.StreamInfo{}, nil
}

func TestEnsureStream_StreamAlreadyExists_NoAdd(t *testing.T) {
	fake := &fakeStreamManager{streamExists: true}
	cfg := DefaultConfig()
	if err := ensureStream(fake, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.added) != 0 {
		t.Errorf("expected no AddStream call, got %d", len(fake.added))
	}
}

func TestEnsureStream_StreamNotFound_Creates(t *testing.T) {
	fake := &fakeStreamManager{}
	cfg := DefaultConfig()
	if err := ensureStream(fake, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fake.added) != 1 {
		t.Fatalf("expected 1 AddStream call, got %d", len(fake.added))
	}
	if fake.added[0].Name != cfg.StreamName {
		t.Errorf("stream name: got %q", fake.added[0].Name)
	}
}

func TestEnsureStream_StreamInfoError_ReturnsError(t *testing.T) {
	fake := &fakeStreamManager{streamInfoErr: errors.New("server error")}
	cfg := DefaultConfig()
	err := ensureStream(fake, cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEnsureStream_AddStreamError_ReturnsError(t *testing.T) {
	fake := &fakeStreamManager{addStreamErr: errors.New("add failed")}
	cfg := DefaultConfig()
	err := ensureStream(fake, cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ---------------------------------------------------------------------------
// NewPublisher public constructor test
// ---------------------------------------------------------------------------

func TestNewPublisher_PublicConstructor_ReturnsNonNil(t *testing.T) {
	// NewPublisher stores the JetStreamContext without calling any methods,
	// so nil is safe to pass for a smoke-test of the constructor itself.
	p := NewPublisher(nil, "out.subject")
	if p == nil {
		t.Fatal("expected non-nil Publisher")
	}
}

// ---------------------------------------------------------------------------
// Client method tests (same package, no real NATS needed)
// ---------------------------------------------------------------------------

func TestClient_JetStream_ReturnsStoredField(t *testing.T) {
	c := &Client{js: nil}
	if c.JetStream() != nil {
		t.Error("expected nil JetStream when field is nil")
	}
}

func TestClient_Close_NilConn_DoesNotPanic(t *testing.T) {
	// Close with a nil conn must not panic; it checks before draining.
	c := &Client{conn: nil}
	c.Close() // must not panic
}

// ---------------------------------------------------------------------------
// runLoop tests (using fake pullSubscription — no real NATS required)
// ---------------------------------------------------------------------------

// fakePullSub implements pullSubscription for testing runLoop.
type fakePullSub struct {
	calls    []fetchCall
	callIdx  int
	unsubbed bool
}

type fetchCall struct {
	msgs []*natspkg.Msg
	err  error
}

func (f *fakePullSub) Fetch(_ int, _ ...natspkg.PullOpt) ([]*natspkg.Msg, error) {
	if f.callIdx >= len(f.calls) {
		// No more prepared responses — block by returning timeout
		return nil, natspkg.ErrTimeout
	}
	c := f.calls[f.callIdx]
	f.callIdx++
	return c.msgs, c.err
}

func (f *fakePullSub) Unsubscribe() error { f.unsubbed = true; return nil }

func TestRunLoop_ContextCancel_ExitsCleanly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	fake := &fakePullSub{}
	err := s.runLoop(ctx, fake)
	if err != nil {
		t.Errorf("expected nil error on ctx cancel, got %v", err)
	}
}

func TestRunLoop_TimeoutContinues(t *testing.T) {
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	fake := &fakePullSub{
		calls: []fetchCall{
			{err: natspkg.ErrTimeout},
		},
	}
	// Let one iteration happen (timeout) then cancel via context.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(600 * time.Millisecond)
		cancel()
	}()
	err := s.runLoop(ctx, fake)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunLoop_ConnectionClosed_ExitsCleanly(t *testing.T) {
	ctx := context.Background()
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	fake := &fakePullSub{
		calls: []fetchCall{
			{err: natspkg.ErrConnectionClosed},
		},
	}
	err := s.runLoop(ctx, fake)
	if err != nil {
		t.Errorf("expected nil on connection closed, got %v", err)
	}
}

func TestRunLoop_FetchError_Warns_ThenContinues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := handlerHarness(&fakeKernel{}, &fakePublisher{})
	fake := &fakePullSub{
		calls: []fetchCall{
			{err: errors.New("transient error")},
		},
	}
	// Cancel context after the loop handles the error and loops back.
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()
	err := s.runLoop(ctx, fake)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunLoop_ProcessesMessage_AcksAfterPublish(t *testing.T) {
	pub := &fakePublisher{}
	k := &fakeKernel{result: kernel.PipelineResult{EventID: "evt_loop"}}
	s := handlerHarness(k, pub)

	input := InputMessage{ID: "evt_loop", Text: "buy tickets", Timestamp: time.Now().UTC()}
	data, _ := json.Marshal(input)
	natMsg := &natspkg.Msg{Data: data, Subject: "praxis.kernel.input"}

	ctx, cancel := context.WithCancel(context.Background())
	fake := &fakePullSub{
		calls: []fetchCall{
			{msgs: []*natspkg.Msg{natMsg}},
			// After the message is processed, subsequent fetches time out so
			// the context cancel can stop the loop.
			{err: natspkg.ErrTimeout},
		},
	}
	go func() {
		time.Sleep(600 * time.Millisecond)
		cancel()
	}()
	err := s.runLoop(ctx, fake)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(pub.published) != 1 {
		t.Errorf("expected 1 published message, got %d", len(pub.published))
	}
}

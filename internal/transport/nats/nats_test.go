package nats

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	natspkg "github.com/nats-io/nats.go"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.URL != "nats://localhost:4222" {
		t.Fatalf("URL = %q", cfg.URL)
	}
	if cfg.StreamName != "PRAXIS" {
		t.Fatalf("StreamName = %q", cfg.StreamName)
	}
	if cfg.InputSubject != "praxis.kernel.input" {
		t.Fatalf("InputSubject = %q", cfg.InputSubject)
	}
	if cfg.OutputSubject != "praxis.kernel.output" {
		t.Fatalf("OutputSubject = %q", cfg.OutputSubject)
	}
	if cfg.DLQSubject != "praxis.kernel.dlq" {
		t.Fatalf("DLQSubject = %q", cfg.DLQSubject)
	}
	if cfg.Durable != "praxis-worker" {
		t.Fatalf("Durable = %q", cfg.Durable)
	}
	if cfg.AckWait != 30*time.Second {
		t.Fatalf("AckWait = %v", cfg.AckWait)
	}
	if cfg.MaxDeliver != 3 {
		t.Fatalf("MaxDeliver = %d", cfg.MaxDeliver)
	}
}

func TestConfigFromEnv(t *testing.T) {
	t.Setenv("NATS_URL", "nats://remote:4222")
	t.Setenv("NATS_STREAM", "MYSTREAM")
	t.Setenv("NATS_INPUT_SUBJECT", "my.input")
	t.Setenv("NATS_OUTPUT_SUBJECT", "my.output")
	t.Setenv("NATS_DLQ_SUBJECT", "my.dlq")
	t.Setenv("NATS_DURABLE", "my-consumer")
	t.Setenv("NATS_ACK_WAIT_SECONDS", "60")
	t.Setenv("NATS_MAX_DELIVER", "5")

	cfg := ConfigFromEnv()
	if cfg.URL != "nats://remote:4222" || cfg.StreamName != "MYSTREAM" {
		t.Fatalf("cfg = %#v", cfg)
	}
	if cfg.InputSubject != "my.input" || cfg.OutputSubject != "my.output" {
		t.Fatalf("cfg subjects = %#v", cfg)
	}
	if cfg.DLQSubject != "my.dlq" {
		t.Fatalf("cfg dlq = %#v", cfg)
	}
	if cfg.Durable != "my-consumer" || cfg.AckWait != 60*time.Second || cfg.MaxDeliver != 5 {
		t.Fatalf("cfg timing = %#v", cfg)
	}
}

func TestInputMessageValidate(t *testing.T) {
	if err := (InputMessage{Text: "x"}).Validate(); !errors.Is(err, ErrEmptyMessageID) {
		t.Fatalf("expected ErrEmptyMessageID, got %v", err)
	}
	if err := (InputMessage{ID: "evt"}).Validate(); !errors.Is(err, ErrEmptyMessageText) {
		t.Fatalf("expected ErrEmptyMessageText, got %v", err)
	}
	if err := (InputMessage{ID: "evt", Text: "ok"}).Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestOutputMessageRoundTrip(t *testing.T) {
	out := OutputMessage{
		InputEventID:  "evt_1",
		CorrelationID: "corr_1",
		Status:        "ok",
		Result:        json.RawMessage(`{"decision":{"id":"dec_1","outcome":"approve"},"actions":[{"id":"act_1"}]}`),
		Metadata:      map[string]string{"chat_id": "42"},
		ProcessedAt:   time.Now().UTC(),
	}

	data, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded OutputMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if decoded.Status != "ok" || len(decoded.Result) == 0 {
		t.Fatalf("decoded = %#v", decoded)
	}
	if decoded.CorrelationID != "corr_1" || decoded.Metadata["chat_id"] != "42" {
		t.Fatalf("decoded metadata = %#v", decoded)
	}
}

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
	}{subj: subj, data: append([]byte(nil), data...)})
	return &natspkg.PubAck{}, nil
}

func TestPublisherPublish(t *testing.T) {
	fake := &fakeJSPublish{}
	p := &Publisher{js: fake, subject: "praxis.kernel.output"}
	out := OutputMessage{InputEventID: "evt_1", Status: "ok", Result: json.RawMessage(`{"event_id":"evt_1"}`)}

	if err := p.Publish(out); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if len(fake.published) != 1 {
		t.Fatalf("publish count = %d", len(fake.published))
	}
	var decoded OutputMessage
	if err := json.Unmarshal(fake.published[0].data, &decoded); err != nil {
		t.Fatalf("decode published payload: %v", err)
	}
	if decoded.InputEventID != "evt_1" {
		t.Fatalf("decoded = %#v", decoded)
	}
}

func TestPublisherPublishFailure(t *testing.T) {
	p := &Publisher{js: &fakeJSPublish{err: errors.New("broker down")}, subject: "out"}
	err := p.Publish(OutputMessage{InputEventID: "evt_x", Status: "error"})
	if !errors.Is(err, ErrPublishFailed) {
		t.Fatalf("expected ErrPublishFailed, got %v", err)
	}
}

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

func TestEnsureStream(t *testing.T) {
	cfg := DefaultConfig()

	if err := ensureStream(&fakeStreamManager{streamExists: true}, cfg); err != nil {
		t.Fatalf("existing stream error = %v", err)
	}

	create := &fakeStreamManager{}
	if err := ensureStream(create, cfg); err != nil {
		t.Fatalf("create stream error = %v", err)
	}
	if len(create.added) != 1 || create.added[0].Name != cfg.StreamName {
		t.Fatalf("added = %#v", create.added)
	}

	if err := ensureStream(&fakeStreamManager{streamInfoErr: errors.New("server error")}, cfg); err == nil {
		t.Fatal("expected stream info error")
	}
	if err := ensureStream(&fakeStreamManager{addStreamErr: errors.New("add failed")}, cfg); err == nil {
		t.Fatal("expected add stream error")
	}
}

func TestClientHelpers(t *testing.T) {
	c := &Client{}
	if c.JetStream() != nil {
		t.Fatal("expected nil JetStream")
	}
	c.Close()
}

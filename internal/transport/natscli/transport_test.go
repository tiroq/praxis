package natscli

import (
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/cli/praxiscli"
)

type fakeClient struct {
	closed bool
}

func (f *fakeClient) Close() { f.closed = true }

type fakeJetStream struct {
	publishedSubject string
	publishedPayload []byte
	lastPublishOpts  int
	publishErr       error
	subscription     syncSubscription
	subscribeErr     error
	lastSubject      string
	lastSubOpts      int
}

func (f *fakeJetStream) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	f.publishedSubject = subj
	f.publishedPayload = append([]byte(nil), data...)
	f.lastPublishOpts = len(opts)
	if f.publishErr != nil {
		return nil, f.publishErr
	}
	return &nats.PubAck{}, nil
}

func (f *fakeJetStream) SubscribeSync(subj string, opts ...nats.SubOpt) (syncSubscription, error) {
	f.lastSubject = subj
	f.lastSubOpts = len(opts)
	if f.subscribeErr != nil {
		return nil, f.subscribeErr
	}
	return f.subscription, nil
}

type fakeSubscription struct {
	msg       transportMessage
	err       error
	closed    bool
	timeouts  []time.Duration
	nextCalls int
}

func (f *fakeSubscription) NextMsg(timeout time.Duration) (transportMessage, error) {
	f.timeouts = append(f.timeouts, timeout)
	f.nextCalls++
	if f.err != nil {
		return nil, f.err
	}
	return f.msg, nil
}

func (f *fakeSubscription) Unsubscribe() error {
	f.closed = true
	return nil
}

type fakeMessage struct {
	data   []byte
	ackErr error
	nakErr error
	acked  bool
	naked  bool
}

func (f *fakeMessage) Data() []byte { return f.data }

func (f *fakeMessage) Ack() error {
	f.acked = true
	return f.ackErr
}

func (f *fakeMessage) Nak() error {
	f.naked = true
	return f.nakErr
}

func TestTransportPublishInput(t *testing.T) {
	js := &fakeJetStream{}
	tr := &transport{client: &fakeClient{}, js: js}
	payload := []byte(`{"text":"hello"}`)

	if err := tr.PublishInput("praxis.kernel.input", payload, "evt_1"); err != nil {
		t.Fatalf("PublishInput() error = %v", err)
	}
	if js.publishedSubject != "praxis.kernel.input" || string(js.publishedPayload) != string(payload) {
		t.Fatalf("published subject/payload = %q %q", js.publishedSubject, string(js.publishedPayload))
	}
	if js.lastPublishOpts != 1 {
		t.Fatalf("publish opts = %d", js.lastPublishOpts)
	}
}

func TestTransportSubscribeOutput(t *testing.T) {
	fakeSub := &fakeSubscription{}
	js := &fakeJetStream{subscription: fakeSub}
	tr := &transport{client: &fakeClient{}, js: js}

	sub, err := tr.SubscribeOutput("praxis.kernel.output", "PRAXIS")
	if err != nil {
		t.Fatalf("SubscribeOutput() error = %v", err)
	}
	if js.lastSubject != "praxis.kernel.output" || js.lastSubOpts != 3 {
		t.Fatalf("subscribe call = subject:%q opts:%d", js.lastSubject, js.lastSubOpts)
	}
	if err := sub.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if !fakeSub.closed {
		t.Fatal("expected unsubscribe")
	}
}

func TestSubscriptionNextMapsErrorsAndReturnsDelivery(t *testing.T) {
	t.Run("timeout maps", func(t *testing.T) {
		sub := &subscription{sub: &fakeSubscription{err: nats.ErrTimeout}}
		if _, err := sub.Next(time.Second); !errors.Is(err, praxiscli.ErrPollTimeout) {
			t.Fatalf("expected ErrPollTimeout, got %v", err)
		}
	})

	t.Run("connection closed maps", func(t *testing.T) {
		sub := &subscription{sub: &fakeSubscription{err: nats.ErrConnectionClosed}}
		if _, err := sub.Next(time.Second); !errors.Is(err, praxiscli.ErrConnectionClosed) {
			t.Fatalf("expected ErrConnectionClosed, got %v", err)
		}
	})

	t.Run("success returns delivery", func(t *testing.T) {
		msg := &fakeMessage{data: []byte("payload")}
		fakeSub := &fakeSubscription{msg: msg}
		sub := &subscription{sub: fakeSub}
		d, err := sub.Next(250 * time.Millisecond)
		if err != nil {
			t.Fatalf("Next() error = %v", err)
		}
		if string(d.Data()) != "payload" {
			t.Fatalf("Data() = %q", string(d.Data()))
		}
		if err := d.Ack(); err != nil || !msg.acked {
			t.Fatalf("Ack() err=%v acked=%v", err, msg.acked)
		}
		if err := d.Nak(); err != nil || !msg.naked {
			t.Fatalf("Nak() err=%v naked=%v", err, msg.naked)
		}
		if fakeSub.nextCalls != 1 || fakeSub.timeouts[0] != 250*time.Millisecond {
			t.Fatalf("timeouts = %#v", fakeSub.timeouts)
		}
	})

	t.Run("other error passes through", func(t *testing.T) {
		want := errors.New("boom")
		sub := &subscription{sub: &fakeSubscription{err: want}}
		if _, err := sub.Next(time.Second); !errors.Is(err, want) {
			t.Fatalf("expected %v, got %v", want, err)
		}
	})
}

func TestTransportClose(t *testing.T) {
	client := &fakeClient{}
	tr := &transport{client: client, js: &fakeJetStream{}}
	tr.Close()
	if !client.closed {
		t.Fatal("expected client.Close()")
	}
}

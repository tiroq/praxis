package natscli

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/cli/praxiscli"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

// NewTransport builds a NATS-backed transport client for praxis CLI.
func NewTransport(cfg natstransport.Config) (praxiscli.Transport, error) {
	client, err := natstransport.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &transport{client: client, js: client.JetStream()}, nil
}

type transport struct {
	client *natstransport.Client
	js     nats.JetStreamContext
}

func (t *transport) PublishInput(subject string, payload []byte, msgID string) error {
	_, err := t.js.Publish(subject, payload, nats.MsgId(msgID))
	return err
}

func (t *transport) SubscribeOutput(subject, stream string) (praxiscli.Subscription, error) {
	sub, err := t.js.SubscribeSync(
		subject,
		nats.DeliverNew(),
		nats.AckExplicit(),
		nats.BindStream(stream),
	)
	if err != nil {
		return nil, err
	}
	return &subscription{sub: sub}, nil
}

func (t *transport) Close() {
	t.client.Close()
}

type subscription struct {
	sub *nats.Subscription
}

func (s *subscription) Next(timeout time.Duration) (praxiscli.Delivery, error) {
	msg, err := s.sub.NextMsg(timeout)
	if err != nil {
		if errors.Is(err, nats.ErrTimeout) {
			return nil, praxiscli.ErrPollTimeout
		}
		if errors.Is(err, nats.ErrConnectionClosed) {
			return nil, praxiscli.ErrConnectionClosed
		}
		return nil, err
	}
	return &delivery{msg: msg}, nil
}

func (s *subscription) Close() error {
	return s.sub.Unsubscribe()
}

type delivery struct {
	msg *nats.Msg
}

func (d *delivery) Data() []byte {
	return d.msg.Data
}

func (d *delivery) Ack() error {
	return d.msg.Ack()
}

func (d *delivery) Nak() error {
	return d.msg.Nak()
}

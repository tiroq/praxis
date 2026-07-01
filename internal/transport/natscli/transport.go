package natscli

import (
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tiroq/praxis/internal/cli/praxiscli"
	natstransport "github.com/tiroq/praxis/internal/transport/nats"
)

type transportClient interface {
	Close()
}

type jetStream interface {
	Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error)
	SubscribeSync(subj string, opts ...nats.SubOpt) (syncSubscription, error)
}

type syncSubscription interface {
	NextMsg(timeout time.Duration) (transportMessage, error)
	Unsubscribe() error
}

type transportMessage interface {
	Data() []byte
	Ack() error
	Nak() error
}

// NewTransport builds a NATS-backed transport client for praxis CLI.
func NewTransport(cfg natstransport.Config) (praxiscli.Transport, error) {
	client, err := natstransport.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &transport{client: client, js: jetStreamAdapter{js: client.JetStream()}}, nil
}

type transport struct {
	client transportClient
	js     jetStream
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
	sub syncSubscription
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
	msg transportMessage
}

func (d *delivery) Data() []byte {
	return d.msg.Data()
}

func (d *delivery) Ack() error {
	return d.msg.Ack()
}

func (d *delivery) Nak() error {
	return d.msg.Nak()
}

type jetStreamAdapter struct {
	js nats.JetStreamContext
}

func (a jetStreamAdapter) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	return a.js.Publish(subj, data, opts...)
}

func (a jetStreamAdapter) SubscribeSync(subj string, opts ...nats.SubOpt) (syncSubscription, error) {
	sub, err := a.js.SubscribeSync(subj, opts...)
	if err != nil {
		return nil, err
	}
	return subscriptionAdapter{sub: sub}, nil
}

type subscriptionAdapter struct {
	sub *nats.Subscription
}

func (a subscriptionAdapter) NextMsg(timeout time.Duration) (transportMessage, error) {
	msg, err := a.sub.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	return messageAdapter{msg: msg}, nil
}

func (a subscriptionAdapter) Unsubscribe() error {
	return a.sub.Unsubscribe()
}

type messageAdapter struct {
	msg *nats.Msg
}

func (a messageAdapter) Data() []byte {
	return a.msg.Data
}

func (a messageAdapter) Ack() error {
	return a.msg.Ack()
}

func (a messageAdapter) Nak() error {
	return a.msg.Nak()
}

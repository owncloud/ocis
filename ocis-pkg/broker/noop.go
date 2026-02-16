package broker

import (
	"errors"

	"go-micro.dev/v4/broker"
)

type NoOp struct{}

func (n NoOp) Init(_ ...broker.Option) error {
	return nil
}

func (n NoOp) Options() broker.Options {
	return broker.Options{}
}

func (n NoOp) Address() string {
	return ""
}

func (n NoOp) Connect() error {
	return nil
}

func (n NoOp) Disconnect() error {
	return nil
}

func (n NoOp) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	return nil
}

func (n NoOp) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, errors.New("not implemented")
}

func (n NoOp) String() string {
	return "NoOp"
}

func NewBroker(_ ...broker.Option) broker.Broker {
	return &NoOp{}
}

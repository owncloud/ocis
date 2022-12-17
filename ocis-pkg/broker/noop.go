package broker

import (
	"errors"

	"go-micro.dev/v4/broker"
)

// NoOp
// FIXME: nolint
// nolint: revive
type NoOp struct{}

// Init
// FIXME: nolint
// nolint: revive
func (n NoOp) Init(_ ...broker.Option) error {
	return nil
}

// Options
// FIXME: nolint
// nolint: revive
func (n NoOp) Options() broker.Options {
	return broker.Options{}
}

// Address
// FIXME: nolint
// nolint: revive
func (n NoOp) Address() string {
	return ""
}

// Connect
// FIXME: nolint
// nolint: revive
func (n NoOp) Connect() error {
	return nil
}

// Disconnect
// FIXME: nolint
// nolint: revive
func (n NoOp) Disconnect() error {
	return nil
}

// Publish
// FIXME: nolint
// nolint: revive
func (n NoOp) Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	return nil
}

// Subscribe
// FIXME: nolint
// nolint: revive
func (n NoOp) Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return nil, errors.New("not implemented")
}

// String implements string interface.
func (n NoOp) String() string {
	return "NoOp"
}

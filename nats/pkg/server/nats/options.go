package nats

import (
	natsServer "github.com/nats-io/nats-server/v2/server"
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

// Option configures the nats server
type Option func(*natsServer.Options, *stanServer.Options)

// Host sets the host URL for the nats server
func Host(url string) Option {
	return func(no *natsServer.Options, _ *stanServer.Options) {
		no.Host = url
	}
}

// Port sets the host URL for the nats server
func Port(port int) Option {
	return func(no *natsServer.Options, _ *stanServer.Options) {
		no.Port = port
	}
}

// NatsOpts allows setting Options from nats package directly
func NatsOpts(opt func(*natsServer.Options)) Option {
	return func(no *natsServer.Options, _ *stanServer.Options) {
		opt(no)
	}
}

// StanOpts allows setting Options from stan package directly
func StanOpts(opt func(*stanServer.Options)) Option {
	return func(_ *natsServer.Options, so *stanServer.Options) {
		opt(so)
	}
}

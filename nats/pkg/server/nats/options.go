package nats

import (
	natsServer "github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-streaming-server/logger"
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

// Port sets the host URL for the nats server
func Logger(logger logger.Logger) Option {
	return func(no *natsServer.Options, so *stanServer.Options) {
		so.CustomLogger = logger
	}
}

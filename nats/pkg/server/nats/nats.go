package nats

import (
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

// RunNatsServer runs the nats streaming server
func RunNatsServer(opts ...Option) (*stanServer.StanServer, error) {
	natsOpts := stanServer.DefaultNatsServerOptions
	stanOpts := stanServer.GetDefaultOptions()

	for _, o := range opts {
		o(&natsOpts, stanOpts)
	}
	s, err := stanServer.RunServerWithOpts(stanOpts, &natsOpts)
	return s, err
}

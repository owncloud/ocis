package nats

import (
	"context"
	// FIXME: nolint
	// nolint: nolintlint
	nserver "github.com/nats-io/nats-server/v2/server"
)

// NATSServer
// FIXME: nolint
// nolint: revive
type NATSServer struct {
	ctx    context.Context
	server *nserver.Server
}

// NewNATSServer initializes a new NATSServer instance.
func NewNATSServer(ctx context.Context, logger nserver.Logger, opts ...NatsOption) (*NATSServer, error) {
	natsOpts := &nserver.Options{}

	for _, o := range opts {
		o(natsOpts)
	}

	// enable JetStream
	natsOpts.JetStream = true

	server, err := nserver.NewServer(natsOpts)
	if err != nil {
		return nil, err
	}

	server.SetLoggerV2(logger, true, true, false)

	return &NATSServer{
		ctx:    ctx,
		server: server,
	}, nil
}

// ListenAndServe runs the NATSServer in a blocking way until the server is shutdown or an error occurs
func (n *NATSServer) ListenAndServe() (err error) {
	go n.server.Start()
	<-n.ctx.Done()
	return nil
}

// Shutdown
// FIXME: nolint
// nolint: revive
func (n *NATSServer) Shutdown() {
	n.server.Shutdown()
}

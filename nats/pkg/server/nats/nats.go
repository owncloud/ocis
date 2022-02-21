package nats

import (
	"context"
	"time"

	nserver "github.com/nats-io/nats-server/v2/server"
)

var NATSListenAndServeLoopTimer = 1 * time.Second

type NATSServer struct {
	ctx             context.Context
	jetStreamConfig *nserver.JetStreamConfig
	server          *nserver.Server
}

func NewNATSServer(ctx context.Context, logger nserver.Logger, opts ...Option) (*NATSServer, error) {
	options := &nserver.Options{}

	for _, o := range opts {
		o(options)
	}

	server, err := nserver.NewServer(
		options,
	)
	if err != nil {
		return nil, err
	}

	server.SetLoggerV2(logger, true, true, false)

	c := &nserver.JetStreamConfig{
		StoreDir: "/tmp/ocis-jetstream", // TODO: configurable
	}

	return &NATSServer{
		ctx:             ctx,
		jetStreamConfig: c,
		server:          server,
	}, nil
}

// ListenAndServe runs the NATSServer in a blocking way until the server is shutdown or an error occurs
func (n *NATSServer) ListenAndServe() (err error) {
	// start NATS first
	go n.server.Start()
	// start NATS JetStream second
	n.server.EnableJetStream(n.jetStreamConfig)
	if err != nil {
		return err
	}

	<-n.ctx.Done()
	return nil
}

func (n *NATSServer) Shutdown() {
	n.server.Shutdown()
}

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

func NewNATSServer(ctx context.Context, logger nserver.Logger, natsOpts []NatsOption, jetstreamOpts []JetStreamOption) (*NATSServer, error) {
	natsOptions := &nserver.Options{}
	jetStreamOptions := &nserver.JetStreamConfig{}

	for _, o := range natsOpts {
		o(natsOptions)
	}

	for _, o := range jetstreamOpts {
		o(jetStreamOptions)
	}

	server, err := nserver.NewServer(natsOptions)
	if err != nil {
		return nil, err
	}

	server.SetLoggerV2(logger, true, true, false)

	return &NATSServer{
		ctx:             ctx,
		jetStreamConfig: jetStreamOptions,
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

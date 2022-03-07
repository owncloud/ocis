package nats

import (
	"context"
	"time"

	natsServer "github.com/nats-io/nats-server/v2/server"
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

var NATSListenAndServeLoopTimer = 1 * time.Second

type NATSServer struct {
	ctx context.Context

	natsOpts *natsServer.Options
	stanOpts *stanServer.Options

	server *stanServer.StanServer
}

// NewNATSServer returns a new NATSServer
func NewNATSServer(ctx context.Context, opts ...Option) (*NATSServer, error) {

	server := &NATSServer{
		ctx:      ctx,
		natsOpts: &stanServer.DefaultNatsServerOptions,
		stanOpts: stanServer.GetDefaultOptions(),
	}

	for _, o := range opts {
		o(server.natsOpts, server.stanOpts)
	}

	return server, nil
}

// ListenAndServe runs the NATSServer in a blocking way until the server is shutdown or an error occurs
func (n *NATSServer) ListenAndServe() (err error) {
	n.server, err = stanServer.RunServerWithOpts(
		n.stanOpts,
		n.natsOpts,
	)
	if err != nil {
		return err
	}

	defer n.Shutdown()

	for {
		// check if NATs server has an encountered an error
		if err := n.server.LastError(); err != nil {
			return err
		}
		// check if the NATs server is still running
		if n.server.State() == stanServer.Shutdown {
			return nil
		}
		// check if context was cancelled
		if n.ctx.Err() != nil {
			return nil
		}
		time.Sleep(NATSListenAndServeLoopTimer)
	}
}

func (n *NATSServer) Shutdown() {
	n.server.Shutdown()
}

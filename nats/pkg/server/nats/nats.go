package nats

import (
	"time"

	natsServer "github.com/nats-io/nats-server/v2/server"
	stanServer "github.com/nats-io/nats-streaming-server/server"
)

type NATSServer struct {
	natsOpts *natsServer.Options
	stanOpts *stanServer.Options

	server *stanServer.StanServer
}

func NewNATSServer(opts ...Option) (*NATSServer, error) {
	server := &NATSServer{
		natsOpts: &stanServer.DefaultNatsServerOptions,
		stanOpts: stanServer.GetDefaultOptions(),
	}

	for _, o := range opts {
		o(server.natsOpts, server.stanOpts)
	}

	return server, nil
}

func (n *NATSServer) ListenAndServe() (err error) {

	n.server, err = stanServer.RunServerWithOpts(
		n.stanOpts,
		n.natsOpts,
	)
	if err != nil {
		return err
	}

	for {
		// check if NATs server has an encountered an error
		if err := n.server.LastError(); err != nil {
			return err
		}
		// check if th NATs server is still running
		if n.server.State() == stanServer.Shutdown {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (n *NATSServer) Shutdown() {
	n.server.Shutdown()
}

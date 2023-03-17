package http

import (
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/broker"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"

	mhttps "github.com/go-micro/plugins/v4/server/http"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/transport"
)

// Service simply wraps the go-micro web service.
type Service struct {
	micro.Service
}

// NewService initializes a new http service.
func NewService(opts ...Option) (Service, error) {
	noopBroker := broker.NoOp{}
	sopts := newOptions(opts...)
	mServer := mhttps.NewServer(server.TLSConfig(sopts.TLSConfig))

	wopts := []micro.Option{
		micro.Server(mServer),
		micro.Broker(noopBroker),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(registry.GetRegistry()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Transport(transport.NewHTTPTransport(transport.TLSConfig(sopts.TLSConfig))),
	}
	if sopts.TLSConfig != nil {
		// mark service in registry as using tls
		wopts = append(wopts, micro.Metadata(map[string]string{"use_tls": "true"}))
	}

	return Service{micro.NewService(wopts...)}, nil
}

package grpc

import (
	"strings"
	"time"

	mgrpcc "github.com/asim/go-micro/plugins/client/grpc/v3"
	mgrpcs "github.com/asim/go-micro/plugins/server/grpc/v3"
	"github.com/asim/go-micro/plugins/wrapper/monitoring/prometheus/v3"
	"github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3"
	"github.com/asim/go-micro/v3"
	"github.com/owncloud/ocis/ocis-pkg/registry"
)

// DefaultClient is a custom ocis grpc configured client.
var DefaultClient = mgrpcc.NewClient()

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)

	sopts.Logger.Info().
		Str("transport", "grpc").
		Str("addr", sopts.Address).
		Msg("starting server")

	mopts := []micro.Option{
		// first add a server because it will reset any options
		micro.Server(mgrpcs.NewServer()),
		// also add a client that can be used after initializing the service
		micro.Client(mgrpcc.NewClient()),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(*registry.GetRegistry()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
	}

	return Service{micro.NewService(mopts...)}
}

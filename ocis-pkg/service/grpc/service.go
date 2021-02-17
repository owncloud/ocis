package grpc

import (
	"strings"
	"time"

	grpcc "github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"

	"github.com/asim/go-micro/plugins/server/grpc/v3"

	"github.com/asim/go-micro/plugins/wrapper/trace/opencensus/v3"
	"github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/ocis-pkg/wrapper/prometheus"
)

// DefaultClient is a custom ocis grpc configured client.
var DefaultClient = newGrpcClient()

func newGrpcClient() client.Client {
	//r := *registry.GetRegistry()

	c := grpcc.NewClient(
	//grpcc.RequestTimeout(10*time.Second),
	//grpcc.Registry(r),
	)
	return c
}

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
		// This needs to be first as it replaces the underlying server
		// which causes any configuration set before it
		// to be discarded
		micro.Server(grpc.NewServer()),
		// TODO(refs) ideally we want to pass micro options from the consumers
		micro.Version(sopts.Version),
		micro.Address(sopts.Address),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Client(DefaultClient),
		micro.Registry(*registry.GetRegistry()),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
	}

	return Service{micro.NewService(mopts...)}
}

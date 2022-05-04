package grpc

import (
	"strings"
	"time"

	mgrpcc "github.com/go-micro/plugins/v4/client/grpc"
	mgrpcs "github.com/go-micro/plugins/v4/server/grpc"
	mbreaker "github.com/go-micro/plugins/v4/wrapper/breaker/gobreaker"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/go-micro/plugins/v4/wrapper/trace/opencensus"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

// DefaultClient is a custom oCIS grpc configured client.
var DefaultClient = getDefaultGrpcClient()

func getDefaultGrpcClient() client.Client {

	reg := registry.GetRegistry()

	return mgrpcc.NewClient(
		client.Registry(reg),
		client.Wrap(mbreaker.NewClientWrapper()),
	)
}

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)

	mopts := []micro.Option{
		// first add a server because it will reset any options
		micro.Server(mgrpcs.NewServer()),
		// also add a client that can be used after initializing the service
		micro.Client(DefaultClient),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(registry.GetRegistry()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
	}

	return Service{micro.NewService(mopts...)}
}

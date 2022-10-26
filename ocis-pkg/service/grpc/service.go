package grpc

import (
	"strings"
	"sync"
	"time"

	mgrpcc "github.com/go-micro/plugins/v4/client/grpc"
	mgrpcs "github.com/go-micro/plugins/v4/server/grpc"
	mbreaker "github.com/go-micro/plugins/v4/wrapper/breaker/gobreaker"
	"github.com/go-micro/plugins/v4/wrapper/monitoring/prometheus"
	"github.com/go-micro/plugins/v4/wrapper/trace/opencensus"
	oregistry "github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

// DefaultClient is a custom oCIS grpc configured client.
var (
	defaultClient client.Client
	once          sync.Once
)

func DefaultClient(registry oregistry.Registry) (client.Client, error) {
	return getDefaultGrpcClient(registry)
}

func getDefaultGrpcClient(registry oregistry.Registry) (client.Client, error) {
	reg, err := oregistry.GetRegistry(registry)
	if err != nil {
		return nil, err
	}

	once.Do(func() {
		defaultClient = mgrpcc.NewClient(
			client.Registry(reg),
			client.Wrap(mbreaker.NewClientWrapper()),
		)
	})
	return defaultClient, nil
}

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(registry oregistry.Registry, opts ...Option) (Service, error) {
	sopts := newOptions(opts...)

	client, err := DefaultClient(registry)
	if err != nil {
		return Service{}, err
	}
	reg, err := oregistry.GetRegistry(registry)
	if err != nil {
		return Service{}, err
	}

	mopts := []micro.Option{
		// first add a server because it will reset any options
		micro.Server(mgrpcs.NewServer()),
		// also add a client that can be used after initializing the service
		micro.Client(client),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
	}

	return Service{micro.NewService(mopts...)}, nil
}

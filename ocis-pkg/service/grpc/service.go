package grpc

import (
	"os"
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"

	etcdr "github.com/micro/go-micro/v2/registry/etcd"
	mdnsr "github.com/micro/go-micro/v2/registry/mdns"

	"github.com/micro/go-micro/v2/registry"

	"github.com/micro/go-plugins/wrapper/trace/opencensus/v2"
	"github.com/owncloud/ocis/ocis-pkg/wrapper/prometheus"
)

// DefaultClient is a custom ocis grpc configured client.
var DefaultClient = newGrpcClient()

func newGrpcClient() mclient.Client {
	addresses := strings.Split(os.Getenv("MICRO_REGISTRY_ADDRESS"), ",")

	var r registry.Registry
	switch os.Getenv("MICRO_REGISTRY") {
	case "etcd":
		r = etcdr.NewRegistry(
			registry.Addrs(addresses...),
		)
	default:
		r = mdnsr.NewRegistry()
	}

	c := grpc.NewClient(
		mclient.RequestTimeout(10*time.Second),
		mclient.Registry(r),
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
		micro.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		micro.Client(DefaultClient),
		micro.Version(sopts.Version),
		micro.Address(sopts.Address),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
	}

	return Service{
		micro.NewService(
			mopts...,
		),
	}
}

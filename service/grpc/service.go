package grpc

import (
	"strings"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/trace/opencensus"
	"github.com/owncloud/ocis-pkg/wrapper/prometheus"
)

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
		Msg("Starting server")

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
		micro.Version(sopts.Version),
		micro.Address(sopts.Address),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		// micro.Flags(sopts.Flags...),
	}

	return Service{
		micro.NewService(
			mopts...,
		),
	}
}

package http

import (
	"github.com/asim/go-micro/v3"
	svc "github.com/owncloud/ocis/graph-explorer/pkg/service/v0"
	"github.com/owncloud/ocis/graph-explorer/pkg/version"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Namespace),
		http.Name("graph-explorer"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			middleware.RealIP,
			middleware.RequestID,
			middleware.NoCache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version(
				"graph-explorer",
				version.String,
			),
			middleware.Logger(
				options.Logger,
			),
		),
	)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	service.Init()
	return service, nil
}

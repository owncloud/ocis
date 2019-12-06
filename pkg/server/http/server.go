package http

import (
	"github.com/owncloud/ocis-phoenix/pkg/assets"
	"github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis-phoenix/pkg/flagset"
	"github.com/owncloud/ocis-phoenix/pkg/service/v0"
	"github.com/owncloud/ocis-phoenix/pkg/version"
	"github.com/owncloud/ocis-pkg/middleware"
	"github.com/owncloud/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace("go.micro.web"),
		http.Name("phoenix"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(flagset.RootWithConfig(config.New())...),
		http.Flags(flagset.ServerWithConfig(config.New())...),
	)

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Cache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version(
				"phoenix",
				version.String,
			),
			middleware.Logger(
				options.Logger,
			),
			middleware.Static(
				assets.New(
					assets.Logger(options.Logger),
					assets.Config(options.Config),
				),
			),
		),
	)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	service.Handle(
		"/",
		handle,
	)

	service.Init()
	return service, nil
}

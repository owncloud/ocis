package http

import (
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/onlyoffice/pkg/assets"
	svc "github.com/owncloud/ocis/onlyoffice/pkg/service/v0"
	"github.com/owncloud/ocis/onlyoffice/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Name(options.Name),
		http.Logger(options.Logger),
		http.Version(version.String),
		http.Namespace(options.Config.HTTP.Namespace),
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
				options.Name,
				version.String,
			),
			middleware.Logger(
				options.Logger,
			),
			middleware.Static(
				options.Config.HTTP.Root,
				assets.New(
					assets.Logger(options.Logger),
					assets.Config(options.Config),
				),
				options.Config.HTTP.CacheTTL,
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
	http.M.Unlock()
	return service, nil
}

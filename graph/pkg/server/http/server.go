package http

import (
	svc "github.com/owncloud/ocis/graph/pkg/service/v0"
	"github.com/owncloud/ocis/graph/pkg/version"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("graph"),
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
				"graph",
				version.String,
			),
			middleware.Logger(
				options.Logger,
			),
			middleware.OpenIDConnect(
				oidc.Endpoint(options.Config.OpenIDConnect.Endpoint),
				oidc.Realm(options.Config.OpenIDConnect.Realm),
				oidc.Insecure(options.Config.OpenIDConnect.Insecure),
				oidc.Logger(options.Logger),
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

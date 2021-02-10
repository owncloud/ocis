package http

import (
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	webmid "github.com/owncloud/ocis/web/pkg/middleware"
	svc "github.com/owncloud/ocis/web/pkg/service/v0"
	"github.com/owncloud/ocis/web/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Namespace),
		http.Name("web"),
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
			webmid.SilentRefresh,
			middleware.Version(
				"web",
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

	service.Handle(
		"/",
		handle,
	)

	service.Init()
	http.M.Unlock()
	return service, nil
}

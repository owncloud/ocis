package http

import (
	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/service/http"
	"github.com/owncloud/ocis/onlyoffice/pkg/assets"
	svc "github.com/owncloud/ocis/onlyoffice/pkg/service/v0"
	"github.com/owncloud/ocis/onlyoffice/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Name("onlyoffice"),
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
			middleware.Cache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version(
				"onlyoffice",
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

	mux := chi.NewMux()

	mux.Use(middleware.Static(
		options.Config.HTTP.Root,
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {})

	service.Handle(
		"/",
		mux,
	)

	service.Init()
	return service, nil
}

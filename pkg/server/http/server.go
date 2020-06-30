package http

import (
	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-accounts/pkg/assets"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-accounts/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) http.Service {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Name(options.Name),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	handler, err := svc.New(svc.Logger(options.Logger), svc.Config(options.Config))
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("could not initialize service handler")
	}

	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Cache)
	mux.Use(middleware.Cors)
	mux.Use(middleware.Secure)

	mux.Use(middleware.Version(
		options.Name,
		version.String,
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Use(middleware.Static(
		options.Config.HTTP.Root,
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		proto.RegisterAccountsServiceWeb(r, handler)
		proto.RegisterGroupsServiceWeb(r, handler)
	})

	service.Handle(
		"/",
		mux,
	)

	if err := service.Init(); err != nil {
		panic(err)
	}
	return service
}

package http

import (
	"github.com/asim/go-micro/v3"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/settings/pkg/assets"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis/settings/pkg/service/v0"
	"github.com/owncloud/ocis/settings/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) http.Service {
	options := newOptions(opts...)

	service := http.NewService(
		http.Logger(options.Logger),
		http.Name(options.Name),
		http.Version(options.Config.Service.Version),
		http.Address(options.Config.HTTP.Addr),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	handle := svc.NewService(options.Config, options.Logger)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	mux := chi.NewMux()

	mux.Use(chimiddleware.RealIP)
	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.Cors)
	mux.Use(middleware.Secure)
	mux.Use(middleware.ExtractAccountUUID(
		account.Logger(options.Logger),
		account.JWTSecret(options.Config.TokenManager.JWTSecret)),
	)

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
		options.Config.HTTP.CacheTTL,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		proto.RegisterBundleServiceWeb(r, handle)
		proto.RegisterValueServiceWeb(r, handle)
		proto.RegisterRoleServiceWeb(r, handle)
		proto.RegisterPermissionServiceWeb(r, handle)
	})

	micro.RegisterHandler(service.Server(), mux)

	return service
}

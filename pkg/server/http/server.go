package http

import (
	"time"

	"github.com/go-chi/chi"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/assets"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-accounts/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/account"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/roles"
	"github.com/owncloud/ocis-pkg/v2/service/http"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
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

	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	rs := settings.NewRoleService("com.owncloud.api.settings", mclient.DefaultClient)
	roleManager := roles.NewManager(
		roles.CacheSize(1024),
		roles.CacheTTL(time.Hour*24*7),
		roles.Logger(options.Logger),
		roles.RoleService(rs),
	)
	handler, err := svc.New(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.RoleManager(&roleManager),
		svc.RoleService(rs),
	)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("could not initialize service handler")
	}

	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Cache)
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

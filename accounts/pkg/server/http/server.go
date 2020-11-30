package http

import (
	_ "net/http/pprof"

	ghttp "net/http"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis/accounts/pkg/assets"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/accounts/pkg/version"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) http.Service {
	options := newOptions(opts...)
	handler := options.Handler

	service := http.NewService(
		http.Logger(options.Logger),
		http.Name(options.Name),
		http.Version(options.Config.Server.Version),
		http.Address(options.Config.HTTP.Addr),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)

	mux := chi.NewMux()

	mux.Use(func(next ghttp.Handler) ghttp.Handler {
		return cmw.Profiler()
	})
	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
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

package http

import (
	"fmt"

	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svc "github.com/owncloud/ocis/v2/services/userlog/pkg/service"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	tlsConfig, err := config.BuildTLSConfig(
		options.Logger,
		options.Config.HTTP.TLS.Enabled,
		options.Config.HTTP.TLS.Cert,
		options.Config.HTTP.TLS.Key,
		options.Config.HTTP.Addr,
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("could not build certificate")
		return http.Service{}, fmt.Errorf("could not build certificate: %w", err)
	}

	service, err := http.NewService(
		http.TLSConfig(tlsConfig),
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("userlog"),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	middlewares := []func(stdhttp.Handler) stdhttp.Handler{
		middleware.TraceContext,
		chimiddleware.RequestID,
		middleware.Version(
			"userlog",
			version.GetString(),
		),
		middleware.Logger(
			options.Logger,
		),
		middleware.ExtractAccountUUID(
			account.Logger(options.Logger),
			account.JWTSecret(options.Config.TokenManager.JWTSecret),
		),
	}

	mux := chi.NewMux()
	mux.Use(middlewares...)

	handle, err := svc.NewUserlogService(
		svc.Logger(options.Logger),
		svc.Consumer(options.Consumer),
		svc.Mux(mux),
		svc.Store(options.Store),
		svc.Config(options.Config),
		svc.HistoryClient(options.HistoryClient),
		svc.GatewayClient(options.GatewayClient),
		svc.RegisteredEvents(options.RegisteredEvents),
	)
	if err != nil {
		return http.Service{}, err
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}

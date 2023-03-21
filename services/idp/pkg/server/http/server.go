package http

import (
	"fmt"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svc "github.com/owncloud/ocis/v2/services/idp/pkg/service/v0"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	tlsConfig, err := config.BuildTLSConfig(
		options.Logger,
		options.Config.HTTP.TLS,
		options.Config.HTTP.TLSCert,
		options.Config.HTTP.TLSKey,
		options.Config.HTTP.Addr,
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("could not build certificate")
		return http.Service{}, fmt.Errorf("could not build certificate: %w", err)
	}

	service, err := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name(options.Config.Service.Name),
		http.Version(version.GetString()),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
		http.TLSConfig(tlsConfig),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			chimiddleware.RealIP,
			chimiddleware.RequestID,
			middleware.TraceContext,
			middleware.NoCache,
			middleware.Secure,
			middleware.Version(
				options.Config.Service.Name,
				version.GetString(),
			),
			middleware.Logger(
				options.Logger,
			),
		),
	)

	{
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLoggingHandler(handle, options.Logger)
	}

	if err := micro.RegisterHandler(service.Server(), handle); err != nil {
		return http.Service{}, err
	}

	return service, nil
}

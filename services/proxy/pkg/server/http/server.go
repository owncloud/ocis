package http

import (
	"fmt"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	chain := options.Middlewares.Then(options.Handler)

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
		http.Name(options.Config.Service.Name),
		http.Version(version.GetString()),
		http.TLSConfig(tlsConfig),
		http.Logger(options.Logger),
		http.Address(options.Config.HTTP.Addr),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Context(options.Context),
		http.Flags(options.Flags...),
	)
	if err != nil {
		options.Logger.Error().
			Err(err).
			Msg("Error initializing http service")
		return http.Service{}, fmt.Errorf("could not initialize http service: %w", err)
	}

	if err := micro.RegisterHandler(service.Server(), chain); err != nil {
		return http.Service{}, err
	}

	return service, nil
}

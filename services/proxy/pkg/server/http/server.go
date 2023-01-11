package http

import (
	"fmt"
	"os"

	pkgcrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)
	l := options.Logger
	httpCfg := options.Config.HTTP

	if options.Config.HTTP.TLS {
		l.Warn().Msgf("No tls certificate provided, using a generated one")
		_, certErr := os.Stat(httpCfg.TLSCert)
		_, keyErr := os.Stat(httpCfg.TLSKey)

		if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
			// GenCert has side effects as it writes 2 files to the binary running location
			if err := pkgcrypto.GenCert(httpCfg.TLSCert, httpCfg.TLSKey, l); err != nil {
				l.Fatal().Err(err).Msgf("Could not generate test-certificate")
				os.Exit(1)
			}
		}
	}
	chain := options.Middlewares.Then(options.Handler)

	service, err := http.NewService(
		http.Name(options.Config.Service.Name),
		http.Version(version.GetString()),
		http.TLSConfig(shared.HTTPServiceTLS{
			Enabled: options.Config.HTTP.TLS,
			Cert:    options.Config.HTTP.TLSCert,
			Key:     options.Config.HTTP.TLSKey,
		}),
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

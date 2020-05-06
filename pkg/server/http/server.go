package http

import (
	"crypto/tls"
	"os"

	svc "github.com/owncloud/ocis-pkg/v2/service/http"
	"github.com/owncloud/ocis-proxy/pkg/crypto"
	"github.com/owncloud/ocis-proxy/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (svc.Service, error) {
	options := newOptions(opts...)
	l := options.Logger
	httpCfg := options.Config.HTTP

	var cer tls.Certificate
	var certErr error

	if httpCfg.TLSCert == "" || httpCfg.TLSKey == "" {
		l.Warn().Msgf("No tls certificate provided, using a generated one")

		// GenCert has side effects as it writes 2 files to the binary running location
		if err := crypto.GenCert(l); err != nil {
			l.Fatal().Err(err).Msgf("Could not generate test-certificate")
		}

		httpCfg.TLSCert = "server.crt"
		httpCfg.TLSKey = "server.key"
	}

	cer, certErr = tls.LoadX509KeyPair(httpCfg.TLSCert, httpCfg.TLSKey)

	if certErr != nil {
		options.Logger.Fatal().Err(certErr).Msg("Could not setup TLS")
		os.Exit(1)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
	chain := options.Middlewares.Then(options.Handler)

	service := svc.NewService(
		svc.Name("web.proxy"),
		svc.TLSConfig(tlsConfig),
		svc.Logger(options.Logger),
		svc.Namespace(options.Namespace),
		svc.Version(version.String),
		svc.Address(options.Config.HTTP.Addr),
		svc.Context(options.Context),
		svc.Flags(options.Flags...),
		svc.Handler(chain),
	)

	if err := service.Init(); err != nil {
		l.Fatal().Err(err).Msgf("Error initializing")
	}

	return service, nil
}

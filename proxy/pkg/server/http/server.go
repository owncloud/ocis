package http

import (
	"crypto/tls"
	"os"

	pkgcrypto "github.com/owncloud/ocis/ocis-pkg/crypto"
	svc "github.com/owncloud/ocis/ocis-pkg/service/http"
	"go-micro.dev/v4"
)

// Server initializes the http service and server.
func Server(opts ...Option) (svc.Service, error) {
	options := newOptions(opts...)
	l := options.Logger
	httpCfg := options.Config.HTTP

	var cer tls.Certificate

	var tlsConfig *tls.Config
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

		cer, certErr = tls.LoadX509KeyPair(httpCfg.TLSCert, httpCfg.TLSKey)
		if certErr != nil {
			options.Logger.Fatal().Err(certErr).Msg("Could not setup TLS")
			os.Exit(1)
		}

		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{cer}}
	}
	chain := options.Middlewares.Then(options.Handler)

	service := svc.NewService(
		svc.Name(options.Config.Service.Name),
		svc.TLSConfig(tlsConfig),
		svc.Logger(options.Logger),
		svc.Namespace(options.Config.Service.Namespace),
		svc.Version(options.Config.Service.Version),
		svc.Address(options.Config.HTTP.Addr),
		svc.Context(options.Context),
		svc.Flags(options.Flags...),
	)

	if err := micro.RegisterHandler(service.Server(), chain); err != nil {
		return svc.Service{}, err
	}

	return service, nil
}

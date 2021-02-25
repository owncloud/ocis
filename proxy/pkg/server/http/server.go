package http

import (
	"crypto/tls"
	"os"

	"github.com/asim/go-micro/v3"

	svc "github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/proxy/pkg/crypto"
)

// Server initializes the http service and server.
func Server(opts ...Option) (svc.Service, error) {
	options := newOptions(opts...)
	l := options.Logger
	httpCfg := options.Config.HTTP

	var cer tls.Certificate
	var certErr error

	var tlsConfig *tls.Config
	if options.Config.HTTP.TLS {
		if httpCfg.TLSCert == "" || httpCfg.TLSKey == "" {
			l.Warn().Msgf("No tls certificate provided, using a generated one")
			_, certErr := os.Stat("./server.crt")
			_, keyErr := os.Stat("./server.key")

			if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
				// GenCert has side effects as it writes 2 files to the binary running location
				if err := crypto.GenCert(l); err != nil {
					l.Fatal().Err(err).Msgf("Could not generate test-certificate")
					os.Exit(1)
				}
			}

			httpCfg.TLSCert = "server.crt"
			httpCfg.TLSKey = "server.key"
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

	micro.RegisterHandler(service.Server(), chain)

	service.Init()
	return service, nil
}

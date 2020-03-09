package http

import (
	"crypto/tls"
	"os"

	"github.com/owncloud/ocis-konnectd/pkg/crypto"
	svc "github.com/owncloud/ocis-konnectd/pkg/service/v0"
	"github.com/owncloud/ocis-konnectd/pkg/version"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	var tlsConfig *tls.Config
	if options.Config.HTTP.TLS {
		if options.Config.HTTP.TLSCert == "" || options.Config.HTTP.TLSKey == "" {
			_, certErr := os.Stat("./server.crt")
			_, keyErr := os.Stat("./server.key")

			if os.IsNotExist(certErr) || os.IsNotExist(keyErr) {
				options.Logger.Info().Msgf("Generating certs")
				if err := crypto.GenCert(options.Logger); err != nil {
					options.Logger.Fatal().Err(err).Msg("Could not setup TLS")
					os.Exit(1)
				}
			}

			options.Config.HTTP.TLSCert = "server.crt"
			options.Config.HTTP.TLSKey = "server.key"
		}

		cer, err := tls.LoadX509KeyPair(options.Config.HTTP.TLSCert, options.Config.HTTP.TLSKey)
		if err != nil {
			options.Logger.Fatal().Err(err).Msg("Could not setup TLS")
			os.Exit(1)
		}

		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	}

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("konnectd"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
		http.TLSConfig(tlsConfig),
	)

	options.Config.Konnectd.Listen = options.Config.HTTP.Addr
	handle := svc.NewService(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.Middleware(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Cache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version(
				"konnectd",
				version.String,
			),
			middleware.Logger(
				options.Logger,
			),
		),
	)

	{
		handle = svc.NewTracing(handle)
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
	}

	service.Handle(
		"/",
		handle,
	)

	service.Init()
	return service, nil
}

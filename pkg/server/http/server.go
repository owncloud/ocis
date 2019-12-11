package http

import (
	"crypto/tls"
	svc "github.com/owncloud/ocis-konnectd/pkg/service/v0"
	"github.com/owncloud/ocis-konnectd/pkg/version"
	"github.com/owncloud/ocis-pkg/middleware"
	"github.com/owncloud/ocis-pkg/service/http"
	"os"
)

// Server initializes the http service and server.
func Server(opts ...Option) (http.Service, error) {
	options := newOptions(opts...)

	cer, err := tls.LoadX509KeyPair(options.Config.HTTP.TLSCert, options.Config.HTTP.TLSKey)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("Could not setup TLS")
		os.Exit(1)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	service := http.NewService(
		http.Logger(options.Logger),
		http.Namespace(options.Config.HTTP.Namespace),
		http.Name("konnectd"),
		http.Version(version.String),
		http.Address(options.Config.HTTP.Addr),
		http.Context(options.Context),
		http.Flags(options.Flags...),
		http.TLSConfig(config),
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
		handle = svc.NewInstrument(handle, options.Metrics)
		handle = svc.NewLogging(handle, options.Logger)
		handle = svc.NewTracing(handle)
	}

	service.Handle(
		"/",
		handle,
	)

	service.Init()
	return service, nil
}

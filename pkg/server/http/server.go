package http

import (
	"crypto/tls"
	"os"

	occrypto "github.com/owncloud/ocis-konnectd/pkg/crypto"
	logger "github.com/owncloud/ocis-pkg/v2/log"
	svc "github.com/owncloud/ocis-pkg/v2/service/http"
	"github.com/owncloud/ocis-proxy/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (svc.Service, error) {
	options := newOptions(opts...)

	// GenCert has side effects as it writes 2 files to the binary running location
	occrypto.GenCert(logger.NewLogger())

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("Could not setup TLS")
		os.Exit(1)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	service := svc.NewService(
		svc.Name("web.proxy"),
		svc.TLSConfig(config),
		svc.Logger(options.Logger),
		svc.Namespace(options.Namespace),
		svc.Version(version.String),
		svc.Address(options.Config.HTTP.Addr),
		svc.Context(options.Context),
		svc.Flags(options.Flags...),
	)

	service.Init()

	return service, nil
}

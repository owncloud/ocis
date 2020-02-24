package http

import (
	"crypto/tls"
	"log"
	"net/http/httputil"
	"net/url"
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
		svc.Logger(options.Logger),
		svc.Namespace(options.Namespace),
		svc.Name("web.proxy"),
		svc.Version(version.String),
		svc.Address(options.Config.HTTP.Addr),
		svc.Context(options.Context),
		svc.TLSConfig(config),
		svc.Flags(options.Flags...),
	)

	phoenixURL, err := url.Parse("http://localhost:9100")
	if err != nil {
		log.Fatal(err)
	}
	konnectdURL, err := url.Parse("http://localhost:9130")
	if err != nil {
		log.Fatal(err)
	}
	revaURL, err := url.Parse("http://localhost:9140")
	if err != nil {
		log.Fatal(err)
	}

	service.Handle("/", httputil.NewSingleHostReverseProxy(phoenixURL))
	service.Handle("/.well-known/openid-configuration", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/jwks.json/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/signin/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/ocs/v1.php/", httputil.NewSingleHostReverseProxy(revaURL))
	service.Handle("/remote.php/webdav/", httputil.NewSingleHostReverseProxy(revaURL))

	service.Init()

	return service, nil
}

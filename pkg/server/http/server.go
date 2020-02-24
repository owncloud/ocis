package http

import (
	"log"
	"net/http/httputil"
	"net/url"

	svc "github.com/owncloud/ocis-pkg/v2/service/http"
	"github.com/owncloud/ocis-proxy/pkg/version"
)

// Server initializes the http service and server.
func Server(opts ...Option) (svc.Service, error) {
	options := newOptions(opts...)

	service := svc.NewService(
		svc.Logger(options.Logger),
		svc.Namespace(options.Namespace),
		svc.Name("web.proxy"),
		svc.Version(version.String),
		svc.Address(options.Config.HTTP.Addr),
		svc.Context(options.Context),
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

	service.Handle("/", httputil.NewSingleHostReverseProxy(phoenixURL))
	service.Handle("/.well-known/openid-configuration", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/jwks.json/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/token/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/userinfo/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/static/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/session/", httputil.NewSingleHostReverseProxy(konnectdURL))
	service.Handle("/konnect/v1/register/", httputil.NewSingleHostReverseProxy(konnectdURL))

	service.Init()

	return service, nil
}

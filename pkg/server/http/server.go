package http

import (
	"encoding/json"
	"net/http"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/memory"
	"github.com/micro/go-plugins/micro/router/v2"
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

	// TODO replace hardcoded adresses with service registry lookups
	routes := []router.Route{
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/.well-known/openid-configuration",
			},
			ProxyURL: router.URL{
				Scheme: "https",
				Host:   "localhost:9130",
				Path:   "/.well-known/openid-configuration",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/signin/v1",
			},
			ProxyURL: router.URL{
				Scheme: "https",
				Host:   "localhost:9130",
				Path:   "/signin/v1",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/konnect/v1",
			},
			ProxyURL: router.URL{
				Scheme: "https",
				Host:   "localhost:9130",
				Path:   "/konnect/v1",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/reva",
			},
			ProxyURL: router.URL{
				Scheme: "http",
				Host:   "localhost:9140",
				Path:   "/",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/graph/",
			},
			ProxyURL: router.URL{
				Scheme: "http",
				Host:   "localhost:9120",
				Path:   "/",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/graph-explorer/",
			},
			ProxyURL: router.URL{
				Scheme: "http",
				Host:   "localhost:9135",
				Path:   "/",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
		{
			Request: router.Request{
				Method: "GET",
				Host:   "localhost:9200",
				Path:   "/",
			},
			ProxyURL: router.URL{
				Scheme: "http",
				Host:   "localhost:9100",
				Path:   "/",
			},
			Weight: 2.0,
			Type:   "proxy",
		},
	}

	apiConfig := map[string]interface{}{
		"api": map[string]interface{}{
			"routes": routes,
		},
	}

	b, _ := json.Marshal(apiConfig)
	m := memory.NewSource(memory.WithJSON(b))
	conf, err := config.NewConfig(config.WithSource(m))
	if err != nil {
		options.Logger.Fatal().
			Err(err).
			Msg("could not parse routes")
	}

	r := router.NewRouter(router.Config(conf))

	wr := r.Handler()
	h := wr(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", 404)
	}))

	service.Handle(
		"/",
		h,
	)

	service.Init()
	return service, nil
}

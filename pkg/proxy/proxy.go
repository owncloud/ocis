package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-proxy/pkg/config"
)

// MultiHostReverseProxy extends httputil to support multiple hosts with diffent policies
type MultiHostReverseProxy struct {
	httputil.ReverseProxy
	Directors map[string]map[string]func(req *http.Request)
	logger    log.Logger
}

// NewMultiHostReverseProxy undocummented
func NewMultiHostReverseProxy(opts ...Option) *MultiHostReverseProxy {
	options := newOptions(opts...)

	reverseProxy := &MultiHostReverseProxy{
		Directors: make(map[string]map[string]func(req *http.Request)),
		logger:    options.Logger,
	}

	if options.Config.Policies == nil {
		reverseProxy.logger.Debug().Msg("config file not provided, using oCIS embedded set of redirects")
		options.Config.Policies = defaultPolicies()
	}

	for _, policy := range options.Config.Policies {
		for _, route := range policy.Routes {
			reverseProxy.logger.Debug().Str("fwd: ", route.Endpoint)
			uri, err := url.Parse(route.Backend)
			if err != nil {
				reverseProxy.logger.
					Fatal().
					Err(err).
					Msgf("malformed url: %v", route.Backend)
			}

			reverseProxy.logger.
				Debug().
				Interface("route", route).
				Msg("adding route")

			reverseProxy.AddHost(policy.Name, uri, route)
		}
	}

	return reverseProxy
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// AddHost undocumented
func (p *MultiHostReverseProxy) AddHost(policy string, target *url.URL, rt config.Route) {
	targetQuery := target.RawQuery
	if p.Directors[policy] == nil {
		p.Directors[policy] = make(map[string]func(req *http.Request))
	}
	p.Directors[policy][rt.Endpoint] = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		// Apache deployments host addresses need to match on req.Host and req.URL.Host
		// see https://stackoverflow.com/questions/34745654/golang-reverseproxy-with-apache2-sni-hostname-error
		if rt.ApacheVHost {
			req.Host = target.Host
		}

		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
}

func (p *MultiHostReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO need to fetch from the accounts service
	var hit bool
	policy := "reva"

	if _, ok := p.Directors[policy]; !ok {
		p.logger.
			Error().
			Msgf("policy %v is not configured", policy)
	}

	for k := range p.Directors[policy] {
		if strings.HasPrefix(r.URL.Path, k) && k != "/" {
			p.Director = p.Directors[policy][k]
			hit = true
			p.logger.
				Debug().
				Str("policy", policy).
				Str("prefix", k).
				Str("path", r.URL.Path).
				Msg("director found")
		}
	}

	// override default director with root. If any
	if !hit && p.Directors[policy]["/"] != nil {
		p.Director = p.Directors[policy]["/"]
	}

	// Call upstream ServeHTTP
	p.ReverseProxy.ServeHTTP(w, r)
}

func defaultPolicies() []config.Policy {
	return []config.Policy{
		config.Policy{
			Name: "reva",
			Routes: []config.Route{
				config.Route{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				config.Route{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint: "/ocs/",
					Backend:  "http://localhost:9140",
				},
				config.Route{
					Endpoint: "/remote.php/",
					Backend:  "http://localhost:9140",
				},
				config.Route{
					Endpoint: "/dav/",
					Backend:  "http://localhost:9140",
				},
				config.Route{
					Endpoint: "/webdav/",
					Backend:  "http://localhost:9140",
				},
			},
		},
		config.Policy{
			Name: "oc10",
			Routes: []config.Route{
				config.Route{
					Endpoint: "/",
					Backend:  "http://localhost:9100",
				},
				config.Route{
					Endpoint: "/.well-known/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint: "/konnect/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint: "/signin/",
					Backend:  "http://localhost:9130",
				},
				config.Route{
					Endpoint:    "/ocs/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				config.Route{
					Endpoint:    "/remote.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				config.Route{
					Endpoint:    "/dav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				config.Route{
					Endpoint:    "/webdav/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
			},
		},
	}
}

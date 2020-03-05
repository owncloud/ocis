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

	for _, policy := range options.Config.Policies {
		for _, route := range policy.Routes {
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

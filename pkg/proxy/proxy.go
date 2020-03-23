package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-proxy/pkg/config"
)

// MultiHostReverseProxy extends httputil to support multiple hosts with diffent policies
type MultiHostReverseProxy struct {
	httputil.ReverseProxy
	Directors map[string]map[config.RouteType]map[string]func(req *http.Request)
	logger    log.Logger
}

// NewMultiHostReverseProxy undocummented
func NewMultiHostReverseProxy(opts ...Option) *MultiHostReverseProxy {
	options := newOptions(opts...)

	rp := &MultiHostReverseProxy{
		Directors: make(map[string]map[config.RouteType]map[string]func(req *http.Request)),
		logger:    options.Logger,
	}

	if options.Config.Policies == nil {
		rp.logger.Info().Str("source", "runtime").Msg("Policies")
		options.Config.Policies = defaultPolicies()
	} else {
		rp.logger.Info().Str("source", "file").Msg("Policies")
	}

	for _, policy := range options.Config.Policies {
		for _, route := range policy.Routes {
			rp.logger.Debug().Str("fwd: ", route.Endpoint)
			uri, err := url.Parse(route.Backend)
			if err != nil {
				rp.logger.
					Fatal().
					Err(err).
					Msgf("malformed url: %v", route.Backend)
			}

			rp.logger.
				Debug().
				Interface("route", route).
				Msg("adding route")

			rp.AddHost(policy.Name, uri, route)
		}
	}

	return rp
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
		p.Directors[policy] = make(map[config.RouteType]map[string]func(req *http.Request))
	}
	routeType := config.DefaultRouteType
	if rt.Type != "" {
		routeType = rt.Type
	}
	if p.Directors[policy][routeType] == nil {
		p.Directors[policy][routeType] = make(map[string]func(req *http.Request))
	}
	p.Directors[policy][routeType][rt.Endpoint] = func(req *http.Request) {
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

	for i := 0; i < len(config.RouteTypes) && !hit; i++ {
		routeType := config.RouteTypes[i]
		var handler func(string, url.URL) bool
		switch routeType {
		case config.QueryRoute:
			handler = p.queryRouteHandler
		case config.RegexRoute:
			handler = p.regexRouteHandler
		case config.PrefixRoute:
			fallthrough
		default:
			handler = p.prefixRouteHandler
		}
		for endpoint := range p.Directors[policy][routeType] {
			if handler(endpoint, *r.URL) {
				p.Director = p.Directors[policy][routeType][endpoint]
				hit = true
				p.logger.
					Info().
					Str("policy", policy).
					Str("prefix", endpoint).
					Str("path", r.URL.Path).
					Str("routeType", string(routeType)).
					Msg("director found")
				break
			}
		}
	}

	// override default director with root. If any
	if !hit && p.Directors[policy][config.PrefixRoute]["/"] != nil {
		p.Director = p.Directors[policy][config.PrefixRoute]["/"]
	}

	// Call upstream ServeHTTP
	p.ReverseProxy.ServeHTTP(w, r)
}

func (p MultiHostReverseProxy) queryRouteHandler(endpoint string, target url.URL) bool {
	u, _ := url.Parse(endpoint)
	if strings.HasPrefix(target.Path, u.Path) && endpoint != "/" {
		query := u.Query()
		if len(query) != 0 {
			rQuery := target.Query()
			match := true
			for k := range query {
				v := query.Get(k)
				rv := rQuery.Get(k)
				if rv != v {
					match = false
					break
				}
			}
			return match
		}
	}
	return false
}

func (p *MultiHostReverseProxy) regexRouteHandler(endpoint string, target url.URL) bool {
	matched, err := regexp.MatchString(endpoint, target.String())
	if err != nil {
		p.logger.Warn().Err(err).Msgf("regex with pattern %s failed", endpoint)
	}
	return matched
}

func (p *MultiHostReverseProxy) prefixRouteHandler(endpoint string, target url.URL) bool {
	return strings.HasPrefix(target.Path, endpoint) && endpoint != "/"
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
					Type:     config.QueryRoute,
					Endpoint: "/remote.php/?preview=1",
					Backend:  "http://localhost:9115",
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
				config.Route{
					Endpoint: "/status.php",
					Backend:  "http://localhost:9140",
				},
				config.Route{
					Endpoint: "/index.php/",
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
				config.Route{
					Endpoint:    "/status.php",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
				config.Route{
					Endpoint:    "/index.php/",
					Backend:     "https://demo.owncloud.com",
					ApacheVHost: true,
				},
			},
		},
	}
}

package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy/policy"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
	"github.com/rs/zerolog"
)

// IdleConnTimeout is the time after which the proxy closes an idle connection
// of an outgoing HTTP client. Keep it lower than the idle timeout of every
// server and network device on the way. The usual lowest value is 60 seconds.
const IdleConnTimeout = 30 * time.Second

// MultiHostReverseProxy extends "httputil" to support multiple hosts with different policies
type MultiHostReverseProxy struct {
	httputil.ReverseProxy
	// Directors holds policy        route type        method    endpoint         Director
	Directors      map[string]map[config.RouteType]map[string]map[string]func(req *http.Request)
	PolicySelector policy.Selector
	logger         log.Logger
	config         *config.Config
}

// NewMultiHostReverseProxy creates a new MultiHostReverseProxy
func NewMultiHostReverseProxy(opts ...Option) (*MultiHostReverseProxy, error) {
	options := newOptions(opts...)

	rp := &MultiHostReverseProxy{
		ReverseProxy: httputil.ReverseProxy{
			ErrorLog: stdlog.New(options.Logger, "", 0),
		},
		Directors: make(map[string]map[config.RouteType]map[string]map[string]func(req *http.Request)),
		logger:    options.Logger,
		config:    options.Config,
	}

	rp.Rewrite = func(r *httputil.ProxyRequest) {
		ri := router.ContextRoutingInfo(r.In.Context())
		ri.Rewrite()(r)
	}

	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		reqLogger := zerolog.Ctx(req.Context())
		if ev := reqLogger.Error(); ev.Enabled() {
			ev.Err(err).Msg("error happened in MultiHostReverseProxy")
		} else {
			rp.logger.Err(err).Msg("error happened in MultiHostReverseProxy")
		}
		rw.WriteHeader(http.StatusBadGateway)
	}

	tlsConf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: options.Config.InsecureBackends, //nolint:gosec
	}
	if options.Config.BackendHTTPSCACert != "" {
		certs := x509.NewCertPool()
		pemData, err := os.ReadFile(options.Config.BackendHTTPSCACert)
		if err != nil {
			return nil, err
		}
		if !certs.AppendCertsFromPEM(pemData) {
			return nil, errors.New("Error initializing LDAP Backend. Adding CA cert failed")
		}
		tlsConf.RootCAs = certs
	}
	rp.Transport = middleware.GetOtelhttpClientTransport(newBackendTransport(tlsConf), options.TraceProvider)
	return rp, nil
}

func (p *MultiHostReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ReverseProxy.ServeHTTP(w, r)
}

// newBackendTransport returns the transport for the requests to the backends.
// It equals http.DefaultTransport except TLSClientConfig, the idle connection
// limits and IdleConnTimeout.
func newBackendTransport(tlsConf *tls.Config) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2: true,
		// Zero means no limit for all hosts together. In a deployment with one
		// pod for each service replica, each pod is one host. A global limit of
		// 100 would leave only a few idle connections for each pod.
		MaxIdleConns: 0,
		// The default is 2 idle connections for each backend. The backends use
		// HTTP/1.1, so each request needs its own connection. With more than 2
		// parallel requests, the transport closes the extra connections after the
		// request and opens new ones for the next requests. This wastes ports.
		MaxIdleConnsPerHost: 100,
		// The oCIS HTTP servers set no idle timeout. Load balancers and ingress
		// controllers between the proxy and a backend often close idle connections
		// after 60 seconds. The proxy closes its idle connections first. If it
		// did not, it would send a request on a closed connection. The transport
		// does not repeat such a request for methods like PUT, POST or PROPFIND.
		IdleConnTimeout:       IdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConf,
	}
}

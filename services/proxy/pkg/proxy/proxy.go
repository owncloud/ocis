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
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy/policy"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
	"github.com/rs/zerolog"
)

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
	// equals http.DefaultTransport except TLSClientConfig
	rp.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConf,
	}
	return rp, nil
}

func (p *MultiHostReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ReverseProxy.ServeHTTP(w, r)
}

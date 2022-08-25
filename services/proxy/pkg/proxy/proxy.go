package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go.opentelemetry.io/otel/attribute"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	pkgtrace "github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy/policy"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
	proxytracing "github.com/owncloud/ocis/v2/services/proxy/pkg/tracing"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
func NewMultiHostReverseProxy(opts ...Option) *MultiHostReverseProxy {
	options := newOptions(opts...)

	rp := &MultiHostReverseProxy{
		Directors: make(map[string]map[config.RouteType]map[string]map[string]func(req *http.Request)),
		logger:    options.Logger,
		config:    options.Config,
	}

	rp.Director = func(r *http.Request) {
		ri := router.ContextRoutingInfo(r.Context())
		ri.Director()(r)
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
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.Config.InsecureBackends, //nolint:gosec
		},
	}
	return rp
}

func (p *MultiHostReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		span trace.Span
	)

	tracer := proxytracing.TraceProvider.Tracer("proxy")
	ctx, span = tracer.Start(ctx, fmt.Sprintf("%s %v", r.Method, r.URL.Path))
	defer span.End()

	span.SetAttributes(
		attribute.KeyValue{
			Key:   "x-request-id",
			Value: attribute.StringValue(chimiddleware.GetReqID(r.Context())),
		})

	pkgtrace.Propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	p.ReverseProxy.ServeHTTP(w, r.WithContext(ctx))
}

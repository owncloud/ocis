package debug

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/zpages"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	graphMiddleware "github.com/owncloud/ocis/v2/services/graph/pkg/middleware"
)

// NewService initializes a new debug service.
func NewService(opts ...Option) *http.Server {
	dopts := newOptions(opts...)
	mux := http.NewServeMux()

	mux.Handle("/metrics", alice.New(
		graphMiddleware.Token(
			dopts.Token,
		),
	).Then(
		promhttp.Handler(),
	))

	if dopts.Health != nil {
		mux.Handle("/healthz", dopts.Health)
	}

	if dopts.Ready != nil {
		mux.Handle("/readyz", dopts.Ready)
	}

	if dopts.ConfigDump != nil {
		mux.Handle("/config", dopts.ConfigDump)
	}

	if dopts.Pprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	if dopts.Zpages {
		h := zpages.NewTracezHandler(zpages.NewSpanProcessor())
		mux.Handle("/debug", h)
	}

	baseCtx := dopts.Context
	if baseCtx == nil {
		baseCtx = context.Background()
	}

	return &http.Server{
		Addr: dopts.Address,
		BaseContext: func(_ net.Listener) context.Context {
			return baseCtx
		},
		Handler: alice.New(
			chimiddleware.RealIP,
			chimiddleware.RequestID,
			middleware.NoCache,
			middleware.Cors(
				cors.AllowedOrigins(dopts.CorsAllowedOrigins),
				cors.AllowedMethods(dopts.CorsAllowedMethods),
				cors.AllowedHeaders(dopts.CorsAllowedHeaders),
				cors.AllowCredentials(dopts.CorsAllowCredentials),
			),
			middleware.Version(
				dopts.Name,
				dopts.Version,
			),
		).Then(
			mux,
		),
	}
}

package debug

import (
	"net/http"
	"net/http/pprof"

	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/justinas/alice"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/zpages"
)

// NewService initializes a new debug service.
func NewService(opts ...Option) *http.Server {
	dopts := newOptions(opts...)

	dopts.Logger.Info().
		Str("transport", "debug").
		Str("addr", dopts.Address).
		Msg("starting server")

	mux := http.NewServeMux()

	mux.Handle("/metrics", alice.New(
		middleware.Token(
			dopts.Token,
		),
	).Then(
		promhttp.Handler(),
	))

	mux.HandleFunc("/healthz", dopts.Health)
	mux.HandleFunc("/readyz", dopts.Ready)

	if dopts.Pprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	if dopts.Zpages {
		zpages.Handle(mux, "/debug")
	}

	return &http.Server{
		Addr: dopts.Address,
		Handler: alice.New(
			chimiddleware.RealIP,
			chimiddleware.RequestID,
			middleware.NoCache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version(
				dopts.Name,
				dopts.Version,
			),
		).Then(
			mux,
		),
	}
}

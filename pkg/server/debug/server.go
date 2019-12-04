package debug

import (
	"io"
	"net/http"
	"net/http/pprof"

	"github.com/justinas/alice"
	"github.com/micro/go-micro/util/log"
	"github.com/owncloud/ocis-graph/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/zpages"
)

type Handler struct {
	*http.ServeMux
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, http.StatusText(http.StatusOK))
}

func (h *Handler) readyz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, http.StatusText(http.StatusOK))
}

func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)
	log.Infof("Server [debug] listening on [%s]", options.Config.Debug.Addr)

	mux := http.NewServeMux()
	handler := &Handler{mux}

	handler.Handle("/metrics", alice.New(
		middleware.Token(
			options.Config.Debug.Token,
		),
	).Then(
		promhttp.Handler(),
	))

	handler.HandleFunc("/healthz", handler.healthz)
	handler.HandleFunc("/readyz", handler.readyz)

	if options.Config.Debug.Pprof {
		handler.HandleFunc("/debug/pprof/", pprof.Index)
		handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
		handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	if options.Config.Debug.Zpages {
		zpages.Handle(mux, "/debug")
	}

	return &http.Server{
		Addr: options.Config.Debug.Addr,
		Handler: alice.New(
			middleware.RealIP,
			middleware.RequestID,
			middleware.Cache,
			middleware.Cors,
			middleware.Secure,
			middleware.Version,
		).Then(
			handler,
		),
	}, nil
}

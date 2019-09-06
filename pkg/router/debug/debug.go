package debug

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis-phoenix/pkg/handler/metrics"
	"github.com/owncloud/ocis-phoenix/pkg/middleware/header"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// debug gets initialized by Router and configures the router.
type debug struct {
	token string
	pprof bool
}

// Router initializes a router for the debug server.
func Router(opts ...Option) *chi.Mux {
	d := new(debug)

	for _, opt := range opts {
		opt(d)
	}

	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(middleware.RealIP)
	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route("/", func(root chi.Router) {
		if d.pprof {
			root.Mount(
				"/debug",
				middleware.Profiler(),
			)
		}

		root.Mount(
			"/metrics",
			metrics.Handler(
				metrics.WithToken(d.token),
			),
		)

		root.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})

		root.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)

			io.WriteString(w, http.StatusText(http.StatusOK))
		})
	})

	return mux
}

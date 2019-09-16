package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis-ocs/pkg/middleware/header"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// server gets initialized by Router and configures the router.
type server struct {
	root string
}

// Router initializes a router for the http server.
func Router(opts ...Option) *chi.Mux {
	s := new(server)

	for _, opt := range opts {
		opt(s)
	}

	mux := chi.NewRouter()

	mux.Use(hlog.NewHandler(log.Logger))
	mux.Use(hlog.RemoteAddrHandler("ip"))
	mux.Use(hlog.URLHandler("path"))
	mux.Use(hlog.MethodHandler("method"))
	mux.Use(hlog.RequestIDHandler("request_id", "Request-Id"))

	mux.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Debug().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	mux.Use(middleware.RealIP)
	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.Route(s.root, func(root chi.Router) {
		root.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)

			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	})

	return mux
}

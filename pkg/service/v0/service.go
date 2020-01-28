package svc

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/owncloud/ocis-konnectd/pkg/config"
	"github.com/owncloud/ocis-konnectd/pkg/log"
	"net/http"
	"stash.kopano.io/kc/konnect/bootstrap"
	kcconfig "stash.kopano.io/kc/konnect/config"
	"stash.kopano.io/kc/konnect/server"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Dummy(http.ResponseWriter, *http.Request)
}

// Konnectd defines implements the business logic for Service.
type Konnectd struct {
	config *config.Config
	mux    http.Handler
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	ctx := context.Background()
	options := newOptions(opts...)
	logger := options.Logger.Logger

	bs, err := bootstrap.Boot(ctx, &options.Config.Konnectd, &kcconfig.Config{
		Logger: log.Wrap(logger),
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not bootstrap konnectd")
	}

	routes := []server.WithRoutes{bs.Managers.Must("identity").(server.WithRoutes)}
	handlers := bs.Managers.Must("handler").(http.Handler)

	return Konnectd{
		config: options.Config,
		mux:    newMux(ctx, routes, handlers, options.Middleware),
	}
}

// newMux initializes the internal konnectd gorilla mux and mounts it in to a ocis chi-router
func newMux(ctx context.Context, r []server.WithRoutes, h http.Handler, middleware []func(http.Handler) http.Handler) http.Handler {
	gm := mux.NewRouter()
	for _, route := range r {
		route.AddRoutes(ctx, gm)
	}

	// Delegate rest to provider which is also a handler.
	if h != nil {
		gm.NotFoundHandler = h
	}

	m := chi.NewMux()
	m.Use(middleware...)
	m.Mount("/", gm)

	return m
}

// ServeHTTP implements the Service interface.
func (g Konnectd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Dummy implements the Service interface.
func (g Konnectd) Dummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

package svc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-ocs/pkg/config"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Dummy(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Ocs{
		config: options.Config,
		mux:    m,
	}

	m.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/", svc.Dummy)
	})

	return svc
}

// Ocs defines implements the business logic for Service.
type Ocs struct {
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (g Ocs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Dummy implements the Service interface.
func (g Ocs) Dummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

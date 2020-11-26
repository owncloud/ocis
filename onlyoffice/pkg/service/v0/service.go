package svc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis/onlyoffice/pkg/config"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Onlyoffice{
		config: options.Config,
		mux:    m,
	}

	return svc
}

// Onlyoffice defines implements the business logic for Service.
type Onlyoffice struct {
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (g Onlyoffice) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

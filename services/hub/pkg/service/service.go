package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	opkgm "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
	"net/http"
)

// Service defines the service handlers.

type Service struct {
	m *chi.Mux
}

// New returns a service implementation for Service.
func New(cfg *config.Config) Service {
	m := chi.NewMux()
	m.Use(
		opkgm.ExtractAccountUUID(
			account.JWTSecret(cfg.TokenManager.JWTSecret),
		),
	)

	m.Route("/hub", func(r chi.Router) {
		r.Route("/sse", ServeSSE)
	})

	svc := Service{
		m: m,
	}

	return svc
}

// ServeHTTP implements the Service interface.
func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

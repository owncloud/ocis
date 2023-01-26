package httpSVC

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/account"
	opkgm "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/services/authz/pkg/config"
)

// Service defines the service handlers.

type Service struct {
	m *chi.Mux
}

// New returns a service implementation for Service.
func New(cfg *config.Config) (Service, error) {

	m := chi.NewMux()
	m.Use(
		opkgm.ExtractAccountUUID(
			account.JWTSecret(cfg.TokenManager.JWTSecret),
		),
	)

	m.Get("/authz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not yet implemented"))
	})

	svc := Service{
		m: m,
	}

	return svc, nil
}

// ServeHTTP implements the Service interface.
func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTP(w, r)
}

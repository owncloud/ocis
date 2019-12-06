package svc

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis-pkg/log"
)

var (
	// ErrConfigInvalid is returned when the config parse is invalid.
	ErrConfigInvalid = `Invalid or missing config`
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Config(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	m := chi.NewMux()
	m.Use(options.Middleware...)

	svc := Phoenix{
		config: options.Config,
		mux:    m,
	}

	m.HandleFunc("/config.json", svc.Config)

	return svc
}

// Phoenix defines implements the business logic for Service.
type Phoenix struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
}

// ServeHTTP implements the Service interface.
func (g Phoenix) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Config implements the Service interface.
func (g Phoenix) Config(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(g.config.Phoenix.Path); os.IsNotExist(err) {
		g.logger.Error().
			Err(err).
			Str("config", g.config.Phoenix.Path).
			Msg("Phoenix config doesn't exist")

		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	payload, err := ioutil.ReadFile(g.config.Phoenix.Path)

	if err != nil {
		g.logger.Error().
			Err(err).
			Str("config", g.config.Phoenix.Path).
			Msg("Failed to read custom config")

		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

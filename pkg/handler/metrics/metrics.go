package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = `Invalid or missing token`
)

// metrics gets initialized by New and provides the handler.
type metrics struct {
	token string
}

// ServeHTTP just implements the http.Handler interface.
func (m metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.token == "" {
		promhttp.Handler().ServeHTTP(w, r)
		return
	}

	header := r.Header.Get("Authorization")

	if header == "" {
		log.Debug().
			Msg("Missing auth header")

		http.Error(w, ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	if header != fmt.Sprintf("Bearer %s", m.token) {
		log.Debug().
			Msg("Invalid token provided")

		http.Error(w, ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	promhttp.Handler().ServeHTTP(w, r)
}

// Handler returns the handler for metrics endpoint.
func Handler(opts ...Option) http.Handler {
	m := new(metrics)

	for _, opt := range opts {
		opt(m)
	}

	return m
}

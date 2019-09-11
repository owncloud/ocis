package config

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

var (
	// ErrConfigInvalid is returned when the config parse is invalid.
	ErrConfigInvalid = `Invalid or missing config`
)

// config gets initialized by New and provides the handler.
type config struct {
	file string
}

// ServeHTTP just implements the http.Handler interface.
func (c config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		log.Error().
			Err(err).
			Str("config", c.file).
			Msg("Phoenix config doesn't exist")

		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	payload, err := ioutil.ReadFile(c.file)

	if err != nil {
		log.Error().
			Err(err).
			Str("config", c.file).
			Msg("Failed to read custom config")

		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// Handler returns the handler for config endpoint.
func Handler(opts ...Option) http.Handler {
	c := new(config)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

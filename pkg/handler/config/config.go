package config

import (
	"encoding/json"
	"fmt"
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
	custom  string
	server  string
	theme   string
	version string
	client  string
	apps    []string
}

// auth is part of the phoenix config repsonse.
type auth struct {
	ClientID string `json:"clientId"`
	URL      string `json:"url"`
	AuthURL  string `json:"authUrl"`
}

// phoenix is part of the phoenix config response.
type phoenix struct {
	Server  string   `json:"server"`
	Theme   string   `json:"theme"`
	Version string   `json:"version"`
	Apps    []string `json:"apps"`
	Auth    auth     `json:"auth"`
}

// ServeHTTP just implements the http.Handler interface.
func (c config) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := c.payload()

	if err != nil {
		http.Error(w, ErrConfigInvalid, http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// payload prepares the phoenix config content.
func (c config) payload() ([]byte, error) {
	if c.custom != "" {
		if _, err := os.Stat(c.custom); os.IsNotExist(err) {
			log.Error().
				Err(err).
				Str("config", c.custom).
				Msg("Custom config doesn't exist")

			return []byte{}, err
		}

		payload, err := ioutil.ReadFile(c.custom)

		if err != nil {
			log.Error().
				Err(err).
				Str("config", c.custom).
				Msg("Failed to read custom config")

			return []byte{}, err
		}

		return payload, nil
	}

	payload, err := json.Marshal(phoenix{
		Server:  c.server,
		Theme:   c.theme,
		Version: c.version,
		Apps:    c.apps,
		Auth: auth{
			ClientID: c.client,
			URL:      fmt.Sprintf("%s/apps/oauth2/api/v1/token", c.server),
			AuthURL:  fmt.Sprintf("%s/apps/oauth2/authorize", c.server),
		},
	})

	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to generate config")

		return []byte{}, err
	}

	return payload, nil
}

// Handler returns the handler for config endpoint.
func Handler(opts ...Option) http.Handler {
	c := new(config)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

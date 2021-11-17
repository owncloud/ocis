package debug

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/proxy/pkg/config"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(options.Config.Service.Version),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(health(options.Config)),
		debug.Ready(ready(options.Config)),
		debug.ConfigDump(configDump(options.Config)),
	), nil
}

// health implements the health check.
func health(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// TODO(tboerger): check if services are up and running

		if _, err := io.WriteString(w, http.StatusText(http.StatusOK)); err != nil {
			panic(err)
		}
	}
}

// ready implements the ready check.
func ready(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		// TODO(tboerger): check if services are up and running

		if _, err := io.WriteString(w, http.StatusText(http.StatusOK)); err != nil {
			panic(err)
		}
	}
}

// configDump implements the config dump
func configDump(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		b, err := json.Marshal(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = w.Write(b)
	}
}

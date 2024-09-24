package debug

import (
	"encoding/json"
	"net/http"

	"github.com/ggwhite/go-masker"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(handlers.Health),
		debug.Ready(handlers.Ready),
		debug.ConfigDump(configDump(options.Config)),
	), nil
}

// configDump implements the config dump
func configDump(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		maskedCfg, err := masker.Struct(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		b, err := json.Marshal(maskedCfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = w.Write(b)
	}
}

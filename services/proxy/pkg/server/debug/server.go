package debug

import (
	"encoding/json"
	"net/http"

	"github.com/ggwhite/go-masker"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("web reachability", checks.NewHTTPCheck(options.Config.HTTP.Addr))

	readyHandlerConfiguration := healthHandlerConfiguration.
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Cluster))

	var configDumpFunc http.HandlerFunc = configDump(options.Config)
	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(handlers.NewCheckHandler(healthHandlerConfiguration)),
		debug.Ready(handlers.NewCheckHandler(readyHandlerConfiguration)),
		debug.ConfigDump(configDumpFunc),
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

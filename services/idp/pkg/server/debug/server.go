package debug

import (
	"io"
	"net/http"
	"net/url"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
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
		debug.Health(health(options.Config, options.Logger)),
		debug.Ready(handlers.Ready),
	), nil
}

// health implements the health check.
func health(cfg *config.Config, l log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetHost, err := url.Parse(cfg.Ldap.URI)
		if err != nil {
			l.Fatal().Err(err).Str("uri", cfg.Ldap.URI).Msg("invalid LDAP URI")
		}
		err = shared.RunChecklist(shared.TCPConnect(targetHost.Host))
		retVal := http.StatusOK
		if err != nil {
			l.Error().Err(err).Msg("Healtcheck failed")
			retVal = http.StatusInternalServerError
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(retVal)

		_, err = io.WriteString(w, http.StatusText(retVal))
		if err != nil {
			l.Fatal().Err(err).Msg("Could not write health check body")
		}
	}
}

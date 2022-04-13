package debug

import (
	"io"
	"net/http"
	"net/url"

	"github.com/owncloud/ocis/idp/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis-pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.String),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(health(options.Config, options.Logger)),
		debug.Ready(ready(options.Config)),
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

// ready implements the ready check.
func ready(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		// if we can call this function, a http(200) is a valid response as
		// there is nothing we can check at this point for IDP
		// if there is a mishap when initializing, there is a minimal (talking ms or ns window)
		// timeframe where this code is callable
		_, err := io.WriteString(w, http.StatusText(http.StatusOK))
		// io.WriteString should not fail but if it does we want to know.
		if err != nil {
			panic(err)
		}
	}
}

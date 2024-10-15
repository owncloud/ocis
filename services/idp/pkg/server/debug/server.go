package debug

import (
	"context"
	"net/http"
	"net/url"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("http reachability", handlers.NewHTTPCheck(options.Config.HTTP.Addr)),
	)

	readinessHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("tcp-check", func(ctx context.Context) error {
				tcpURL := options.Config.Ldap.URI
				u, err := url.Parse(options.Config.Ldap.URI)
				if err != nil {
					return err
				}
				if u.Host != "" {
					tcpURL = u.Host
				}

				return handlers.NewTCPCheck(tcpURL)(ctx)
			}).
			WithInheritedChecksFrom(healthHandler.Conf),
	)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(healthHandler),
		debug.Ready(readinessHandler),
	), nil
}

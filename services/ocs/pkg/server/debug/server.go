package debug

import (
	"context"
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("web reachability", checks.NewHTTPCheck(options.Config.HTTP.Addr))

	readyHandlerConfiguration := healthHandlerConfiguration.
		WithCheck("nats reachability", func(ctx context.Context) error {
			if len(options.Config.SigningKeys.Nodes) > 0 {
				return checks.NewNatsCheck(options.Config.SigningKeys.Nodes[0])(ctx)
			}
			return nil
		})

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
		debug.CorsAllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
		debug.CorsAllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
		debug.CorsAllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
		debug.CorsAllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
	), nil
}

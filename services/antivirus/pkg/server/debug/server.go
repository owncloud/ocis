package debug

import (
	"context"
	"net/http"

	"github.com/dutchcoders/go-clamd"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger),
	)

	readyHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("nats reachability", handlers.NewNatsCheck(options.Config.Events.Cluster)).
			WithCheck("antivirus reachability", func(ctx context.Context) error {
				cfg := options.Config
				switch cfg.Scanner.Type {
				default:
					// there is not av configured, so we panic
					panic("no antivirus configured")
				case "clamav":
					return clamd.NewClamd(cfg.Scanner.ClamAV.Socket).Ping()
				case "icap":
					return handlers.NewTCPCheck(cfg.Scanner.ICAP.URL)(ctx)
				}
			}),
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
		debug.Ready(readyHandler),
	), nil
}

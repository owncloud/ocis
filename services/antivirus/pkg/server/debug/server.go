package debug

import (
	"context"
	"errors"
	"net/http"

	"github.com/dutchcoders/go-clamd"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	readyHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Endpoint)).
		WithCheck("antivirus reachability", func(ctx context.Context) error {
			cfg := options.Config
			switch cfg.Scanner.Type {
			default:
				return errors.New("no antivirus configured")
			case "clamav":
				return clamd.NewClamd(cfg.Scanner.ClamAV.Socket).Ping()
			case "icap":
				return checks.NewTCPCheck(cfg.Scanner.ICAP.URL)(ctx)
			}
		})

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Ready(handlers.NewCheckHandler(readyHandlerConfiguration)),
	), nil
}

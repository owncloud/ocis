package debug

import (
	"context"
	"net/http"
	"net/url"

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
		WithCheck("grpc reachability", checks.NewGRPCCheck(options.Config.GRPC.Addr))

	readyHandlerConfiguration := healthHandlerConfiguration.
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Endpoint)).
		WithCheck("tika-check", func(ctx context.Context) error {
			if options.Config.Extractor.Type == "tika" {
				u, err := url.Parse(options.Config.Extractor.Tika.TikaURL)
				if err != nil {
					return err
				}
				return checks.NewTCPCheck(u.Host)(ctx)
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
	), nil
}

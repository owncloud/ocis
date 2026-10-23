package debug

import (
	"context"
	"net/http"
	"net/url"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/reva/v2/pkg/events/stream"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("grpc reachability", checks.NewGRPCCheck(options.Config.GRPC.Addr))

	readyHandlerConfiguration := healthHandlerConfiguration.
		WithCheck("nats reachability", checks.NewNatsCheckWithConfig(options.Config.Events.Endpoint, stream.NatsConfig{
			TLSInsecure:          options.Config.Events.TLSInsecure,
			TLSRootCACertificate: options.Config.Events.TLSRootCACertificate,
			EnableTLS:            options.Config.Events.EnableTLS,
			AuthUsername:         options.Config.Events.AuthUsername,
			AuthPassword:         options.Config.Events.AuthPassword,
			CredsFile:            options.Config.Events.CredsFile,
			Creds:                options.Config.Events.Creds,
		})).
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

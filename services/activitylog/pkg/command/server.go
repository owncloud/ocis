package command

import (
	"context"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/store"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/service"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"
)

var _registeredEvents = []events.Unmarshaller{
	events.PostprocessingFinished{},
}

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			tracerProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			defer cancel()

			evStream, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			evStore := store.Create(
				store.Store(cfg.Store.Store),
				store.TTL(cfg.Store.TTL),
				store.Size(cfg.Store.Size),
				microstore.Nodes(cfg.Store.Nodes...),
				microstore.Database(cfg.Store.Database),
				microstore.Table(cfg.Store.Table),
				store.Authentication(cfg.Store.AuthUsername, cfg.Store.AuthPassword),
			)

			{
				svc, err := service.New(
					service.Logger(logger),
					service.Config(cfg),
					service.TraceProvider(tracerProvider),
					service.Stream(evStream),
					service.RegisteredEvents(_registeredEvents),
					service.Store(evStore),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(func() error {
					return svc.Run()
				}, func(err error) {
					logger.Error().
						Str("transport", "http").
						Err(err).
						Msg("Shutting down server")

					cancel()
					os.Exit(1)
				})
			}

			{
				server := debug.NewService(
					debug.Logger(logger),
					debug.Name(cfg.Service.Name),
					debug.Version(version.GetString()),
					debug.Address(cfg.Debug.Addr),
					debug.Token(cfg.Debug.Token),
					debug.Pprof(cfg.Debug.Pprof),
					debug.Zpages(cfg.Debug.Zpages),
					debug.Health(handlers.Health),
					debug.Ready(handlers.Ready),
				)

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

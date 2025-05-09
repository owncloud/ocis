package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/owncloud/reva/v2/pkg/events/stream"
	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"

	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/logging"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/service"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			var cancel context.CancelFunc
			ctx := cfg.Context
			if ctx == nil {
				ctx, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}

			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			gr := runner.NewGroup()
			{
				connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
				bus, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Postprocessing.Events))
				if err != nil {
					return err
				}

				st := store.Create(
					store.Store(cfg.Store.Store),
					store.TTL(cfg.Store.TTL),
					microstore.Nodes(cfg.Store.Nodes...),
					microstore.Database(cfg.Store.Database),
					microstore.Table(cfg.Store.Table),
					store.Authentication(cfg.Store.AuthUsername, cfg.Store.AuthPassword),
				)

				svc, err := service.NewPostprocessingService(ctx, bus, logger, st, traceProvider, cfg.Postprocessing)
				if err != nil {
					return err
				}

				gr.Add(runner.New("postprocessing_svc", func() error {
					return svc.Run()
				}, func() {
					svc.Close()
				}))
			}

			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)
				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.NewGolangHttpServerRunner("postprocessing_debug", debugServer))
			}

			grResults := gr.Run(ctx)

			// return the first non-nil error found in the results
			for _, grResult := range grResults {
				if grResult.RunnerError != nil {
					return grResult.RunnerError
				}
			}
			return nil
		},
	}
}

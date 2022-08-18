package command

import (
	"context"
	"fmt"
	"os"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/search/pkg/logging"
	"github.com/owncloud/ocis/v2/services/search/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/search/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/search/pkg/server/grpc"
	"github.com/owncloud/ocis/v2/services/search/pkg/tracing"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
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
			err := tracing.Configure(cfg)
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
			defer cancel()

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			grpcServer, err := grpc.Server(
				grpc.Config(cfg),
				grpc.Logger(logger),
				grpc.Name(cfg.Service.Name),
				grpc.Context(ctx),
				grpc.Metrics(mtrcs),
			)

			if err != nil {
				logger.Info().Err(err).Str("transport", "grpc").Msg("Failed to initialize server")
				return err
			}

			gr.Add(grpcServer.Run, func(_ error) {
				logger.Error().
					Err(err).
					Str("server", "grpc").
					Msg("shutting down server")
				cancel()
			})

			server, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(server.ListenAndServe, func(_ error) {
				_ = server.Shutdown(ctx)
				cancel()
			})

			return gr.Run()
		},
	}
}

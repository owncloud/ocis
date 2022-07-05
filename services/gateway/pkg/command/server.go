package command

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/cs3org/reva/v2/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/external"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/config"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/logging"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/revaconfig"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/tracing"
	"github.com/urfave/cli/v2"
)

// Server is the entry point for the server command.
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
			err := tracing.Configure(cfg, logger)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)

			defer cancel()

			pidFile := path.Join(os.TempDir(), "revad-"+cfg.Service.Name+"-"+uuid.Must(uuid.NewV4()).String()+".pid")

			rcfg := revaconfig.GatewayConfigFromStruct(cfg, logger)

			gr.Add(func() error {
				runtime.RunWithOptions(rcfg, pidFile, runtime.WithLogger(&logger.Logger))
				return nil
			}, func(err error) {
				logger.Error().
					Err(err).
					Str("server", cfg.Service.Name).
					Msg("Shutting down server")

				cancel()
			})

			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)

			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				cancel()
			})

			if err := external.RegisterGRPCEndpoint(
				ctx,
				cfg.GRPC.Namespace+"."+cfg.Service.Name,
				uuid.Must(uuid.NewV4()).String(),
				cfg.GRPC.Addr,
				version.GetString(),
				logger,
			); err != nil {
				logger.Fatal().Err(err).Msg("failed to register the grpc endpoint")
			}

			return gr.Run()
		},
	}
}

// defineContext sets the context for the extension. If there is a context configured it will create a new child from it,
// if not, it will create a root context that can be cancelled.
func defineContext(cfg *config.Config) (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		if cfg.Context == nil {
			return context.WithCancel(context.Background())
		}
		return context.WithCancel(cfg.Context)
	}()
}

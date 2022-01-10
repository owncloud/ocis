package command

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/gateway/pkg/config"
	"github.com/owncloud/ocis/gateway/pkg/config/parser"
	"github.com/owncloud/ocis/gateway/pkg/config/reva"
	"github.com/owncloud/ocis/gateway/pkg/logging"
	"github.com/owncloud/ocis/gateway/pkg/metrics"
	"github.com/owncloud/ocis/gateway/pkg/server/debug"
	"github.com/owncloud/ocis/gateway/pkg/tracing"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/storage/pkg/service/external"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start %s extension without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return parser.ParseConfig(cfg)
		},
		Action: func(c *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			err := tracing.Configure(cfg)
			if err != nil {
				return err
			}
			gr := run.Group{}
			ctx, cancel := defineContext(cfg)
			mtrcs := metrics.New()

			defer cancel()

			mtrcs.BuildInfo.WithLabelValues(version.String).Set(1)

			uuid := uuid.Must(uuid.NewV4())

			pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

			rcfg, err := reva.Config(cfg)
			if err != nil {
				logger.Error().Err(err).Str("server", "reva configuration").Msg("Failed to initialize server")
				return err
			}
			//logger.Debug().
			//	Str("server", "gateway").
			//	Interface("reva-config", rcfg).
			//	Msg("config")

			defer cancel()

			gr.Add(func() error {
				err := external.RegisterGRPCEndpoint(
					ctx,
					cfg.GRPC.Namespace,
					uuid.String(),
					cfg.GRPC.Addr,
					version.String,
					logger,
				)

				if err != nil {
					return err
				}

				runtime.RunWithOptions(
					rcfg,
					pidFile,
					runtime.WithLogger(&logger.Logger),
				)
				return nil
			}, func(_ error) {
				logger.Info().
					Str("server", c.Command.Name).
					Msg("Shutting down server")

				cancel()
			})

			// TODO: what is this?
			//if !cfg.Reva.Gateway.Supervised {
			//	sync.Trap(&gr, cancel)
			//}

			// prepare a debug server and add it to the group run.
			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Error().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(debugServer.ListenAndServe, func(_ error) {
				_ = debugServer.Shutdown(ctx)
				cancel()
			})

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

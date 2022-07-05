package command

import (
	"context"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/micro/ocdav"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/ocdav/pkg/config"
	"github.com/owncloud/ocis/v2/services/ocdav/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/ocdav/pkg/logging"
	"github.com/owncloud/ocis/v2/services/ocdav/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/ocdav/pkg/tracing"
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

			gr.Add(func() error {

				opts := []ocdav.Option{
					ocdav.Name(cfg.HTTP.Namespace + "." + cfg.Service.Name),
					ocdav.Version(version.GetString()),
					ocdav.Context(ctx),
					ocdav.Logger(logger.Logger),
					ocdav.Address(cfg.HTTP.Addr),
					ocdav.FilesNamespace(cfg.FilesNamespace),
					ocdav.WebdavNamespace(cfg.WebdavNamespace),
					ocdav.SharesNamespace(cfg.SharesNamespace),
					ocdav.Timeout(cfg.Timeout),
					ocdav.Insecure(cfg.Insecure),
					ocdav.PublicURL(cfg.PublicURL),
					ocdav.Prefix(cfg.HTTP.Prefix),
					ocdav.GatewaySvc(cfg.Reva.Address),
					ocdav.JWTSecret(cfg.TokenManager.JWTSecret),
					ocdav.ProductName(cfg.Status.ProductName),
					ocdav.ProductVersion(cfg.Status.ProductVersion),
					ocdav.Product(cfg.Status.Product),
					ocdav.Version(cfg.Status.Version),
					ocdav.VersionString(cfg.Status.VersionString),
					ocdav.Edition(cfg.Status.Edition),
					ocdav.MachineAuthAPIKey(cfg.MachineAuthAPIKey),
					// ocdav.FavoriteManager() // FIXME needs a proper persistence implementation https://github.com/owncloud/ocis/issues/1228
					// ocdav.LockSystem(), // will default to the CS3 lock system
					// ocdav.TLSConfig() // tls config for the http server
				}

				if cfg.Tracing.Enabled {
					opts = append(opts, ocdav.Tracing(cfg.Tracing.Endpoint, cfg.Tracing.Collector))
				}

				s, err := ocdav.Service(opts...)
				if err != nil {
					return err
				}

				return s.Run()
			}, func(err error) {
				logger.Error().
					Err(err).
					Str("server", c.Command.Name).
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

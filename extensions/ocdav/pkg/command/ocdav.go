package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/micro/ocdav"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/ocdav/pkg/config"
	"github.com/owncloud/ocis/extensions/ocdav/pkg/config/parser"
	"github.com/owncloud/ocis/extensions/ocdav/pkg/server/debug"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/owncloud/ocis/ocis-pkg/tracing"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// OCDav is the entrypoint for the ocdav command.
// TODO move ocdav cmd to a separate service
func OCDav(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "ocdav",
		Usage: "start ocdav service",
		Before: func(ctx *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			logCfg := cfg.Logging
			logger := log.NewLogger(
				log.Level(logCfg.Level),
				log.File(logCfg.File),
				log.Pretty(logCfg.Pretty),
				log.Color(logCfg.Color),
			)
			tracing.Configure(cfg.Tracing.Enabled, cfg.Tracing.Type, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			//metrics     = metrics.New()

			defer cancel()

			gr.Add(func() error {
				s, err := ocdav.Service(
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
					// ocdav.FavoriteManager() // FIXME needs a proper persistence implementation
					// ocdav.LockSystem(), // will default to the CS3 lock system
					// ocdav.TLSConfig() // tls config for the http server
				)
				if err != nil {
					return err
				}

				return s.Run()
			}, func(err error) {
				logger.Info().Err(err).Str("server", c.Command.Name).Msg("Shutting down server")
				cancel()
			})

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", "debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					cancel()
				})
			}

			if !cfg.Supervised {
				sync.Trap(&gr, cancel)
			}

			return gr.Run()
		},
	}
}

// OCDavSutureService allows for the ocdav command to be embedded and supervised by a suture supervisor tree.
type OCDavSutureService struct {
	cfg *config.Config
}

// NewOCDav creates a new ocdav.OCDavSutureService
func NewOCDav(cfg *ociscfg.Config) suture.Service {
	cfg.OCDav.Commons = cfg.Commons
	return OCDavSutureService{
		cfg: cfg.OCDav,
	}
}

func (s OCDavSutureService) Serve(ctx context.Context) error {
	// s.cfg.Reva.Frontend.Context = ctx
	cmd := OCDav(s.cfg)
	f := &flag.FlagSet{}
	cmdFlags := cmd.Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if cmd.Before != nil {
		if err := cmd.Before(cliCtx); err != nil {
			return err
		}
	}
	if err := cmd.Action(cliCtx); err != nil {
		return err
	}

	return nil
}

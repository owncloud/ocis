package command

import (
	"context"
	"flag"

	"github.com/cs3org/reva/v2/pkg/micro/ocdav"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/extensions/storage/pkg/config"
	"github.com/owncloud/ocis/extensions/storage/pkg/server/debug"
	"github.com/owncloud/ocis/extensions/storage/pkg/tracing"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/sync"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// OCDav is the entrypoint for the ocdav command.
// TODO move ocdav cmd to a separate service
func OCDav(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "ocdav",
		Usage: "start ocdav service",
		Before: func(c *cli.Context) error {
			if err := loadUserAgent(c, cfg); err != nil {
				return err
			}
			return ParseConfig(c, cfg, "ocdav")
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			tracing.Configure(cfg, logger)

			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			//metrics     = metrics.New()

			defer cancel()

			gr.Add(func() error {
				s, err := ocdav.Service(
					ocdav.Context(ctx),
					ocdav.Logger(logger.Logger),
					ocdav.Address(cfg.OCDav.Addr),
					ocdav.FilesNamespace(cfg.OCDav.FilesNamespace),
					ocdav.WebdavNamespace(cfg.OCDav.WebdavNamespace),
					ocdav.SharesNamespace(cfg.OCDav.SharesNamespace),
					ocdav.Timeout(cfg.OCDav.Timeout),
					ocdav.Insecure(cfg.OCDav.Insecure),
					ocdav.PublicURL(cfg.OCDav.PublicURL),
					ocdav.Prefix(cfg.OCDav.Prefix),
					ocdav.GatewaySvc(cfg.OCDav.GatewaySVC),
					ocdav.JWTSecret(cfg.OCDav.JWTSecret),
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
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.OCDav.DebugAddr),
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

			if !cfg.Reva.Frontend.Supervised {
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
	cfg.Storage.Commons = cfg.Commons
	return OCDavSutureService{
		cfg: cfg.Storage,
	}
}

func (s OCDavSutureService) Serve(ctx context.Context) error {
	s.cfg.Reva.Frontend.Context = ctx
	f := &flag.FlagSet{}
	cmdFlags := OCDav(s.cfg).Flags
	for k := range cmdFlags {
		if err := cmdFlags[k].Apply(f); err != nil {
			return err
		}
	}
	cliCtx := cli.NewContext(nil, f, nil)
	if OCDav(s.cfg).Before != nil {
		if err := OCDav(s.cfg).Before(cliCtx); err != nil {
			return err
		}
	}
	if err := OCDav(s.cfg).Action(cliCtx); err != nil {
		return err
	}

	return nil
}

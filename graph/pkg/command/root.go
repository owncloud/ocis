package command

import (
	"context"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/thejerf/suture/v4"

	"github.com/owncloud/ocis/graph/pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis-graph command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-graph",
		Version:  version.String,
		Usage:    "Serve Graph API for oCIS",
		Compiled: version.Compiled(),
		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Before: func(c *cli.Context) error {
			cfg.Server.Version = version.String
			return ParseConfig(c, cfg)
		},
		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
		},
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help,h",
		Usage: "Show the help",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version,v",
		Usage: "Print the version",
	}

	return app.Run(os.Args)
}

// NewLogger initializes a service-specific logger instance.
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("graph"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads graph configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	conf, err := ociscfg.BindSourcesToStructs("graph", cfg)
	if err != nil {
		return err
	}
	conf.LoadOSEnv(config.GetEnv(), false)
	bindings := config.StructMappings(cfg)
	return ociscfg.BindEnv(conf, bindings)
}

// SutureService allows for the graph command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new graph.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if (cfg.Accounts.Log == shared.Log{}) {
		cfg.Accounts.Log = cfg.Log
	}
	return SutureService{
		cfg: cfg.Graph,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}

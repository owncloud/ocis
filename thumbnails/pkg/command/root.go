package command

import (
	"context"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/thejerf/suture/v4"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis-thumbnails command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "ocis-thumbnails",
		Version:  version.String,
		Usage:    "Example usage",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},

		Before: func(c *cli.Context) error {
			cfg.Server.Version = version.String
			return nil
		},

		Commands: []*cli.Command{
			Server(cfg),
			Health(cfg),
			PrintVersion(cfg),
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
		log.Name("thumbnails"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads configuration from Viper known paths.
// ParseConfig loads glauth configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	conf, err := ociscfg.BindSourcesToStructs("thumbnails", cfg)
	if err != nil {
		return err
	}
	conf.LoadOSEnv(config.GetEnv(), false)
	bindings := config.StructMappings(cfg)
	return ociscfg.BindEnv(conf, bindings)
}

// SutureService allows for the thumbnails command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	cfg *config.Config
}

// NewSutureService creates a new thumbnails.SutureService
func NewSutureService(cfg *ociscfg.Config) suture.Service {
	if (cfg.Thumbnails.Log == shared.Log{}) {
		cfg.Thumbnails.Log = cfg.Log
	}
	return SutureService{
		cfg: cfg.Thumbnails,
	}
}

func (s SutureService) Serve(ctx context.Context) error {
	s.cfg.Context = ctx
	if err := Execute(s.cfg); err != nil {
		return err
	}

	return nil
}

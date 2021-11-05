package command

import (
	"os"

	"github.com/owncloud/ocis/ocis/pkg/flagset"

	"github.com/owncloud/ocis/ocis-pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis",
		Version:  version.String,
		Usage:    "ownCloud Infinite Scale Stack",
		Compiled: version.Compiled(),
		Flags:    flagset.RootWithConfig(cfg),
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},
		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
	}

	for _, fn := range register.Commands {
		app.Commands = append(
			app.Commands,
			fn(cfg),
		)
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

// NewLogger initializes a service-specific logger instance
func NewLogger(cfg *config.Config) log.Logger {
	return log.NewLogger(
		log.Name("ocis"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

// ParseConfig loads ocis configuration from known paths.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	conf, err := ociscfg.BindSourcesToStructs("ocis", cfg)
	if err != nil {
		return err
	}

	conf.LoadOSEnv(config.GetEnv(), false)

	if err = cfg.UnmapEnv(conf); err != nil {
		return err
	}

	return nil
}

package command

import (
	"errors"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis command.
func Execute() error {
	cfg := config.DefaultConfig()

	app := &cli.App{
		Name:     "ocis",
		Version:  version.String,
		Usage:    "ownCloud Infinite Scale Stack",
		Compiled: version.Compiled(),

		Before: func(c *cli.Context) error {
			//cfg.Service.Version = version.String
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

// ParseConfig loads ocis configuration.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	_, err := ociscfg.BindSourcesToStructs("ocis", cfg)
	if err != nil {
		return err
	}

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	return nil
}

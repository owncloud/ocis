package command

import (
	"errors"
	"os"

	"github.com/owncloud/ocis/ocis-pkg/clihelper"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the ocis command.
func Execute() error {
	cfg := config.DefaultConfig()

	app := clihelper.DefaultApp(&cli.App{
		Name:  "ocis",
		Usage: "ownCloud Infinite Scale Stack",

		Before: func(c *cli.Context) error {
			// TODO: what do do?
			//cfg.Service.Version = version.String
			return ParseConfig(c, cfg)
		},
	})

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

	return app.Run(os.Args)
}

// ParseConfig loads ocis configuration.
func ParseConfig(c *cli.Context, cfg *config.Config) error {
	_, err := config.BindSourcesToStructs("ocis", cfg)
	if err != nil {
		return err
	}

	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &shared.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil && cfg.Commons == nil {
		cfg.Log = &shared.Log{}
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

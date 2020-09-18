package command

import (
	"os"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/runtime"
	"github.com/owncloud/ocis/ocis/pkg/version"
)

// Execute is the entry point for the ocis-ocis command.
func Execute() error {
	cfg := config.New()

	app := &cli.App{
		Name:     "ocis",
		Version:  version.String,
		Usage:    "ownCloud Infinite Scale Stack",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Flags: flagset.RootWithConfig(cfg),
	}

	for _, fn := range register.Commands {
		app.Commands = append(
			app.Commands,
			fn(cfg),
		)
	}

	runtime.AddMicroPlatform(app)

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
	)
}

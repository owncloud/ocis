package command

import (
	"os"

	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/version"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// Execute is the entry point for the storage command.
func Execute(cfg *config.Config) error {
	app := &cli.App{
		Name:     "storage",
		Version:  version.String,
		Usage:    "Storage service for oCIS",
		Compiled: version.Compiled(),

		Authors: []*cli.Author{
			{
				Name:  "ownCloud GmbH",
				Email: "support@owncloud.com",
			},
		},
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg, "storage")
		},

		Commands: []*cli.Command{
			Frontend(cfg),
			Gateway(cfg),
			Users(cfg),
			Groups(cfg),
			AppProvider(cfg),
			AuthBasic(cfg),
			AuthBearer(cfg),
			Sharing(cfg),
			StorageHome(cfg),
			StorageUsers(cfg),
			StoragePublicLink(cfg),
			StorageMetadata(cfg),
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
		log.Name("storage"),
		log.Level(cfg.Log.Level),
		log.Pretty(cfg.Log.Pretty),
		log.Color(cfg.Log.Color),
		log.File(cfg.Log.File),
	)
}

package command

import (
	"os"

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
			cfg.Service.Version = version.String
			return ParseConfig(c, cfg, "_")
		},

		Commands: []*cli.Command{
			Frontend(cfg),
			Gateway(cfg),
			Users(cfg),
			Groups(cfg),
			AppProvider(cfg),
			AuthBasic(cfg),
			AuthBearer(cfg),
			AuthMachine(cfg),
			Sharing(cfg),
			StorageUsers(cfg),
			StorageShares(cfg),
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

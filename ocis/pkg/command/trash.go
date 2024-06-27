package command

import (
	"fmt"
	"github.com/owncloud/ocis/v2/ocis/pkg/trash"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

func TrashCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "trash",
		Usage: "ocis trash functionality",
		Subcommands: []*cli.Command{
			TrashPurgeOrphanedDirsCommand(cfg),
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg, true))
		},
		Action: func(_ *cli.Context) error {
			fmt.Println("Read the docs")
			return nil
		},
	}
}

func TrashPurgeOrphanedDirsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-orphaned-dirs",
		Usage: "purge orphaned directories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "basepath",
				Aliases:  []string{"p"},
				Usage:    "the basepath of the decomposedfs (e.g. /var/tmp/ocis/storage/users)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "do not delete anything, just print what would be deleted",
				Value: true,
			},
		},
		Action: func(c *cli.Context) error {
			basePath := c.String("basepath")
			if basePath == "" {
				fmt.Println("basepath is required")
				return cli.ShowCommandHelp(c, "trash")
			}

			if err := trash.PurgeTrashOrphanedPaths(basePath, c.Bool("dry-run")); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}

func init() {
	register.AddCommand(TrashCommand)
}

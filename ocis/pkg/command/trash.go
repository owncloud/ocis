package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/owncloud/ocis/v2/ocis/pkg/trash"
	ocisbs "github.com/owncloud/reva/v2/pkg/storage/fs/ocis/blobstore"
	s3bs "github.com/owncloud/reva/v2/pkg/storage/fs/s3ng/blobstore"

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
			TrashPurgeEmptyDirsCommand(cfg),
			TrashPurgeOrphanedFilesCommand(cfg),
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

func TrashPurgeEmptyDirsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-empty-dirs",
		Usage: "purge empty directories",
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

			if err := trash.PurgeTrashEmptyPaths(basePath, c.Bool("dry-run")); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}

// TrashPurgeOrphanedFilesCommand removes the files of nodes that do not exist
func TrashPurgeOrphanedFilesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-orphaned-files",
		Usage: "purge lock files and revisions of nodes that do not exist",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "basepath",
				Aliases:  []string{"p"},
				Usage:    "the basepath of the decomposedfs (e.g. /var/tmp/ocis/storage/users)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "blobstore",
				Aliases: []string{"b"},
				Usage:   "the blobstore type. Can be (none, ocis, s3ng). With none, orphaned revisions stay on the disk",
				Value:   "ocis",
			},
			&cli.DurationFlag{
				Name:  "min-age",
				Usage: "do not remove files that are younger than this duration",
				Value: time.Hour,
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "do not delete anything, just print what would be deleted",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "print each file that is removed or kept. Without it, only errors and the summary are printed",
			},
		},
		Action: func(c *cli.Context) error {
			basePath := c.String("basepath")
			if basePath == "" {
				fmt.Println("basepath is required")
				return cli.ShowCommandHelp(c, "purge-orphaned-files")
			}

			// without a blobstore, orphaned revisions stay on the disk
			var bs trash.BlobDeleter
			if c.String("blobstore") != "none" {
				var err error
				if bs, err = newBlobDeleter(cfg, c.String("blobstore"), basePath); err != nil {
					fmt.Println(err)
					return err
				}
			}

			if err := trash.PurgeOrphanedFiles(basePath, trash.PurgeOrphansOptions{
				Blobstore: bs,
				MinAge:    c.Duration("min-age"),
				DryRun:    c.Bool("dry-run"),
				Verbose:   c.Bool("verbose"),
				Out:       c.App.Writer,
			}); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}

// newBlobDeleter returns the blobstore of the given type
func newBlobDeleter(cfg *config.Config, blobstore, basePath string) (trash.BlobDeleter, error) {
	switch blobstore {
	case "s3ng":
		return s3bs.New(
			cfg.StorageUsers.Drivers.S3NG.Endpoint,
			cfg.StorageUsers.Drivers.S3NG.Region,
			cfg.StorageUsers.Drivers.S3NG.Bucket,
			cfg.StorageUsers.Drivers.S3NG.AccessKey,
			cfg.StorageUsers.Drivers.S3NG.SecretKey,
			s3bs.Options{},
		)
	case "ocis":
		return ocisbs.New(basePath)
	default:
		return nil, errors.New("blobstore type not supported")
	}
}

func init() {
	register.AddCommand(TrashCommand)
}

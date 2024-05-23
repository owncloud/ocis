package command

import (
	"errors"
	"fmt"

	ocisbs "github.com/cs3org/reva/v2/pkg/storage/fs/ocis/blobstore"
	s3bs "github.com/cs3org/reva/v2/pkg/storage/fs/s3ng/blobstore"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/backup"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// BackupCommand is the entrypoint for the backup command
func BackupCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "backup",
		Usage: "ocis backup functionality",
		Subcommands: []*cli.Command{
			ConsistencyCommand(cfg),
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg, false))
		},
		Action: func(_ *cli.Context) error {
			fmt.Println("Read the docs")
			return nil
		},
	}
}

// ConsistencyCommand is the entrypoint for the consistency Command
func ConsistencyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "consistency",
		Usage: "check backup consistency",
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
				Usage:   "the blobstore type. Can be (ocis, s3ng). Default ocis",
				Value:   "ocis",
			},
		},
		Action: func(c *cli.Context) error {
			basePath := c.String("basepath")
			if basePath == "" {
				fmt.Println("basepath is required")
				return cli.ShowCommandHelp(c, "consistency")
			}

			var (
				bs  backup.ListBlobstore
				err error
			)
			switch c.String("blobstore") {
			case "s3ng":
				bs, err = s3bs.New(
					cfg.StorageUsers.Drivers.S3NG.Endpoint,
					cfg.StorageUsers.Drivers.S3NG.Region,
					cfg.StorageUsers.Drivers.S3NG.Bucket,
					cfg.StorageUsers.Drivers.S3NG.AccessKey,
					cfg.StorageUsers.Drivers.S3NG.SecretKey,
					s3bs.Options{},
				)
			case "ocis":
				bs, err = ocisbs.New(basePath)
			default:
				err = errors.New("blobstore type not supported")
			}
			if err != nil {
				fmt.Println(err)
				return err
			}
			if err := backup.CheckSpaceConsistency(basePath, bs); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}

func init() {
	register.AddCommand(BackupCommand)
}

package command

import (
	"fmt"
	"path/filepath"

	"github.com/cs3org/reva/v2/pkg/storage/fs/ocis/blobstore"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
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
		Action: func(c *cli.Context) error {
			fmt.Println("Read the docs")
			return nil
		},
	}
}

// ConsistencyCommand is the entrypoint for the consistency Command
func ConsistencyCommand(_ *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "consistency",
		Usage: "check backup consistency",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "blobstore",
				Aliases: []string{"b"},
				Usage:   "the blobstore type. Can be (ocis, s3ng). Default ocis",
			},
		},
		Action: func(c *cli.Context) error {
			basePath := "/home/jkoberg/.ocis/storage/users"

			// TODO: switch for s3ng blobstore
			bs, err := blobstore.New(basePath)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if err := backup.CheckSpaceConsistency(filepath.Join(basePath, "spaces/23/ebf113-76d4-43c0-8594-df974b02cd74"), bs); err != nil {
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

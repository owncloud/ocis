package command

import (
	"github.com/owncloud/ocis/extensions/storage-metadata/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageMetadataCommand is the entrypoint for the storage-metadata command.
func StorageMetadataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-metadata",
		Usage:    "start storage and data service for metadata",
		Category: "extensions",
		Action: func(c *cli.Context) error {
			origCmd := command.StorageMetadata(cfg.StorageMetadata)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageMetadataCommand)
}

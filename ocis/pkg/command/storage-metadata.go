package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/storage-metadata/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageMetadataCommand is the entrypoint for the StorageMetadata command.
func StorageMetadataCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageMetadata.Service.Name,
		Usage:    subcommandDescription(cfg.StorageMetadata.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.StorageMetadata.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StorageMetadata),
	}
}

func init() {
	register.AddCommand(StorageMetadataCommand)
}

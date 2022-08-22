package command

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/command/helper"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/owncloud/ocis/v2/services/storage-shares/pkg/command"
	"github.com/urfave/cli/v2"
)

// StorageSharesCommand is the entrypoint for the StorageShares command.
func StorageSharesCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.StorageShares.Service.Name,
		Usage:    helper.SubcommandDescription(cfg.StorageShares.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			configlog.Error(parser.ParseConfig(cfg, true))
			cfg.StorageShares.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.StorageShares),
	}
}

func init() {
	register.AddCommand(StorageSharesCommand)
}

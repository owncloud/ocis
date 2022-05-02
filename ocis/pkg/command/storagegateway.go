package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/gateway/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// StorageGatewayCommand is the entrypoint for the reva-gateway command.
func StorageGatewayCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "storage-gateway",
		Usage:    "start storage gateway",
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Gateway.Commons = cfg.Commons
			return nil
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Gateway(cfg.Gateway)
			return handleOriginalAction(c, origCmd)
		},
	}
}

func init() {
	register.AddCommand(StorageGatewayCommand)
}

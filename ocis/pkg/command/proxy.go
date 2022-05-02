package command

import (
	"fmt"

	"github.com/owncloud/ocis/extensions/proxy/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// ProxyCommand is the entry point for the proxy command.
func ProxyCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     cfg.Proxy.Service.Name,
		Usage:    subcommandDescription(cfg.Proxy.Service.Name),
		Category: "extensions",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			cfg.Proxy.Commons = cfg.Commons
			return nil
		},
		Subcommands: command.GetCommands(cfg.Proxy),
	}
}

func init() {
	register.AddCommand(ProxyCommand)
}

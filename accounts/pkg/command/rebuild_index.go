package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
)

func DeleteIndex(cdf *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "add",
		Usage:   "Create a new account",
		Aliases: []string{"create", "a"},
		//Flags:   flagset.AddAccountWithConfig(cfg, a),
	}
}

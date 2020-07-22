// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-store/pkg/command"
	svcconfig "github.com/owncloud/ocis-store/pkg/config"
	"github.com/owncloud/ocis-store/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "store",
		Usage:    "Start a go-micro store",
		Category: "Runtime",
		Flags:    flagset.ServerWithConfig(cfg.Store),
		Action: func(ctx *cli.Context) error {
			storeCommand := command.Server(configureStore(cfg))

			if err := storeCommand.Before(ctx); err != nil {
				return err
			}

			return cli.HandleAction(storeCommand.Action, ctx)
		},
	}
}

func configureStore(cfg *config.Config) *svcconfig.Config {
	cfg.Store.Log.Level = cfg.Log.Level
	cfg.Store.Log.Pretty = cfg.Log.Pretty
	cfg.Store.Log.Color = cfg.Log.Color

	return cfg.Store
}

func init() {
	register.AddCommand(StoreCommand)
}

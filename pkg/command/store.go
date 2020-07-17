// +build !simple

package command

import (
	"context"
	"strings"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	gostore "github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/store/file"
	"github.com/owncloud/ocis-ocs/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// StoreCommand is the entrypoint for the ocs command.
func StoreCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "store",
		Usage:    "Start a go-micro store",
		Category: "Runtime",
		Flags:    flagset.ServerWithConfig(cfg.OCS),
		Action: func(ctx *cli.Context) error {
			file.DefaultDir = "/var/tmp/ocis/store"
			store := file.NewStore(
				gostore.Database("ocis"),
				gostore.Table("ocis"),
				gostore.WithContext(context.Background()),
			)

			mopts := []micro.Option{
				micro.Name(
					strings.Join(
						[]string{
							"com.owncloud",
							"store",
						},
						".",
					),
				),
				micro.RegisterTTL(time.Second * 30),
				micro.RegisterInterval(time.Second * 10),
				micro.Store(store),
			}

			service := micro.NewService(mopts...)
			service.Init()
			return service.Run()
		},
	}
}

func init() {
	register.AddCommand(StoreCommand)
}

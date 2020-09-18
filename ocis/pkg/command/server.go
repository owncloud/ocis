// +build !simple

package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis/pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/runtime"
	"github.com/owncloud/ocis/ocis/pkg/tracing"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "Start fullstack server",
		Category: "Fullstack",
		Flags:    flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			if err := tracing.Start(cfg); err != nil {
				return err
			}

			r := runtime.New()
			// TODO temporary service startup selection. Should go away and the runtime should take care of it.
			return r.Start(append([]string{
				"accounts",
				"settings",
				"konnectd",
				"proxy",
				"ocs",
				"phoenix",
				"glauth",
				"webdav",
				"store",
				"thumbnails",
				"reva-frontend",
				"reva-gateway",
				"reva-users",
				"reva-auth-basic",
				"reva-auth-bearer",
				"reva-storage-home",
				"reva-storage-home-data",
				"reva-storage-eos",
				"reva-storage-eos-data",
				"reva-storage-oc",
				"reva-storage-oc-data",
				"reva-storage-public-link",
			}, runtime.MicroServices...)...)
		},
	}
}

func init() {
	register.AddCommand(Server)
}

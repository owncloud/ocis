//go:build !simple
// +build !simple

package command

import (
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/flagset"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/runtime"
	"github.com/urfave/cli/v2"
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
			r := runtime.New(cfg)
			return r.Start()
		},
	}
}

func init() {
	register.AddCommand(Server)
}

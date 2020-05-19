// +build !simple

package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"github.com/owncloud/ocis/pkg/micro/runtime"
	"github.com/owncloud/ocis/pkg/register"
	"github.com/owncloud/ocis/pkg/tracing"
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
			logger := NewLogger(cfg)

			if err := tracing.Start(cfg); err != nil {
				return err
			}

			runtime := runtime.New(
				runtime.Services(append(runtime.MicroServices, runtime.Extensions...)),
				runtime.Logger(logger),
				runtime.MicroRuntime(cmd.DefaultCmd.Options().Runtime),
				runtime.Context(c),
			)

			runtime.Start()

			return nil
		},
	}
}

func init() {
	register.AddCommand(Server)
}

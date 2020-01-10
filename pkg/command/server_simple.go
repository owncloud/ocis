// +build simple

package command

import (
	"strings"

	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"github.com/owncloud/ocis/pkg/micro/runtime"
	"github.com/owncloud/ocis/pkg/register"
	"github.com/owncloud/ocis/pkg/tracing"
)

// Simple is the entrypoint for the server command.
func Simple(cfg *config.Config) cli.Command {
	return cli.Command{
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
				runtime.Services(
					append(
						runtime.RuntimeServices,
						[]string{
							"hello",
							"konnectd",
							"phoenix",
						}...,
					),
				),
				runtime.Logger(logger),
				runtime.MicroRuntime(cmd.DefaultCmd.Options().Runtime),
			)

			// fork uses the micro runtime to fork go-micro services
			runtime.Start()

			// trap blocks until a kill signal is sent
			runtime.Trap()

			return nil
		},
	}
}

func init() {
	register.AddCommand(Simple)
}

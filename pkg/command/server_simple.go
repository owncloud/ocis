// +build simple

package command

import (
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/config/cmd"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/flagset"
	"github.com/owncloud/ocis/pkg/micro/runtime"
	"github.com/owncloud/ocis/pkg/register"
	"github.com/owncloud/ocis/pkg/tracing"
)

var (
	// SimpleRuntimeServices declares which services will be started for the fullstack server
	SimpleRuntimeServices = []string{
		"hello",
		"phoenix",
		"reva-frontend",
		"reva-gateway",
		"reva-users",
		"reva-auth-basic",
		"reva-auth-bearer",
	}
)

// Simple is the entrypoint for the server command. It is the `ocis server` subcommand overloaded with a different set of services
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
				runtime.Logger(logger),
				runtime.Services(append(runtime.RuntimeServices, SimpleRuntimeServices...)),
				runtime.MicroRuntime(cmd.DefaultCmd.Options().Runtime),
			)

			{
				runtime.Start()
				runtime.Trap()
			}

			return nil
		},
	}
}

func init() {
	register.AddCommand(Simple)
}

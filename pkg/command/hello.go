package command

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-hello/pkg/command"
	svcconfig "github.com/owncloud/ocis-hello/pkg/config"
	"github.com/owncloud/ocis-hello/pkg/flagset"
	"github.com/owncloud/ocis/pkg/config"
	"github.com/owncloud/ocis/pkg/register"
)

// HelloCommand is the entrypoint for the hello command.
func HelloCommand(cfg *config.Config) cli.Command {
	return cli.Command{
		Name:     "hello",
		Usage:    "Start hello server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Hello),
		Action: func(c *cli.Context) error {
			scfg := configureHello(cfg)

			return cli.HandleAction(
				command.Server(scfg).Action,
				c,
			)
		},
	}
}

func configureHello(cfg *config.Config) *svcconfig.Config {
	cfg.Hello.Log.Level = cfg.Log.Level
	cfg.Hello.Log.Pretty = cfg.Log.Pretty
	cfg.Hello.Log.Color = cfg.Log.Color
	cfg.Hello.Tracing.Enabled = false
	cfg.Hello.HTTP.Addr = "localhost:9105"
	cfg.Hello.HTTP.Root = "/"
	cfg.Hello.GRPC.Addr = "localhost:9106"

	return cfg.Hello
}

func init() {
	register.AddCommand(HelloCommand)
}

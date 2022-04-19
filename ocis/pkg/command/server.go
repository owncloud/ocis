package command

import (
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/shared"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/runtime"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "start a fullstack server (runtime and all extensions in supervised mode)",
		Category: "fullstack",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Value:       cfg.ConfigFile,
				Usage:       "config file to be loaded by the extension",
				Destination: &cfg.ConfigFile,
			},
		},
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				logger := log.NewLogger(
					log.Name("oCIS"),
				)
				logger.Error().Err(err).Msg("couldn't find the specified config file")
			}
			return err
		},
		Action: func(c *cli.Context) error {

			cfg.Commons = &shared.Commons{
				Log: cfg.Log,
			}

			r := runtime.New(cfg)
			return r.Start()
		},
	}
}

func init() {
	register.AddCommand(Server)
}

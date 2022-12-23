package command

import (
	"context"
	"fmt"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/http"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
	"github.com/owncloud/ocis/v2/services/hub/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/hub/pkg/service"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", "hub"),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var (
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
			)

			defer cancel()

			httpService, err := http.NewService(
				http.Name(cfg.Service.Name),
				http.Namespace(cfg.HTTP.Namespace),
				http.Version(version.GetString()),
				http.Address(cfg.HTTP.Addr),
				http.Context(ctx),
			)
			if err != nil {
				return err
			}

			if err := micro.RegisterHandler(httpService.Server(), service.New(cfg)); err != nil {
				return err
			}

			return httpService.Run()
		},
	}
}

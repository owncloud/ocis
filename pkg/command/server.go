package command

import (
	"context"
	"fmt"
	"github.com/owncloud/ocis-accounts/pkg/flagset"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/server/grpc"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	oclog "github.com/owncloud/ocis-pkg/v2/log"
)

var (
	logger oclog.Logger
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "an accounts backend manager (driver) needs to be specified. By default the service uses the filesystem as storage",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			logger = oclog.NewLogger(
				oclog.Name(cfg.Server.Name),
				oclog.Level("info"),
				oclog.Color(true),
				oclog.Pretty(true),
			)
			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()
			service := grpc.NewService(
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Server.Name),
				grpc.Namespace(cfg.Server.Namespace),
				grpc.Address(cfg.Server.Address),
			)

			gr.Add(func() error {
				logger.Info().Str("service", service.Name()).Msg("Reporting settings bundle to account service")
				go svc.RegisterSettingsBundles(&logger)
				return service.Run()
			}, func(_ error) {
				fmt.Println("shutting down grpc server")
				cancel()
			})

			run.SignalHandler(ctx, syscall.SIGKILL)
			return gr.Run()
		},
	}
}

package command

import (
	"context"
	"syscall"

	"github.com/owncloud/ocis-accounts/pkg/flagset"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/server/grpc"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "uses an LDAP server as the storage backend",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(c *cli.Context) error {
			return ParseConfig(c, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)
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
				grpc.Flags(flagset.RootWithConfig(config.New())),
			)

			gr.Add(func() error {
				logger.Info().Str("service", service.Name()).Msg("Reporting settings bundle to account service")
				go svc.RegisterSettingsBundles(&logger)
				return service.Run()
			}, func(err error) {
				if err != nil {
					logger.Error().Err(err).Msg("account service died")
				} else {
					logger.Info().
						Str("service", service.Name()).
						Msg("Shutting down server")
				}
				cancel()
			})

			run.SignalHandler(ctx, syscall.SIGKILL)
			return gr.Run()
		},
	}
}

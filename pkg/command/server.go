package command

import (
	"context"
	"fmt"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/micro/grpc"
	olog "github.com/owncloud/ocis-pkg/v2/log"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "an accounts backend manager (driver) needs to be specified. By default the service uses the filesystem as storage",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "manager",
				DefaultText: "filesystem",
				Usage:       "accounts backend manager",
				Value:       "filesystem",
				EnvVars:     []string{"OCIS_ACCOUNTS_MANAGER"},
				Destination: &cfg.Manager,
			},
			&cli.StringFlag{
				Name:        "mount-path",
				Usage:       "mounting point (necessary when manager=filesystem)",
				EnvVars:     []string{"OCIS_ACCOUNTS_MOUNT_PATH"},
				Destination: &cfg.MountPath,
			},
			&cli.StringFlag{
				Name:        "name",
				Value:       "accounts",
				DefaultText: "accounts",
				Usage:       "service name",
				EnvVars:     []string{"OCIS_ACCOUNTS_NAME"},
				Destination: &cfg.Server.Name,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Aliases:     []string{"ns"},
				Value:       "com.owncloud",
				DefaultText: "com.owncloud",
				Usage:       "namespace",
				EnvVars:     []string{"OCIS_ACCOUNTS_NAMESPACE"},
				Destination: &cfg.Server.Namespace,
			},
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"addr"},
				Value:       "localhost:9180",
				DefaultText: "localhost:9180",
				Usage:       "service endpoint",
				EnvVars:     []string{"OCIS_ACCOUNTS_ADDRESS"},
				Destination: &cfg.Server.Address,
			},
		},
		Action: func(c *cli.Context) error {
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())
			l := olog.NewLogger(
				olog.Name(cfg.Server.Name),
			)

			defer cancel()
			service := grpc.NewService(
				grpc.Logger(l),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Server.Name),
				grpc.Namespace(cfg.Server.Namespace),
				grpc.Address(cfg.Server.Address),
			)

			gr.Add(func() error {
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

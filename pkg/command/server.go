package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/micro/grpc"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	baseDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	return &cli.Command{
		Name:  "server",
		Usage: "Start accounts service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "manager",
				DefaultText: "filesystem",
				Usage:       "store controller driver. eg: filesystem",
				Value:       "filesystem",
				EnvVars:     []string{"OCIS_ACCOUNTS_MANAGER"},
				Destination: &cfg.Manager,
			},
			&cli.StringFlag{
				Name:        "mount-path",
				DefaultText: "binary default running location",
				Usage:       "where to mount the ocis accounts store",
				Value:       baseDir,
				EnvVars:     []string{"OCIS_ACCOUNTS_MOUNT_PATH"},
				Destination: &cfg.MountPath,
			},
			&cli.StringFlag{
				Name:        "name",
				Value:       "accounts",
				Destination: &cfg.Server.Name,
			},
			&cli.StringFlag{
				Name:        "namespace",
				Value:       "com.owncloud",
				Destination: &cfg.Server.Namespace,
			},
			&cli.StringFlag{
				Name:        "address",
				Value:       "localhost:9999",
				Destination: &cfg.Server.Address,
			},
		},
		Action: func(c *cli.Context) error {
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()
			service := grpc.NewService(
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

package command

import (
	"context"
	"fmt"
	"syscall"

	"github.com/micro/cli"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/micro/grpc"
)

// Server is the entry point for the server command.
func Server() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Start accounts service",
		Action: func(c *cli.Context) error {
			gr := run.Group{}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()
			service := grpc.NewService(ctx)

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

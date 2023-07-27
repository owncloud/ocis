package command

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v4/client"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/config/parser"
)

// Index is the entrypoint for the server command.
func Index(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "index",
		Usage:    "index the files for one one more users",
		Category: "index management",
		Aliases:  []string{"i"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "space",
				Aliases:  []string{"s"},
				Required: true,
				Usage:    "the id of the space to travers and index the files of",
			},
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Required: true,
				Usage:    "the username of the user that shall be used to access the files",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(ctx *cli.Context) error {
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			grpcClient, err := grpc.NewClient(
				append(grpc.GetClientOptions(cfg.GRPCClientTLS),
					grpc.WithTraceProvider(traceProvider),
				)...,
			)
			if err != nil {
				return err
			}

			c := searchsvc.NewSearchProviderService("com.owncloud.api.search", grpcClient)
			_, err = c.IndexSpace(context.Background(), &searchsvc.IndexSpaceRequest{
				SpaceId: ctx.String("space"),
				UserId:  ctx.String("user"),
			}, func(opts *client.CallOptions) { opts.RequestTimeout = 10 * time.Minute })
			if err != nil {
				fmt.Println("failed to index space: " + err.Error())
				return err
			}
			return nil
		},
	}
}

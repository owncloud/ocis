package command

import (
	"context"
	"errors"
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

// _noTimeout stands in for "no limit" when --timeout is 0: go-micro treats a
// zero RequestTimeout as an immediately expired deadline, so a far-off one is
// used instead.
const _noTimeout = 100 * 365 * 24 * time.Hour

// Index is the entrypoint for the index command.
func Index(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "index",
		Usage:    "index the files for one one more users",
		Category: "index management",
		Aliases:  []string{"i"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "space",
				Aliases: []string{"s"},
				Usage:   "the id of the space to travers and index the files of. This or --all-spaces is required.",
			},
			&cli.BoolFlag{
				Name:  "all-spaces",
				Usage: "index all spaces instead. This or --space is required.",
			},
			&cli.DurationFlag{
				Name:  "timeout",
				Usage: "how long to wait for the indexing to finish before giving up, e.g. 2h. The walk is cancelled when the timeout expires. 0 waits indefinitely.",
				Value: 10 * time.Minute,
			},
		},
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(ctx *cli.Context) error {
			if ctx.String("space") == "" && !ctx.Bool("all-spaces") {
				return errors.New("either --space or --all-spaces is required")
			}

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

			// The service walks the whole space inside this request, so the
			// request timeout bounds the walk: when it expires the walk is
			// cancelled and everything after that point stays unindexed.
			timeout := ctx.Duration("timeout")
			if timeout <= 0 {
				timeout = _noTimeout
			}

			c := searchsvc.NewSearchProviderService("com.owncloud.api.search", grpcClient)
			_, err = c.IndexSpace(context.Background(), &searchsvc.IndexSpaceRequest{
				SpaceId: ctx.String("space"),
			}, func(opts *client.CallOptions) { opts.RequestTimeout = timeout })
			if err != nil {
				fmt.Println("failed to index space: " + err.Error())
				return err
			}
			return nil
		},
	}
}

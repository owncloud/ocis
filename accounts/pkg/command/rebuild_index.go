package command

import (
	"context"
	"fmt"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/owncloud/ocis/accounts/pkg/config"
	index "github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

// DeleteIndex rebuilds the entire configured index.
func DeleteIndex(cdf *config.Config) *cli.Command {
	return &cli.Command{
		Name:    "rebuildIndex",
		Usage:   "Rebuilds the service's index",
		Aliases: []string{"rebuild", "ri"},
		Action: func(ctx *cli.Context) error {
			idxSvcID := "com.owncloud.api.accounts"
			idxSvc := index.NewIndexService(idxSvcID, grpc.NewClient())

			rsp, err := idxSvc.RebuildIndex(context.Background(), &index.RebuildIndexRequest{})
			if err != nil {
				return err
			}

			fmt.Printf("deleted: %+v", rsp.Indices)
			return nil
		},
	}
}

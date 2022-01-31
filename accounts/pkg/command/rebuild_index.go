package command

import (
	"context"
	"fmt"

	accountssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/accounts/v0"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/urfave/cli/v2"
	merrors "go-micro.dev/v4/errors"
)

// RebuildIndex rebuilds the entire configured index.
func RebuildIndex(cdf *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "rebuildIndex",
		Usage:    "rebuilds the service's index, i.e. deleting and then re-adding all existing documents",
		Category: "account management",
		Aliases:  []string{"rebuild", "ri"},
		Action: func(ctx *cli.Context) error {
			idxSvcID := "com.owncloud.api.accounts"
			idxSvc := accountssvc.NewIndexService(idxSvcID, grpc.NewClient())

			_, err := idxSvc.RebuildIndex(context.Background(), &accountssvc.RebuildIndexRequest{})
			if err != nil {
				fmt.Println(merrors.FromError(err).Detail)
				return err
			}

			fmt.Println("index rebuilt successfully")
			return nil
		},
	}
}

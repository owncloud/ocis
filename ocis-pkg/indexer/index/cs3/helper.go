package cs3

import (
	"context"
	"fmt"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/utils"
)

func deleteIndexRoot(ctx context.Context, storageProvider provider.ProviderAPIClient, spaceid, indexRootDir string) error {
	res, err := storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: spaceid,
				OpaqueId:  spaceid,
			},
			Path: utils.MakeRelativePath(indexRootDir),
		},
	})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("error deleting index root dir: %v", indexRootDir)
	}

	return nil
}

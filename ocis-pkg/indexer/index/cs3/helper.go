package cs3

import (
	"context"
	"fmt"
	"path"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
)

func deleteIndexRoot(ctx context.Context, storageProvider provider.ProviderAPIClient, indexRootDir string) error {
	res, err := storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{Path: path.Join("/meta", indexRootDir)},
	})
	if err != nil {
		return err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("error deleting index root dir: %v", indexRootDir)
	}

	return nil
}

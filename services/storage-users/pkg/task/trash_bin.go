package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	apiGateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	apiRpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	apiProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/utils"
	"google.golang.org/grpc/metadata"
)

// PurgeTrashBin can be used to purge space trash-bin's,
// the provided executantId must have space access and delete permissions.
// removeBefore specifies how long an item must be in the trash-bin to be deleted,
// items that stay there for a shorter time are ignored.
func PurgeTrashBin(gw apiGateway.GatewayAPIClient, clientSecret, executantId string, removeBefore time.Time) error {
	executantToken, err := getToken(gw, clientSecret, executantId)
	if err != nil {
		return err
	}

	executantCtx := metadata.AppendToOutgoingContext(context.Background(), ctxpkg.TokenHeader, executantToken)
	listStorageSpacesResponse, err := gw.ListStorageSpaces(executantCtx, &apiProvider.ListStorageSpacesRequest{
		// ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE cant handle `AND` (`project` && `personal`) chains,
		// therefore we need to request all space types and filter them out later (`virtual`).
	})
	if err != nil {
		return err
	}

	for _, storageSpace := range listStorageSpacesResponse.StorageSpaces {
		var (
			err                   error
			impersonationUserID   string
			storageSpaceReference = &apiProvider.Reference{
				ResourceId: storageSpace.GetRoot(),
			}
		)

		switch storageSpace.GetSpaceType() {
		case "personal":
			impersonationUserID = storageSpace.GetOwner().GetId().GetOpaqueId()
		case "project":
			opaqueGrants, ok := storageSpace.GetOpaque().GetMap()["grants"]
			if !ok {
				err = errors.New("no grants")
				break
			}

			var permissionsMap map[string]*apiProvider.ResourcePermissions
			err = json.Unmarshal(opaqueGrants.Value, &permissionsMap)
			if err != nil {
				break
			}

			for id, permissions := range permissionsMap {
				if permissions.Delete != true {
					continue
				}

				impersonationUserID = id
				break
			}
		default:
			continue
		}

		if impersonationUserID == "" {
			err = fmt.Errorf("cant impersonate space user for space: %s", storageSpace.GetId())
		}

		if err != nil {
			return err
		}

		impersonationToken, err := getToken(gw, clientSecret, impersonationUserID)
		if err != nil {
			return err
		}

		impersonatedCtx := metadata.AppendToOutgoingContext(context.Background(), ctxpkg.TokenHeader, impersonationToken)
		listRecycleResponse, err := gw.ListRecycle(impersonatedCtx, &apiProvider.ListRecycleRequest{Ref: storageSpaceReference})

		for _, recycleItem := range listRecycleResponse.GetRecycleItems() {
			doDelete := utils.TSToUnixNano(recycleItem.DeletionTime) < utils.TSToUnixNano(utils.TimeToTS(removeBefore))
			if !doDelete {
				continue
			}

			purgeRecycleResponse, err := gw.PurgeRecycle(impersonatedCtx, &apiProvider.PurgeRecycleRequest{
				Ref: storageSpaceReference,
				Key: recycleItem.Key,
			})

			if purgeRecycleResponse.GetStatus().GetCode() != apiRpc.Code_CODE_OK {
				return errtypes.NewErrtypeFromStatus(purgeRecycleResponse.Status)
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

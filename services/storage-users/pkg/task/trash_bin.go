package task

import (
	"fmt"
	"time"

	apiGateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	apiUser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	apiRpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	apiProvider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
)

// PurgeTrashBin can be used to purge space trash-bin's,
// the provided executantID must have space access.
// removeBefore specifies how long an item must be in the trash-bin to be deleted,
// items that stay there for a shorter time are ignored and kept in place.
func PurgeTrashBin(executantID *apiUser.UserId, deleteBefore time.Time, spaceType SpaceType, gatewaySelector pool.Selectable[apiGateway.GatewayAPIClient], machineAuthAPIKey string) error {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return err
	}

	executantCtx, _, err := utils.Impersonate(executantID, gatewayClient, machineAuthAPIKey)
	if err != nil {
		return err
	}

	listStorageSpacesResponse, err := gatewayClient.ListStorageSpaces(executantCtx, &apiProvider.ListStorageSpacesRequest{
		Filters: []*apiProvider.ListStorageSpacesRequest_Filter{
			{
				Type: apiProvider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &apiProvider.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: string(spaceType),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	for _, storageSpace := range listStorageSpacesResponse.StorageSpaces {
		var (
			err                   error
			impersonationID       *apiUser.UserId
			storageSpaceReference = &apiProvider.Reference{
				ResourceId: storageSpace.GetRoot(),
			}
		)

		switch SpaceType(storageSpace.GetSpaceType()) {
		case Personal:
			impersonationID = storageSpace.GetOwner().GetId()
		case Project:
			var permissionsMap map[string]*apiProvider.ResourcePermissions
			err := utils.ReadJSONFromOpaque(storageSpace.GetOpaque(), "grants", &permissionsMap)
			if err != nil {
				break
			}

			for id, permissions := range permissionsMap {
				if !permissions.Delete {
					continue
				}

				impersonationID = &apiUser.UserId{
					OpaqueId: id,
				}
				break
			}
		default:
			continue
		}

		if err != nil {
			return err
		}

		if impersonationID == nil {
			return fmt.Errorf("can't impersonate space user for space: %s", storageSpace.GetId().GetOpaqueId())
		}

		impersonatedCtx, _, err := utils.Impersonate(impersonationID, gatewayClient, machineAuthAPIKey)
		if err != nil {
			return err
		}

		listRecycleResponse, err := gatewayClient.ListRecycle(impersonatedCtx, &apiProvider.ListRecycleRequest{Ref: storageSpaceReference})
		if err != nil {
			return err
		}

		for _, recycleItem := range listRecycleResponse.GetRecycleItems() {
			doDelete := utils.TSToUnixNano(recycleItem.DeletionTime) < utils.TSToUnixNano(utils.TimeToTS(deleteBefore))
			if !doDelete {
				continue
			}

			purgeRecycleResponse, err := gatewayClient.PurgeRecycle(impersonatedCtx, &apiProvider.PurgeRecycleRequest{
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

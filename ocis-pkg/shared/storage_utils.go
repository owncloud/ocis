package shared

import (
	"context"
	"errors"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/reva/v2/pkg/appctx"
	"github.com/owncloud/reva/v2/pkg/utils"
)

const (
	_spaceTypePersonal = "personal"
	_spaceStateTrashed = "trashed"
)

var (
	// ErrNotFound is returned when a personal space not found.
	ErrNotFound = errors.New("personal space not found")
)

// DisablePersonalSpace disables (deletes) the personal storage space for the given userID.
// If the personal space is already deleted (trashed), it is a no-op.
// Returns an error if the space cannot be found or the deletion fails.
func DisablePersonalSpace(ctx context.Context, client gateway.GatewayAPIClient, userID string) error {
	logger := appctx.GetLogger(ctx)
	sp, err := getPersonalSpace(ctx, client, userID)
	if err == ErrNotFound {
		logger.Debug().Str("userID", userID).Msg("no personal space found to delete")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve personal space: %w", err)
	}
	if isTrashed(sp) {
		logger.Debug().Str("userID", userID).Msg("the personal space already deleted")
		return nil
	}

	dRes, derr := client.DeleteStorageSpace(ctx, &storageprovider.DeleteStorageSpaceRequest{
		Id: &storageprovider.StorageSpaceId{OpaqueId: sp.GetId().GetOpaqueId()},
	})
	if derr != nil {
		return fmt.Errorf("failed to disable personal space: %w", derr)
	}
	if dRes.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
		return errorcode.NewFromStatusCode(dRes.GetStatus().GetCode(), dRes.GetStatus().GetMessage())
	}
	return nil
}

// RestorePersonalSpace ensures that a personal storage space exists and is enabled for the given user.
// If the personal space is found and is trashed, it restores it.
// Returns an error if the operation fails.
func RestorePersonalSpace(ctx context.Context, client gateway.GatewayAPIClient, userID string) error {
	sp, err := getPersonalSpace(ctx, client, userID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	if sp != nil && isTrashed(sp) {
		req := &storageprovider.UpdateStorageSpaceRequest{
			Opaque: utils.AppendPlainToOpaque(nil, "restore", "true"),
			StorageSpace: &storageprovider.StorageSpace{
				Id:   &storageprovider.StorageSpaceId{OpaqueId: sp.GetId().GetOpaqueId()},
				Root: sp.GetRoot(),
			},
		}
		uRes, uerr := client.UpdateStorageSpace(ctx, req)
		if uerr != nil {
			return fmt.Errorf("failed to enable personal space: %w", uerr)
		}
		if uRes.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
			return errorcode.NewFromStatusCode(uRes.GetStatus().GetCode(), uRes.GetStatus().GetMessage())
		}
		return nil
	}
	return nil
}

func getPersonalSpace(ctx context.Context, client gateway.GatewayAPIClient, userID string) (*storageprovider.StorageSpace, error) {
	lspr, err := client.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{
		Opaque: utils.AppendPlainToOpaque(nil, "unrestricted", "T"),
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			listStorageSpacesUserFilter(userID),
			listStorageSpacesTypeFilter(_spaceTypePersonal)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve personal space: %w", err)
	}
	if lspr.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
		return nil, fmt.Errorf("failed to retrieve personal space: %s", lspr.GetStatus().GetMessage())
	}
	if len(lspr.GetStorageSpaces()) > 1 {
		return nil, errors.New("retrieved more than one personal space")
	}
	if len(lspr.GetStorageSpaces()) != 1 {
		return nil, ErrNotFound
	}
	return lspr.GetStorageSpaces()[0], nil
}

func listStorageSpacesUserFilter(id string) *storageprovider.ListStorageSpacesRequest_Filter {
	return &storageprovider.ListStorageSpacesRequest_Filter{
		Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_USER,
		Term: &storageprovider.ListStorageSpacesRequest_Filter_User{
			User: &userv1beta1.UserId{
				OpaqueId: id,
			},
		},
	}
}

func listStorageSpacesTypeFilter(spaceType string) *storageprovider.ListStorageSpacesRequest_Filter {
	return &storageprovider.ListStorageSpacesRequest_Filter{
		Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
		Term: &storageprovider.ListStorageSpacesRequest_Filter_SpaceType{
			SpaceType: spaceType,
		},
	}
}

func isTrashed(sp *storageprovider.StorageSpace) bool {
	return utils.ReadPlainFromOpaque(sp.GetOpaque(), _spaceStateTrashed) == _spaceStateTrashed
}

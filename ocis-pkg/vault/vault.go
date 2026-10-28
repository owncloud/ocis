// Package vault provides helpers for resources in the vault storage provider.
package vault

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/permission"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
)

// IsVaultStorage reports whether the storage ID belongs to the vault storage provider.
func IsVaultStorage(storageID string) bool {
	return storageID == utils.VaultStorageProviderID
}

// IsVaultSpace reports whether the space ID belongs to the vault storage provider.
// It returns true for a space ID that it cannot parse, so that callers fail closed.
func IsVaultSpace(spaceID string) bool {
	rid, err := storagespace.ParseID(spaceID)
	if err != nil {
		return true
	}
	return IsVaultStorage(rid.GetStorageId())
}

// HasVaultMode reports whether the user holds the VaultMode permission.
// A non-nil error means undetermined, so callers must treat the user as not allowed.
func HasVaultMode(ctx context.Context, gwc gateway.GatewayAPIClient, userID *user.UserId) (bool, error) {
	res, err := gwc.CheckPermission(ctx, &permissions.CheckPermissionRequest{
		SubjectRef: &permissions.SubjectReference{
			Spec: &permissions.SubjectReference_UserId{UserId: userID},
		},
		Permission: permission.VaultMode,
	})
	if err != nil {
		return false, err
	}
	switch code := res.GetStatus().GetCode(); code {
	case rpc.Code_CODE_OK:
		return true, nil
	case rpc.Code_CODE_PERMISSION_DENIED:
		return false, nil
	default:
		return false, fmt.Errorf("vault permission check returned status %q: %s", code, res.GetStatus().GetMessage())
	}
}

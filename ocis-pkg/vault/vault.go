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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

// maxConcurrentChecks bounds how many VaultMode checks run at the same time, so that a large
// group does not send all its checks to the gateway at once.
const maxConcurrentChecks = 10

const tracerName = "github.com/owncloud/ocis/v2/ocis-pkg/vault"

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

// FilterVaultModeUsers returns the users from userIDs who hold the VaultMode permission, in the
// same order. A group can contain users without the permission, because the gateway does not
// check the group members. The checks run in parallel, with a bound. When a check fails, the
// function calls onError for that user and leaves the user out, so that the caller fails closed.
func FilterVaultModeUsers(ctx context.Context, gwc gateway.GatewayAPIClient, userIDs []*user.UserId, onError func(*user.UserId, error)) []*user.UserId {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(tracerName).Start(ctx, "FilterVaultModeUsers")
	defer span.End()
	span.SetAttributes(attribute.Int("vault.users.checked", len(userIDs)))

	allowed := make([]bool, len(userIDs))
	errs := make([]error, len(userIDs))
	var eg errgroup.Group
	eg.SetLimit(maxConcurrentChecks)
	for i, id := range userIDs {
		eg.Go(func() error {
			allowed[i], errs[i] = HasVaultMode(ctx, gwc, id)
			return nil
		})
	}
	_ = eg.Wait()

	filtered := make([]*user.UserId, 0, len(userIDs))
	failed := 0
	for i, id := range userIDs {
		if errs[i] != nil {
			failed++
			if onError != nil {
				onError(id, errs[i])
			}
			continue
		}
		if allowed[i] {
			filtered = append(filtered, id)
		}
	}
	span.SetAttributes(
		attribute.Int("vault.users.allowed", len(filtered)),
		attribute.Int("vault.users.failed", failed),
	)
	return filtered
}

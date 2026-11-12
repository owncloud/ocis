package auth

import (
	"context"
	"sync"
	"time"

	"github.com/bluele/gcache"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/permission"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/pkg/errors"
)

const vaultModeRequiredMsg = "vault mode permission required to access vault storage"

var (
	vaultModeCache     gcache.Cache
	vaultModeCacheTTL  time.Duration
	vaultModeCacheOnce sync.Once

	// getGatewayClient is a variable, so that tests can replace the gateway client.
	getGatewayClient = func(addr string) (gatewayv1beta1.GatewayAPIClient, error) {
		return pool.GetGatewayServiceClient(addr)
	}
)

// initVaultModeCache creates the cache for the VaultMode answers. Only the first call has an effect.
func initVaultModeCache(size int, ttl time.Duration) {
	vaultModeCacheOnce.Do(func() {
		vaultModeCache = gcache.New(size).LRU().Build()
		vaultModeCacheTTL = ttl
	})
}

// hasVaultMode reports whether the user holds the VaultMode permission.
// The vault storage provider uses it to deny access to users who cannot use vault mode,
// for example group members who received a vault share. A non-nil error means undetermined,
// so the caller must deny the request. Only definitive answers are cached.
func hasVaultMode(ctx context.Context, u *userpb.User, gatewayAddr string) (bool, error) {
	id := u.GetId().GetOpaqueId()
	if v, err := vaultModeCache.Get(id); err == nil {
		return v.(bool), nil
	}

	client, err := getGatewayClient(gatewayAddr)
	if err != nil {
		return false, err
	}
	res, err := client.CheckPermission(ctx, &permissions.CheckPermissionRequest{
		SubjectRef: &permissions.SubjectReference{
			Spec: &permissions.SubjectReference_UserId{UserId: u.GetId()},
		},
		Permission: permission.VaultMode,
	})
	if err != nil {
		return false, errors.Wrap(err, "auth: error calling CheckPermission")
	}

	var allowed bool
	switch code := res.GetStatus().GetCode(); code {
	case rpc.Code_CODE_OK:
		allowed = true
	case rpc.Code_CODE_PERMISSION_DENIED:
		allowed = false
	default:
		return false, errors.Errorf("auth: vault permission check returned status %q: %s", code, res.GetStatus().GetMessage())
	}

	_ = vaultModeCache.SetWithExpire(id, allowed, vaultModeCacheTTL)
	return allowed, nil
}

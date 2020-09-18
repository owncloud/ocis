package roles

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/owncloud/ocis-pkg/v2/middleware"
)

// ReadRoleIDsFromContext extracts roleIDs from the metadata context and returns them as []string
func ReadRoleIDsFromContext(ctx context.Context) (roleIDs []string, ok bool) {
	roleIDsJSON, ok := metadata.Get(ctx, middleware.RoleIDs)
	if !ok {
		return nil, false
	}
	err := json.Unmarshal([]byte(roleIDsJSON), &roleIDs)
	if err != nil {
		return nil, false
	}
	return roleIDs, true
}

package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/micro/go-micro/v2/metadata"
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/roles"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

func Roles(log log.Logger, rs settings.RoleService, cache *roles.Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get roleIDs from context
			roleIDs, ok := ReadRoleIDsFromContext(r.Context())
			if !ok {
				log.Debug().Msg("failed to read roleIDs from context")
				next.ServeHTTP(w, r)
				return
			}

			// check which roles are not cached, yet
			lookup := make([]string, 0)
			for _, roleID := range roleIDs {
				if hit := cache.Get(roleID); hit == nil {
					lookup = append(lookup, roleID)
				}
			}

			// fetch roles
			if len(lookup) > 0 {
				request := &settings.ListBundlesRequest{
					BundleIds: lookup,
				}
				res, err := rs.ListRoles(r.Context(), request)
				if err != nil {
					log.Debug().Err(err).Msg("failed to fetch roles by roleIDs")
					next.ServeHTTP(w, r)
					return
				}
				for _, role := range res.Bundles {
					cache.Set(role.Id, role)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ReadRoleIDsFromContext extracts roleIDs from the metadata context and returns them as []string
func ReadRoleIDsFromContext(ctx context.Context) (roleIDs []string, ok bool) {
	roleIDsJson, ok := metadata.Get(ctx, RoleIDs)
	if !ok {
		return nil, false
	}
	err := json.Unmarshal([]byte(roleIDsJson), &roleIDs)
	if err != nil {
		return nil, false
	}
	return roleIDs, true
}

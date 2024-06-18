package middleware

import (
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
)

// RequireAdmin middleware is used to require the user in context to be an admin / have account management permissions
func RequireAdmin(rm *roles.Manager, logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		l := logger.With().Str("middleware", "requireAdmin").Logger()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := revactx.ContextGetUser(r.Context())
			if !ok {
				errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if u.Id == nil || u.Id.OpaqueId == "" {
				errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user is missing an id")
				return
			}
			// get roles from context
			roleIDs, ok := roles.ReadRoleIDsFromContext(r.Context())
			if !ok {
				l.Debug().Str("userid", u.Id.OpaqueId).Msg("No roles in context, contacting settings service")
				var err error
				roleIDs, err = rm.FindRoleIDsForUser(r.Context(), u.Id.OpaqueId)
				if err != nil {
					l.Error().Err(err).Str("userid", u.Id.OpaqueId).Msg("Failed to get roles for user")
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
				if len(roleIDs) == 0 {
					l.Error().Err(err).Str("userid", u.Id.OpaqueId).Msg("No roles assigned to user")
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
			}

			// check if permission is present in roles of the authenticated account
			if rm.FindPermissionByID(r.Context(), roleIDs, settings.AccountManagementPermissionID) != nil {
				next.ServeHTTP(w, r)
				return
			}

			errorcode.AccessDenied.Render(w, r, http.StatusForbidden, "Forbidden")
		})
	}
}

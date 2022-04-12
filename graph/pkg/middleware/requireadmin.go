package middleware

import (
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/service/v0"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/roles"
)

// RequireAdmin middleware is used to require the user in context to be an admin / have account management permissions
func RequireAdmin(rm *roles.Manager, logger log.Logger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
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
				logger.Debug().Str("userid", u.Id.OpaqueId).Msg("No roles in context, contacting settings service")
				var err error
				roleIDs, err = rm.FindRoleIDsForUser(r.Context(), u.Id.OpaqueId)
				if err != nil {
					logger.Err(err).Str("userid", u.Id.OpaqueId).Msg("failed to get roles for user")
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
				if len(roleIDs) == 0 {
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
			}

			// check if permission is present in roles of the authenticated account
			if rm.FindPermissionByID(r.Context(), roleIDs, accounts.AccountManagementPermissionID) != nil {
				next.ServeHTTP(w, r)
				return
			}

			errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
		})
	}
}

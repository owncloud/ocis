package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	accounts "github.com/owncloud/ocis/accounts/pkg/service/v0"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// RequireAdmin middleware is used to require the user in context to be an admin / have account management permissions
func RequireAdmin(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	mustRender := func(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
		if err := render.Render(w, r, renderer); err != nil {
			opt.Logger.Err(err).Msgf("failed to write response for ocs request %s on %s", r.Method, r.URL)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// get roles from context
			roleIDs, ok := roles.ReadRoleIDsFromContext(r.Context())
			if !ok {
				mustRender(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized"))
				return
			}

			// check if permission is present in roles of the authenticated account
			if opt.RoleManager.FindPermissionByID(r.Context(), roleIDs, accounts.AccountManagementPermissionID) != nil {
				next.ServeHTTP(w, r)
				return
			}

			mustRender(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized"))
		})
	}
}

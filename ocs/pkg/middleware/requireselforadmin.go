package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	accounts "github.com/owncloud/ocis/accounts/pkg/service/v0"
	"github.com/owncloud/ocis/ocis-pkg/roles"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// RequireSelfOrAdmin middleware is used to require the requesting user to be an admin or the requested user himself
func RequireSelfOrAdmin(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			u, ok := user.ContextGetUser(r.Context())
			if !ok {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized")))
				return
			}
			if u.Id == nil || u.Id.OpaqueId == "" {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "user is missing an id")))
				return
			}
			// get roles from context
			roleIDs, ok := roles.ReadRoleIDsFromContext(r.Context())
			if !ok {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized")))
				return
			}

			// check if account management permission is present in roles of the authenticated account
			if opt.RoleManager.FindPermissionByID(r.Context(), roleIDs, accounts.AccountManagementPermissionID) != nil {
				next.ServeHTTP(w, r)
				return
			}

			// check if self management permission is present in roles of the authenticated account
			if opt.RoleManager.FindPermissionByID(r.Context(), roleIDs, accounts.SelfManagementPermissionID) != nil {
				userid := chi.URLParam(r, "userid")
				if userid == "" || userid == u.Id.OpaqueId || userid == u.Username {
					next.ServeHTTP(w, r)
					return
				}
			}

			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized")))

		})
	}
}

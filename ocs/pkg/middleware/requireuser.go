package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// RequireUser middleware is used to require a user in context
func RequireUser() func(next http.Handler) http.Handler {

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

			next.ServeHTTP(w, r)

		})
	}
}

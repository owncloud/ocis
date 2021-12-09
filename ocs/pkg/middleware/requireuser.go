package middleware

import (
	"net/http"

	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// RequireUser middleware is used to require a user in context
func RequireUser(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	mustRender := func(w http.ResponseWriter, r *http.Request, renderer render.Renderer) {
		if err := render.Render(w, r, renderer); err != nil {
			opt.Logger.Err(err).Msgf("failed to write response for ocs request %s on %s", r.Method, r.URL)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			u, ok := revactx.ContextGetUser(r.Context())
			if !ok {
				mustRender(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "Unauthorized"))
				return
			}
			if u.Id == nil || u.Id.OpaqueId == "" {
				mustRender(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "user is missing an id"))
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}

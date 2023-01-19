package authorization

import (
	"context"
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// Authorizer is the common interface implemented by all request authorizers.
type Authorizer interface {
	// Authorize is used to authorize incoming HTTP requests.
	Authorize(context.Context, Info) (bool, error)
}

// Info contains every data that is needed to decide if the request should pass or not
type Info struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Authorization is a higher order authorization middleware.
func Authorization(logger log.Logger, auths []Authorizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := Info{
				Method: r.Method,
				Path:   r.URL.Path,
			}

			for _, a := range auths {
				allowed, err := a.Authorize(r.Context(), info)
				if err != nil {
					logger.Err(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if !allowed {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

package authorization

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"net/http"
)

// Authorizer is the common interface implemented by all request authorizers.
type Authorizer interface {
	// Authorize is used to authorize incoming HTTP requests.
	Authorize(r *http.Request) (bool, error)
}

// Authorization is a higher order authorization middleware.
func Authorization(logger log.Logger, auths []Authorizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, a := range auths {
				allowed, err := a.Authorize(r)
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

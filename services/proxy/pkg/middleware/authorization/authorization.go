package authorization

import (
	"bytes"
	"context"
	"io"
	"net/http"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// Authorizer is the common interface implemented by all request authorizers.
type Authorizer interface {
	// Authorize is used to authorize incoming HTTP requests.
	Authorize(context.Context, Info) (bool, error)
}

// Info contains every data that is needed to decide if the request should pass or not
type Info struct {
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   []byte      `json:"body"`
	User   userpb.User `json:"user"`
}

// Authorization is a higher order authorization middleware.
func Authorization(logger log.Logger, auths []Authorizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := Info{
				Method: r.Method,
				Path:   r.URL.Path,
			}

			if r.Body != nil {
				info.Body, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(info.Body))
			}

			if user, ok := revactx.ContextGetUser(r.Context()); ok {
				info.User = *user
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

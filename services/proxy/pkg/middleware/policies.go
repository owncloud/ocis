package middleware

import (
	"net/http"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	pMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
	pService "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
)

// Policies verifies if a request is granted or not.
func Policies(logger log.Logger, enabled bool, qs string) func(next http.Handler) http.Handler {
	pClient := pService.NewPoliciesProviderService("com.owncloud.api.policies", grpc.DefaultClient())

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled {
				next.ServeHTTP(w, r)
				return
			}

			req := &pService.EvaluateRequest{
				Query: qs,
				Environment: &pMessage.Environment{
					Request: &pMessage.Request{
						Method: r.Method,
						Path:   r.URL.Path,
					},
					Stage: pMessage.Stage_STAGE_HTTP,
				},
			}

			if user, ok := revactx.ContextGetUser(r.Context()); ok {
				req.Environment.User = &pMessage.User{
					Id: &pMessage.User_ID{
						OpaqueId: user.GetId().GetOpaqueId(),
					},
					Username:    user.GetUsername(),
					Mail:        user.GetMail(),
					DisplayName: user.GetDisplayName(),
					Groups:      user.GetGroups(),
				}
			}

			rsp, err := pClient.Evaluate(r.Context(), req)
			if err != nil {
				logger.Err(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !rsp.Result {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

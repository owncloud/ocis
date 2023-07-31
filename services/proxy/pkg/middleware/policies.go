package middleware

import (
	"net/http"
	"path"
	"time"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	pMessage "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
	pService "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	tusd "github.com/tus/tusd/pkg/handler"
	"go-micro.dev/v4/client"
)

type (
	// RequestDenied struct for OdataErrorMain
	RequestDenied struct {
		Error RequestDeniedError `json:"error"`
	}

	// RequestDeniedError struct for RequestDenied
	RequestDeniedError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		// The structure of this object is service-specific
		Innererror map[string]interface{} `json:"innererror,omitempty"`
	}
)

const DeniedMessage = "Operation denied due to security policies"

// Policies verifies if a request is granted or not.
func Policies(logger log.Logger, qs string, grpcClient client.Client) func(next http.Handler) http.Handler {
	pClient := pService.NewPoliciesProviderService("com.owncloud.api.policies", grpcClient)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if qs == "" {
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

			meta := tusd.ParseMetadataHeader(r.Header.Get(net.HeaderUploadMetadata))
			req.Environment.Resource = &pMessage.Resource{
				Name: meta["filename"],
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
				logger.Err(err).Msg("error evaluating request")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !rsp.Result {
				RenderError(w, r, req, http.StatusForbidden, DeniedMessage)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RenderError writes a Policies ErrorObject to the response writer
func RenderError(w http.ResponseWriter, r *http.Request, evaluateReq *pService.EvaluateRequest, status int, msg string) {
	filename := evaluateReq.Environment.GetResource().GetName()
	if filename == "" {
		filename = path.Base(evaluateReq.Environment.GetRequest().GetPath())
	}

	innererror := map[string]interface{}{
		"date": time.Now().UTC().Format(time.RFC3339),
	}

	innererror["request-id"] = middleware.GetReqID(r.Context())
	innererror["method"] = evaluateReq.Environment.GetRequest().GetMethod()
	innererror["filename"] = filename
	innererror["path"] = evaluateReq.Environment.GetRequest().GetPath()

	resp := &RequestDenied{
		Error: RequestDeniedError{
			Code:       "deniedByPolicy",
			Message:    msg,
			Innererror: innererror,
		},
	}
	render.Status(r, status)
	render.JSON(w, r, resp)
}

package middleware

import (
	"net/http"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	tokenpkg "github.com/cs3org/reva/pkg/token"
	"google.golang.org/grpc/metadata"
)

// CreateHome provides a middleware which sends a CreateHome request to the reva gateway
func CreateHome(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-access-token")

			// we need to pass the token to authenticate the CreateHome request.
			//ctx := tokenpkg.ContextSetToken(r.Context(), token)
			ctx := metadata.AppendToOutgoingContext(r.Context(), tokenpkg.TokenHeader, token)

			createHomeReq := &provider.CreateHomeRequest{}
			createHomeRes, err := opt.RevaGatewayClient.CreateHome(ctx, createHomeReq)

			if err != nil {
				opt.Logger.Err(err).Msg("error calling CreateHome")
			}

			if createHomeRes.Status.Code != rpc.Code_CODE_OK {
				err := status.NewErrorFromCode(createHomeRes.Status.Code, "gateway")
				opt.Logger.Err(err).Msg("error calling Createhome")
			}

			next.ServeHTTP(w, r)
		})
	}
}

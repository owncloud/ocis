package middleware

import (
	"net/http"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	tokenpkg "github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"google.golang.org/grpc/metadata"
)

// CreateHome provides a middleware which sends a CreateHome request to the reva gateway
func CreateHome(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accounts := opt.AccountsClient

			tokenManager, err := jwt.New(map[string]interface{}{
				"secret": opt.TokenManagerConfig.JWTSecret,
			})
			if err != nil {
				opt.Logger.Error().Err(err).Msg("error creating a token manager")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			token := r.Header.Get("x-access-token")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := tokenManager.DismantleToken(r.Context(), token)
			if err != nil {
				opt.Logger.Err(err).Msg("error getting user from access token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = accounts.GetAccount(r.Context(), &proto.GetAccountRequest{
				Id: user.Id.OpaqueId,
			})
			if err != nil {
				e := errors.Parse(err.Error())
				if e.Code == http.StatusNotFound {
					opt.Logger.Debug().Msgf("account with id %s not found", user.Id.OpaqueId)
					next.ServeHTTP(w, r)
					return
				}
				opt.Logger.Err(err).Msgf("error getting user with id %s from accounts service", user.Id.OpaqueId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

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

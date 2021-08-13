package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/auth/scope"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"google.golang.org/grpc/metadata"
)

// authOptions initializes the available default options.
func authOptions(opts ...account.Option) account.Options {
	opt := account.Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Auth provides a middleware to authenticate requestrs using the x-access-token header value
// and write it to the context. If there is no x-access-token the middleware prevents access and renders a json document.
func Auth(opts ...account.Option) func(http.Handler) http.Handler {
	opt := authOptions(opts...)
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret":  opt.JWTSecret,
		"expires": int64(60),
	})
	if err != nil {
		opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			t := r.Header.Get("x-access-token")
			if t == "" {
				errorcode.InvalidAuthenticationToken.Render(w, r, http.StatusUnauthorized, "Access token is empty.")
				/* msgraph error for GET https://graph.microsoft.com/v1.0/me
				{
					"error":
				 	{
						"code":"InvalidAuthenticationToken",
						"message":"Access token is empty.",
						"innerError":{
							"date":"2021-07-09T14:40:51",
							"request-id":"bb12f7db-b4c4-43a9-ba4b-31676aeed019",
							"client-request-id":"bb12f7db-b4c4-43a9-ba4b-31676aeed019"
						}
					}
				}
				*/
				return
			}

			u, tokenScope, err := tokenManager.DismantleToken(r.Context(), t)
			if err != nil {
				errorcode.InvalidAuthenticationToken.Render(w, r, http.StatusUnauthorized, "invalid token")
				return
			}
			if ok, err := scope.VerifyScope(tokenScope, r); err != nil || !ok {
				opt.Logger.Error().Err(err).Msg("verifying scope failed")
				errorcode.InvalidAuthenticationToken.Render(w, r, http.StatusUnauthorized, "verifying scope failed")
				return
			}

			ctx = revactx.ContextSetToken(ctx, t)
			ctx = revactx.ContextSetUser(ctx, u)
			ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

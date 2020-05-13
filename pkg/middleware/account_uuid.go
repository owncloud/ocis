package middleware

import (
	"context"
	"net/http"

	revauser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	mclient "github.com/micro/go-micro/v2/client"
	acc "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-pkg/v2/log"
	ocisoidc "github.com/owncloud/ocis-pkg/v2/oidc"
)

// AccountMiddlewareOption defines a single option function.
type AccountMiddlewareOption func(o *AccountMiddlewareOptions)

// AccountMiddlewareOptions defines the available options for this package.
type AccountMiddlewareOptions struct {
	// Logger to use for logging, must be set
	Logger log.Logger
}

// Logger provides a function to set the logger option.
func Logger(l log.Logger) AccountMiddlewareOption {
	return func(o *AccountMiddlewareOptions) {
		o.Logger = l
	}
}

// AccountUUID provides a middleware which mints a jwt and adds it to the proxied request based
// on the oidc-claims
func AccountUUID(opts ...AccountMiddlewareOption) func(next http.Handler) http.Handler {
	opt := AccountMiddlewareOptions{}
	for _, o := range opts {
		o(&opt)
	}

	return func(next http.Handler) http.Handler {
		// TODO: handle error
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  "Pive-Fumkiu4",
			"expires": int64(60),
		})

		if err != nil {
			opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
		// https://github.com/owncloud/ocis-proxy/issues/38
		accounts := acc.NewAccountsService("com.owncloud.accounts", mclient.DefaultClient)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := opt.Logger
			claims, ok := r.Context().Value(ClaimsKey).(ocisoidc.StandardClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			var uuid string
			entry, err := svcCache.Get(AccountsKey, claims.Email)
			if err != nil {
				l.Debug().Msgf("No cache entry for %v", claims.Email)
				resp, err := accounts.Get(context.Background(), &acc.GetRequest{
					Email: claims.Email,
				})

				if err != nil {
					l.Error().Err(err).Str("email", claims.Email).Msgf("Error fetching from accounts-service")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				err = svcCache.Set(AccountsKey, claims.Email, resp.Payload.Account.Uuid)
				if err != nil {
					l.Err(err).Str("email", claims.Email).Msgf("Could not cache user")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				uuid = resp.Payload.Account.Uuid
			}

			uuid, ok = entry.V.(string)
			if !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			l.Debug().Interface("claims", claims).Interface("uuid", uuid).Msgf("Associated claims with uuid")
			token, err := tokenManager.MintToken(r.Context(), &revauser.User{
				Id: &revauser.UserId{
					OpaqueId: uuid,
				},
				Username: claims.Email,
			})

			if err != nil {
				l.Error().Err(err).Msgf("Could not mint token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("x-access-token", token)
			next.ServeHTTP(w, r)
		})
	}
}

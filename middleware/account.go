package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/owncloud/ocis-pkg/v2/account"
)

// newAccountOptions initializes the available default options.
func newAccountOptions(opts ...account.Option) account.Options {
	opt := account.Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// AccountID serves as key for the account uuid in the context
const AccountID string = "Account-Id"

// ExtractAccountUUID provides a middleware to extract the account uuid from the x-access-token header value
// and write it to the context. If there is no x-access-token the middleware is omitted.
func ExtractAccountUUID(opts ...account.Option) func(http.Handler) http.Handler {
	opt := newAccountOptions(opts...)
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret":  opt.JWTSecret,
		"expires": int64(60),
	})
	if err != nil {
		opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-access-token")
			if len(token) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			user, err := tokenManager.DismantleToken(r.Context(), token)
			if err != nil {
				opt.Logger.Error().Err(err)
				return
			}

			// Important: user.Id.OpaqueId is the AccountUUID. Set this way in the account uuid middleware in ocis-proxy.
			// https://github.com/owncloud/ocis-proxy/blob/ea254d6036592cf9469d757d1295e0c4309d1e63/pkg/middleware/account_uuid.go#L109
			ctx := metadata.Set(r.Context(), AccountID, user.Id.OpaqueId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package middleware

import (
	"net/http"

	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/cs3org/reva/pkg/user"
)

// AccessToken middleware is used to set the user from an x-access-token to the context
func AccessToken(opts ...Option) func(next http.Handler) http.Handler {
	opt := newOptions(opts...)

	return func(next http.Handler) http.Handler {
		// TODO: handle error
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  opt.TokenManagerConfig.JWTSecret,
			"expires": int64(60),
		})
		if err != nil {
			opt.Logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("x-access-token")
			if token != "" {
				u, err := tokenManager.DismantleToken(r.Context(), token)
				if err != nil {
					opt.Logger.Error().Err(err).Msg("could not dismantle token")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// store user in context for request
				r = r.WithContext(user.ContextSetUser(r.Context(), u))
			}

			next.ServeHTTP(w, r)
		})
	}
}

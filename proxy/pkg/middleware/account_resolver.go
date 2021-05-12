package middleware

import (
	"github.com/cs3org/reva/pkg/auth/scope"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
	"net/http"

	tokenPkg "github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	revauser "github.com/cs3org/reva/pkg/user"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
)

// AccountResolver provides a middleware which mints a jwt and adds it to the proxied request based
// on the oidc-claims
func AccountResolver(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret":  options.TokenManagerConfig.JWTSecret,
			"expires": int64(60),
		})
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return &accountResolver{
			next:                  next,
			logger:                logger,
			tokenManager:          tokenManager,
			userProvider:          options.UserProvider,
			autoProvisionAccounts: options.AutoprovisionAccounts,
		}
	}
}

type accountResolver struct {
	next                  http.Handler
	logger                log.Logger
	tokenManager          tokenPkg.Manager
	userProvider          backend.UserBackend
	autoProvisionAccounts bool
}

func (m accountResolver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	claims := oidc.FromContext(req.Context())
	u, ok := revauser.ContextGetUser(req.Context())

	if claims == nil && !ok {
		m.next.ServeHTTP(w, req)
		return
	}

	if u == nil && claims != nil {
		var claim, value string
		switch {
		case claims.Email != "":
			claim, value = "mail", claims.Email
		case claims.PreferredUsername != "":
			claim, value = "username", claims.PreferredUsername
		case claims.OcisID != "":
			//claim, value = "id", claims.OcisID
		default:
			// TODO allow lookup by custom claim, eg an id ... or sub
			m.logger.Error().Msg("Could not lookup account, no mail or preferred_username claim set")
			w.WriteHeader(http.StatusInternalServerError)
		}

		var err error
		u, err = m.userProvider.GetUserByClaims(req.Context(), claim, value, true)

		if m.autoProvisionAccounts && err == backend.ErrAccountNotFound {
			m.logger.Debug().Interface("claims", claims).Interface("user", u).Msgf("User by claim not found... autoprovisioning.")
			u, err = m.userProvider.CreateUserFromClaims(req.Context(), claims)
		}

		if err == backend.ErrAccountNotFound || err == backend.ErrAccountDisabled {
			m.logger.Debug().Interface("claims", claims).Interface("user", u).Msgf("Unautorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			m.logger.Error().Err(err).Msg("Could not get user by claim")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.logger.Debug().Interface("claims", claims).Interface("user", u).Msgf("associated claims with uuid")
	}

	s, err := scope.GetOwnerScope()
	if err != nil {
		m.logger.Error().Err(err).Msgf("could not get owner scope")
		return
	}
	token, err := m.tokenManager.MintToken(req.Context(), u, s)
	if err != nil {
		m.logger.Error().Err(err).Msgf("could not mint token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set(tokenPkg.TokenHeader, token)

	m.next.ServeHTTP(w, req)
}

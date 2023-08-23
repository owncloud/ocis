package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
)

// AccountResolver provides a middleware which mints a jwt and adds it to the proxied request based
// on the oidc-claims
func AccountResolver(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		return &accountResolver{
			next:                  next,
			logger:                logger,
			userProvider:          options.UserProvider,
			userOIDCClaim:         options.UserOIDCClaim,
			userCS3Claim:          options.UserCS3Claim,
			userRoleAssigner:      options.UserRoleAssigner,
			autoProvisionAccounts: options.AutoprovisionAccounts,
		}
	}
}

type accountResolver struct {
	next                  http.Handler
	logger                log.Logger
	userProvider          backend.UserBackend
	userRoleAssigner      userroles.UserRoleAssigner
	autoProvisionAccounts bool
	userOIDCClaim         string
	userCS3Claim          string
}

// from https://codereview.stackexchange.com/a/280193
func splitWithEscaping(s string, separator string, escapeString string) []string {
	a := strings.Split(s, separator)

	for i := len(a) - 2; i >= 0; i-- {
		if strings.HasSuffix(a[i], escapeString) {
			a[i] = a[i][:len(a[i])-len(escapeString)] + separator + a[i+1]
			a = append(a[:i+1], a[i+2:]...)
		}
	}
	return a
}

func readUserIDClaim(path string, claims map[string]interface{}) (string, error) {
	// happy path
	value, _ := claims[path].(string)
	if value != "" {
		return value, nil
	}

	// try splitting path at .
	segments := splitWithEscaping(path, ".", "\\")
	subclaims := claims
	lastSegment := len(segments) - 1
	for i := range segments {
		if i < lastSegment {
			if castedClaims, ok := subclaims[segments[i]].(map[string]interface{}); ok {
				subclaims = castedClaims
			} else if castedClaims, ok := subclaims[segments[i]].(map[interface{}]interface{}); ok {
				subclaims = make(map[string]interface{}, len(castedClaims))
				for k, v := range castedClaims {
					if s, ok := k.(string); ok {
						subclaims[s] = v
					} else {
						return "", fmt.Errorf("could not walk claims path, key '%v' is not a string", k)
					}
				}
			}
		} else {
			if value, _ = subclaims[segments[i]].(string); value != "" {
				return value, nil
			}
		}
	}

	return value, fmt.Errorf("claim path '%s' not set or empty", path)
}

// TODO do not use the context to store values: https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
func (m accountResolver) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	claims := oidc.FromContext(ctx)
	user, ok := revactx.ContextGetUser(ctx)
	token := ""
	// TODO what if an X-Access-Token is set? happens eg for download requests to the /data endpoint in the reva frontend

	if claims == nil && !ok {
		m.next.ServeHTTP(w, req)
		return
	}

	if user == nil && claims != nil {
		value, err := readUserIDClaim(m.userOIDCClaim, claims)
		if err != nil {
			m.logger.Error().Err(err).Msg("could not read user id claim")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		user, token, err = m.userProvider.GetUserByClaims(req.Context(), m.userCS3Claim, value)

		if errors.Is(err, backend.ErrAccountNotFound) {
			m.logger.Debug().Str("claim", m.userOIDCClaim).Str("value", value).Msg("User by claim not found")
			if !m.autoProvisionAccounts {
				m.logger.Debug().Interface("claims", claims).Msg("Autoprovisioning disabled")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			m.logger.Debug().Interface("claims", claims).Msg("Autoprovisioning user")
			user, err = m.userProvider.CreateUserFromClaims(req.Context(), claims)
			if err != nil {
				m.logger.Error().Err(err).Msg("Autoprovisioning user failed")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			user, token, err = m.userProvider.GetUserByClaims(req.Context(), "userid", user.Id.OpaqueId)
			if err != nil {
				m.logger.Error().Err(err).Str("userid", user.Id.OpaqueId).Msg("Error getting token for autoprovisioned user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		if errors.Is(err, backend.ErrAccountDisabled) {
			m.logger.Debug().Interface("claims", claims).Msg("Disabled")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			m.logger.Error().Err(err).Msg("Could not get user by claim")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// resolve the user's roles
		user, err = m.userRoleAssigner.UpdateUserRoleAssignment(ctx, user, claims)
		if err != nil {
			m.logger.Error().Err(err).Msg("Could not get user roles")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// add user to context for selectors
		ctx = revactx.ContextSetUser(ctx, user)
		req = req.WithContext(ctx)

		m.logger.Debug().Interface("claims", claims).Interface("user", user).Msg("associated claims with user")
	} else if user != nil {
		var err error
		_, token, err = m.userProvider.GetUserByClaims(req.Context(), "username", user.Username)

		if errors.Is(err, backend.ErrAccountDisabled) {
			m.logger.Debug().Interface("user", user).Msg("Disabled")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			m.logger.Error().Err(err).Msg("Could not get user by claim")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	req.Header.Set(revactx.TokenHeader, token)

	m.next.ServeHTTP(w, req)
}

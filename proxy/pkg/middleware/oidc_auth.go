package middleware

import (
	"context"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	gOidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*gOidc.UserInfo, error)
}

func OIDCAuth(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)

	return func(next http.Handler) http.Handler {
		return &oidcAuth{
			next:         next,
			logger:       options.Logger,
			providerFunc: options.OIDCProviderFunc,
			httpClient:   options.HTTPClient,
			oidcIss: options.OIDCIss,
		}
	}
}

type oidcAuth struct {
	next         http.Handler
	logger       log.Logger
	provider     OIDCProvider
	providerFunc func() (OIDCProvider, error)
	httpClient   *http.Client
	oidcIss        string
}

func (m oidcAuth) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	if m.provider == nil {
		// Lazily initialize a provider

		// provider needs to be cached as when it is created
		// it will fetch the keys from the issuer using the .well-known
		// endpoint
		provider, err := m.providerFunc()
		if err != nil {
			m.logger.Error().Err(err).Msg("could not initialize oidcAuth provider")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		m.provider = provider
	}

	token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

	// TODO cache userinfo for access token if we can determine the expiry (which works in case it is a jwt based access token)
	oauth2Token := &oauth2.Token{
		AccessToken: token,
	}

	userInfo, err := m.provider.UserInfo(
		context.WithValue(req.Context(), oauth2.HTTPClient, m.httpClient),
		oauth2.StaticTokenSource(oauth2Token),
	)
	if err != nil {
		m.logger.Error().Err(err).Str("token", token).Msg("Failed to get userinfo")
		http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
		return
	}

	var claims oidc.StandardClaims
	if err := userInfo.Claims(&claims); err != nil {
		m.logger.Error().Err(err).Interface("userinfo", userInfo).Msg("failed to unmarshal userinfo claims")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: This should be read from the token instead of config
	claims.Iss = m.oidcIss

	// inject claims to the request context for the account_uuid middleware.
	req = req.WithContext(oidc.NewContext(req.Context(), &claims))

	m.logger.Debug().Interface("claims", claims).Interface("userInfo", userInfo).Msg("unmarshalled userinfo")

	// store claims in context
	// uses the original context, not the one with probably reduced security
	m.next.ServeHTTP(w, req.WithContext(oidc.NewContext(req.Context(), &claims)))
}

func (m oidcAuth) shouldServe(req *http.Request) bool {
	header := req.Header.Get("Authorization")

	if m.oidcIss == "" {
		return false
	}

	// todo: looks dirty, check later
	for _, ignoringPath := range []string{"/konnect/v1/userinfo"} {
		if req.URL.Path == ignoringPath {
			return false
		}
	}

	if !strings.HasPrefix(header, "Bearer ") {
		return false
	}

	return true
}

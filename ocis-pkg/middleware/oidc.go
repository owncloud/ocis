package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// newOidcOptions initializes the available default options.
func newOidcOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// OidcAuth provides a middleware to authenticate a bearer auth with an OpenID Connect identity provider
// It will put all claims provided by the userinfo endpoint in the context
func OidcAuth(opts ...Option) func(http.Handler) http.Handler {
	opt := newOidcOptions(opts...)

	// TODO use a micro store cache option

	var JWKS *keyfunc.JWKS
	getKeyfuncOnce := sync.Once{}
	issuer := "https://cloud.ocis.test"

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")
			switch {
			case strings.HasPrefix(authHeader, "Bearer "):
				getKeyfuncOnce.Do(func() {
					JWKS = getKeyfunc(opt.Logger, issuer, &http.Client{}, config.JWKS{
						RefreshInterval:   60, // minutes
						RefreshRateLimit:  60, // seconds
						RefreshTimeout:    10, // seconds
						RefreshUnknownKID: true,
					})
				})
				if JWKS == nil {
					return
				}

				jwtToken, err := jwt.Parse(strings.TrimPrefix(authHeader, "Bearer "), JWKS.Keyfunc)
				if err != nil {
					opt.Logger.Info().Err(err).Msg("Failed to parse/verify the access token.")
					return
				}
				opt.Logger.Debug().Interface("access token", &jwtToken).Msg("parsed access token")

				if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
					ctx = oidc.NewContext(ctx, claims)
				}

			default:
				// do nothing
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type jwksJSON struct {
	JWKSURL string `json:"jwks_uri"`
}

func getKeyfunc(log log.Logger, issuer string, client *http.Client, JwksOptions config.JWKS) *keyfunc.JWKS {
	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"

	resp, err := client.Get(wellKnown)
	if err != nil {
		log.Error().Err(err).Msg("Failed to set request for .well-known/openid-configuration")
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to read discovery response body")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("status", resp.Status).Str("body", string(body)).Msg("error requesting openid-configuration")
		return nil
	}

	var j jwksJSON
	err = json.Unmarshal(body, &j)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode provider openid-configuration")
		return nil
	}
	log.Debug().Str("jwks", j.JWKSURL).Msg("discovered jwks endpoint")
	options := keyfunc.Options{
		Client: client,
		RefreshErrorHandler: func(err error) {
			log.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
		},
		RefreshInterval:   time.Minute * time.Duration(JwksOptions.RefreshInterval),
		RefreshRateLimit:  time.Second * time.Duration(JwksOptions.RefreshRateLimit),
		RefreshTimeout:    time.Second * time.Duration(JwksOptions.RefreshTimeout),
		RefreshUnknownKID: JwksOptions.RefreshUnknownKID,
	}
	JWKS, err := keyfunc.Get(j.JWKSURL, options)
	if err != nil {
		JWKS = nil
		log.Error().Err(err).Msg("Failed to create JWKS from resource at the given URL.")
		return nil
	}
	return JWKS
}

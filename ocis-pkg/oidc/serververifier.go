package oidc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

// AccessTokenVerifier verifies bearer access tokens presented to the resource
// server. This is server-side verification of a token minted by the IDP, not
// an OIDC client operation.
type AccessTokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (RegClaimsWithSID, jwt.MapClaims, error)
}

type accessTokenVerifier struct {
	// Logger to use for logging, must be set
	Logger log.Logger

	issuer                     string
	provider                   *ProviderMetadata
	providerLock               *sync.Mutex
	skipIssuerValidation       bool
	accessTokenVerifyMethod    string
	accessTokenVerifyAudiences []string

	JWKSOptions config.JWKS
	JWKS        *keyfunc.JWKS
	jwksLock    *sync.Mutex

	httpClient *http.Client
}

// NewAccessTokenVerifier returns an AccessTokenVerifier for the given issuer.
func NewAccessTokenVerifier(opts ...Option) AccessTokenVerifier {
	options := newOptions(opts...)

	return &accessTokenVerifier{
		Logger:                     options.Logger,
		issuer:                     options.OIDCIssuer,
		httpClient:                 options.HTTPClient,
		accessTokenVerifyMethod:    options.AccessTokenVerifyMethod,
		accessTokenVerifyAudiences: options.AccessTokenVerifyAudiences,
		JWKSOptions:                options.JWKSOptions,
		JWKS:                       options.JWKS,
		providerLock:               &sync.Mutex{},
		jwksLock:                   &sync.Mutex{},
		provider:                   options.ProviderMetadata,
	}
}

func (v *accessTokenVerifier) lookupWellKnownOpenidConfiguration(ctx context.Context) error {
	v.providerLock.Lock()
	defer v.providerLock.Unlock()
	if v.provider == nil {
		wellKnown := strings.TrimSuffix(v.issuer, "/") + wellknownPath
		req, err := tracing.GetNewRequest(ctx, http.MethodGet, wellKnown, nil)
		if err != nil {
			return err
		}
		resp, err := v.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s: %s", resp.Status, body)
		}

		var p ProviderMetadata
		err = unmarshalResp(resp, body, &p)
		if err != nil {
			return fmt.Errorf("oidc: failed to decode provider discovery object: %v", err)
		}

		if !v.skipIssuerValidation && p.Issuer != v.issuer {
			return fmt.Errorf("oidc: issuer did not match the issuer returned by provider, expected %q got %q", v.issuer, p.Issuer)
		}
		v.provider = &p
	}
	return nil
}

func (v *accessTokenVerifier) getKeyfunc() *keyfunc.JWKS {
	v.jwksLock.Lock()
	defer v.jwksLock.Unlock()
	if v.JWKS == nil {
		var err error
		v.Logger.Debug().Str("jwks", v.provider.JwksURI).Msg("discovered jwks endpoint")
		options := keyfunc.Options{
			Client: v.httpClient,
			RefreshErrorHandler: func(err error) {
				v.Logger.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
			},
			RefreshInterval:   time.Minute * time.Duration(v.JWKSOptions.RefreshInterval),
			RefreshRateLimit:  time.Second * time.Duration(v.JWKSOptions.RefreshRateLimit),
			RefreshTimeout:    time.Second * time.Duration(v.JWKSOptions.RefreshTimeout),
			RefreshUnknownKID: v.JWKSOptions.RefreshUnknownKID,
		}
		v.JWKS, err = keyfunc.Get(v.provider.JwksURI, options)
		if err != nil {
			v.JWKS = nil
			v.Logger.Error().Err(err).Msg("Failed to create JWKS from resource at the given URL.")
			return nil
		}
	}
	return v.JWKS
}

// VerifyAccessToken verifies the given bearer access token according to the
// configured verification method.
func (v *accessTokenVerifier) VerifyAccessToken(ctx context.Context, token string) (RegClaimsWithSID, jwt.MapClaims, error) {
	if err := v.lookupWellKnownOpenidConfiguration(ctx); err != nil {
		return RegClaimsWithSID{}, jwt.MapClaims{}, err
	}
	switch v.accessTokenVerifyMethod {
	case config.AccessTokenVerificationJWT:
		return v.verifyAccessTokenJWT(token)
	case config.AccessTokenVerificationNone:
		v.Logger.Debug().Msg("Access Token verification disabled")
		return RegClaimsWithSID{}, jwt.MapClaims{}, nil
	default:
		v.Logger.Error().Str("access_token_verify_method", v.accessTokenVerifyMethod).Msg("Unknown Access Token verification setting")
		return RegClaimsWithSID{}, jwt.MapClaims{}, errors.New("unknown Access Token Verification method")
	}
}

// verifyAccessTokenJWT tries to parse and verify the access token as a JWT.
func (v *accessTokenVerifier) verifyAccessTokenJWT(token string) (RegClaimsWithSID, jwt.MapClaims, error) {
	var claims RegClaimsWithSID
	mapClaims := jwt.MapClaims{}
	jwks := v.getKeyfunc()
	if jwks == nil {
		return claims, mapClaims, errors.New("error initializing jwks keyfunc")
	}

	issuer := v.issuer
	if v.provider.AccessTokenIssuer != "" {
		// AD FS .well-known/openid-configuration has an optional `access_token_issuer` which takes precedence over `issuer`
		// See https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-oidce/586de7dd-3385-47c7-93a2-935d9e90441c
		issuer = v.provider.AccessTokenIssuer
	}

	_, err := jwt.ParseWithClaims(token, &claims, jwks.Keyfunc, jwt.WithIssuer(issuer))
	if err != nil {
		return claims, mapClaims, err
	}
	_, _, err = new(jwt.Parser).ParseUnverified(token, mapClaims)
	// TODO: decode mapClaims to sth readable
	v.Logger.Debug().Interface("access token", &claims).Msg("parsed access token")
	if err != nil {
		v.Logger.Info().Err(err).Msg("Failed to parse/verify the access token.")
		return claims, mapClaims, err
	}

	// Verify the token was issued for this instance. Keycloak puts a generic
	// value (e.g. "account") in "aud" and the real client id in "azp", so we
	// accept the token when the allowlist matches either claim.
	if err := v.verifyAudience(claims.Audience, azpFromClaims(mapClaims)); err != nil {
		v.Logger.Info().Err(err).Msg("Access token rejected: audience not allowed.")
		return claims, mapClaims, err
	}

	return claims, mapClaims, nil
}

// verifyAudience checks that the token's "aud" or "azp" claim matches the
// configured allowlist. An empty allowlist disables the check.
func (v *accessTokenVerifier) verifyAudience(audiences jwt.ClaimStrings, azp string) error {
	if len(v.accessTokenVerifyAudiences) == 0 {
		return nil
	}
	for _, allowed := range v.accessTokenVerifyAudiences {
		allowed = strings.TrimSpace(allowed)
		if allowed == "" {
			continue
		}
		if allowed == azp || slices.Contains(audiences, allowed) {
			return nil
		}
	}
	return errors.New("oidc: access token audience is not allowed")
}

func azpFromClaims(claims jwt.MapClaims) string {
	if azp, ok := claims[Azp].(string); ok {
		return azp
	}
	return ""
}

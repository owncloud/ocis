package oidc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	gOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2"
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*UserInfo, error)
	VerifyAccessToken(ctx context.Context, token string) (jwt.RegisteredClaims, []string, error)
	VerifyLogoutToken(ctx context.Context, token string) (*LogoutToken, error)
}

// KeySet is a set of publc JSON Web Keys that can be used to validate the signature
// of JSON web tokens. This is expected to be backed by a remote key set through
// provider metadata discovery or an in-memory set of keys delivered out-of-band.
type KeySet interface {
	// VerifySignature parses the JSON web token, verifies the signature, and returns
	// the raw payload. Header and claim fields are validated by other parts of the
	// package. For example, the KeySet does not need to check values such as signature
	// algorithm, issuer, and audience since the IDTokenVerifier validates these values
	// independently.
	//
	// If VerifySignature makes HTTP requests to verify the token, it's expected to
	// use any HTTP client associated with the context through ClientContext.
	VerifySignature(ctx context.Context, jwt string) (payload []byte, err error)
}

type oidcClient struct {
	// Logger to use for logging, must be set
	Logger log.Logger

	clientID                string
	skipClientIDCheck       bool
	issuer                  string
	provider                *ProviderMetadata
	providerLock            *sync.Mutex
	skipIssuerValidation    bool
	accessTokenVerifyMethod string
	remoteKeySet            KeySet // TODO replace usage with keyfunc?
	algorithms              []string

	JWKSOptions config.JWKS
	JWKS        *keyfunc.JWKS
	jwksLock    *sync.Mutex

	httpClient *http.Client
}

// supportedAlgorithms is a list of algorithms explicitly supported by this
// package. If a provider supports other algorithms, such as HS256 or none,
// those values won't be passed to the IDTokenVerifier.
var supportedAlgorithms = map[string]bool{
	RS256: true,
	RS384: true,
	RS512: true,
	ES256: true,
	ES384: true,
	ES512: true,
	PS256: true,
	PS384: true,
	PS512: true,
}

// NewOIDCClient returns an OIDC client for the given issuer
func NewOIDCClient(opts ...Option) OIDCProvider {
	options := newOptions(opts...)

	return &oidcClient{
		Logger:                  options.Logger,
		issuer:                  options.OidcIssuer,
		httpClient:              options.HTTPClient,
		accessTokenVerifyMethod: options.AccessTokenVerifyMethod,
		JWKSOptions:             options.JWKSOptions, // TODO I don't like that we pass down config options ...
		providerLock:            &sync.Mutex{},
		jwksLock:                &sync.Mutex{},
		clientID:                options.ClientID,
		skipClientIDCheck:       options.SkipClientIDCheck,
	}
}

func (c *oidcClient) lookupWellKnownOpenidConfiguration(ctx context.Context) error {
	c.providerLock.Lock()
	defer c.providerLock.Unlock()
	if c.provider == nil {
		wellKnown := strings.TrimSuffix(c.issuer, "/") + wellknownPath
		req, err := http.NewRequest("GET", wellKnown, nil)
		if err != nil {
			return err
		}
		resp, err := c.httpClient.Do(req.WithContext(ctx))
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

		if !c.skipIssuerValidation && p.Issuer != c.issuer {
			return fmt.Errorf("oidc: issuer did not match the issuer returned by provider, expected %q got %q", c.issuer, p.Issuer)
		}
		var algs []string
		for _, a := range p.IDTokenSigningAlgValuesSupported {
			if supportedAlgorithms[a] {
				algs = append(algs, a)
			}
		}
		c.provider = &p
		c.algorithms = algs
		c.remoteKeySet = gOidc.NewRemoteKeySet(ctx, p.JwksURI)
	}
	return nil
}

func (c *oidcClient) getKeyfunc() *keyfunc.JWKS {
	c.jwksLock.Lock()
	defer c.jwksLock.Unlock()
	if c.JWKS == nil {
		var err error
		c.Logger.Debug().Str("jwks", c.provider.JwksURI).Msg("discovered jwks endpoint")
		options := keyfunc.Options{
			Client: c.httpClient,
			RefreshErrorHandler: func(err error) {
				c.Logger.Error().Err(err).Msg("There was an error with the jwt.Keyfunc")
			},
			RefreshInterval:   time.Minute * time.Duration(c.JWKSOptions.RefreshInterval),
			RefreshRateLimit:  time.Second * time.Duration(c.JWKSOptions.RefreshRateLimit),
			RefreshTimeout:    time.Second * time.Duration(c.JWKSOptions.RefreshTimeout),
			RefreshUnknownKID: c.JWKSOptions.RefreshUnknownKID,
		}
		c.JWKS, err = keyfunc.Get(c.provider.JwksURI, options)
		if err != nil {
			c.JWKS = nil
			c.Logger.Error().Err(err).Msg("Failed to create JWKS from resource at the given URL.")
			return nil
		}
	}
	return c.JWKS
}

type stringAsBool bool

func (sb *stringAsBool) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "true", `"true"`:
		*sb = true
	case "false", `"false"`:
		*sb = false
	default:
		return errors.New("invalid value for boolean")
	}
	return nil
}

// UserInfo represents the OpenID Connect userinfo claims.
type UserInfo struct {
	Subject       string `json:"sub"`
	Profile       string `json:"profile"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`

	claims []byte
}

type userInfoRaw struct {
	Subject string `json:"sub"`
	Profile string `json:"profile"`
	Email   string `json:"email"`
	// Handle providers that return email_verified as a string
	// https://forums.aws.amazon.com/thread.jspa?messageID=949441&#949441 and
	// https://discuss.elastic.co/t/openid-error-after-authenticating-against-aws-cognito/206018/11
	EmailVerified stringAsBool `json:"email_verified"`
}

// Claims unmarshals the raw JSON object claims into the provided object.
func (u *UserInfo) Claims(v interface{}) error {
	if u.claims == nil {
		return errors.New("oidc: claims not set")
	}
	return json.Unmarshal(u.claims, v)
}

func (c *oidcClient) UserInfo(ctx context.Context, tokenSource oauth2.TokenSource) (*UserInfo, error) {
	if err := c.lookupWellKnownOpenidConfiguration(ctx); err != nil {
		return nil, err
	}

	if c.provider.UserinfoEndpoint == "" {
		if c.provider.UserinfoEndpoint == "" {
			return nil, errors.New("oidc: user info endpoint is not supported by this provider")
		}
	}

	req, err := http.NewRequest("GET", c.provider.UserinfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("oidc: create GET request: %v", err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("oidc: get access token: %v", err)
	}
	token.SetAuthHeader(req)

	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	ct := resp.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(ct)
	if parseErr == nil && mediaType == "application/jwt" {
		payload, err := c.remoteKeySet.VerifySignature(ctx, string(body))
		if err != nil {
			return nil, fmt.Errorf("oidc: invalid userinfo jwt signature %v", err)
		}
		body = payload
	}

	var userInfo userInfoRaw
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("oidc: failed to decode userinfo: %v", err)
	}
	return &UserInfo{
		Subject:       userInfo.Subject,
		Profile:       userInfo.Profile,
		Email:         userInfo.Email,
		EmailVerified: bool(userInfo.EmailVerified),
		claims:        body,
	}, nil
}

func (c *oidcClient) VerifyAccessToken(ctx context.Context, token string) (jwt.RegisteredClaims, []string, error) {
	var mapClaims []string
	if err := c.lookupWellKnownOpenidConfiguration(ctx); err != nil {
		return jwt.RegisteredClaims{}, mapClaims, err
	}
	switch c.accessTokenVerifyMethod {
	case config.AccessTokenVerificationJWT:
		return c.verifyAccessTokenJWT(token)
	case config.AccessTokenVerificationNone:
		c.Logger.Debug().Msg("Access Token verification disabled")
		return jwt.RegisteredClaims{}, mapClaims, nil
	default:
		c.Logger.Error().Str("access_token_verify_method", c.accessTokenVerifyMethod).Msg("Unknown Access Token verification setting")
		return jwt.RegisteredClaims{}, mapClaims, errors.New("unknown Access Token Verification method")
	}
}

// verifyAccessTokenJWT tries to parse and verify the access token as a JWT.
func (c *oidcClient) verifyAccessTokenJWT(token string) (jwt.RegisteredClaims, []string, error) {
	var claims jwt.RegisteredClaims
	var mapClaims []string
	jwks := c.getKeyfunc()
	if jwks == nil {
		return claims, mapClaims, errors.New("error initializing jwks keyfunc")
	}

	_, err := jwt.ParseWithClaims(token, &claims, jwks.Keyfunc)
	if err != nil {
		return claims, mapClaims, err
	}
	_, mapClaims, err = new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	// TODO: decode mapClaims to sth readable
	c.Logger.Debug().Interface("access token", &claims).Msg("parsed access token")
	if err != nil {
		c.Logger.Info().Err(err).Msg("Failed to parse/verify the access token.")
		return claims, mapClaims, err
	}

	if !claims.VerifyIssuer(c.issuer, true) {
		vErr := jwt.ValidationError{}
		vErr.Inner = jwt.ErrTokenInvalidIssuer
		vErr.Errors |= jwt.ValidationErrorIssuer
		return claims, mapClaims, vErr
	}

	return claims, mapClaims, nil
}

func (c *oidcClient) VerifyLogoutToken(ctx context.Context, rawToken string) (*LogoutToken, error) {
	jws, err := jose.ParseSigned(rawToken)
	if err != nil {
		return nil, err
	}
	// Throw out tokens with invalid claims before trying to verify the token. This lets
	// us do cheap checks before possibly re-syncing keys.
	payload, err := parseJWT(rawToken)
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt: %v", err)
	}
	var token LogoutToken
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, fmt.Errorf("oidc: failed to unmarshal claims: %v", err)
	}

	//4. Verify that the Logout Token contains a sub Claim, a sid Claim, or both.
	if token.Subject == "" && token.SessionId == "" {
		return nil, fmt.Errorf("oidc: logout token must contain either sub or sid and MAY contain both")
	}
	//5. Verify that the Logout Token contains an events Claim whose value is JSON object containing the member name http://schemas.openid.net/event/backchannel-logout.
	if token.Events.Event == nil {
		return nil, fmt.Errorf("oidc: logout token must contain logout event")
	}
	//6. Verify that the Logout Token does not contain a nonce Claim.
	type nonce struct {
		Nonce *string `json:"nonce"`
	}
	var n nonce
	json.Unmarshal(payload, &n)
	if n.Nonce != nil {
		return nil, fmt.Errorf("oidc: nonce on logout token MUST NOT be present")
	}
	// Check issuer.
	if !c.skipIssuerValidation && token.Issuer != c.issuer {
		return nil, fmt.Errorf("oidc: id token issued by a different provider, expected %q got %q", c.issuer, token.Issuer)
	}

	// If a client ID has been provided, make sure it's part of the audience. SkipClientIDCheck must be true if ClientID is empty.
	//
	// This check DOES NOT ensure that the ClientID is the party to which the ID Token was issued (i.e. Authorized party).
	if !c.skipClientIDCheck {
		if c.clientID != "" {
			if !contains(token.Audience, c.clientID) {
				return nil, fmt.Errorf("oidc: expected audience %q got %q", c.clientID, token.Audience)
			}
		} else {
			return nil, fmt.Errorf("oidc: invalid configuration, clientID must be provided or SkipClientIDCheck must be set")
		}
	}

	switch len(jws.Signatures) {
	case 0:
		return nil, fmt.Errorf("oidc: id token not signed")
	case 1:
	default:
		return nil, fmt.Errorf("oidc: multiple signatures on id token not supported")
	}

	sig := jws.Signatures[0]
	supportedSigAlgs := c.algorithms
	if len(supportedSigAlgs) == 0 {
		supportedSigAlgs = []string{RS256}
	}

	if !contains(supportedSigAlgs, sig.Header.Algorithm) {
		return nil, fmt.Errorf("oidc: id token signed with unsupported algorithm, expected %q got %q", supportedSigAlgs, sig.Header.Algorithm)
	}

	gotPayload, err := c.remoteKeySet.VerifySignature(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %v", err)
	}

	// Ensure that the payload returned by the square actually matches the payload parsed earlier.
	if !bytes.Equal(gotPayload, payload) {
		return nil, errors.New("oidc: internal error, payload parsed did not match previous payload")
	}

	return &token, nil
}

func unmarshalResp(r *http.Response, body []byte, v interface{}) error {
	err := json.Unmarshal(body, &v)
	if err == nil {
		return nil
	}
	ct := r.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(ct)
	if parseErr == nil && mediaType == "application/json" {
		return fmt.Errorf("got Content-Type = application/json, but could not unmarshal as JSON: %v", err)
	}
	return fmt.Errorf("expected Content-Type = application/json, got %q: %v", ct, err)
}

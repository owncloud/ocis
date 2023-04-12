package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"sync"

	gOidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"golang.org/x/oauth2"
)

// OIDCProvider used to mock the oidc provider during tests
type OIDCProvider interface {
	UserInfo(ctx context.Context, ts oauth2.TokenSource) (*UserInfo, error)
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
	issuer               string
	provider             *ProviderMetadata
	providerLock         *sync.Mutex
	skipIssuerValidation bool
	remoteKeySet         KeySet
	algorithms           []string
	// Logger to use for logging, must be set
	Logger log.Logger

	client *http.Client
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
		Logger:       options.Logger,
		issuer:       options.OidcIssuer,
		client:       options.HTTPClient,
		providerLock: &sync.Mutex{},
	}
}

func (c *oidcClient) lookupWellKnownOpenidConfiguration(ctx context.Context) error {
	c.providerLock.Lock()
	defer c.providerLock.Unlock()
	if c.provider == nil {
		wellKnown := strings.TrimSuffix(c.issuer, "/") + "/.well-known/openid-configuration"
		req, err := http.NewRequest("GET", wellKnown, nil)
		if err != nil {
			return err
		}
		resp, err := c.client.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
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

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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

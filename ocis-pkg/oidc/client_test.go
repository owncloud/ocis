package oidc_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"golang.org/x/oauth2"
)

type signingKey struct {
	priv interface{}
	jwks *keyfunc.JWKS
}

func TestLogoutVerify(t *testing.T) {
	tests := []logoutVerificationTest{
		{
			name: "good token",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo",
				"sub": "248289761001",
				"aud": "s6BhdRkqt3",
				"iat": 1471566154,
				"jti": "bWJq",
				"sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
				"events": map[string]interface{}{
					"http://schemas.openid.net/event/backchannel-logout": struct{}{},
				},
			}),
			signKey: newRSAKey(t),
		},
		{
			name:   "invalid issuer",
			issuer: "https://bar",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo1",
				"sub": "248289761001",
				"events": map[string]interface{}{
					"http://schemas.openid.net/event/backchannel-logout": struct{}{},
				},
			}),
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "invalid sig",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo",
				"sub": "248289761001",
				"aud": "s6BhdRkqt3",
				"iat": 1471566154,
				"jti": "bWJq",
				"sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
				"events": map[string]interface{}{
					"http://schemas.openid.net/event/backchannel-logout": struct{}{},
				},
			}),
			signKey:         newRSAKey(t),
			verificationKey: newRSAKey(t),
			wantErr:         true,
		},
		{
			name: "no sid and no sub",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo",
				"aud": "s6BhdRkqt3",
				"iat": 1471566154,
				"jti": "bWJq",
				"events": map[string]interface{}{
					"http://schemas.openid.net/event/backchannel-logout": struct{}{},
				},
			}),
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "Prohibited nonce present",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss":   "https://foo",
				"sub":   "248289761001",
				"aud":   "s6BhdRkqt3",
				"iat":   1471566154,
				"jti":   "bWJq",
				"sid":   "08a5019c-17e1-4977-8f42-65a12843ea02",
				"nonce": "123",
				"events": map[string]interface{}{
					"http://schemas.openid.net/event/backchannel-logout": struct{}{},
				},
			}),
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "Wrong Event string",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo",
				"sub": "248289761001",
				"aud": "s6BhdRkqt3",
				"iat": 1471566154,
				"jti": "bWJq",
				"sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
				"events": map[string]interface{}{
					"http://blah.blah.blash/event/backchannel-logout": struct{}{},
				},
			}),
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "No Event string",
			logoutToken: jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": "https://foo",
				"sub": "248289761001",
				"aud": "s6BhdRkqt3",
				"iat": 1471566154,
				"jti": "bWJq",
				"sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
			}),
			signKey: newRSAKey(t),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

type logoutVerificationTest struct {
	// Name of the subtest.
	name string

	// If not provided defaults to "https://foo"
	issuer string

	// JWT payload (just the claims).
	logoutToken *jwt.Token

	// Key to sign the ID Token with.
	signKey *signingKey
	// If not provided defaults to signKey. Only useful when
	// testing invalid signatures.
	verificationKey *signingKey

	wantErr bool
}

func (v logoutVerificationTest) runGetToken(t *testing.T) (*oidc.LogoutToken, error) {
	//	token := v.signKey.sign(t, []byte(v.logoutToken))
	v.logoutToken.Header["kid"] = "1"
	token, err := v.logoutToken.SignedString(v.signKey.priv)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	issuer := "https://foo"
	var jwks *keyfunc.JWKS
	if v.verificationKey == nil {
		jwks = v.signKey.jwks
	} else {
		jwks = v.verificationKey.jwks
	}

	pm := oidc.ProviderMetadata{}
	verifier := oidc.NewOIDCClient(
		oidc.WithOidcIssuer(issuer),
		oidc.WithJWKS(jwks),
		oidc.WithProviderMetadata(&pm),
	)

	return verifier.VerifyLogoutToken(ctx, token)
}

func (l logoutVerificationTest) run(t *testing.T) {
	_, err := l.runGetToken(t)
	if err != nil && !l.wantErr {
		t.Errorf("%v", err)
	}
	if err == nil && l.wantErr {
		t.Errorf("expected error")
	}
}

// A userinfo call that returns a transient status (429/5xx) must be classified as
// temporarily unavailable so the proxy answers 503, not 401 (#12999). A genuine
// auth status (401/403) must stay a plain error that maps to 401.
func TestUserInfoStatusClassification(t *testing.T) {
	tests := []struct {
		status    int
		transient bool
	}{
		{http.StatusUnauthorized, false},
		{http.StatusForbidden, false},
		{http.StatusInternalServerError, true},
		{http.StatusBadGateway, true},
		{http.StatusServiceUnavailable, true},
		{http.StatusTooManyRequests, true},
	}
	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			}))
			defer srv.Close()

			c := oidc.NewOIDCClient(
				oidc.WithProviderMetadata(&oidc.ProviderMetadata{UserinfoEndpoint: srv.URL}),
				oidc.WithHTTPClient(srv.Client()),
			)

			_, err := c.UserInfo(context.Background(),
				oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "t"}))
			if err == nil {
				t.Fatalf("expected an error for status %d", tt.status)
			}
			if got := errors.Is(err, oidc.ErrTemporarilyUnavailable); got != tt.transient {
				t.Fatalf("status %d: transient = %v, want %v (err=%v)", tt.status, got, tt.transient, err)
			}
		})
	}
}

// A transient failure while discovering the IdP (well-known/JWKS) during access-token
// verification must be classified as temporarily unavailable, not surfaced as a 401 (#12999).
func TestVerifyAccessTokenTransientDiscovery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := oidc.NewOIDCClient(
		oidc.WithOidcIssuer(srv.URL),
		oidc.WithHTTPClient(srv.Client()),
		oidc.WithAccessTokenVerifyMethod(config.AccessTokenVerificationJWT),
	)

	_, _, err := c.VerifyAccessToken(context.Background(), "any.token.here")
	if err == nil {
		t.Fatal("expected an error for a 503 during discovery")
	}
	if !errors.Is(err, oidc.ErrTemporarilyUnavailable) {
		t.Fatalf("discovery 503 must be transient, got %v", err)
	}
}

func newRSAKey(t testing.TB) *signingKey {
	priv, err := rsa.GenerateKey(rand.Reader, 1028)
	if err != nil {
		t.Fatal(err)
	}
	givenKey := keyfunc.NewGivenRSA(
		&priv.PublicKey,
		keyfunc.GivenKeyOptions{Algorithm: jwt.SigningMethodRS256.Alg()},
	)
	jwks := keyfunc.NewGiven(
		map[string]keyfunc.GivenKey{
			"1": givenKey,
		},
	)

	return &signingKey{priv, jwks}
}

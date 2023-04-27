package oidc_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"

	goidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"gopkg.in/square/go-jose.v2"
)

type signingKey struct {
	keyID string // optional
	priv  interface{}
	pub   interface{}
	alg   jose.SignatureAlgorithm
}

// sign creates a JWS using the private key from the provided payload.
func (s *signingKey) sign(t testing.TB, payload []byte) string {
	privKey := &jose.JSONWebKey{Key: s.priv, Algorithm: string(s.alg), KeyID: s.keyID}

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: s.alg, Key: privKey}, nil)
	if err != nil {
		t.Fatal(err)
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}

	data, err := jws.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func (s *signingKey) jwk() jose.JSONWebKey {
	return jose.JSONWebKey{Key: s.pub, Use: "sig", Algorithm: string(s.alg), KeyID: s.keyID}
}

func TestLogoutVerify(t *testing.T) {
	tests := []logoutVerificationTest{
		{
			name: "good token",
			logoutToken: ` {
							   "iss": "https://foo",
							   "sub": "248289761001",
							   "aud": "s6BhdRkqt3",
							   "iat": 1471566154,
							   "jti": "bWJq",
							   "sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
							   "events": {
								 "http://schemas.openid.net/event/backchannel-logout": {}
								 }
							  }`,
			signKey: newRSAKey(t),
		},
		{
			name:        "invalid issuer",
			issuer:      "https://bar",
			logoutToken: `{"iss":"https://foo"}`,
			config: goidc.Config{
				SkipExpiryCheck: true,
			},
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "invalid sig",
			logoutToken: `{
							   "iss": "https://foo",
							   "sub": "248289761001",
							   "aud": "s6BhdRkqt3",
							   "iat": 1471566154,
							   "jti": "bWJq",
							   "sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
							   "events": {
								 "http://schemas.openid.net/event/backchannel-logout": {}
								 }
							  }`,
			config: goidc.Config{
				SkipExpiryCheck: true,
			},
			signKey:         newRSAKey(t),
			verificationKey: newRSAKey(t),
			wantErr:         true,
		},
		{
			name: "no sid and no sub",
			logoutToken: ` {
								"iss": "https://foo",
							   "aud": "s6BhdRkqt3",
							   "iat": 1471566154,
							   "jti": "bWJq",
							   "events": {
								 "http://schemas.openid.net/event/backchannel-logout": {}
								 }
							  }`,
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "Prohibited nonce present",
			logoutToken: ` {
							   	"iss": "https://foo",
								"sub": "248289761001",
							   	"aud": "s6BhdRkqt3",
							   	"iat": 1471566154,
							   	"jti": "bWJq",
								"nonce" : "prohibited",
							   	"events": {
								 "http://schemas.openid.net/event/backchannel-logout": {}
								 }
							  }`,
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "Wrong Event string",
			logoutToken: ` {
							   "iss": "https://foo",
							   "sub": "248289761001",
							   "aud": "s6BhdRkqt3",
							   "iat": 1471566154,
							   "jti": "bWJq",
							   "sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
							   "events": {
								 "not a logout event": {}
								 }
							  }`,
			signKey: newRSAKey(t),
			wantErr: true,
		},
		{
			name: "No Event string",
			logoutToken: ` {
							   "iss": "https://foo",
							   "sub": "248289761001",
							   "aud": "s6BhdRkqt3",
							   "iat": 1471566154,
							   "jti": "bWJq",
							   "sid": "08a5019c-17e1-4977-8f42-65a12843ea02",
							  }`,
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
	logoutToken string

	// Key to sign the ID Token with.
	signKey *signingKey
	// If not provided defaults to signKey. Only useful when
	// testing invalid signatures.
	verificationKey *signingKey

	config  goidc.Config
	wantErr bool
}

type testVerifier struct {
	jwk jose.JSONWebKey
}

func (t *testVerifier) VerifySignature(ctx context.Context, jwt string) ([]byte, error) {
	jws, err := jose.ParseSigned(jwt)
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt: %v", err)
	}
	return jws.Verify(&t.jwk)
}

func (v logoutVerificationTest) runGetToken(t *testing.T) (*oidc.LogoutToken, error) {
	token := v.signKey.sign(t, []byte(v.logoutToken))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	issuer := "https://foo"
	var ks goidc.KeySet
	if v.verificationKey == nil {
		ks = &testVerifier{v.signKey.jwk()}
	} else {
		ks = &testVerifier{v.verificationKey.jwk()}
	}

	pm := oidc.ProviderMetadata{}
	verifier := oidc.NewOIDCClient(
		oidc.WithOidcIssuer(issuer),
		oidc.WithKeySet(ks),
		oidc.WithConfig(&v.config),
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

func newRSAKey(t testing.TB) *signingKey {
	priv, err := rsa.GenerateKey(rand.Reader, 1028)
	if err != nil {
		t.Fatal(err)
	}
	return &signingKey{"", priv, priv.Public(), jose.RS256}
}

func newECDSAKey(t *testing.T) *signingKey {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return &signingKey{"", priv, priv.Public(), jose.ES256}
}

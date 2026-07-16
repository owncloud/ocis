package oidc_test

import (
	"context"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
)

func TestVerifyAccessTokenAudience(t *testing.T) {
	const issuer = "https://foo"
	key := newRSAKey(t)

	tests := []struct {
		name string
		// allowlist configured on the verifier (empty disables the check)
		verifyAud []string
		// claims of the presented access token
		aud     interface{}
		azp     string
		wantErr bool
	}{
		{
			name:      "check disabled accepts foreign audience",
			verifyAud: nil,
			aud:       "account",
			azp:       "evil-app",
			wantErr:   false,
		},
		{
			name:      "rejects foreign azp with generic aud (keycloak style)",
			verifyAud: []string{"web"},
			aud:       "account",
			azp:       "evil-app",
			wantErr:   true,
		},
		{
			name:      "accepts matching azp with generic aud (keycloak style)",
			verifyAud: []string{"web"},
			aud:       "account",
			azp:       "web",
			wantErr:   false,
		},
		{
			name:      "accepts matching aud (builtin idp style)",
			verifyAud: []string{"web"},
			aud:       "web",
			azp:       "",
			wantErr:   false,
		},
		{
			name:      "accepts when one of multiple audiences matches",
			verifyAud: []string{"desktop", "web"},
			aud:       []string{"account", "web"},
			azp:       "evil-app",
			wantErr:   false,
		},
		{
			name:      "rejects when neither aud nor azp match",
			verifyAud: []string{"web", "desktop"},
			aud:       []string{"account", "other"},
			azp:       "evil-app",
			wantErr:   true,
		},
		{
			name:      "rejects when azp is empty and aud does not match",
			verifyAud: []string{"web", "desktop"},
			aud:       []string{"account", "other"},
			azp:       "",
			wantErr:   true,
		},
		{
			name:      "tolerates whitespace around configured audience",
			verifyAud: []string{" web ", "desktop"},
			aud:       "account",
			azp:       "web",
			wantErr:   false,
		},
		{
			name:      "blank configured audience does not match missing azp",
			verifyAud: []string{"  "},
			aud:       "account",
			azp:       "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"iss": issuer,
				"sub": "einstein",
				"aud": tc.aud,
			}
			if tc.azp != "" {
				claims["azp"] = tc.azp
			}
			tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
			tok.Header["kid"] = "1"
			signed, err := tok.SignedString(key.priv)
			if err != nil {
				t.Fatal(err)
			}

			verifier := oidc.NewAccessTokenVerifier(
				oidc.WithOidcIssuer(issuer),
				oidc.WithJWKS(key.jwks),
				oidc.WithProviderMetadata(&oidc.ProviderMetadata{}),
				oidc.WithAccessTokenVerifyMethod(config.AccessTokenVerificationJWT),
				oidc.WithAccessTokenVerifyAudiences(tc.verifyAud),
			)

			_, _, err = verifier.VerifyAccessToken(context.Background(), signed)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

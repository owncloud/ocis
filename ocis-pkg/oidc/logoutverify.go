package oidc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/square/go-jose.v2"

	gOidc "github.com/coreos/go-oidc/v3/oidc"
)

// This adds the ability to verify Logout Tokens as specified in https://openid.net/specs/openid-connect-backchannel-1_0.html

// LogoutTokenVerifier provides verification for Logout Tokens.
type LogoutTokenVerifier struct {
	keySet gOidc.KeySet
	config *gOidc.Config
	issuer string
}

func NewLogoutVerifier(issuerURL string, keySet gOidc.KeySet, config *gOidc.Config) *LogoutTokenVerifier {
	return &LogoutTokenVerifier{keySet: keySet, config: config, issuer: issuerURL}
}

//Upon receiving a logout request at the back-channel logout URI, the RP MUST validate the Logout Token as follows:
//
//1. If the Logout Token is encrypted, decrypt it using the keys and algorithms that the Client specified during Registration that the OP was to use to encrypt ID Tokens. If ID Token encryption was negotiated with the OP at Registration time and the Logout Token is not encrypted, the RP SHOULD reject it.
//2. Validate the Logout Token signature in the same way that an ID Token signature is validated, with the following refinements.
//3. Validate the iss, aud, and iat Claims in the same way they are validated in ID Tokens.
//4. Verify that the Logout Token contains a sub Claim, a sid Claim, or both.
//5. Verify that the Logout Token contains an events Claim whose value is JSON object containing the member name http://schemas.openid.net/event/backchannel-logout.
//6. Verify that the Logout Token does not contain a nonce Claim.
//7. Optionally verify that another Logout Token with the same jti value has not been recently received.
//If any of the validation steps fails, reject the Logout Token and return an HTTP 400 Bad Request error. Otherwise, proceed to perform the logout actions.

// Verify verifies a Logout token according to Specs
func (v *LogoutTokenVerifier) Verify(ctx context.Context, rawIDToken string) (*LogoutToken, error) {
	jws, err := jose.ParseSigned(rawIDToken)
	if err != nil {
		return nil, err
	}
	// Throw out tokens with invalid claims before trying to verify the token. This lets
	// us do cheap checks before possibly re-syncing keys.
	payload, err := parseJWT(rawIDToken)
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
	if !v.config.SkipIssuerCheck && token.Issuer != v.issuer {
		return nil, fmt.Errorf("oidc: id token issued by a different provider, expected %q got %q", v.issuer, token.Issuer)
	}

	// If a client ID has been provided, make sure it's part of the audience. SkipClientIDCheck must be true if ClientID is empty.
	//
	// This check DOES NOT ensure that the ClientID is the party to which the ID Token was issued (i.e. Authorized party).
	if !v.config.SkipClientIDCheck {
		if v.config.ClientID != "" {
			if !contains(token.Audience, v.config.ClientID) {
				return nil, fmt.Errorf("oidc: expected audience %q got %q", v.config.ClientID, token.Audience)
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
	supportedSigAlgs := v.config.SupportedSigningAlgs
	if len(supportedSigAlgs) == 0 {
		supportedSigAlgs = []string{gOidc.RS256}
	}

	if !contains(supportedSigAlgs, sig.Header.Algorithm) {
		return nil, fmt.Errorf("oidc: id token signed with unsupported algorithm, expected %q got %q", supportedSigAlgs, sig.Header.Algorithm)
	}

	gotPayload, err := v.keySet.VerifySignature(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %v", err)
	}

	// Ensure that the payload returned by the square actually matches the payload parsed earlier.
	if !bytes.Equal(gotPayload, payload) {
		return nil, errors.New("oidc: internal error, payload parsed did not match previous payload")
	}

	return &token, nil
}

func parseJWT(p string) ([]byte, error) {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", len(parts))
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %v", err)
	}
	return payload, nil
}

func contains(sli []string, ele string) bool {
	for _, s := range sli {
		if s == ele {
			return true
		}
	}
	return false
}

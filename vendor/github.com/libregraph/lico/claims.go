/*
 * Copyright 2017-2021 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package lico

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"

	"github.com/libregraph/lico/oidc"
	"github.com/libregraph/lico/oidc/payload"
)

// Access token claims used.
const (
	RefClaim              = "lg.r"
	IdentityClaim         = "lg.i"
	IdentityProviderClaim = "lg.p"
	ScopesClaim           = "scp"
)

// Identifier identity sub claims used.
const (
	IdentifiedUserClaim        = "us"
	IdentifiedUserIDClaim      = "id"
	IdentifiedUsernameClaim    = "un"
	IdentifiedDisplayNameClaim = "dn"
	IdentifiedData             = "da"
	IdentifiedUserIsGuest      = "gu"
)

// TokenType defines the token type value.
type TokenTypeValue string

// Is compares the associated TokenTypeValue to the provided one.
func (ttv TokenTypeValue) Is(value TokenTypeValue) bool {
	return ttv == value
}

// The known token type values.
const (
	TokenTypeIDToken      TokenTypeValue = "" // Just a placeholder, not actually set in ID Tokens.
	TokenTypeAccessToken  TokenTypeValue = "1"
	TokenTypeRefreshToken TokenTypeValue = "2"
)

// AccessTokenClaims define the claims found in access tokens issued.
type AccessTokenClaims struct {
	jwt.StandardClaims

	TokenType TokenTypeValue `json:"lg.t"`

	AuthorizedClaimsRequest *payload.ClaimsRequest `json:"lg.acr,omitempty"`

	AuthorizedScopesList payload.ScopesValue `json:"scp"`

	IdentityClaims   jwt.MapClaims `json:"lg.i"`
	IdentityProvider string        `json:"lg.p,omitempty"`

	*oidc.SessionClaims
}

// Valid implements the jwt.Claims interface.
func (c AccessTokenClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.IdentityClaims != nil {
		if err := c.IdentityClaims.Valid(); err != nil {
			return err
		}
	}
	if c.TokenType.Is(TokenTypeAccessToken) {
		return nil
	}
	return errors.New("not an access token")
}

// AuthorizedScopes returns a map with scope keys and true value of all scopes
// set in the accociated access token.
func (c AccessTokenClaims) AuthorizedScopes() map[string]bool {
	authorizedScopes := make(map[string]bool)
	for _, scope := range c.AuthorizedScopesList {
		authorizedScopes[scope] = true
	}

	return authorizedScopes
}

// RefreshTokenClaims define the claims used by refresh tokens.
type RefreshTokenClaims struct {
	jwt.StandardClaims

	TokenType TokenTypeValue `json:"lg.t"`

	ApprovedScopesList payload.ScopesValue `json:"scp"`

	ApprovedClaimsRequest *payload.ClaimsRequest `json:"lg.acr,omitempty"`
	Ref                   string                 `json:"lg.r"`

	IdentityClaims   jwt.MapClaims `json:"lg.i"`
	IdentityProvider string        `json:"lg.p,omitempty"`
}

// Valid implements the jwt.Claims interface.
func (c RefreshTokenClaims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.IdentityClaims != nil {
		if err := c.IdentityClaims.Valid(); err != nil {
			return err
		}
	}
	if c.TokenType.Is(TokenTypeRefreshToken) {
		return nil
	}
	return errors.New("not a refresh token")
}

// NumericIDClaims define the claims used with the konnect/id scope.
type NumericIDClaims struct {
	// NOTE(longsleep): Always keep these claims compatible with the GitLab API
	// https://docs.gitlab.com/ce/api/users.html#for-user.
	NumericID         int64  `json:"id,omitempty"`
	NumericIDUsername string `json:"username,omitempty"`
}

// Valid implements the jwt.Claims interface.
func (c NumericIDClaims) Valid() error {
	if c.NumericIDUsername == "" {
		return errors.New("username claim not valid")
	}
	return nil
}

// UniqueUserIDClaims define the claims used with the konnect/uuid scope.
type UniqueUserIDClaims struct {
	UniqueUserID string `json:"lg.uuid,omitempty"`
}

// Valid implements the jwt.Claims interface.
func (c UniqueUserIDClaims) Valid() error {
	if c.UniqueUserID == "" {
		return errors.New("lg.uuid claim not valid")
	}
	return nil
}

/*
 * Copyright 2017-2019 Kopano and its licensors
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

package payload

import (
	"encoding/json"
	"strings"

	"stash.kopano.io/kgol/oidc-go"
)

var scopedClaims = map[string]string{
	oidc.NameClaim:              oidc.ScopeProfile,
	oidc.FamilyNameClaim:        oidc.ScopeProfile,
	oidc.GivenNameClaim:         oidc.ScopeProfile,
	oidc.MiddleNameClaim:        oidc.ScopeProfile,
	oidc.PreferredUsernameClaim: oidc.ScopeProfile,
	oidc.ProfileClaim:           oidc.ScopeProfile,
	oidc.PictureClaim:           oidc.ScopeProfile,
	oidc.WebsiteClaim:           oidc.ScopeProfile,
	oidc.GenderClaim:            oidc.ScopeProfile,
	oidc.BirthdateClaim:         oidc.ScopeProfile,
	oidc.ZoneinfoClaim:          oidc.ScopeProfile,
	oidc.UpdatedAtClaim:         oidc.ScopeProfile,

	oidc.EmailClaim:         oidc.ScopeEmail,
	oidc.EmailVerifiedClaim: oidc.ScopeEmail,
}

// GetScopeForClaim returns the known scope if any for the provided claim name.
func GetScopeForClaim(claim string) (string, bool) {
	scope, ok := scopedClaims[claim]
	return scope, ok
}

// ClaimsRequest define the base claims structure for OpenID Connect claims
// request parameter value as specified at
// https://openid.net/specs/openid-connect-core-1_0.html#ClaimsParameter - in
// addition a Konnect specific pass thru value can be used to pass through any
// application specific values to access and reqfresh tokens.
type ClaimsRequest struct {
	UserInfo *ClaimsRequestMap `json:"userinfo,omitempty"`
	IDToken  *ClaimsRequestMap `json:"id_token,omitempty"`
	Passthru json.RawMessage   `json:"passthru,omitempty"`
}

// ApplyScopes removes all claims requests from the accociated claims request
// which are not mapped to one of the provided approved scopes.
func (cr *ClaimsRequest) ApplyScopes(approvedScopes map[string]bool) error {
	if cr.UserInfo != nil {
		for claim := range *cr.UserInfo {
			if approved := approvedScopes[scopedClaims[claim]]; !approved {
				delete(*cr.UserInfo, claim)
			}
		}
	}
	if cr.IDToken != nil {
		for claim := range *cr.IDToken {
			if approved := approvedScopes[scopedClaims[claim]]; !approved {
				delete(*cr.IDToken, claim)
			}
		}
	}

	return nil
}

// Scopes adds all scopes of the accociated claims requests claims to
// the provied scopes mapping safe the scopes already defined in the provided
// excluded scopes mapping.
func (cr *ClaimsRequest) Scopes(excludedScopes map[string]bool) []string {
	scopesMap := make(map[string]bool)

	if cr.UserInfo != nil {
		for claim := range *cr.UserInfo {
			scope := scopedClaims[claim]
			if _, excluded := excludedScopes[scope]; !excluded {
				scopesMap[scope] = true
			}
		}
	}
	if cr.IDToken != nil {
		for claim := range *cr.IDToken {
			scope := scopedClaims[claim]
			if _, excluded := excludedScopes[scope]; !excluded {
				scopesMap[scope] = true
			}
		}
	}

	scopes := make([]string, 0)
	for scope := range scopesMap {
		scopes = append(scopes, scope)
	}

	return scopes
}

// ClaimsRequestMap defines a mapping of claims request values used with
// OpenID Connect claims request parameter values.
type ClaimsRequestMap map[string]*ClaimsRequestValue

// ScopesMap returns a map of scopes defined by the claims in tha associated map.
func (crm *ClaimsRequestMap) ScopesMap(excludedScopes map[string]bool) map[string]bool {
	scopesMap := make(map[string]bool)

	for claim := range *crm {
		scope := scopedClaims[claim]
		if _, excluded := excludedScopes[scope]; !excluded {
			scopesMap[scope] = true
		}
	}

	return scopesMap
}

// Get returns the accociated maps claim value identified by the provided name.
func (crm ClaimsRequestMap) Get(claim string) (*ClaimsRequestValue, bool) {
	value, ok := crm[claim]

	return value, ok
}

// GetStringValue returns the accociated maps claim value identified by the
// provided name as string value.
func (crm ClaimsRequestMap) GetStringValue(claim string) (string, bool) {
	value, ok := crm.Get(claim)
	if !ok {
		return "", false
	}

	s, ok := value.Value.(string)
	return s, ok
}

// ClaimsRequestValue is the claims request detail definition of an OpenID
// Connect claims request parameter value.
type ClaimsRequestValue struct {
	Essential bool          `json:"essential,omitempty"`
	Value     interface{}   `json:"value,omitempty"`
	Values    []interface{} `json:"values,omitempty"`
}

// Match returns true of the provided value is contained inside the accociated
// request values values or value.
func (crv *ClaimsRequestValue) Match(value interface{}) bool {
	if len(crv.Values) == 0 {
		return value == crv.Value
	}
	for _, v := range crv.Values {
		if v == value {
			return true
		}
	}

	return false
}

// ScopesValue is a string array with JSON marshal to/from a space separated
// single string value.
type ScopesValue []string

func (sv ScopesValue) MarshalJSON() ([]byte, error) {
	result := strings.Join(sv, " ")
	return json.Marshal(&result)
}

func (sv *ScopesValue) UnmarshalJSON(data []byte) error {
	var parsed string
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return err
	}

	result := ScopesValue(strings.Split(parsed, " "))
	*sv = result
	return nil
}

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

package authorities

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-jose/go-jose/v3"
)

// Supported Authority kind string values.
const (
	AuthorityTypeOIDC  = "oidc"
	AuthorityTypeSAML2 = "saml2"
)

type authorityRegistrationData struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AuthorityType string `json:"authority_type"`

	Iss string `json:"iss"`

	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`

	EntityID string `json:"entity_id"`

	Trusted  bool  `json:"trusted"`
	Insecure bool  `json:"insecure"`
	Default  bool  `json:"default"`
	Discover *bool `json:"discover"`

	Scopes              []string `json:"scopes"`
	ResponseType        string   `json:"response_type"`
	ResponseMode        string   `json:"response_mode"`
	CodeChallengeMethod string   `json:"code_challenge_method"`

	RawMetadataEndpoint      string `json:"metadata_endpoint"`
	RawAuthorizationEndpoint string `json:"authorization_endpoint"`
	RawTokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint         string `json:"user_info_endpoint"`

	JWKS *jose.JSONWebKeySet `json:"jwks"`

	IdentityClaimName string `json:"identity_claim_name"`

	IdentityAliases       map[string]string `json:"identity_aliases"`
	IdentityAliasRequired bool              `json:"identity_alias_required"`

	EndSessionEnabled bool `json:"end_session_enabled"`
}

type authorityRegistryData struct {
	Authorities []*authorityRegistrationData `json:"authorities"`
}

// AuthorityRegistration defines an authority with its properties.
type AuthorityRegistration interface {
	ID() string
	Name() string
	AuthorityType() string

	Authority() *Details
	Issuer() string

	Validate() error

	Initialize(ctx context.Context, registry *Registry) error

	MakeRedirectAuthenticationRequestURL(state string) (*url.URL, map[string]interface{}, error)
	MakeRedirectEndSessionRequestURL(ref interface{}, state string) (*url.URL, map[string]interface{}, error)
	MakeRedirectEndSessionResponseURL(req interface{}, state string) (*url.URL, map[string]interface{}, error)

	ParseStateResponse(req *http.Request, state string, extra map[string]interface{}) (interface{}, error)

	ValidateIdpEndSessionRequest(req interface{}, state string) (bool, error)
	ValidateIdpEndSessionResponse(res interface{}, state string) (bool, error)

	IdentityClaimValue(data interface{}) (string, map[string]interface{}, error)

	Metadata() AuthorityMetadata
}

type AuthorityMetadata interface {
}

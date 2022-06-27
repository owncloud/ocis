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

	"gopkg.in/square/go-jose.v2"
)

// Supported Authority kind string values.
const (
	AuthorityTypeOIDC  = "oidc"
	AuthorityTypeSAML2 = "saml2"
)

type authorityRegistrationData struct {
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	AuthorityType string `yaml:"authority_type"`

	Iss string `yaml:"iss"`

	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	EntityID string `yaml:"entity_id"`

	Trusted  bool  `yaml:"trusted"`
	Insecure bool  `yaml:"insecure"`
	Default  bool  `yaml:"default"`
	Discover *bool `yaml:"discover"`

	Scopes              []string `yaml:"scopes"`
	ResponseType        string   `yaml:"response_type"`
	CodeChallengeMethod string   `yaml:"code_challenge_method"`

	RawMetadataEndpoint      string `yaml:"metadata_endpoint"`
	RawAuthorizationEndpoint string `yaml:"authorization_endpoint"`

	JWKS *jose.JSONWebKeySet `yaml:"jwks"`

	IdentityClaimName string `yaml:"identity_claim_name"`

	IdentityAliases       map[string]string `yaml:"identity_aliases,flow"`
	IdentityAliasRequired bool              `yaml:"identity_alias_required"`

	EndSessionEnabled bool `yaml:"end_session_enabled"`
}

type authorityRegistryData struct {
	Authorities []*authorityRegistrationData `yaml:"authorities,flow"`
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

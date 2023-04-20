/*
 * Copyright 2019 Kopano
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

// WellKnown defines the OpenID Connect 1.0 discovery provider meta data as
// specified at https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type WellKnown struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	UserInfoEndpoint       string   `json:"userinfo_endpoint"`
	EndSessionEndpoint     string   `json:"end_session_endpoint"`
	RegistrationEndpoint   string   `json:"registration_endpoint,omitempty"`
	CheckSessionIframe     string   `json:"check_session_iframe,omitempty"`
	JwksURI                string   `json:"jwks_uri"`
	ScopesSupported        []string `json:"scopes_supported"`
	ResponseTypesSupported []string `json:"response_types_supported"`
	SubjectTypesSupported  []string `json:"subject_types_supported"`

	IDTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
	UserInfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported"`
	RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported"`
	TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`

	ClaimsParameterSupported bool     `json:"claims_parameter_supported"`
	ClaimsSupported          []string `json:"claims_supported"`

	RequestParameterSupported    bool `json:"request_parameter_supported"`
	RequestURIParameterSupported bool `json:"request_uri_parameter_supported"`
}

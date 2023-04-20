// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package wellknown

import (
	"encoding/json"
	"net/http"

	"github.com/cs3org/reva/v2/pkg/appctx"
)

func (s *svc) doOpenidConfiguration(w http.ResponseWriter, r *http.Request) {
	log := appctx.GetLogger(r.Context())
	pm := &ProviderMetadata{
		Issuer:                s.conf.Issuer,
		AuthorizationEndpoint: s.conf.AuthorizationEndpoint,
		JwksURI:               s.conf.JwksURI,
		TokenEndpoint:         s.conf.TokenEndpoint,
		RevocationEndpoint:    s.conf.RevocationEndpoint,
		IntrospectionEndpoint: s.conf.IntrospectionEndpoint,
		UserinfoEndpoint:      s.conf.UserinfoEndpoint,
		EndSessionEndpoint:    s.conf.EndSessionEndpoint,
	}

	b, err := json.Marshal(pm)
	if err != nil {
		log.Error().Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Error().Err(err).Msg("Error writing response")
		return
	}
}

// The ProviderMetadata describes an idp.
// see https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
// TODO(labkode): do we really need it to validate the token and get user claims?
type ProviderMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	// claims_parameter_supported
	ClaimsSupported []string `json:"claims_supported,omitempty"`
	// grant_types_supported
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`
	Issuer                           string   `json:"issuer,omitempty"`
	JwksURI                          string   `json:"jwks_uri,omitempty"`
	// registration_endpoint
	// request_object_signing_alg_values_supported
	// request_parameter_supported
	// request_uri_parameter_supported
	// require_request_uri_registration
	// response_modes_supported
	ResponseTypesSupported []string `json:"response_types_supported,omitempty"`
	ScopesSupported        []string `json:"scopes_supported,omitempty"`
	SubjectTypesSupported  []string `json:"subject_types_supported,omitempty"`
	TokenEndpoint          string   `json:"token_endpoint,omitempty"`
	// token_endpoint_auth_methods_supported
	// token_endpoint_auth_signing_alg_values_supported
	UserinfoEndpoint string `json:"userinfo_endpoint,omitempty"`
	// userinfo_signing_alg_values_supported
	// code_challenge_methods_supported
	IntrospectionEndpoint string `json:"introspection_endpoint,omitempty"`
	// introspection_endpoint_auth_methods_supported
	// introspection_endpoint_auth_signing_alg_values_supported
	RevocationEndpoint string `json:"revocation_endpoint,omitempty"`
	// revocation_endpoint_auth_methods_supported
	// revocation_endpoint_auth_signing_alg_values_supported
	// id_token_encryption_alg_values_supported
	// id_token_encryption_enc_values_supported
	// userinfo_encryption_alg_values_supported
	// userinfo_encryption_enc_values_supported
	// request_object_encryption_alg_values_supported
	// request_object_encryption_enc_values_supported
	CheckSessionIframe string `json:"check_session_iframe,omitempty"`
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`
	// claim_types_supported
}

package oidc

const (
	Iss               = "iss"
	Sub               = "sub"
	Email             = "email"
	Name              = "name"
	PreferredUsername = "preferred_username"
	UIDNumber         = "uidnumber"
	GIDNumber         = "gidnumber"
	Groups            = "groups"
	OwncloudUUID      = "ownclouduuid"
	OcisRoutingPolicy = "ocis.routing.policy"
)

// The ProviderMetadata describes an idp.
// see https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type ProviderMetadata struct {
	AuthorizationEndpoint string `json:"authorization_endpoint,omitempty"`
	//claims_parameter_supported
	ClaimsSupported []string `json:"claims_supported,omitempty"`
	//grant_types_supported
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`
	Issuer                           string   `json:"issuer,omitempty"`
	JwksURI                          string   `json:"jwks_uri,omitempty"`
	//registration_endpoint
	//request_object_signing_alg_values_supported
	//request_parameter_supported
	//request_uri_parameter_supported
	//require_request_uri_registration
	//response_modes_supported
	ResponseTypesSupported []string `json:"response_types_supported,omitempty"`
	ScopesSupported        []string `json:"scopes_supported,omitempty"`
	SubjectTypesSupported  []string `json:"subject_types_supported,omitempty"`
	TokenEndpoint          string   `json:"token_endpoint,omitempty"`
	//token_endpoint_auth_methods_supported
	//token_endpoint_auth_signing_alg_values_supported
	UserinfoEndpoint string `json:"userinfo_endpoint,omitempty"`
	//userinfo_signing_alg_values_supported
	//code_challenge_methods_supported
	IntrospectionEndpoint string `json:"introspection_endpoint,omitempty"`
	//introspection_endpoint_auth_methods_supported
	//introspection_endpoint_auth_signing_alg_values_supported
	RevocationEndpoint string `json:"revocation_endpoint,omitempty"`
	//revocation_endpoint_auth_methods_supported
	//revocation_endpoint_auth_signing_alg_values_supported
	//id_token_encryption_alg_values_supported
	//id_token_encryption_enc_values_supported
	//userinfo_encryption_alg_values_supported
	//userinfo_encryption_enc_values_supported
	//request_object_encryption_alg_values_supported
	//request_object_encryption_enc_values_supported
	CheckSessionIframe string `json:"check_session_iframe,omitempty"`
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`
	//claim_types_supported
}

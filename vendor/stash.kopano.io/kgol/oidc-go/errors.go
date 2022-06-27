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

// OAuth2 error codes.
const (
	ErrorCodeOAuth2UnsupportedResponseType = "unsupported_response_type"
	ErrorCodeOAuth2InvalidRequest          = "invalid_request"
	ErrorCodeOAuth2InvalidToken            = "invalid_token"
	ErrorCodeOAuth2InsufficientScope       = "insufficient_scope"
	ErrorCodeOAuth2InvalidGrant            = "invalid_grant"
	ErrorCodeOAuth2UnsupportedGrantType    = "unsupported_grant_type"
	ErrorCodeOAuth2AccessDenied            = "access_denied"
	ErrorCodeOAuth2ServerError             = "server_error"
	ErrorCodeOAuth2TemporarilyUnavailable  = "temporarily_unavailable"
)

// OIDC error codes.
const (
	ErrorCodeOIDCInteractionRequired = "interaction_required"
	ErrorCodeOIDCLoginRequired       = "login_required"
	ErrorCodeOIDCConsentRequired     = "consent_required"

	ErrorCodeOIDCRequestNotSupported      = "request_not_supported"
	ErrorCodeOIDCInvalidRequestObject     = "invalid_request_object"
	ErrorCodeOIDCRequestURINotSupported   = "request_uri_not_supported"
	ErrorCodeOIDCRegistrationNotSupported = "registration_not_supported"

	ErrorCodeOIDCInvalidRedirectURI    = "invalid_redirect_uri"
	ErrorCodeOIDCInvalidClientMetadata = "invalid_client_metadata"
)

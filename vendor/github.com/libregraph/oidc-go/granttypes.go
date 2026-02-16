/*
 * Copyright 2017-2019 Kopano
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

const (
	// GrantTypeAuthorizationCode is the string value for the
	// OAuth2 authroization code token request grant type.
	GrantTypeAuthorizationCode = "authorization_code"

	// GrantTypeImplicit is the string value for the OAuth2 id_token, token
	// id_token token request grant type.
	GrantTypeImplicit = "implicit"

	// GrantTypeRefreshToken is the string value for the OAuth2 refresh_token
	// token request grant_type.
	GrantTypeRefreshToken = "refresh_token"
)

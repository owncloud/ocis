/*
 * Copyright 2017 Kopano
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

// OIDC response types and flows.
const (
	ResponseTypeCode             = "code"                // OIDC code flow
	ResponseTypeIDTokenToken     = "id_token token"      // OIDC implicit flow
	ResponseTypeIDToken          = "id_token"            // OIDC implicit flow
	ResponseTypeCodeIDToken      = "code id_token"       // OIDC hybrid flow
	ResponseTypeCodeToken        = "code token"          // OIDC hybrid flow
	ResponseTypeCodeIDTokenToken = "code id_token token" // OIDC hybrid flow
	ResponseTypeToken            = "token"               // OAuth2

	ResponseModeFragment = "fragment"
	ResponseModeQuery    = "query"

	FlowCode     = "code"
	FlowImplicit = "implicit"
	FlowHybrid   = "hybrid"
)

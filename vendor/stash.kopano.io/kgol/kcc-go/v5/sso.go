/*
 * Copyright 2018-2019 Kopano and its licensors
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

package kcc

// SSOType is the type of SSO to use with single sign on.
type SSOType string

func (sst SSOType) String() string {
	return string(sst)
}

// Known Kopano SSO types.
const (
	KOPANO_SSO_TYPE_NTML   SSOType = "NTLM"
	KOPANO_SSO_TYPE_KCOIDC SSOType = "KCOIDC"
	KOPANO_SSO_TYPE_KRB5   SSOType = ""
)

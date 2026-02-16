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

// Standard claims as used in JSON Web Tokens.
const (
	IssuerIdentifierClaim  = "iss"
	SubjectIdentifierClaim = "sub"
	AudienceClaim          = "aud"
	ExpirationClaim        = "exp"
	IssuedAtClaim          = "iat"
)

// Additional claims as defined by OIDC.
const (
	NameClaim              = "name"
	FamilyNameClaim        = "family_name"
	GivenNameClaim         = "given_name"
	MiddleNameClaim        = "middle_name"
	NicknameClaim          = "nickname"
	PreferredUsernameClaim = "preferred_username"
	ProfileClaim           = "profile"
	PictureClaim           = "picture"
	WebsiteClaim           = "website"
	GenderClaim            = "gender"
	BirthdateClaim         = "birthdate"
	ZoneinfoClaim          = "zoneinfo"
	LocaleClaim            = "locale"
	UpdatedAtClaim         = "updated_at"

	EmailClaim         = "email"
	EmailVerifiedClaim = "email_verified"

	AuthTimeClaim = "auth_time"
)

// Additional claims as defined by OIDC extensions.
const (
	SessionIDClaim = "sid"
)

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

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
)

// Code challenge methods implemented by Konnect. See https://tools.ietf.org/html/rfc7636.
const (
	PlainCodeChallengeMethod = "plain"
	S256CodeChallengeMethod  = "S256"
)

// ValidateCodeChallenge implements https://tools.ietf.org/html/rfc7636#section-4.6
// code challenge verification.
func ValidateCodeChallenge(challenge string, method string, verifier string) error {
	if method == "" {
		// We default to S256CodeChallengeMethod.
		method = S256CodeChallengeMethod
	}

	computed, err := MakeCodeChallenge(method, verifier)
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare([]byte(challenge), []byte(computed)) != 1 {
		return errors.New("invalid code challenge")
	}
	return nil
}

// MakeCodeChallenge creates a code challenge using the provided method and
// verifier for https://tools.ietf.org/html/rfc7636#section-4.6 verification.
func MakeCodeChallenge(method string, verifier string) (string, error) {
	if verifier == "" {
		return "", errors.New("invalid verifier")
	}

	switch method {
	case PlainCodeChallengeMethod:
		// Challenge is verifier.
		return verifier, nil
	case S256CodeChallengeMethod:
		// BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
		// Base64 encoding using the URL- and filename-safe character set
		// defined in Section 5 of [RFC4648], with all trailing '='
		// characters omitted (as permitted by Section 3.2 of [RFC4648]) and
		// without the inclusion of any line breaks, whitespace, or other
		// additional characters.
		sum := sha256.Sum256([]byte(verifier))
		return base64.RawURLEncoding.EncodeToString(sum[:]), nil
	}

	return "", errors.New("transform algorithm not supported")
}

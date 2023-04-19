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

package signing

import (
	"crypto"
	"encoding/base64"

	jwk "github.com/mendsley/gojwk"
	"golang.org/x/crypto/ed25519"
)

// JWKFromPublicKey creates a JWK from a public key
func JWKFromPublicKey(key crypto.PublicKey) (*jwk.Key, error) {
	switch key := key.(type) {
	case ed25519.PublicKey:
		jwt := &jwk.Key{
			Kty: "OKP",
			Crv: "Ed25519",
			X:   base64.RawURLEncoding.EncodeToString(key),
		}

		return jwt, nil

	default:
		return jwk.PublicKey(key)
	}

}

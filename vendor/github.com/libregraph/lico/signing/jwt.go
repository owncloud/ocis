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
	"crypto/rand"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/ed25519"
)

// Errors used by this package.
var (
	ErrEdDSAVerification = errors.New("eddsa: verification error")
)

// SigningMethodEdwardsCurve implements the EdDSA family of signing methods.
type SigningMethodEdwardsCurve struct {
	Name string
}

// Specific instances for EdDSA
var (
	SigningMethodEdDSA *SigningMethodEdwardsCurve
)

func init() {
	// EdDSA with Ed25519 https://tools.ietf.org/html/rfc8037#section-3.1.
	SigningMethodEdDSA = &SigningMethodEdwardsCurve{"EdDSA"}
	jwt.RegisterSigningMethod(SigningMethodEdDSA.Alg(), func() jwt.SigningMethod {
		return SigningMethodEdDSA
	})
}

// Alg implements the jwt.SigningMethod interface.
func (m *SigningMethodEdwardsCurve) Alg() string {
	return m.Name
}

// Verify implements the jwt.SigningMethod interface.
func (m *SigningMethodEdwardsCurve) Verify(signingString, signature string, key interface{}) error {
	var err error

	// Decode the signature
	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	// Get the key
	switch k := key.(type) {
	case ed25519.PublicKey:
		if len(k) != ed25519.PublicKeySize {
			return jwt.ErrInvalidKey
		}
		if verifystatus := ed25519.Verify(k, []byte(signingString), sig); verifystatus == true {
			return nil
		} else {
			return ErrEdDSAVerification
		}

	default:
		return jwt.ErrInvalidKeyType
	}
}

// Sign implements the jwt.SigningMethod interface.
func (m *SigningMethodEdwardsCurve) Sign(signingString string, key interface{}) (string, error) {
	switch k := key.(type) {
	case ed25519.PrivateKey:
		if len(k) != ed25519.PrivateKeySize {
			return "", jwt.ErrInvalidKey
		}
		if s, err := k.Sign(rand.Reader, []byte(signingString), crypto.Hash(0)); err == nil {
			// We serialize the signature.
			return jwt.EncodeSegment(s), nil
		} else {
			return "", err
		}

	default:
		return "", jwt.ErrInvalidKeyType
	}
}

/*
 * Copyright 2020 Kopano and its licensors
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

package samlext

import (
	"crypto"
	"crypto/rsa"
	_ "crypto/sha1" // Import all supported hashers.
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"strings"

	dsig "github.com/russellhaering/goxmldsig"
)

// VerifySignedHTTPRedirectQuery implements validation for signed SAML HTTP
// redirect binding parameters provides via URL query.
func VerifySignedHTTPRedirectQuery(kind string, rawQuery string, sigAlg string, signature []byte, pubKey crypto.PublicKey) error {
	var hasher crypto.Hash

	// Validate signature.
	switch sigAlg {
	case dsig.RSASHA1SignatureMethod:
		hasher = crypto.SHA1

	case dsig.RSASHA256SignatureMethod:
		hasher = crypto.SHA256

	case dsig.RSASHA512SignatureMethod:
		hasher = crypto.SHA512

	default:
		return fmt.Errorf("unsupported sig alg: %v", sigAlg)
	}

	if len(signature) == 0 {
		return fmt.Errorf("signature data is empty")
	}

	// The signed data format goes like this:
	// SAMLRequest=urlencode(base64(deflate($xml)))&RelayState=urlencode($(relay_state))&SigAlg=urlencode($sig_alg)
	// We rebuild it ourselves from the raw request, to avoid differences when url decoding/encoding.
	signedQuery := func(query string) string {
		m := make(map[string]string)
		for query != "" {
			key := query
			if i := strings.IndexAny(key, "&;"); i >= 0 {
				key, query = key[:i], key[i+1:]
			} else {
				query = ""
			}
			if key == "" {
				continue
			}
			value := ""
			if i := strings.Index(key, "="); i >= 0 {
				key, value = key[:i], key[i+1:]
			}
			m[key] = value // Support only one value, but thats ok since we only want one and if really someone signed multiple values then its fine to fail.
		}
		s := new(strings.Builder)
		for idx, key := range []string{kind, "RelayState", "SigAlg"} {
			if value, ok := m[key]; ok {
				if idx > 0 {
					s.WriteString("&")
				}
				s.WriteString(key)
				s.WriteString("=")
				s.WriteString(value)
			}
		}
		return s.String()
	}(rawQuery)

	// Create hash for the alg.
	hash := hasher.New()
	if _, hashErr := hash.Write([]byte(signedQuery)); hashErr != nil {
		return fmt.Errorf("failed to hash: %w", hashErr)
	}
	hashed := hash.Sum(nil)

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid RSA public key")
	}

	// NOTE(longsleep): All sig algs above, use PKCS1v15 with RSA.
	if verifyErr := rsa.VerifyPKCS1v15(rsaPubKey, hasher, hashed, signature); verifyErr != nil {
		return fmt.Errorf("signature verification failed: %w", verifyErr)
	}
	return nil
}

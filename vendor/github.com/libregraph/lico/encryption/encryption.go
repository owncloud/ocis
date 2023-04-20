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

package encryption

import (
	"fmt"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// KeySize is the size of the keys created by GenerateKey()
	KeySize = 32
	// NonceSize is the size of the nonces created by GenerateNonce()
	NonceSize = 24
)

// Encrypt generates a random nonce and encrypts the input using nacl.secretbox
// package. We store the nonce in the first 24 bytes of the encrypted text.
func Encrypt(msg []byte, key *[KeySize]byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, err
	}

	return encryptWithNonce(msg, nonce, key)
}

func encryptWithNonce(msg []byte, nonce *[NonceSize]byte, key *[KeySize]byte) ([]byte, error) {
	encrypted := secretbox.Seal(nonce[:], msg, nonce, key)
	return encrypted, nil
}

// Decrypt extracts the nonce from the encrypted text, and attempts to decrypt
// with nacl.box.
func Decrypt(msg []byte, key *[KeySize]byte) ([]byte, error) {
	if len(msg) < (NonceSize + secretbox.Overhead) {
		return nil, fmt.Errorf("wrong length of ciphertext")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], msg[:NonceSize])
	decrypted, ok := secretbox.Open(nil, msg[NonceSize:], &nonce, key)
	if !ok {
		return nil, fmt.Errorf("decryption failed")
	}

	return decrypted, nil
}

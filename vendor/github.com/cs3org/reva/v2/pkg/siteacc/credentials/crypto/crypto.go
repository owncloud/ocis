// Copyright 2018-2020 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/pkg/errors"
)

const (
	passphraseLength = 32
)

// EncodeString encodes a string using AES and returns the base64-encoded result.
func EncodeString(s string, passphrase string) (string, error) {
	if len(s) == 0 || len(passphrase) == 0 {
		return "", nil
	}
	passphrase = normalizePassphrase(passphrase)

	gcm, err := createGCM([]byte(passphrase))
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "unable to generate nonce")
	}
	encryptedData := gcm.Seal(nonce, nonce, []byte(s), nil)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// DecodeString decodes a base64-encoded string encoded with AES.
func DecodeString(s string, passphrase string) (string, error) {
	if len(s) == 0 || len(passphrase) == 0 {
		return "", nil
	}
	data, _ := base64.StdEncoding.DecodeString(s)
	passphrase = normalizePassphrase(passphrase)

	gcm, err := createGCM([]byte(passphrase))
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(s) < nonceSize {
		return "", errors.Errorf("input string length too short")
	}
	nonce, data := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode string")
	}
	return string(plain), nil
}

func createGCM(passphrase []byte) (cipher.AEAD, error) {
	c, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate cipher")
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate GCM")
	}
	return gcm, nil
}

func normalizePassphrase(passphrase string) string {
	if len(passphrase) > passphraseLength {
		passphrase = passphrase[:passphraseLength]
	} else if len(passphrase) < passphraseLength {
		for i := len(passphrase); i < passphraseLength; i++ {
			passphrase += "#"
		}
	}
	return passphrase
}

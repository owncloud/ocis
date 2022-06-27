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

package credentials

import (
	"github.com/cs3org/reva/v2/pkg/siteacc/credentials/crypto"
	"github.com/pkg/errors"
)

// Credentials stores and en-/decrypts credentials
type Credentials struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

// Get decrypts and retrieves the stored credentials.
func (creds *Credentials) Get(passphrase string) (string, string, error) {
	id, err := crypto.DecodeString(creds.ID, passphrase)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to decode ID")
	}
	secret, err := crypto.DecodeString(creds.Secret, passphrase)
	if err != nil {
		return "", "", errors.Wrap(err, "unable to decode secret")
	}
	return id, secret, nil
}

// Set encrypts and sets new credentials.
func (creds *Credentials) Set(id, secret string, passphrase string) error {
	if s, err := crypto.EncodeString(id, passphrase); err == nil {
		creds.ID = s
	} else {
		return errors.Wrap(err, "unable to encode ID")
	}
	if s, err := crypto.EncodeString(secret, passphrase); err == nil {
		creds.Secret = s
	} else {
		return errors.Wrap(err, "unable to encode secret")
	}
	return nil
}

// IsValid checks whether the credentials are valid.
func (creds *Credentials) IsValid() bool {
	return len(creds.ID) > 0 && len(creds.Secret) > 0
}

// Clear resets the credentials.
func (creds *Credentials) Clear() {
	creds.ID = ""
	creds.Secret = ""
}

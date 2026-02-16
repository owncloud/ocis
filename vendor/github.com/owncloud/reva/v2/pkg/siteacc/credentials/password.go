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
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Password holds a hash password alongside its salt value.
type Password struct {
	Value string `json:"value"`
}

const (
	passwordMinLength = 8
)

// Set sets a new password by hashing the plaintext version using bcrypt.
func (password *Password) Set(pwd string) error {
	if err := VerifyPassword(pwd); err != nil {
		return errors.Wrap(err, "invalid password")
	}

	pwdData, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "unable to generate password hash")
	}
	password.Value = string(pwdData)
	return nil
}

// Compare checks whether the given password string equals the stored one.
func (password *Password) Compare(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password.Value), []byte(pwd)) == nil
}

// IsValid checks whether the password is valid.
func (password *Password) IsValid() bool {
	// bcrypt hashes are in the form of $[version]$[cost]$[22 character salt][31 character hash], so they have a minimum length of 58
	return len(password.Value) > 58 && strings.Count(password.Value, "$") >= 3
}

// Clear resets the password.
func (password *Password) Clear() {
	password.Value = ""
}

// VerifyPassword checks whether the given password abides to the enforced password strength.
func VerifyPassword(pwd string) error {
	if len(pwd) < passwordMinLength {
		return errors.Errorf("the password must be at least 8 characters long")
	}
	if !strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") {
		return errors.Errorf("the password must contain at least one lowercase letter")
	}
	if !strings.ContainsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return errors.Errorf("the password must contain at least one uppercase letter")
	}
	if !strings.ContainsAny(pwd, "0123456789") {
		return errors.Errorf("the password must contain at least one digit")
	}

	return nil
}

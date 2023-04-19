// Copyright 2018-2020 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this filePath except in compliance with the License.
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

package data

import (
	"strings"
	"time"

	"github.com/cs3org/reva/v2/pkg/siteacc/credentials"
	"github.com/pkg/errors"

	"github.com/cs3org/reva/v2/pkg/utils"
)

// Account represents a single site account.
type Account struct {
	Email       string `json:"email"`
	Title       string `json:"title"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Site        string `json:"site"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phoneNumber"`

	Password credentials.Password `json:"password"`

	DateCreated  time.Time `json:"dateCreated"`
	DateModified time.Time `json:"dateModified"`

	Data     AccountData     `json:"data"`
	Settings AccountSettings `json:"settings"`
}

// AccountData holds additional data for a site account.
type AccountData struct {
	GOCDBAccess bool `json:"gocdbAccess"`
	SiteAccess  bool `json:"siteAccess"`
}

// AccountSettings holds additional settings for a site account.
type AccountSettings struct {
	ReceiveAlerts bool `json:"receiveAlerts"`
}

// Accounts holds an array of site accounts.
type Accounts = []*Account

// Update copies the data of the given account to this account.
func (acc *Account) Update(other *Account, setPassword bool, copyData bool) error {
	if err := other.verify(false, false); err != nil {
		return errors.Wrap(err, "unable to update account data")
	}

	// Manually update fields
	acc.Title = other.Title
	acc.FirstName = other.FirstName
	acc.LastName = other.LastName
	acc.Role = other.Role
	acc.PhoneNumber = other.PhoneNumber

	if setPassword && other.Password.Value != "" {
		// If a password was provided, use that as the new one
		if err := acc.UpdatePassword(other.Password.Value); err != nil {
			return errors.Wrap(err, "unable to update account data")
		}
	}

	if copyData {
		acc.Data = other.Data
	}

	return nil
}

// Configure copies the settings of the given account to this account.
func (acc *Account) Configure(other *Account) error {
	// Simply copy the stored settings
	acc.Settings = other.Settings

	return nil
}

// UpdatePassword assigns a new password to the account, hashing it first.
func (acc *Account) UpdatePassword(pwd string) error {
	if err := acc.Password.Set(pwd); err != nil {
		return errors.Wrap(err, "unable to update the user password")
	}
	return nil
}

// Clone creates a copy of the account; if erasePassword is set to true, the password will be cleared in the cloned object.
func (acc *Account) Clone(erasePassword bool) *Account {
	clone := *acc

	if erasePassword {
		clone.Password.Clear()
	}

	return &clone
}

// CheckScopeAccess checks whether the user can access the specified scope.
func (acc *Account) CheckScopeAccess(scope string) bool {
	hasAccess := false

	switch strings.ToLower(scope) {
	case ScopeDefault:
		hasAccess = true

	case ScopeGOCDB:
		hasAccess = acc.Data.GOCDBAccess

	case ScopeSite:
		hasAccess = acc.Data.SiteAccess
	}

	return hasAccess
}

// Cleanup trims all string entries.
func (acc *Account) Cleanup() {
	acc.Email = strings.TrimSpace(acc.Email)
	acc.Title = strings.TrimSpace(acc.Title)
	acc.FirstName = strings.TrimSpace(acc.FirstName)
	acc.LastName = strings.TrimSpace(acc.LastName)
	acc.Site = strings.TrimSpace(acc.Site)
	acc.Role = strings.TrimSpace(acc.Role)
	acc.PhoneNumber = strings.TrimSpace(acc.PhoneNumber)
}

func (acc *Account) verify(isNewAccount, verifyPassword bool) error {
	if acc.Email == "" {
		return errors.Errorf("no email address provided")
	} else if !utils.IsEmailValid(acc.Email) {
		return errors.Errorf("invalid email address: %v", acc.Email)
	}

	if acc.FirstName == "" {
		return errors.Errorf("no first name provided")
	} else if !utils.IsValidName(acc.FirstName) {
		return errors.Errorf("first name contains invalid characters: %v", acc.FirstName)
	}

	if acc.LastName == "" {
		return errors.Errorf("no last name provided")
	} else if !utils.IsValidName(acc.LastName) {
		return errors.Errorf("last name contains invalid characters: %v", acc.LastName)
	}

	if isNewAccount && acc.Site == "" {
		return errors.Errorf("no site provided")
	}

	if acc.Role == "" {
		return errors.Errorf("no role provided")
	} else if !utils.IsValidName(acc.Role) {
		return errors.Errorf("role contains invalid characters: %v", acc.Role)
	}

	if acc.PhoneNumber != "" && !utils.IsValidPhoneNumber(acc.PhoneNumber) {
		return errors.Errorf("invalid phone number provided")
	}

	if verifyPassword {
		if !acc.Password.IsValid() {
			return errors.Errorf("no valid password set")
		}
	}

	return nil
}

// NewAccount creates a new site account.
func NewAccount(email string, title, firstName, lastName string, site, role string, phoneNumber string, password string) (*Account, error) {
	t := time.Now()

	acc := &Account{
		Email:        email,
		Title:        title,
		FirstName:    firstName,
		LastName:     lastName,
		Site:         site,
		Role:         role,
		PhoneNumber:  phoneNumber,
		DateCreated:  t,
		DateModified: t,
		Data: AccountData{
			GOCDBAccess: false,
			SiteAccess:  false,
		},
		Settings: AccountSettings{
			ReceiveAlerts: true,
		},
	}

	// Set the user password, which also makes sure that the given password is strong enough
	if err := acc.UpdatePassword(password); err != nil {
		return nil, err
	}

	if err := acc.verify(true, true); err != nil {
		return nil, err
	}

	return acc, nil
}

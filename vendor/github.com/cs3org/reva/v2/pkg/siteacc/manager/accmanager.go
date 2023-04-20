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

package manager

import (
	"strings"
	"sync"
	"time"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/cs3org/reva/v2/pkg/siteacc/email"
	"github.com/cs3org/reva/v2/pkg/siteacc/manager/gocdb"
	"github.com/cs3org/reva/v2/pkg/smtpclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-password/password"
)

const (
	// FindByEmail holds the string value of the corresponding search criterium.
	FindByEmail = "email"
)

// AccountsManager is responsible for all site account related tasks.
type AccountsManager struct {
	conf *config.Configuration
	log  *zerolog.Logger

	storage data.Storage

	accounts          data.Accounts
	accountsListeners []AccountsListener

	smtp *smtpclient.SMTPCredentials

	mutex sync.RWMutex
}

func (mngr *AccountsManager) initialize(storage data.Storage, conf *config.Configuration, log *zerolog.Logger) error {
	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	mngr.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	mngr.log = log

	if storage == nil {
		return errors.Errorf("no storage provided")
	}
	mngr.storage = storage

	mngr.accounts = make(data.Accounts, 0, 32) // Reserve some space for accounts
	mngr.readAllAccounts()

	// Register accounts listeners
	if listener, err := gocdb.NewListener(mngr.conf, mngr.log); err == nil {
		mngr.accountsListeners = append(mngr.accountsListeners, listener)
	} else {
		return errors.Wrap(err, "unable to create the GOCDB accounts listener")
	}

	// Create the SMTP client
	if conf.Email.SMTP != nil {
		mngr.smtp = smtpclient.NewSMTPCredentials(conf.Email.SMTP)
	}

	return nil
}

func (mngr *AccountsManager) readAllAccounts() {
	if accounts, err := mngr.storage.ReadAccounts(); err == nil {
		mngr.accounts = *accounts
	} else {
		// Just warn when not being able to read accounts
		mngr.log.Warn().Err(err).Msg("error while reading accounts")
	}
}

func (mngr *AccountsManager) writeAllAccounts() {
	if err := mngr.storage.WriteAccounts(&mngr.accounts); err != nil {
		// Just warn when not being able to write accounts
		mngr.log.Warn().Err(err).Msg("error while writing accounts")
	}
}

func (mngr *AccountsManager) findAccount(by string, value string) (*data.Account, error) {
	if len(value) == 0 {
		return nil, errors.Errorf("no search value specified")
	}

	var account *data.Account
	switch strings.ToLower(by) {
	case FindByEmail:
		account = mngr.findAccountByPredicate(func(account *data.Account) bool { return strings.EqualFold(account.Email, value) })

	default:
		return nil, errors.Errorf("invalid search type %v", by)
	}

	if account != nil {
		return account, nil
	}

	return nil, errors.Errorf("no user found matching the specified criteria")
}

func (mngr *AccountsManager) findAccountByPredicate(predicate func(*data.Account) bool) *data.Account {
	for _, account := range mngr.accounts {
		if predicate(account) {
			return account
		}
	}
	return nil
}

// CreateAccount creates a new account; if an account with the same email address already exists, an error is returned.
func (mngr *AccountsManager) CreateAccount(accountData *data.Account) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	// Accounts must be unique (identified by their email address)
	if account, _ := mngr.findAccount(FindByEmail, accountData.Email); account != nil {
		return errors.Errorf("an account with the specified email address already exists")
	}

	if account, err := data.NewAccount(accountData.Email, accountData.Title, accountData.FirstName, accountData.LastName, accountData.Site, accountData.Role, accountData.PhoneNumber, accountData.Password.Value); err == nil {
		mngr.accounts = append(mngr.accounts, account)
		mngr.storage.AccountAdded(account)
		mngr.writeAllAccounts()

		mngr.sendEmail(account, nil, email.SendAccountCreated)
		mngr.callListeners(account, AccountsListener.AccountCreated)
	} else {
		return errors.Wrap(err, "error while creating account")
	}

	return nil
}

// UpdateAccount updates the account identified by the account email; if no such account exists, an error is returned.
func (mngr *AccountsManager) UpdateAccount(accountData *data.Account, setPassword bool, copyData bool) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	account, err := mngr.findAccount(FindByEmail, accountData.Email)
	if err != nil {
		return errors.Wrap(err, "user to update not found")
	}

	if err := account.Update(accountData, setPassword, copyData); err == nil {
		account.DateModified = time.Now()

		mngr.storage.AccountUpdated(account)
		mngr.writeAllAccounts()

		mngr.callListeners(account, AccountsListener.AccountUpdated)
	} else {
		return errors.Wrap(err, "error while updating account")
	}

	return nil
}

// ConfigureAccount configures the account identified by the account email; if no such account exists, an error is returned.
func (mngr *AccountsManager) ConfigureAccount(accountData *data.Account) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	account, err := mngr.findAccount(FindByEmail, accountData.Email)
	if err != nil {
		return errors.Wrap(err, "user to configure not found")
	}

	if err := account.Configure(accountData); err == nil {
		account.DateModified = time.Now()

		mngr.storage.AccountUpdated(account)
		mngr.writeAllAccounts()

		mngr.callListeners(account, AccountsListener.AccountUpdated)
	} else {
		return errors.Wrap(err, "error while configuring account")
	}

	return nil
}

// ResetPassword resets the password for the given user.
func (mngr *AccountsManager) ResetPassword(name string) error {
	account, err := mngr.findAccount(FindByEmail, name)
	if err != nil {
		return errors.Wrap(err, "user to reset password for not found")
	}
	accountUpd := account.Clone(true)
	accountUpd.Password.Value = password.MustGenerate(defaultPasswordLength, 2, 0, false, true)

	err = mngr.UpdateAccount(accountUpd, true, false)
	if err == nil {
		mngr.sendEmail(accountUpd, nil, email.SendPasswordReset)
	}

	return err
}

// FindAccount is used to find an account by various criteria. The account is cloned to prevent data changes.
func (mngr *AccountsManager) FindAccount(by string, value string) (*data.Account, error) {
	return mngr.FindAccountEx(by, value, true)
}

// FindAccountEx is used to find an account by various criteria and optionally clone the account.
func (mngr *AccountsManager) FindAccountEx(by string, value string, cloneAccount bool) (*data.Account, error) {
	mngr.mutex.RLock()
	defer mngr.mutex.RUnlock()

	account, err := mngr.findAccount(by, value)
	if err != nil {
		return nil, err
	}

	if cloneAccount {
		account = account.Clone(false)
	}

	return account, nil
}

// GrantSiteAccess sets the Site access status of the account identified by the account email; if no such account exists, an error is returned.
func (mngr *AccountsManager) GrantSiteAccess(accountData *data.Account, grantAccess bool) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	account, err := mngr.findAccount(FindByEmail, accountData.Email)
	if err != nil {
		return errors.Wrap(err, "no account with the specified email exists")
	}

	return mngr.grantAccess(account, &account.Data.SiteAccess, grantAccess, email.SendSiteAccessGranted)
}

// GrantGOCDBAccess sets the GOCDB access status of the account identified by the account email; if no such account exists, an error is returned.
func (mngr *AccountsManager) GrantGOCDBAccess(accountData *data.Account, grantAccess bool) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	account, err := mngr.findAccount(FindByEmail, accountData.Email)
	if err != nil {
		return errors.Wrap(err, "no account with the specified email exists")
	}

	return mngr.grantAccess(account, &account.Data.GOCDBAccess, grantAccess, email.SendGOCDBAccessGranted)
}

// RemoveAccount removes the account identified by the account email; if no such account exists, an error is returned.
func (mngr *AccountsManager) RemoveAccount(accountData *data.Account) error {
	mngr.mutex.Lock()
	defer mngr.mutex.Unlock()

	for i, account := range mngr.accounts {
		if strings.EqualFold(account.Email, accountData.Email) {
			mngr.accounts = append(mngr.accounts[:i], mngr.accounts[i+1:]...)
			mngr.storage.AccountRemoved(account)
			mngr.writeAllAccounts()

			mngr.callListeners(account, AccountsListener.AccountRemoved)
			return nil
		}
	}

	return errors.Errorf("no account with the specified email exists")
}

// SendContactForm sends a generic email to the ScienceMesh admins.
func (mngr *AccountsManager) SendContactForm(account *data.Account, subject, message string) {
	mngr.sendEmail(account, map[string]string{"Subject": subject, "Message": message}, email.SendContactForm)
}

// CloneAccounts retrieves all accounts currently stored by cloning the data, thus avoiding race conflicts and making outside modifications impossible.
func (mngr *AccountsManager) CloneAccounts(erasePasswords bool) data.Accounts {
	mngr.mutex.RLock()
	defer mngr.mutex.RUnlock()

	clones := make(data.Accounts, 0, len(mngr.accounts))
	for _, acc := range mngr.accounts {
		clones = append(clones, acc.Clone(erasePasswords))
	}

	return clones
}

func (mngr *AccountsManager) grantAccess(account *data.Account, accessFlag *bool, grantAccess bool, emailFunc email.SendFunction) error {
	accessOld := *accessFlag
	*accessFlag = grantAccess

	mngr.storage.AccountUpdated(account)
	mngr.writeAllAccounts()

	if *accessFlag && *accessFlag != accessOld {
		mngr.sendEmail(account, nil, emailFunc)
	}

	mngr.callListeners(account, AccountsListener.AccountUpdated)

	return nil
}

func (mngr *AccountsManager) callListeners(account *data.Account, cb AccountsListenerCallback) {
	for _, listener := range mngr.accountsListeners {
		cb(listener, account)
	}
}

func (mngr *AccountsManager) sendEmail(account *data.Account, params map[string]string, sendFunc email.SendFunction) {
	_ = sendFunc(account, []string{account.Email, mngr.conf.Email.NotificationsMail}, params, *mngr.conf)
}

// NewAccountsManager creates a new accounts manager instance.
func NewAccountsManager(storage data.Storage, conf *config.Configuration, log *zerolog.Logger) (*AccountsManager, error) {
	mngr := &AccountsManager{}
	if err := mngr.initialize(storage, conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the accounts manager")
	}
	return mngr, nil
}

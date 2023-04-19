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

package gocdb

import (
	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// AccountsListener is the GOCDB accounts listener.
type AccountsListener struct {
	conf *config.Configuration
	log  *zerolog.Logger
}

func (listener *AccountsListener) initialize(conf *config.Configuration, log *zerolog.Logger) error {
	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	listener.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	listener.log = log

	return nil
}

// AccountCreated is called whenever an account was created.
func (listener *AccountsListener) AccountCreated(account *data.Account) {
	listener.updateGOCDB(account, false)
}

// AccountUpdated is called whenever an account was updated.
func (listener *AccountsListener) AccountUpdated(account *data.Account) {
	listener.updateGOCDB(account, false)
}

// AccountRemoved is called whenever an account was removed.
func (listener *AccountsListener) AccountRemoved(account *data.Account) {
	listener.updateGOCDB(account, true)
}

func (listener *AccountsListener) updateGOCDB(account *data.Account, forceRemoval bool) {
	if account != nil && account.Data.GOCDBAccess && !forceRemoval {
		if err := writeAccount(account, opCreateOrUpdate, listener.conf.GOCDB.WriteURL, listener.conf.GOCDB.APIKey); err != nil {
			listener.log.Err(err).Str("userid", account.Email).Msg("unable to update GOCDB account")
		}
	} else {
		// Errors while deleting an account are ignored (account might not exist at all, for example)
		_ = writeAccount(account, opDelete, listener.conf.GOCDB.WriteURL, listener.conf.GOCDB.APIKey)
	}
}

// NewListener creates a new GOCDB accounts listener.
func NewListener(conf *config.Configuration, log *zerolog.Logger) (*AccountsListener, error) {
	listener := &AccountsListener{}
	if err := listener.initialize(conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the GOCDB accounts listener")
	}
	return listener, nil
}

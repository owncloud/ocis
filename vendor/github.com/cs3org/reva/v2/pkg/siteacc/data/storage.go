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

// Storage defines the interface for sites and accounts storages.
type Storage interface {
	// ReadSites reads all stored sites into the given data object.
	ReadSites() (*Sites, error)
	// WriteSites writes all stored sites from the given data object.
	WriteSites(sites *Sites) error

	// SiteAdded is called when a site has been added.
	SiteAdded(site *Site)
	// SiteUpdated is called when a site has been updated.
	SiteUpdated(site *Site)
	// SiteRemoved is called when a site has been removed.
	SiteRemoved(site *Site)

	// ReadAccounts reads all stored accounts into the given data object.
	ReadAccounts() (*Accounts, error)
	// WriteAccounts writes all stored accounts from the given data object.
	WriteAccounts(accounts *Accounts) error

	// AccountAdded is called when an account has been added.
	AccountAdded(account *Account)
	// AccountUpdated is called when an account has been updated.
	AccountUpdated(account *Account)
	// AccountRemoved is called when an account has been removed.
	AccountRemoved(account *Account)
}

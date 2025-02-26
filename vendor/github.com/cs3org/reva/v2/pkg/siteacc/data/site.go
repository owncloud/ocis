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

package data

import (
	"github.com/cs3org/reva/v2/pkg/siteacc/credentials"
	"github.com/pkg/errors"
)

// Site represents the global site-specific settings stored in the service.
type Site struct {
	ID string `json:"id"`

	Config SiteConfiguration `json:"config"`
}

// SiteConfiguration stores the global configuration of a site.
type SiteConfiguration struct {
	TestClientCredentials credentials.Credentials `json:"testClientCredentials"`
}

// Sites holds an array of sites.
type Sites = []*Site

// Update copies the data of the given site to this site.
func (site *Site) Update(other *Site, credsPassphrase string) error {
	if other.Config.TestClientCredentials.IsValid() {
		// If credentials were provided, use those as the new ones
		if err := site.UpdateTestClientCredentials(other.Config.TestClientCredentials.ID, other.Config.TestClientCredentials.Secret, credsPassphrase); err != nil {
			return err
		}
	}

	return nil
}

// UpdateTestClientCredentials assigns new test client credentials, encrypting the information first.
func (site *Site) UpdateTestClientCredentials(id, secret string, passphrase string) error {
	if err := site.Config.TestClientCredentials.Set(id, secret, passphrase); err != nil {
		return errors.Wrap(err, "unable to update the test client credentials")
	}
	return nil
}

// Clone creates a copy of the site; if eraseCredentials is set to true, the (test user) credentials will be cleared in the cloned object.
func (site *Site) Clone(eraseCredentials bool) *Site {
	clone := *site

	if eraseCredentials {
		clone.Config.TestClientCredentials.Clear()
	}

	return &clone
}

// NewSite creates a new site.
func NewSite(id string) (*Site, error) {
	site := &Site{
		ID: id,
		Config: SiteConfiguration{
			TestClientCredentials: credentials.Credentials{},
		},
	}
	return site, nil
}

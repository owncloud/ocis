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

package config

import (
	"strings"

	"github.com/cs3org/reva/v2/pkg/smtpclient"
)

// Configuration holds the general service configuration.
type Configuration struct {
	Prefix string `mapstructure:"prefix"`

	Security struct {
		CredentialsPassphrase string `mapstructure:"creds_passphrase"`
	} `mapstructure:"security"`

	Storage struct {
		Driver string `mapstructure:"driver"`

		File struct {
			SitesFile    string `mapstructure:"sites_file"`
			AccountsFile string `mapstructure:"accounts_file"`
		} `mapstructure:"file"`
	} `mapstructure:"storage"`

	Email struct {
		SMTP              *smtpclient.SMTPCredentials `mapstructure:"smtp"`
		NotificationsMail string                      `mapstructure:"notifications_mail"`
	} `mapstructure:"email"`

	Mentix struct {
		URL                      string `mapstructure:"url"`
		DataEndpoint             string `mapstructure:"data_endpoint"`
		SiteRegistrationEndpoint string `mapstructure:"sitereg_endpoint"`
	} `mapstructure:"mentix"`

	Webserver struct {
		URL string `mapstructure:"url"`

		SessionTimeout      int  `mapstructure:"session_timeout"`
		VerifyRemoteAddress bool `mapstructure:"verify_remote_address"`
		LogSessions         bool `mapstructure:"log_sessions"`
	} `mapstructure:"webserver"`

	GOCDB struct {
		URL      string `mapstructure:"url"`
		WriteURL string `mapstructure:"write_url"`

		APIKey string `mapstructure:"apikey"`
	} `mapstructure:"gocdb"`
}

// Cleanup cleans up certain settings, normalizing them.
func (cfg *Configuration) Cleanup() {
	// Ensure the webserver URL ends with a slash
	if cfg.Webserver.URL != "" && !strings.HasSuffix(cfg.Webserver.URL, "/") {
		cfg.Webserver.URL += "/"
	}

	// Ensure the GOCDB URL ends with a slash
	if cfg.GOCDB.URL != "" && !strings.HasSuffix(cfg.GOCDB.URL, "/") {
		cfg.GOCDB.URL += "/"
	}

	// Ensure the GOCDB Write URL ends with a slash
	if cfg.GOCDB.WriteURL != "" && !strings.HasSuffix(cfg.GOCDB.WriteURL, "/") {
		cfg.GOCDB.WriteURL += "/"
	}
}

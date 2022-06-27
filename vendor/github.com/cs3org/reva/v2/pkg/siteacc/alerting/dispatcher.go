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

package alerting

import (
	"strings"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/cs3org/reva/v2/pkg/siteacc/email"
	"github.com/cs3org/reva/v2/pkg/smtpclient"
	"github.com/pkg/errors"
	"github.com/prometheus/alertmanager/template"
	"github.com/rs/zerolog"
)

// Dispatcher is used to dispatch Prometheus alerts via email.
type Dispatcher struct {
	conf *config.Configuration
	log  *zerolog.Logger

	smtp *smtpclient.SMTPCredentials
}

func (dispatcher *Dispatcher) initialize(conf *config.Configuration, log *zerolog.Logger) error {
	if conf == nil {
		return errors.Errorf("no configuration provided")
	}
	dispatcher.conf = conf

	if log == nil {
		return errors.Errorf("no logger provided")
	}
	dispatcher.log = log

	// Create the SMTP client
	if conf.Email.SMTP != nil {
		dispatcher.smtp = smtpclient.NewSMTPCredentials(conf.Email.SMTP)
	}

	return nil
}

// DispatchAlerts sends the provided alert(s) via email to the appropriate recipients.
func (dispatcher *Dispatcher) DispatchAlerts(alerts *template.Data, accounts data.Accounts) error {
	for _, alert := range alerts.Alerts {
		siteID, ok := alert.Labels["site_id"]
		if !ok {
			continue
		}

		// Dispatch the alert to all accounts configured to receive it
		for _, account := range accounts {
			if strings.EqualFold(account.Site, siteID) /* && account.Settings.ReceiveAlerts */ { // TODO: Uncomment if alert notifications aren't mandatory anymore
				if err := dispatcher.dispatchAlert(alert, account); err != nil {
					// Log errors only
					dispatcher.log.Err(err).Str("id", alert.Fingerprint).Str("recipient", account.Email).Msg("unable to dispatch alert to user")
				}
			}
		}

		// Dispatch the alert to the global receiver (if set)
		if dispatcher.conf.Email.NotificationsMail != "" {
			globalAccount := data.Account{ // On-the-fly account representing the "global alerts receiver"
				Email:     dispatcher.conf.Email.NotificationsMail,
				FirstName: "ScienceMesh",
				LastName:  "Global Alerts receiver",
				Site:      "Global",
				Role:      "Alerts receiver",
				Settings: data.AccountSettings{
					ReceiveAlerts: true,
				},
			}
			if err := dispatcher.dispatchAlert(alert, &globalAccount); err != nil {
				dispatcher.log.Err(err).Str("id", alert.Fingerprint).Str("recipient", globalAccount.Email).Msg("unable to dispatch alert to global alerts receiver")
			}
		}
	}
	return nil
}

func (dispatcher *Dispatcher) dispatchAlert(alert template.Alert, account *data.Account) error {
	alertValues := map[string]string{
		"Status":      alert.Status,
		"StartDate":   alert.StartsAt.String(),
		"EndDate":     alert.EndsAt.String(),
		"Fingerprint": alert.Fingerprint,

		"Name":     alert.Labels["alertname"],
		"Service":  alert.Labels["service_type"],
		"Instance": alert.Labels["instance"],
		"Job":      alert.Labels["job"],
		"Severity": alert.Labels["severity"],
		"Site":     alert.Labels["site"],
		"SiteID":   alert.Labels["site_id"],

		"Description": alert.Annotations["description"],
		"Summary":     alert.Annotations["summary"],
	}

	return email.SendAlertNotification(account, []string{account.Email}, alertValues, *dispatcher.conf)
}

// NewDispatcher creates a new dispatcher instance.
func NewDispatcher(conf *config.Configuration, log *zerolog.Logger) (*Dispatcher, error) {
	dispatcher := &Dispatcher{}
	if err := dispatcher.initialize(conf, log); err != nil {
		return nil, errors.Wrap(err, "unable to initialize the alerts dispatcher")
	}
	return dispatcher, nil
}

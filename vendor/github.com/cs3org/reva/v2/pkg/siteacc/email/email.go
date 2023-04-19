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

package email

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/cs3org/reva/v2/pkg/siteacc/config"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/cs3org/reva/v2/pkg/smtpclient"
	"github.com/pkg/errors"
)

type emailData struct {
	Account *data.Account

	AccountsAddress string
	GOCDBAddress    string

	Params map[string]string
}

// SendFunction is the definition of email send functions.
type SendFunction = func(*data.Account, []string, map[string]string, config.Configuration) error

func getEmailData(account *data.Account, conf config.Configuration, params map[string]string) *emailData {
	return &emailData{
		Account:         account,
		AccountsAddress: conf.Webserver.URL,
		GOCDBAddress:    conf.GOCDB.URL,
		Params:          params,
	}
}

// SendAccountCreated sends an email about account creation.
func SendAccountCreated(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	return send(recipients, "ScienceMesh: Site Administrator Account created", accountCreatedTemplate, getEmailData(account, conf, params), conf.Email.SMTP)
}

// SendSiteAccessGranted sends an email about granted Site access.
func SendSiteAccessGranted(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	return send(recipients, "ScienceMesh: Site access granted", siteAccessGrantedTemplate, getEmailData(account, conf, params), conf.Email.SMTP)
}

// SendGOCDBAccessGranted sends an email about granted GOCDB access.
func SendGOCDBAccessGranted(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	return send(recipients, "ScienceMesh: GOCDB access granted", gocdbAccessGrantedTemplate, getEmailData(account, conf, params), conf.Email.SMTP)
}

// SendPasswordReset sends an email containing the user's new password.
func SendPasswordReset(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	return send(recipients, "ScienceMesh: Password reset", passwordResetTemplate, getEmailData(account, conf, params), conf.Email.SMTP)
}

// SendContactForm sends a generic contact form to the ScienceMesh admins.
func SendContactForm(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	return send(recipients, "ScienceMesh: Contact form", contactFormTemplate, getEmailData(account, conf, params), conf.Email.SMTP)
}

// SendAlertNotification sends an alert via email.
func SendAlertNotification(account *data.Account, recipients []string, params map[string]string, conf config.Configuration) error {
	subject := params["Summary"]
	tpl := alertFiringNotificationTemplate
	if strings.EqualFold(params["Status"], "resolved") {
		tpl = alertResolvedNotificationTemplate
		subject += " [RESOLVED]"
	}
	return send(recipients, "ScienceMesh Alert: "+subject, tpl, getEmailData(account, conf, params), conf.Email.SMTP)
}

func send(recipients []string, subject string, bodyTemplate string, data interface{}, smtp *smtpclient.SMTPCredentials) error {
	// Do not fail if no SMTP client or recipient is given
	if smtp == nil {
		return nil
	}

	tpl := template.New("email")
	prepareEmailTemplate(tpl)

	if _, err := tpl.Parse(bodyTemplate); err != nil {
		return errors.Wrap(err, "error while parsing email template")
	}

	var body bytes.Buffer
	if err := tpl.Execute(&body, data); err != nil {
		return errors.Wrap(err, "error while executing email template")
	}

	for _, recipient := range recipients {
		if len(recipient) == 0 {
			continue
		}

		// Send the mail w/o blocking the main thread
		go func(recipient string) {
			_ = smtp.SendMail(recipient, subject, body.String())
		}(recipient)
	}

	return nil
}

func prepareEmailTemplate(tpl *template.Template) {
	// Add some custom helper functions to the template
	tpl.Funcs(template.FuncMap{
		"indent": func(n int, s string) string {
			lines := make([]string, 0, 10)
			for _, line := range strings.Split(s, "\n") {
				line = strings.TrimSpace(line)
				line = strings.Repeat(" ", n) + line
				lines = append(lines, line)
			}
			return strings.Join(lines, "\n")
		},
	})
}

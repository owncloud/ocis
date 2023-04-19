// Copyright 2018-2021 CERN
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

package smtpclient

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// SMTPCredentials stores the credentials required to connect to an SMTP server.
type SMTPCredentials struct {
	SenderLogin    string `mapstructure:"sender_login" docs:";The login to be used by sender."`
	SenderMail     string `mapstructure:"sender_mail" docs:";The email to be used to send mails."`
	SenderPassword string `mapstructure:"sender_password" docs:";The sender's password."`
	SMTPServer     string `mapstructure:"smtp_server" docs:";The hostname of the SMTP server."`
	SMTPPort       int    `mapstructure:"smtp_port" docs:"587;The port on which the SMTP daemon is running."`
	DisableAuth    bool   `mapstructure:"disable_auth" docs:"false;Whether to disable SMTP auth."`
	LocalName      string `mapstructure:"local_name" docs:";The host name to be used for unauthenticated SMTP."`
}

// NewSMTPCredentials creates a new SMTPCredentials object with the details of the passed object with sane defaults.
func NewSMTPCredentials(c *SMTPCredentials) *SMTPCredentials {
	creds := c

	if creds.SMTPPort == 0 {
		creds.SMTPPort = 587
	}
	if !creds.DisableAuth && creds.SenderPassword == "" {
		creds.SenderPassword = os.Getenv("REVA_SMTP_SENDER_PASSWORD")
	}
	if creds.LocalName == "" {
		tokens := strings.Split(creds.SenderMail, "@")
		creds.LocalName = tokens[len(tokens)-1]
	}
	if creds.SenderLogin == "" {
		creds.SenderLogin = creds.SenderMail
	}
	return creds
}

// SendMail allows sending mails using a set of client credentials.
func (creds *SMTPCredentials) SendMail(recipient, subject, body string) error {

	headers := map[string]string{
		"From":                      creds.SenderMail,
		"To":                        recipient,
		"Subject":                   subject,
		"Date":                      time.Now().Format(time.RFC1123Z),
		"Message-ID":                uuid.New().String(),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding": "base64",
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	if creds.DisableAuth {
		return creds.sendMailSMTP(recipient, subject, message)
	}
	return creds.sendMailAuthSMTP(recipient, subject, message)
}

func (creds *SMTPCredentials) sendMailAuthSMTP(recipient, subject, message string) error {

	auth := smtp.PlainAuth("", creds.SenderLogin, creds.SenderPassword, creds.SMTPServer)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", creds.SMTPServer, creds.SMTPPort),
		auth,
		creds.SenderMail,
		[]string{recipient},
		[]byte(message),
	)
	if err != nil {
		err = errors.Wrap(err, "smtpclient: error sending mail")
		return err
	}

	return nil
}

func (creds *SMTPCredentials) sendMailSMTP(recipient, subject, message string) error {

	c, err := smtp.Dial(fmt.Sprintf("%s:%d", creds.SMTPServer, creds.SMTPPort))
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Hello(creds.LocalName); err != nil {
		return err
	}
	if err = c.Mail(creds.SenderMail); err != nil {
		return err
	}
	if err = c.Rcpt(recipient); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	buf := bytes.NewBufferString(message)
	if _, err = buf.WriteTo(wc); err != nil {
		return err
	}

	return nil
}

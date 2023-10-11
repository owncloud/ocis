// Copyright 2018-2023 CERN
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

package sciencemesh

import (
	"bytes"
	"html/template"
	"io"
	"os"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

type emailParams struct {
	User             *userpb.User
	Token            string
	MeshDirectoryURL string
	InviteLink       string
}

const defaultSubject = `ScienceMesh: {{.User.DisplayName}} wants to collaborate with you`

const defaultBody = `Hi

{{.User.DisplayName}} ({{.User.Mail}}) wants to start sharing OCM resources with you.
To accept the invite, please visit the following URL:
{{.InviteLink}}

Alternatively, you can visit your mesh provider and use the following details:
Token: {{.Token}}
ProviderDomain: {{.User.Id.Idp}}

Best,
The ScienceMesh team`

func (h *tokenHandler) sendEmail(recipient string, obj *emailParams) error {
	subj, err := h.generateEmailSubject(obj)
	if err != nil {
		return err
	}

	body, err := h.generateEmailBody(obj)
	if err != nil {
		return err
	}

	return h.smtpCredentials.SendMail(recipient, subj, body)
}

func (h *tokenHandler) generateEmailSubject(obj *emailParams) (string, error) {
	var buf bytes.Buffer
	err := h.tplSubj.Execute(&buf, obj)
	return buf.String(), err
}

func (h *tokenHandler) generateEmailBody(obj *emailParams) (string, error) {
	var buf bytes.Buffer
	err := h.tplBody.Execute(&buf, obj)
	return buf.String(), err
}

func (h *tokenHandler) initBodyTemplate(bodyTemplPath string) error {
	var body string
	if bodyTemplPath == "" {
		body = defaultBody
	} else {
		f, err := os.Open(bodyTemplPath)
		if err != nil {
			return err
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		body = string(data)
	}

	tpl, err := template.New("tpl_body").Parse(body)
	if err != nil {
		return err
	}

	h.tplBody = tpl
	return nil
}

func (h *tokenHandler) initSubjectTemplate(subjTempl string) error {
	var subj string
	if subjTempl == "" {
		subj = defaultSubject
	} else {
		subj = subjTempl
	}

	tpl, err := template.New("tpl_subj").Parse(subj)
	if err != nil {
		return err
	}
	h.tplSubj = tpl
	return nil
}

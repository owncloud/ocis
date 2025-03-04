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

const accountCreatedTemplate = `
Dear {{.Account.FirstName}} {{.Account.LastName}},

Your ScienceMesh Site Administrator Account has been successfully created!

Log in to your account by visiting the user account panel:
{{.AccountsAddress}}

Using this panel, you can modify your information, request access to the GOCDB, and more. 

Kind regards,
The ScienceMesh Team
`

const siteAccessGrantedTemplate = `
Dear {{.Account.FirstName}} {{.Account.LastName}},

You have been granted access to the global configuration of your site.

Log in to your account to access this configuration:
{{.AccountsAddress}} 

Kind regards,
The ScienceMesh Team
`

const gocdbAccessGrantedTemplate = `
Dear {{.Account.FirstName}} {{.Account.LastName}},

You have been granted access to the ScienceMesh GOCDB instance:
{{.GOCDBAddress}}

Simply use your regular ScienceMesh Site Administrator Account credentials to log in to the GOCDB. 

Kind regards,
The ScienceMesh Team
`

const passwordResetTemplate = `
Dear {{.Account.FirstName}} {{.Account.LastName}},

Your password has been successfully reset!
Your new password is: {{.Account.Password.Value}}

We recommend to change this password immediately after logging in.

Kind regards,
The ScienceMesh Team
`

const contactFormTemplate = `
{{.Account.FirstName}} {{.Account.LastName}} ({{.Account.Email}}) has sent the following message:

{{.Params.Subject}}
---------------------------------------------------------------------------------------------------

{{.Params.Message}}
`

const alertFiringNotificationTemplate = `
Site '{{.Params.Site}}' has generated an alert:

  Type:     {{.Params.Name}}
  Service:  {{.Params.Service}}
  Instance: {{.Params.Instance}}
  Job:      {{.Params.Job}}
  Severity: {{.Params.Severity}}

{{.Params.Description | indent 2}}

{{.Params.StartDate}} ({{.Params.Fingerprint}})
`

const alertResolvedNotificationTemplate = `
Site '{{.Params.Site}}' has resolved an alert:

  Type:     {{.Params.Name}}
  Service:  {{.Params.Service}}
  Instance: {{.Params.Instance}}
  Job:      {{.Params.Job}}
  Severity: {{.Params.Severity}}

{{.Params.Description | indent 2}}

{{.Params.StartDate}} - {{.Params.EndDate}} ({{.Params.Fingerprint}})
`

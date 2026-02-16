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

const (
	// EndpointAdministration is the endpoint path of the web interface administration panel.
	EndpointAdministration = "/admin"
	// EndpointAccount is the endpoint path of the web interface account panel.
	EndpointAccount = "/account"

	// EndpointList is the endpoint path for listing all stored accounts.
	EndpointList = "/list"
	// EndpointFind is the endpoint path for finding accounts.
	EndpointFind = "/find"

	// EndpointCreate is the endpoint path for account creation.
	EndpointCreate = "/create"
	// EndpointUpdate is the endpoint path for account updates.
	EndpointUpdate = "/update"
	// EndpointConfigure is the endpoint path for account configuration.
	EndpointConfigure = "/configure"
	// EndpointRemove is the endpoint path for account removal.
	EndpointRemove = "/remove"

	// EndpointSiteGet is the endpoint path for retrieving site data.
	EndpointSiteGet = "/site-get"
	// EndpointSiteConfigure is the endpoint path for site configuration.
	EndpointSiteConfigure = "/site-configure"

	// EndpointLogin is the endpoint path for (internal) user login.
	EndpointLogin = "/login"
	// EndpointLogout is the endpoint path for (internal) user logout.
	EndpointLogout = "/logout"
	// EndpointResetPassword is the endpoint path for resetting user passwords
	EndpointResetPassword = "/reset-password"
	// EndpointContact is the endpoint path for sending contact emails
	EndpointContact = "/contact"

	// EndpointVerifyUserToken is the endpoint path for user token validation.
	EndpointVerifyUserToken = "/verify-user-token"

	// EndpointGrantSiteAccess is the endpoint path for granting or revoking Site access.
	EndpointGrantSiteAccess = "/grant-site-access"
	// EndpointGrantGOCDBAccess is the endpoint path for granting or revoking GOCDB access.
	EndpointGrantGOCDBAccess = "/grant-gocdb-access"

	// EndpointDispatchAlert is the endpoint path for dispatching alerts from Prometheus.
	EndpointDispatchAlert = "/dispatch-alert"
)

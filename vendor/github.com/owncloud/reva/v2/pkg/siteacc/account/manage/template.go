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

package manage

const tplJavaScript = `
function handleAccountSettings() {
	setState(STATE_STATUS, "Redirecting to the account settings...");
	window.location.replace("{{getServerAddress}}/account/?path=settings");
}

function handleEditAccount() {
	setState(STATE_STATUS, "Redirecting to the account editor...");
	window.location.replace("{{getServerAddress}}/account/?path=edit");
}

function handleSiteSettings() {
	setState(STATE_STATUS, "Redirecting to the site settings...");
	window.location.replace("{{getServerAddress}}/account/?path=site");
}

function handleRequestAccess(scope) {
	setState(STATE_STATUS, "Redirecting to the contact form...");		
	window.location.replace("{{getServerAddress}}/account/?path=contact&subject=" + encodeURIComponent("Request " + scope + " access"));
}

function handleLogout() {
	var xhr = new XMLHttpRequest();
    xhr.open("GET", "{{getServerAddress}}/logout");
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	setState(STATE_STATUS, "Logging out...");

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Done! Redirecting...");
			window.location.replace("{{getServerAddress}}/account/?path=login");
		} else {
			setState(STATE_ERROR, "An error occurred while logging out: " + this.responseText);
		}
	}
    
    xhr.send();
}
`

const tplStyleSheet = `
html * {
	font-family: arial !important;
}
button {
	min-width: 170px;
}
`

const tplBody = `
<div>
	<p><strong>Hello {{.Account.FirstName}} {{.Account.LastName}},</strong></p>
	<p>On this page, you can manage your ScienceMesh Site Administrator Account. This includes editing your personal information, requesting access to the GOCDB and more.</p>
</div>
<div>&nbsp;</div>
<div>
	<strong>Personal information:</strong>
	<ul style="margin-top: 0em;">
		<li>Name: <em>{{.Account.Title}}. {{.Account.FirstName}} {{.Account.LastName}}</em></li>
		<li>Email: <em><a href="mailto:{{.Account.Email}}">{{.Account.Email}}</a></em></li>
		<li>ScienceMesh Site: <em>{{getSiteName .Account.Site false}} ({{getSiteName .Account.Site true}})</em></li>
		<li>Role: <em>{{.Account.Role}}</em></li>
		{{if .Account.PhoneNumber}}
		<li>Phone: <em>{{.Account.PhoneNumber}}</em></li>
		{{end}}
	</ul>
</div>
<div>
	<strong>Account data:</strong>
	<ul style="margin-top: 0em;">	
		<li>Site access: <em>{{if .Account.Data.SiteAccess}}Granted{{else}}Not granted{{end}}</em></li>
		<li>GOCDB access: <em>{{if .Account.Data.GOCDBAccess}}Granted{{else}}Not granted{{end}}</em></li>	
	</ul>
</div>
<div>
	<form id="form" method="POST" class="box" style="width: 100%;">
		<div>
			<button type="button" onClick="handleAccountSettings();">Account settings</button>
			<button type="button" onClick="handleEditAccount();">Edit account</button>
			<span style="width: 25px;">&nbsp;</span>
			
			{{if .Account.Data.SiteAccess}}
			<button type="button" onClick="handleSiteSettings();">Site settings</button>
			<span style="width: 25px;">&nbsp;</span>
			{{end}}	

			<button type="button" onClick="handleLogout();" style="float: right;">Logout</button>
		</div>
		<div style="margin-top: 0.5em;">
			<button type="button" onClick="handleRequestAccess('Site');" {{if .Account.Data.SiteAccess}}disabled{{end}}>Request Site access</button>
			<button type="button" onClick="handleRequestAccess('GOCDB');" {{if .Account.Data.GOCDBAccess}}disabled{{end}}>Request GOCDB access</button>	
		</div>
	</form>
</div>
<div style="font-size: 90%; margin-top: 1em;">
	<div>
		<div>Notes:</div>
		<ul style="margin-top: 0em;">
			<li>The <em>Site access</em> allows you to access and modify the global configuration of your site.</li>
			<li>The <em>GOCDB access</em> allows you to log into the central database where all site metadata is stored.</li>
		</ul>
	</div>
	<div>
		<div>Quick links:</div>
		<ul style="margin-top: 0em;">
			<li><a href="https://gocdb.sciencemesh.uni-muenster.de" target="_blank">Central Database (GOCDB)</a></li>
			<li><a href="https://developer.sciencemesh.io/docs/technical-documentation/central-database/" target="_blank">Central Database documentation</a></li>
		</ul>
	</div>
</div>
`

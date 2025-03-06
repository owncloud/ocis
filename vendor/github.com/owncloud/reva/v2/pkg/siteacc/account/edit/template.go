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

package edit

const tplJavaScript = `
function verifyForm(formData) {
	if (formData.getTrimmed("fname") == "") {
		setState(STATE_ERROR, "Please specify your first name.", "form", "fname", true);
		return false;
	}

	if (formData.getTrimmed("lname") == "") {
		setState(STATE_ERROR, "Please specify your last name.", "form", "lname", true);	
		return false;
	}

	if (formData.getTrimmed("role") == "") {
		setState(STATE_ERROR, "Please specify your role within your site.", "form", "role", true);
		return false;
	}

	if (formData.get("password") != "") {
		if (formData.get("password2") == "") {
			setState(STATE_ERROR, "Please confirm your new password.", "form", "password2", true);
			return false;
		}
	
		if (formData.get("password") != formData.get("password2")) {
			setState(STATE_ERROR, "The entered passwords do not match.", "form", "password2", true);
			return false;
		}
	}

	return true;
}

function handleAction(action) {
	const formData = new FormData(document.querySelector("form"));
	if (!verifyForm(formData)) {
		return;
	}

	setState(STATE_STATUS, "Updating account... this should only take a moment.", "form", null, false);

	var xhr = new XMLHttpRequest();
    xhr.open("POST", "{{getServerAddress}}/" + action);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Your account was successfully updated!", "form", null, true);
		} else {
			var resp = JSON.parse(this.responseText);
			setState(STATE_ERROR, "An error occurred while trying to update your account:<br><em>" + resp.error + "</em>", "form", null, true);
		}
	}

	var postData = {
		"title": formData.getTrimmed("title"),
		"firstName": formData.getTrimmed("fname"),
		"lastName": formData.getTrimmed("lname"),
		"role": formData.getTrimmed("role"),
		"phoneNumber": formData.getTrimmed("phone"),
		"password": {
			"value": formData.get("password")
		}
    };

    xhr.send(JSON.stringify(postData));
}
`

const tplStyleSheet = `
html * {
	font-family: arial !important;
}

.mandatory {
	color: red;
	font-weight: bold;
}
`

const tplBody = `
<div>
	<p>Edit your ScienceMesh Site Administrator Account information below.</p>
	<p>Please note that you cannot modify your email address using this form.</p>
</div>
<div>&nbsp;</div>
<div>
	<form id="form" method="POST" class="box container-inline" style="width: 100%;"  onSubmit="handleAction('update?invoker=user'); return false;">
		<div style="grid-row: 1;"><label for="title">Title: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 2;">
			<select id="title" name="title">
			{{$title := .Account.Title}}
			{{range .Titles}}
			<option value="{{.}}" {{if eq . $title}}selected{{end}}>{{.}}.</option>
			{{end}}
			</select>
		</div>

		<div style="grid-row: 3;"><label for="fname">First name: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 4;"><input type="text" id="fname" name="fname" value="{{.Account.FirstName}}"/></div>
		<div style="grid-row: 3;"><label for="lname">Last name: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 4;"><input type="text" id="lname" name="lname" value="{{.Account.LastName}}"/></div>

		<div style="grid-row: 5;"><label for="role">Role: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 6;"><input type="text" id="role" name="role" placeholder="Site administrator" value="{{.Account.Role}}"/></div>
		<div style="grid-row: 5;"><label for="phone">Phone number:</label></div>
		<div style="grid-row: 6;"><input type="text" id="phone" name="phone" placeholder="+49 030 123456" value="{{.Account.PhoneNumber}}"/></div>

		<div style="grid-row: 7;">&nbsp;</div>

		<div style="grid-row: 8; grid-column: 1 / span 2;">If you want to change your password, fill out the fields below. Otherwise, leave them empty to keep your current one.</div>
		<div style="grid-row: 9;"><label for="password">New password:</label></div>
		<div style="grid-row: 10;"><input type="password" id="password" name="password" autocomplete="new-password"/></div>
		<div style="grid-row: 9"><label for="password2">Confirm new password:</label></div>
		<div style="grid-row: 10;"><input type="password" id="password2" name="password2" autocomplete="new-password"/></div>

		<div style="grid-row: 11; font-style: italic; font-size: 0.8em;">
			The password must fulfil the following criteria:
			<ul style="margin-top: 0em;">
				<li>Must be at least 8 characters long</li>
				<li>Must contain at least 1 lowercase letter</li>
				<li>Must contain at least 1 uppercase letter</li>
				<li>Must contain at least 1 digit</li>
			</ul>
		</div>

		<div style="grid-row: 12; align-self: center;">
			Fields marked with <span class="mandatory">*</span> are mandatory.
		</div>
		<div style="grid-row: 12; grid-column: 2; text-align: right;">
			<button type="reset">Reset</button>
			<button type="submit" style="font-weight: bold;">Save</button>
		</div>
	</form>
</div>
<div>
	<p>Go <a href="{{getServerAddress}}/account/?path=manage">back</a> to the main account page.</p>
</div>
`

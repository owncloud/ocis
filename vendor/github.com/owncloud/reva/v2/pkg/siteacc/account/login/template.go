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

package login

const tplJavaScript = `
function verifyForm(formData, requirePassword = true) {
	if (formData.getTrimmed("email") == "") {
		setState(STATE_ERROR, "Please enter your email address.", "form", "email", true);
		return false;
	}

	if (requirePassword) {
		if (formData.get("password") == "") {
			setState(STATE_ERROR, "Please enter your password.", "form", "password", true);
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

	setState(STATE_STATUS, "Logging in... this should only take a moment.", "form", null, false);

	var xhr = new XMLHttpRequest();
    xhr.open("POST", "{{getServerAddress}}/" + action);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Your login was successful! Redirecting...");
			window.location.replace("{{getServerAddress}}/account/?path=manage");
		} else {
			var resp = JSON.parse(this.responseText);
			setState(STATE_ERROR, "An error occurred while trying to login your account:<br><em>" + resp.error + "</em>", "form", null, true);
		}
	}

	var postData = {
        "email": formData.getTrimmed("email"),
		"password": {
			"value": formData.get("password")
		}
    };

    xhr.send(JSON.stringify(postData));
}

function handleResetPassword() {
	const formData = new FormData(document.querySelector("form"));
	if (!verifyForm(formData, false)) {
		return;
	}

	setState(STATE_STATUS, "Resetting password... this should only take a moment.", "form", null, false);

	var xhr = new XMLHttpRequest();
    xhr.open("POST", "{{getServerAddress}}/reset-password");
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Your password was successfully reset! Please check your inbox for your new password.", "form", null, true);
		} else {
			var resp = JSON.parse(this.responseText);
			setState(STATE_ERROR, "An error occurred while trying to reset your password:<br><em>" + resp.error + "</em>", "form", null, true);
		}
	}

	var postData = {
        "email": formData.get("email")
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
	<p>Login to your ScienceMesh Site Administrator Account using the form below.</p>
</div>
<div>&nbsp;</div>
<div>
	<form id="form" method="POST" class="box container-inline" style="width: 100%;" onSubmit="handleAction('login'); return false;">
		<div style="grid-row: 1;"><label for="email">Email address: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 2;"><input type="text" id="email" name="email" placeholder="me@example.com"/></div>
		<div style="grid-row: 1;"><label for="password">Password: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 2;"><input type="password" id="password" name="password"/></div>
		<div style="grid-row: 3; grid-column: 2; font-style: italic; font-size: 0.8em;">
			Forgot your password? Click <a href="#" onClick="handleResetPassword();">here</a> to reset it.
		</div>

		<div style="grid-row: 4; align-self: center;">
			Fields marked with <span class="mandatory">*</span> are mandatory.
		</div>
		<div style="grid-row: 4; grid-column: 2; text-align: right;">
			<button type="reset">Reset</button>
			<button type="submit" style="font-weight: bold;">Login</button>
		</div>	
	</form>	
</div>
<div>
	<p>Don't' have an account yet? Register <a href="{{getServerAddress}}/account/?path=register">here</a>.</p>
</div>
`

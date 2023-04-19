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

package settings

const tplJavaScript = `
function verifyForm(formData) {
	return true;
}

function handleAction(action) {
	const formData = new FormData(document.querySelector("form"));
	if (!verifyForm(formData)) {
		return;
	}

	setState(STATE_STATUS, "Configuring account... this should only take a moment.", "form", null, false);

	var xhr = new XMLHttpRequest();
    xhr.open("POST", "{{getServerAddress}}/" + action);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Your account was successfully configured!", "form", null, true);
		} else {
			var resp = JSON.parse(this.responseText);
			setState(STATE_ERROR, "An error occurred while trying to configure your account:<br><em>" + resp.error + "</em>", "form", null, true);
		}
	}

	var postData = {
		"settings": {
			"receiveAlerts": (formData.get("rcvAlerts") === "on")
		}
    };

    xhr.send(JSON.stringify(postData));
}
`

const tplStyleSheet = `
html * {
	font-family: arial !important;
}

input[type="checkbox"] {
	width: auto;
}
`

const tplBody = `
<div>
	<p>Configure your ScienceMesh Site Administrator Account below.</p>	
</div>
<div>&nbsp;</div>
<div>
	<form id="form" method="POST" class="box container-inline" style="width: 100%;" onSubmit="handleAction('configure?invoker=user'); return false;">
		<div style="grid-row: 1; grid-column: 1 / span 2;">
			<h3>Notification settings</h3>
			<hr>
		</div>
	
		<div style="grid-row: 2; grid-column: 1 / span 2;">
			<input type="checkbox" id="rcvAlerts" name="rcvAlerts" value="on" checked disabled/>
			<label for="rcvAlerts" style="font-weight: normal;">Receive email notifications about site alerts <em>(mandatory; always on)</em></label>
		</div>

		<div style="grid-row: 3; grid-column: 2; text-align: right;">
			<button type="reset">Reset</button>
			<button type="submit" style="font-weight: bold;">Save</button>
		</div>
	</form>
</div>
<div>
	<p>Go <a href="{{getServerAddress}}/account/?path=manage">back</a> to the main account page.</p>
</div>
`

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

package contact

const tplJavaScript = `
function verifyForm(formData) {
	if (formData.getTrimmed("subject") == "") {
		setState(STATE_ERROR, "Please enter a subject.", "form", "subject", true);
		return false;
	}

	if (formData.getTrimmed("message") == "") {
		setState(STATE_ERROR, "Please enter a message.", "form", "message", true);	
		return false;
	}

	return true;
}

function handleAction(action) {
	const formData = new FormData(document.querySelector("form"));
	if (!verifyForm(formData)) {
		return;
	}

	setState(STATE_STATUS, "Sending message... this should only take a moment.", "form", null, false);

	var xhr = new XMLHttpRequest();
    xhr.open("POST", "{{getServerAddress}}/" + action);
    xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');

	xhr.onload = function() {
		if (this.status == 200) {
			setState(STATE_SUCCESS, "Your message was successfully sent! A copy of the message has been sent to your email address.");
		} else {
			var resp = JSON.parse(this.responseText);
			setState(STATE_ERROR, "An error occurred while trying to send your message:<br><em>" + resp.error + "</em>", "form", null, true);
		}
	}

	var postData = {
		"subject": formData.getTrimmed("subject"),
		"message": formData.getTrimmed("message")
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
	<p>Contact the ScienceMesh administration using the form below.</p>
	<p style="margin-bottom: 0em;">Please include as much information as possible in your request, especially:</p>
	<ul style="margin-top: 0em;">
		<li>The site your request refers to (if not obvious from your account information)</li>
		<li>Your role within the ScienceMesh site (e.g., administrator, operational team member, etc.)</li>
		<li>Any specific reasons for your request</li>
		<li>Anything else that might help to process your request</li>
	</ul>
</div>
<div>&nbsp;</div>
<div>
	<form id="form" method="POST" class="box container-inline" style="width: 100%;" onSubmit="handleAction('contact'); return false;">
		<div style="grid-row: 1;"><label for="subject">Subject: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 2; grid-column: 1 / span 2;"><input type="text" id="subject" name="subject" {{if .Params.Subject}}value="{{.Params.Subject}}" readonly{{end}}/></div>

		<div style="grid-row: 3;"><label for="message">Message: <span class="mandatory">*</span></label></div>
		<div style="grid-row: 4; grid-column: 1 / span 2;">
			<textarea rows="10" id="message" name="message" style="box-sizing: border-box; width: 100%;" {{if .Params.Info}}placeholder="{{.Params.Info}}"{{end}}>{{if .Params.Message}}{{.Params.Message}}{{end}}</textarea>   
		</div>

		<div style="grid-row: 5; align-self: center;">
			Fields marked with <span class="mandatory">*</span> are mandatory.
		</div>
		<div style="grid-row: 5; grid-column: 2; text-align: right;">
			<button type="reset">Reset</button>
			<button type="submit" style="font-weight: bold;">Send</button>
		</div>
	</form>
</div>
<div>
	<p>Go <a href="{{getServerAddress}}/account/?path=manage">back</a> to the main account page.</p>
</div>
`

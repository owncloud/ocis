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

package html

const panelTemplate = `
<!DOCTYPE html>
<html>
<head>	
	<script>
		const STATE_NONE = 0
		const STATE_STATUS = 1
		const STATE_SUCCESS = 2
		const STATE_ERROR = 3

		function enableForm(id, enable) {
			var form = document.getElementById(id);
			var elements = form.elements;
			for (var i = 0, len = elements.length; i < len; ++i) {
				elements[i].disabled = !enable;
			}
		}

		function setElementVisibility(id, visible) {
			var elem = document.getElementById(id);
			if (visible) {			
				elem.classList.add("visible");
				elem.classList.remove("hidden");
			} else {
				elem.classList.remove("visible");
				elem.classList.add("hidden");
			}
		}

		function setState(state, msg = "", formId = null, focusElem = null, formState = null) {
			setElementVisibility("status", state == STATE_STATUS);
			setElementVisibility("success", state == STATE_SUCCESS);
			setElementVisibility("error", state == STATE_ERROR);

			var elem = null;
			switch (state) {
			case STATE_STATUS:
				elem = document.getElementById("status");	
				break;

			case STATE_SUCCESS:
				elem = document.getElementById("success");	
				break;

			case STATE_ERROR:
				elem = document.getElementById("error");	
				break;
			}

			if (elem !== null) {
				elem.innerHTML = msg;
			}

			if (formId !== null && formState !== null) {
				enableForm(formId, formState);
			}

			if (focusElem !== null) {
				var elem = document.getElementById(focusElem);
				elem.focus();
			}
		}

		FormData.prototype.getTrimmed = function(id) {
			var val = this.get(id);

			if (val != null) {
				return val.trim();
			} else {
				return "";
			}
		}

		$(CONTENT_JAVASCRIPT)
	</script>
	<style>
		form {
			border-color: lightgray !important;
		}
		button {
			min-width: 140px;
		}
		input {
			width: 95%;
		}
		label {
			font-weight: bold;
		}
		h1 {
			text-align: center;
		}

		.box {
			width: 100%;
			border: 1px solid black;
			border-radius: 10px;
			padding: 10px;
		}
		.container {
			width: 900px;
			display: grid;
			grid-gap: 10px;
			position: absolute;
			left: 50%;
			transform: translate(-50%, 0);
		}
		.container-inline {
			display: inline-grid;
			grid-gap: 10px;
		}
		.status {
			border-color: #F7B22A;
			background: #FFEABF;
		}
		.success {
			border-color: #3CAC3A;
			background: #D3EFD2;
		}
		.error {
			border-color: #F20000;
			background: #F4D0D0;
		}
		.visible {
			display: block;
		}
		.hidden {
			display: none;
		}

		$(CONTENT_STYLESHEET)
	</style>
	<title>$(TITLE)</title>
</head>
<body>

<div class="container">
	<div><h1>$(CAPTION)</h1></div>
	
	$(CONTENT_BODY)
	
	<div id="status" class="box status hidden">
	</div>
	<div id="success" class="box success hidden">
	</div>
	<div id="error" class="box error hidden">
	</div>
</div>
</body>
</html>
`

/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package identifier

import (
	"html/template"
	"net/http"
	"net/url"

	"github.com/libregraph/lico/version"
)

var trampolinTemplate = template.Must(template.New("trampolin").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body trampolin="{{.URI}}">
<script src="trampolin/trampolin.js?v={{.Version}}"></script>
<noscript>Javascript is required for this app.</noscript>
</body>
</html>
`))

var trampolinScript = []byte(`(function(window) {
window.location.replace(document.body.getAttribute('trampolin'));
}(window));
`)

var trampolinVersion = url.QueryEscape(version.Version)

type trampolinData struct {
	URI     string
	Version string
}

func (i *Identifier) writeTrampolinHTML(rw http.ResponseWriter, req *http.Request, uri *url.URL) {
	data := &trampolinData{
		URI:     uri.String(),
		Version: trampolinVersion,
	}

	rw.Header().Set("Content-Type", "text/html")
	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Expires", "0")
	err := trampolinTemplate.Execute(rw, data)
	if err != nil {
		i.logger.WithError(err).Errorln("failed to write trampolin")
	}
}

func (i *Identifier) writeTrampolinScript(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/javascript")
	rw.Header().Set("Cache-Control", "public, max-age=31536000, immutable")

	rw.Write(trampolinScript)
}

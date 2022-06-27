/*
 * Copyright 2017-2020 Kopano and its licensors
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

package samlext

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"net/url"

	"github.com/beevik/etree"
	"github.com/crewjam/saml"
)

type LogoutRequest struct {
	*saml.LogoutRequest
}

// Redirect returns a URL suitable for using the redirect binding with the response.
func (req *LogoutRequest) Redirect(relayState string) *url.URL {
	w := &bytes.Buffer{}
	w1 := base64.NewEncoder(base64.StdEncoding, w)
	w2, _ := flate.NewWriter(w1, 9)
	doc := etree.NewDocument()
	doc.SetRoot(req.Element())
	if _, err := doc.WriteTo(w2); err != nil {
		panic(err)
	}
	w2.Close()
	w1.Close()

	rv, _ := url.Parse(req.Destination)

	query := rv.Query()
	query.Set("SAMLRequest", w.String())
	if relayState != "" {
		query.Set("RelayState", relayState)
	}
	rv.RawQuery = query.Encode()

	return rv
}

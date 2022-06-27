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
	"encoding/xml"
	"fmt"
	"net/url"

	"github.com/crewjam/saml"
	"stash.kopano.io/kgol/rndm"
)

func MakeLogoutResponse(sp *saml.ServiceProvider, req *saml.LogoutRequest, status *saml.Status, binding string) (*LogoutResponse, error) {

	res := &LogoutResponse{&saml.LogoutResponse{
		ID:           fmt.Sprintf("id-%x", rndm.GenerateRandomBytes(20)),
		InResponseTo: req.ID,

		Version:      "2.0",
		IssueInstant: saml.TimeNow(),
		Destination:  sp.GetSLOBindingLocation(binding),

		Issuer: &saml.Issuer{
			Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
			Value:  firstSet(sp.EntityID, sp.MetadataURL.String()),
		},
	}}

	if status != nil {
		res.LogoutResponse.Status = *status
	}

	return res, nil
}

func firstSet(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

type LogoutResponse struct {
	*saml.LogoutResponse
}

// Redirect returns a URL suitable for using the redirect binding with the response.
func (res *LogoutResponse) Redirect(relayState string) *url.URL {
	w := &bytes.Buffer{}
	w1 := base64.NewEncoder(base64.StdEncoding, w)
	w2, _ := flate.NewWriter(w1, 9)
	e := xml.NewEncoder(w2)
	if err := e.Encode(res); err != nil {
		panic(err)
	}
	w2.Close()
	w1.Close()

	rv, _ := url.Parse(res.Destination)

	query := rv.Query()
	query.Set("SAMLResponse", w.String())
	if relayState != "" {
		query.Set("RelayState", relayState)
	}
	rv.RawQuery = query.Encode()

	return rv
}

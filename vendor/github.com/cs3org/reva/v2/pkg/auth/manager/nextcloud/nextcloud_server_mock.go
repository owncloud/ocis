// Copyright 2018-2021 CERN
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

package nextcloud

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

// Response contains data for the Nextcloud mock server to respond
// and to switch to a new server state
type Response struct {
	code           int
	body           string
	newServerState string
}

const serverStateError = "ERROR"
const serverStateEmpty = "EMPTY"
const serverStateHome = "HOME"

var serverState = serverStateEmpty

var responses = map[string]Response{
	`POST /apps/sciencemesh/~einstein/api/auth/Authenticate {"clientID":"einstein","clientSecret":"relativity"}`: {200, `{"user":{"id":{"idp":"some-idp","opaque_id":"some-opaque-user-id","type":1}},"scopes":{"user":{"resource":{"decoder":"json","value":"eyJyZXNvdXJjZV9pZCI6eyJzdG9yYWdlX2lkIjoic3RvcmFnZS1pZCIsIm9wYXF1ZV9pZCI6Im9wYXF1ZS1pZCJ9LCJwYXRoIjoic29tZS9maWxlL3BhdGgudHh0In0="},"role":1}}}`, serverStateHome},
}

// GetNextcloudServerMock returns a handler that pretends to be a remote Nextcloud server
func GetNextcloudServerMock(called *[]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, r.Body)
		if err != nil {
			panic("Error reading response into buffer")
		}
		var key = fmt.Sprintf("%s %s %s", r.Method, r.URL, buf.String())
		*called = append(*called, key)
		response := responses[key]
		if (response == Response{}) {
			key = fmt.Sprintf("%s %s %s %s", r.Method, r.URL, buf.String(), serverState)
			response = responses[key]
		}
		if (response == Response{}) {
			fmt.Printf("%s %s %s %s", r.Method, r.URL, buf.String(), serverState)
			response = Response{500, fmt.Sprintf("response not defined! %s", key), serverStateEmpty}
		}
		serverState = responses[key].newServerState
		if serverState == `` {
			serverState = serverStateError
		}
		w.WriteHeader(response.code)
		// w.Header().Set("Etag", "mocker-etag")
		_, err = w.Write([]byte(responses[key].body))
		if err != nil {
			panic(err)
		}
	})
}

// TestingHTTPClient thanks to https://itnext.io/how-to-stub-requests-to-remote-hosts-with-go-6c2c1db32bf2
// Ideally, this function would live in tests/helpers, but
// if we put it there, it gets excluded by .dockerignore, and the
// Docker build fails (see https://github.com/cs3org/reva/issues/1999)
// So putting it here for now - open to suggestions if someone knows
// a better way to inject this.
func TestingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

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

package sender

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rhttp"

	"github.com/pkg/errors"
)

const createOCMCoreShareEndpoint = "shares"

func getOCMEndpoint(originProvider *ocmprovider.ProviderInfo) (string, error) {
	for _, s := range originProvider.Services {
		if s.Endpoint.Type.Name == "OCM" {
			return s.Endpoint.Path, nil
		}
	}
	return "", errors.New("json: ocm endpoint not specified for mesh provider")
}

// Send executes the POST to the OCM shares endpoint to create the share at the
// remote site.
func Send(requestBodyMap map[string]interface{}, pi *ocmprovider.ProviderInfo) error {
	requestBody, err := json.Marshal(requestBodyMap)
	if err != nil {
		err = errors.Wrap(err, "error marshalling request body")
		return err
	}
	ocmEndpoint, err := getOCMEndpoint(pi)
	if err != nil {
		return err
	}
	u, err := url.Parse(ocmEndpoint)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, createOCMCoreShareEndpoint)
	recipientURL := u.String()

	req, err := http.NewRequest("POST", recipientURL, strings.NewReader(string(requestBody)))
	if err != nil {
		return errors.Wrap(err, "sender: error framing post request")
	}
	req.Header.Set("Content-Type", "application/json; param=value")
	client := rhttp.GetHTTPClient(
		rhttp.Timeout(5 * time.Second),
	)

	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "sender: error sending post request")
		return err
	}

	defer resp.Body.Close()
	if (resp.StatusCode != http.StatusCreated) && (resp.StatusCode != http.StatusOK) {
		respBody, e := io.ReadAll(resp.Body)
		if e != nil {
			e = errors.Wrap(e, "sender: error reading request body")
			return e
		}
		err = errors.Wrap(fmt.Errorf("%s: %s", resp.Status, string(respBody)), "sender: error sending create ocm core share post request")
		return err
	}
	return nil
}

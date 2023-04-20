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

package gocdb

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cs3org/reva/v2/pkg/mentix/utils/network"
	"github.com/cs3org/reva/v2/pkg/siteacc/data"
	"github.com/pkg/errors"
)

const (
	opCreateOrUpdate = "CreateOrUpdate"
	opDelete         = "Delete"
)

type writeAccountUserData struct {
	Email       string `json:"Email"`
	FirstName   string `json:"FirstName"`
	LastName    string `json:"LastName"`
	PhoneNumber string `json:"PhoneNumber"`
}

type writeAccountData struct {
	APIKey    string `json:"APIKey"`
	Operation string `json:"Operation"`

	Data writeAccountUserData `json:"Data"`
}

func writeAccount(account *data.Account, operation string, address string, apiKey string) error {
	// Fill in the data to send
	userData := getWriteAccountData(account)
	userData.APIKey = apiKey
	userData.Operation = operation

	// Send the data to the GOCDB endpoint
	endpointURL, err := network.GenerateURL(address, "/ext/v1/user", network.URLParams{})
	if err != nil {
		return errors.Wrap(err, "unable to generate the GOCDB URL")
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return errors.Wrap(err, "unable to marshal the user data")
	}

	req, err := http.NewRequest(http.MethodPost, endpointURL.String(), bytes.NewReader(jsonData))
	if err != nil {
		return errors.Wrap(err, "unable to create HTTP request")
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to send data to endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		msg, _ := io.ReadAll(resp.Body)
		return errors.Errorf("unable to perform request: %v", string(msg))
	}

	return nil
}

func getWriteAccountData(account *data.Account) *writeAccountData {
	return &writeAccountData{
		Data: writeAccountUserData{
			Email:       account.Email,
			FirstName:   account.FirstName,
			LastName:    account.LastName,
			PhoneNumber: account.PhoneNumber,
		},
	}
}

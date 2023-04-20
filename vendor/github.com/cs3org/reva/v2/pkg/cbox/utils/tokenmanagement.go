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

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cs3org/reva/v2/pkg/rhttp"
)

// APITokenManager stores config related to api management
type APITokenManager struct {
	oidcToken OIDCToken
	conf      *config
	client    *http.Client
}

// OIDCToken stores the OIDC token used to authenticate requests to the REST API service
type OIDCToken struct {
	sync.Mutex          // concurrent access to apiToken and tokenExpirationTime
	apiToken            string
	tokenExpirationTime time.Time
}

type config struct {
	TargetAPI         string
	OIDCTokenEndpoint string
	ClientID          string
	ClientSecret      string
}

// InitAPITokenManager initializes a new APITokenManager
func InitAPITokenManager(targetAPI, oidcTokenEndpoint, clientID, clientSecret string) *APITokenManager {
	return &APITokenManager{
		conf: &config{
			TargetAPI:         targetAPI,
			OIDCTokenEndpoint: oidcTokenEndpoint,
			ClientID:          clientID,
			ClientSecret:      clientSecret,
		},
		client: rhttp.GetHTTPClient(
			rhttp.Timeout(10*time.Second),
			rhttp.Insecure(true),
		),
	}
}

func (a *APITokenManager) renewAPIToken(ctx context.Context, forceRenewal bool) error {
	// Received tokens have an expiration time of 20 minutes.
	// Take a couple of seconds as buffer time for the API call to complete
	if forceRenewal || a.oidcToken.tokenExpirationTime.Before(time.Now().Add(time.Second*time.Duration(2))) {
		token, expiration, err := a.getAPIToken(ctx)
		if err != nil {
			return err
		}

		a.oidcToken.Lock()
		defer a.oidcToken.Unlock()

		a.oidcToken.apiToken = token
		a.oidcToken.tokenExpirationTime = expiration
	}
	return nil
}

func (a *APITokenManager) getAPIToken(ctx context.Context) (string, time.Time, error) {

	params := url.Values{
		"grant_type": {"client_credentials"},
		"audience":   {a.conf.TargetAPI},
	}

	httpReq, err := http.NewRequest("POST", a.conf.OIDCTokenEndpoint, strings.NewReader(params.Encode()))
	if err != nil {
		return "", time.Time{}, err
	}
	httpReq.SetBasicAuth(a.conf.ClientID, a.conf.ClientSecret)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	httpRes, err := a.client.Do(httpReq)
	if err != nil {
		return "", time.Time{}, err
	}
	defer httpRes.Body.Close()

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", time.Time{}, err
	}
	if httpRes.StatusCode < 200 || httpRes.StatusCode > 299 {
		return "", time.Time{}, errors.New("rest: get token endpoint returned " + httpRes.Status)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", time.Time{}, err
	}

	expirationSecs := result["expires_in"].(float64)
	expirationTime := time.Now().Add(time.Second * time.Duration(expirationSecs))
	return result["access_token"].(string), expirationTime, nil
}

// SendAPIGetRequest makes an API GET Request to the passed URL
func (a *APITokenManager) SendAPIGetRequest(ctx context.Context, url string, forceRenewal bool) (map[string]interface{}, error) {
	err := a.renewAPIToken(ctx, forceRenewal)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// We don't need to take the lock when reading apiToken, because if we reach here,
	// the token is valid at least for a couple of seconds. Even if another request modifies
	// the token and expiration time while this request is in progress, the current token will still be valid.
	httpReq.Header.Set("Authorization", "Bearer "+a.oidcToken.apiToken)

	httpRes, err := a.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusUnauthorized {
		// The token is no longer valid, try renewing it
		return a.SendAPIGetRequest(ctx, url, true)
	}
	if httpRes.StatusCode < 200 || httpRes.StatusCode > 299 {
		return nil, errors.New("rest: API request returned " + httpRes.Status)
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

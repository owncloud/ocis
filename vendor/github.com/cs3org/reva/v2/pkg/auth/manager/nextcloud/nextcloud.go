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

// Package nextcloud verifies a clientID and clientSecret against a Nextcloud backend.
package nextcloud

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/auth"
	"github.com/cs3org/reva/v2/pkg/auth/manager/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	registry.Register("nextcloud", New)
}

// Manager is the Nextcloud-based implementation of the auth.Manager interface
// see https://github.com/cs3org/reva/blob/v1.13.0/pkg/auth/auth.go#L32-L35
type Manager struct {
	client       *http.Client
	sharedSecret string
	endPoint     string
}

// AuthManagerConfig contains config for a Nextcloud-based AuthManager
type AuthManagerConfig struct {
	EndPoint     string `mapstructure:"endpoint" docs:";The Nextcloud backend endpoint for user check"`
	SharedSecret string `mapstructure:"shared_secret"`
	MockHTTP     bool   `mapstructure:"mock_http"`
}

// Action describes a REST request to forward to the Nextcloud backend
type Action struct {
	verb     string
	username string
	argS     string
}

func (c *AuthManagerConfig) init() {
}

func parseConfig(m map[string]interface{}) (*AuthManagerConfig, error) {
	c := &AuthManagerConfig{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New returns an auth manager implementation that verifies against a Nextcloud backend.
func New(m map[string]interface{}) (auth.Manager, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}
	c.init()

	return NewAuthManager(c)
}

// NewAuthManager returns a new Nextcloud-based AuthManager
func NewAuthManager(c *AuthManagerConfig) (*Manager, error) {
	var client *http.Client
	if c.MockHTTP {
		// called := make([]string, 0)
		// nextcloudServerMock := GetNextcloudServerMock(&called)
		// client, _ = TestingHTTPClient(nextcloudServerMock)

		// Wait for SetHTTPClient to be called later
		client = nil
	} else {
		if len(c.EndPoint) == 0 {
			return nil, errors.New("Please specify 'endpoint' in '[grpc.services.authprovider.auth_managers.nextcloud]'")
		}
		client = &http.Client{}
	}

	return &Manager{
		endPoint:     c.EndPoint, // e.g. "http://nc/apps/sciencemesh/"
		sharedSecret: c.SharedSecret,
		client:       client,
	}, nil
}

// Configure method as defined in https://github.com/cs3org/reva/blob/v1.13.0/pkg/auth/auth.go#L32-L35
func (am *Manager) Configure(ml map[string]interface{}) error {
	return nil
}

// SetHTTPClient sets the HTTP client
func (am *Manager) SetHTTPClient(c *http.Client) {
	am.client = c
}

func (am *Manager) do(ctx context.Context, a Action) (int, []byte, error) {
	log := appctx.GetLogger(ctx)
	url := am.endPoint + "~" + a.username + "/api/auth/" + a.verb
	log.Info().Msgf("am.do %s %s %s", url, a.argS, am.sharedSecret)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(a.argS))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("X-Reva-Secret", am.sharedSecret)

	req.Header.Set("Content-Type", "application/json")
	resp, err := am.client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	log.Info().Msgf("am.do response %d %s", resp.StatusCode, body)
	return resp.StatusCode, body, nil
}

// Authenticate method as defined in https://github.com/cs3org/reva/blob/28500a8/pkg/auth/auth.go#L31-L33
func (am *Manager) Authenticate(ctx context.Context, clientID, clientSecret string) (*user.User, map[string]*authpb.Scope, error) {
	type paramsObj struct {
		ClientID     string `json:"clientID"`
		ClientSecret string `json:"clientSecret"`
		// Scope        authpb.Scope
	}
	bodyObj := &paramsObj{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		// Scope: authpb.Scope{
		// 	Resource: &types.OpaqueEntry{
		// 		Decoder: "json",
		// 		Value:   []byte(`{"resource_id":{"storage_id":"storage-id","opaque_id":"opaque-id"},"path":"some/file/path.txt"}`),
		// 	},
		// 	Role: authpb.Role_ROLE_OWNER,
		// },
	}
	bodyStr, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, nil, err
	}
	log := appctx.GetLogger(ctx)
	log.Info().Msgf("Authenticate %s %s", clientID, bodyStr)

	statusCode, body, err := am.do(ctx, Action{"Authenticate", clientID, string(bodyStr)})

	if err != nil {
		return nil, nil, err
	}

	if statusCode != 200 {
		return nil, nil, errors.New("Username/password not recognized by Nextcloud backend")
	}

	type resultsObj struct {
		User   user.User                `json:"user"`
		Scopes map[string]*authpb.Scope `json:"scopes"`
	}
	result := &resultsObj{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, err
	}
	var pointersMap = make(map[string]*authpb.Scope)
	for k := range result.Scopes {
		scope := result.Scopes[k]
		pointersMap[k] = scope
	}
	return &result.User, pointersMap, nil
}

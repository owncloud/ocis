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

package json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/invite"
	"github.com/cs3org/reva/v2/pkg/ocm/invite/manager/registry"
	"github.com/cs3org/reva/v2/pkg/ocm/invite/token"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const acceptInviteEndpoint = "invites/accept"

type inviteModel struct {
	File          string
	Invites       map[string]*invitepb.InviteToken `json:"invites"`
	AcceptedUsers map[string][]*userpb.User        `json:"accepted_users"`
}

type manager struct {
	config     *config
	sync.Mutex // concurrent access to the file
	model      *inviteModel
	client     *http.Client
}

type config struct {
	File                string `mapstructure:"file"`
	Expiration          string `mapstructure:"expiration"`
	InsecureConnections bool   `mapstructure:"insecure_connections"`
}

func init() {
	registry.Register("json", New)
}

func (c *config) init() error {
	if c.File == "" {
		c.File = "/var/tmp/reva/ocm-invites.json"
	}

	if c.Expiration == "" {
		c.Expiration = token.DefaultExpirationTime
	}
	return nil
}

// New returns a new invite manager object.
func New(m map[string]interface{}) (invite.Manager, error) {

	config, err := parseConfig(m)
	if err != nil {
		err = errors.Wrap(err, "error parsing config for json invite manager")
		return nil, err
	}
	err = config.init()
	if err != nil {
		err = errors.Wrap(err, "error setting config defaults for json invite manager")
		return nil, err
	}

	// load or create file
	model, err := loadOrCreate(config.File)
	if err != nil {
		err = errors.Wrap(err, "error loading the file containing the invites")
		return nil, err
	}

	manager := &manager{
		config: config,
		model:  model,
		client: rhttp.GetHTTPClient(
			rhttp.Timeout(5*time.Second),
			rhttp.Insecure(config.InsecureConnections),
		),
	}

	return manager, nil
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	return c, nil
}

func loadOrCreate(file string) (*inviteModel, error) {

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		if err := os.WriteFile(file, []byte("{}"), 0700); err != nil {
			err = errors.Wrap(err, "error creating the invite storage file: "+file)
			return nil, err
		}
	}

	fd, err := os.OpenFile(file, os.O_CREATE, 0644)
	if err != nil {
		err = errors.Wrap(err, "error opening the invite storage file: "+file)
		return nil, err
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		err = errors.Wrap(err, "error reading the data")
		return nil, err
	}

	model := &inviteModel{}
	if err := json.Unmarshal(data, model); err != nil {
		err = errors.Wrap(err, "error decoding invite data to json")
		return nil, err
	}

	if model.Invites == nil {
		model.Invites = make(map[string]*invitepb.InviteToken)
	}
	if model.AcceptedUsers == nil {
		model.AcceptedUsers = make(map[string][]*userpb.User)
	}

	model.File = file
	return model, nil
}

func (model *inviteModel) Save() error {
	data, err := json.Marshal(model)
	if err != nil {
		err = errors.Wrap(err, "error encoding invite data to json")
		return err
	}

	if err := os.WriteFile(model.File, data, 0644); err != nil {
		err = errors.Wrap(err, "error writing invite data to file: "+model.File)
		return err
	}

	return nil
}

func (m *manager) GenerateToken(ctx context.Context) (*invitepb.InviteToken, error) {

	contexUser := ctxpkg.ContextMustGetUser(ctx)
	inviteToken, err := token.CreateToken(m.config.Expiration, contexUser.GetId())
	if err != nil {
		return nil, err
	}

	// Store token data
	m.Lock()
	defer m.Unlock()

	m.model.Invites[inviteToken.GetToken()] = inviteToken
	if err := m.model.Save(); err != nil {
		err = errors.Wrap(err, "error saving model")
		return nil, err
	}

	return inviteToken, nil
}

func (m *manager) ForwardInvite(ctx context.Context, invite *invitepb.InviteToken, originProvider *ocmprovider.ProviderInfo) error {

	contextUser := ctxpkg.ContextMustGetUser(ctx)
	recipientProvider := contextUser.GetId().GetIdp()

	requestBody := url.Values{
		"token":             {invite.GetToken()},
		"userID":            {contextUser.GetId().GetOpaqueId()},
		"recipientProvider": {recipientProvider},
		"email":             {contextUser.GetMail()},
		"name":              {contextUser.GetDisplayName()},
	}

	ocmEndpoint, err := getOCMEndpoint(originProvider)
	if err != nil {
		return err
	}
	u, err := url.Parse(ocmEndpoint)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, acceptInviteEndpoint)
	recipientURL := u.String()

	req, err := http.NewRequest("POST", recipientURL, strings.NewReader(requestBody.Encode()))
	if err != nil {
		return errors.Wrap(err, "json: error framing post request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := m.client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "json: error sending post request")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, e := io.ReadAll(resp.Body)
		if e != nil {
			return errors.Wrap(e, "json: error reading request body")
		}
		return errors.Wrap(fmt.Errorf("%s: %s", resp.Status, string(respBody)), "json: error sending accept post request")
	}

	return nil
}

func (m *manager) AcceptInvite(ctx context.Context, invite *invitepb.InviteToken, remoteUser *userpb.User) error {

	m.Lock()
	defer m.Unlock()

	inviteToken, err := m.getTokenIfValid(invite)
	if err != nil {
		return err
	}

	currUser := inviteToken.GetUserId()

	// do not allow the user who created the token to accept it
	if remoteUser.Id.Idp == currUser.Idp && remoteUser.Id.OpaqueId == currUser.OpaqueId {
		return errors.New("json: token creator and recipient are the same")
	}

	for _, acceptedUser := range m.model.AcceptedUsers[currUser.GetOpaqueId()] {
		if acceptedUser.Id.GetOpaqueId() == remoteUser.Id.OpaqueId && acceptedUser.Id.GetIdp() == remoteUser.Id.Idp {
			return errors.New("json: user already added to accepted users")
		}
	}
	m.model.AcceptedUsers[currUser.GetOpaqueId()] = append(m.model.AcceptedUsers[currUser.GetOpaqueId()], remoteUser)
	if err := m.model.Save(); err != nil {
		err = errors.Wrap(err, "json: error saving model")
		return err
	}
	return nil
}

func (m *manager) GetAcceptedUser(ctx context.Context, remoteUserID *userpb.UserId) (*userpb.User, error) {

	userKey := ctxpkg.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	for _, acceptedUser := range m.model.AcceptedUsers[userKey] {
		if (acceptedUser.Id.GetOpaqueId() == remoteUserID.OpaqueId) && (remoteUserID.Idp == "" || acceptedUser.Id.GetIdp() == remoteUserID.Idp) {
			return acceptedUser, nil
		}
	}
	return nil, errtypes.NotFound(remoteUserID.OpaqueId)
}

func (m *manager) FindAcceptedUsers(ctx context.Context, query string) ([]*userpb.User, error) {
	users := []*userpb.User{}
	userKey := ctxpkg.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	for _, acceptedUser := range m.model.AcceptedUsers[userKey] {
		if query == "" || userContains(acceptedUser, query) {
			users = append(users, acceptedUser)
		}
	}
	return users, nil
}

func userContains(u *userpb.User, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(u.Username), query) || strings.Contains(strings.ToLower(u.DisplayName), query) ||
		strings.Contains(strings.ToLower(u.Mail), query) || strings.Contains(strings.ToLower(u.Id.OpaqueId), query)
}

func (m *manager) getTokenIfValid(token *invitepb.InviteToken) (*invitepb.InviteToken, error) {
	inviteToken, ok := m.model.Invites[token.GetToken()]
	if !ok {
		return nil, errors.New("json: invalid token")
	}

	if uint64(time.Now().Unix()) > inviteToken.Expiration.Seconds {
		return nil, errors.New("json: token expired")
	}
	return inviteToken, nil
}

func getOCMEndpoint(originProvider *ocmprovider.ProviderInfo) (string, error) {
	for _, s := range originProvider.Services {
		if s.Endpoint.Type.Name == "OCM" {
			return s.Endpoint.Path, nil
		}
	}
	return "", errors.New("json: ocm endpoint not specified for mesh provider")
}

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

package memory

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	ocmprovider "github.com/cs3org/go-cs3apis/cs3/ocm/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/ocm/invite"
	"github.com/cs3org/reva/v2/pkg/ocm/invite/manager/registry"
	"github.com/cs3org/reva/v2/pkg/ocm/invite/token"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const acceptInviteEndpoint = "invites/accept"

func init() {
	registry.Register("memory", New)
}

func (c *config) init() {
	if c.Expiration == "" {
		c.Expiration = token.DefaultExpirationTime
	}
}

// New returns a new invite manager.
func New(m map[string]interface{}) (invite.Manager, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error creating a new manager")
		return nil, err
	}
	c.init()

	return &manager{
		Invites:       sync.Map{},
		AcceptedUsers: sync.Map{},
		Config:        c,
		Client: rhttp.GetHTTPClient(
			rhttp.Timeout(5*time.Second),
			rhttp.Insecure(c.InsecureConnections),
		),
	}, nil
}

type manager struct {
	Invites       sync.Map
	AcceptedUsers sync.Map
	Client        *http.Client
	Config        *config
}

type config struct {
	Expiration          string `mapstructure:"expiration"`
	InsecureConnections bool   `mapstructure:"insecure_connections"`
}

func (m *manager) GenerateToken(ctx context.Context) (*invitepb.InviteToken, error) {

	ctxUser := ctxpkg.ContextMustGetUser(ctx)
	inviteToken, err := token.CreateToken(m.Config.Expiration, ctxUser.GetId())
	if err != nil {
		return nil, errors.Wrap(err, "memory: error creating token")
	}

	m.Invites.Store(inviteToken.GetToken(), inviteToken)
	return inviteToken, nil
}

func (m *manager) ForwardInvite(ctx context.Context, invite *invitepb.InviteToken, originProvider *ocmprovider.ProviderInfo) error {

	contextUser := ctxpkg.ContextMustGetUser(ctx)
	requestBody := url.Values{
		"token":             {invite.GetToken()},
		"userID":            {contextUser.GetId().GetOpaqueId()},
		"recipientProvider": {contextUser.GetId().GetIdp()},
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

	resp, err := m.Client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "memory: error sending post request")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.Wrap(errors.New(resp.Status), "memory: error sending accept post request")
		return err
	}

	return nil
}

func (m *manager) AcceptInvite(ctx context.Context, invite *invitepb.InviteToken, remoteUser *userpb.User) error {
	inviteToken, err := m.getTokenIfValid(invite)
	if err != nil {
		return err
	}

	currUser := inviteToken.GetUserId()

	// do not allow the user who created the token to accept it
	if remoteUser.Id.Idp == currUser.Idp && remoteUser.Id.OpaqueId == currUser.OpaqueId {
		return errors.New("memory: token creator and recipient are the same")
	}

	usersList, ok := m.AcceptedUsers.Load(currUser)
	acceptedUsers := usersList.([]*userpb.User)
	if ok {
		for _, acceptedUser := range acceptedUsers {
			if acceptedUser.Id.GetOpaqueId() == remoteUser.Id.OpaqueId && acceptedUser.Id.GetIdp() == remoteUser.Id.Idp {
				return errors.New("memory: user already added to accepted users")
			}
		}

		acceptedUsers = append(acceptedUsers, remoteUser)
		m.AcceptedUsers.Store(currUser.GetOpaqueId(), acceptedUsers)
	} else {
		acceptedUsers := []*userpb.User{remoteUser}
		m.AcceptedUsers.Store(currUser.GetOpaqueId(), acceptedUsers)
	}
	return nil
}

func (m *manager) GetAcceptedUser(ctx context.Context, remoteUserID *userpb.UserId) (*userpb.User, error) {
	currUser := ctxpkg.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	usersList, ok := m.AcceptedUsers.Load(currUser)
	if !ok {
		return nil, errtypes.NotFound(remoteUserID.OpaqueId)
	}

	acceptedUsers := usersList.([]*userpb.User)
	for _, acceptedUser := range acceptedUsers {
		if (acceptedUser.Id.GetOpaqueId() == remoteUserID.OpaqueId) && (remoteUserID.Idp == "" || acceptedUser.Id.GetIdp() == remoteUserID.Idp) {
			return acceptedUser, nil
		}
	}
	return nil, errtypes.NotFound(remoteUserID.OpaqueId)
}

func (m *manager) FindAcceptedUsers(ctx context.Context, query string) ([]*userpb.User, error) {
	currUser := ctxpkg.ContextMustGetUser(ctx).GetId().GetOpaqueId()
	usersList, ok := m.AcceptedUsers.Load(currUser)
	if !ok {
		return []*userpb.User{}, nil
	}

	users := []*userpb.User{}
	acceptedUsers := usersList.([]*userpb.User)
	for _, acceptedUser := range acceptedUsers {
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
	tokenInterface, ok := m.Invites.Load(token.GetToken())
	if !ok {
		return nil, errors.New("memory: invalid token")
	}

	inviteToken := tokenInterface.(*invitepb.InviteToken)
	if uint64(time.Now().Unix()) > inviteToken.Expiration.Seconds {
		return nil, errors.New("memory: token expired")
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

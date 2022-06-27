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

package ocminvitemanager

import (
	"context"

	invitepb "github.com/cs3org/go-cs3apis/cs3/ocm/invite/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/invite"
	"github.com/cs3org/reva/v2/pkg/ocm/invite/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func init() {
	rgrpc.Register("ocminvitemanager", New)
}

type config struct {
	Driver  string                            `mapstructure:"driver"`
	Drivers map[string]map[string]interface{} `mapstructure:"drivers"`
}

type service struct {
	conf *config
	im   invite.Manager
}

func (c *config) init() {
	if c.Driver == "" {
		c.Driver = "json"
	}
}

func (s *service) Register(ss *grpc.Server) {
	invitepb.RegisterInviteAPIServer(ss, s)
}

func getInviteManager(c *config) (invite.Manager, error) {
	if f, ok := registry.NewFuncs[c.Driver]; ok {
		return f(c.Drivers[c.Driver])
	}
	return nil, errtypes.NotFound("driver not found: " + c.Driver)
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New creates a new OCM invite manager svc
func New(m map[string]interface{}, ss *grpc.Server) (rgrpc.Service, error) {

	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}
	c.init()

	im, err := getInviteManager(c)
	if err != nil {
		return nil, err
	}

	service := &service{
		conf: c,
		im:   im,
	}
	return service, nil
}

func (s *service) Close() error {
	return nil
}

func (s *service) UnprotectedEndpoints() []string {
	return []string{"/cs3.ocm.invite.v1beta1.InviteAPI/AcceptInvite"}
}

func (s *service) GenerateInviteToken(ctx context.Context, req *invitepb.GenerateInviteTokenRequest) (*invitepb.GenerateInviteTokenResponse, error) {
	token, err := s.im.GenerateToken(ctx)
	if err != nil {
		return &invitepb.GenerateInviteTokenResponse{
			Status: status.NewInternal(ctx, "error generating invite token"),
		}, nil
	}

	return &invitepb.GenerateInviteTokenResponse{
		Status:      status.NewOK(ctx),
		InviteToken: token,
	}, nil
}

func (s *service) ForwardInvite(ctx context.Context, req *invitepb.ForwardInviteRequest) (*invitepb.ForwardInviteResponse, error) {
	err := s.im.ForwardInvite(ctx, req.InviteToken, req.OriginSystemProvider)
	if err != nil {
		return &invitepb.ForwardInviteResponse{
			Status: status.NewInternal(ctx, "error forwarding invite"),
		}, nil
	}

	return &invitepb.ForwardInviteResponse{
		Status: status.NewOK(ctx),
	}, nil
}

func (s *service) AcceptInvite(ctx context.Context, req *invitepb.AcceptInviteRequest) (*invitepb.AcceptInviteResponse, error) {
	err := s.im.AcceptInvite(ctx, req.InviteToken, req.RemoteUser)
	if err != nil {
		return &invitepb.AcceptInviteResponse{
			Status: status.NewInternal(ctx, "error accepting invite"),
		}, nil
	}

	return &invitepb.AcceptInviteResponse{
		Status: status.NewOK(ctx),
	}, nil
}

func (s *service) GetAcceptedUser(ctx context.Context, req *invitepb.GetAcceptedUserRequest) (*invitepb.GetAcceptedUserResponse, error) {
	remoteUser, err := s.im.GetAcceptedUser(ctx, req.RemoteUserId)
	if err != nil {
		return &invitepb.GetAcceptedUserResponse{
			Status: status.NewInternal(ctx, "error fetching remote user details"),
		}, nil
	}

	return &invitepb.GetAcceptedUserResponse{
		Status:     status.NewOK(ctx),
		RemoteUser: remoteUser,
	}, nil
}

func (s *service) FindAcceptedUsers(ctx context.Context, req *invitepb.FindAcceptedUsersRequest) (*invitepb.FindAcceptedUsersResponse, error) {
	acceptedUsers, err := s.im.FindAcceptedUsers(ctx, req.Filter)
	if err != nil {
		return &invitepb.FindAcceptedUsersResponse{
			Status: status.NewInternal(ctx, "error finding remote users"),
		}, nil
	}

	return &invitepb.FindAcceptedUsersResponse{
		Status:        status.NewOK(ctx),
		AcceptedUsers: acceptedUsers,
	}, nil
}

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

package ocmcore

import (
	"context"
	"encoding/json"
	"fmt"

	ocmcore "github.com/cs3org/go-cs3apis/cs3/ocm/core/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/ocm/share"
	"github.com/cs3org/reva/v2/pkg/ocm/share/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func init() {
	rgrpc.Register("ocmcore", New)
}

type config struct {
	Driver  string                            `mapstructure:"driver"`
	Drivers map[string]map[string]interface{} `mapstructure:"drivers"`
}

type service struct {
	conf *config
	sm   share.Manager
}

func (c *config) init() {
	if c.Driver == "" {
		c.Driver = "json"
	}
}

func (s *service) Register(ss *grpc.Server) {
	ocmcore.RegisterOcmCoreAPIServer(ss, s)
}

func getShareManager(c *config) (share.Manager, error) {
	if f, ok := registry.NewFuncs[c.Driver]; ok {
		return f(c.Drivers[c.Driver])
	}
	return nil, errtypes.NotFound(fmt.Sprintf("driver not found: %s", c.Driver))
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New creates a new ocm core svc
func New(m map[string]interface{}, ss *grpc.Server) (rgrpc.Service, error) {

	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}
	c.init()

	sm, err := getShareManager(c)
	if err != nil {
		return nil, err
	}

	service := &service{
		conf: c,
		sm:   sm,
	}

	return service, nil
}

func (s *service) Close() error {
	return nil
}

func (s *service) UnprotectedEndpoints() []string {
	return []string{"/cs3.ocm.core.v1beta1.OcmCoreAPI/CreateOCMCoreShare"}
}

// CreateOCMCoreShare is called when an OCM request comes into this reva instance from
func (s *service) CreateOCMCoreShare(ctx context.Context, req *ocmcore.CreateOCMCoreShareRequest) (*ocmcore.CreateOCMCoreShareResponse, error) {
	resource := &provider.ResourceId{
		StorageId: "remote",
		OpaqueId:  req.Name,
	}

	var resourcePermissions *provider.ResourcePermissions
	permOpaque, ok := req.Protocol.Opaque.Map["permissions"]
	if !ok {
		return &ocmcore.CreateOCMCoreShareResponse{
			Status: status.NewInternal(ctx, "resource permissions not set"),
		}, nil
	}
	switch permOpaque.Decoder {
	case "json":
		err := json.Unmarshal(permOpaque.Value, &resourcePermissions)
		if err != nil {
			return &ocmcore.CreateOCMCoreShareResponse{
				Status: status.NewInternal(ctx, "error decoding resource permissions"),
			}, nil
		}
	default:
		return &ocmcore.CreateOCMCoreShareResponse{
			Status: status.NewInternal(ctx, "invalid opaque entry decoder"),
		}, nil
	}

	var token string
	tokenOpaque, ok := req.Protocol.Opaque.Map["token"]
	if !ok {
		return &ocmcore.CreateOCMCoreShareResponse{
			Status: status.NewInternal(ctx, "token not set"),
		}, nil
	}
	switch tokenOpaque.Decoder {
	case "plain":
		token = string(tokenOpaque.Value)
	default:
		return &ocmcore.CreateOCMCoreShareResponse{
			Status: status.NewInternal(ctx, "invalid opaque entry decoder"),
		}, nil
	}

	grant := &ocm.ShareGrant{
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_USER,
			// For now, we only support user shares.
			// TODO (ishank011): To be updated once this is decided.
			Id: &provider.Grantee_UserId{UserId: req.ShareWith},
			// passing this in grant.Grantee.Opaque because ShareGrant itself doesn't have a root opaque.
			Opaque: &typespb.Opaque{
				Map: map[string]*typespb.OpaqueEntry{
					"remoteShareId": {
						Decoder: "plain",
						Value:   []byte(req.ProviderId),
					},
				},
			},
		},
		Permissions: &ocm.SharePermissions{
			Permissions: resourcePermissions,
		},
	}

	var shareType ocm.Share_ShareType
	switch req.Protocol.Name {
	case "datatx":
		shareType = ocm.Share_SHARE_TYPE_TRANSFER
	default:
		shareType = ocm.Share_SHARE_TYPE_REGULAR
	}

	share, err := s.sm.Share(ctx, resource, grant, req.Name, nil, "", req.Owner, token, shareType)

	if err != nil {
		return &ocmcore.CreateOCMCoreShareResponse{
			Status: status.NewInternal(ctx, "error creating ocm core share"),
		}, nil
	}

	res := &ocmcore.CreateOCMCoreShareResponse{
		Status:  status.NewOK(ctx),
		Id:      share.Id.OpaqueId,
		Created: share.Ctime,
	}
	return res, nil
}

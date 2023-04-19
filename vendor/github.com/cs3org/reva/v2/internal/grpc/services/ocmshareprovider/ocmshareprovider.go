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

package ocmshareprovider

import (
	"context"

	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
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
	rgrpc.Register("ocmshareprovider", New)
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
	ocm.RegisterOcmAPIServer(ss, s)
}

func getShareManager(c *config) (share.Manager, error) {
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

// New creates a new ocm share provider svc
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
	return []string{}
}

// Note: this is for outgoing OCM shares
// This function is used when you for instance
// call `ocm-share-create` in reva-cli.
// For incoming OCM shares from internal/http/services/ocmd/shares.go
// there is the very similar but slightly different function
// CreateOCMCoreShare (the "Core" somehow means "incoming").
// So make sure to keep in mind the difference between this file for outgoing:
// internal/grpc/services/ocmshareprovider/ocmshareprovider.go
// and the other one for incoming:
// internal/grpc/service/ocmcore/ocmcore.go
// Both functions end up calling the same s.sm.Share function
// on the OCM share manager:
// pkg/ocm/share/manager/{json|nextcloud|...}
func (s *service) CreateOCMShare(ctx context.Context, req *ocm.CreateOCMShareRequest) (*ocm.CreateOCMShareResponse, error) {
	if req.Opaque == nil {
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "can't find resource permissions"),
		}, nil
	}

	var permissions string
	permOpaque, ok := req.Opaque.Map["permissions"]
	if !ok {
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "resource permissions not set"),
		}, nil
	}
	switch permOpaque.Decoder {
	case "plain":
		permissions = string(permOpaque.Value)
	default:
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "invalid opaque entry decoder"),
		}, nil
	}

	var name string
	nameOpaque, ok := req.Opaque.Map["name"]
	if !ok {
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "resource name not set"),
		}, nil
	}
	switch nameOpaque.Decoder {
	case "plain":
		name = string(nameOpaque.Value)
	default:
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "invalid opaque entry decoder"),
		}, nil
	}

	// discover share type
	sharetype := ocm.Share_SHARE_TYPE_REGULAR
	// FIXME: https://github.com/cs3org/reva/issues/2402
	protocol, ok := req.Opaque.Map["protocol"]
	if ok {
		switch protocol.Decoder {
		case "plain":
			if string(protocol.Value) == "datatx" {
				sharetype = ocm.Share_SHARE_TYPE_TRANSFER
			}
		default:
			return &ocm.CreateOCMShareResponse{
				Status: status.NewInternal(ctx, "error creating share"),
			}, nil
		}
		// token = protocol FIXME!
	}

	var sharedSecret string
	share, err := s.sm.Share(ctx, req.ResourceId, req.Grant, name, req.RecipientMeshProvider, permissions, nil, sharedSecret, sharetype)

	if err != nil {
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "error creating share"),
		}, nil
	}

	res := &ocm.CreateOCMShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

func (s *service) RemoveOCMShare(ctx context.Context, req *ocm.RemoveOCMShareRequest) (*ocm.RemoveOCMShareResponse, error) {
	err := s.sm.Unshare(ctx, req.Ref)
	if err != nil {
		return &ocm.RemoveOCMShareResponse{
			Status: status.NewInternal(ctx, "error removing share"),
		}, nil
	}

	return &ocm.RemoveOCMShareResponse{
		Status: status.NewOK(ctx),
	}, nil
}

func (s *service) GetOCMShare(ctx context.Context, req *ocm.GetOCMShareRequest) (*ocm.GetOCMShareResponse, error) {
	share, err := s.sm.GetShare(ctx, req.Ref)
	if err != nil {
		return &ocm.GetOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting share"),
		}, nil
	}

	return &ocm.GetOCMShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}, nil
}

func (s *service) ListOCMShares(ctx context.Context, req *ocm.ListOCMSharesRequest) (*ocm.ListOCMSharesResponse, error) {
	shares, err := s.sm.ListShares(ctx, req.Filters) // TODO(labkode): add filter to share manager
	if err != nil {
		return &ocm.ListOCMSharesResponse{
			Status: status.NewInternal(ctx, "error listing shares"),
		}, nil
	}

	res := &ocm.ListOCMSharesResponse{
		Status: status.NewOK(ctx),
		Shares: shares,
	}
	return res, nil
}

func (s *service) UpdateOCMShare(ctx context.Context, req *ocm.UpdateOCMShareRequest) (*ocm.UpdateOCMShareResponse, error) {
	_, err := s.sm.UpdateShare(ctx, req.Ref, req.Field.GetPermissions()) // TODO(labkode): check what to update
	if err != nil {
		return &ocm.UpdateOCMShareResponse{
			Status: status.NewInternal(ctx, "error updating share"),
		}, nil
	}

	res := &ocm.UpdateOCMShareResponse{
		Status: status.NewOK(ctx),
	}
	return res, nil
}

func (s *service) ListReceivedOCMShares(ctx context.Context, req *ocm.ListReceivedOCMSharesRequest) (*ocm.ListReceivedOCMSharesResponse, error) {
	shares, err := s.sm.ListReceivedShares(ctx)
	if err != nil {
		return &ocm.ListReceivedOCMSharesResponse{
			Status: status.NewInternal(ctx, "error listing received shares"),
		}, nil
	}

	res := &ocm.ListReceivedOCMSharesResponse{
		Status: status.NewOK(ctx),
		Shares: shares,
	}
	return res, nil
}

func (s *service) UpdateReceivedOCMShare(ctx context.Context, req *ocm.UpdateReceivedOCMShareRequest) (*ocm.UpdateReceivedOCMShareResponse, error) {
	_, err := s.sm.UpdateReceivedShare(ctx, req.Share, req.UpdateMask) // TODO(labkode): check what to update
	if err != nil {
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: status.NewInternal(ctx, "error updating received share"),
		}, nil
	}

	res := &ocm.UpdateReceivedOCMShareResponse{
		Status: status.NewOK(ctx),
	}
	return res, nil
}

func (s *service) GetReceivedOCMShare(ctx context.Context, req *ocm.GetReceivedOCMShareRequest) (*ocm.GetReceivedOCMShareResponse, error) {
	share, err := s.sm.GetReceivedShare(ctx, req.Ref)
	if err != nil {
		return &ocm.GetReceivedOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting received share"),
		}, nil
	}

	res := &ocm.GetReceivedOCMShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

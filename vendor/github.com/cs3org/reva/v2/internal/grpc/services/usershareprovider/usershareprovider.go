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

package usershareprovider

import (
	"context"
	"regexp"
	"slices"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/conversions"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/permission"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/manager/registry"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/cs3org/reva/v2/pkg/utils"
)

func init() {
	rgrpc.Register("usershareprovider", NewDefault)
}

type config struct {
	Driver                string                            `mapstructure:"driver"`
	Drivers               map[string]map[string]interface{} `mapstructure:"drivers"`
	GatewayAddr           string                            `mapstructure:"gateway_addr"`
	AllowedPathsForShares []string                          `mapstructure:"allowed_paths_for_shares"`
}

func (c *config) init() {
	if c.Driver == "" {
		c.Driver = "json"
	}
}

type service struct {
	sm                    share.Manager
	gatewaySelector       pool.Selectable[gateway.GatewayAPIClient]
	allowedPathsForShares []*regexp.Regexp
}

func getShareManager(c *config) (share.Manager, error) {
	if f, ok := registry.NewFuncs[c.Driver]; ok {
		return f(c.Drivers[c.Driver])
	}
	return nil, errtypes.NotFound("driver not found: " + c.Driver)
}

// TODO(labkode): add ctx to Close.
func (s *service) Close() error {
	return nil
}

func (s *service) UnprotectedEndpoints() []string {
	return []string{}
}

func (s *service) Register(ss *grpc.Server) {
	collaboration.RegisterCollaborationAPIServer(ss, s)
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New creates a new user share provider svc initialized from defaults
func NewDefault(m map[string]interface{}, ss *grpc.Server) (rgrpc.Service, error) {

	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	c.init()

	sm, err := getShareManager(c)
	if err != nil {
		return nil, err
	}

	allowedPathsForShares := make([]*regexp.Regexp, 0, len(c.AllowedPathsForShares))
	for _, s := range c.AllowedPathsForShares {
		regex, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		allowedPathsForShares = append(allowedPathsForShares, regex)
	}

	gatewaySelector, err := pool.GatewaySelector(sharedconf.GetGatewaySVC(c.GatewayAddr))
	if err != nil {
		return nil, err
	}

	return New(gatewaySelector, sm, allowedPathsForShares), nil
}

// New creates a new user share provider svc
func New(gatewaySelector pool.Selectable[gateway.GatewayAPIClient], sm share.Manager, allowedPathsForShares []*regexp.Regexp) rgrpc.Service {
	service := &service{
		sm:                    sm,
		gatewaySelector:       gatewaySelector,
		allowedPathsForShares: allowedPathsForShares,
	}

	return service
}

func (s *service) isPathAllowed(path string) bool {
	if len(s.allowedPathsForShares) == 0 {
		return true
	}
	for _, reg := range s.allowedPathsForShares {
		if reg.MatchString(path) {
			return true
		}
	}
	return false
}

func (s *service) CreateShare(ctx context.Context, req *collaboration.CreateShareRequest) (*collaboration.CreateShareResponse, error) {
	log := appctx.GetLogger(ctx)
	user := ctxpkg.ContextMustGetUser(ctx)

	// Grants must not allow grant permissions
	if HasGrantPermissions(req.GetGrant().GetPermissions().GetPermissions()) {
		return &collaboration.CreateShareResponse{
			Status: status.NewInvalidArg(ctx, "resharing not supported"),
		}, nil
	}

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	// check if the user has the permission to create shares at all
	ok, err := utils.CheckPermission(ctx, permission.WriteShare, gatewayClient)
	if err != nil {
		return &collaboration.CreateShareResponse{
			Status: status.NewInternal(ctx, "failed check user permission to write public link"),
		}, err
	}
	if !ok {
		return &collaboration.CreateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "no permission to create public links"),
		}, nil
	}

	if req.GetGrant().GetGrantee().GetType() == provider.GranteeType_GRANTEE_TYPE_USER && req.GetGrant().GetGrantee().GetUserId().GetIdp() == "" {
		// use logged in user Idp as default.
		req.GetGrant().GetGrantee().Id = &provider.Grantee_UserId{
			UserId: &userpb.UserId{
				OpaqueId: req.GetGrant().GetGrantee().GetUserId().GetOpaqueId(),
				Idp:      user.GetId().GetIdp(),
				Type:     userpb.UserType_USER_TYPE_PRIMARY},
		}
	}

	sRes, err := gatewayClient.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: req.GetResourceInfo().GetId()}})
	if err != nil {
		log.Err(err).Interface("resource_id", req.GetResourceInfo().GetId()).Msg("failed to stat resource to share")
		return &collaboration.CreateShareResponse{
			Status: status.NewInternal(ctx, "failed to stat shared resource"),
		}, err
	}
	// the user needs to have the AddGrant permissions on the Resource to be able to create a share
	if !sRes.GetInfo().GetPermissionSet().AddGrant {
		return &collaboration.CreateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "no permission to add grants on shared resource"),
		}, err
	}
	// check if the share creator has sufficient permissions to do so.
	if shareCreationAllowed := conversions.SufficientCS3Permissions(
		sRes.GetInfo().GetPermissionSet(),
		req.GetGrant().GetPermissions().GetPermissions(),
	); !shareCreationAllowed {
		return &collaboration.CreateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "insufficient permissions to create that kind of share"),
		}, nil
	}
	// check if the requested permission are plausible for the Resource
	if sRes.GetInfo().GetType() == provider.ResourceType_RESOURCE_TYPE_FILE {
		if newPermissions := req.GetGrant().GetPermissions().GetPermissions(); newPermissions.GetCreateContainer() || newPermissions.GetMove() || newPermissions.GetDelete() {
			return &collaboration.CreateShareResponse{
				Status: status.NewInvalid(ctx, "cannot set the requested permissions on that type of resource"),
			}, nil
		}
	}

	if !s.isPathAllowed(req.GetResourceInfo().GetPath()) {
		return &collaboration.CreateShareResponse{
			Status: status.NewFailedPrecondition(ctx, nil, "share creation is not allowed for the specified path"),
		}, nil
	}

	createdShare, err := s.sm.Share(ctx, req.GetResourceInfo(), req.GetGrant())
	if err != nil {
		return &collaboration.CreateShareResponse{
			Status: status.NewStatusFromErrType(ctx, "error creating share", err),
		}, nil
	}

	return &collaboration.CreateShareResponse{
		Status: status.NewOK(ctx),
		Share:  createdShare,
	}, nil
}

func HasGrantPermissions(p *provider.ResourcePermissions) bool {
	return p.GetAddGrant() || p.GetUpdateGrant() || p.GetRemoveGrant() || p.GetDenyGrant()
}

func (s *service) RemoveShare(ctx context.Context, req *collaboration.RemoveShareRequest) (*collaboration.RemoveShareResponse, error) {
	log := appctx.GetLogger(ctx)
	user := ctxpkg.ContextMustGetUser(ctx)
	share, err := s.sm.GetShare(ctx, req.Ref)
	if err != nil {
		return &collaboration.RemoveShareResponse{
			Status: status.NewInternal(ctx, "error getting share"),
		}, nil
	}

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	sRes, err := gatewayClient.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: share.GetResourceId()}})
	if err != nil {
		log.Err(err).Interface("resource_id", share.GetResourceId()).Msg("failed to stat shared resource")
		return &collaboration.RemoveShareResponse{
			Status: status.NewInternal(ctx, "failed to stat shared resource"),
		}, err
	}
	// the requesting user needs to be either the Owner/Creator of the share or have the RemoveGrant permissions on the Resource
	switch {
	case utils.UserEqual(user.GetId(), share.GetCreator()) || utils.UserEqual(user.GetId(), share.GetOwner()):
		fallthrough
	case sRes.GetInfo().GetPermissionSet().RemoveGrant:
		break
	default:
		return &collaboration.RemoveShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "no permission to remove grants on shared resource"),
		}, err
	}

	err = s.sm.Unshare(ctx, req.Ref)
	if err != nil {
		return &collaboration.RemoveShareResponse{
			Status: status.NewInternal(ctx, "error removing share"),
		}, nil
	}

	o := utils.AppendJSONToOpaque(nil, "resourceid", share.GetResourceId())
	if user := share.GetGrantee().GetUserId(); user != nil {
		o = utils.AppendJSONToOpaque(o, "granteeuserid", user)
	} else {
		o = utils.AppendJSONToOpaque(o, "granteegroupid", share.GetGrantee().GetGroupId())
	}

	return &collaboration.RemoveShareResponse{
		Opaque: o,
		Status: status.NewOK(ctx),
	}, nil
}

func (s *service) GetShare(ctx context.Context, req *collaboration.GetShareRequest) (*collaboration.GetShareResponse, error) {
	share, err := s.sm.GetShare(ctx, req.Ref)
	if err != nil {
		var st *rpc.Status
		switch err.(type) {
		case errtypes.IsNotFound:
			st = status.NewNotFound(ctx, err.Error())
		default:
			st = status.NewInternal(ctx, err.Error())
		}
		return &collaboration.GetShareResponse{
			Status: st,
		}, nil
	}

	return &collaboration.GetShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}, nil
}

func (s *service) ListShares(ctx context.Context, req *collaboration.ListSharesRequest) (*collaboration.ListSharesResponse, error) {
	shares, err := s.sm.ListShares(ctx, req.Filters) // TODO(labkode): add filter to share manager
	if err != nil {
		return &collaboration.ListSharesResponse{
			Status: status.NewInternal(ctx, "error listing shares"),
		}, nil
	}

	res := &collaboration.ListSharesResponse{
		Status: status.NewOK(ctx),
		Shares: shares,
	}
	return res, nil
}

func (s *service) UpdateShare(ctx context.Context, req *collaboration.UpdateShareRequest) (*collaboration.UpdateShareResponse, error) {
	log := appctx.GetLogger(ctx)
	user := ctxpkg.ContextMustGetUser(ctx)

	// Grants must not allow grant permissions
	if HasGrantPermissions(req.GetShare().GetPermissions().GetPermissions()) {
		return &collaboration.UpdateShareResponse{
			Status: status.NewInvalidArg(ctx, "resharing not supported"),
		}, nil
	}

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	// check if the user has the permission to create shares at all
	ok, err := utils.CheckPermission(ctx, permission.WriteShare, gatewayClient)
	if err != nil {
		return &collaboration.UpdateShareResponse{
			Status: status.NewInternal(ctx, "failed check user permission to write share"),
		}, err
	}
	if !ok {
		return &collaboration.UpdateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "no permission to create user share"),
		}, nil
	}

	// Read share from backend. We need the shared resource's id for STATing it, it might not be in
	// the incoming request
	currentShare, err := s.sm.GetShare(ctx,
		&collaboration.ShareReference{
			Spec: &collaboration.ShareReference_Id{
				Id: req.GetShare().GetId(),
			},
		},
	)
	if err != nil {
		var st *rpc.Status
		switch err.(type) {
		case errtypes.IsNotFound:
			st = status.NewNotFound(ctx, err.Error())
		default:
			st = status.NewInternal(ctx, err.Error())
		}
		return &collaboration.UpdateShareResponse{
			Status: st,
		}, nil
	}

	sRes, err := gatewayClient.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{ResourceId: currentShare.GetResourceId()}})
	if err != nil {
		log.Err(err).Interface("resource_id", req.GetShare().GetResourceId()).Msg("failed to stat resource to share")
		return &collaboration.UpdateShareResponse{
			Status: status.NewInternal(ctx, "failed to stat shared resource"),
		}, err
	}
	// the requesting user needs to be either the Owner/Creator of the share or have the UpdateGrant permissions on the Resource
	switch {
	case utils.UserEqual(user.GetId(), currentShare.GetCreator()) || utils.UserEqual(user.GetId(), currentShare.GetOwner()):
		fallthrough
	case sRes.GetInfo().GetPermissionSet().UpdateGrant:
		break
	default:
		return &collaboration.UpdateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "no permission to remove grants on shared resource"),
		}, err
	}

	// If this is a permissions update, check if user's permissions on the resource are sufficient to set the desired permissions
	var newPermissions *provider.ResourcePermissions
	if slices.Contains(req.GetUpdateMask().GetPaths(), "permissions") {
		newPermissions = req.GetShare().GetPermissions().GetPermissions()
	} else {
		newPermissions = req.GetField().GetPermissions().GetPermissions()
	}
	if newPermissions != nil && !conversions.SufficientCS3Permissions(sRes.GetInfo().GetPermissionSet(), newPermissions) {
		return &collaboration.UpdateShareResponse{
			Status: status.NewPermissionDenied(ctx, nil, "insufficient permissions to create that kind of share"),
		}, nil
	}

	// check if the requested permission are plausible for the Resource
	// do we need more here?
	if sRes.GetInfo().GetType() == provider.ResourceType_RESOURCE_TYPE_FILE {
		if newPermissions.GetCreateContainer() || newPermissions.GetMove() || newPermissions.GetDelete() {
			return &collaboration.UpdateShareResponse{
				Status: status.NewInvalid(ctx, "cannot set the requested permissions on that type of resource"),
			}, nil
		}
	}

	share, err := s.sm.UpdateShare(ctx, req.Ref, req.Field.GetPermissions(), req.Share, req.UpdateMask) // TODO(labkode): check what to update
	if err != nil {
		return &collaboration.UpdateShareResponse{
			Status: status.NewInternal(ctx, "error updating share"),
		}, nil
	}

	res := &collaboration.UpdateShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

func (s *service) ListReceivedShares(ctx context.Context, req *collaboration.ListReceivedSharesRequest) (*collaboration.ListReceivedSharesResponse, error) {
	// For the UI add a filter to not display the denial shares
	foundExclude := false
	for _, f := range req.Filters {
		if f.Type == collaboration.Filter_TYPE_EXCLUDE_DENIALS {
			foundExclude = true
			break
		}
	}
	if !foundExclude {
		req.Filters = append(req.Filters, &collaboration.Filter{Type: collaboration.Filter_TYPE_EXCLUDE_DENIALS})
	}

	var uid userpb.UserId
	_ = utils.ReadJSONFromOpaque(req.Opaque, "userid", &uid)
	shares, err := s.sm.ListReceivedShares(ctx, req.Filters, &uid) // TODO(labkode): check what to update
	if err != nil {
		return &collaboration.ListReceivedSharesResponse{
			Status: status.NewInternal(ctx, "error listing received shares"),
		}, nil
	}

	res := &collaboration.ListReceivedSharesResponse{
		Status: status.NewOK(ctx),
		Shares: shares,
	}
	return res, nil
}

func (s *service) GetReceivedShare(ctx context.Context, req *collaboration.GetReceivedShareRequest) (*collaboration.GetReceivedShareResponse, error) {
	log := appctx.GetLogger(ctx)

	share, err := s.sm.GetReceivedShare(ctx, req.Ref)
	if err != nil {
		log.Err(err).Msg("error getting received share")
		switch err.(type) {
		case errtypes.NotFound:
			return &collaboration.GetReceivedShareResponse{
				Status: status.NewNotFound(ctx, "error getting received share"),
			}, nil
		default:
			return &collaboration.GetReceivedShareResponse{
				Status: status.NewInternal(ctx, "error getting received share"),
			}, nil
		}
	}

	res := &collaboration.GetReceivedShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

func (s *service) UpdateReceivedShare(ctx context.Context, req *collaboration.UpdateReceivedShareRequest) (*collaboration.UpdateReceivedShareResponse, error) {

	if req.Share == nil {
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalid(ctx, "updating requires a received share object"),
		}, nil
	}
	if req.Share.Share == nil {
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalid(ctx, "share missing"),
		}, nil
	}
	if req.Share.Share.Id == nil {
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalid(ctx, "share id missing"),
		}, nil
	}
	if req.Share.Share.Id.OpaqueId == "" {
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalid(ctx, "share id empty"),
		}, nil
	}

	var uid userpb.UserId
	_ = utils.ReadJSONFromOpaque(req.Opaque, "userid", &uid)
	share, err := s.sm.UpdateReceivedShare(ctx, req.Share, req.UpdateMask, &uid)
	if err != nil {
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInternal(ctx, "error updating received share"),
		}, nil
	}

	res := &collaboration.UpdateReceivedShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

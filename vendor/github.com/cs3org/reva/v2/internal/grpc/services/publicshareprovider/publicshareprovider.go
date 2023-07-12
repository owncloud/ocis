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

package publicshareprovider

import (
	"context"
	"regexp"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	"github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func init() {
	rgrpc.Register("publicshareprovider", New)
}

type config struct {
	Driver                         string                            `mapstructure:"driver"`
	Drivers                        map[string]map[string]interface{} `mapstructure:"drivers"`
	AllowedPathsForShares          []string                          `mapstructure:"allowed_paths_for_shares"`
	WriteableShareMustHavePassword bool                              `mapstructure:"writeable_share_must_have_password"`
}

func (c *config) init() {
	if c.Driver == "" {
		c.Driver = "json"
	}
}

type service struct {
	conf                  *config
	sm                    publicshare.Manager
	allowedPathsForShares []*regexp.Regexp
}

func getShareManager(c *config) (publicshare.Manager, error) {
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
	return []string{"/cs3.sharing.link.v1beta1.LinkAPI/GetPublicShareByToken"}
}

func (s *service) Register(ss *grpc.Server) {
	link.RegisterLinkAPIServer(ss, s)
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New creates a new user share provider svc
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

	allowedPathsForShares := make([]*regexp.Regexp, 0, len(c.AllowedPathsForShares))
	for _, s := range c.AllowedPathsForShares {
		regex, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		allowedPathsForShares = append(allowedPathsForShares, regex)
	}

	service := &service{
		conf:                  c,
		sm:                    sm,
		allowedPathsForShares: allowedPathsForShares,
	}

	return service, nil
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

func (s *service) CreatePublicShare(ctx context.Context, req *link.CreatePublicShareRequest) (*link.CreatePublicShareResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Info().Str("publicshareprovider", "create").Msg("create public share")

	if !s.isPathAllowed(req.ResourceInfo.Path) {
		return &link.CreatePublicShareResponse{
			Status: status.NewInvalid(ctx, "share creation is not allowed for the specified path"),
		}, nil
	}

	grant := req.GetGrant()
	if grant != nil && s.conf.WriteableShareMustHavePassword &&
		publicshare.IsWriteable(grant.GetPermissions()) && grant.Password == "" {
		return &link.CreatePublicShareResponse{
			Status: status.NewInvalid(ctx, "writeable shares must have a password protection"),
		}, nil
	}

	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		log.Error().Msg("error getting user from context")
	}

	share, err := s.sm.CreatePublicShare(ctx, u, req.ResourceInfo, req.Grant)
	if err != nil {
		log.Debug().Err(err).Str("createShare", "shares").Msg("error connecting to storage provider")
	}

	res := &link.CreatePublicShareResponse{
		Status: status.NewOK(ctx),
		Share:  share,
	}
	return res, nil
}

func (s *service) RemovePublicShare(ctx context.Context, req *link.RemovePublicShareRequest) (*link.RemovePublicShareResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Info().Str("publicshareprovider", "remove").Msg("remove public share")

	user := ctxpkg.ContextMustGetUser(ctx)
	err := s.sm.RevokePublicShare(ctx, user, req.Ref)
	if err != nil {
		return &link.RemovePublicShareResponse{
			Status: status.NewInternal(ctx, "error deleting public share"),
		}, err
	}
	return &link.RemovePublicShareResponse{
		Status: status.NewOK(ctx),
	}, nil
}

func (s *service) GetPublicShareByToken(ctx context.Context, req *link.GetPublicShareByTokenRequest) (*link.GetPublicShareByTokenResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Debug().Msg("getting public share by token")

	// there are 2 passes here, and the second request has no password
	found, err := s.sm.GetPublicShareByToken(ctx, req.GetToken(), req.GetAuthentication(), req.GetSign())
	switch v := err.(type) {
	case nil:
		return &link.GetPublicShareByTokenResponse{
			Status: status.NewOK(ctx),
			Share:  found,
		}, nil
	case errtypes.InvalidCredentials:
		return &link.GetPublicShareByTokenResponse{
			Status: status.NewPermissionDenied(ctx, v, "wrong password"),
		}, nil
	case errtypes.NotFound:
		return &link.GetPublicShareByTokenResponse{
			Status: status.NewNotFound(ctx, "unknown token"),
		}, nil
	default:
		return &link.GetPublicShareByTokenResponse{
			Status: status.NewInternal(ctx, "unexpected error"),
		}, nil
	}
}

func (s *service) GetPublicShare(ctx context.Context, req *link.GetPublicShareRequest) (*link.GetPublicShareResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Info().Str("publicshareprovider", "get").Msg("get public share")

	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		log.Error().Msg("error getting user from context")
	}

	ps, err := s.sm.GetPublicShare(ctx, u, req.Ref, req.GetSign())
	switch {
	case err != nil:
		var st *rpc.Status
		switch err.(type) {
		case errtypes.IsNotFound:
			st = status.NewNotFound(ctx, err.Error())
		default:
			st = status.NewInternal(ctx, err.Error())
		}
		return &link.GetPublicShareResponse{
			Status: st,
		}, nil
	case ps == nil:
		return &link.GetPublicShareResponse{
			Status: status.NewNotFound(ctx, "not found"),
		}, nil
	default:
		return &link.GetPublicShareResponse{
			Status: status.NewOK(ctx),
			Share:  ps,
		}, nil
	}
}

func (s *service) ListPublicShares(ctx context.Context, req *link.ListPublicSharesRequest) (*link.ListPublicSharesResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Info().Str("publicshareprovider", "list").Msg("list public share")
	user, _ := ctxpkg.ContextGetUser(ctx)

	shares, err := s.sm.ListPublicShares(ctx, user, req.Filters, req.GetSign())
	if err != nil {
		log.Err(err).Msg("error listing shares")
		return &link.ListPublicSharesResponse{
			Status: status.NewInternal(ctx, "error listing public shares"),
		}, nil
	}

	res := &link.ListPublicSharesResponse{
		Status: status.NewOK(ctx),
		Share:  shares,
	}
	return res, nil
}

func (s *service) UpdatePublicShare(ctx context.Context, req *link.UpdatePublicShareRequest) (*link.UpdatePublicShareResponse, error) {
	log := appctx.GetLogger(ctx)
	log.Info().Str("publicshareprovider", "update").Msg("update public share")

	u, ok := ctxpkg.ContextGetUser(ctx)
	if !ok {
		log.Error().Msg("error getting user from context")
	}

	updateR, err := s.sm.UpdatePublicShare(ctx, u, req)
	if err != nil {
		if errors.Is(err, publicshare.ErrShareNeedsPassword) {
			return &link.UpdatePublicShareResponse{
				Status: status.NewInvalid(ctx, err.Error()),
			}, nil
		}
		return &link.UpdatePublicShareResponse{
			Status: status.NewInternal(ctx, err.Error()),
		}, nil
	}

	res := &link.UpdatePublicShareResponse{
		Status: status.NewOK(ctx),
		Share:  updateR,
	}
	return res, nil
}

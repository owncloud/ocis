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

package gateway

import (
	"context"
	"path"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storage/utils/grants"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

// TODO(labkode): add multi-phase commit logic when commit share or commit ref is enabled.
func (s *svc) CreateShare(ctx context.Context, req *collaboration.CreateShareRequest) (*collaboration.CreateShareResponse, error) {
	// Don't use the share manager when sharing a space root
	if !s.c.UseCommonSpaceRootShareLogic && refIsSpaceRoot(req.ResourceInfo.Id) {
		return s.addSpaceShare(ctx, req)
	}
	return s.addShare(ctx, req)
}

func (s *svc) RemoveShare(ctx context.Context, req *collaboration.RemoveShareRequest) (*collaboration.RemoveShareResponse, error) {
	key := req.Ref.GetKey()
	if !s.c.UseCommonSpaceRootShareLogic && shareIsSpaceRoot(key) {
		return s.removeSpaceShare(ctx, key.ResourceId, key.Grantee)
	}
	return s.removeShare(ctx, req)
}

// TODO(labkode): we need to validate share state vs storage grant and storage ref
// If there are any inconsistencies, the share needs to be flag as invalid and a background process
// or active fix needs to be performed.
func (s *svc) GetShare(ctx context.Context, req *collaboration.GetShareRequest) (*collaboration.GetShareResponse, error) {
	return s.getShare(ctx, req)
}

func (s *svc) getShare(ctx context.Context, req *collaboration.GetShareRequest) (*collaboration.GetShareResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("getShare: failed to get user share provider")
		return &collaboration.GetShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.GetShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetShare")
	}

	return res, nil
}

// TODO(labkode): read GetShare comment.
func (s *svc) ListShares(ctx context.Context, req *collaboration.ListSharesRequest) (*collaboration.ListSharesResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("ListShares: failed to get user share provider")
		return &collaboration.ListSharesResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.ListShares(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling ListShares")
	}

	return res, nil
}

func (s *svc) UpdateShare(ctx context.Context, req *collaboration.UpdateShareRequest) (*collaboration.UpdateShareResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("UpdateShare: failed to get user share provider")
		return &collaboration.UpdateShareResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}
	res, err := c.UpdateShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling UpdateShare")
	}

	// if we don't need to commit we return earlier
	if !s.c.CommitShareToStorageGrant && !s.c.CommitShareToStorageRef {
		return res, nil
	}

	// TODO(labkode): if both commits are enabled they could be done concurrently.

	if s.c.CommitShareToStorageGrant {
		updateGrantStatus, err := s.updateGrant(ctx, res.GetShare().GetResourceId(),
			res.GetShare().GetGrantee(),
			res.GetShare().GetPermissions().GetPermissions())

		if err != nil {
			return nil, errors.Wrap(err, "gateway: error calling updateGrant")
		}

		if updateGrantStatus.Code != rpc.Code_CODE_OK {
			return &collaboration.UpdateShareResponse{
				Status: updateGrantStatus,
				Share:  res.Share,
			}, nil
		}
	}

	s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), res.Share.ResourceId)
	return res, nil
}

// TODO(labkode): listing received shares just goes to the user share manager and gets the list of
// received shares. The display name of the shares should be the a friendly name, like the basename
// of the original file.
func (s *svc) ListReceivedShares(ctx context.Context, req *collaboration.ListReceivedSharesRequest) (*collaboration.ListReceivedSharesResponse, error) {
	logger := appctx.GetLogger(ctx)
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("ListReceivedShares: failed to get user share provider")
		return &collaboration.ListReceivedSharesResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	res, err := c.ListReceivedShares(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling ListReceivedShares")
	}
	return res, nil
}

func (s *svc) GetReceivedShare(ctx context.Context, req *collaboration.GetReceivedShareRequest) (*collaboration.GetReceivedShareResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("GetReceivedShare: failed to get user share provider")
		return &collaboration.GetReceivedShareResponse{
			Status: status.NewInternal(ctx, "error getting received share"),
		}, nil
	}

	res, err := c.GetReceivedShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetReceivedShare")
	}

	return res, nil
}

// When updating a received share:
// if the update contains update for displayName:
//   1) if received share is mounted: we also do a rename in the storage
//   2) if received share is not mounted: we only rename in user share provider.
func (s *svc) UpdateReceivedShare(ctx context.Context, req *collaboration.UpdateReceivedShareRequest) (*collaboration.UpdateReceivedShareResponse, error) {
	ctx, span := appctx.GetTracerProvider(ctx).Tracer("gateway").Start(ctx, "Gateway.UpdateReceivedShare")
	defer span.End()

	// sanity checks
	switch {
	case req.GetShare() == nil:
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalidArg(ctx, "updating requires a received share object"),
		}, nil
	case req.GetShare().GetShare() == nil:
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalidArg(ctx, "share missing"),
		}, nil
	case req.GetShare().GetShare().GetId() == nil:
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalidArg(ctx, "share id missing"),
		}, nil
	case req.GetShare().GetShare().GetId().GetOpaqueId() == "":
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInvalidArg(ctx, "share id empty"),
		}, nil
	}

	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("UpdateReceivedShare: failed to get user share provider")
		return &collaboration.UpdateReceivedShareResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), req.Share.Share.ResourceId)
	return c.UpdateReceivedShare(ctx, req)
	/*
		    TODO: Leftover from master merge. Do we need this?
			if err != nil {
				appctx.GetLogger(ctx).
					Err(err).
					Msg("UpdateReceivedShare: failed to get user share provider")
				return &collaboration.UpdateReceivedShareResponse{
					Status: status.NewInternal(ctx, "error getting share provider client"),
				}, nil
			}
			// check if we have a resource id in the update response that we can use to update references
			if res.GetShare().GetShare().GetResourceId() == nil {
				log.Err(err).Msg("gateway: UpdateReceivedShare must return a ResourceId")
				return &collaboration.UpdateReceivedShareResponse{
					Status: &rpc.Status{
						Code: rpc.Code_CODE_INTERNAL,
					},
				}, nil
			}

			// properties are updated in the order they appear in the field mask
			// when an error occurs the request ends and no further fields are updated
			for i := range req.UpdateMask.Paths {
				switch req.UpdateMask.Paths[i] {
				case "state":
					switch req.GetShare().GetState() {
					case collaboration.ShareState_SHARE_STATE_ACCEPTED:
						rpcStatus := s.createReference(ctx, res.GetShare().GetShare().GetResourceId())
						if rpcStatus.Code != rpc.Code_CODE_OK {
							return &collaboration.UpdateReceivedShareResponse{Status: rpcStatus}, nil
						}
					case collaboration.ShareState_SHARE_STATE_REJECTED:
						rpcStatus := s.removeReference(ctx, res.GetShare().GetShare().ResourceId)
						if rpcStatus.Code != rpc.Code_CODE_OK && rpcStatus.Code != rpc.Code_CODE_NOT_FOUND {
							return &collaboration.UpdateReceivedShareResponse{Status: rpcStatus}, nil
						}
					}
				case "mount_point":
					// TODO(labkode): implementing updating mount point
					err = errtypes.NotSupported("gateway: update of mount point is not yet implemented")
					return &collaboration.UpdateReceivedShareResponse{
						Status: status.NewUnimplemented(ctx, err, "error updating received share"),
					}, nil
				default:
					return nil, errtypes.NotSupported("updating " + req.UpdateMask.Paths[i] + " is not supported")
				}
			}
			return res, nil
	*/
}

func (s *svc) removeReference(ctx context.Context, resourceID *provider.ResourceId) *rpc.Status {
	log := appctx.GetLogger(ctx)

	idReference := &provider.Reference{ResourceId: resourceID}
	storageProvider, _, err := s.find(ctx, idReference)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", idReference).
			Msg("removeReference: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found")
		}
		return status.NewInternal(ctx, "error finding storage provider")
	}

	statRes, err := storageProvider.Stat(ctx, &provider.StatRequest{Ref: idReference})
	if err != nil {
		log.Error().Err(err).Interface("reference", idReference).Msg("removeReference: error calling Stat")
		return status.NewInternal(ctx, "gateway: error calling Stat for the share resource id: "+resourceID.String())
	}

	// FIXME how can we delete a reference if the original resource was deleted?
	if statRes.Status.Code != rpc.Code_CODE_OK {
		log.Error().Interface("status", statRes.Status).Interface("reference", idReference).Msg("removeReference: error calling Stat")
		return status.NewInternal(ctx, "could not delete share reference")
	}

	homeRes, err := s.GetHome(ctx, &provider.GetHomeRequest{})
	if err != nil {
		return status.NewInternal(ctx, "could not delete share reference")
	}

	sharePath := path.Join(homeRes.Path, s.c.ShareFolder, path.Base(statRes.Info.Path))
	log.Debug().Str("share_path", sharePath).Msg("remove reference of share")

	sharePathRef := &provider.Reference{Path: sharePath}
	homeProvider, providerInfo, err := s.find(ctx, sharePathRef)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", sharePathRef).
			Msg("removeReference: failed to get storage provider for share ref")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found")
		}
		return status.NewInternal(ctx, "error finding storage provider")
	}

	var (
		root      *provider.ResourceId
		mountPath string
	)
	for _, space := range decodeSpaces(providerInfo) {
		mountPath = decodePath(space)
		root = space.Root
		break // TODO can there be more than one space for a path?
	}

	ref := unwrap(sharePathRef, mountPath, root)

	deleteReq := &provider.DeleteRequest{
		Opaque: &typesv1beta1.Opaque{
			Map: map[string]*typesv1beta1.OpaqueEntry{
				// This signals the storageprovider that we want to delete the share reference and not the underlying file.
				"deleting_shared_resource": {},
			},
		},
		Ref: ref,
	}

	deleteResp, err := homeProvider.Delete(ctx, deleteReq)
	if err != nil {
		return status.NewInternal(ctx, "could not delete share reference")
	}

	switch deleteResp.Status.Code {
	case rpc.Code_CODE_OK:
		// we can continue deleting the reference
	case rpc.Code_CODE_NOT_FOUND:
		// This is fine, we wanted to delete it anyway
		return status.NewOK(ctx)
	default:
		return status.NewInternal(ctx, "could not delete share reference")
	}

	log.Debug().Str("share_path", sharePath).Msg("share reference successfully removed")

	return status.NewOK(ctx)
}

func (s *svc) denyGrant(ctx context.Context, id *provider.ResourceId, g *provider.Grantee, opaque *typesv1beta1.Opaque) (*rpc.Status, error) {
	ref := &provider.Reference{
		ResourceId: id,
	}

	grantReq := &provider.DenyGrantRequest{
		Ref:     ref,
		Grantee: g,
		Opaque:  opaque,
	}

	c, _, err := s.find(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", ref).
			Msg("denyGrant: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found"), nil
		}
		return status.NewInternal(ctx, "error finding storage provider"), nil
	}

	grantRes, err := c.DenyGrant(ctx, grantReq)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling DenyGrant")
	}
	return grantRes.Status, nil
}

func (s *svc) addGrant(ctx context.Context, id *provider.ResourceId, g *provider.Grantee, p *provider.ResourcePermissions, opaque *typesv1beta1.Opaque) (*rpc.Status, error) {
	ref := &provider.Reference{
		ResourceId: id,
	}

	grantReq := &provider.AddGrantRequest{
		Ref: ref,
		Grant: &provider.Grant{
			Grantee:     g,
			Permissions: p,
		},
		Opaque: opaque,
	}

	c, _, err := s.find(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", ref).
			Msg("addGrant: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found"), nil
		}
		return status.NewInternal(ctx, "error finding storage provider"), nil
	}

	grantRes, err := c.AddGrant(ctx, grantReq)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling AddGrant")
	}
	return grantRes.Status, nil
}

func (s *svc) updateGrant(ctx context.Context, id *provider.ResourceId, g *provider.Grantee, p *provider.ResourcePermissions) (*rpc.Status, error) {
	ref := &provider.Reference{
		ResourceId: id,
	}
	grantReq := &provider.UpdateGrantRequest{
		Ref: ref,
		Grant: &provider.Grant{
			Grantee:     g,
			Permissions: p,
		},
	}

	c, _, err := s.find(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", ref).
			Msg("updateGrant: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found"), nil
		}
		return status.NewInternal(ctx, "error finding storage provider"), nil
	}

	grantRes, err := c.UpdateGrant(ctx, grantReq)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling UpdateGrant")
	}
	if grantRes.Status.Code != rpc.Code_CODE_OK {
		return status.NewInternal(ctx,
			"error committing share to storage grant"), nil
	}

	return status.NewOK(ctx), nil
}

func (s *svc) removeGrant(ctx context.Context, id *provider.ResourceId, g *provider.Grantee, p *provider.ResourcePermissions) (*rpc.Status, error) {
	ref := &provider.Reference{
		ResourceId: id,
	}

	grantReq := &provider.RemoveGrantRequest{
		Ref: ref,
		Grant: &provider.Grant{
			Grantee:     g,
			Permissions: p,
		},
	}

	c, _, err := s.find(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", ref).
			Msg("removeGrant: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return status.NewNotFound(ctx, "storage provider not found"), nil
		}
		return status.NewInternal(ctx, "error finding storage provider"), nil
	}

	grantRes, err := c.RemoveGrant(ctx, grantReq)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling RemoveGrant")
	}
	if grantRes.Status.Code != rpc.Code_CODE_OK {
		return status.NewInternal(ctx,
			"error removing storage grant"), nil
	}

	return status.NewOK(ctx), nil
}

func (s *svc) listGrants(ctx context.Context, id *provider.ResourceId) (*provider.ListGrantsResponse, error) {
	ref := &provider.Reference{
		ResourceId: id,
	}

	grantReq := &provider.ListGrantsRequest{
		Ref: ref,
	}

	c, _, err := s.find(ctx, ref)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Interface("reference", ref).
			Msg("listGrants: failed to get storage provider")
		if _, ok := err.(errtypes.IsNotFound); ok {
			return &provider.ListGrantsResponse{
				Status: status.NewNotFound(ctx, "storage provider not found"),
			}, nil
		}
		return &provider.ListGrantsResponse{
			Status: status.NewInternal(ctx, "error finding storage provider"),
		}, nil
	}

	grantRes, err := c.ListGrants(ctx, grantReq)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling ListGrants")
	}
	if grantRes.Status.Code != rpc.Code_CODE_OK {
		return &provider.ListGrantsResponse{Status: status.NewInternal(ctx,
				"error listing storage grants"),
			},
			nil
	}
	return grantRes, nil
}

func (s *svc) addShare(ctx context.Context, req *collaboration.CreateShareRequest) (*collaboration.CreateShareResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("CreateShare: failed to get user share provider")
		return &collaboration.CreateShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}
	// TODO the user share manager needs to be able to decide if the current user is allowed to create that share (and not eg. incerase permissions)
	// jfd: AFAICT this can only be determined by a storage driver - either the storage provider is queried first or the share manager needs to access the storage using a storage driver
	res, err := c.CreateShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling CreateShare")
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return res, nil
	}
	// if we don't need to commit we return earlier
	if !s.c.CommitShareToStorageGrant && !s.c.CommitShareToStorageRef {
		return res, nil
	}

	// TODO(labkode): if both commits are enabled they could be done concurrently.
	if s.c.CommitShareToStorageGrant {
		// If the share is a denial we call  denyGrant instead.
		var status *rpc.Status
		if grants.PermissionsEqual(req.Grant.Permissions.Permissions, &provider.ResourcePermissions{}) {
			status, err = s.denyGrant(ctx, req.ResourceInfo.Id, req.Grant.Grantee, nil)
			if err != nil {
				return nil, errors.Wrap(err, "gateway: error denying grant in storage")
			}
		} else {
			status, err = s.addGrant(ctx, req.ResourceInfo.Id, req.Grant.Grantee, req.Grant.Permissions.Permissions, nil)
			if err != nil {
				return nil, errors.Wrap(err, "gateway: error adding grant to storage")
			}
		}

		switch status.Code {
		case rpc.Code_CODE_OK:
			s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), req.ResourceInfo.Id)
		case rpc.Code_CODE_UNIMPLEMENTED:
			appctx.GetLogger(ctx).Debug().Interface("status", status).Interface("req", req).Msg("storing grants not supported, ignoring")
		default:
			return &collaboration.CreateShareResponse{
				Status: status,
			}, err
		}
	}
	return res, nil
}

func (s *svc) addSpaceShare(ctx context.Context, req *collaboration.CreateShareRequest) (*collaboration.CreateShareResponse, error) {
	// If the share is a denial we call  denyGrant instead.
	var st *rpc.Status
	var err error
	// TODO: change CS3 APIs
	opaque := typesv1beta1.Opaque{
		Map: map[string]*typesv1beta1.OpaqueEntry{
			"spacegrant": {},
		},
	}
	if grants.PermissionsEqual(req.Grant.Permissions.Permissions, &provider.ResourcePermissions{}) {
		st, err = s.denyGrant(ctx, req.ResourceInfo.Id, req.Grant.Grantee, &opaque)
		if err != nil {
			return nil, errors.Wrap(err, "gateway: error denying grant in storage")
		}
	} else {
		st, err = s.addGrant(ctx, req.ResourceInfo.Id, req.Grant.Grantee, req.Grant.Permissions.Permissions, &opaque)
		if err != nil {
			return nil, errors.Wrap(err, "gateway: error adding grant to storage")
		}
	}

	switch st.Code {
	case rpc.Code_CODE_OK:
		s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), req.ResourceInfo.Id)
		s.cache.RemoveListStorageProviders(req.ResourceInfo.Id)
	case rpc.Code_CODE_UNIMPLEMENTED:
		appctx.GetLogger(ctx).Debug().Interface("status", st).Interface("req", req).Msg("storing grants not supported, ignoring")
	default:
		return &collaboration.CreateShareResponse{
			Status: st,
		}, err
	}

	return &collaboration.CreateShareResponse{
		Status: status.NewOK(ctx),
		Share: &collaboration.Share{
			ResourceId:  req.ResourceInfo.Id,
			Permissions: &collaboration.SharePermissions{Permissions: req.Grant.Permissions.GetPermissions()},
			Grantee:     req.Grant.Grantee,
		},
	}, nil
}

func (s *svc) removeShare(ctx context.Context, req *collaboration.RemoveShareRequest) (*collaboration.RemoveShareResponse, error) {
	c, err := pool.GetUserShareProviderClient(s.c.UserShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).
			Err(err).
			Msg("RemoveShare: failed to get user share provider")
		return &collaboration.RemoveShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	// if we need to commit the share, we need the resource it points to.
	var share *collaboration.Share
	// FIXME: I will cause a panic as share will be nil when I'm false
	if s.c.CommitShareToStorageGrant || s.c.CommitShareToStorageRef {
		getShareReq := &collaboration.GetShareRequest{
			Ref: req.Ref,
		}
		getShareRes, err := c.GetShare(ctx, getShareReq)
		if err != nil {
			return nil, errors.Wrap(err, "gateway: error calling GetShare")
		}

		if getShareRes.Status.Code != rpc.Code_CODE_OK {
			res := &collaboration.RemoveShareResponse{
				Status: status.NewInternal(ctx,
					"error getting share when committing to the storage"),
			}
			return res, nil
		}
		share = getShareRes.Share
	}

	res, err := c.RemoveShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling RemoveShare")
	}

	// we do not want to remove the reference if it is a reshare
	if utils.UserEqual(share.Owner, share.Creator) {
		s.removeReference(ctx, share.ResourceId)
	}

	// if we don't need to commit we return earlier
	if !s.c.CommitShareToStorageGrant && !s.c.CommitShareToStorageRef {
		return res, nil
	}

	// TODO(labkode): if both commits are enabled they could be done concurrently.
	if s.c.CommitShareToStorageGrant {
		removeGrantStatus, err := s.removeGrant(ctx, share.ResourceId, share.Grantee, share.Permissions.Permissions)
		if err != nil {
			return nil, errors.Wrap(err, "gateway: error removing grant from storage")
		}
		if removeGrantStatus.Code != rpc.Code_CODE_OK {
			return &collaboration.RemoveShareResponse{
				Status: removeGrantStatus,
			}, err
		}
	}

	s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), share.ResourceId)
	return res, nil
}

func (s *svc) removeSpaceShare(ctx context.Context, ref *provider.ResourceId, grantee *provider.Grantee) (*collaboration.RemoveShareResponse, error) {
	listGrantRes, err := s.listGrants(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error getting grant to remove from storage")
	}
	var permissions *provider.ResourcePermissions
	for _, g := range listGrantRes.Grants {
		if isEqualGrantee(g.Grantee, grantee) {
			permissions = g.Permissions
		}
	}
	if permissions == nil {
		return nil, errors.New("gateway: error getting grant to remove from storage")
	}
	removeGrantStatus, err := s.removeGrant(ctx, ref, grantee, permissions)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error removing grant from storage")
	}
	if removeGrantStatus.Code != rpc.Code_CODE_OK {
		return &collaboration.RemoveShareResponse{
			Status: removeGrantStatus,
		}, err
	}
	s.cache.RemoveStat(ctxpkg.ContextMustGetUser(ctx), ref)
	s.cache.RemoveListStorageProviders(ref)
	return &collaboration.RemoveShareResponse{Status: status.NewOK(ctx)}, nil
}

func refIsSpaceRoot(ref *provider.ResourceId) bool {
	if ref == nil {
		return false
	}
	if ref.StorageId == "" || ref.OpaqueId == "" {
		return false
	}
	_, sid := storagespace.SplitStorageID(ref.GetStorageId())
	return sid == ref.OpaqueId
}

func shareIsSpaceRoot(key *collaboration.ShareKey) bool {
	if key == nil {
		return false
	}
	return refIsSpaceRoot(key.ResourceId)
}

func isEqualGrantee(a, b *provider.Grantee) bool {
	// Ideally we would want to use utils.GranteeEqual()
	// but the grants stored in the decomposedfs aren't complete (missing usertype and idp)
	// because of that the check would fail so we can only check the ... for now.
	if a.Type != b.Type {
		return false
	}

	var aID, bID string
	switch a.Type {
	case provider.GranteeType_GRANTEE_TYPE_GROUP:
		aID = a.GetGroupId().OpaqueId
		bID = b.GetGroupId().OpaqueId
	case provider.GranteeType_GRANTEE_TYPE_USER:
		aID = a.GetUserId().OpaqueId
		bID = b.GetUserId().OpaqueId
	}
	return aID == bID
}

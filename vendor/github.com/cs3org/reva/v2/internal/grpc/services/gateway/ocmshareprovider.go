// Copyright 2018-2023 CERN
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
	"fmt"
	"net/url"
	"path"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	datatx "github.com/cs3org/go-cs3apis/cs3/tx/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/pkg/errors"
)

// TODO(labkode): add multi-phase commit logic when commit share or commit ref is enabled.
func (s *svc) CreateOCMShare(ctx context.Context, req *ocm.CreateOCMShareRequest) (*ocm.CreateOCMShareResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		return &ocm.CreateOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.CreateOCMShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling CreateShare")
	}

	status, err := s.addGrant(ctx, req.ResourceId, req.Grantee, req.AccessMethods[0].GetWebdavOptions().Permissions, req.Expiration, nil)
	if err != nil {
		appctx.GetLogger(ctx).Debug().Interface("status", status).Interface("req", req).Msg(err.Error())
		return nil, errors.Wrap(err, "gateway: error adding grant to storage")
	}

	switch status.Code {
	case rpc.Code_CODE_OK:
		s.statCache.RemoveStatContext(ctx, ctxpkg.ContextMustGetUser(ctx).GetId(), req.ResourceId)
	case rpc.Code_CODE_UNIMPLEMENTED:
		appctx.GetLogger(ctx).Debug().Interface("status", status).Interface("req", req).Msg("storing grants not supported, ignoring")
	default:
		appctx.GetLogger(ctx).Debug().Interface("status", status).Interface("req", req).Msg("storing grants is not successful")
		return &ocm.CreateOCMShareResponse{
			Status: status,
		}, err
	}

	return res, nil
}

func (s *svc) RemoveOCMShare(ctx context.Context, req *ocm.RemoveOCMShareRequest) (*ocm.RemoveOCMShareResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		return &ocm.RemoveOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.RemoveOCMShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling RemoveShare")
	}

	return res, nil
}

// TODO(labkode): we need to validate share state vs storage grant and storage ref
// If there are any inconsistencies, the share needs to be flag as invalid and a background process
// or active fix needs to be performed.
func (s *svc) GetOCMShare(ctx context.Context, req *ocm.GetOCMShareRequest) (*ocm.GetOCMShareResponse, error) {
	return s.getOCMShare(ctx, req)
}

func (s *svc) getOCMShare(ctx context.Context, req *ocm.GetOCMShareRequest) (*ocm.GetOCMShareResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.GetOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.GetOCMShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetShare")
	}

	return res, nil
}

func (s *svc) GetOCMShareByToken(ctx context.Context, req *ocm.GetOCMShareByTokenRequest) (*ocm.GetOCMShareByTokenResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetOCMShareProviderClient")
	}

	res, err := c.GetOCMShareByToken(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetOCMShareByToken")
	}

	return res, nil
}

// TODO(labkode): read GetShare comment.
func (s *svc) ListOCMShares(ctx context.Context, req *ocm.ListOCMSharesRequest) (*ocm.ListOCMSharesResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.ListOCMSharesResponse{
			Status: status.NewInternal(ctx, "error getting user share provider client"),
		}, nil
	}

	res, err := c.ListOCMShares(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling ListShares")
	}

	return res, nil
}

func (s *svc) UpdateOCMShare(ctx context.Context, req *ocm.UpdateOCMShareRequest) (*ocm.UpdateOCMShareResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.UpdateOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	res, err := c.UpdateOCMShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling UpdateShare")
	}

	return res, nil
}

func (s *svc) ListReceivedOCMShares(ctx context.Context, req *ocm.ListReceivedOCMSharesRequest) (*ocm.ListReceivedOCMSharesResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.ListReceivedOCMSharesResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	res, err := c.ListReceivedOCMShares(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling ListReceivedShares")
	}

	return res, nil
}

func (s *svc) UpdateReceivedOCMShare(ctx context.Context, req *ocm.UpdateReceivedOCMShareRequest) (*ocm.UpdateReceivedOCMShareResponse, error) {
	log := appctx.GetLogger(ctx)
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	// retrieve the current received share
	getShareReq := &ocm.GetReceivedOCMShareRequest{
		Ref: &ocm.ShareReference{
			Spec: &ocm.ShareReference_Id{
				Id: req.Share.Id,
			},
		},
	}
	getShareRes, err := s.GetReceivedOCMShare(ctx, getShareReq)
	if err != nil {
		log.Error().Err(err).Msg("gateway: error calling GetReceivedOCMShare")
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: &rpc.Status{
				Code: rpc.Code_CODE_INTERNAL,
			},
		}, nil
	}
	if getShareRes.Status.Code != rpc.Code_CODE_OK {
		log.Error().Msg("gateway: error calling GetReceivedOCMShare")
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: &rpc.Status{
				Code:    rpc.Code_CODE_INTERNAL,
				Message: "gateway: error calling GetReceivedOCMShare",
			},
		}, nil
	}
	share := getShareRes.Share
	if share == nil {
		log.Error().Err(err).Msg("gateway: got a nil share from GetReceivedOCMShare")
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: &rpc.Status{
				Code:    rpc.Code_CODE_INTERNAL,
				Message: "gateway: got a nil share from GetReceivedOCMShare",
			},
		}, nil
	}

	res, err := c.UpdateReceivedOCMShare(ctx, req)
	if err != nil {
		log.Err(err).Msg("gateway: error calling UpdateReceivedShare")
		return &ocm.UpdateReceivedOCMShareResponse{
			Status: &rpc.Status{
				Code: rpc.Code_CODE_INTERNAL,
			},
		}, nil
	}

	for i := range req.UpdateMask.Paths {
		switch req.UpdateMask.Paths[i] {
		case "state":
			switch req.GetShare().GetState() {
			case ocm.ShareState_SHARE_STATE_ACCEPTED:
				// for a transfer this is handled elsewhere
			case ocm.ShareState_SHARE_STATE_PENDING:
				// currently no consequences
			case ocm.ShareState_SHARE_STATE_REJECTED:
				// TODO
				return res, nil
			}
		case "mount_point":
			// TODO(labkode): implementing updating mount point
			err = errtypes.NotSupported("gateway: update of mount point is not yet implemented")
			return &ocm.UpdateReceivedOCMShareResponse{
				Status: status.NewUnimplemented(ctx, err, "error updating received share"),
			}, nil
		default:
			return nil, errtypes.NotSupported("updating " + req.UpdateMask.Paths[i] + " is not supported")
		}
	}
	// handle transfer in case it has not already been accepted
	if s.isTransferShare(share) && req.GetShare().State == ocm.ShareState_SHARE_STATE_ACCEPTED {
		if share.State == ocm.ShareState_SHARE_STATE_ACCEPTED {
			log.Err(err).Msg("gateway: error calling UpdateReceivedShare, share already accepted.")
			return &ocm.UpdateReceivedOCMShareResponse{
				Status: &rpc.Status{
					Code:    rpc.Code_CODE_FAILED_PRECONDITION,
					Message: "Share already accepted.",
				},
			}, err
		}
		// get provided destination path
		transferDestinationPath, err := s.getTransferDestinationPath(ctx, req)
		if err != nil {
			if err != nil {
				log.Err(err).Msg("gateway: error calling UpdateReceivedShare")
				return &ocm.UpdateReceivedOCMShareResponse{
					Status: &rpc.Status{
						Code: rpc.Code_CODE_INTERNAL,
					},
				}, err
			}
		}

		error := s.handleTransfer(ctx, share, transferDestinationPath)
		if error != nil {
			log.Err(error).Msg("gateway: error handling transfer in UpdateReceivedShare")
			return &ocm.UpdateReceivedOCMShareResponse{
				Status: &rpc.Status{
					Code: rpc.Code_CODE_INTERNAL,
				},
			}, error
		}
	}
	return res, nil
}

func (s *svc) handleTransfer(ctx context.Context, share *ocm.ReceivedShare, transferDestinationPath string) error {
	log := appctx.GetLogger(ctx)

	protocol, ok := s.getTransferProtocol(share)
	if !ok {
		return errors.New("gateway: unable to retrieve transfer protocol")
	}
	sourceURI := protocol.SourceUri

	// get the webdav endpoint of the grantee's idp
	var granteeIdp string
	if share.GetGrantee().Type == provider.GranteeType_GRANTEE_TYPE_USER {
		granteeIdp = share.GetGrantee().GetUserId().Idp
	}
	if share.GetGrantee().Type == provider.GranteeType_GRANTEE_TYPE_GROUP {
		granteeIdp = share.GetGrantee().GetGroupId().Idp
	}
	destWebdavEndpoint, err := s.getWebdavEndpoint(ctx, granteeIdp)
	if err != nil {
		log.Err(err).Msg("gateway: error calling UpdateReceivedShare")
		return err
	}
	destWebdavEndpointURL, err := url.Parse(destWebdavEndpoint)
	if err != nil {
		log.Err(err).Msg("gateway: error calling UpdateReceivedShare: unable to parse webdav endpoint \"" + destWebdavEndpoint + "\" into URL structure")
		return err
	}
	destWebdavHost, err := s.getWebdavHost(ctx, granteeIdp)
	if err != nil {
		log.Err(err).Msg("gateway: error calling UpdateReceivedShare")
		return err
	}
	var dstWebdavURLString string
	if strings.Contains(destWebdavHost, "://") {
		dstWebdavURLString = destWebdavHost
	} else {
		dstWebdavURLString = "http://" + destWebdavHost
	}
	dstWebdavHostURL, err := url.Parse(dstWebdavURLString)
	if err != nil {
		log.Err(err).Msg("gateway: error calling UpdateReceivedShare: unable to parse webdav service host \"" + dstWebdavURLString + "\" into URL structure")
		return err
	}
	destServiceHost := dstWebdavHostURL.Host + dstWebdavHostURL.Path
	// optional prefix must only appear in target url path:
	// http://...token...@reva.eu/prefix/?name=remote.php/webdav/home/...
	destEndpointPath := strings.TrimPrefix(destWebdavEndpointURL.Path, dstWebdavHostURL.Path)
	destEndpointScheme := destWebdavEndpointURL.Scheme
	destToken := ctxpkg.ContextMustGetToken(ctx)
	destPath := path.Join(destEndpointPath, transferDestinationPath, path.Base(share.Name))
	destTargetURI := fmt.Sprintf("%s://%s@%s?name=%s", destEndpointScheme, destToken, destServiceHost, destPath)
	// var destUri string
	req := &datatx.CreateTransferRequest{
		SrcTargetUri:  sourceURI,
		DestTargetUri: destTargetURI,
		ShareId:       share.Id,
	}

	res, err := s.CreateTransfer(ctx, req)
	if err != nil {
		return err
	}
	log.Info().Msgf("gateway: CreateTransfer: %v", res.TxInfo)
	return nil
}

func (s *svc) isTransferShare(share *ocm.ReceivedShare) bool {
	_, ok := s.getTransferProtocol(share)
	return ok
}

func (s *svc) getTransferDestinationPath(ctx context.Context, req *ocm.UpdateReceivedOCMShareRequest) (string, error) {
	log := appctx.GetLogger(ctx)
	// the destination path is not part of any protocol, but an opaque field
	destPathOpaque, ok := req.GetOpaque().GetMap()["transfer_destination_path"]
	if ok {
		switch destPathOpaque.Decoder {
		case "plain":
			if string(destPathOpaque.Value) != "" {
				return string(destPathOpaque.Value), nil
			}
		default:
			return "", errtypes.NotSupported("decoder of opaque entry 'transfer_destination_path' not recognized: " + destPathOpaque.Decoder)
		}
	}
	log.Info().Msg("destination path not provided, trying default transfer destination folder")
	if s.c.DataTransfersFolder == "" {
		return "", errtypes.NotSupported("no destination path provided and default transfer destination folder is not set")
	}
	return s.c.DataTransfersFolder, nil
}

func (s *svc) GetReceivedOCMShare(ctx context.Context, req *ocm.GetReceivedOCMShareRequest) (*ocm.GetReceivedOCMShareResponse, error) {
	c, err := pool.GetOCMShareProviderClient(s.c.OCMShareProviderEndpoint)
	if err != nil {
		appctx.GetLogger(ctx).Error().Err(err).Msg("error calling GetOCMShareProviderClient")
		return &ocm.GetReceivedOCMShareResponse{
			Status: status.NewInternal(ctx, "error getting share provider client"),
		}, nil
	}

	res, err := c.GetReceivedOCMShare(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "gateway: error calling GetReceivedShare")
	}

	return res, nil
}

func (s *svc) getTransferProtocol(share *ocm.ReceivedShare) (*ocm.TransferProtocol, bool) {
	for _, p := range share.Protocols {
		if d, ok := p.Term.(*ocm.Protocol_TransferOptions); ok {
			return d.TransferOptions, true
		}
	}
	return nil, false
}

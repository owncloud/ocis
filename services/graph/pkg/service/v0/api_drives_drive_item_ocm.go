package svc

import (
	"context"
	"errors"

	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

var (
	// ErrUnmountOCMShare is returned when unmounting a share fails
	ErrUnmountOCMShare = errorcode.New(errorcode.InvalidRequest, "unmounting ocm share failed")

	// ErrMountOCMShare is returned when mounting a share fails
	ErrMountOCMShare = errorcode.New(errorcode.InvalidRequest, "mounting ocm share failed")
)

type (
	// UpdateOCMShareClosure is a closure that injects required updates into the update request
	UpdateOCMShareClosure func(share *ocm.ReceivedShare, request *ocm.UpdateReceivedOCMShareRequest)
)

// GetOCMSharesForResource returns all federated shares for a given resourceID
func (s DrivesDriveItemService) GetOCMSharesForResource(ctx context.Context, resourceID *storageprovider.ResourceId) ([]*ocm.ReceivedShare, error) {
	// Find all accepted shares for this resource
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	receivedOCMSharesResponse, err := gatewayClient.ListReceivedOCMShares(ctx, &ocm.ListReceivedOCMSharesRequest{
		/* ocm has no filters, yet
		Filters: append([]*collaboration.Filter{
			{
				Type: collaboration.Filter_TYPE_RESOURCE_ID,
				Term: &collaboration.Filter_ResourceId{
					ResourceId: resourceID,
				},
			},
		}, filters...),
		*/
	})
	switch {
	case err != nil:
		return nil, err
	case len(receivedOCMSharesResponse.GetShares()) == 0:
		return nil, ErrNoShares
	default:
		return receivedOCMSharesResponse.GetShares(), errorcode.FromCS3Status(receivedOCMSharesResponse.GetStatus(), err)
	}
}

// UpdateShares updates multiple shares;
// it could happen that some shares are updated and some are not,
// this will return a list of updated shares and a list of errors;
// there is no guarantee that all updates are successful
func (s DrivesDriveItemService) UpdateOCMShares(ctx context.Context, shares []*ocm.ReceivedShare, updater UpdateOCMShareClosure) ([]*ocm.ReceivedShare, error) {
	errs := make([]error, 0, len(shares))
	updatedShares := make([]*ocm.ReceivedShare, 0, len(shares))

	for _, share := range shares {
		err := s.UpdateOCMShare(
			ctx,
			share,
			updater,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		updatedShares = append(updatedShares, share)
	}

	return updatedShares, errors.Join(errs...)
}

// UpdateOCMShare updates a single share
func (s DrivesDriveItemService) UpdateOCMShare(ctx context.Context, share *ocm.ReceivedShare, updater UpdateOCMShareClosure) error {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return err
	}

	updateReceivedOCMShareRequest := &ocm.UpdateReceivedOCMShareRequest{
		Share: &ocm.ReceivedShare{
			Id: share.GetId(),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{}},
	}

	switch updater {
	case nil:
		return ErrNoUpdater
	default:
		updater(share, updateReceivedOCMShareRequest)
	}

	if len(updateReceivedOCMShareRequest.GetUpdateMask().GetPaths()) == 0 {
		return ErrNoUpdates
	}

	updateReceivedOCMShareResponse, err := gatewayClient.UpdateReceivedOCMShare(ctx, updateReceivedOCMShareRequest)
	return errorcode.FromCS3Status(updateReceivedOCMShareResponse.GetStatus(), err)
}

func (s DrivesDriveItemService) MountOCMShare(ctx context.Context, resourceID *storageprovider.ResourceId /*, name string*/) ([]*ocm.ReceivedShare, error) {
	/*
		if filepath.IsAbs(name) {
			return nil, ErrAbsoluteNamePath
		}

		if name != "" {
			name = filepath.Clean(name)
		}
	*/

	shares, err := s.GetOCMSharesForResource(ctx, resourceID)
	if err != nil {
		return nil, err
	}

	availableShares := make([]*ocm.ReceivedShare, 0, len(shares))
	mountedShares := make([]*ocm.ReceivedShare, 0, 1)
	for _, v := range shares {
		switch v.GetState() {
		case ocm.ShareState_SHARE_STATE_ACCEPTED:
			mountedShares = append(mountedShares, v)
		case ocm.ShareState_SHARE_STATE_PENDING, ocm.ShareState_SHARE_STATE_REJECTED:
			availableShares = append(availableShares, v)
		}
	}
	if len(availableShares) == 0 {
		if len(mountedShares) > 0 {
			return nil, ErrAlreadyMounted
		}
		return nil, ErrNoShares
	}

	updatedShares, err := s.UpdateOCMShares(ctx, availableShares, func(share *ocm.ReceivedShare, request *ocm.UpdateReceivedOCMShareRequest) {
		request.Share.State = ocm.ShareState_SHARE_STATE_ACCEPTED
		request.UpdateMask.Paths = append(request.UpdateMask.Paths, _fieldMaskPathState)

		// only update if mountPoint name is not empty and the path has changed
		/* ocm shares have no mount point???
		if name != "" {
			mountPoint := share.GetMountPoint()
			if mountPoint == nil {
				mountPoint = &storageprovider.Reference{}
			}

			if filepath.Clean(mountPoint.GetPath()) != name {
				mountPoint.Path = name
				request.Share.MountPoint = mountPoint
				request.UpdateMask.Paths = append(request.UpdateMask.Paths, _fieldMaskPathMountPoint)
			}
		}
		*/
	})

	errs, ok := err.(interface{ Unwrap() []error })
	if ok && len(errs.Unwrap()) == len(availableShares) {
		// none of the received ocm shares could be accepted.
		// this is an error, return it.
		return nil, err
	}

	return updatedShares, nil
}

package svc

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
)

const (
	_fieldMaskPathState      = "state"
	_fieldMaskPathMountPoint = "mount_point"
)

// DrivesDriveItemProvider is the interface that needs to be implemented by the individual space service
type DrivesDriveItemProvider interface {
	MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error)
	UnmountShare(ctx context.Context, resourceID storageprovider.ResourceId) error
}

// DrivesDriveItemService contains the production business logic for everything that relates to drives
type DrivesDriveItemService struct {
	logger          log.Logger
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
	identityCache   identity.IdentityCache
}

// NewDrivesDriveItemService creates a new DrivesDriveItemService
func NewDrivesDriveItemService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], identityCache identity.IdentityCache) (DrivesDriveItemService, error) {
	return DrivesDriveItemService{
		logger:          log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
		gatewaySelector: gatewaySelector,
		identityCache:   identityCache,
	}, nil
}

// UnmountShare unmounts a share from the share-jail
func (s DrivesDriveItemService) UnmountShare(ctx context.Context, resourceID storageprovider.ResourceId) error {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return err
	}

	// This is a bit of a hack. We should not rely on a specific format of the item id.
	// But currently there is no other way to get the ShareID.
	shareId := resourceID.GetOpaqueId()

	// Now, find out the resourceID of the shared resource
	getReceivedShareResponse, err := gatewayClient.GetReceivedShare(ctx,
		&collaboration.GetReceivedShareRequest{
			Ref: &collaboration.ShareReference{
				Spec: &collaboration.ShareReference_Id{
					Id: &collaboration.ShareId{
						OpaqueId: shareId,
					},
				},
			},
		},
	)
	if errCode := errorcode.FromCS3Status(getReceivedShareResponse.GetStatus(), err); errCode != nil {
		s.logger.Debug().Err(errCode).
			Str("shareid", shareId).
			Msg("failed to read share")
		return errCode
	}

	// Find all accepted shares for this resource
	gatewayClient, err = s.gatewaySelector.Next()
	if err != nil {
		return err
	}
	receivedSharesResponse, err := gatewayClient.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{
		Filters: []*collaboration.Filter{
			{
				Type: collaboration.Filter_TYPE_STATE,
				Term: &collaboration.Filter_State{
					State: collaboration.ShareState_SHARE_STATE_ACCEPTED,
				},
			},
			{
				Type: collaboration.Filter_TYPE_RESOURCE_ID,
				Term: &collaboration.Filter_ResourceId{
					ResourceId: getReceivedShareResponse.GetShare().GetShare().GetResourceId(),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	var errs []error

	// Reject all the shares for this resource
	for _, receivedShare := range receivedSharesResponse.GetShares() {
		receivedShare.State = collaboration.ShareState_SHARE_STATE_REJECTED

		updateReceivedShareRequest := &collaboration.UpdateReceivedShareRequest{
			Share:      receivedShare,
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{_fieldMaskPathState}},
		}

		_, err := gatewayClient.UpdateReceivedShare(ctx, updateReceivedShareRequest)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}

// MountShare mounts a share
func (s DrivesDriveItemService) MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error) {
	if filepath.IsAbs(name) {
		return libregraph.DriveItem{}, errorcode.New(errorcode.InvalidRequest, "name cannot be an absolute path")
	}
	name = filepath.Clean(name)

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return libregraph.DriveItem{}, err
	}

	// Get all shares that the user has received for this resource. There might be multiple
	receivedSharesResponse, err := gatewayClient.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{
		Filters: []*collaboration.Filter{
			{
				Type: collaboration.Filter_TYPE_STATE,
				Term: &collaboration.Filter_State{
					State: collaboration.ShareState_SHARE_STATE_PENDING,
				},
			},
			{
				Type: collaboration.Filter_TYPE_STATE,
				Term: &collaboration.Filter_State{
					State: collaboration.ShareState_SHARE_STATE_REJECTED,
				},
			},
			{
				Type: collaboration.Filter_TYPE_RESOURCE_ID,
				Term: &collaboration.Filter_ResourceId{
					ResourceId: &resourceID,
				},
			},
		},
	})
	if err != nil {
		return libregraph.DriveItem{}, err
	}

	var errs []error

	var acceptedShares []*collaboration.ReceivedShare

	// try to accept all the received shares for this resource. So that the stat is in sync across all
	// shares
	for _, receivedShare := range receivedSharesResponse.GetShares() {
		updateMask := &fieldmaskpb.FieldMask{Paths: []string{_fieldMaskPathState}}
		receivedShare.State = collaboration.ShareState_SHARE_STATE_ACCEPTED

		// only update if mountPoint name is not empty and the path has changed
		if name != "" {
			mountPoint := receivedShare.GetMountPoint()
			if mountPoint == nil {
				mountPoint = &storageprovider.Reference{}
			}

			if filepath.Clean(mountPoint.GetPath()) != name {
				mountPoint.Path = name
				receivedShare.MountPoint = mountPoint
				updateMask.Paths = append(updateMask.Paths, _fieldMaskPathMountPoint)
			}
		}

		updateReceivedShareRequest := &collaboration.UpdateReceivedShareRequest{
			Share:      receivedShare,
			UpdateMask: updateMask,
		}

		gatewayClient, err = s.gatewaySelector.Next()
		if err != nil {
			return libregraph.DriveItem{}, err
		}
		updateReceivedShareResponse, err := gatewayClient.UpdateReceivedShare(ctx, updateReceivedShareRequest)
		switch errCode := errorcode.FromCS3Status(updateReceivedShareResponse.GetStatus(), err); {
		case errCode == nil:
			acceptedShares = append(acceptedShares, updateReceivedShareResponse.GetShare())
		default:
			// Just log at debug level here. If a single accept for any of the received shares failed this
			// is not a critical problem. We mainly need to handle the case where all accepts fail. (Outside
			// the loop)
			s.logger.Debug().Err(errCode).
				Str("shareid", receivedShare.GetShare().GetId().String()).
				Str("resourceid", receivedShare.GetShare().GetResourceId().String()).
				Msg("failed to accept share")
			errs = append(errs, errCode)
		}
	}

	if len(receivedSharesResponse.GetShares()) == len(errs) {
		// none of the received shares could be accepted. This is an error. Return it.
		return libregraph.DriveItem{}, errors.Join(errs...)
	}

	// As the accepted shares are all for the same resource they should collapse to a single driveitem
	items, err := cs3ReceivedSharesToDriveItems(ctx, &s.logger, gatewayClient, s.identityCache, acceptedShares)
	switch {
	case err != nil:
		return libregraph.DriveItem{}, err
	case len(items) != 1:
		return libregraph.DriveItem{}, errorcode.New(errorcode.GeneralException, "failed to convert accepted shares into drive-item")
	}
	return items[0], nil
}

// DrivesDriveItemApi is the api that registers the http endpoints which expose needed operation to the graph api.
// the business logic is delegated to the space service and further down to the cs3 client.
type DrivesDriveItemApi struct {
	logger                 log.Logger
	drivesDriveItemService DrivesDriveItemProvider
}

// NewDrivesDriveItemApi creates a new DrivesDriveItemApi
func NewDrivesDriveItemApi(drivesDriveItemService DrivesDriveItemProvider, logger log.Logger) (DrivesDriveItemApi, error) {
	return DrivesDriveItemApi{
		logger:                 log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		drivesDriveItemService: drivesDriveItemService,
	}, nil
}

// DeleteDriveItem deletes a drive item
func (api DrivesDriveItemApi) DeleteDriveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driveID, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		msg := "invalid driveID or itemID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	if !IsShareJail(driveID) {
		msg := "invalid driveID, must be share jail"
		api.logger.Debug().Interface("driveID", driveID).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	if err := api.drivesDriveItemService.UnmountShare(ctx, itemID); err != nil {
		msg := "unmounting share failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusOK)
}

// CreateDriveItem creates a drive item
func (api DrivesDriveItemApi) CreateDriveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		api.logger.Debug().Err(err).Msg("invalid driveID")
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, "invalid driveID")
		return
	}

	if !IsShareJail(driveID) {
		msg := "invalid driveID, must be share jail"
		api.logger.Debug().Interface("driveID", driveID).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	requestDriveItem := libregraph.DriveItem{}
	if err := StrictJSONUnmarshal(r.Body, &requestDriveItem); err != nil {
		msg := "invalid request body"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	remoteItem := requestDriveItem.GetRemoteItem()
	resourceId, err := storagespace.ParseID(remoteItem.GetId())
	if err != nil {
		msg := "invalid remote item id"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	mountShareResponse, err := api.drivesDriveItemService.
		MountShare(ctx, resourceId, requestDriveItem.GetName())
	if err != nil {
		msg := "mounting share failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, mountShareResponse)
}

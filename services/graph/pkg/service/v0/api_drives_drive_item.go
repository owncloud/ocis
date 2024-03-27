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
	_fieldMaskPathHidden     = "hidden"
)

var (
	ErrNoChanges     = errors.New("no changes specified")
	ErrNotAShareJail = errors.New("id does not belong to a share jail")
)

// DrivesDriveItemProvider is the interface that needs to be implemented by the individual space service
type DrivesDriveItemProvider interface {
	MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) ([]*collaboration.ReceivedShare, error)
	UnmountShare(ctx context.Context, shareID *collaboration.ShareId) error
	UpdateShare(ctx context.Context, shareID *collaboration.ShareId, instructions *ShareUpdateInstruction) (*collaboration.ReceivedShare, error)
}

// ShareUpdateInstruction is a helper struct to build the update instruction for a share
type ShareUpdateInstruction struct {
	changes    []string
	hidden     bool
	mountPoint *storageprovider.Reference
	state      collaboration.ShareState
}

// State sets the state of the share
func (i *ShareUpdateInstruction) State(state collaboration.ShareState) *ShareUpdateInstruction {
	i.changes = append(i.changes, _fieldMaskPathState)
	i.state = state

	return i
}

// MountPoint sets the mount point of the share
func (i *ShareUpdateInstruction) MountPoint(mountPoint *storageprovider.Reference) *ShareUpdateInstruction {
	i.changes = append(i.changes, _fieldMaskPathMountPoint)
	i.mountPoint = mountPoint

	return i
}

// Hidden sets the visibility of the share
func (i *ShareUpdateInstruction) Hidden(hidden bool) *ShareUpdateInstruction {
	i.changes = append(i.changes, _fieldMaskPathHidden)
	i.hidden = hidden

	return i
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

// UpdateShare updates the visibility of a share
func (s DrivesDriveItemService) UpdateShare(ctx context.Context, shareID *collaboration.ShareId, instructions *ShareUpdateInstruction) (*collaboration.ReceivedShare, error) {
	if len(instructions.changes) == 0 {
		return nil, ErrNoChanges
	}

	updateReceivedShareRequest := &collaboration.UpdateReceivedShareRequest{
		Share: &collaboration.ReceivedShare{
			Share: &collaboration.Share{
				Id: shareID,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{}},
	}

	for _, c := range instructions.changes {
		switch c {
		case _fieldMaskPathState:
			updateReceivedShareRequest.Share.State = instructions.state
		case _fieldMaskPathMountPoint:
			updateReceivedShareRequest.Share.MountPoint = instructions.mountPoint
		case _fieldMaskPathHidden:
			updateReceivedShareRequest.Share.Hidden = instructions.hidden
		default:
			continue
		}

		updateReceivedShareRequest.UpdateMask.Paths = append(updateReceivedShareRequest.UpdateMask.Paths, c)
	}

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	updateReceivedShareResponse, err := gatewayClient.UpdateReceivedShare(ctx, updateReceivedShareRequest)
	return updateReceivedShareResponse.GetShare(), errorcode.FromCS3Status(updateReceivedShareResponse.GetStatus(), err)
}

// UnmountShare unmounts a share from the sharejail
func (s DrivesDriveItemService) UnmountShare(ctx context.Context, shareID *collaboration.ShareId) error {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return err
	}

	// Now, find out the resourceID of the shared resource
	getReceivedShareResponse, err := gatewayClient.GetReceivedShare(ctx,
		&collaboration.GetReceivedShareRequest{
			Ref: &collaboration.ShareReference{
				Spec: &collaboration.ShareReference_Id{
					Id: shareID,
				},
			},
		},
	)
	if err := errorcode.FromCS3Status(getReceivedShareResponse.GetStatus(), err); err != nil {
		s.logger.Debug().Err(err).
			Str("shareID", shareID.GetOpaqueId()).
			Msg("failed to read share")
		return err
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
	if len(receivedSharesResponse.GetShares()) == 0 {
		return errorcode.New(errorcode.InvalidRequest, "invalid itemID")
	}

	var errs []error

	// Reject all the shares for this resource
	for _, receivedShare := range receivedSharesResponse.GetShares() {
		if _, err := s.UpdateShare(
			ctx,
			receivedShare.GetShare().GetId(),
			(&ShareUpdateInstruction{}).
				State(collaboration.ShareState_SHARE_STATE_REJECTED),
		); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// MountShare mounts a share
func (s DrivesDriveItemService) MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) ([]*collaboration.ReceivedShare, error) {
	if filepath.IsAbs(name) {
		return nil, errorcode.New(errorcode.InvalidRequest, "name cannot be an absolute path")
	}

	if name != "" {
		name = filepath.Clean(name)
	}

	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return nil, err
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
	switch {
	case err != nil:
		return nil, err
	case len(receivedSharesResponse.GetShares()) == 0:
		return nil, errorcode.New(errorcode.InvalidRequest, "invalid itemID")
	}

	receivedShares := receivedSharesResponse.GetShares()
	errs := make([]error, 0, len(receivedShares))
	acceptedShares := make([]*collaboration.ReceivedShare, 0, len(receivedShares))

	// try to accept all the received shares for this resource. So that the stat is in sync across all
	// shares
	for _, receivedShare := range receivedShares {
		updateInstruction := &ShareUpdateInstruction{}
		updateInstruction.State(collaboration.ShareState_SHARE_STATE_ACCEPTED)

		// only update if mountPoint name is not empty and the path has changed
		if name != "" {
			mountPoint := receivedShare.GetMountPoint()
			if mountPoint == nil {
				mountPoint = &storageprovider.Reference{}
			}

			if filepath.Clean(mountPoint.GetPath()) != name {
				mountPoint.Path = name
				updateInstruction.MountPoint(mountPoint)
			}
		}

		updatedShare, err := s.UpdateShare(
			ctx,
			receivedShare.GetShare().GetId(),
			updateInstruction,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		acceptedShares = append(acceptedShares, updatedShare)
	}

	if len(receivedSharesResponse.GetShares()) == len(errs) {
		// none of the received shares could be accepted. This is an error. Return it.
		return acceptedShares, errors.Join(errs...)
	}

	return acceptedShares, nil
}

// DrivesDriveItemApi is the api that registers the http endpoints which expose needed operation to the graph api.
// the business logic is delegated to the space service and further down to the cs3 client.
type DrivesDriveItemApi struct {
	logger                 log.Logger
	drivesDriveItemService DrivesDriveItemProvider
	baseGraphService       BaseGraphProvider
}

// NewDrivesDriveItemApi creates a new DrivesDriveItemApi
func NewDrivesDriveItemApi(drivesDriveItemService DrivesDriveItemProvider, baseGraphService BaseGraphProvider, logger log.Logger) (DrivesDriveItemApi, error) {
	return DrivesDriveItemApi{
		logger:                 log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		drivesDriveItemService: drivesDriveItemService,
		baseGraphService:       baseGraphService,
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
		api.logger.Debug().Interface("driveID", driveID).Msg(ErrNotAShareJail.Error())
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, ErrNotAShareJail.Error())
		return
	}

	shareID, err := GetShareID(itemID)
	if err != nil {
		msg := "invalid shareID"
		api.logger.Debug().Interface("driveID", driveID).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
	}

	if err := api.drivesDriveItemService.UnmountShare(ctx, shareID); err != nil {
		msg := "unmounting share failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// UpdateDriveItem updates a drive item, currently only the visibility of the share is updated
func (api DrivesDriveItemApi) UpdateDriveItem(w http.ResponseWriter, r *http.Request) {
	driveID, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		msg := "invalid driveID or itemID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	if !IsShareJail(driveID) {
		api.logger.Debug().Interface("driveID", driveID).Msg(ErrNotAShareJail.Error())
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, ErrNotAShareJail.Error())
		return
	}

	shareID, err := GetShareID(itemID)
	if err != nil {
		msg := "invalid shareID"
		api.logger.Debug().Interface("driveID", driveID).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
	}

	requestDriveItem := libregraph.DriveItem{}
	if err := StrictJSONUnmarshal(r.Body, &requestDriveItem); err != nil {
		msg := "invalid request body"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	updatedShare, err := api.drivesDriveItemService.UpdateShare(
		r.Context(),
		shareID,
		(&ShareUpdateInstruction{}).
			Hidden(requestDriveItem.GetUIHidden()),
	)

	if err != nil {
		msg := "failed to update share"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedShare)

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
		api.logger.Debug().Interface("driveID", driveID).Msg(ErrNotAShareJail.Error())
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, ErrNotAShareJail.Error())
		return
	}

	requestDriveItem := libregraph.DriveItem{}
	if err := StrictJSONUnmarshal(r.Body, &requestDriveItem); err != nil {
		msg := "invalid request body"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, msg)
		return
	}

	remoteItem := requestDriveItem.GetRemoteItem()
	resourceId, err := storagespace.ParseID(remoteItem.GetId())
	if err != nil {
		msg := "invalid remote item id"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, msg)
		return
	}

	mountedShares, err := api.drivesDriveItemService.
		MountShare(ctx, resourceId, requestDriveItem.GetName())
	if err != nil {
		msg := "mounting share failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, msg)
		return
	}

	driveItems, err := api.baseGraphService.CS3ReceivedSharesToDriveItems(ctx, mountedShares)
	switch {
	case err != nil:
		break
	case len(driveItems) != 1:
		err = errorcode.New(errorcode.GeneralException, "failed to convert accepted shares into drive-item")
	}
	if err != nil {
		msg := "converting received shares to drive items failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, driveItems[0])
}

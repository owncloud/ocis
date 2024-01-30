package svc

import (
	"context"
	"errors"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

// DrivesDriveItemProvider is the interface that needs to be implemented by the individual space service
type DrivesDriveItemProvider interface {
	MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error)
}

// DrivesDriveItemService contains the production business logic for everything that relates to drives
type DrivesDriveItemService struct {
	logger          log.Logger
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// NewDrivesDriveItemService creates a new DrivesDriveItemService
func NewDrivesDriveItemService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (DrivesDriveItemService, error) {
	return DrivesDriveItemService{
		logger:          log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
		gatewaySelector: gatewaySelector,
	}, nil
}

// MountShare mounts a share
func (s DrivesDriveItemService) MountShare(ctx context.Context, resourceID storageprovider.ResourceId, name string) (libregraph.DriveItem, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return libregraph.DriveItem{}, err
	}

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

	for _, receivedShare := range receivedSharesResponse.GetShares() {
		updateMask := &fieldmaskpb.FieldMask{Paths: []string{"state"}}
		receivedShare.State = collaboration.ShareState_SHARE_STATE_ACCEPTED

		// only update if mountPoint name is not empty and the path has changed
		if name != "" {
			mountPoint := receivedShare.GetMountPoint()
			if mountPoint == nil {
				mountPoint = &storageprovider.Reference{}
			}

			newPath := utils.MakeRelativePath(name)
			if mountPoint.GetPath() != newPath {
				mountPoint.Path = newPath
				receivedShare.MountPoint = mountPoint
				updateMask.Paths = append(updateMask.Paths, "mount_point")
			}
		}

		updateReceivedShareRequest := &collaboration.UpdateReceivedShareRequest{
			Share:      receivedShare,
			UpdateMask: updateMask,
		}

		updateReceivedShareResponse, err := gatewayClient.UpdateReceivedShare(ctx, updateReceivedShareRequest)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// fixMe: send to nirvana, wait for toDriverItem func
		_ = updateReceivedShareResponse
	}

	// fixMe: return a concrete driveItem
	return libregraph.DriveItem{}, errors.Join(errs...)
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

// Routes returns the routes that should be registered for this api
func (api DrivesDriveItemApi) Routes() []Route {
	return []Route{
		{http.MethodPost, "/v1beta1/drives/{driveID}/items/{itemID}/children", api.CreateChildren},
	}
}

// CreateChildren exposes the CreateChildren operation of the space service as an http endpoint
func (api DrivesDriveItemApi) CreateChildren(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driveID, itemID, err := GetDriveAndItemIDParam(r, &api.logger)
	if err != nil {
		msg := "invalid driveID or itemID"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, msg)
		return
	}

	if !IsShareJail(driveID) || !IsShareJail(itemID) {
		msg := "invalid driveID or itemID, must be share jail"
		api.logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg(msg)
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

	responseDriveItem, err := api.drivesDriveItemService.
		MountShare(ctx, resourceId, remoteItem.GetName())
	if err != nil {
		msg := "mounting share failed"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, responseDriveItem)
}

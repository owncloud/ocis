package svc

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

// DrivesDriveItemProvider is the interface that needs to be implemented by the individual space service
type DrivesDriveItemProvider interface {
	CreateChildren(ctx context.Context, driveId, itemId storageprovider.ResourceId, driveItem libregraph.DriveItem) (libregraph.DriveItem, error)
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

// CreateChildren is currently only used for accepting pending//dangling shares.
// fixMe: currently the driveItem is not used, why is it needed?
func (s DrivesDriveItemService) CreateChildren(ctx context.Context, driveId, itemId storageprovider.ResourceId, _ libregraph.DriveItem) (libregraph.DriveItem, error) {
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
					ResourceId: &storageprovider.ResourceId{
						StorageId: driveId.GetStorageId(),
						SpaceId:   driveId.GetSpaceId(),
						OpaqueId:  itemId.GetOpaqueId(),
					},
				},
			},
		},
	})

	for _, receivedShare := range receivedSharesResponse.GetShares() {
		mountPoint := receivedShare.GetMountPoint()
		if mountPoint == nil {
			// fixMe: should not happen, add exception handling
			continue
		}

		receivedShare.State = collaboration.ShareState_SHARE_STATE_ACCEPTED
		receivedShare.MountPoint = &storageprovider.Reference{
			// keep the original mount point path,
			// custom path handling should come here later
			Path: mountPoint.GetPath(),
		}

		updateReceivedShareRequest := &collaboration.UpdateReceivedShareRequest{
			Share: receivedShare,
			// mount_point currently contain no changes, this is for future use
			// state changes from pending or rejected to accept
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"state", "mount_point"}},
		}

		//fixMe:  should be processed in parallel
		updateReceivedShareResponse, err := gatewayClient.UpdateReceivedShare(ctx, updateReceivedShareRequest)
		if err != nil {
			// fixMe: should not happen, add exception handling
			continue
		}

		// fixMe: send to nirvana, add status handling
		_ = updateReceivedShareResponse
	}

	// fixMe: return a concrete driveItem
	return libregraph.DriveItem{}, nil
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
		errorcode.RenderError(w, r, err)
		return
	}

	driveItem, err := api.drivesDriveItemService.
		CreateChildren(ctx, driveID, itemID, libregraph.DriveItem{})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, driveItem)
}

package svc

import (
	"context"
	"net/http"

	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

type DrivesDriveItemServicer interface {
	CreateChildren(ctx context.Context, driveId, itemId storageprovider.ResourceId, driveItem libregraph.DriveItem) (libregraph.DriveItem, error)
}

type DrivesDriveItemService struct {
	logger log.Logger
}

func NewDrivesDriveItemService(logger log.Logger) (DrivesDriveItemService, error) {
	return DrivesDriveItemService{
		logger: log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemService").Logger()},
	}, nil
}

func (s DrivesDriveItemService) CreateChildren(ctx context.Context, driveId, itemId storageprovider.ResourceId, driveItem libregraph.DriveItem) (libregraph.DriveItem, error) {
	return libregraph.DriveItem{}, nil
}

type DrivesDriveItemApi struct {
	logger                 log.Logger
	drivesDriveItemService DrivesDriveItemServicer
}

func NewDrivesDriveItemApi(drivesDriveItemService DrivesDriveItemServicer, logger log.Logger) (DrivesDriveItemApi, error) {
	return DrivesDriveItemApi{
		logger:                 log.Logger{Logger: logger.With().Str("graph api", "DrivesDriveItemApi").Logger()},
		drivesDriveItemService: drivesDriveItemService,
	}, nil
}

func (api DrivesDriveItemApi) Routes() []Route {
	return []Route{
		{http.MethodPost, "/v1beta1/drives/{driveID}/items/{itemID}/children", api.CreateChildren},
	}
}

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

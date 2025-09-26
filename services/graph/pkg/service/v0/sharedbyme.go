package svc

import (
	"net/http"

	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

type driveItemsByResourceID map[string]libregraph.DriveItem

// GetSharedByMe implements the Service interface (/me/drives/sharedByMe endpoint)
func (g Graph) GetSharedByMe(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msg("Calling GetRootDriveChildren")
	ctx := r.Context()

	driveItems := make(driveItemsByResourceID)
	var err error
	driveItems, err = g.listUserShares(ctx, nil, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	if g.config.IncludeOCMSharees {
		driveItems, err = g.listOCMShares(ctx, nil, driveItems)
		if err != nil {
			errorcode.RenderError(w, r, err)
			return
		}
	}

	driveItems, err = g.listPublicShares(ctx, nil, driveItems)
	if err != nil {
		errorcode.RenderError(w, r, err)
		return
	}

	res := make([]libregraph.DriveItem, 0, len(driveItems))
	for _, v := range driveItems {
		res = append(res, v)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: res})
}

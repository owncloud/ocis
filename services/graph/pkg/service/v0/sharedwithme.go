package svc

import (
	"context"
	"net/http"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

// ListSharedWithMe lists the files shared with the current user.
func (g Graph) ListSharedWithMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	driveItems, err := g.listSharedWithMe(ctx)
	if err != nil {
		g.logger.Error().Err(err).Msg("listSharedWithMe failed")
		errorcode.RenderError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: driveItems})
}

// listSharedWithMe is a helper function that lists the drive items shared with the current user.
func (g Graph) listSharedWithMe(ctx context.Context) ([]libregraph.DriveItem, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, err
	}

	listReceivedSharesResponse, err := gatewayClient.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{})
	if errCode := errorcode.FromCS3Status(listReceivedSharesResponse.GetStatus(), err); errCode != nil {
		g.logger.Error().Err(err).Msg("listing shares failed")
		return nil, *errCode
	}

	return cs3ReceivedSharesToDriveItems(ctx, g.logger, gatewayClient, g.identityCache, g.config.FilesSharing.EnableResharing, listReceivedSharesResponse.GetShares())
}

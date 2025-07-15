package svc

import (
	"context"
	"net/http"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	ocm "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
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
	if err := errorcode.FromCS3Status(listReceivedSharesResponse.GetStatus(), err); err != nil {
		g.logger.Error().Err(err).Msg("listing shares failed")
		return nil, err
	}
	availableRoles := unifiedrole.GetRoles(unifiedrole.RoleFilterIDs(g.config.UnifiedRoles.AvailableRoles...))
	driveItems, err := cs3ReceivedSharesToDriveItems(ctx, g.logger, gatewayClient, g.identityCache, listReceivedSharesResponse.GetShares(), availableRoles)
	if err != nil {
		g.logger.Error().Err(err).Msg("could not convert received shares to drive items")
		return nil, err
	}

	if g.config.IncludeOCMSharees {
		listReceivedOCMSharesResponse, err := gatewayClient.ListReceivedOCMShares(ctx, &ocm.ListReceivedOCMSharesRequest{})
		if err := errorcode.FromCS3Status(listReceivedSharesResponse.GetStatus(), err); err != nil {
			g.logger.Error().Err(err).Msg("listing shares failed")
			return nil, err
		}
		ocmDriveItems, err := cs3ReceivedOCMSharesToDriveItems(ctx, g.logger, gatewayClient, g.identityCache, listReceivedOCMSharesResponse.GetShares(), availableRoles)
		if err != nil {
			g.logger.Error().Err(err).Msg("could not convert received ocm shares to drive items")
			return nil, err
		}
		driveItems = append(driveItems, ocmDriveItems...)
	}

	return driveItems, err
}

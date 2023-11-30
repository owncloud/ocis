package svc

import (
	"context"
	"errors"
	"net/http"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/storagespace"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
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

func (g Graph) listUserShares(ctx context.Context, filters []*collaboration.Filter, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}

	concreteFilters := []*collaboration.Filter{
		share.UserGranteeFilter(),
		share.GroupGranteeFilter(),
	}
	concreteFilters = append(concreteFilters, filters...)

	lsUserSharesRequest := collaboration.ListSharesRequest{
		Filters: concreteFilters,
	}

	lsUserSharesResponse, err := gatewayClient.ListShares(ctx, &lsUserSharesRequest)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if statusCode := lsUserSharesResponse.GetStatus().GetCode(); statusCode != rpc.Code_CODE_OK {
		return driveItems, errorcode.New(cs3StatusToErrCode(statusCode), lsUserSharesResponse.Status.Message)
	}
	driveItems, err = g.cs3UserSharesToDriveItems(ctx, lsUserSharesResponse.Shares, driveItems)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	return driveItems, nil
}

func (g Graph) listPublicShares(ctx context.Context, filters []*link.ListPublicSharesRequest_Filter, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {

	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}

	var concreteFilters []*link.ListPublicSharesRequest_Filter
	concreteFilters = append(concreteFilters, filters...)

	req := link.ListPublicSharesRequest{
		Filters: concreteFilters,
	}

	lsPublicSharesResponse, err := gatewayClient.ListPublicShares(ctx, &req)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	if statusCode := lsPublicSharesResponse.GetStatus().GetCode(); statusCode != rpc.Code_CODE_OK {
		return driveItems, errorcode.New(cs3StatusToErrCode(statusCode), lsPublicSharesResponse.Status.Message)
	}
	driveItems, err = g.cs3PublicSharesToDriveItems(ctx, lsPublicSharesResponse.Share, driveItems)
	if err != nil {
		return driveItems, errorcode.New(errorcode.GeneralException, err.Error())
	}
	return driveItems, nil

}

func (g Graph) cs3UserSharesToDriveItems(ctx context.Context, shares []*collaboration.Share, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	for _, s := range shares {
		g.logger.Debug().Interface("CS3 UserShare", s).Msg("Got Share")
		resIDStr := storagespace.FormatResourceID(*s.ResourceId)
		item, ok := driveItems[resIDStr]
		if !ok {
			itemptr, err := g.getDriveItem(ctx, storageprovider.Reference{ResourceId: s.ResourceId})
			if err != nil {
				g.logger.Debug().Err(err).Interface("Share", s.ResourceId).Msg("could not stat share, skipping")
				continue
			}
			item = *itemptr
		}
		perm, err := g.cs3UserShareToPermission(ctx, s)

		var errcode errorcode.Error
		switch {
		case errors.As(err, &errcode) && errcode.GetCode() == errorcode.ItemNotFound:
			// The Grantee couldn't be found (user/group does not exist anymore)
			continue
		case err != nil:
			return driveItems, err
		}
		item.Permissions = append(item.Permissions, *perm)
		driveItems[resIDStr] = item
	}
	return driveItems, nil
}

func (g Graph) cs3UserShareToPermission(ctx context.Context, share *collaboration.Share) (*libregraph.Permission, error) {
	perm := libregraph.Permission{}
	perm.SetRoles([]string{})
	perm.SetId(share.Id.OpaqueId)
	grantedTo := libregraph.SharePointIdentitySet{}
	var li libregraph.Identity
	switch share.GetGrantee().GetType() {
	case storageprovider.GranteeType_GRANTEE_TYPE_USER:
		user, err := g.identityCache.GetUser(ctx, share.Grantee.GetUserId().GetOpaqueId())
		switch {
		case errors.Is(err, identity.ErrNotFound):
			g.logger.Warn().Str("userid", share.Grantee.GetUserId().GetOpaqueId()).Msg("User not found by id")
			// User does not seem to exist anymore, don't add a permission for this
			return nil, errorcode.New(errorcode.ItemNotFound, "grantee does not exist")
		case err != nil:
			return nil, errorcode.New(errorcode.GeneralException, err.Error())
		default:
			li.SetDisplayName(user.GetDisplayName())
			li.SetId(user.GetId())
			grantedTo.SetUser(li)
		}
	case storageprovider.GranteeType_GRANTEE_TYPE_GROUP:
		group, err := g.identityCache.GetGroup(ctx, share.Grantee.GetGroupId().GetOpaqueId())
		switch {
		case errors.Is(err, identity.ErrNotFound):
			g.logger.Warn().Str("groupid", share.Grantee.GetGroupId().GetOpaqueId()).Msg("Group not found by id")
			// Group not seem to exist anymore, don't add a permission for this
			return nil, errorcode.New(errorcode.ItemNotFound, "grantee does not exist")
		case err != nil:
			return nil, errorcode.New(errorcode.GeneralException, err.Error())
		default:
			li.SetDisplayName(group.GetDisplayName())
			li.SetId(group.GetId())
			grantedTo.SetGroup(li)
		}
	}

	// set expiration date
	if share.GetExpiration() != nil {
		perm.SetExpirationDateTime(cs3TimestampToTime(share.GetExpiration()))
	}
	role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(
		*share.GetPermissions().GetPermissions(),
		unifiedrole.UnifiedRoleConditionGrantee,
		g.config.FilesSharing.EnableResharing,
	)
	if role != nil {
		perm.SetRoles([]string{role.GetId()})
	} else {
		actions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(*share.GetPermissions().GetPermissions())
		perm.SetLibreGraphPermissionsActions(actions)
		perm.SetRoles(nil)
	}
	perm.SetGrantedToV2(grantedTo)
	return &perm, nil
}

func (g Graph) cs3PublicSharesToDriveItems(ctx context.Context, shares []*link.PublicShare, driveItems driveItemsByResourceID) (driveItemsByResourceID, error) {
	for _, s := range shares {
		g.logger.Debug().Interface("CS3 PublicShare", s).Msg("Got Share")
		resIDStr := storagespace.FormatResourceID(*s.ResourceId)
		item, ok := driveItems[resIDStr]
		if !ok {
			itemptr, err := g.getDriveItem(ctx, storageprovider.Reference{ResourceId: s.ResourceId})
			if err != nil {
				g.logger.Debug().Err(err).Interface("Share", s.ResourceId).Msg("could not stat share, skipping")
				continue
			}
			item = *itemptr
		}
		perm, err := g.libreGraphPermissionFromCS3PublicShare(s)
		if err != nil {
			g.logger.Error().Err(err).Interface("Link", s.ResourceId).Msg("could not convert link to libregraph")
			return driveItems, err
		}

		item.Permissions = append(item.Permissions, *perm)
		driveItems[resIDStr] = item
	}

	return driveItems, nil
}

func cs3StatusToErrCode(code rpc.Code) (errcode errorcode.ErrorCode) {
	switch code {
	case rpc.Code_CODE_UNAUTHENTICATED:
		errcode = errorcode.Unauthenticated
	case rpc.Code_CODE_PERMISSION_DENIED:
		errcode = errorcode.AccessDenied
	case rpc.Code_CODE_NOT_FOUND:
		errcode = errorcode.ItemNotFound
	case rpc.Code_CODE_LOCKED:
		errcode = errorcode.ItemIsLocked
	case rpc.Code_CODE_INVALID_ARGUMENT:
		errcode = errorcode.InvalidRequest
	case rpc.Code_CODE_FAILED_PRECONDITION:
		errcode = errorcode.InvalidRequest
	default:
		errcode = errorcode.GeneralException
	}
	return errcode
}

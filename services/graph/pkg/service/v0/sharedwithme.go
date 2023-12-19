package svc

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
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

	var driveItems []libregraph.DriveItem
	for _, receivedShare := range listReceivedSharesResponse.GetShares() {
		statRequest := &storageprovider.StatRequest{}

		switch receivedShare.GetState() {
		case collaboration.ShareState_SHARE_STATE_ACCEPTED:
			statRequest.Ref = &storageprovider.Reference{
				ResourceId: &storageprovider.ResourceId{
					StorageId: utils.ShareStorageProviderID,
					OpaqueId:  receivedShare.GetShare().GetId().GetOpaqueId(),
					SpaceId:   utils.ShareStorageSpaceID,
				},
			}
		case collaboration.ShareState_SHARE_STATE_PENDING:
			// return no remoteItem
			fallthrough
		case collaboration.ShareState_SHARE_STATE_REJECTED:
			// what to return here? same as pending?
			statRequest.Ref = &storageprovider.Reference{
				ResourceId: receivedShare.GetShare().GetResourceId(),
			}
		default:
			continue
		}

		statResponse, err := gatewayClient.Stat(ctx, statRequest)
		if errCode := errorcode.FromCS3Status(statResponse.GetStatus(), err); errCode != nil {
			g.logger.Error().Err(err).Msg("could not stat")
			continue
		}

		var commonResourceOwner *libregraph.Identity
		if userID := statResponse.GetInfo().GetOwner(); userID != nil {
			if user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId()); err != nil {
				g.logger.Error().Err(err).Msg("could not get user")
				continue
			} else {
				commonResourceOwner = &libregraph.Identity{
					DisplayName: user.GetDisplayName(),
					Id:          libregraph.PtrString(user.GetId()),
				}
			}
		}

		var commonShareCreator *libregraph.Identity
		if userID := receivedShare.GetShare().GetCreator(); userID != nil {
			if user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId()); err != nil {
				g.logger.Error().Err(err).Msg("could not get user")
				continue
			} else {
				commonShareCreator = &libregraph.Identity{
					DisplayName: user.GetDisplayName(),
					Id:          libregraph.PtrString(user.GetId()),
				}
			}
		}

		var commonPermission *libregraph.Permission
		{
			permission := libregraph.NewPermission()

			if id := receivedShare.GetShare().GetId().GetOpaqueId(); id != "" {
				permission.SetId(id)
			}

			if permissionSet := statResponse.GetInfo().GetPermissionSet(); permissionSet != nil {
				if actions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(*permissionSet); len(actions) > 0 {
					permission.SetLibreGraphPermissionsActions(actions)
				}

				if role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(
					*permissionSet,
					unifiedrole.UnifiedRoleConditionGrantee,
					g.config.FilesSharing.EnableResharing,
				); role != nil {
					permission.SetRoles([]string{role.GetId()})
				}
			}

			if expiration := receivedShare.GetShare().GetExpiration(); expiration != nil {
				permission.SetExpirationDateTime(cs3TimestampToTime(expiration))
			}

			switch grantee := receivedShare.GetShare().GetGrantee(); {
			case grantee.GetType() == storageprovider.GranteeType_GRANTEE_TYPE_USER:
				permission.SetGrantedToV2(libregraph.SharePointIdentitySet{
					User: &libregraph.Identity{
						Id: conversions.ToPointer(grantee.GetUserId().GetOpaqueId()),
					},
				})
			case grantee.GetType() == storageprovider.GranteeType_GRANTEE_TYPE_GROUP:
				permission.SetGrantedToV2(libregraph.SharePointIdentitySet{
					Group: &libregraph.Identity{
						Id: conversions.ToPointer(grantee.GetGroupId().GetOpaqueId()),
					},
				})
			}

			if !reflect.ValueOf(*permission).IsZero() {
				commonPermission = permission
			}
		}

		driveItem := libregraph.NewDriveItem()
		{
			if commonShareCreator != nil {
				driveItem.SetCreatedBy(libregraph.IdentitySet{
					User: commonShareCreator,
				})
			}

			{
				parentReference := libregraph.NewItemReference()

				if spaceType := statResponse.GetInfo().GetSpace().GetSpaceType(); spaceType != "" {
					parentReference.SetDriveType(spaceType)
				}

				if root := statResponse.GetInfo().GetSpace().GetRoot(); root != nil {
					parentReference.SetDriveId(storagespace.FormatResourceID(*root))
				}

				if !reflect.ValueOf(*parentReference).IsZero() {
					driveItem.ParentReference = parentReference
				}
			}
		}

		switch receivedShare.GetState() {
		case collaboration.ShareState_SHARE_STATE_ACCEPTED:
			if resourceID := statRequest.GetRef().GetResourceId(); resourceID != nil {
				driveItem.SetId(storagespace.FormatResourceID(*resourceID))
			}

			if name := receivedShare.GetMountPoint().GetPath(); name != "" {
				driveItem.SetName(receivedShare.GetMountPoint().GetPath())
			}

			if mTime := receivedShare.GetShare().GetMtime(); mTime != nil {
				driveItem.SetLastModifiedDateTime(cs3TimestampToTime(mTime))
			}

			if cTime := receivedShare.GetShare().GetCtime(); cTime != nil {
				driveItem.SetCreatedDateTime(cs3TimestampToTime(cTime))
			}

			{
				remoteItem := libregraph.NewRemoteItem()

				if id := statResponse.GetInfo().GetId(); id != nil {
					remoteItem.SetId(storagespace.FormatResourceID(*id))
				}

				if mTime := statResponse.GetInfo().GetMtime(); mTime != nil {
					remoteItem.SetLastModifiedDateTime(cs3TimestampToTime(mTime))
				}

				if name := statResponse.GetInfo().GetName(); name != "" {
					remoteItem.SetName(name)
				}

				if size := statResponse.GetInfo().GetSize(); size != 0 {
					remoteItem.SetSize(int64(size))
				}

				if etag := statResponse.GetInfo().GetEtag(); etag != "" {
					remoteItem.SetETag(strings.Trim(etag, "\""))
				}

				if commonResourceOwner != nil {
					remoteItem.SetCreatedBy(libregraph.IdentitySet{
						User: commonResourceOwner,
					})
				}

				switch info := statResponse.GetInfo(); {
				case info.GetType() == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER:
					remoteItem.Folder = libregraph.NewFolder()
				case info.GetType() == storageprovider.ResourceType_RESOURCE_TYPE_FILE:
					openGraphFile := libregraph.NewOpenGraphFile()

					if mimeType := info.GetMimeType(); mimeType != "" {
						openGraphFile.MimeType = &mimeType
					}

					remoteItem.File = openGraphFile
				case info.GetType() == storageprovider.ResourceType_RESOURCE_TYPE_INVALID:
					g.logger.Info().Interface("info", info).Msg("invalid resource type")
				}

				if commonPermission != nil {
					remoteItem.SetPermissions([]libregraph.Permission{*commonPermission})
				}

				{
					shared := libregraph.NewShared()

					if cTime := receivedShare.GetShare().GetCtime(); cTime != nil {
						shared.SetSharedDateTime(cs3TimestampToTime(cTime))
					}

					if commonResourceOwner != nil {
						shared.SetOwner(libregraph.IdentitySet{
							User: commonResourceOwner,
						})
					}

					if commonShareCreator != nil {
						shared.SetSharedBy(libregraph.IdentitySet{
							User: commonShareCreator,
						})
					}

					if !reflect.ValueOf(*shared).IsZero() {
						remoteItem.SetShared(*shared)
					}
				}

				if !reflect.ValueOf(*remoteItem).IsZero() {
					driveItem.SetRemoteItem(*remoteItem)
				}
			}
		case collaboration.ShareState_SHARE_STATE_PENDING:
			fallthrough
		case collaboration.ShareState_SHARE_STATE_REJECTED:
			if id := statResponse.GetInfo().GetId(); id != nil {
				driveItem.SetId(storagespace.FormatResourceID(*id))
			}

			if name := statResponse.GetInfo().GetName(); name != "" {
				driveItem.SetName(name)
			}

			if mTime := statResponse.GetInfo().GetMtime(); mTime != nil {
				driveItem.SetLastModifiedDateTime(cs3TimestampToTime(mTime))
			}

			if commonPermission != nil {
				driveItem.SetPermissions([]libregraph.Permission{*commonPermission})
			}
		}

		driveItems = append(driveItems, *driveItem)
	}

	return driveItems, nil
}

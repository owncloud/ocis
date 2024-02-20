package svc

import (
	"context"
	"net/http"
	"reflect"

	cs3User "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"golang.org/x/sync/errgroup"

	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"

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
	if errCode := errorcode.FromCS3Status(listReceivedSharesResponse.GetStatus(), err); errCode != nil {
		g.logger.Error().Err(err).Msg("listing shares failed")
		return nil, *errCode
	}

	return g.cs3ReceivedSharesToDriveItems(ctx, listReceivedSharesResponse.GetShares())
}

func (g Graph) cs3ReceivedSharesToDriveItems(ctx context.Context, receivedShares []*collaboration.ReceivedShare) ([]libregraph.DriveItem, error) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, err
	}

	// doStat is a helper function that stat a resource.
	doStat := func(resourceId *storageprovider.ResourceId) (*storageprovider.StatResponse, error) {
		shareStat, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{
			Ref: &storageprovider.Reference{ResourceId: resourceId},
		})
		switch errCode := errorcode.FromCS3Status(shareStat.GetStatus(), err); {
		case errCode == nil:
			break
		// skip ItemNotFound shares, they might have been deleted in the meantime or orphans.
		case errCode.GetCode() == errorcode.ItemNotFound:
			return nil, nil
		default:
			g.logger.Error().Err(errCode).Msg("could not stat")
			return nil, errCode
		}

		return shareStat, nil
	}

	ch := make(chan libregraph.DriveItem)
	group := new(errgroup.Group)
	// Set max concurrency
	group.SetLimit(10)

	receivedSharesByResourceID := make(map[string][]*collaboration.ReceivedShare, len(receivedShares))
	for _, receivedShare := range receivedShares {
		rIDStr := storagespace.FormatResourceID(*receivedShare.GetShare().GetResourceId())
		receivedSharesByResourceID[rIDStr] = append(receivedSharesByResourceID[rIDStr], receivedShare)
	}

	for _, receivedSharesForResource := range receivedSharesByResourceID {
		receivedShares := receivedSharesForResource

		group.Go(func() error {
			var err error // redeclare
			shareStat, err := doStat(receivedShares[0].GetShare().GetResourceId())
			if shareStat == nil || err != nil {
				return err
			}

			driveItem := libregraph.NewDriveItem()

			permissions := make([]libregraph.Permission, 0, len(receivedShares))

			var oldestReceivedShare *collaboration.ReceivedShare
			for _, receivedShare := range receivedShares {
				switch {
				case oldestReceivedShare == nil:
					fallthrough
				case utils.TSToTime(receivedShare.GetShare().GetCtime()).Before(utils.TSToTime(oldestReceivedShare.GetShare().GetCtime())):
					oldestReceivedShare = receivedShare
				}

				permission, err := g.cs3ReceivedShareToLibreGraphPermissions(ctx, receivedShare)
				if err != nil {
					return err
				}

				// If at least one of the shares was accepted, we consider the driveItem's synchronized
				// flag enabled.
				// Also we use the Mountpoint name of the first accepted mountpoint as the name of
				// of the driveItem
				if receivedShare.GetState() == collaboration.ShareState_SHARE_STATE_ACCEPTED {
					driveItem.SetClientSynchronize(true)
					if name := receivedShare.GetMountPoint().GetPath(); name != "" && driveItem.GetName() == "" {
						driveItem.SetName(receivedShare.GetMountPoint().GetPath())
					}
				}

				// if at least one share is marked as hidden, consider the whole driveItem to be hidden
				if receivedShare.GetHidden() {
					driveItem.SetUIHidden(true)
				}

				if userID := receivedShare.GetShare().GetCreator(); userID != nil {
					identity, err := g.cs3UserIdToIdentity(ctx, userID)
					if err != nil {
						g.logger.Warn().Err(err).Str("userid", userID.String()).Msg("could not get creator of the share")
					}

					permission.SetInvitation(
						libregraph.SharingInvitation{
							InvitedBy: &libregraph.IdentitySet{
								User: &identity,
							},
						},
					)
				}
				permissions = append(permissions, *permission)

			}

			// To stay compatible with the usershareprovider and the webdav
			// service the id of the driveItem is composed of the StorageID and
			// SpaceID of the sharestorage appended with the opaque ID of
			// the oldest share for the resource:
			// '<sharestorageid>$<sharespaceid>!<share-opaque-id>
			// Note: This means that the driveitem ID will change when the oldest
			//   shared is removed. It would be good to have are more stable ID here (e.g.
			//   derived from the shared resource's ID. But as we need to use the same
			//   ID across all services this means we needed to make similar adjustments
			//   to the sharejail (usershareprovider, webdav). Which we can't currently do
			//   as some clients rely on the IDs used there having a special format.
			driveItem.SetId(storagespace.FormatResourceID(storageprovider.ResourceId{
				StorageId: utils.ShareStorageProviderID,
				OpaqueId:  oldestReceivedShare.GetShare().GetId().GetOpaqueId(),
				SpaceId:   utils.ShareStorageSpaceID,
			}))

			if !driveItem.HasUIHidden() {
				driveItem.SetUIHidden(false)
			}
			if !driveItem.HasClientSynchronize() {
				driveItem.SetClientSynchronize(false)
				if name := shareStat.GetInfo().GetName(); name != "" {
					driveItem.SetName(name)
				}
			}

			remoteItem := libregraph.NewRemoteItem()
			{
				if id := shareStat.GetInfo().GetId(); id != nil {
					remoteItem.SetId(storagespace.FormatResourceID(*id))
				}

				if name := shareStat.GetInfo().GetName(); name != "" {
					remoteItem.SetName(name)
				}

				if etag := shareStat.GetInfo().GetEtag(); etag != "" {
					remoteItem.SetETag(etag)
				}

				if mTime := shareStat.GetInfo().GetMtime(); mTime != nil {
					remoteItem.SetLastModifiedDateTime(cs3TimestampToTime(mTime))
				}

				if size := shareStat.GetInfo().GetSize(); size != 0 {
					remoteItem.SetSize(int64(size))
				}

				parentReference := libregraph.NewItemReference()
				if spaceType := shareStat.GetInfo().GetSpace().GetSpaceType(); spaceType != "" {
					parentReference.SetDriveType(spaceType)
				}

				if root := shareStat.GetInfo().GetSpace().GetRoot(); root != nil {
					parentReference.SetDriveId(storagespace.FormatResourceID(*root))
				}
				if !reflect.ValueOf(*parentReference).IsZero() {
					remoteItem.ParentReference = parentReference
				}

			}

			// the parentReference of the outer driveItem should be the drive
			// containing the mountpoint i.e. the share jail
			driveItem.ParentReference = libregraph.NewItemReference()
			driveItem.ParentReference.SetDriveType("virtual")
			driveItem.ParentReference.SetDriveId(storagespace.FormatStorageID(utils.ShareStorageProviderID, utils.ShareStorageSpaceID))
			driveItem.ParentReference.SetId(storagespace.FormatResourceID(storageprovider.ResourceId{
				StorageId: utils.ShareStorageProviderID,
				OpaqueId:  utils.ShareStorageSpaceID,
				SpaceId:   utils.ShareStorageSpaceID,
			}))
			if etag := shareStat.GetInfo().GetEtag(); etag != "" {
				driveItem.SetETag(etag)
			}

			// connect the dots
			{
				if mTime := shareStat.GetInfo().GetMtime(); mTime != nil {
					t := cs3TimestampToTime(mTime)

					driveItem.SetLastModifiedDateTime(t)
					remoteItem.SetLastModifiedDateTime(t)
				}

				if size := shareStat.GetInfo().GetSize(); size != 0 {
					s := int64(size)

					driveItem.SetSize(s)
					remoteItem.SetSize(s)
				}

				if userID := shareStat.GetInfo().GetOwner(); userID != nil && userID.Type != cs3User.UserType_USER_TYPE_SPACE_OWNER {
					identity, err := g.cs3UserIdToIdentity(ctx, userID)
					if err != nil {
						// TODO: define a proper error behavior here. We don't
						// want the whole request to fail just because a single
						// resource owner couldn't be resolved. But, should be
						// really return the affect share in the response?
						// For now we just log a warning. The returned
						// identitySet will just contain the userid.
						g.logger.Warn().Err(err).Str("userid", userID.String()).Msg("could not get owner of shared resource")
					}

					remoteItem.SetCreatedBy(libregraph.IdentitySet{User: &identity})
					driveItem.SetCreatedBy(libregraph.IdentitySet{User: &identity})
				}
				switch info := shareStat.GetInfo(); {
				case info.GetType() == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER:
					folder := libregraph.NewFolder()

					remoteItem.Folder = folder
					driveItem.Folder = folder
				case info.GetType() == storageprovider.ResourceType_RESOURCE_TYPE_FILE:
					file := libregraph.NewOpenGraphFile()

					if mimeType := info.GetMimeType(); mimeType != "" {
						file.MimeType = &mimeType
					}

					remoteItem.File = file
					driveItem.File = file
				}

				if len(permissions) > 0 {
					remoteItem.Permissions = permissions
				}

				if !reflect.ValueOf(*remoteItem).IsZero() {
					driveItem.RemoteItem = remoteItem
				}
			}

			ch <- *driveItem

			return nil
		})
	}

	// wait for concurrent requests to finish
	go func() {
		err = group.Wait()
		close(ch)
	}()

	driveItems := make([]libregraph.DriveItem, 0, len(receivedSharesByResourceID))
	for di := range ch {
		driveItems = append(driveItems, di)
	}

	return driveItems, err
}

func (g Graph) cs3ReceivedShareToLibreGraphPermissions(ctx context.Context, receivedShare *collaboration.ReceivedShare) (*libregraph.Permission, error) {
	permission := libregraph.NewPermission()
	if id := receivedShare.GetShare().GetId().GetOpaqueId(); id != "" {
		permission.SetId(id)
	}

	if expiration := receivedShare.GetShare().GetExpiration(); expiration != nil {
		permission.SetExpirationDateTime(cs3TimestampToTime(expiration))
	}

	if permissionSet := receivedShare.GetShare().GetPermissions().GetPermissions(); permissionSet != nil {
		role := unifiedrole.CS3ResourcePermissionsToUnifiedRole(
			*permissionSet,
			unifiedrole.UnifiedRoleConditionGrantee,
			g.config.FilesSharing.EnableResharing,
		)

		if role != nil {
			permission.SetRoles([]string{role.GetId()})
		}

		actions := unifiedrole.CS3ResourcePermissionsToLibregraphActions(*permissionSet)

		// actions only make sense if no role is set
		if role == nil && len(actions) > 0 {
			permission.SetLibreGraphPermissionsActions(actions)
		}
	}

	switch grantee := receivedShare.GetShare().GetGrantee(); {
	case grantee.GetType() == storageprovider.GranteeType_GRANTEE_TYPE_USER:
		user, err := g.identityCache.GetUser(ctx, grantee.GetUserId().GetOpaqueId())
		if err != nil {
			g.logger.Error().Err(err).Msg("could not get user")
			return nil, err
		}

		permission.SetGrantedToV2(libregraph.SharePointIdentitySet{
			User: &libregraph.Identity{
				DisplayName: user.GetDisplayName(),
				Id:          user.Id,
			},
		})
	case grantee.GetType() == storageprovider.GranteeType_GRANTEE_TYPE_GROUP:
		group, err := g.identityCache.GetGroup(ctx, grantee.GetGroupId().GetOpaqueId())
		if err != nil {
			g.logger.Error().Err(err).Msg("could not get group")
			return nil, err
		}

		permission.SetGrantedToV2(libregraph.SharePointIdentitySet{
			Group: &libregraph.Identity{
				DisplayName: group.GetDisplayName(),
				Id:          group.Id,
			},
		})
	}

	return permission, nil
}

func (g Graph) cs3UserIdToIdentity(ctx context.Context, cs3UserID *cs3User.UserId) (libregraph.Identity, error) {
	identity := libregraph.Identity{
		Id: libregraph.PtrString(cs3UserID.GetOpaqueId()),
	}
	var err error
	if cs3UserID.GetType() != cs3User.UserType_USER_TYPE_SPACE_OWNER {
		var user libregraph.User
		user, err = g.identityCache.GetUser(ctx, cs3UserID.GetOpaqueId())
		if err == nil {
			identity.SetDisplayName(user.GetDisplayName())
		}
	}
	return identity, err
}

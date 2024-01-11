package svc

import (
	"context"
	"net/http"
	"reflect"
	"slices"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"golang.org/x/sync/errgroup"

	"github.com/cs3org/reva/v2/pkg/storagespace"

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

	group := new(errgroup.Group)
	receivedShares := listReceivedSharesResponse.GetShares()
	driveItems := make([]libregraph.DriveItem, len(receivedShares))

	for i, receivedShare := range receivedShares {
		i, receivedShare := i, receivedShare

		group.Go(func() error {
			shareStat, err := doStat(receivedShare.GetShare().GetResourceId())
			if shareStat == nil || err != nil {
				return err
			}

			permission, err := g.cs3ShareToLibreGraphPermissions(ctx, receivedShare.GetShare(), shareStat.GetInfo())
			if err != nil {
				return err
			}

			parentReference := libregraph.NewItemReference()
			{
				if spaceType := shareStat.GetInfo().GetSpace().GetSpaceType(); spaceType != "" {
					parentReference.SetDriveType(spaceType)
				}

				if root := shareStat.GetInfo().GetSpace().GetRoot(); root != nil {
					parentReference.SetDriveId(storagespace.FormatResourceID(*root))
				}
			}

			shared := libregraph.NewShared()
			{
				if cTime := receivedShare.GetShare().GetCtime(); cTime != nil {
					shared.SetSharedDateTime(cs3TimestampToTime(cTime))
				}
			}

			remoteItem := libregraph.NewRemoteItem()
			{
				remoteItem.SetUiHidden(receivedShare.GetHidden())
				remoteItem.SetClientSynchronize(receivedShare.GetState() == collaboration.ShareState_SHARE_STATE_ACCEPTED)

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
			}

			driveItem := libregraph.NewDriveItem()

			// handle share state related stuff
			switch receivedShare.GetState() {
			case collaboration.ShareState_SHARE_STATE_ACCEPTED:

				driveItem.SetId(storagespace.FormatResourceID(storageprovider.ResourceId{
					StorageId: utils.ShareStorageProviderID,
					OpaqueId:  receivedShare.GetShare().GetId().GetOpaqueId(),
					SpaceId:   utils.ShareStorageSpaceID,
				}))

				if name := receivedShare.GetMountPoint().GetPath(); name != "" {
					driveItem.SetName(receivedShare.GetMountPoint().GetPath())
				}

				if etag := shareStat.GetInfo().GetEtag(); etag != "" {
					driveItem.SetETag(etag)
				}
			case collaboration.ShareState_SHARE_STATE_PENDING:
				fallthrough
			case collaboration.ShareState_SHARE_STATE_REJECTED:
				if name := shareStat.GetInfo().GetName(); name != "" {
					driveItem.SetName(name)
				}
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

				if userID := shareStat.GetInfo().GetOwner(); userID != nil {
					user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId())
					if err != nil {
						g.logger.Error().Err(err).Msg("could not get user")
						return err
					}

					identitySet := libregraph.IdentitySet{
						User: &libregraph.Identity{
							DisplayName: user.GetDisplayName(),
							Id:          libregraph.PtrString(user.GetId()),
						},
					}

					remoteItem.SetCreatedBy(identitySet)
				}

				if userID := receivedShare.GetShare().GetOwner(); userID != nil {
					user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId())
					if err != nil {
						g.logger.Error().Err(err).Msg("could not get user")
						return err
					}

					identitySet := libregraph.IdentitySet{
						User: &libregraph.Identity{
							DisplayName: user.GetDisplayName(),
							Id:          libregraph.PtrString(user.GetId()),
						},
					}

					shared.SetOwner(identitySet)
				}

				if userID := receivedShare.GetShare().GetCreator(); userID != nil {
					user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId())
					if err != nil {
						g.logger.Error().Err(err).Msg("could not get user")
						return err
					}

					identitySet := libregraph.IdentitySet{
						User: &libregraph.Identity{
							DisplayName: user.GetDisplayName(),
							Id:          libregraph.PtrString(user.GetId()),
						},
					}

					driveItem.SetCreatedBy(identitySet)
					shared.SetSharedBy(identitySet)

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

				if !reflect.ValueOf(*shared).IsZero() {
					remoteItem.Shared = shared
				}

				if !reflect.ValueOf(*permission).IsZero() {
					permissions := []libregraph.Permission{*permission}

					remoteItem.Permissions = permissions
				}

				if !reflect.ValueOf(*parentReference).IsZero() {
					remoteItem.ParentReference = parentReference
					driveItem.ParentReference = parentReference
				}

				if !reflect.ValueOf(*remoteItem).IsZero() {
					driveItem.RemoteItem = remoteItem
				}
			}

			driveItems[i] = *driveItem

			return nil
		})
	}

	// wait for concurrent requests to finish
	err = group.Wait()

	// filter out empty drive items
	return slices.Clip(slices.DeleteFunc(driveItems, func(item libregraph.DriveItem) bool {
		return reflect.ValueOf(item).IsZero()
	})), err
}

func (g Graph) cs3ShareToLibreGraphPermissions(ctx context.Context, share *collaboration.Share, shareStatInfo *storageprovider.ResourceInfo) (*libregraph.Permission, error) {
	permission := libregraph.NewPermission()
	if id := share.GetId().GetOpaqueId(); id != "" {
		permission.SetId(id)
	}

	if expiration := share.GetExpiration(); expiration != nil {
		permission.SetExpirationDateTime(cs3TimestampToTime(expiration))
	}

	if permissionSet := shareStatInfo.GetPermissionSet(); permissionSet != nil {
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

	switch grantee := share.GetGrantee(); {
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

package svc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3User "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	v1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"golang.org/x/sync/errgroup"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

// StrictJSONUnmarshal is a wrapper around json.Unmarshal that returns an error if the json contains unknown fields.
func StrictJSONUnmarshal(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// IsSpaceRoot returns true if the resourceID is a space root.
func IsSpaceRoot(rid *storageprovider.ResourceId) bool {
	if rid == nil {
		return false
	}
	if rid.GetSpaceId() == "" || rid.GetOpaqueId() == "" {
		return false
	}

	return rid.GetSpaceId() == rid.GetOpaqueId()
}

// GetDriveAndItemIDParam parses the driveID and itemID from the request,
// validates the common fields and returns the parsed IDs if ok.
func GetDriveAndItemIDParam(r *http.Request, logger *log.Logger) (storageprovider.ResourceId, storageprovider.ResourceId, error) {
	empty := storageprovider.ResourceId{}

	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		logger.Debug().Err(err).Msg("could not parse driveID")
		return empty, empty, errorcode.New(errorcode.InvalidRequest, "invalid driveID")
	}

	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		logger.Debug().Err(err).Msg("could not parse itemID")
		return empty, empty, errorcode.New(errorcode.InvalidRequest, "invalid itemID")
	}

	if itemID.GetOpaqueId() == "" {
		logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg("empty item opaqueID")
		return empty, empty, errorcode.New(errorcode.InvalidRequest, "invalid itemID")
	}

	if driveID.GetStorageId() != itemID.GetStorageId() || driveID.GetSpaceId() != itemID.GetSpaceId() {
		logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg("driveID and itemID do not match")
		return empty, empty, errorcode.New(errorcode.ItemNotFound, "driveID and itemID do not match")
	}

	return driveID, itemID, nil
}

// GetGatewayClient returns a gateway client from the gatewaySelector.
func (g Graph) GetGatewayClient(w http.ResponseWriter, r *http.Request) (gateway.GatewayAPIClient, bool) {
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		g.logger.Debug().Err(err).Msg("selecting gatewaySelector failed")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return nil, false
	}

	return gatewayClient, true
}

// IsShareJail returns true if given id is a share jail id.
func IsShareJail(id storageprovider.ResourceId) bool {
	return id.GetStorageId() == utils.ShareStorageProviderID && id.GetSpaceId() == utils.ShareStorageSpaceID
}

// userIdToIdentity looks the user for the supplied id using the cache and returns it
// as a libregraph.Identity
func userIdToIdentity(ctx context.Context, cache identity.IdentityCache, userID string) (libregraph.Identity, error) {
	identity := libregraph.Identity{
		Id: libregraph.PtrString(userID),
	}
	user, err := cache.GetUser(ctx, userID)
	if err == nil {
		identity.SetDisplayName(user.GetDisplayName())
	}
	return identity, err
}

// cs3UserIdToIdentity looks up the user for the supplied cs3 userid using the cache and returns it
// as a libregraph.Identity. Skips the user lookup if the id type is USER_TYPE_SPACE_OWNER
func cs3UserIdToIdentity(ctx context.Context, cache identity.IdentityCache, cs3UserID *cs3User.UserId) (libregraph.Identity, error) {
	if cs3UserID.GetType() != cs3User.UserType_USER_TYPE_SPACE_OWNER {
		return userIdToIdentity(ctx, cache, cs3UserID.GetOpaqueId())
	}
	return libregraph.Identity{Id: libregraph.PtrString(cs3UserID.GetOpaqueId())}, nil
}

// groupIdToIdentity looks up the group for the supplied cs3 groupid using the cache and returns it
// as a libregraph.Identity.
func groupIdToIdentity(ctx context.Context, cache identity.IdentityCache, groupID string) (libregraph.Identity, error) {
	identity := libregraph.Identity{
		Id: libregraph.PtrString(groupID),
	}
	group, err := cache.GetGroup(ctx, groupID)
	if err == nil {
		identity.SetDisplayName(group.GetDisplayName())
	}
	return identity, err
}

func cs3ReceivedSharesToDriveItems(ctx context.Context,
	logger *log.Logger,
	gatewayClient gateway.GatewayAPIClient,
	identityCache identity.IdentityCache,
	resharing bool,
	receivedShares []*collaboration.ReceivedShare) ([]libregraph.DriveItem, error) {

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
			logger.Error().Err(errCode).Msg("could not stat")
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

				permission, err := cs3ReceivedShareToLibreGraphPermissions(ctx, logger, identityCache, resharing, receivedShare)
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
					identity, err := cs3UserIdToIdentity(ctx, identityCache, userID)
					if err != nil {
						logger.Warn().Err(err).Str("userid", userID.String()).Msg("could not get creator of the share")
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
			//   share is removed. It would be good to have are more stable ID here (e.g.
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
					identity, err := cs3UserIdToIdentity(ctx, identityCache, userID)
					if err != nil {
						// TODO: define a proper error behavior here. We don't
						// want the whole request to fail just because a single
						// resource owner couldn't be resolved. But, should we
						// really return the affected share in the response?
						// For now we just log a warning. The returned
						// identitySet will just contain the userid.
						logger.Warn().Err(err).Str("userid", userID.String()).Msg("could not get owner of shared resource")
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

	var err error
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

func cs3ReceivedShareToLibreGraphPermissions(ctx context.Context, logger *log.Logger,
	identityCache identity.IdentityCache, resharing bool, receivedShare *collaboration.ReceivedShare) (*libregraph.Permission, error) {
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
			resharing,
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
		user, err := cs3UserIdToIdentity(ctx, identityCache, grantee.GetUserId())
		if err != nil {
			logger.Error().Err(err).Msg("could not get user")
			return nil, err
		}
		permission.SetGrantedToV2(libregraph.SharePointIdentitySet{User: &user})
	case grantee.GetType() == storageprovider.GranteeType_GRANTEE_TYPE_GROUP:
		group, err := groupIdToIdentity(ctx, identityCache, grantee.GetGroupId().GetOpaqueId())
		if err != nil {
			logger.Error().Err(err).Msg("could not get group")
			return nil, err
		}
		permission.SetGrantedToV2(libregraph.SharePointIdentitySet{Group: &group})
	}

	return permission, nil
}

func (g Graph) getCS3Client(gwc gateway.GatewayAPIClient, root *storageprovider.ResourceId) (context.Context, *metadata.CS3, error) {
	mdc := metadata.NewCS3(g.config.Reva.Address, "com.owncloud.api.storage-users")
	mdc.SpaceRoot = root
	ctx, err := utils.GetServiceUserContext(g.config.ServiceAccount.ServiceAccountID, gwc, g.config.ServiceAccount.ServiceAccountSecret)
	return ctx, mdc, err
}

func applyTemplate(ctx context.Context, mdc *metadata.CS3, gwc gateway.GatewayAPIClient, root *storageprovider.ResourceId, fsys fs.ReadDirFS) error {
	entries, err := fsys.ReadDir(".")
	if err != nil {
		return err
	}

	updateSpaceRequest := &storageprovider.UpdateStorageSpaceRequest{
		// Prepare the object to apply the diff from. The properties on StorageSpace will overwrite
		// the original storage space.
		StorageSpace: &storageprovider.StorageSpace{
			Id: &storageprovider.StorageSpaceId{
				OpaqueId: storagespace.FormatResourceID(*root),
			},
			Root: root,
		},
	}

	updateSpaceRequest.StorageSpace.Opaque, err = uploadFolder(ctx, mdc, "", "", updateSpaceRequest.StorageSpace.Opaque, fsys, entries)
	if err != nil {
		return err
	}

	if len(updateSpaceRequest.StorageSpace.Opaque.Map) == 0 {
		return nil
	}

	resp, err := gwc.UpdateStorageSpace(ctx, updateSpaceRequest)
	switch {
	case err != nil:
		return err
	case resp.Status.Code == rpc.Code_CODE_OK:
		return nil
	default:
		return errors.New(resp.Status.Message)
	}
}

func uploadFolder(ctx context.Context, mdc *metadata.CS3, pathOnDisc, pathOnSpace string, opaque *v1beta1.Opaque, fsys fs.ReadDirFS, entries []os.DirEntry) (*v1beta1.Opaque, error) {
	for _, entry := range entries {
		spacePath := filepath.Join(pathOnSpace, entry.Name())
		discPath := filepath.Join(pathOnDisc, entry.Name())

		if entry.IsDir() {
			err := mdc.MakeDirIfNotExist(ctx, spacePath)
			if err != nil {
				return opaque, err
			}

			entries, err := fsys.ReadDir(discPath)
			if err != nil {
				return opaque, err
			}

			opaque, err = uploadFolder(ctx, mdc, discPath, spacePath, opaque, fsys, entries)
			if err != nil {
				return opaque, err
			}
			continue
		}

		b, err := fs.ReadFile(fsys, discPath)
		if err != nil {
			return opaque, err
		}

		if err := mdc.SimpleUpload(ctx, spacePath, b); err != nil {
			return opaque, err
		}

		// TODO: use upload to avoid second stat
		i, err := mdc.Stat(ctx, spacePath)
		if err != nil {
			return opaque, err
		}

		identifier := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		opaque = utils.AppendPlainToOpaque(opaque, identifier, storagespace.FormatResourceID(*i.Id))
	}

	return opaque, nil
}

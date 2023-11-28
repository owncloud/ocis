package svc

import (
	"context"
	"net/http"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
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
	if err != nil {
		g.logger.Error().Err(err).Msg("listing shares failed")
		return nil, errorcode.New(errorcode.GeneralException, err.Error())
	}

	switch listReceivedSharesResponse.Status.Code {
	case rpc.Code_CODE_NOT_FOUND:
		return nil, identity.ErrNotFound
	}

	var driveItems []libregraph.DriveItem
	for _, receivedShare := range listReceivedSharesResponse.GetShares() {
		share := receivedShare.GetShare()
		if share == nil {
			g.logger.Error().Interface("ListReceivedShares", listReceivedSharesResponse).Msg("unexpected empty ReceivedShare.Share")
			continue
		}

		driveItem := &libregraph.DriveItem{}

		statResponse, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: share.GetResourceId()}})
		if err != nil {
			g.logger.Error().Err(err).Msg("could not stat")
			continue
		}
		if statResponse.GetStatus().GetCode() != rpc.Code_CODE_OK {
			g.logger.Error().Err(err).Msg("invalid stat response")
			continue
		}
		resourceInfo := statResponse.GetInfo()

		var driveOwner *libregraph.Identity
		if userID := statResponse.GetInfo().GetOwner(); userID != nil {
			if user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId()); err != nil {
				g.logger.Error().Err(err).Msg("could not get user")
				continue
			} else {
				driveOwner = &libregraph.Identity{
					DisplayName: user.GetDisplayName(),
					Id:          libregraph.PtrString(user.GetId()),
				}
			}
		}

		var shareCreator *libregraph.Identity
		if userID := share.GetCreator(); userID != nil {
			if user, err := g.identityCache.GetUser(ctx, userID.GetOpaqueId()); err != nil {
				g.logger.Error().Err(err).Msg("could not get user")
				continue
			} else {
				shareCreator = &libregraph.Identity{
					DisplayName: user.GetDisplayName(),
					Id:          libregraph.PtrString(user.GetId()),
				}
			}
		}

		if cTime := share.GetCtime(); cTime != nil {
			driveItem.CreatedDateTime = libregraph.PtrTime(cs3TimestampToTime(cTime))
		}

		driveItem.ETag = libregraph.PtrString(strings.Trim(statResponse.GetInfo().GetEtag(), "\""))

		if id := share.GetId().GetOpaqueId(); id != "" {
			driveItem.Id = libregraph.PtrString(id)
		}

		if mTime := share.GetMtime(); mTime != nil {
			driveItem.LastModifiedDateTime = libregraph.PtrTime(cs3TimestampToTime(mTime))
		}

		if name := resourceInfo.GetName(); name != "" {
			driveItem.Name = libregraph.PtrString(name)
		}

		{
			addParentReference := false
			parentReference := &libregraph.ItemReference{}

			if id := share.GetId().GetOpaqueId(); id != "" {
				parentReference.DriveId = libregraph.PtrString(id)
				addParentReference = true
			}

			if addParentReference {
				driveItem.ParentReference = parentReference
			}
		}

		{
			remoteItem := &libregraph.RemoteItem{}

			if id := resourceInfo.GetId(); id != nil {
				remoteItem.Id = libregraph.PtrString(storagespace.FormatResourceID(*id))
			}

			if mTime := resourceInfo.GetMtime(); mTime != nil {
				remoteItem.LastModifiedDateTime = libregraph.PtrTime(cs3TimestampToTime(mTime))
			}

			if name := resourceInfo.GetName(); name != "" {
				remoteItem.Name = libregraph.PtrString(name)
			}

			// fixMe:
			// - negative permission could distort the size, am i right?
			remoteItem.Size = libregraph.PtrInt64(int64(resourceInfo.GetSize()))

			remoteItem.CreatedBy = &libregraph.IdentitySet{
				User: driveOwner,
			}

			{

				addFileSystemInfo := false
				fileSystemInfo := &libregraph.FileSystemInfo{}

				if cTime := share.GetCtime(); cTime != nil {
					// fixMe:
					// - ms uses the root resource ctime for that,
					//   the stat response does not contain any information about this, use share instead?
					fileSystemInfo.CreatedDateTime = libregraph.PtrTime(cs3TimestampToTime(cTime))
					addFileSystemInfo = true
				}

				if mTime := resourceInfo.GetMtime(); mTime != nil {
					fileSystemInfo.LastModifiedDateTime = libregraph.PtrTime(cs3TimestampToTime(mTime))
					addFileSystemInfo = true
				}

				if addFileSystemInfo {
					remoteItem.FileSystemInfo = fileSystemInfo
				}
			}

			switch resourceInfo.GetType() {
			case storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER:
				remoteItem.Folder = &libregraph.Folder{}
			case storageprovider.ResourceType_RESOURCE_TYPE_FILE:
				openGraphFile := &libregraph.OpenGraphFile{}

				if mimeType := resourceInfo.GetMimeType(); mimeType != "" {
					openGraphFile.MimeType = libregraph.PtrString(mimeType)
				}

				remoteItem.File = openGraphFile
			case storageprovider.ResourceType_RESOURCE_TYPE_INVALID:
				g.logger.Error().Msg("invalid resource type")
				continue
			}

			{
				addShared := false
				shared := &libregraph.Shared{
					Owner: &libregraph.IdentitySet{
						User: shareCreator,
					},
					SharedBy: &libregraph.IdentitySet{
						User: shareCreator,
					},
				}

				if cTime := share.GetCtime(); cTime != nil {
					shared.SharedDateTime = libregraph.PtrTime(cs3TimestampToTime(cTime))
					addShared = true
				}

				if shareCreator != nil {
					shared.Owner.User = shareCreator
					shared.SharedBy.User = shareCreator
					addShared = true
				}

				if addShared {
					remoteItem.Shared = shared
				}
			}

			driveItem.RemoteItem = remoteItem
		}

		driveItems = append(driveItems, *driveItem)
	}

	return driveItems, nil
}

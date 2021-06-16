package svc

import (
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"

	msgraph "github.com/owncloud/ocis/graph/pkg/openapi/v0"
	opengraph "github.com/owncloud/ocis/graph/pkg/openapi/v0"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
)

// GetDrives implements the Service interface.
func (g Graph) GetDrives(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msg("Calling GetDrives")
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &storageprovider.ListStorageSpacesRequest{

		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: "1284d238-aa92-42ce-bdc4-0b0000009157!*", // FIXME dynamically discover home and other storages ... actually the storage registry should provide the list of storage spaces
					},
				},
			},
		},
	}

	res, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msg("error sending list storage spaces grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.HandleErrorStatus(&g.logger.Logger, w, res.Status)
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := formatDrives(wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

func (g Graph) RootRouter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		restOfURL := chi.URLParam(r, "*")
		switch {
		case strings.HasSuffix(restOfURL, "/children"):
			g.GetChildren(w, r)
		default:
			g.GetDriveItem(w, r)
		}
	})
}

// GetDrive implements the Service interface.
func (g Graph) GetDrive(w http.ResponseWriter, r *http.Request) {
	driveID := chi.URLParam(r, "drive-id")
	if driveID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest)
		return
	}
	g.logger.Debug().Str("drive-id", driveID).Msg("Calling GetDrive")
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &storageprovider.ListStorageSpacesRequest{
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: driveID,
					},
				},
			},
		},
	}

	res, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msg("error sending list storage spaces grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.HandleErrorStatus(&g.logger.Logger, w, res.Status)
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := formatDrives(wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

// GetPersonalDriveChildren implements the Service interface.
func (g Graph) GetPersonalDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msgf("Calling GetPersonalDriveChildren")

	ctx := r.Context()

	fn := "/home"

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &storageprovider.ListContainerRequest{Ref: &storageprovider.Reference{Path: fn}}

	res, err := client.ListContainer(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error sending list container grpc request %s", fn)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.HandleErrorStatus(&g.logger.Logger, w, res.Status)
		return
	}

	files, err := formatDriveItems(res.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error encoding response as json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

// GetChildren implements the Service interface.
func (g Graph) GetChildren(w http.ResponseWriter, r *http.Request) {
	driveID := chi.URLParam(r, "drive-id")
	relPath := chi.URLParam(r, "relative-path")

	g.logger.Debug().Str("drive-id", driveID).Str("relative-path", relPath).Msg("Calling GetDrive")
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	parts := strings.SplitN(driveID, "!", 2)
	if len(parts) != 2 {
		g.logger.Err(err).Msg("invalid drive id, must contain a !")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := &storageprovider.ListContainerRequest{Ref: &storageprovider.Reference{
		ResourceId: &storageprovider.ResourceId{
			StorageId: parts[0],
			OpaqueId:  parts[1],
		},
		Path: relPath,
	}}

	res, err := client.ListContainer(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error sending list container grpc request %s", relPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.HandleErrorStatus(&g.logger.Logger, w, res.Status)
		return
	}

	files, err := formatDriveItems(res.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error encoding response as json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

// GetPersonalDriveChildren implements the Service interface.
func (g Graph) GetDriveItem(w http.ResponseWriter, r *http.Request) {
	g.logger.Debug().Msgf("Calling GetPersonalDriveChildren")

	ctx := r.Context()

	fn := "/home"

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req := &storageprovider.ListContainerRequest{Ref: &storageprovider.Reference{Path: fn}}

	res, err := client.ListContainer(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error sending list container grpc request %s", fn)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		errorcode.HandleErrorStatus(&g.logger.Logger, w, res.Status)
		return
	}

	files, err := formatDriveItems(res.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error encoding response as json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

func cs3TimestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

func cs3ResourceToDriveItem(res *storageprovider.ResourceInfo) (*msgraph.DriveItem, error) {
	size := new(int64)
	*size = int64(res.Size) // uint64 -> int :boom:
	name := path.Base(res.Path)

	id := res.Id.StorageId + "!" + res.Id.OpaqueId

	driveItem := &msgraph.DriveItem{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Id: &id,
			},
			Name: &name,
			// Note: The eTag and cTag properties work differently on containers (folders). The cTag
			// value is modified when content or metadata of any descendant of the folder is changed.
			// The eTag value is only modified when the folder's properties are changed, except for
			// properties that are derived from descendants (like childCount or lastModifiedDateTime).
			ETag: &res.Etag, // should we drop the enclosing "? ms graph api does not have them
		},
		Size: size, // oc:size
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		driveItem.BaseItem.LastModifiedDateTime = &lastModified
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE {
		driveItem.File = &msgraph.File{
			// TODO add file facet https://docs.microsoft.com/en-us/graph/api/resources/file?view=graph-rest-1.0
			// hashes	Hashes	Hashes of the file's binary content, if available. Read-only.
			// set below
			// mimeType	string	The MIME type for the file. This is determined by logic on the server and might not be the value provided when the file was uploaded. Read-only.
			MimeType: &res.MimeType,
		}
		// oc:checksum
		switch res.Checksum.Type {
		case storageprovider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_SHA1:
			driveItem.File.Hashes = &msgraph.Hashes{
				Sha1Hash: &res.Checksum.Sum,
			}
			// cTag	String	An eTag for the content of the item. This eTag is not changed if only the metadata is changed. Note This property is not returned if the item is a folder. Read-only.
			driveItem.CTag = &res.Checksum.Sum
		case storageprovider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_ADLER32:
			// TODO add to opengraph
		case storageprovider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_MD5:
			// TODO add to opengraph
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &msgraph.Folder{
			// TODO add folder facet properties https://docs.microsoft.com/en-us/graph/api/resources/folder?view=graph-rest-1.0
			// childCount	Int32	Number of children contained immediately within this container.
			// view	folderView	A collection of properties defining the recommended view for the folder.
		}
		// Note: The eTag and cTag properties work differently on containers (folders). The cTag
		// value is modified when content or metadata of any descendant of the folder is changed.
		// The eTag value is only modified when the folder's properties are changed, except for
		// properties that are derived from descendants (like childCount or lastModifiedDateTime).
		driveItem.CTag = &res.Etag
	}
	// in comparison to propfind
	// oc:permissions?
	//   permissions	permission collection	The set of permissions for the item. Read-only. Nullable.
	//   https://docs.microsoft.com/en-us/graph/api/driveitem-list-permissions?view=graph-rest-1.0&tabs=http

	// http://open-collaboration-services.org/ns/share-permissions
	// oc:public-link-permission
	// oc:public-link-item-type
	// oc:public-link-share-datetime
	// oc:public-link-share-owner
	// oc:public-link-expiration
	//  remoteItem	remoteItem	Remote item data, if the item is shared from a drive other than the one being accessed. Read-only.
	//  shared	shared	Indicates that the item has been shared with others and provides information about the shared state of the item. Read-only.
	//  -> used to indicate a file has been shared in the row of files, when selecting the file
	//     GET /me/drive/items/{item-id}/permissions can be used to list the actual permissions

	// oc:favorite
	// oc:owner-id
	// oc:share-types
	// oc:owner-display-name
	// oc:downloadURL -> webdavURL?
	// TODO driveItem.WebDavUrl
	//   @microsoft.graph.downloadUrl	string	A URL that can be used to download this file's content. Authentication is not required with this URL. Read-only.
	//   webDavUrl	String	WebDAV compatible URL for the item.
	//   webUrl	String	URL that displays the resource in the browser. Read-only.
	// oc:privatelink
	// oc:dDC for the desktop? used?
	// oc:data-fingerprint
	//   used by admins to indicate a backup has been restored,
	//   can only occur on the root node
	//   server implementation in https://github.com/owncloud/core/pull/24054
	//   see https://doc.owncloud.com/server/admin_manual/configuration/server/occ_command.html#maintenance-commands
	//   TODO(jfd): double check the client behavior with reva on backup restore

	// quota? -> drive property https://docs.microsoft.com/en-us/graph/api/resources/drive?view=graph-rest-1.0#properties
	// quota-used-bytes
	// quota-available-bytes

	// TODO image facet
	// TODO photo facet
	// TODO root If this property is non-null, it indicates that the driveItem is the top-most driveItem in the drive.

	// TODO created by
	// TODO last modified by
	// TODO parentReference

	// TODO thumbnails	thumbnailSet collection	Collection containing ThumbnailSet objects associated with the item. For more info, see getting thumbnails. Read-only. Nullable.
	//   https://docs.microsoft.com/en-us/graph/api/driveitem-list-thumbnails?view=graph-rest-1.0&tabs=http
	// TODO versions	driveItemVersion collection	The list of previous versions of the item. For more info, see getting previous versions. Read-only. Nullable.
	//   https://docs.microsoft.com/en-us/graph/api/driveitem-list-versions?view=graph-rest-1.0

	return driveItem, nil
}

func formatDriveItems(mds []*storageprovider.ResourceInfo) ([]*msgraph.DriveItem, error) {
	responses := make([]*msgraph.DriveItem, 0, len(mds))
	for i := range mds {
		res, err := cs3ResourceToDriveItem(mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func cs3StorageSpaceToDrive(baseUrl *url.URL, space *storageprovider.StorageSpace) (*msgraph.Drive, error) {
	rootId := space.Root.StorageId + "!" + space.Root.OpaqueId
	drive := &msgraph.Drive{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Id: &space.Id.OpaqueId,
			},
			Name: &space.Name,
			//"createdDateTime": "string (timestamp)", // TODO read from StorageSpace ... needs Opaque for now
			//"description": "string", // TODO read from StorageSpace ... needs Opaque for now
		},
		Owner: &msgraph.IdentitySet{
			User: &msgraph.Identity{
				Id: &space.Owner.Id.OpaqueId,
				// DisplayName: , TODO read and cache from users provider
			},
		},

		DriveType: &space.SpaceType,
		Root: &msgraph.DriveItem{
			BaseItem: msgraph.BaseItem{
				Entity: msgraph.Entity{
					Id: &rootId,
				},
			},
		},
	}

	if baseUrl != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		webDavURL := baseUrl.String() + rootId
		drive.Root.WebDavUrl = &webDavURL
	}

	if space.Mtime != nil {
		lastModified := cs3TimestampToTime(space.Mtime)
		drive.BaseItem.LastModifiedDateTime = &lastModified
	}
	if space.Quota != nil {
		// FIXME use https://github.com/owncloud/open-graph-api and return proper int64
		var t int64
		if space.Quota.QuotaMaxBytes > math.MaxInt64 {
			t = math.MaxInt64
		} else {
			t = int64(space.Quota.QuotaMaxBytes)
		}
		drive.Quota = &msgraph.Quota{
			Total: &t,
		}
	}
	// FIXME use coowner from https://github.com/owncloud/open-graph-api

	return drive, nil
}

func formatDrives(baseUrl *url.URL, mds []*storageprovider.StorageSpace) ([]*opengraph.Drive, error) {
	responses := make([]*msgraph.Drive, 0, len(mds))
	for i := range mds {
		res, err := cs3StorageSpaceToDrive(baseUrl, mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

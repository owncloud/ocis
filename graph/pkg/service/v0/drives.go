package svc

import (
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/render"
	"google.golang.org/grpc/metadata"

	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/token"

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
	// TODO refactor this: forward the token to the next service
	t := r.Header.Get("x-access-token")
	if t == "" {
		g.logger.Error().Msg("no access token provided in request")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ctx = token.ContextSetToken(ctx, t)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-access-token", t)

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

// GetRootDriveChildren implements the Service interface.
func (g Graph) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msgf("Calling GetRootDriveChildren")

	ctx := r.Context()

	fn := "/home"

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	g.logger.Info().Msgf("provides access token %v", ctx)

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
	name := strings.TrimPrefix(res.Path, "/home/")

	id := res.Id.StorageId + "!" + res.Id.NodeId

	driveItem := &msgraph.DriveItem{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Id: &id,
			},
			Name: &name,
			ETag: &res.Etag,
		},
		Size: size,
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		driveItem.BaseItem.LastModifiedDateTime = &lastModified
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE {
		driveItem.File = &msgraph.File{
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &msgraph.Folder{}
	}
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
	rootId := space.Root.StorageId + "!" + space.Root.NodeId
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

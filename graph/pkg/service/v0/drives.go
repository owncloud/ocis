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
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

func getToken(r *http.Request) string {
	// 1. check Authorization header
	hdr := r.Header.Get("Authorization")
	t := strings.TrimPrefix(hdr, "Bearer ")
	if t != "" {
		return t
	}
	// TODO 2. check form encoded body parameter for POST requests, see https://tools.ietf.org/html/rfc6750#section-2.2

	// 3. check uri query parameter, see https://tools.ietf.org/html/rfc6750#section-2.3
	tokens, ok := r.URL.Query()["access_token"]
	if !ok || len(tokens[0]) < 1 {
		return ""
	}

	return tokens[0]
}

// GetDrives implements the Service interface.
func (g Graph) GetDrives(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msgf("Calling GetDrives")
	accessToken := getToken(r)
	if accessToken == "" {
		g.logger.Error().Msg("no access token provided in request")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := r.Header.Get("x-access-token")
	ctx = token.ContextSetToken(ctx, t)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-access-token", t)

	req := &storageprovider.ListStorageSpacesRequest{

		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: "1284d238-aa92-42ce-bdc4-0b0000009157", // FIXME dynamically discover home and other storages ... actually the storage registry should provide the list of storage spaces
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
		g.logger.Error().Err(err).Msg("error calling grpc list storage spaces")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error parsing url", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := formatDrives(wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error encoding response as json %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

// GetRootDriveChildren implements the Service interface.
func (g Graph) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msgf("Calling GetRootDriveChildren")
	accessToken := getToken(r)
	if accessToken == "" {
		g.logger.Error().Msg("no access token provided in request")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ctx := r.Context()

	fn := g.config.WebdavNamespace

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t := r.Header.Get("x-access-token")
	ctx = token.ContextSetToken(ctx, t)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-access-token", t)

	g.logger.Info().Msgf("provides access token %v", ctx)

	ref := &storageprovider.Reference{
		Path: fn,
	}

	req := &storageprovider.ListContainerRequest{
		Ref: ref,
	}
	res, err := client.ListContainer(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Msgf("error sending list container grpc request %s", fn)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Msgf("error calling grpc list container %s", fn)
		w.WriteHeader(http.StatusInternalServerError)
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
	size := new(int)
	*size = int(res.Size) // uint64 -> int :boom:
	name := strings.TrimPrefix(res.Path, "/home/")

	driveItem := &msgraph.DriveItem{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Object: msgraph.Object{},
				ID:     &res.Id.OpaqueId,
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
	rootId := space.Root.StorageId + "!" + space.Root.OpaqueId
	drive := &msgraph.Drive{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				ID: &space.Id.OpaqueId,
			},
			Name: &space.Name,
			//"createdDateTime": "string (timestamp)", // TODO read from StorageSpace ... needs Opaque for now
			//"description": "string", // TODO read from StorageSpace ... needs Opaque for now
		},
		Owner: &msgraph.IdentitySet{
			User: &msgraph.Identity{
				ID: &space.Owner.Id.OpaqueId,
				// DisplayName: , TODO read and cache from users provider
			},
		},

		DriveType: &space.SpaceType,
		Root: &msgraph.DriveItem{
			BaseItem: msgraph.BaseItem{
				Entity: msgraph.Entity{
					ID: &rootId,
				},
			},
		},
	}

	if baseUrl != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		webDavURL := baseUrl.String() + rootId
		drive.Root.WebDavURL = &webDavURL
	}

	if space.Mtime != nil {
		lastModified := cs3TimestampToTime(space.Mtime)
		drive.BaseItem.LastModifiedDateTime = &lastModified
	}
	if space.Quota != nil {
		// FIXME use https://github.com/owncloud/open-graph-api and return proper int64
		var t int
		if space.Quota.QuotaMaxBytes > math.MaxInt32 {
			t = math.MaxInt32
		} else {
			t = int(space.Quota.QuotaMaxBytes)
		}
		drive.Quota = &msgraph.Quota{
			Total: &t,
		}
	}
	// FIXME use coowner from https://github.com/owncloud/open-graph-api

	return drive, nil
}

func formatDrives(baseUrl *url.URL, mds []*storageprovider.StorageSpace) ([]*msgraph.Drive, error) {
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

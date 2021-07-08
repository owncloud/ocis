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

	msgraph "github.com/owncloud/open-graph-api-go"
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
	g.logger.Info().Msg("Calling GetDrives")
	if getToken(r) == "" {
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

	res, err := client.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{})
	if err != nil {
		g.logger.Error().Err(err).Msg("error sending list storage spaces grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO handle not found and other status codes
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Interface("status", res.Status).Msg("error calling grpc list storage spaces")
		w.WriteHeader(http.StatusInternalServerError)
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
	g.logger.Info().Msg("Calling GetRootDriveChildren")
	if getToken(r) == "" {
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

	g.logger.Info().Interface("context", ctx).Msg("provides access token")

	ref := &storageprovider.Reference{
		Path: fn,
	}

	req := &storageprovider.ListContainerRequest{
		Ref: ref,
	}
	res, err := client.ListContainer(ctx, req)
	if err != nil {
		g.logger.Error().Err(err).Str("path", fn).Msg("error sending list container grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO handle not found and other status codes
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		g.logger.Error().Err(err).Str("path", fn).Interface("status", res.Status).Msg("error calling grpc list container")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files, err := formatDriveItems(res.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
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

	driveItem := &msgraph.DriveItem{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Id: &res.Id.OpaqueId,
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
		driveItem.File = &msgraph.OpenGraphFile{ // FIXME We cannot use msgraph.File here because the openapi codegenerator autodetects 'File' as a go type ...
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

func cs3StorageSpaceToDrive(baseURL *url.URL, space *storageprovider.StorageSpace) (*msgraph.Drive, error) {
	rootID := space.Root.StorageId + "!" + space.Root.OpaqueId
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
					Id: &rootID,
				},
			},
		},
	}

	if baseURL != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		webDavURL := baseURL.String() + rootID
		drive.Root.WebDavUrl = &webDavURL
	}

	if space.Mtime != nil {
		lastModified := cs3TimestampToTime(space.Mtime)
		drive.BaseItem.LastModifiedDateTime = &lastModified
	}
	if space.Quota != nil {
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

func formatDrives(baseURL *url.URL, mds []*storageprovider.StorageSpace) ([]*msgraph.Drive, error) {
	responses := make([]*msgraph.Drive, 0, len(mds))
	for i := range mds {
		res, err := cs3StorageSpaceToDrive(baseURL, mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

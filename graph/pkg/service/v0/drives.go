package svc

import (
	"github.com/go-chi/render"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
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

	fn := "/home"

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get reva token
	authReq := &gateway.AuthenticateRequest{
		Type:         "bearer",
		ClientSecret: accessToken,
	}

	authRes, _ := client.Authenticate(ctx, authReq);
	ctx = token.ContextSetToken(ctx, authRes.Token)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-access-token", authRes.Token)

	g.logger.Info().Msgf("provides access token %v", ctx)

	ref := &storageprovider.Reference{
		Spec: &storageprovider.Reference_Path{Path: fn},
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


func cs3ResourceToDriveItem(res *storageprovider.ResourceInfo) (*msgraph.DriveItem, error) {
	size := new(int)
	*size = int(res.Size) // uint64 -> int :boom:
	name := strings.TrimPrefix(res.Path, "/home/")
	lastModified := new(time.Time)
	*lastModified = time.Unix(int64(res.Mtime.Seconds), int64(res.Mtime.Nanos))

	driveItem := &msgraph.DriveItem{
		BaseItem: msgraph.BaseItem{
			Entity: msgraph.Entity{
				Object: msgraph.Object{},
				ID:     &res.Id.OpaqueId,
			},
			Name: &name,
			LastModifiedDateTime: lastModified,
			ETag: &res.Etag,
		},
		Size: size,
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE {
		driveItem.File = &msgraph.File{
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &msgraph.Folder{
		}
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

package svc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"time"

	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/utils"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
)

// GetRootDriveChildren implements the Service interface.
func (g Graph) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetRootDriveChildren")
	ctx := r.Context()

	client := g.GetGatewayClient()

	res, err := client.GetHome(ctx, &storageprovider.GetHomeRequest{})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error sending get home grpc request")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
			return
		}
		g.logger.Error().Err(err).Msg("error sending get home grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	lRes, err := client.ListContainer(ctx, &storageprovider.ListContainerRequest{
		Ref: &storageprovider.Reference{
			Path: res.Path,
		},
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error sending list container grpc request")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
			return
		}
		if res.Status.Code == cs3rpc.Code_CODE_PERMISSION_DENIED {
			// TODO check if we should return 404 to not disclose existing items
			errorcode.AccessDenied.Render(w, r, http.StatusForbidden, res.Status.Message)
			return
		}
		g.logger.Error().Err(err).Msg("error sending list container grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	files, err := formatDriveItems(lRes.Infos)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

func (g Graph) getDriveItem(ctx context.Context, root *storageprovider.ResourceId, relativePath string) (*libregraph.DriveItem, error) {

	client := g.GetGatewayClient()

	ref := &storageprovider.Reference{
		ResourceId: root,
		// the path is always relative to the root of the drive, not the location of the .config/ocis/space.yaml file
		Path: utils.MakeRelativePath(relativePath),
	}
	res, err := client.Stat(ctx, &storageprovider.StatRequest{Ref: ref})
	if err != nil {
		return nil, err
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		return nil, fmt.Errorf("could not stat %s: %s", ref, res.Status.Message)
	}

	return cs3ResourceToDriveItem(res.Info)
}

func formatDriveItems(mds []*storageprovider.ResourceInfo) ([]*libregraph.DriveItem, error) {
	responses := make([]*libregraph.DriveItem, 0, len(mds))
	for i := range mds {
		res, err := cs3ResourceToDriveItem(mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func cs3TimestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

func cs3ResourceToDriveItem(res *storageprovider.ResourceInfo) (*libregraph.DriveItem, error) {
	size := new(int64)
	*size = int64(res.Size) // TODO lurking overflow: make size of libregraph drive item use uint64

	driveItem := &libregraph.DriveItem{
		Id:   &res.Id.OpaqueId,
		Size: size,
	}

	if name := path.Base(res.Path); name != "" {
		driveItem.Name = &name
	}
	if res.Etag != "" {
		driveItem.ETag = &res.Etag
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		driveItem.LastModifiedDateTime = &lastModified
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE && res.MimeType != "" {
		// We cannot use a libregraph.File here because the openapi codegenerator autodetects 'File' as a go type ...
		driveItem.File = &libregraph.OpenGraphFile{
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &libregraph.Folder{}
	}
	return driveItem, nil
}

func (g Graph) getDriveMetadata(ctx context.Context, space *storageprovider.StorageSpace) (map[string]string, error) {
	client := g.GetGatewayClient()
	sResp, err := client.Stat(
		ctx,
		&storageprovider.StatRequest{
			Ref: &storageprovider.Reference{
				ResourceId: &storageprovider.ResourceId{
					StorageId: space.Root.StorageId,
					OpaqueId:  space.Root.OpaqueId,
				},
			},
		},
	)
	if err != nil {
		g.logger.Debug().Err(err).Interface("space", space).Msg("transport error")
		return nil, err
	}
	if sResp.Status.Code != cs3rpc.Code_CODE_OK {
		g.logger.Debug().Interface("space", space).Msg("space not found")
		return nil, fmt.Errorf("space not found")
	}
	md := sResp.Info.ArbitraryMetadata.GetMetadata()
	if md == nil {
		g.logger.Error().Err(err).Interface("space", space).Msg("could not read metadata from space")
		return nil, fmt.Errorf("could not read metadata from space")
	}
	return md, nil
}

func (g Graph) setSpaceMetadata(ctx context.Context, metadata map[string]string, resourceID *storageprovider.ResourceId) error {
	if len(metadata) == 0 {
		return nil
	}
	client := g.GetGatewayClient()
	resMd, err := client.SetArbitraryMetadata(
		ctx,
		&storageprovider.SetArbitraryMetadataRequest{
			ArbitraryMetadata: &storageprovider.ArbitraryMetadata{
				Metadata: metadata,
			},
			Ref: &storageprovider.Reference{
				ResourceId: resourceID,
			},
		},
	)
	if err != nil {
		g.logger.Error().Msg("transport error, could not set metadata")
		return err
	}
	if resMd.Status.Code != cs3rpc.Code_CODE_OK {
		g.logger.Error().Msg("could not set metadata")
		return fmt.Errorf("could not set metadata")
	}
	return nil
}

func (g Graph) GetSpecialSpaceItems(ctx context.Context, baseURL *url.URL, spaceRootID *storageprovider.ResourceId, metadata map[string]string) []libregraph.DriveItem {
	var spaceItems []libregraph.DriveItem
	if metadata == nil {
		return spaceItems
	}
	if readmePath, ok := metadata[ReadmePathAttrName]; ok {
		readmeItem, err := g.getDriveItem(ctx, spaceRootID, readmePath)
		if err != nil {
			g.logger.Error().Err(err).Str(ReadmePathAttrName, readmePath).Msg("Could not get readme Item")
		} else {
			readmeItem.SpecialFolder = &libregraph.SpecialFolder{Name: libregraph.PtrString(ReadmePathSpecialFolderName)}
			readmeItem.WebDavUrl = libregraph.PtrString(baseURL.String() + filepath.Join(spaceRootID.OpaqueId, readmePath))
			spaceItems = append(spaceItems, *readmeItem)
		}
	}
	if imagePath, ok := metadata[SpaceImageAttrName]; ok {
		imageItem, err := g.getDriveItem(ctx, spaceRootID, imagePath)
		if err != nil {
			g.logger.Error().Err(err).Str(SpaceImageAttrName, imagePath).Msg("Could not get image Item")
		} else {
			imageItem.SpecialFolder = &libregraph.SpecialFolder{Name: libregraph.PtrString(SpaceImageSpecialFolderName)}
			imageItem.WebDavUrl = libregraph.PtrString(baseURL.String() + filepath.Join(spaceRootID.OpaqueId, imagePath))
			spaceItems = append(spaceItems, *imageItem)
		}
	}
	return spaceItems
}

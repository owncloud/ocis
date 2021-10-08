package svc

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/CiscoM31/godata"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	sproto "github.com/owncloud/ocis/settings/pkg/proto/v0"
	settingsSvc "github.com/owncloud/ocis/settings/pkg/service/v0"
	msgraph "github.com/owncloud/open-graph-api-go"
)

// GetDrives implements the Service interface.
func (g Graph) GetDrives(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetDrives")
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Err(err).Msg("error getting grpc client")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := client.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{
		// TODO add filters?
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error sending list storage spaces grpc request")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			// return an empty list
			render.Status(r, http.StatusOK)
			render.JSON(w, r, &listResponse{})
			return
		}
		g.logger.Error().Err(err).Msg("error sending list storage spaces grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	files, err := formatDrives(wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: files})
}

// GetRootDriveChildren implements the Service interface.
func (g Graph) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	g.logger.Info().Msg("Calling GetRootDriveChildren")
	ctx := r.Context()

	client, err := g.GetClient()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not get client")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

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

// CreateDrive creates a storage drive (space).
func (g Graph) CreateDrive(w http.ResponseWriter, r *http.Request) {
	us, ok := ctxpkg.ContextGetUser(r.Context())
	if !ok {
		errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "invalid user")
		return
	}

	s := sproto.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)

	_, err := s.GetPermissionByID(r.Context(), &sproto.GetPermissionByIDRequest{
		PermissionId: settingsSvc.CreateSpacePermissionID,
	})
	if err != nil {
		// if the permission is not existing for the user in context we can assume we don't have it. Return 401.
		errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "insufficient permissions to create a space.")
		return
	}

	client, err := g.GetClient()
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	drive := msgraph.Drive{}
	if err := json.NewDecoder(r.Body).Decode(&drive); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Errorf("invalid schema definition").Error())
		return
	}
	spaceName := *drive.Name
	if spaceName == "" {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, fmt.Errorf("invalid name").Error())
		return
	}

	var driveType string
	if drive.DriveType != nil {
		driveType = *drive.DriveType
	}
	switch driveType {
	case "":
		driveType = "project"
	case "share":
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Errorf("drives of type share cannot be created via this api").Error())
	}

	var quota uint64
	if drive.Quota != nil && drive.Quota.Total != nil {
		quota = uint64(*drive.Quota.Total)
	} else {
		quota = 65536 // set default quota if no value was sent.
	}

	csr := provider.CreateStorageSpaceRequest{
		Owner: us,
		Type:  driveType,
		Name:  spaceName,
		Quota: &provider.Quota{
			QuotaMaxBytes: quota,
		},
	}

	resp, err := client.CreateStorageSpace(r.Context(), &csr)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if resp.GetStatus().GetCode() != v1beta11.Code_CODE_OK {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, fmt.Errorf("").Error())
	}
}

func (g Graph) UpdateDrive(w http.ResponseWriter, r *http.Request) {
	// wildcards however addressed here is not yet supported. We want to address drives by their unique
	// identifiers. Any open queries need to be implemented. Same applies for sub-entities.
	// For further reading: http://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part2-url-conventions.html#sec_AddressingaSubsetofaCollection

	// strip "/graph/v1.0/" out and parse the rest. This is how godata input is expected.
	//https://github.com/CiscoM31/godata/blob/d70e191d2908191623be84401fecc40d6af4afde/url_parser_test.go#L10
	sanitized := strings.TrimLeft(r.URL.Path, "/graph/v1.0/")

	req, err := godata.ParseRequest(sanitized, r.URL.Query(), true)
	if err != nil {
		panic(err)
	}

	if req.FirstSegment.Identifier.Get() == "" {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Errorf("identifier cannot be empty").Error())
		return
	}

	drive := msgraph.Drive{}
	if err = json.NewDecoder(r.Body).Decode(&drive); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", r.Body).Error())
		return
	}

	identifierParts := strings.Split(req.FirstSegment.Identifier.Get(), "!")
	if len(identifierParts) != 2 {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Errorf("invalid resource id: %v", req.FirstSegment.Identifier.Get()).Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	storageID, opaqueID := identifierParts[0], identifierParts[1]

	client, err := g.GetClient()
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	updateSpaceRequest := &provider.UpdateStorageSpaceRequest{
		// Prepare the object to apply the diff from. The properties on StorageSpace will overwrite
		// the original storage space.
		StorageSpace: &provider.StorageSpace{
			Root: &provider.ResourceId{
				StorageId: storageID,
				OpaqueId:  opaqueID,
			},
			Name: *drive.Name,
		},
	}

	resp, err := client.UpdateStorageSpace(r.Context(), updateSpaceRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if resp.GetStatus().GetCode() != v1beta11.Code_CODE_OK {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, fmt.Errorf("").Error())
	}

	w.WriteHeader(http.StatusNoContent)
}

func cs3TimestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

func cs3ResourceToDriveItem(res *storageprovider.ResourceInfo) (*msgraph.DriveItem, error) {
	size := new(int64)
	*size = int64(res.Size) // uint64 -> int :boom:
	name := path.Base(res.Path)

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

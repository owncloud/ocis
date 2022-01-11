package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/pkg/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	sproto "github.com/owncloud/ocis/settings/pkg/proto/v0"
	settingsSvc "github.com/owncloud/ocis/settings/pkg/service/v0"

	merrors "go-micro.dev/v4/errors"
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

	permissions := make(map[string]struct{}, 1)
	s := sproto.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)

	_, err = s.GetPermissionByID(ctx, &sproto.GetPermissionByIDRequest{
		PermissionId: settingsSvc.ListAllSpacesPermissionID,
	})

	// No error means the user has the permission
	if err == nil {
		permissions[settingsSvc.ListAllSpacesPermissionName] = struct{}{}
	}
	value, err := json.Marshal(permissions)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := client.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{
		Opaque: &types.Opaque{Map: map[string]*types.OpaqueEntry{
			"permissions": {
				Decoder: "json",
				Value:   value,
			},
		}},
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
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	files, err := g.formatDrives(ctx, wdu, res.StorageSpaces)
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
	drive := libregraph.Drive{}
	if err := json.NewDecoder(r.Body).Decode(&drive); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, "invalid schema definition")
		return
	}
	spaceName := *drive.Name
	if spaceName == "" {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "invalid name")
		return
	}

	var driveType string
	if drive.DriveType != nil {
		driveType = *drive.DriveType
	}
	switch driveType {
	case "", "project":
		driveType = "project"
	default:
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Sprintf("drives of type %s cannot be created via this api", driveType))
		return
	}

	csr := storageprovider.CreateStorageSpaceRequest{
		Owner: us,
		Type:  driveType,
		Name:  spaceName,
		Quota: getQuota(drive.Quota, g.config.Spaces.DefaultQuota),
	}

	resp, err := client.CreateStorageSpace(r.Context(), &csr)
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if resp.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "")
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	newDrive, err := cs3StorageSpaceToDrive(wdu, resp.StorageSpace)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing space")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newDrive)
}

func (g Graph) UpdateDrive(w http.ResponseWriter, r *http.Request) {
	driveID, err := url.PathUnescape(chi.URLParam(r, "driveID"))
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping drive id failed")
		return
	}

	if driveID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing drive id")
		return
	}

	drive := libregraph.Drive{}
	if err = json.NewDecoder(r.Body).Decode(&drive); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", r.Body))
		return
	}

	root := &storageprovider.ResourceId{}

	identifierParts := strings.Split(driveID, "!")
	switch len(identifierParts) {
	case 1:
		root.StorageId, root.OpaqueId = identifierParts[0], identifierParts[0]
	case 2:
		root.StorageId, root.OpaqueId = identifierParts[0], identifierParts[1]
	default:
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid resource id: %v", driveID))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client, err := g.GetClient()
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	updateSpaceRequest := &storageprovider.UpdateStorageSpaceRequest{
		// Prepare the object to apply the diff from. The properties on StorageSpace will overwrite
		// the original storage space.
		StorageSpace: &storageprovider.StorageSpace{
			Id: &storageprovider.StorageSpaceId{
				OpaqueId: root.StorageId + "!" + root.OpaqueId,
			},
			Root: root,
		},
	}

	if drive.Name != nil {
		updateSpaceRequest.StorageSpace.Name = *drive.Name
	}

	if drive.Quota.HasTotal() {
		user := ctxpkg.ContextMustGetUser(r.Context())
		canSetSpaceQuota, err := canSetSpaceQuota(r.Context(), user)
		if err != nil {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		if !canSetSpaceQuota {
			errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "user is not allowed to set the space quota")
			return
		}
		updateSpaceRequest.StorageSpace.Quota = &storageprovider.Quota{
			QuotaMaxBytes: uint64(*drive.Quota.Total),
		}
	}

	resp, err := client.UpdateStorageSpace(r.Context(), updateSpaceRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.GetStatus().GetCode() != cs3rpc.Code_CODE_OK {
		switch resp.Status.GetCode() {
		case cs3rpc.Code_CODE_NOT_FOUND:
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, resp.GetStatus().GetMessage())
			return
		default:
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, resp.GetStatus().GetMessage())
			return
		}
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	updatedDrive, err := cs3StorageSpaceToDrive(wdu, resp.StorageSpace)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing space")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, updatedDrive)
}

func cs3TimestampToTime(t *types.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), int64(t.Nanos))
}

func cs3ResourceToDriveItem(res *storageprovider.ResourceInfo) (*libregraph.DriveItem, error) {
	size := new(int64)
	*size = int64(res.Size) // uint64 -> int :boom:
	name := path.Base(res.Path)

	driveItem := &libregraph.DriveItem{
		Id:   &res.Id.OpaqueId,
		Name: &name,
		ETag: &res.Etag,
		Size: size,
	}
	if res.Mtime != nil {
		lastModified := cs3TimestampToTime(res.Mtime)
		driveItem.LastModifiedDateTime = &lastModified
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_FILE {
		driveItem.File = &libregraph.OpenGraphFile{ // FIXME We cannot use libregraph.File here because the openapi codegenerator autodetects 'File' as a go type ...
			MimeType: &res.MimeType,
		}
	}
	if res.Type == storageprovider.ResourceType_RESOURCE_TYPE_CONTAINER {
		driveItem.Folder = &libregraph.Folder{}
	}
	return driveItem, nil
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

func cs3StorageSpaceToDrive(baseURL *url.URL, space *storageprovider.StorageSpace) (*libregraph.Drive, error) {
	rootID := space.Root.StorageId + "!" + space.Root.OpaqueId
	if space.Root.StorageId == space.Root.OpaqueId {
		// omit opaqueid
		rootID = space.Root.StorageId
	}

	drive := &libregraph.Drive{
		Id:   &rootID,
		Name: &space.Name,
		//"createdDateTime": "string (timestamp)", // TODO read from StorageSpace ... needs Opaque for now
		//"description": "string", // TODO read from StorageSpace ... needs Opaque for now
		DriveType: &space.SpaceType,
		Root: &libregraph.DriveItem{
			Id: &rootID,
		},
	}

	if baseURL != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		webDavURL := baseURL.String() + rootID
		drive.Root.WebDavUrl = &webDavURL
	}

	// TODO The public space has no owner ... should we even show it?
	if space.Owner != nil && space.Owner.Id != nil {
		drive.Owner = &libregraph.IdentitySet{
			User: &libregraph.Identity{
				Id: &space.Owner.Id.OpaqueId,
				// DisplayName: , TODO read and cache from users provider
			},
		}
	}
	if space.Mtime != nil {
		lastModified := cs3TimestampToTime(space.Mtime)
		drive.LastModifiedDateTime = &lastModified
	}
	if space.Quota != nil {
		var t int64
		if space.Quota.QuotaMaxBytes > math.MaxInt64 {
			t = math.MaxInt64
		} else {
			t = int64(space.Quota.QuotaMaxBytes)
		}
		drive.Quota = &libregraph.Quota{
			Total: &t,
		}
	}
	// FIXME use coowner from https://github.com/owncloud/open-graph-api

	return drive, nil
}

func (g Graph) formatDrives(ctx context.Context, baseURL *url.URL, mds []*storageprovider.StorageSpace) ([]*libregraph.Drive, error) {
	responses := make([]*libregraph.Drive, 0, len(mds))
	for i := range mds {
		res, err := cs3StorageSpaceToDrive(baseURL, mds[i])
		if err != nil {
			return nil, err
		}
		res.Quota, err = g.getDriveQuota(ctx, mds[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

func (g Graph) getDriveQuota(ctx context.Context, space *storageprovider.StorageSpace) (*libregraph.Quota, error) {
	client, err := g.GetClient()
	if err != nil {
		g.logger.Error().Err(err).Msg("error creating grpc client")
		return nil, err
	}

	req := &gateway.GetQuotaRequest{
		Ref: &storageprovider.Reference{
			ResourceId: &storageprovider.ResourceId{
				StorageId: space.Root.StorageId,
				OpaqueId:  space.Root.OpaqueId,
			},
			Path: ".",
		},
	}
	res, err := client.GetQuota(ctx, req)
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg("error sending get quota grpc request")
		return nil, nil
	case res.Status.Code == cs3rpc.Code_CODE_UNIMPLEMENTED:
		// TODO well duh
		return nil, nil
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		g.logger.Error().Err(err).Msg("error sending sending get quota grpc request")
		return nil, err
	}

	total := int64(res.TotalBytes)

	used := int64(res.UsedBytes)
	remaining := total - used
	qta := libregraph.Quota{
		Remaining: &remaining,
		Total:     &total,
		Used:      &used,
	}
	state := calculateQuotaState(total, used)
	qta.State = &state

	return &qta, nil
}

func calculateQuotaState(total int64, used int64) (state string) {
	percent := (float64(used) / float64(total)) * 100

	switch {
	case percent <= float64(75):
		return "normal"
	case percent <= float64(90):
		return "nearing"
	case percent <= float64(99):
		return "critical"
	default:
		return "exceeded"
	}
}

func getQuota(quota *libregraph.Quota, defaultQuota string) *storageprovider.Quota {
	switch {
	case quota != nil && quota.Total != nil:
		if q := *quota.Total; q >= 0 {
			return &storageprovider.Quota{QuotaMaxBytes: uint64(q)}
		}
		fallthrough
	case defaultQuota != "":
		if q, err := strconv.ParseInt(defaultQuota, 10, 64); err == nil && q >= 0 {
			return &storageprovider.Quota{QuotaMaxBytes: uint64(q)}
		}
		fallthrough
	default:
		return nil
	}
}

func canSetSpaceQuota(ctx context.Context, user *userv1beta1.User) (bool, error) {
	settingsService := sproto.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)
	_, err := settingsService.GetPermissionByID(ctx, &sproto.GetPermissionByIDRequest{PermissionId: settingsSvc.SetSpaceQuotaPermissionID})
	if err != nil {
		merror := merrors.FromError(err)
		if merror.Status == http.StatusText(http.StatusNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

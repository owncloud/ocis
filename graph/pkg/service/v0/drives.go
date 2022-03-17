package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/CiscoM31/godata"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/cs3org/reva/v2/pkg/utils/resourceid"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	v0 "github.com/owncloud/ocis/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/settings/v0"
	settingsServiceExt "github.com/owncloud/ocis/settings/pkg/service/v0"
	merrors "go-micro.dev/v4/errors"
)

// GetDrives implements the Service interface.
func (g Graph) GetDrives(w http.ResponseWriter, r *http.Request) {
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	// Parse the request with odata parser
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		g.logger.Err(err).Interface("query", r.URL.Query()).Msg("query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	g.logger.Info().Interface("query", r.URL.Query()).Msg("Calling GetDrives")
	ctx := r.Context()

	filters, err := generateCs3Filters(odataReq)
	if err != nil {
		g.logger.Err(err).Interface("query", r.URL.Query()).Msg("query error")
		errorcode.NotSupported.Render(w, r, http.StatusNotImplemented, err.Error())
		return
	}
	res, err := g.ListStorageSpacesWithFilters(ctx, filters)
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg(ListStorageSpacesTransportErr)
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			// return an empty list
			render.Status(r, http.StatusOK)
			render.JSON(w, r, &listResponse{})
			return
		}
		g.logger.Error().Err(err).Msg(ListStorageSpacesReturnsErr)
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	spaces, err := g.formatDrives(ctx, wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response as json")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	spaces, err = sortSpaces(odataReq, spaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error sorting the spaces list")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: spaces})
}

// GetSingleDrive does a lookup of a single space by spaceId
func (g Graph) GetSingleDrive(w http.ResponseWriter, r *http.Request) {
	driveID := chi.URLParam(r, "driveID")
	if driveID == "" {
		err := fmt.Errorf("no valid space id retrieved")
		g.logger.Err(err)
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	g.logger.Info().Str("driveID", driveID).Msg("Calling GetSingleDrive")
	ctx := r.Context()

	filters := []*storageprovider.ListStorageSpacesRequest_Filter{listStorageSpacesIDFilter(driveID)}
	res, err := g.ListStorageSpacesWithFilters(ctx, filters)
	switch {
	case err != nil:
		g.logger.Error().Err(err).Msg(ListStorageSpacesTransportErr)
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			// the client is doing a lookup for a specific space, therefore we need to return
			// not found to the caller
			g.logger.Error().Str("driveID", driveID).Msg(fmt.Sprintf(NoSpaceFoundMessage, driveID))
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, fmt.Sprintf(NoSpaceFoundMessage, driveID))
			return
		}
		g.logger.Error().Err(err).Msg(ListStorageSpacesReturnsErr)
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing url")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	spaces, err := g.formatDrives(ctx, wdu, res.StorageSpaces)
	if err != nil {
		g.logger.Error().Err(err).Msg("error encoding response")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	switch num := len(spaces); {
	case num == 0:
		g.logger.Error().Str("driveID", driveID).Msg("no space found")
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, fmt.Sprintf(NoSpaceFoundMessage, driveID))
		return
	case num == 1:
		render.Status(r, http.StatusOK)
		render.JSON(w, r, spaces[0])
	default:
		g.logger.Error().Int("number", num).Msg("expected to find a single space but found more")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "expected to find a single space but found more")
		return
	}
}

func canCreateSpace(ctx context.Context, ownPersonalHome bool) bool {
	s := settingssvc.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)

	pr, err := s.GetPermissionByID(ctx, &settingssvc.GetPermissionByIDRequest{
		PermissionId: settingsServiceExt.CreateSpacePermissionID,
	})
	if err != nil || pr.Permission == nil {
		return false
	}
	// TODO @C0rby shouldn't the permissions service check this? aka shouldn't we call CheckPermission?
	if pr.Permission.Constraint == v0.Permission_CONSTRAINT_OWN && !ownPersonalHome {
		return false
	}
	return true
}

// CreateDrive creates a storage drive (space).
func (g Graph) CreateDrive(w http.ResponseWriter, r *http.Request) {
	us, ok := ctxpkg.ContextGetUser(r.Context())
	if !ok {
		errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "invalid user")
		return
	}

	// TODO determine if the user tries to create his own personal space and pass that as a boolean
	if !canCreateSpace(r.Context(), false) {
		// if the permission is not existing for the user in context we can assume we don't have it. Return 401.
		errorcode.GeneralException.Render(w, r, http.StatusUnauthorized, "insufficient permissions to create a space.")
		return
	}

	client := g.GetGatewayClient()
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

	if drive.Description != nil {
		csr.Opaque = utils.AppendPlainToOpaque(csr.Opaque, "description", *drive.Description)
	}

	if drive.DriveAlias != nil {
		csr.Opaque = utils.AppendPlainToOpaque(csr.Opaque, "spaceAlias", *drive.DriveAlias)
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

	newDrive, err := g.cs3StorageSpaceToDrive(wdu, resp.StorageSpace)
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

	client := g.GetGatewayClient()

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

	// Note: this is the Opaque prop of the request
	if restore, _ := strconv.ParseBool(r.Header.Get("restore")); restore {
		updateSpaceRequest.Opaque = &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"restore": {
					Decoder: "plain",
					Value:   []byte("true"),
				},
			},
		}
	}

	if drive.Description != nil {
		updateSpaceRequest.StorageSpace.Opaque = utils.AppendPlainToOpaque(updateSpaceRequest.StorageSpace.Opaque, "description", *drive.Description)
	}

	if drive.DriveAlias != nil {
		updateSpaceRequest.StorageSpace.Opaque = utils.AppendPlainToOpaque(updateSpaceRequest.StorageSpace.Opaque, "spaceAlias", *drive.DriveAlias)
	}

	for _, special := range drive.Special {
		if special.Id != nil {
			updateSpaceRequest.StorageSpace.Opaque = utils.AppendPlainToOpaque(updateSpaceRequest.StorageSpace.Opaque, *special.SpecialFolder.Name, *special.Id)
		}
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
		case cs3rpc.Code_CODE_PERMISSION_DENIED:
			errorcode.NotAllowed.Render(w, r, http.StatusForbidden, resp.GetStatus().GetMessage())
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

	spaces, err := g.formatDrives(r.Context(), wdu, []*storageprovider.StorageSpace{resp.StorageSpace})
	if err != nil {
		g.logger.Error().Err(err).Msg("error parsing space")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, spaces[0])
}

func (g Graph) formatDrives(ctx context.Context, baseURL *url.URL, storageSpaces []*storageprovider.StorageSpace) ([]*libregraph.Drive, error) {
	responses := make([]*libregraph.Drive, 0, len(storageSpaces))
	for _, storageSpace := range storageSpaces {
		res, err := g.cs3StorageSpaceToDrive(baseURL, storageSpace)
		if err != nil {
			return nil, err
		}
		res.Special = g.GetExtendedSpaceProperties(ctx, baseURL, storageSpace)
		res.Quota, err = g.getDriveQuota(ctx, storageSpace)
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	return responses, nil
}

// ListStorageSpacesWithFilters List Storage Spaces using filters
func (g Graph) ListStorageSpacesWithFilters(ctx context.Context, filters []*storageprovider.ListStorageSpacesRequest_Filter) (*storageprovider.ListStorageSpacesResponse, error) {
	client := g.GetGatewayClient()

	permissions := make(map[string]struct{}, 1)
	s := settingssvc.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)

	_, err := s.GetPermissionByID(ctx, &settingssvc.GetPermissionByIDRequest{
		PermissionId: settingsServiceExt.ListAllSpacesPermissionID,
	})

	// No error means the user has the permission
	if err == nil {
		permissions[settingsServiceExt.ListAllSpacesPermissionName] = struct{}{}
	}
	value, err := json.Marshal(permissions)
	if err != nil {
		return nil, err
	}

	res, err := client.ListStorageSpaces(ctx, &storageprovider.ListStorageSpacesRequest{
		Opaque: &types.Opaque{Map: map[string]*types.OpaqueEntry{
			"permissions": {
				Decoder: "json",
				Value:   value,
			},
		}},
		Filters: filters,
	})
	return res, err
}

func (g Graph) cs3StorageSpaceToDrive(baseURL *url.URL, space *storageprovider.StorageSpace) (*libregraph.Drive, error) {
	if space.Root == nil {
		return nil, fmt.Errorf("space has no root")
	}
	rootID := resourceid.OwnCloudResourceIDWrap(space.Root)

	var permissions []libregraph.Permission
	if space.Opaque != nil {
		var m map[string]*storageprovider.ResourcePermissions
		entry, ok := space.Opaque.Map["grants"]
		if ok {
			err := json.Unmarshal(entry.Value, &m)
			if err != nil {
				g.logger.Error().
					Err(err).
					Str("space", space.Root.OpaqueId).
					Msg("failed to read spaces grants")
			}
		}
		if len(m) != 0 {
			managerIdentities := []libregraph.IdentitySet{}
			editorIdentities := []libregraph.IdentitySet{}
			viewerIdentities := []libregraph.IdentitySet{}

			for id, perm := range m {
				// This temporary variable is necessary since we need to pass a pointer to the
				// libregraph.Identity and if we pass the pointer from the loop every identity
				// will have the same id.
				tmp := id
				identity := libregraph.IdentitySet{User: &libregraph.Identity{Id: &tmp}}
				switch {
				case perm.AddGrant:
					managerIdentities = append(managerIdentities, identity)
				case perm.InitiateFileUpload:
					editorIdentities = append(editorIdentities, identity)
				case perm.Stat:
					viewerIdentities = append(viewerIdentities, identity)
				}
			}

			permissions = make([]libregraph.Permission, 0, 3)
			if len(managerIdentities) != 0 {
				permissions = append(permissions, libregraph.Permission{
					GrantedTo: managerIdentities,
					Roles:     []string{"manager"},
				})
			}
			if len(editorIdentities) != 0 {
				permissions = append(permissions, libregraph.Permission{
					GrantedTo: editorIdentities,
					Roles:     []string{"editor"},
				})
			}
			if len(viewerIdentities) != 0 {
				permissions = append(permissions, libregraph.Permission{
					GrantedTo: viewerIdentities,
					Roles:     []string{"viewer"},
				})
			}
		}
	}

	drive := &libregraph.Drive{
		Id:   &space.Root.StorageId,
		Name: &space.Name,
		//"createdDateTime": "string (timestamp)", // TODO read from StorageSpace ... needs Opaque for now
		//"description": "string", // TODO read from StorageSpace ... needs Opaque for now
		DriveType: &space.SpaceType,
		Root: &libregraph.DriveItem{
			Id:          &rootID,
			Permissions: permissions,
		},
	}
	if space.Opaque != nil {
		if description, ok := space.Opaque.Map["description"]; ok {
			drive.Description = libregraph.PtrString(string(description.Value))
		}

		if alias, ok := space.Opaque.Map["spaceAlias"]; ok {
			drive.DriveAlias = libregraph.PtrString(string(alias.Value))
		}

		if v, ok := space.Opaque.Map["trashed"]; ok {
			deleted := &libregraph.Deleted{}
			deleted.SetState(string(v.Value))
			drive.Root.Deleted = deleted
		}

		if entry, ok := space.Opaque.Map["etag"]; ok {
			drive.Root.ETag = libregraph.PtrString(string(entry.Value))
		}
	}

	if baseURL != nil {
		// TODO read from StorageSpace ... needs Opaque for now
		// TODO how do we build the url?
		// for now: read from request
		webDavURL := baseURL.String() + space.Root.StorageId
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

func (g Graph) getDriveQuota(ctx context.Context, space *storageprovider.StorageSpace) (*libregraph.Quota, error) {
	client := g.GetGatewayClient()

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
		g.logger.Error().Err(err).Msg("could not call GetQuota")
		return nil, nil
	case res.Status.Code == cs3rpc.Code_CODE_UNIMPLEMENTED:
		// TODO well duh
		return nil, nil
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		g.logger.Error().Err(err).Msg("error sending get quota grpc request")
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
	settingsService := settingssvc.NewPermissionService("com.owncloud.api.settings", grpc.DefaultClient)
	_, err := settingsService.GetPermissionByID(ctx, &settingssvc.GetPermissionByIDRequest{PermissionId: settingsServiceExt.SetSpaceQuotaPermissionID})
	if err != nil {
		merror := merrors.FromError(err)
		if merror.Status == http.StatusText(http.StatusNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func generateCs3Filters(request *godata.GoDataRequest) ([]*storageprovider.ListStorageSpacesRequest_Filter, error) {
	var filters []*storageprovider.ListStorageSpacesRequest_Filter
	if request.Query.Filter != nil {
		if request.Query.Filter.Tree.Token.Value == "eq" {
			switch request.Query.Filter.Tree.Children[0].Token.Value {
			case "driveType":
				filters = append(filters, listStorageSpacesTypeFilter(strings.Trim(request.Query.Filter.Tree.Children[1].Token.Value, "'")))
			case "id":
				filters = append(filters, listStorageSpacesIDFilter(strings.Trim(request.Query.Filter.Tree.Children[1].Token.Value, "'")))
			}
		} else {
			err := fmt.Errorf("unsupported filter operand: %s", request.Query.Filter.Tree.Token.Value)
			return nil, err
		}
	}
	return filters, nil
}

func listStorageSpacesIDFilter(id string) *storageprovider.ListStorageSpacesRequest_Filter {
	return &storageprovider.ListStorageSpacesRequest_Filter{
		Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
		Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
			Id: &storageprovider.StorageSpaceId{
				OpaqueId: id,
			},
		},
	}
}

func listStorageSpacesTypeFilter(spaceType string) *storageprovider.ListStorageSpacesRequest_Filter {
	return &storageprovider.ListStorageSpacesRequest_Filter{
		Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
		Term: &storageprovider.ListStorageSpacesRequest_Filter_SpaceType{
			SpaceType: spaceType,
		},
	}
}

func (g Graph) DeleteDrive(w http.ResponseWriter, r *http.Request) {
	driveID, err := url.PathUnescape(chi.URLParam(r, "driveID"))
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping drive id failed")
		return
	}

	if driveID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing drive id")
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

	purge := parsePurgeHeader(r.Header)

	var opaque *types.Opaque
	if purge {
		opaque = &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"purge": {},
			},
		}
	}

	dRes, err := g.gatewayClient.DeleteStorageSpace(r.Context(), &storageprovider.DeleteStorageSpaceRequest{
		Opaque: opaque,
		Id: &storageprovider.StorageSpaceId{
			OpaqueId: root.StorageId,
		},
	})
	switch {
	case dRes.Status.Code == cs3rpc.Code_CODE_INVALID_ARGUMENT:
		errorcode.GeneralException.Render(w, r, http.StatusBadRequest, dRes.Status.Message)
		w.WriteHeader(http.StatusBadRequest)
		return
	case err != nil || dRes.Status.Code != cs3rpc.Code_CODE_OK:
		g.logger.Error().Err(err).Msg("error deleting storage space")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func sortSpaces(req *godata.GoDataRequest, spaces []*libregraph.Drive) ([]*libregraph.Drive, error) {
	var sorter sort.Interface
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return spaces, nil
	}
	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case "name":
		sorter = spacesByName{spaces}
	case "lastModifiedDateTime":
		sorter = spacesByLastModifiedDateTime{spaces}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == "asc" {
		sorter = sort.Reverse(sorter)
	}
	sort.Sort(sorter)
	return spaces, nil
}

package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"

	"google.golang.org/grpc/metadata"
)

// SpaceRole defines the user role on space
type SpaceRole func(*storageprovider.ResourcePermissions) bool

// Possible roles in spaces
var (
	AllRole     SpaceRole = func(perms *storageprovider.ResourcePermissions) bool { return true }
	ViewerRole  SpaceRole = func(perms *storageprovider.ResourcePermissions) bool { return perms.Stat }
	EditorRole  SpaceRole = func(perms *storageprovider.ResourcePermissions) bool { return perms.InitiateFileUpload }
	ManagerRole SpaceRole = func(perms *storageprovider.ResourcePermissions) bool { return perms.DenyGrant }
)

var _errStatusCodeTmpl = "unexpected status code while %s: %v"

// Package error checkers
var (
	IsErrNotFound         = func(err error) bool { return IsStatusCodeError(err, rpc.Code_CODE_NOT_FOUND) }
	IsErrPermissionDenied = func(err error) bool { return IsStatusCodeError(err, rpc.Code_CODE_PERMISSION_DENIED) }
)

// GetServiceUserContext returns an authenticated context of the given service user
func GetServiceUserContext(serviceUserID string, gwc gateway.GatewayAPIClient, serviceUserSecret string) (context.Context, error) {
	ctx := context.Background()
	authRes, err := gwc.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "serviceaccounts",
		ClientId:     serviceUserID,
		ClientSecret: serviceUserSecret,
	})
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode("authenticating service user", authRes.GetStatus().GetCode()); err != nil {
		return nil, err
	}

	return metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, authRes.Token), nil
}

// GetUser gets the specified user
func GetUser(userID *user.UserId, gwc gateway.GatewayAPIClient) (*user.User, error) {
	getUserResponse, err := gwc.GetUser(context.Background(), &user.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode("getting user", getUserResponse.GetStatus().GetCode()); err != nil {
		return nil, err
	}

	return getUserResponse.GetUser(), nil
}

// GetSpace returns the given space
func GetSpace(ctx context.Context, spaceID string, gwc gateway.GatewayAPIClient) (*storageprovider.StorageSpace, error) {
	res, err := gwc.ListStorageSpaces(ctx, listStorageSpaceRequest(spaceID))
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode("getting space", res.GetStatus().GetCode()); err != nil {
		return nil, err
	}

	if len(res.StorageSpaces) == 0 {
		return nil, statusCodeError{"getting space", rpc.Code_CODE_NOT_FOUND}
	}

	return res.StorageSpaces[0], nil
}

// GetGroupMembers returns all members of the given group
func GetGroupMembers(ctx context.Context, groupID string, gwc gateway.GatewayAPIClient) ([]string, error) {
	r, err := gwc.GetGroup(ctx, &group.GetGroupRequest{GroupId: &group.GroupId{OpaqueId: groupID}})
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode("getting group", r.GetStatus().GetCode()); err != nil {
		return nil, err
	}

	users := make([]string, 0, len(r.GetGroup().GetMembers()))
	for _, u := range r.GetGroup().GetMembers() {
		users = append(users, u.GetOpaqueId())
	}

	return users, nil
}

// ResolveID returns either the given userID or all members of the given groupID (if userID is nil)
func ResolveID(ctx context.Context, userid *user.UserId, groupid *group.GroupId, gwc gateway.GatewayAPIClient) ([]string, error) {
	if userid != nil {
		return []string{userid.GetOpaqueId()}, nil
	}

	if ctx == nil {
		return nil, errors.New("need ctx to resolve group id")
	}

	return GetGroupMembers(ctx, groupid.GetOpaqueId(), gwc)
}

// GetSpaceMembers returns all members of the given space that have at least the given role. `nil` role will be interpreted as all
func GetSpaceMembers(ctx context.Context, spaceID string, gwc gateway.GatewayAPIClient, role SpaceRole) ([]string, error) {
	if ctx == nil {
		return nil, errors.New("need authenticated context to find space members")
	}

	space, err := GetSpace(ctx, spaceID, gwc)
	if err != nil {
		return nil, err
	}

	var users []string
	switch space.SpaceType {
	case "personal":
		users = append(users, space.GetOwner().GetId().GetOpaqueId())
	case "project":
		if users, err = gatherProjectSpaceMembers(ctx, space, gwc, role); err != nil {
			return nil, err
		}
	default:
		// TODO: shares? other space types?
		return nil, fmt.Errorf("unsupported space type: %s", space.SpaceType)
	}

	return users, nil
}

// GetResourceByID is a convenience method to get a resource by its resourceID
func GetResourceByID(ctx context.Context, resourceid *storageprovider.ResourceId, gwc gateway.GatewayAPIClient) (*storageprovider.ResourceInfo, error) {
	return GetResource(ctx, &storageprovider.Reference{ResourceId: resourceid}, gwc)
}

// GetResource returns a resource by reference
func GetResource(ctx context.Context, ref *storageprovider.Reference, gwc gateway.GatewayAPIClient) (*storageprovider.ResourceInfo, error) {
	res, err := gwc.Stat(ctx, &storageprovider.StatRequest{Ref: ref})
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode("getting resource", res.GetStatus().GetCode()); err != nil {
		return nil, err
	}

	return res.GetInfo(), nil
}

// IsStatusCodeError returns true if `err` was caused because of status code `code`
func IsStatusCodeError(err error, code rpc.Code) bool {
	sce, ok := err.(statusCodeError)
	if !ok {
		return false
	}

	return sce.code == code
}

func checkStatusCode(reason string, code rpc.Code) error {
	if code == rpc.Code_CODE_OK {
		return nil
	}
	return statusCodeError{reason, code}
}

func gatherProjectSpaceMembers(ctx context.Context, space *storageprovider.StorageSpace, gwc gateway.GatewayAPIClient, role SpaceRole) ([]string, error) {
	var permissionsMap map[string]*storageprovider.ResourcePermissions
	if err := ReadJSONFromOpaque(space.GetOpaque(), "grants", &permissionsMap); err != nil {
		return nil, err
	}

	groupsMap := make(map[string]struct{})
	if opaqueGroups, ok := space.Opaque.Map["groups"]; ok {
		_ = json.Unmarshal(opaqueGroups.GetValue(), &groupsMap)
	}

	if role == nil {
		role = AllRole
	}

	// we use a map to avoid duplicates
	usermap := make(map[string]struct{})
	for id, perm := range permissionsMap {
		if !role(perm) {
			continue
		}

		if _, isGroup := groupsMap[id]; !isGroup {
			usermap[id] = struct{}{}
			continue
		}

		usrs, err := GetGroupMembers(ctx, id, gwc)
		if err != nil {
			// TODO: continue?
			return nil, err
		}

		for _, u := range usrs {
			usermap[u] = struct{}{}
		}
	}

	users := make([]string, 0, len(usermap))
	for id := range usermap {
		users = append(users, id)
	}

	return users, nil
}

func listStorageSpaceRequest(spaceID string) *storageprovider.ListStorageSpacesRequest {
	return &storageprovider.ListStorageSpacesRequest{
		Opaque: AppendPlainToOpaque(nil, "unrestricted", "true"),
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: spaceID,
					},
				},
			},
		},
	}
}

// statusCodeError is a helper struct to return errors
type statusCodeError struct {
	reason string
	code   rpc.Code
}

// Error implements error interface
func (sce statusCodeError) Error() string {
	return fmt.Sprintf(_errStatusCodeTmpl, sce.reason, sce.code)
}

package backend

import (
	"context"
	"encoding/json"
	"errors"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	"google.golang.org/grpc"
)

var (
	// ErrAccountNotFound account not found
	ErrAccountNotFound = errors.New("user not found")
	// ErrAccountDisabled account disabled
	ErrAccountDisabled = errors.New("account disabled")
	// ErrNotSupported operation not supported by user-backend
	ErrNotSupported = errors.New("operation not supported")
)

// UserBackend allows the proxy to retrieve users from different user-backends (accounts-service, CS3)
type UserBackend interface {
	GetUserByClaims(ctx context.Context, claim, value string, withRoles bool) (*cs3.User, error)
	Authenticate(ctx context.Context, username string, password string) (*cs3.User, error)
	CreateUserFromClaims(ctx context.Context, claims *oidc.StandardClaims) (*cs3.User, error)
	GetUserGroups(ctx context.Context, userID string)
}

// RevaAuthenticator helper interface to mock auth-method from reva gateway-client.
type RevaAuthenticator interface {
	Authenticate(ctx context.Context, in *gateway.AuthenticateRequest, opts ...grpc.CallOption) (*gateway.AuthenticateResponse, error)
}

// loadRolesIDs returns the role-ids assigned to an user
func loadRolesIDs(ctx context.Context, opaqueUserID string, rs settings.RoleService) ([]string, error) {
	req := &settings.ListRoleAssignmentsRequest{AccountUuid: opaqueUserID}
	assignmentResponse, err := rs.ListRoleAssignments(ctx, req)

	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0)

	for _, assignment := range assignmentResponse.Assignments {
		roleIDs = append(roleIDs, assignment.RoleId)
	}

	return roleIDs, nil
}

// encodeRoleIDs encoded the given role id's in to reva-specific format to be able to mint a token from them
func encodeRoleIDs(roleIDs []string) (*types.OpaqueEntry, error) {
	roleIDsJSON, err := json.Marshal(roleIDs)
	if err != nil {
		return nil, err
	}

	return &types.OpaqueEntry{
		Decoder: "json",
		Value:   roleIDsJSON,
	}, nil
}

package backend

import (
	"context"
	"fmt"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	settings "github.com/owncloud/ocis/settings/pkg/proto/v0"
	settingsSvc "github.com/owncloud/ocis/settings/pkg/service/v0"
)

type cs3backend struct {
	userProvider        cs3.UserAPIClient
	settingsRoleService settings.RoleService
	authProvider        RevaAuthenticator
	logger              log.Logger
}

// NewCS3UserBackend creates a user-provider which fetches users from a CS3 UserBackend
func NewCS3UserBackend(up cs3.UserAPIClient, rs settings.RoleService, ap RevaAuthenticator, logger log.Logger) UserBackend {
	return &cs3backend{
		userProvider:        up,
		settingsRoleService: rs,
		authProvider:        ap,
		logger:              logger,
	}
}

func (c *cs3backend) GetUserByClaims(ctx context.Context, claim, value string, withRoles bool) (*cs3.User, error) {
	res, err := c.userProvider.GetUserByClaim(ctx, &cs3.GetUserByClaimRequest{
		Claim: claim,
		Value: value,
	})

	switch {
	case err != nil:
		return nil, fmt.Errorf("could not get user by claim %v with value %v: %w", claim, value, err)
	case res.Status.Code != rpcv1beta1.Code_CODE_OK:
		if res.Status.Code == rpcv1beta1.Code_CODE_NOT_FOUND {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("could not get user by claim %v with value %v : %w ", claim, value, err)
	}

	user := res.User

	if !withRoles {
		return user, nil
	}

	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, c.settingsRoleService)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not load roles")
	}

	if len(roleIDs) == 0 {
		roleIDs = append(roleIDs, settingsSvc.BundleUUIDRoleUser, settingsSvc.SelfManagementPermissionID)
		// if roles are empty, assume we haven't seen the user before and assign a default user role. At least until
		// proper roles are provided. See https://github.com/owncloud/ocis/issues/1825 for more context.
		//return user, nil
	}

	enc, err := encodeRoleIDs(roleIDs)
	if err != nil {
		c.logger.Error().Err(err).Msg("Could not encode loaded roles")
	}

	if user.Opaque == nil {
		user.Opaque = &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"roles": enc,
			},
		}
	} else {
		user.Opaque.Map["roles"] = enc
	}

	return res.User, nil
}

func (c *cs3backend) Authenticate(ctx context.Context, username string, password string) (*cs3.User, error) {
	res, err := c.authProvider.Authenticate(ctx, &gateway.AuthenticateRequest{
		ClientId:     username,
		ClientSecret: password,
	})

	switch {
	case err != nil:
		return nil, fmt.Errorf("could not authenticate with username and password user: %s, %w", username, err)
	case res.Status.Code != rpcv1beta1.Code_CODE_OK:
		return nil, fmt.Errorf("could not authenticate with username and password user: %s, got code: %d", username, res.Status.Code)
	}

	return res.User, nil
}

func (c *cs3backend) CreateUserFromClaims(ctx context.Context, claims *oidc.StandardClaims) (*cs3.User, error) {
	return nil, fmt.Errorf("CS3 Backend does not support creating users from claims")
}

func (c cs3backend) GetUserGroups(ctx context.Context, userID string) {
	panic("implement me")
}

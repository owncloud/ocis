package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/token"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	settingsService "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/selector"
)

type cs3backend struct {
	graphSelector       selector.Selector
	settingsRoleService settingssvc.RoleService
	authProvider        RevaAuthenticator
	oidcISS             string
	machineAuthAPIKey   string
	tokenManager        token.Manager
	logger              log.Logger
}

// NewCS3UserBackend creates a user-provider which fetches users from a CS3 UserBackend
func NewCS3UserBackend(rs settingssvc.RoleService, ap RevaAuthenticator, machineAuthAPIKey string, oidcISS string, tokenManager token.Manager, logger log.Logger) UserBackend {
	reg := registry.GetRegistry()
	sel := selector.NewSelector(selector.Registry(reg))
	return &cs3backend{
		graphSelector:       sel,
		settingsRoleService: rs,
		authProvider:        ap,
		oidcISS:             oidcISS,
		machineAuthAPIKey:   machineAuthAPIKey,
		tokenManager:        tokenManager,
		logger:              logger,
	}
}

func (c *cs3backend) GetUserByClaims(ctx context.Context, claim, value string, withRoles bool) (*cs3.User, string, error) {
	res, err := c.authProvider.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     claim + ":" + value,
		ClientSecret: c.machineAuthAPIKey,
	})

	switch {
	case err != nil:
		return nil, "", fmt.Errorf("could not get user by claim %v with value %v: %w", claim, value, err)
	case res.Status.Code != rpcv1beta1.Code_CODE_OK:
		if res.Status.Code == rpcv1beta1.Code_CODE_NOT_FOUND {
			return nil, "", ErrAccountNotFound
		}
		return nil, "", fmt.Errorf("could not get user by claim %v with value %v : %w ", claim, value, err)
	}

	user := res.User

	if !withRoles {
		return user, res.Token, nil
	}

	var roleIDs []string
	if user.Id.Type != cs3.UserType_USER_TYPE_LIGHTWEIGHT {
		roleIDs, err = loadRolesIDs(ctx, user.Id.OpaqueId, c.settingsRoleService)
		if err != nil {
			var merr *merrors.Error
			if errors.As(err, &merr) && merr.Code == http.StatusNotFound {
				// This user doesn't have a role assignment yet. Assign a
				// default user role. At least until proper roles are provided. See
				// https://github.com/owncloud/ocis/v2/issues/1825 for more context.
				if user.Id.Type == cs3.UserType_USER_TYPE_PRIMARY {
					c.logger.Info().Str("userid", user.Id.OpaqueId).Msg("user has no role assigned, assigning default user role")
					_, err := c.settingsRoleService.AssignRoleToUser(ctx, &settingssvc.AssignRoleToUserRequest{
						AccountUuid: user.Id.OpaqueId,
						RoleId:      settingsService.BundleUUIDRoleUser,
					})
					if err != nil {
						c.logger.Error().Err(err).Msg("Could not add default role")
						return nil, "", err
					}
					roleIDs = append(roleIDs, settingsService.BundleUUIDRoleUser)
				}
			} else {
				c.logger.Error().Err(err).Msgf("Could not load roles")
				return nil, "", err
			}
		}
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

	return user, res.Token, nil
}

func (c *cs3backend) Authenticate(ctx context.Context, username string, password string) (*cs3.User, string, error) {
	res, err := c.authProvider.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "basic",
		ClientId:     username,
		ClientSecret: password,
	})

	switch {
	case err != nil:
		return nil, "", fmt.Errorf("could not authenticate with username and password user: %s, %w", username, err)
	case res.Status.Code != rpcv1beta1.Code_CODE_OK:
		return nil, "", fmt.Errorf("could not authenticate with username and password user: %s, got code: %d", username, res.Status.Code)
	}

	return res.User, res.Token, nil
}

// CreateUserFromClaims creates a new user via libregraph users API, taking the
// attributes from the provided `claims` map. On success it returns the new
// user. If the user already exist this is not considered an error and the
// function will just return the existing user.
func (c *cs3backend) CreateUserFromClaims(ctx context.Context, claims map[string]interface{}) (*cs3.User, error) {
	newctx := context.Background()
	token, err := c.generateAutoProvisionAdminToken(newctx)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error generating token for autoprovisioning user.")
		return nil, err
	}
	lgClient, err := c.setupLibregraphClient(ctx, token)
	if err != nil {
		c.logger.Error().Err(err).Msg("Error setting up libregraph client.")
		return nil, err
	}

	newUser, err := c.libregraphUserFromClaims(newctx, claims)
	if err != nil {
		c.logger.Error().Err(err).Interface("claims", claims).Msg("Error creating user from claims")
		return nil, fmt.Errorf("Error creating user from claims: %w", err)
	}

	req := lgClient.UsersApi.CreateUser(newctx).User(newUser)

	created, resp, err := req.Execute()
	var reread bool
	if err != nil {
		if resp == nil {
			return nil, err
		}

		// If the user already exists here, some other request did already create it in parallel.
		// So just issue a Debug message and ignore the libregraph error otherwise
		var lerr error
		if reread, lerr = c.isAlreadyExists(resp); lerr != nil {
			c.logger.Error().Err(lerr).Msg("extracting error from ibregraph response body failed.")
			return nil, err
		}
		if !reread {
			c.logger.Error().Err(err).Msg("Error creating user")
			return nil, err
		}
	}

	// User has been created meanwhile, re-read it to get the user id
	if reread {
		c.logger.Debug().Msg("User already exist, re-reading via libregraph")
		gureq := lgClient.UserApi.GetUser(newctx, newUser.GetOnPremisesSamAccountName())
		created, resp, err = gureq.Execute()
		if err != nil {
			c.logger.Error().Err(err).Msg("Error trying to re-read user from graphAPI")
			return nil, err
		}
	}

	cs3UserCreated := c.cs3UserFromLibregraph(newctx, created)

	return &cs3UserCreated, nil
}

func (c cs3backend) GetUserGroups(ctx context.Context, userID string) {
	panic("implement me")
}

func (c cs3backend) setupLibregraphClient(ctx context.Context, cs3token string) (*libregraph.APIClient, error) {
	// Use micro registry to resolve next graph service endpoint
	next, err := c.graphSelector.Select("com.owncloud.graph.graph")
	if err != nil {
		c.logger.Debug().Err(err).Msg("setupLibregraphClient: error during Select")
		return nil, err
	}
	node, err := next()
	if err != nil {
		c.logger.Debug().Err(err).Msg("setupLibregraphClient: error getting next Node")
		return nil, err
	}
	lgconf := libregraph.NewConfiguration()
	lgconf.Servers = libregraph.ServerConfigurations{
		{
			URL: fmt.Sprintf("%s://%s/graph/v1.0", node.Metadata["protocol"], node.Address),
		},
	}

	lgconf.DefaultHeader = map[string]string{revactx.TokenHeader: cs3token}
	return libregraph.NewAPIClient(lgconf), nil
}

func (c cs3backend) isAlreadyExists(resp *http.Response) (bool, error) {
	oDataErr := libregraph.NewOdataErrorWithDefaults()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Debug().Err(err).Msg("Error trying to read libregraph response")
		return false, err
	}
	err = json.Unmarshal(body, oDataErr)
	if err != nil {
		c.logger.Debug().Err(err).Msg("Error unmarshalling libregraph response")
		return false, err
	}

	if oDataErr.Error.Code == errorcode.NameAlreadyExists.String() {
		return true, nil
	}
	return false, nil
}

func (c cs3backend) libregraphUserFromClaims(ctx context.Context, claims map[string]interface{}) (libregraph.User, error) {
	var ok bool
	var dn, mail, username string
	user := libregraph.User{}
	if dn, ok = claims[oidc.Name].(string); !ok {
		return user, fmt.Errorf("Missing claim '%s'", oidc.Name)
	}
	if mail, ok = claims[oidc.Email].(string); !ok {
		return user, fmt.Errorf("Missing claim '%s'", oidc.Email)
	}
	if username, ok = claims[oidc.PreferredUsername].(string); !ok {
		c.logger.Warn().Str("claim", oidc.PreferredUsername).Msg("Missing claim for username, falling back to email address")
		username = mail
	}
	user.DisplayName = &dn
	user.OnPremisesSamAccountName = &username
	user.Mail = &mail
	return user, nil
}

func (c cs3backend) cs3UserFromLibregraph(ctx context.Context, lu *libregraph.User) cs3.User {
	cs3id := cs3.UserId{
		Type: cs3.UserType_USER_TYPE_PRIMARY,
		Idp:  c.oidcISS,
	}

	cs3id.OpaqueId = lu.GetId()

	cs3user := cs3.User{
		Id: &cs3id,
	}
	cs3user.Username = lu.GetOnPremisesSamAccountName()
	cs3user.DisplayName = lu.GetDisplayName()
	cs3user.Mail = lu.GetMail()
	return cs3user
}

// This returns an hardcoded internal User, that is privileged to create new User via
// the Graph API. This user is needed for autoprovisioning of users from incoming OIDC
// claims.
func getAutoProvisionUserCreator() (*cs3.User, error) {
	encRoleID, err := encodeRoleIDs([]string{settingsService.BundleUUIDRoleAdmin})
	if err != nil {
		return nil, err
	}

	autoProvisionUserCreator := &cs3.User{
		DisplayName: "Autoprovision User",
		Username:    "autoprovisioner",
		Id: &cs3.UserId{
			Idp:      "internal",
			OpaqueId: "autoprov-user-id00-0000-000000000000",
		},
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"roles": encRoleID,
			},
		},
	}
	return autoProvisionUserCreator, nil
}

func (c cs3backend) generateAutoProvisionAdminToken(ctx context.Context) (string, error) {
	userCreator, err := getAutoProvisionUserCreator()
	if err != nil {
		return "", err
	}

	s, err := scope.AddOwnerScope(nil)
	if err != nil {
		c.logger.Error().Err(err).Msg("could not get owner scope")
		return "", err
	}

	token, err := c.tokenManager.MintToken(ctx, userCreator, s)
	if err != nil {
		c.logger.Error().Err(err).Msg("could not mint token")
		return "", err
	}
	return token, nil
}

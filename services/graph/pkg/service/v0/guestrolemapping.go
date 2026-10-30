package svc

import (
	"context"
	"slices"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/utils"
	"go-micro.dev/v4/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/rolemapping"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

// guestRoleMapping maps a value to the role a multi-instance guest gets, see
// OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING. The proxy applies it to the roles claim at login.
type guestRoleMapping struct {
	value    rolemapping.Matcher
	roleName string
}

// parseGuestRoleMapping parses the pairs like the proxy does, and compiles their values once.
func parseGuestRoleMapping(pairs []string) ([]guestRoleMapping, error) {
	mappings := make([]guestRoleMapping, 0, len(pairs))
	for _, pair := range pairs {
		value, roleName, err := rolemapping.ParsePair(pair)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, guestRoleMapping{value: rolemapping.NewMatcher(value), roleName: roleName})
	}
	return mappings, nil
}

// guestVaultRoleID returns the id of the role the first mapping that matches the userClass gives,
// like the proxy picks it at login, if that role grants vault mode. It returns "" otherwise.
func guestVaultRoleID(mappings []guestRoleMapping, userClass []string, roles []*settingsmsg.Bundle) string {
	for _, m := range mappings {
		if !slices.ContainsFunc(userClass, m.value.Matches) {
			continue
		}
		for _, role := range roles {
			if role.GetName() == m.roleName && roleGrantsVaultMode(role) {
				return role.GetId()
			}
		}
		return ""
	}
	return ""
}

// guestVaultEligible reports whether the guest role mapping gives the user's userClass a role with
// vault mode. A failed lookup leaves them out, like members, who get their role from the role mapping
// of members at login.
func (g Graph) guestVaultEligible(ctx context.Context, userID string, roles []*settingsmsg.Bundle) bool {
	userClass, member, err := g.identityBackend.GetUserClass(ctx, userID)
	if err != nil || member {
		return false
	}
	return guestVaultRoleID(g.guestRoleMapping, userClass, roles) != ""
}

// withGuestVaultRole creates a vault share for the user with createShare. The vault share check only
// reads role assignments, and a user of another instance only gets a role here at their first login
// here. So if the user has no role yet, withGuestVaultRole first assigns the role the guest role
// mapping gives their userClass, if it grants vault mode, and removes it again if the share fails.
func (g Graph) withGuestVaultRole(ctx context.Context, userID string, createShare func() error) error {
	logger := g.logger.SubloggerWithRequestID(ctx)
	role, err := g.assignGuestVaultRole(ctx, userID)
	if err != nil {
		// without the role the vault share check refuses the share
		logger.Warn().Err(err).Str("userid", userID).Msg("could not assign the guest vault role")
	}
	if err := createShare(); err != nil {
		role.undo()
		return err
	}
	if err := role.finish(ctx); err != nil {
		logger.Warn().Err(err).Str("userid", userID).Msg("could not adjust the personal space to the guest vault role")
	}
	return nil
}

// guestVaultRole is a guest vault role assigned for a vault share that is about to be created. The
// zero value is no role, its methods do nothing then.
type guestVaultRole struct {
	g             Graph
	serviceCtx    context.Context
	gatewayClient gateway.GatewayAPIClient
	userID        string
	roleID        string
	assignmentID  string
}

// assignGuestVaultRole assigns a guest who has no role here yet the role the guest role mapping gives
// their userClass, if that role grants vault mode. It returns the zero value if it assigned no role.
//
// Only the mapping counts. A guest no mapping matches gets OCIS_MULTI_INSTANCE_GUEST_ROLE at login,
// which is not assigned here: such a guest can only receive vault shares after their first login.
func (g Graph) assignGuestVaultRole(ctx context.Context, userID string) (guestVaultRole, error) {
	// like the proxy, only for multi-instance guests
	if !g.config.MultiInstance.Enabled || len(g.guestRoleMapping) == 0 {
		return guestVaultRole{}, nil
	}
	assignments, err := g.roleService.ListRoleAssignments(ctx, &settingssvc.ListRoleAssignmentsRequest{AccountUuid: userID})
	if err != nil {
		return guestVaultRole{}, err
	}
	if len(assignments.GetAssignments()) > 0 {
		return guestVaultRole{}, nil
	}

	userClass, member, err := g.identityBackend.GetUserClass(ctx, userID)
	if err != nil {
		return guestVaultRole{}, err
	}
	// members get their role from the role mapping of members at login, not from the guest mapping
	if member {
		return guestVaultRole{}, nil
	}
	roles, err := g.roleService.ListRoles(ctx, &settingssvc.ListBundlesRequest{})
	if err != nil {
		return guestVaultRole{}, err
	}
	roleID := guestVaultRoleID(g.guestRoleMapping, userClass, roles.GetBundles())
	if roleID == "" {
		return guestVaultRole{}, nil
	}

	// assigning roles needs the role management permission. Start from a new context, the settings
	// service would check the role ids of the sharing user in the request context instead.
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return guestVaultRole{}, err
	}
	serviceCtx, err := utils.GetServiceUserContextWithContext(context.Background(), gatewayClient, g.config.ServiceAccount.ServiceAccountID, g.config.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return guestVaultRole{}, err
	}
	serviceCtx = metadata.Set(serviceCtx, middleware.AccountID, g.config.ServiceAccount.ServiceAccountID)
	assignment, err := g.roleService.AssignRoleToUser(serviceCtx, &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      roleID,
	})
	if err != nil {
		return guestVaultRole{}, err
	}
	return guestVaultRole{
		g:             g,
		serviceCtx:    serviceCtx,
		gatewayClient: gatewayClient,
		userID:        userID,
		roleID:        roleID,
		assignmentID:  assignment.GetAssignment().GetId(),
	}, nil
}

// undo removes the role, the user had none before
func (r guestVaultRole) undo() {
	if r.roleID == "" {
		return
	}
	r.g.revertRoleAssignment(r.serviceCtx, r.userID, "", r.assignmentID)
}

// finish adjusts the personal space to the role and announces the change, like the other role
// assignments. The share exists, so the role stays even if the personal space can't be adjusted.
func (r guestVaultRole) finish(ctx context.Context) error {
	if r.roleID == "" {
		return nil
	}
	canCreateDrives, err := r.g.checkUserPermission(r.serviceCtx, "Drives.Create", r.userID, r.gatewayClient)
	switch {
	case err != nil:
		// leave the personal space as it is when it's unclear what the role allows
	case canCreateDrives:
		err = shared.RestorePersonalSpace(r.serviceCtx, r.gatewayClient, r.userID)
	default:
		err = shared.DisablePersonalSpace(r.serviceCtx, r.gatewayClient, r.userID)
	}

	oldRole := ""
	e := events.UserFeatureChanged{
		UserID: r.userID,
		Features: []events.UserFeature{
			{
				Name:     "roleChanged",
				Value:    r.roleID,
				OldValue: &oldRole,
			},
		},
		Timestamp: utils.TSNow(),
	}
	if currentUser, ok := revactx.ContextGetUser(ctx); ok {
		e.Executant = currentUser.GetId()
	}
	r.g.publishEvent(ctx, e)
	return err
}

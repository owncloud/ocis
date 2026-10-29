package svc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/utils"
	"go-micro.dev/v4/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

// guestRoleMapping maps a value to the role a multi-instance guest gets, see
// OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING. The proxy applies it to the roles claim at login.
type guestRoleMapping struct {
	value    string
	roleName string
}

// parseGuestRoleMapping parses 'value:rolename' pairs like the proxy does. It splits at the last
// colon, since role names can't contain one.
func parseGuestRoleMapping(pairs []string) ([]guestRoleMapping, error) {
	mappings := make([]guestRoleMapping, 0, len(pairs))
	for _, pair := range pairs {
		i := strings.LastIndex(pair, ":")
		value, roleName := "", ""
		if i >= 0 {
			value, roleName = strings.TrimSpace(pair[:i]), strings.TrimSpace(pair[i+1:])
		}
		if value == "" || roleName == "" {
			return nil, fmt.Errorf("invalid guest role mapping '%s', expected 'claimvalue:rolename'", pair)
		}
		mappings = append(mappings, guestRoleMapping{value: value, roleName: roleName})
	}
	return mappings, nil
}

// guestRoleName returns the role of the first mapping that matches one of the values, like the
// proxy does at login. It returns "" when no mapping matches.
func guestRoleName(mappings []guestRoleMapping, values []string) string {
	for _, m := range mappings {
		if matchesGuestRoleValue(m.value, values) {
			return m.roleName
		}
	}
	return ""
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
		logger.Warn().Err(err).Str("userid", userID).Msg("could not finish the guest vault role, it was removed again")
	}
	return nil
}

// guestVaultRole is a guest vault role assigned for a vault share that is about to be created
type guestVaultRole struct {
	g             Graph
	serviceCtx    context.Context
	gatewayClient gateway.GatewayAPIClient
	userID        string
	roleID        string
	assignmentID  string
}

// assignGuestVaultRole assigns a guest who has no role here yet the role the guest role mapping gives
// their userClass, if that role grants vault mode. It returns nil if it assigned no role.
//
// Only the mapping counts. A guest no mapping matches gets OCIS_MULTI_INSTANCE_GUEST_ROLE at login,
// which is not assigned here: such a guest can only receive vault shares after their first login.
func (g Graph) assignGuestVaultRole(ctx context.Context, userID string) (*guestVaultRole, error) {
	// like the proxy, only for multi-instance guests
	if !g.config.MultiInstance.Enabled || len(g.guestRoleMapping) == 0 {
		return nil, nil
	}
	assignments, err := g.roleService.ListRoleAssignments(ctx, &settingssvc.ListRoleAssignmentsRequest{AccountUuid: userID})
	if err != nil {
		return nil, err
	}
	if len(assignments.GetAssignments()) > 0 {
		return nil, nil
	}

	userClass, member, err := g.identityBackend.GetUserClass(ctx, userID)
	if err != nil {
		return nil, err
	}
	// members get their role from the role mapping of members at login, not from the guest mapping
	if member {
		return nil, nil
	}
	roleName := guestRoleName(g.guestRoleMapping, userClass)
	if roleName == "" {
		return nil, nil
	}
	roles, err := g.roleService.ListRoles(ctx, &settingssvc.ListBundlesRequest{})
	if err != nil {
		return nil, err
	}
	roleID := ""
	for _, role := range roles.GetBundles() {
		if role.GetName() == roleName && roleGrantsVaultMode(role) {
			roleID = role.GetId()
			break
		}
	}
	if roleID == "" {
		return nil, nil
	}

	// assigning roles needs the role management permission. Start from a new context, the settings
	// service would check the role ids of the sharing user in the request context instead.
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	serviceCtx, err := utils.GetServiceUserContextWithContext(context.Background(), gatewayClient, g.config.ServiceAccount.ServiceAccountID, g.config.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return nil, err
	}
	serviceCtx = metadata.Set(serviceCtx, middleware.AccountID, g.config.ServiceAccount.ServiceAccountID)
	assignment, err := g.roleService.AssignRoleToUser(serviceCtx, &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      roleID,
	})
	if err != nil {
		return nil, err
	}
	return &guestVaultRole{
		g:             g,
		serviceCtx:    serviceCtx,
		gatewayClient: gatewayClient,
		userID:        userID,
		roleID:        roleID,
		assignmentID:  assignment.GetAssignment().GetId(),
	}, nil
}

// undo removes the role, the user had none before
func (r *guestVaultRole) undo() {
	if r == nil {
		return
	}
	r.g.revertRoleAssignment(r.serviceCtx, r.userID, "", r.assignmentID)
}

// finish adjusts the personal space to the role and announces the change, like the other role
// assignments. If the personal space can't be adjusted, it removes the role again.
func (r *guestVaultRole) finish(ctx context.Context) error {
	if r == nil {
		return nil
	}
	canCreateDrives, err := r.g.checkUserPermission(r.serviceCtx, "Drives.Create", r.userID, r.gatewayClient)
	if err == nil {
		if canCreateDrives {
			err = shared.RestorePersonalSpace(r.serviceCtx, r.gatewayClient, r.userID)
		} else {
			err = shared.DisablePersonalSpace(r.serviceCtx, r.gatewayClient, r.userID)
		}
	}
	if err != nil {
		r.undo()
		return err
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
	return nil
}

// matchesGuestRoleValue matches like the proxy's role mapping: the mapping value is an exact value
// or a regular expression that must match a whole value.
func matchesGuestRoleValue(mappingValue string, values []string) bool {
	rx, err := regexp.Compile("^(?:" + mappingValue + ")$")
	for _, v := range values {
		if v == mappingValue || (err == nil && rx.MatchString(v)) {
			return true
		}
	}
	return false
}

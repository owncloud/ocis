package svc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/owncloud/reva/v2/pkg/utils"
	"go-micro.dev/v4/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
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

// assignGuestVaultRole assigns a user who has no role here yet the role their first login here will
// give them, if that role grants vault mode. The vault share check only reads role assignments, so
// without it a user of another instance can't receive a vault share before logging in here.
func (g Graph) assignGuestVaultRole(ctx context.Context, userID string) error {
	if len(g.guestRoleMapping) == 0 {
		return nil
	}
	assignments, err := g.roleService.ListRoleAssignments(ctx, &settingssvc.ListRoleAssignmentsRequest{AccountUuid: userID})
	if err != nil {
		return err
	}
	if len(assignments.GetAssignments()) > 0 {
		return nil
	}

	// GetPreciseUser looks a user up by id when there is no instance name
	_, userClass, err := g.identityBackend.GetPreciseUser(ctx, userID, "", nil)
	if err != nil {
		return err
	}
	roleName := guestRoleName(g.guestRoleMapping, userClass)
	if roleName == "" {
		return nil
	}
	roles, err := g.roleService.ListRoles(ctx, &settingssvc.ListBundlesRequest{})
	if err != nil {
		return err
	}
	roleID := ""
	for _, role := range roles.GetBundles() {
		if role.GetName() == roleName && roleGrantsVaultMode(role) {
			roleID = role.GetId()
			break
		}
	}
	if roleID == "" {
		return nil
	}

	// assigning roles needs the role management permission. Start from a new context, the settings
	// service would check the role ids of the sharing user in the request context instead.
	gatewayClient, err := g.gatewaySelector.Next()
	if err != nil {
		return err
	}
	serviceCtx, err := utils.GetServiceUserContextWithContext(context.Background(), gatewayClient, g.config.ServiceAccount.ServiceAccountID, g.config.ServiceAccount.ServiceAccountSecret)
	if err != nil {
		return err
	}
	serviceCtx = metadata.Set(serviceCtx, middleware.AccountID, g.config.ServiceAccount.ServiceAccountID)
	_, err = g.roleService.AssignRoleToUser(serviceCtx, &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      roleID,
	})
	return err
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

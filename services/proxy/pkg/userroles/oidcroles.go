package userroles

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"sync"
	"time"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
	"go-micro.dev/v4/metadata"
)

// ErrNoRoleAssigned is returned by UpdateUserRoleAssignment when the user's claims
// yield no ocis role and no default role is configured. It is a property of the
// user's account rather than a server-side failure, so callers should report it as
// such instead of as a generic internal error.
var ErrNoRoleAssigned = errors.New("no role in claim maps to an ocis role and no default role is configured")

// ErrRolesClaimNotSet is returned by extractRoles when the token carries no roles for
// this user at all: the configured claim is absent, or present and null. The token is
// well formed and simply says nothing about roles, which is the normal shape for users
// federated into the IDP from an external directory. It is deliberately distinct from
// the errors returned for a claim that is present and unreadable - only one of the two
// is a sign that something is wrong.
var ErrRolesClaimNotSet = errors.New("no roles claim in user claims")

type oidcRoleAssigner struct {
	Options
}

// NewOIDCRoleAssigner returns an implementation of the UserRoleAssigner interface
func NewOIDCRoleAssigner(opts ...Option) UserRoleAssigner {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}

	return oidcRoleAssigner{
		Options: opt,
	}
}

func extractRoles(rolesClaim string, claims map[string]interface{}) (map[string]struct{}, error) {

	claimRoles := map[string]struct{}{}
	// happy path
	value, _ := claims[rolesClaim].(string)
	if value != "" {
		claimRoles[value] = struct{}{}
		return claimRoles, nil
	}

	segments := oidc.SplitWithEscaping(rolesClaim, ".", "\\")
	claim, err := oidc.WalkSegments(segments, claims)
	if err != nil {
		// WalkSegments reports a missing intermediate segment and a segment that is
		// present but is not an object with the same "unsupported type" error. Only
		// the first is a token that is merely silent about roles.
		if claimPathAbsent(segments, claims) {
			return nil, ErrRolesClaimNotSet
		}
		return nil, err
	}

	switch v := claim.(type) {
	case nil:
		// The path resolved, but there is nothing at the end of it.
		return nil, ErrRolesClaimNotSet
	case []string:
		for _, cr := range v {
			claimRoles[cr] = struct{}{}
		}
	case []interface{}:
		for _, cri := range v {
			cr, ok := cri.(string)
			if !ok {
				err := errors.New("invalid role in claims")
				return nil, err
			}

			claimRoles[cr] = struct{}{}
		}
	case string:
		claimRoles[v] = struct{}{}
	default:
		return nil, fmt.Errorf("roles claim %q holds an unusable %T", rolesClaim, claim)
	}

	return claimRoles, nil
}

// claimPathAbsent reports whether the dotted claim path is simply not present in the
// token, as opposed to present and holding something we cannot read. Every segment
// before the last must exist and be an object for the walk to continue: a missing
// segment means the token is silent about roles, while a segment that exists but is
// not an object means the token says something about them that we cannot make sense
// of. It mirrors the traversal in oidc.WalkSegments, which reports both as one error.
func claimPathAbsent(segments []string, claims map[string]interface{}) bool {
	for i := 0; i < len(segments)-1; i++ {
		switch next := claims[segments[i]].(type) {
		case nil:
			return true
		case map[string]interface{}:
			claims = next
		case map[interface{}]interface{}:
			converted := make(map[string]interface{}, len(next))
			for k, v := range next {
				s, ok := k.(string)
				if !ok {
					return false
				}
				converted[s] = v
			}
			claims = converted
		default:
			return false
		}
	}
	return claims[segments[len(segments)-1]] == nil
}

// matchesClaimMapping returns true if the provided mapping pattern matches at least
// one of the values present in claimRoles. It supports:
// - exact match when ClaimValue is a literal equal to a claim value
// - regex match when ClaimValue is a regex pattern (e.g. "ocis-user-.*")
// The regex is matched against the entire claim value, not a substring.
func matchesClaimMapping(mappingValue string, claimRoles map[string]struct{}) bool {
	if _, ok := claimRoles[mappingValue]; ok {
		return true
	}

	rx, err := regexp.Compile("^(?:" + mappingValue + ")$")
	if err != nil {
		return false
	}
	for cr := range claimRoles {
		if rx.MatchString(cr) {
			return true
		}
	}
	return false
}

// UpdateUserRoleAssignment assigns the role "User" to the supplied user. Unless the user
// already has a different role assigned.
func (ra oidcRoleAssigner) UpdateUserRoleAssignment(ctx context.Context, user *cs3user.User, claims map[string]interface{}, overwriteRole string) (*cs3user.User, error) {
	userID := user.GetId().GetOpaqueId()
	logger := ra.logger.SubloggerWithRequestID(ctx).With().Str("userid", userID).Logger()
	roleNamesToRoleIDs, err := ra.roleNamesToRoleIDs(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error mapping role names to role ids")
		return nil, err
	}

	roleIDFromClaim := roleNamesToRoleIDs[overwriteRole]
	if overwriteRole == "" {
		// A token that says nothing about roles and a token we cannot read are two
		// different situations and are handled differently: the first falls back to
		// the configured default role, the second is reported.
		claimRoles, err := extractRoles(ra.rolesClaim, claims)
		switch {
		case errors.Is(err, ErrRolesClaimNotSet):
			// The token is well formed and carries no roles for this user. This is the
			// common case when users are federated into the IDP from an external
			// directory and simply have no role attached, and it is what this login
			// path exists to serve: fall through to the role mapping and the default
			// role below.
			logger.Debug().Str("Claim", ra.rolesClaim).Msg("No roles claim in user claims")
			claimRoles = nil
		case err != nil:
			// The claim is present and unreadable, so something is wrong with the
			// token or with PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM. Signing the user in on
			// the default role would paper over it, so report it instead.
			logger.Error().Err(err).Str("Claim", ra.rolesClaim).Msg("Could not extract roles from claims")
			return nil, err
		case len(claimRoles) == 0:
			logger.Debug().Str("Claim", ra.rolesClaim).Msg("No roles set in claim")
		}

		// the roleMapping config is supposed to have the role mappings ordered from the highest privileged role
		// down to the lowest privileged role. Since ocis currently only can handle a single role assignment we
		// pick the highest privileged role that matches a value from the claims
		for _, mapping := range ra.Options.roleMapping {
			if matchesClaimMapping(mapping.ClaimValue, claimRoles) {
				logger.Debug().Str("ocisRole", mapping.RoleName).Str("role id", roleNamesToRoleIDs[mapping.RoleName]).Msg("first matching role")
				roleIDFromClaim = roleNamesToRoleIDs[mapping.RoleName]
				break
			}
		}

		if roleIDFromClaim == "" {
			roleIDFromClaim, err = ra.fallbackRoleID(logger, roleNamesToRoleIDs, claimRoles)
			if err != nil {
				return nil, err
			}
		}
	}

	assignedRoles, err := loadRolesIDs(ctx, userID, ra.roleService)
	if err != nil {
		logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}
	if len(assignedRoles) > 1 {
		logger.Error().Str("userID", userID).Int("numRoles", len(assignedRoles)).Msg("The user has too many roles assigned")
	}
	logger.Debug().Interface("assignedRoleIds", assignedRoles).Msg("Currently assigned roles")

	if len(assignedRoles) != 1 || (assignedRoles[0] != roleIDFromClaim) {
		logger.Debug().Interface("assignedRoleIds", assignedRoles).Interface("newRoleId", roleIDFromClaim).Msg("Updating role assignment for user")
		var oldRole string
		if len(assignedRoles) > 0 {
			oldRole = assignedRoles[0]
		}
		newctx, err := ra.prepareAdminContext(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating admin context")
			return nil, err
		}
		assignResp, err := ra.roleService.AssignRoleToUser(newctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: userID,
			RoleId:      roleIDFromClaim,
		})
		if err != nil {
			logger.Error().Err(err).Msg("Role assignment failed")
			return nil, err
		}

		userID := user.GetId().GetOpaqueId()
		client, err := ra.gatewaySelector.Next()
		if err != nil {
			return nil, err
		}

        // TODO: check if it's ctx or newctx
		canCreateDrives, err := utils.CheckPermission(revactx.ContextSetUser(ctx, user), "Drives.Create", client)
		if err != nil {
			// The permission could not be determined. Fail closed: leave the personal
			// space untouched rather than disabling (trashing) it on an indeterminate
			// result, and revert the role assignment so the user is not left in a
			// half-applied state. The role transition is retried on the next login.
			logger.Error().Any("userID", userID).Err(err).Msg("could not determine Drives.Create permission, reverting role assignment and leaving personal space unchanged")
			ra.revertRoleAssignment(newctx, userID, oldRole, assignResp.GetAssignment().GetId())
			return nil, err
		}
		if canCreateDrives {
			libreUser := identity.CreateUserModelFromCS3(user)
			err = shared.RestorePersonalSpace(newctx, client, libreUser.GetId())
			if err != nil {
				logger.Error().Any("userID", userID).Err(err).Msg("can't ensure the personal space")
				return nil, err
			}
		} else {
			err := shared.DisablePersonalSpace(newctx, client, userID)
			if err != nil {
				logger.Error().Any("userID", userID).Err(err).Msg("can't disable the personal space")
				return nil, err
			}
		}
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", []string{roleIDFromClaim})
	return user, nil
}

// fallbackRoleID resolves the configured default role for a user whose claims matched
// no role mapping. It returns ErrNoRoleAssigned when no default role is configured, and
// a descriptive error when one is configured but does not exist in the settings
// service - a misconfiguration is worth reporting differently from a user without a
// role, because only one of the two is fixed by editing the deployment.
func (ra oidcRoleAssigner) fallbackRoleID(logger zerolog.Logger, roleNamesToRoleIDs map[string]string, claimRoles map[string]struct{}) (string, error) {
	if ra.defaultRole == "" {
		logger.Error().
			Str("claim", ra.rolesClaim).
			Strs("claimValues", slices.Sorted(maps.Keys(claimRoles))).
			Interface("roleMapping", ra.roleMapping).
			Msg("No role mapping matched the user's claim and PROXY_ROLE_ASSIGNMENT_OIDC_DEFAULT_ROLE is not set. " +
				"Add a matching entry to the role mapping, or set a default role to let such users log in.")
		return "", ErrNoRoleAssigned
	}

	roleID := roleNamesToRoleIDs[ra.defaultRole]
	if roleID == "" {
		err := fmt.Errorf("the configured default role %q does not exist", ra.defaultRole)
		logger.Error().Err(err).
			Strs("knownRoles", slices.Sorted(maps.Keys(roleNamesToRoleIDs))).
			Msg("PROXY_ROLE_ASSIGNMENT_OIDC_DEFAULT_ROLE names a role that the settings service does not know")
		return "", err
	}

	logger.Warn().
		Str("claim", ra.rolesClaim).
		Strs("claimValues", slices.Sorted(maps.Keys(claimRoles))).
		Str("ocisRole", ra.defaultRole).
		Msg("No role mapping matched the user's claim, assigning the configured default role. " +
			"Add a matching entry to PROXY_ROLE_ASSIGNMENT_OIDC_ROLE_MAPPING if this user should get a different role.")
	return roleID, nil
}

// ApplyUserRole it looks up the user's role in the settings service and adds it
// user's opaque data
func (ra oidcRoleAssigner) ApplyUserRole(ctx context.Context, user *cs3user.User) (*cs3user.User, error) {
	roleIDs, err := loadRolesIDs(ctx, user.Id.OpaqueId, ra.roleService)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Could not load roles")
		return nil, err
	}

	user.Opaque = utils.AppendJSONToOpaque(user.Opaque, "roles", roleIDs)
	return user, nil
}

func (ra oidcRoleAssigner) prepareAdminContext(ctx context.Context) (context.Context, error) {
	gatewayClient, err := ra.gatewaySelector.Next()
	if err != nil {
		ra.logger.Error().Err(err).Msg("could not select next gateway client")
		return nil, err
	}
	newctx, err := utils.GetServiceUserContextWithContext(ctx, gatewayClient, ra.serviceAccount.ServiceAccountID, ra.serviceAccount.ServiceAccountSecret)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error preparing request context for provisioning role assignments.")
		return nil, err
	}

	newctx = metadata.Set(newctx, middleware.AccountID, ra.serviceAccount.ServiceAccountID)
	return newctx, nil
}

type roleNameToIDCache struct {
	roleNameToID map[string]string
	lastRead     time.Time
	lock         sync.RWMutex
}

var roleNameToID roleNameToIDCache

func (ra oidcRoleAssigner) roleNamesToRoleIDs(ctx context.Context) (map[string]string, error) {
	cacheTTL := 5 * time.Minute
	roleNameToID.lock.RLock()

	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < cacheTTL {
		defer roleNameToID.lock.RUnlock()
		return roleNameToID.roleNameToID, nil
	}
	ra.logger.Debug().Msg("refreshing roles ids")

	// cache needs Refresh get a write lock
	roleNameToID.lock.RUnlock()
	roleNameToID.lock.Lock()
	defer roleNameToID.lock.Unlock()

	// check again, another goroutine might have updated while we "upgraded" the lock
	if !roleNameToID.lastRead.IsZero() && time.Since(roleNameToID.lastRead) < cacheTTL {
		return roleNameToID.roleNameToID, nil
	}

	// Get all roles to find the role IDs.
	// To list roles we need some elevated access to the settings service
	// prepare a new request context for that until we have service accounts
	newctx, err := ra.prepareAdminContext(ctx)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Error creating admin context")
		return nil, err
	}

	req := &settingssvc.ListBundlesRequest{}
	res, err := ra.roleService.ListRoles(newctx, req)
	if err != nil {
		ra.logger.Error().Err(err).Msg("Failed to list all roles")
		return map[string]string{}, err
	}

	newIDs := map[string]string{}
	for _, role := range res.Bundles {
		ra.logger.Debug().Str("role", role.Name).Str("id", role.Id).Msg("Got Role")
		newIDs[role.Name] = role.Id
	}
	ra.logger.Debug().Interface("roleMap", newIDs).Msg("Role Name to role ID map")
	roleNameToID.roleNameToID = newIDs
	roleNameToID.lastRead = time.Now()
	return roleNameToID.roleNameToID, nil
}

// revertRoleAssignment best-effort restores the user's previous role after a
// failed reconciliation, so the user is not left with the new role applied but
// the personal space unreconciled. A user has exactly one role, so this means
// re-assigning the previous role, or removing the new assignment if there was
// none. Failures are only logged; the next login reconciles idempotently.
func (ra oidcRoleAssigner) revertRoleAssignment(ctx context.Context, userID, oldRoleID, newAssignmentID string) {
	if oldRoleID != "" {
		if _, err := ra.roleService.AssignRoleToUser(ctx, &settingssvc.AssignRoleToUserRequest{
			AccountUuid: userID,
			RoleId:      oldRoleID,
		}); err != nil {
			ra.logger.Error().Any("userID", userID).Str("roleID", oldRoleID).Err(err).Msg("could not revert role assignment to previous role")
		}
		return
	}
	if _, err := ra.roleService.RemoveRoleFromUser(ctx, &settingssvc.RemoveRoleFromUserRequest{
		Id: newAssignmentID,
	}); err != nil {
		ra.logger.Error().Any("userID", userID).Str("assignmentID", newAssignmentID).Err(err).Msg("could not revert role assignment by removing the new assignment")
	}
}

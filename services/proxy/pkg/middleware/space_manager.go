package middleware

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/claimsmapper"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/reva/v2/pkg/conversions"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// SpaceManager return a middleware that manages space memberships
func SpaceManager(cfg config.ClaimSpaceManagement, opts ...Option) func(next http.Handler) http.Handler {
	options := newOptions(opts...)
	logger := options.Logger

	var cm claimsmapper.ClaimsMapper
	if cfg.Enabled {
		cm = claimsmapper.NewClaimsMapper(cfg.Regexp, cfg.Mapping)
	}

	return func(next http.Handler) http.Handler {
		return &claimSpaceManager{
			next:                 next,
			logger:               logger,
			gws:                  options.RevaGatewaySelector,
			mapper:               cm,
			serviceAccountID:     options.ServiceAccountID,
			serviceAccountSecret: options.ServiceAccountSecret,
			claimName:            cfg.Claim,
			enabled:              cfg.Enabled,
		}
	}
}

type claimSpaceManager struct {
	next                 http.Handler
	logger               log.Logger
	gws                  pool.Selectable[gateway.GatewayAPIClient]
	mapper               claimsmapper.ClaimsMapper
	serviceAccountID     string
	serviceAccountSecret string
	claimName            string
	enabled              bool
}

func (csm claimSpaceManager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer csm.next.ServeHTTP(w, req)

	if !csm.enabled {
		return
	}

	userid, spaceAssignments := csm.evaluateContext(req.Context())
	if userid == "" {
		// no user in context, we omit this request
		return
	}

	ctx, gwc, err := csm.getCtx()
	if err != nil {
		csm.logger.Error().Err(err).Msg("could not get service user context")
		return
	}

	// get all project spaces
	res, err := gwc.ListStorageSpaces(ctx, listStorageSpaceRequest())
	if err != nil {
		csm.logger.Error().Err(err).Msg("error doing grpc request")
		return
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		csm.logger.Error().Str("message", res.GetStatus().GetMessage()).Msg("unexpected status code doing listspaces request")
		return
	}

	for _, s := range res.GetStorageSpaces() {
		hasAccess, actualPerms, err := getSpaceMemberStatus(s, userid)
		if err != nil {
			csm.logger.Error().Err(err).Msg("error extracting space member")
			continue
		}

		desiredRole := conversions.RoleFromName(spaceAssignments[s.GetRoot().GetOpaqueId()])
		shouldHaveAccess := desiredRole.Name != conversions.RoleUnknown

		switch {
		case shouldHaveAccess && !hasAccess:
			// add user to space
			res, err := gwc.CreateShare(ctx, createShareRequest(userid, s, desiredRole.CS3ResourcePermissions()))
			if err != nil {
				csm.logger.Error().Err(err).Msg("error adding space member")
				continue
			}
			if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
				csm.logger.Error().Str("message", res.GetStatus().GetMessage()).Msg("unexpected status code doing createshare request")
				continue
			}

		case !shouldHaveAccess && hasAccess:
			// remove user from space
			res, err := gwc.RemoveShare(ctx, removeShareRequest(userid, s))
			if err != nil {
				csm.logger.Error().Err(err).Msg("error removing space member")
				continue
			}
			if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
				csm.logger.Error().Str("message", res.GetStatus().GetMessage()).Msg("unexpected status code doing removeshare request")
				continue
			}

		case shouldHaveAccess && hasAccess && !permissionsEqual(actualPerms, desiredRole.CS3ResourcePermissions()):
			// update user permissions
			res, err := gwc.UpdateShare(ctx, updateShareRequest(userid, s, desiredRole.CS3ResourcePermissions()))
			if err != nil {
				csm.logger.Error().Err(err).Msg("error updating space member")
				continue
			}
			if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
				csm.logger.Error().Str("message", res.GetStatus().GetMessage()).Msg("unexpected status code doing updateshare request")
				continue
			}
		}
	}
}

// returns the service user context and the gateway client
func (csm claimSpaceManager) getCtx() (context.Context, gateway.GatewayAPIClient, error) {
	gwc, err := csm.gws.Next()
	if err != nil {
		csm.logger.Error().Err(err).Msg("could not get gateway client")
		return nil, nil, err
	}
	ctx, err := utils.GetServiceUserContext(csm.serviceAccountID, gwc, csm.serviceAccountSecret)
	return ctx, gwc, err
}

// returns the userid and the space assignments from the context
func (csm claimSpaceManager) evaluateContext(ctx context.Context) (string, map[string]string) {
	u, _ := revactx.ContextGetUser(ctx)
	return u.GetId().GetOpaqueId(), csm.getSpaceAssignments(ctx)
}

// returns a map[spaceID]role
func (csm claimSpaceManager) getSpaceAssignments(ctx context.Context) map[string]string {
	claims := oidc.FromContext(ctx)
	values, ok := claims[csm.claimName].([]any)
	if !ok {
		csm.logger.Error().Interface("claims", claims).Str("claimname", csm.claimName).Msg("configured claims are not an array")
	}

	assignments := make(map[string]string)
	for _, ent := range values {
		e, ok := ent.(string)
		if !ok {
			csm.logger.Error().Interface("assignment", ent).Msg("assignment is not a string")
			continue
		}

		match, spaceid, role := csm.mapper.Exec(e)
		if !match {
			continue
		}
		assignments[spaceid] = chooseRole(role, assignments[spaceid])
	}

	return assignments
}

// will return the role with the highest permissions.
func chooseRole(roleA, roleB string) string {
	if roleA == "" {
		return roleB
	}

	if roleB == "" {
		return roleA
	}

	permsA := conversions.RoleFromName(roleA).CS3ResourcePermissions()
	permsB := conversions.RoleFromName(roleB).CS3ResourcePermissions()

	if conversions.SufficientCS3Permissions(permsA, permsB) {
		return roleA
	}
	// Note: This could be an issue if roleB does not contain roleA
	return roleB

}

func getSpaceMemberStatus(space *storageprovider.StorageSpace, userid string) (bool, *storageprovider.ResourcePermissions, error) {
	var permissionsMap map[string]*storageprovider.ResourcePermissions
	if err := utils.ReadJSONFromOpaque(space.GetOpaque(), "grants", &permissionsMap); err != nil {
		return false, nil, err
	}

	for id, perm := range permissionsMap {
		if id == userid {
			return true, perm, nil
		}
	}
	return false, nil, nil
}

func permissionsEqual(p1, p2 *storageprovider.ResourcePermissions) bool {
	if !conversions.SufficientCS3Permissions(p1, p2) {
		return false
	}
	if !conversions.SufficientCS3Permissions(p2, p1) {
		return false
	}
	return true
}

func listStorageSpaceRequest() *storageprovider.ListStorageSpacesRequest {
	return &storageprovider.ListStorageSpacesRequest{
		Opaque: utils.AppendPlainToOpaque(nil, "unrestricted", "true"),
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: "project",
				},
			},
		},
	}
}

func createShareRequest(userid string, space *storageprovider.StorageSpace, perms *storageprovider.ResourcePermissions) *collaboration.CreateShareRequest {
	return &collaboration.CreateShareRequest{
		ResourceInfo: space.GetRootInfo(),
		Grant: &collaboration.ShareGrant{
			Grantee: &storageprovider.Grantee{
				Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
				Id: &storageprovider.Grantee_UserId{UserId: &userpb.UserId{
					OpaqueId: userid,
				}},
			},
			Permissions: &collaboration.SharePermissions{
				Permissions: perms,
			},
		},
	}
}

func removeShareRequest(userid string, space *storageprovider.StorageSpace) *collaboration.RemoveShareRequest {
	return &collaboration.RemoveShareRequest{
		Ref: &collaboration.ShareReference{
			Spec: &collaboration.ShareReference_Key{
				Key: &collaboration.ShareKey{
					ResourceId: space.GetRoot(),
					Grantee:    buildGrantee(userid)},
			},
		},
	}
}

func updateShareRequest(userid string, s *storageprovider.StorageSpace, perms *storageprovider.ResourcePermissions) *collaboration.UpdateShareRequest {
	o := &types.Opaque{
		Map: map[string]*types.OpaqueEntry{
			"spacegrant": {},
		},
	}
	o = utils.AppendPlainToOpaque(o, "spacetype", "project")
	return &collaboration.UpdateShareRequest{
		Share: &collaboration.Share{
			ResourceId:  s.GetRoot(),
			Grantee:     buildGrantee(userid),
			Permissions: &collaboration.SharePermissions{Permissions: perms},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"permissions"},
		},
		Opaque: o,
	}

}

func buildGrantee(userid string) *storageprovider.Grantee {
	return &storageprovider.Grantee{
		Type: storageprovider.GranteeType_GRANTEE_TYPE_USER,
		Id: &storageprovider.Grantee_UserId{
			UserId: &userpb.UserId{
				OpaqueId: userid,
			},
		},
	}
}

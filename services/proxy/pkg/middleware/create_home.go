package middleware

import (
	"net/http"
	"strconv"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	"google.golang.org/grpc/metadata"
)

// CreateHome provides a middleware which sends a CreateHome request to the reva gateway
func CreateHome(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		return &createHome{
			next:                next,
			logger:              logger,
			revaGatewaySelector: options.RevaGatewaySelector,
			roleQuotas:          options.RoleQuotas,
		}
	}
}

type createHome struct {
	next                http.Handler
	logger              log.Logger
	revaGatewaySelector pool.Selectable[gateway.GatewayAPIClient]
	roleQuotas          map[string]uint64
}

func (m createHome) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	token := req.Header.Get("x-access-token")

	// we need to pass the token to authenticate the CreateHome request.
	//ctx := tokenpkg.ContextSetToken(r.Context(), token)
	ctx := metadata.AppendToOutgoingContext(req.Context(), revactx.TokenHeader, token)

	createHomeReq := &provider.CreateHomeRequest{}
	u, ok := revactx.ContextGetUser(ctx)
	if ok {
		roleIDs, err := m.getUserRoles(u)
		if err != nil {
			m.logger.Error().Err(err).Str("userid", u.Id.OpaqueId).Msg("failed to get roles for user")
			errorcode.GeneralException.Render(w, req, http.StatusInternalServerError, "Unauthorized")
			return
		}
		if limit, hasLimit := m.checkRoleQuotaLimit(roleIDs); hasLimit {
			createHomeReq.Opaque = utils.AppendPlainToOpaque(nil, "quota", strconv.FormatUint(limit, 10))
		}
	}

	client, err := m.revaGatewaySelector.Next()
	if err != nil {
		m.logger.Err(err).Msg("error selecting next gateway client")
	} else {
		createHomeRes, err := client.CreateHome(ctx, createHomeReq)
		switch {
		case err != nil:
			m.logger.Err(err).Msg("error calling CreateHome")
		case createHomeRes.GetStatus().GetCode() == rpc.Code_CODE_OK:
			m.logger.Debug().Interface("userID", u.GetId().GetOpaqueId()).Msg("personal space created")
		case createHomeRes.GetStatus().GetCode() == rpc.Code_CODE_ALREADY_EXISTS:
			m.logger.Info().Interface("userID", u.GetId().GetOpaqueId()).Interface("status", createHomeRes.GetStatus()).Msg("===== personal space already exists")
		default:
			m.logger.Error().Interface("userID", u.GetId().GetOpaqueId()).Interface("status", createHomeRes.GetStatus()).Msg("===== personal space creation failed")
		}

		// TODO: issue: The vault personal space can not be created after regular personal space creation.
		// If the regular personal cometed out the valut creation pass
		//
		// Create vault personal space
		// Inject storage_id into opaque for vault personal space
		createHomeReq.Opaque = utils.AppendPlainToOpaque(createHomeReq.Opaque, "storage_id", utils.VaultStorageProviderID)
		cpsRes, err := client.CreateHome(ctx, createHomeReq)
		switch {
		case err != nil:
			m.logger.Err(err).Msg("error calling CreateHome for vault personal")
		case cpsRes.GetStatus().GetCode() == rpc.Code_CODE_OK:
			m.logger.Debug().Interface("userID", u.GetId().GetOpaqueId()).Msg("vault personal space created")
		case cpsRes.GetStatus().GetCode() == rpc.Code_CODE_ALREADY_EXISTS:
			m.logger.Info().Interface("userID", u.GetId().GetOpaqueId()).Interface("status", cpsRes.GetStatus()).Msg("=====+ vault personal space already exists")
		default:
			m.logger.Error().Interface("userID", u.GetId().GetOpaqueId()).Interface("status", cpsRes.GetStatus()).Msg("=====! vault personal space creation failed")
		}

	}

	m.next.ServeHTTP(w, req)
}

func (m createHome) shouldServe(req *http.Request) bool {
	return req.Header.Get("x-access-token") != ""
}

func (m createHome) getUserRoles(user *userv1beta1.User) ([]string, error) {
	var roleIDs []string
	if err := utils.ReadJSONFromOpaque(user.Opaque, "roles", &roleIDs); err != nil {
		return nil, err
	}

	tmp := make(map[string]struct{})
	for _, id := range roleIDs {
		tmp[id] = struct{}{}
	}

	dedup := make([]string, 0, len(tmp))
	for k := range tmp {
		dedup = append(dedup, k)
	}
	return dedup, nil
}

func (m createHome) checkRoleQuotaLimit(roleIDs []string) (uint64, bool) {
	if len(roleIDs) == 0 {
		return 0, false
	}
	id := roleIDs[0] // At the moment a user can only have one role.
	quota, ok := m.roleQuotas[id]
	return quota, ok
}

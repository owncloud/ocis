package middleware

import (
	"net/http"
	"strconv"
	"sync"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/metadata"
)

// CreateHome provides a middleware which sends a CreateHome request to the reva gateway.
// It deduplicates concurrent requests for the same user and caches successful results
// for the lifetime of the process.
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

	knownHomes sync.Map           // map[userID]struct{} — users whose home is confirmed
	flight     singleflight.Group // collapses concurrent CreateHome calls per user
}

func (m *createHome) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	token := req.Header.Get("X-Access-Token")

	// we need to pass the token to authenticate the CreateHome request.
	ctx := metadata.AppendToOutgoingContext(req.Context(), revactx.TokenHeader, token)

	u, ok := revactx.ContextGetUser(ctx)
	if !ok {
		m.next.ServeHTTP(w, req)
		return
	}

	userID := u.Id.OpaqueId

	// Fast path: home already known to exist for this user.
	if _, ok := m.knownHomes.Load(userID); ok {
		m.next.ServeHTTP(w, req)
		return
	}

	createHomeReq := &provider.CreateHomeRequest{}
	roleIDs, err := m.getUserRoles(u)
	if err != nil {
		m.logger.Error().Err(err).Str("userid", userID).Msg("failed to get roles for user")
		errorcode.GeneralException.Render(w, req, http.StatusInternalServerError, "Unauthorized")
		return
	}
	if limit, hasLimit := m.checkRoleQuotaLimit(roleIDs); hasLimit {
		createHomeReq.Opaque = utils.AppendPlainToOpaque(nil, "quota", strconv.FormatUint(limit, 10))
	}

	// Deduplicate concurrent CreateHome calls for the same user.
	_, err, _ = m.flight.Do(userID, func() (interface{}, error) {
		client, err := m.revaGatewaySelector.Next()
		if err != nil {
			m.logger.Err(err).Msg("error selecting next gateway client")
			return nil, err
		}

		createHomeRes, err := client.CreateHome(ctx, createHomeReq)
		if err != nil {
			m.logger.Err(err).Msg("error calling CreateHome")
			return nil, err
		}

		if createHomeRes.Status.Code != rpc.Code_CODE_OK && createHomeRes.Status.Code != rpc.Code_CODE_ALREADY_EXISTS {
			err := status.NewErrorFromCode(createHomeRes.Status.Code, "gateway")
			m.logger.Err(err).Msg("error when calling CreateHome")
			return nil, err
		}

		return true, nil
	})

	// Only cache on success — if CreateHome failed, retry on next request.
	if err == nil {
		m.knownHomes.Store(userID, struct{}{})
	}

	m.next.ServeHTTP(w, req)
}

func (m *createHome) shouldServe(req *http.Request) bool {
	return req.Header.Get("X-Access-Token") != ""
}

func (m *createHome) getUserRoles(user *userv1beta1.User) ([]string, error) {
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

func (m *createHome) checkRoleQuotaLimit(roleIDs []string) (uint64, bool) {
	if len(roleIDs) == 0 {
		return 0, false
	}
	id := roleIDs[0] // At the moment a user can only have one role.
	quota, ok := m.roleQuotas[id]
	return quota, ok
}

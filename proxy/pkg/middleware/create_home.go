package middleware

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/status"
	tokenPkg "github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"google.golang.org/grpc/metadata"
	"net/http"
	microErrors "github.com/micro/go-micro/v2/errors"
)

func CreateHome(optionSetters ...Option) func(next http.Handler) http.Handler {
	options := newOptions(optionSetters...)
	logger := options.Logger

	return func(next http.Handler) http.Handler {
		tokenManager, err := jwt.New(map[string]interface{}{
			"secret": options.TokenManagerConfig.JWTSecret,
		})
		if err != nil {
			logger.Fatal().Err(err).Msgf("Could not initialize token-manager")
		}

		return &createHome{
			next:           next,
			logger:         logger,
			accountsClient: options.AccountsClient,
			tokenManager:   tokenManager,
			revaGatewayClient: options.RevaGatewayClient,
		}
	}
}

type createHome struct {
	next              http.Handler
	logger            log.Logger
	accountsClient    accounts.AccountsService
	tokenManager      tokenPkg.Manager
	revaGatewayClient gateway.GatewayAPIClient
}

func (m createHome) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !m.shouldServe(req) {
		m.next.ServeHTTP(w, req)
		return
	}

	token := req.Header.Get("x-access-token")

	user, err := m.tokenManager.DismantleToken(req.Context(), token)
	if err != nil {
		m.logger.Logger.Err(err).Msg("error getting user from access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = m.accountsClient.GetAccount(req.Context(), &accounts.GetAccountRequest{
		Id: user.Id.OpaqueId,
	})

	if err != nil {
		e := microErrors.Parse(err.Error())

		if e.Code == http.StatusNotFound {
			m.logger.Debug().Msgf("account with id %s not found", user.Id.OpaqueId)
			m.next.ServeHTTP(w, req)
			return
		}

		m.logger.Err(err).Msgf("error getting user with id %s from accounts service", user.Id.OpaqueId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// we need to pass the token to authenticate the CreateHome request.
	//ctx := tokenpkg.ContextSetToken(r.Context(), token)
	ctx := metadata.AppendToOutgoingContext(req.Context(), tokenPkg.TokenHeader, token)

	createHomeReq := &provider.CreateHomeRequest{}
	createHomeRes, err := m.revaGatewayClient.CreateHome(ctx, createHomeReq)

	if err != nil {
		m.logger.Err(err).Msg("error calling CreateHome")
	} else if createHomeRes.Status.Code != rpc.Code_CODE_OK {
		err := status.NewErrorFromCode(createHomeRes.Status.Code, "gateway")
		m.logger.Err(err).Msg("error when calling Createhome")
	}

	m.next.ServeHTTP(w, req)
}

func (m createHome) shouldServe(req *http.Request) bool {
	return req.Header.Get("x-access-token") != ""
}

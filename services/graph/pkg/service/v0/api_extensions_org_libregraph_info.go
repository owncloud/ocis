package svc

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"google.golang.org/grpc/metadata"

	"github.com/cs3org/reva/v2/pkg/conversions"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

// TemporaryTokenInfoResponse is a response struct for the TokenInfo method,
// it will go away once its part of the libregraph API
type TemporaryTokenInfoResponse struct {
	ID          string `json:"id"`
	HasPassword bool   `json:"has_password"`
	IsInternal  bool   `json:"is_internal"`
}

// ExtensionsOrgLibregraphInfoProvider is the interface that defines all methods that are needed to provide the information
type ExtensionsOrgLibregraphInfoProvider interface {
	TokenInfo(ctx context.Context, token, password string) (TemporaryTokenInfoResponse, error)
}

// ExtensionsOrgLibregraphInfoService implements the libregraph API
type ExtensionsOrgLibregraphInfoService struct {
	logger          log.Logger
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// NewExtensionsOrgLibregraphInfoService initializes a new libregraph service
func NewExtensionsOrgLibregraphInfoService(logger log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (ExtensionsOrgLibregraphInfoService, error) {
	return ExtensionsOrgLibregraphInfoService{
		logger:          log.Logger{Logger: logger.With().Str("graph api", "ExtensionsOrgLibregraphInfoService").Logger()},
		gatewaySelector: gatewaySelector,
	}, nil
}

// TokenInfo returns information about a token
func (s ExtensionsOrgLibregraphInfoService) TokenInfo(ctx context.Context, token, password string) (TemporaryTokenInfoResponse, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return TemporaryTokenInfoResponse{}, err
	}

	authenticateResponse, err := gatewayClient.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "publicshares",
		ClientId:     token,
		ClientSecret: "password|" + password,
	})
	// could authenticateResponse be a nil?
	if errCode := errorcode.FromCS3Status(authenticateResponse.GetStatus(), err); errCode != nil {
		return TemporaryTokenInfoResponse{}, errCode
	}

	ctx = ctxpkg.ContextSetToken(ctx, authenticateResponse.GetToken())
	ctx = ctxpkg.ContextSetUser(ctx, authenticateResponse.GetUser())
	ctx = metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, authenticateResponse.GetToken())

	getPublicShareResponse, err := gatewayClient.GetPublicShare(
		ctx,
		&link.GetPublicShareRequest{
			Ref: &link.PublicShareReference{
				Spec: &link.PublicShareReference_Token{
					Token: token,
				},
			},
		},
	)
	if errCode := errorcode.FromCS3Status(getPublicShareResponse.GetStatus(), err); errCode != nil {
		return TemporaryTokenInfoResponse{}, errCode
	}

	return TemporaryTokenInfoResponse{
		ID:          storagespace.FormatResourceID(*getPublicShareResponse.Share.GetResourceId()),
		HasPassword: getPublicShareResponse.GetShare().GetPasswordProtected(),
		IsInternal: func() bool {
			return conversions.RoleFromResourcePermissions(getPublicShareResponse.Share.Permissions.GetPermissions(), true).OCSPermissions() == 0
		}(),
	}, nil
}

// ExtensionsOrgLibregraphInfoApi is the API that exposes the extensions libregraph info API
type ExtensionsOrgLibregraphInfoApi struct {
	logger                              log.Logger
	extensionsOrgLibregraphInfoProvider ExtensionsOrgLibregraphInfoProvider
}

// NewExtensionsOrgLibregraphInfoApi initializes a new libregraph info API
func NewExtensionsOrgLibregraphInfoApi(extensionsOrgLibregraphInfoProvider ExtensionsOrgLibregraphInfoProvider, logger log.Logger) (ExtensionsOrgLibregraphInfoApi, error) {
	return ExtensionsOrgLibregraphInfoApi{
		logger:                              log.Logger{Logger: logger.With().Str("graph api", "ExtensionsOrgLibregraphInfoApi").Logger()},
		extensionsOrgLibregraphInfoProvider: extensionsOrgLibregraphInfoProvider,
	}, nil
}

// TokenInfo returns information about a token
func (api ExtensionsOrgLibregraphInfoApi) TokenInfo(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		api.logger.Debug().Msg("no token provided")
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, "no token provided")
		return
	}

	_, pw, _ := r.BasicAuth()
	tokenInfo, err := api.extensionsOrgLibregraphInfoProvider.TokenInfo(r.Context(), token, pw)
	if err != nil {
		msg := "could not get token info"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tokenInfo)
}

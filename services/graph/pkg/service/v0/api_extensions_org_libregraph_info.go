package svc

import (
	"context"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
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
	ID          string `json:"id,omitempty"`
	HasPassword bool   `json:"hasPassword"`
	IsInternal  bool   `json:"isInternal"`
}

// ExtensionsOrgLibregraphInfoProvider is the interface that defines all methods that are needed to provide the information
type ExtensionsOrgLibregraphInfoProvider interface {
	TokenInfo(ctx context.Context, infoToken, password string) (TemporaryTokenInfoResponse, error)
}

// ExtensionsOrgLibregraphInfoServiceOptions defines the required parameters to create a new libregraph info service
type ExtensionsOrgLibregraphInfoServiceOptions struct {
	logger          log.Logger
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// validate checks if the required parameters are set
func (o ExtensionsOrgLibregraphInfoServiceOptions) validate() error {
	return nil
}

// WithLogger sets the logger option
func (o ExtensionsOrgLibregraphInfoServiceOptions) WithLogger(logger log.Logger) ExtensionsOrgLibregraphInfoServiceOptions {
	o.logger = log.Logger{Logger: logger.With().Str("graph api", "ExtensionsOrgLibregraphInfoService").Logger()}
	return o
}

// WithGatewaySelector sets the gatewaySelector option
func (o ExtensionsOrgLibregraphInfoServiceOptions) WithGatewaySelector(gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) ExtensionsOrgLibregraphInfoServiceOptions {
	o.gatewaySelector = gatewaySelector
	return o
}

// ExtensionsOrgLibregraphInfoService implements the libregraph API
type ExtensionsOrgLibregraphInfoService struct {
	logger          log.Logger
	gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
}

// NewExtensionsOrgLibregraphInfoService initializes a new libregraph service
func NewExtensionsOrgLibregraphInfoService(options ExtensionsOrgLibregraphInfoServiceOptions) (ExtensionsOrgLibregraphInfoService, error) {
	return ExtensionsOrgLibregraphInfoService{
		logger:          options.logger,
		gatewaySelector: options.gatewaySelector,
	}, options.validate()
}

// TokenInfo returns information about a token
func (s ExtensionsOrgLibregraphInfoService) TokenInfo(ctx context.Context, infoToken, password string) (TemporaryTokenInfoResponse, error) {
	gatewayClient, err := s.gatewaySelector.Next()
	if err != nil {
		return TemporaryTokenInfoResponse{}, err
	}

	authenticateResponse, err := gatewayClient.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "publicshares",
		ClientId:     infoToken,
		ClientSecret: "password|" + password,
	})
	switch errCode := errorcode.FromCS3Status(authenticateResponse.GetStatus(), err); {
	case errCode == nil:
		break
	// AccessDenied is a special case, we need to return the password-protected status
	case errCode.GetCode() == errorcode.AccessDenied:
		return TemporaryTokenInfoResponse{HasPassword: true}, nil
	default:
		return TemporaryTokenInfoResponse{}, errCode
	}

	createdCTX := ctxpkg.ContextSetToken(context.Background(), authenticateResponse.GetToken())
	createdCTX = ctxpkg.ContextSetUser(createdCTX, authenticateResponse.GetUser())
	createdCTX = metadata.AppendToOutgoingContext(createdCTX, ctxpkg.TokenHeader, authenticateResponse.GetToken())

	getPublicShareResponse, err := gatewayClient.GetPublicShare(
		createdCTX,
		&link.GetPublicShareRequest{
			Ref: &link.PublicShareReference{
				Spec: &link.PublicShareReference_Token{
					Token: infoToken,
				},
			},
		},
	)
	if errCode := errorcode.FromCS3Status(getPublicShareResponse.GetStatus(), err); errCode != nil {
		return TemporaryTokenInfoResponse{}, errCode
	}

	temporaryTokenInfoResponse := TemporaryTokenInfoResponse{
		ID:          storagespace.FormatResourceID(*getPublicShareResponse.Share.GetResourceId()),
		HasPassword: getPublicShareResponse.GetShare().GetPasswordProtected(),
		IsInternal: func() bool {
			return conversions.RoleFromResourcePermissions(getPublicShareResponse.Share.Permissions.GetPermissions(), true).OCSPermissions() == 0
		}(),
	}

	if !temporaryTokenInfoResponse.IsInternal {
		return temporaryTokenInfoResponse, nil
	}

	statResponse, err := gatewayClient.Stat(
		ctx,
		&storageprovider.StatRequest{
			Ref: &storageprovider.Reference{
				ResourceId: getPublicShareResponse.GetShare().GetResourceId(),
			},
		},
	)
	switch errCode := errorcode.FromCS3Status(statResponse.GetStatus(), err); {
	case errCode == nil:
		break
	case errCode.GetCode() == errorcode.ItemNotFound:
		temporaryTokenInfoResponse.ID = ""
		return temporaryTokenInfoResponse, nil
	default:
		return TemporaryTokenInfoResponse{}, errCode
	}

	return temporaryTokenInfoResponse, err
}

// ExtensionsOrgLibregraphInfoApiOptions defines the required parameters to create a new libregraph info API
type ExtensionsOrgLibregraphInfoApiOptions struct {
	logger                              log.Logger
	extensionsOrgLibregraphInfoProvider ExtensionsOrgLibregraphInfoProvider
}

// validate checks if the required parameters are set
func (o ExtensionsOrgLibregraphInfoApiOptions) validate() error {
	return nil
}

// WithLogger sets the logger option
func (o ExtensionsOrgLibregraphInfoApiOptions) WithLogger(logger log.Logger) ExtensionsOrgLibregraphInfoApiOptions {
	o.logger = log.Logger{Logger: logger.With().Str("graph api", "ExtensionsOrgLibregraphInfoApi").Logger()}
	return o
}

// WithExtensionsOrgLibregraphInfoProvider sets the extensionsOrgLibregraphInfoProvider option
func (o ExtensionsOrgLibregraphInfoApiOptions) WithExtensionsOrgLibregraphInfoProvider(extensionsOrgLibregraphInfoProvider ExtensionsOrgLibregraphInfoProvider) ExtensionsOrgLibregraphInfoApiOptions {
	o.extensionsOrgLibregraphInfoProvider = extensionsOrgLibregraphInfoProvider
	return o
}

// ExtensionsOrgLibregraphInfoApi is the API that exposes the extensions libregraph info API
type ExtensionsOrgLibregraphInfoApi struct {
	logger                              log.Logger
	extensionsOrgLibregraphInfoProvider ExtensionsOrgLibregraphInfoProvider
}

// NewExtensionsOrgLibregraphInfoApi initializes a new libregraph info API
func NewExtensionsOrgLibregraphInfoApi(options ExtensionsOrgLibregraphInfoApiOptions) (ExtensionsOrgLibregraphInfoApi, error) {
	return ExtensionsOrgLibregraphInfoApi{
		logger:                              options.logger,
		extensionsOrgLibregraphInfoProvider: options.extensionsOrgLibregraphInfoProvider,
	}, options.validate()
}

// TokenInfo returns information about a token
func (api ExtensionsOrgLibregraphInfoApi) TokenInfo(w http.ResponseWriter, r *http.Request) {
	infoToken := chi.URLParam(r, "token")
	if infoToken == "" {
		api.logger.Debug().Msg("no token provided")
		errorcode.InvalidRequest.Render(w, r, http.StatusUnprocessableEntity, "no token provided")
		return
	}

	_, pw, _ := r.BasicAuth()
	tokenInfo, err := api.extensionsOrgLibregraphInfoProvider.TokenInfo(r.Context(), infoToken, pw)
	if err != nil {
		msg := "could not get token info"
		api.logger.Debug().Err(err).Msg(msg)
		errorcode.InvalidRequest.Render(w, r, http.StatusFailedDependency, msg)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tokenInfo)
}

package helpers

import (
	"context"
	"errors"

	registryv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/mime"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
)

// RegisterOcisService will register this service.
// There are no explicit requirements for the context, and it will be passed
// without changes to the underlying RegisterService method.
func RegisterOcisService(ctx context.Context, cfg *config.Config, logger log.Logger) error {
	svc := registry.BuildGRPCService(cfg.GRPC.Namespace+"."+cfg.Service.Name+"."+cfg.App.Name, cfg.GRPC.Protocol, cfg.GRPC.Addr, version.GetString())
	return registry.RegisterService(ctx, logger, svc, cfg.Debug.Addr)
}

// RegisterAppProvider will register this service as app provider in REVA.
// The GatewayAPIClient is expected to be provided via `helpers.GetCS3apiClient`.
// The appUrls are expected to be provided via `helpers.GetAppURLs`
//
// Note that this method doesn't provide a re-registration mechanism, so it
// will register the service once
func RegisterAppProvider(
	ctx context.Context,
	cfg *config.Config,
	logger log.Logger,
	gws pool.Selectable[gatewayv1beta1.GatewayAPIClient],
	appUrls map[string]map[string]string,
) error {
	mimeTypesMap := make(map[string]bool)
	for _, extensions := range appUrls {
		for ext := range extensions {
			m := mime.Detect(false, ext)
			// skip the default
			if m == "application/octet-stream" {
				continue
			}
			mimeTypesMap[m] = true
		}
	}

	mimeTypes := make([]string, 0, len(mimeTypesMap))
	for m := range mimeTypesMap {
		mimeTypes = append(mimeTypes, m)
	}

	logger.Debug().
		Str("AppName", cfg.App.Name).
		Strs("Mimetypes", mimeTypes).
		Msg("Registering mimetypes in the app provider")

	// TODO: an added app provider shouldn't last forever. Instead the registry should use a TTL
	// and delete providers that didn't register again. If an app provider dies or get's disconnected,
	// the users will be no longer available to choose to open a file with it (currently, opening a file just fails)
	req := &registryv1beta1.AddAppProviderRequest{
		Provider: &registryv1beta1.ProviderInfo{
			Name:        cfg.App.Name,
			Description: cfg.App.Description,
			Icon:        cfg.App.Icon,
			Address:     cfg.GRPC.Namespace + "." + cfg.Service.Name + "." + cfg.App.Name,
			MimeTypes:   mimeTypes,
			ProductName: cfg.App.Product,
		},
	}
	gwc, err := gws.Next()
	if err != nil {
		return err
	}
	resp, err := gwc.AddAppProvider(ctx, req)
	if err != nil {
		logger.Error().Err(err).Msg("AddAppProvider failed")
		return err
	}

	if resp.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		logger.Error().Str("status_code", resp.GetStatus().GetCode().String()).Msg("AddAppProvider failed")
		return errors.New("status code != CODE_OK")
	}

	return nil
}

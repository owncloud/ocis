package app

import (
	"context"
	"errors"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/internal/logging"

	registryv1beta1 "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"google.golang.org/grpc"
)

type DemoApp struct {
	gwc        gatewayv1beta1.GatewayAPIClient
	grpcServer *grpc.Server

	AppURLs map[string]map[string]string

	Config *config.Config

	Logger log.Logger
}

func New(cfg *config.Config) (*DemoApp, error) {
	app := &DemoApp{
		Config: cfg,
	}

	err := envdecode.Decode(app)
	if err != nil {
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return nil, err
		}
	}

	app.Logger = logging.Configure("wopiserver", defaults.FullDefaultConfig().Log)

	return app, nil
}

func (app *DemoApp) GetCS3apiClient() error {
	// establish a connection to the cs3 api endpoint
	// in this case a REVA gateway, started by oCIS
	gwc, err := pool.GetGatewayServiceClient(app.Config.CS3Api.Gateway.Name)
	if err != nil {
		return err
	}
	app.gwc = gwc

	return nil
}

func (app *DemoApp) RegisterOcisService(ctx context.Context) error {
	svc := registry.BuildGRPCService(app.Config.Service.Name, uuid.Must(uuid.NewV4()).String(), app.Config.GRPC.Addr, "0.0.0")
	return registry.RegisterService(ctx, svc, app.Logger)
}

func (app *DemoApp) RegisterDemoApp(ctx context.Context) error {
	mimeTypesMap := make(map[string]bool)
	for _, extensions := range app.AppURLs {
		for ext := range extensions {
			m := mime.Detect(false, ext)
			mimeTypesMap[m] = true
		}
	}

	mimeTypes := make([]string, 0, len(mimeTypesMap))
	for m := range mimeTypesMap {
		mimeTypes = append(mimeTypes, m)
	}

	// TODO: REVA has way to filter supported mimetypes (do we need to implement it here or is it in the registry?)

	// TODO: an added app provider shouldn't last forever. Instead the registry should use a TTL
	// and delete providers that didn't register again. If an app provider dies or get's disconnected,
	// the users will be no longer available to choose to open a file with it (currently, opening a file just fails)
	req := &registryv1beta1.AddAppProviderRequest{
		Provider: &registryv1beta1.ProviderInfo{
			Name:        app.Config.App.Name,
			Description: app.Config.App.Description,
			Icon:        app.Config.App.Icon,
			Address:     app.Config.Service.Name,
			MimeTypes:   mimeTypes,
		},
	}

	resp, err := app.gwc.AddAppProvider(ctx, req)
	if err != nil {
		app.Logger.Error().Err(err).Msg("AddAppProvider failed")
		return err
	}

	if resp.Status.Code != rpcv1beta1.Code_CODE_OK {
		app.Logger.Error().Str("status_code", resp.Status.Code.String()).Msg("AddAppProvider failed")
		return errors.New("status code != CODE_OK")
	}

	return nil
}

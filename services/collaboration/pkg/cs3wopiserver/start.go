package cs3wopiserver

import (
	"context"

	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/internal/app"
)

func Start(cfg *config.Config) (*app.DemoApp, error) {
	ctx := context.Background()

	app, err := app.New(cfg)
	if err != nil {
		return nil, err
	}

	if err := app.RegisterOcisService(ctx); err != nil {
		return nil, err
	}

	if err := app.WopiDiscovery(ctx); err != nil {
		return nil, err
	}

	if err := app.GetCS3apiClient(); err != nil {
		return nil, err
	}

	if err := app.RegisterDemoApp(ctx); err != nil {
		return nil, err
	}

	// NOTE:
	// GRPC and HTTP server are started using the standard
	// `ocis collaboration server` command through the usual means

	// TODO:
	// "app" initialization needs to be moved

	return app, nil
}

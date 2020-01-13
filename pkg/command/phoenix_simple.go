// +build simple

package command

import (
	"os"

	svcconfig "github.com/owncloud/ocis-phoenix/pkg/config"
	"github.com/owncloud/ocis/pkg/config"
)

func configurePhoenix(cfg *config.Config) *svcconfig.Config {
	cfg.Phoenix.Log.Level = cfg.Log.Level
	cfg.Phoenix.Log.Pretty = cfg.Log.Pretty
	cfg.Phoenix.Log.Color = cfg.Log.Color
	
	// TODO they will be overriden in the extensions service.Init(), tracked in https://github.com/owncloud/ocis/issues/81
	//cfg.Phoenix.Phoenix.Config.OpenIDConnect.Authority = "http://localhost:9135"
	//cfg.Phoenix.Phoenix.Config.OpenIDConnect.MetadataURL = "http://localhost:9135/.well-known/openid-configuration"

    if len(os.Getenv("PHOENIX_OIDC_METADATA_URL")) == 0 {
        os.Setenv("PHOENIX_OIDC_METADATA_URL", "http://localhost:9135/.well-known/openid-configuration")
    }
    if len(os.Getenv("PHOENIX_OIDC_AUTHORITY")) == 0 {
		os.Setenv("PHOENIX_OIDC_AUTHORITY", "http://localhost:9135")
	}

	// disable built in apps
	cfg.Phoenix.Phoenix.Config.Apps = []string{}
	// enable ocis-hello extension
	cfg.Phoenix.Phoenix.Config.ExternalApps = []svcconfig.ExternalApp{
		svcconfig.ExternalApp{
			ID: "hello",
			Path: "http://localhost:9105/hello.js",
			Config: map[string]interface{}{
				"url": "http://localhost:9105",
			},
		},
	}

	return cfg.Phoenix
}

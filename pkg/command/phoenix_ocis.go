// +build !simple

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

	if len(os.Getenv("PHOENIX_OIDC_METADATA_URL")) == 0 {
		os.Setenv("PHOENIX_OIDC_METADATA_URL", "http://localhost:20080/.well-known/openid-configuration")
	}
	if len(os.Getenv("PHOENIX_OIDC_AUTHORITY")) == 0 {
		os.Setenv("PHOENIX_OIDC_AUTHORITY", "http://localhost:20080")
	}
	if len(os.Getenv("PHOENIX_WEB_CONFIG_SERVER")) == 0 {
		os.Setenv("PHOENIX_WEB_CONFIG_SERVER", "http://localhost:20080")
	}

	return cfg.Phoenix
}

package command

import (
	"github.com/spf13/viper"
)

// init defined the default options for viper.
func init() {
	viper.SetDefault("debug.addr", "0.0.0.0:8090")
	viper.SetDefault("debug.token", "")
	viper.SetDefault("debug.pprof", false)

	viper.SetDefault("http.addr", "0.0.0.0:8080")
	viper.SetDefault("http.root", "/")

	viper.SetDefault("asset.path", "")

	viper.SetDefault("config.custom", "")
	viper.SetDefault("config.server", "")
	viper.SetDefault("config.theme", "owncloud")
	viper.SetDefault("config.version", "0.1.0")
	viper.SetDefault("config.client", "")
	viper.SetDefault("config.apps", []string{"files"})
}

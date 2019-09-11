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
	viper.SetDefault("config.file", "")
}

package command

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("server.addr", "0.0.0.0:8080")
	viper.SetDefault("metrics.addr", "0.0.0.0:8090")
}

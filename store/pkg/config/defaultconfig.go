package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9464",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9460",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "store",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "store"),
	}
}

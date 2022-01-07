package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "settings",
		},
		Debug: Debug{
			Addr:   "127.0.0.1:9194",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9190",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800, // 7 days
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9191",
			Namespace: "com.owncloud.api",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
		},
		DataPath: path.Join(defaults.BaseDataPath(), "settings"),
		Asset: Asset{
			Path: "",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
	}
}

package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9182",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9181",
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
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "accounts",
		},
		Asset: Asset{},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		HashDifficulty:     11,
		DemoUsersAndGroups: true,
		Repo: Repo{
			Backend: "CS3",
			Disk: Disk{
				Path: path.Join(defaults.BaseDataPath(), "accounts"),
			},
			CS3: CS3{
				ProviderAddr: "localhost:9215",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
		Index: Index{
			UID: UIDBound{
				Lower: 0,
				Upper: 1000,
			},
			GID: GIDBound{
				Lower: 0,
				Upper: 1000,
			},
		},
		ServiceUser: ServiceUser{
			UUID:     "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			Username: "",
			UID:      0,
			GID:      0,
		},
	}
}

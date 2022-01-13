package config

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9114",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9110",
			Root:      "/ocs",
			Namespace: "com.owncloud.web",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Service: Service{
			Name: "ocs",
		},

		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		AccountBackend:     "accounts",
		Reva:               Reva{Address: "127.0.0.1:9142"},
		StorageUsersDriver: "ocis",
		MachineAuthAPIKey:  "change-me-please",
		IdentityManagement: IdentityManagement{
			Address: "https://localhost:9200",
		},
	}
}

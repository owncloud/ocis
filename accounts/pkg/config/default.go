package config

// DefaultConfig are values stored in the flagset, but moved to a struct.
func DefaultConfig() Config {
	return Config{
		LDAP: LDAP{},
		HTTP: HTTP{
			Addr:      "0.0.0.0:9181",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800,
		},
		GRPC: GRPC{
			Addr:      "0.0.0.0:9180",
			Namespace: "com.owncloud.api",
		},
		Server: Server{
			Name:           "accounts",
			HashDifficulty: 11,
		},
		Asset: Asset{
			Path: "",
		},
		Log: Log{},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Repo: Repo{
			Disk: Disk{
				Path: "",
			},
			CS3: CS3{
				ProviderAddr: "localhost:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
		Index: Index{
			UID: Bound{
				Lower: 0,
				Upper: 1000,
			},
			GID: Bound{
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
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "accounts",
		},
		Context:    nil,
		Supervised: false,
	}
}

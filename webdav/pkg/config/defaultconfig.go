package config

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9119",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9115",
			Root:      "/",
			Namespace: "com.owncloud.web",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Service: Service{
			Name: "webdav",
		},
		OcisPublicURL:   "https://127.0.0.1:9200",
		WebdavNamespace: "/users/{{.Id.OpaqueId}}",
		RevaGateway:     "127.0.0.1:9142",
	}
}

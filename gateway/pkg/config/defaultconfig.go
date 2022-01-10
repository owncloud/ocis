package config

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "127.0.0.1:9143",
			Token: "",
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9142",
			Namespace: "com.owncloud.api",
			Network:   "tcp",
		},
		Service: Service{
			Name: "gateway",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
		},
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		StorageRegistry: StorageRegistry{
			Driver:       "static",
			HomeProvider: "/home",
			Storages: Storages{
				StorageHome: StorageHome{
					MountPath:     "1284d238-aa92-42ce-bdc4-0b0000009154",
					AlternativeID: "/home",
				},
				StorageUsers: StorageUsers{
					MountPath: "/users",
					MountID:   "1284d238-aa92-42ce-bdc4-0b0000009157",
				},
				StoragePublicShare: StoragePublicShare{
					MountPath: "",
					MountID:   "",
				},
			},
		},
		ServiceMap: ServiceMap{
			AuthRegistryAddr:    "127.0.0.1:9142",
			StorageRegistryAddr: "127.0.0.1:9142",
			AppRegistryAddr:     "127.0.0.1:9142",

			PreferenceAddr:    "127.0.0.1:9144",
			UserProviderAddr:  "127.0.0.1:9144",
			GroupProviderAddr: "127.0.0.1:9144",

			UserShareProviderAddr:   "127.0.0.1:9150",
			PublicShareProviderAddr: "127.0.0.1:9150",
			OCMShareProviderAddr:    "127.0.0.1:9150",

			AuthBasicAddr:        "127.0.0.1:9146",
			AuthBearerAddr:       "127.0.0.1:9166",
			AuthMachineAddr:      "127.0.0.1:9148",
			AuthPublicSharesAddr: "127.0.0.1:9178",
		},
	}
}

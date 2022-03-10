package config

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "audit",
		},
		Events: Events{
			Endpoint:      "127.0.0.1:9233",
			Cluster:       "test-cluster",
			ConsumerGroup: "audit",
		},
		Auditlog: Auditlog{
			LogToConsole: true,
			Format:       "json",
		},
	}
}

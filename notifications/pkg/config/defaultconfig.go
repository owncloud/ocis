package config

// NOTE: Most of this configuration is not needed to keep it as simple as possible
// TODO: Clean up unneeded configuration

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "notifications",
		},
		Notifications: Notifications{
			SMTP: SMTP{
				Host:     "127.0.0.1",
				Port:     "1025",
				Sender:   "god@example.com",
				Password: "godisdead",
			},
			Events: Events{
				Endpoint:      "127.0.0.1:9233",
				Cluster:       "test-cluster",
				ConsumerGroup: "notifications",
			},
			RevaGateway:       "127.0.0.1:9142",
			MachineAuthSecret: "change-me-please",
		},
	}
}

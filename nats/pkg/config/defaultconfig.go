package config

// NOTE: Most of this configuration is not needed to keep it as simple as possible
// TODO: Clean up unneeded configuration

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "nats",
		},
		Nats: Nats{
			Host: "127.0.0.1",
			Port: 4222,
		},
	}
}

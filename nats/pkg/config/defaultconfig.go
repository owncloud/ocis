package config

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "nats",
		},
	}
}

package config

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `yaml:"address" env:"REVA_GATEWAY" desc:"The CS3 gateway endpoint"`
}

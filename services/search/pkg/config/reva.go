package config

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `yaml:"address" env:"OCIS_REVA_GATEWAY" desc:"The CS3 gateway endpoint." introductionVersion:"pre5.0"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;SEARCH_JWT_SECRET" desc:"The secret to mint and validate jwt tokens." introductionVersion:"pre5.0"`
}

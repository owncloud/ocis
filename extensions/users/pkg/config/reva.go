package config

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `yaml:"address" env:"REVA_GATEWAY" desc:"The CS3 gateway endpoint"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;USERS_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

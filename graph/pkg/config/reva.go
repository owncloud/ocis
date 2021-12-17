package config

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address" env:"REVA_GATEWAY"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;OCS_JWT_SECRET"`
}

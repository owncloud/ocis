package config

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;SETTINGS_JWT_SECRET"`
}

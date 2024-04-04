package config

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;OCS_JWT_SECRET" desc:"The secret to mint and validate jwt tokens." introductionVersion:"pre5.0"`
}

type API struct {
	ShowUserEmailInResults bool `yaml:"show_email_in_results" env:"OCIS_SHOW_USER_EMAIL_IN_RESULTS" desc:"Mask user email addresses in responses." introductionVersion:"5.1"`
}

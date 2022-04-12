package config

// Service defines the available service configuration.
type Service struct {
	Name             string `yaml:"-"`
	PasswordResetURI string `yaml:"password_reset_uri" env:"IDP_PASSWORD_RESET_URI" desc:"The URI where a user can reset their password."`
}

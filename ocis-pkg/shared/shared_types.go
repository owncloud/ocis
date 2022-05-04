package shared

// EnvBinding represents a direct binding from an env variable to a go kind. Along with gookit/config, its primal goal
// is to unpack environment variables into a Go value. We do so with reflection, and this data structure is just a step
// in between.
type EnvBinding struct {
	EnvVars     []string    // name of the environment var.
	Destination interface{} // pointer to the original config value to modify.
}

// Log defines the available logging configuration.
type Log struct {
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL"`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY"`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR"`
	File   string `yaml:"file" env:"OCIS_LOG_FILE"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET" desc:"The secret to mint jwt tokens."`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `yaml:"address" env:"REVA_GATEWAY"`
}

// Commons holds configuration that are common to all extensions. Each extension can then decide whether
// to overwrite its values.
type Commons struct {
	Log               *Log          `yaml:"log"`
	Tracing           *Tracing      `yaml:"tracing"`
	OcisURL           string        `yaml:"ocis_url" env:"OCIS_URL"`
	TokenManager      *TokenManager `yaml:"token_manager"`
	Reva              *Reva         `yaml:"reva"`
	MachineAuthAPIKey string        `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
	TransferSecret    string        `yaml:"transfer_secret,omitempty" env:"REVA_TRANSFER_SECRET"`
	MetadataUserID    string        `yaml:"metadata_user_id" env:"METADATA_USER_ID"`
	SystemUserAPIKey  string        `yaml:"system_user_api_key" env:"SYSTEM_USER_API_KEY"`
}

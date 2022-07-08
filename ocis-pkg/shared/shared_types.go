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
	Level  string `yaml:"level" env:"OCIS_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `yaml:"pretty" env:"OCIS_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `yaml:"color" env:"OCIS_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `yaml:"file" env:"OCIS_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE" desc:"The type of tracing. Defaults to \"\", which is the same as \"jaeger\". Allowed tracing types are \"jaeger\" and \"\" as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mask:"password" yaml:"jwt_secret" env:"OCIS_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `yaml:"address" env:"REVA_GATEWAY" desc:"The CS3 gateway endpoint."`
}

// Commons holds configuration that are common to all extensions. Each extension can then decide whether
// to overwrite its values.
type Commons struct {
	Log               *Log          `yaml:"log"`
	Tracing           *Tracing      `yaml:"tracing"`
	OcisURL           string        `yaml:"ocis_url" env:"OCIS_URL" desc:"URL, where oCIS is reachable for users."`
	TokenManager      *TokenManager `mask:"struct" yaml:"token_manager"`
	Reva              *Reva         `yaml:"reva"`
	MachineAuthAPIKey string        `mask:"password" yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`
	TransferSecret    string        `mask:"password" yaml:"transfer_secret,omitempty" env:"REVA_TRANSFER_SECRET"`
	SystemUserID      string        `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID" desc:"ID of the oCIS storage-system system user. Admins need to set the ID for the storage-system system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
	SystemUserAPIKey  string        `mask:"password" yaml:"system_user_api_key" env:"SYSTEM_USER_API_KEY"`
	AdminUserID       string        `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID" desc:"ID of a user, that should receive admin privileges."`
}

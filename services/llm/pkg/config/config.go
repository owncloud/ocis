package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	TokenManager *TokenManager `yaml:"token_manager"`

	Store     Store     `yaml:"store"`
	RateLimit RateLimit `yaml:"rate_limit"`
	LLM       LLM       `yaml:"llm"`

	Context context.Context `yaml:"-"`
}

// TokenManager is the config for using the reva token manager.
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;LLM_JWT_SECRET" desc:"The secret to mint and validate jwt tokens." introductionVersion:"8.3.0"`
}

// Store configures the store used for per-user rate limiting.
type Store struct {
	Store                string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;LLM_STORE" desc:"The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details." introductionVersion:"8.3.0"`
	Nodes                []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;LLM_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"8.3.0"`
	Database             string        `yaml:"database" env:"LLM_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"8.3.0"`
	Table                string        `yaml:"table" env:"LLM_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"8.3.0"`
	TTL                  time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;LLM_STORE_TTL" desc:"Time to live for rate-limit entries in the store. See the Environment Variable Types description for more details." introductionVersion:"8.3.0"`
	AuthUsername         string        `yaml:"username" env:"OCIS_PERSISTENT_STORE_AUTH_USERNAME;LLM_STORE_AUTH_USERNAME" desc:"The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	AuthPassword         string        `yaml:"password" env:"OCIS_PERSISTENT_STORE_AUTH_PASSWORD;LLM_STORE_AUTH_PASSWORD" desc:"The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	EnableTLS            bool          `yaml:"enable_tls" env:"OCIS_PERSISTENT_STORE_ENABLE_TLS;LLM_STORE_ENABLE_TLS" desc:"Activate TLS for the connection to the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	TLSInsecure          bool          `yaml:"tls_insecure" env:"OCIS_PERSISTENT_STORE_TLS_INSECURE;LLM_STORE_TLS_INSECURE" desc:"Disable TLS certificate verification for the store connection. Only applies when store type 'nats-js-kv' is configured. Do not enable this in production because it disables authentication of the NATS server. Use it only for testing with self-signed certificates." introductionVersion:"8.3.0"`
	TLSRootCACertificate string        `yaml:"tls_root_ca_certificate" env:"OCIS_PERSISTENT_STORE_TLS_ROOT_CA_CERTIFICATE;LLM_STORE_TLS_ROOT_CA_CERTIFICATE" desc:"Path to the PEM-encoded root CA certificate for the store TLS connection. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
}

// RateLimit configures the per-user request rate limit.
type RateLimit struct {
	Window      time.Duration `yaml:"window" env:"LLM_RATE_LIMIT_WINDOW" desc:"The sliding time window over which requests per user are counted. See the Environment Variable Types description for more details." introductionVersion:"8.3.0"`
	MaxRequests int           `yaml:"max_requests" env:"LLM_RATE_LIMIT_MAX_REQUESTS" desc:"The maximum number of chat completion requests a single user may make within the configured window. Set to 0 to disable rate limiting." introductionVersion:"8.3.0"`
}

// LLM configures the upstream OpenAI-compatible LLM endpoint.
type LLM struct {
	Endpoint       string        `yaml:"endpoint" env:"LLM_ENDPOINT" desc:"The base URL of the OpenAI-compatible LLM endpoint, without the trailing '/chat/completions' path, e.g. 'https://api.openai.com/v1'." introductionVersion:"8.3.0"`
	APIKey         string        `yaml:"api_key" env:"LLM_API_KEY" desc:"The API key used to authenticate with the LLM endpoint. Leave empty for LLM endpoints that require no authentication." introductionVersion:"8.3.0" mask:"password"`
	Model          string        `yaml:"model" env:"LLM_MODEL" desc:"When set, overrides the 'model' field of every client request, regardless of what the client sent." introductionVersion:"8.3.0"`
	MaxTokensLimit int           `yaml:"max_tokens_limit" env:"LLM_MAX_TOKENS_LIMIT" desc:"Hard ceiling on the 'max_tokens' value forwarded to the LLM. A client-requested value above this is clamped to it." introductionVersion:"8.3.0"`
	Timeout        time.Duration `yaml:"timeout" env:"LLM_TIMEOUT" desc:"Timeout for the upstream LLM request. See the Environment Variable Types description for more details." introductionVersion:"8.3.0"`
	MaxBodyBytes   int64         `yaml:"max_body_bytes" env:"LLM_MAX_BODY_BYTES" desc:"Maximum request body size in bytes the service will buffer from a client before rejecting the request." introductionVersion:"8.3.0"`
}

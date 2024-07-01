package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"go-micro.dev/v4/client"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP           HTTPConfig            `yaml:"http"`
	Middleware     Middleware            `yaml:"middleware"`
	GRPC           GRPCConfig            `yaml:"grpc"`
	GRPCClientTLS  *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	GrpcClient     client.Client         `yaml:"-"`
	ServiceAccount ServiceAccount        `yaml:"service_account"`

	Reva                         *shared.Reva                 `yaml:"reva"`
	OCMD                         OCMD                         `yaml:"ocmd"`
	ScienceMesh                  ScienceMesh                  `yaml:"sciencemesh"`
	OCMInviteManager             OCMInviteManager             `yaml:"ocm_invite_manager"`
	OCMProviderAuthorizerDriver  string                       `yaml:"ocm_provider_authorizer_driver" env:"SHARING_OCM_PROVIDER_AUTHORIZER_DRIVER" desc:"Driver to be used to persist ocm invites. Supported value is only 'json'." introductionVersion:"5.0"`
	OCMProviderAuthorizerDrivers OCMProviderAuthorizerDrivers `yaml:"ocm_provider_authorizer_drivers"`
	OCMShareProvider             OCMShareProvider             `yaml:"ocm_share_provider"`
	OCMCore                      OCMCore                      `yaml:"ocm_core"`
	OCMStorageProvider           OCMStorageProvider           `yaml:"ocm_storage_provider"`

	Supervised bool            `yaml:"-"`
	Context    context.Context `yaml:"-"`
}

// HTTPConfig defines the available http configuration.
type HTTPConfig struct {
	Addr      string `yaml:"addr" env:"OCM_HTTP_ADDR" desc:"The bind address of the HTTP service." introductionVersion:"5.0"`
	Namespace string `yaml:"-"`
	Protocol  string `yaml:"protocol" env:"OCM_HTTP_PROTOCOL" desc:"The transport protocol of the HTTP service." introductionVersion:"5.0"`
	Prefix    string `yaml:"prefix" env:"OCM_HTTP_PREFIX" desc:"The path prefix where OCM can be accessed (defaults to /)." introductionVersion:"5.0"`
	CORS      CORS   `yaml:"cors"`
}

// Middleware configures reva middlewares.
type Middleware struct {
	Auth Auth `yaml:"auth"`
}

// Auth configures reva http auth middleware.
type Auth struct {
	CredentialsByUserAgent map[string]string `yaml:"credentials_by_user_agent"`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ID     string `yaml:"service_account_id" env:"OCIS_SERVICE_ACCOUNT_ID;OCM_SERVICE_ACCOUNT_ID" desc:"The ID of the service account the service should use. See the 'auth-service' service description for more details." introductionVersion:"5.0"`
	Secret string `yaml:"service_account_secret" env:"OCIS_SERVICE_ACCOUNT_SECRET;OCM_SERVICE_ACCOUNT_SECRET" desc:"The service account secret." introductionVersion:"5.0"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;OCM_CORS_ALLOW_ORIGINS" desc:"A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details." introductionVersion:"5.0"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;OCM_CORS_ALLOW_METHODS" desc:"A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details." introductionVersion:"5.0"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;OCM_CORS_ALLOW_HEADERS" desc:"A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details." introductionVersion:"5.0"`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;OCM_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials." introductionVersion:"5.0"`
}

// GRPCConfig defines the available grpc configuration.
type GRPCConfig struct {
	Addr      string                 `ocisConfig:"addr" env:"OCM_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"5.0"`
	Namespace string                 `ocisConfig:"-" yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
	Protocol  string                 `yaml:"protocol" env:"OCM_GRPC_PROTOCOL" desc:"The transport protocol of the GRPC service." introductionVersion:"5.0"`
}

type ScienceMesh struct {
	Prefix           string `yaml:"prefix" env:"OCM_SCIENCEMESH_PREFIX" desc:"URL path prefix for the ScienceMesh service. Note that the string must not start with '/'." introductionVersion:"5.0"`
	MeshDirectoryURL string `yaml:"science_mesh_directory_url" env:"OCM_MESH_DIRECTORY_URL" desc:"URL of the mesh directory service." introductionVersion:"5.0"`
}

type OCMD struct {
	Prefix                     string `yaml:"prefix" env:"OCM_OCMD_PREFIX" desc:"URL path prefix for the OCMD service. Note that the string must not start with '/'." introductionVersion:"5.0"`
	ExposeRecipientDisplayName bool   `yaml:"expose_recipient_display_name" env:"OCM_OCMD_EXPOSE_RECIPIENT_DISPLAY_NAME" desc:"Expose the display name of OCM share recipients." introductionVersion:"5.0"`
}

type OCMInviteManager struct {
	Driver          string                  `yaml:"driver" env:"OCM_OCM_INVITE_MANAGER_DRIVER" desc:"Driver to be used to persist OCM invites. Supported value is only 'json'." introductionVersion:"5.0"`
	Drivers         OCMInviteManagerDrivers `yaml:"drivers"`
	TokenExpiration time.Duration           `yaml:"token_expiration" env:"OCM_OCM_INVITE_MANAGER_TOKEN_EXPIRATION" desc:"Expiry duration for invite tokens." introductionVersion:"6.0.1"`
	Timeout         time.Duration           `yaml:"timeout" env:"OCM_OCM_INVITE_MANAGER_TIMEOUT" desc:"Timeout specifies a time limit for requests made to OCM endpoints." introductionVersion:"6.0.1"`
	Insecure        bool                    `yaml:"insecure" env:"OCM_OCM_INVITE_MANAGER_INSECURE" desc:"Disable TLS certificate validation for the OCM connections. Do not set this in production environments." introductionVersion:"5.0"`
}

type OCMInviteManagerDrivers struct {
	JSON OCMInviteManagerJSONDriver `yaml:"json"`
}

type OCMInviteManagerJSONDriver struct {
	File string `yaml:"file" env:"OCM_OCM_INVITE_MANAGER_JSON_FILE" desc:"Path to the JSON file where OCM invite data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/storage." introductionVersion:"5.0"`
}

type OCMProviderAuthorizerDrivers struct {
	JSON OCMProviderAuthorizerJSONDriver `yaml:"json"`
}

type OCMProviderAuthorizerJSONDriver struct {
	Providers             string `yaml:"providers" env:"OCM_OCM_PROVIDER_AUTHORIZER_PROVIDERS_FILE" desc:"Path to the JSON file where ocm invite data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/storage." introductionVersion:"5.0"`
	VerifyRequestHostname bool   `yaml:"verify_request_hostname" env:"OCM_OCM_PROVIDER_AUTHORIZER_VERIFY_REQUEST_HOSTNAME" desc:"Verify the hostname of the incoming request against the hostname of the OCM provider." introductionVersion:"5.0"`
}

type OCMCore struct {
	Driver  string         `yaml:"driver" env:"OCM_OCM_CORE_DRIVER" desc:"Driver to be used for the OCM core. Supported value is only 'json'." introductionVersion:"5.0"`
	Drivers OCMCoreDrivers `yaml:"drivers"`
}
type OCMStorageProvider struct {
	Insecure    bool   `yaml:"insecure" env:"OCM_OCM_STORAGE_PROVIDER_INSECURE" desc:"Disable TLS certificate validation for the OCM connections. Do not set this in production environments." introductionVersion:"5.0"`
	StorageRoot string `yaml:"storage_root" env:"OCM_OCM_STORAGE_PROVIDER_STORAGE_ROOT" desc:"Directory where the ocm storage provider persists its data like tus upload info files." introductionVersion:"5.0"`
}

type OCMCoreDrivers struct {
	JSON OCMCoreJSONDriver `yaml:"json"`
}

type OCMCoreJSONDriver struct {
	File string `yaml:"file" env:"OCM_OCM_CORE_JSON_FILE" desc:"Path to the JSON file where OCM share data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/storage." introductionVersion:"5.0"`
}

type OCMShareProvider struct {
	Driver         string                  `yaml:"driver" env:"OCM_OCM_SHARE_PROVIDER_DRIVER" desc:"Driver to be used for the OCM share provider. Supported value is only 'json'." introductionVersion:"5.0"`
	Drivers        OCMShareProviderDrivers `yaml:"drivers"`
	Insecure       bool                    `yaml:"insecure" env:"OCM_OCM_SHARE_PROVIDER_INSECURE" desc:"Disable TLS certificate validation for the OCM connections. Do not set this in production environments." introductionVersion:"5.0"`
	WebappTemplate string                  `yaml:"webapp_template" env:"OCM_WEBAPP_TEMPLATE" desc:"Template for the webapp url." introductionVersion:"5.0"`
}

type OCMShareProviderDrivers struct {
	JSON OCMShareProviderJSONDriver `yaml:"json"`
}

type OCMShareProviderJSONDriver struct {
	File string `yaml:"file" env:"OCM_OCM_SHAREPROVIDER_JSON_FILE" desc:"Path to the JSON file where OCM share data will be stored. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH:/storage." introductionVersion:"5.0"`
}

package init

type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret"`
}

type InsecureService struct {
	Insecure bool
}

type ProxyService struct {
	OIDC             InsecureProxyOIDC `yaml:"oidc"`
	InsecureBackends bool              `yaml:"insecure_backends"`
	ServiceAccount   ServiceAccount    `yaml:"service_account"`
}

type InsecureProxyOIDC struct {
	Insecure bool `yaml:"insecure"`
}

type LdapSettings struct {
	BindPassword string `yaml:"bind_password"`
}
type LdapBasedService struct {
	Ldap LdapSettings
}

type Events struct {
	TLSInsecure bool `yaml:"tls_insecure"`
}
type GraphApplication struct {
	ID string `yaml:"id"`
}

type GraphService struct {
	Application    GraphApplication
	Events         Events
	Spaces         InsecureService
	Identity       LdapBasedService
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type ServiceUserPasswordsSettings struct {
	AdminPassword string `yaml:"admin_password"`
	IdmPassword   string `yaml:"idm_password"`
	RevaPassword  string `yaml:"reva_password"`
	IdpPassword   string `yaml:"idp_password"`
}
type IdmService struct {
	ServiceUserPasswords ServiceUserPasswordsSettings `yaml:"service_user_passwords"`
}

type SettingsService struct {
	ServiceAccountIDs []string `yaml:"service_account_ids"`
}

type FrontendService struct {
	AppHandler     InsecureService `yaml:"app_handler"`
	Archiver       InsecureService
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type OcmService struct {
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type AuthbasicService struct {
	AuthProviders LdapBasedService `yaml:"auth_providers"`
}

type AuthProviderSettings struct {
	Oidc InsecureService
}
type AuthbearerService struct {
	AuthProviders AuthProviderSettings `yaml:"auth_providers"`
}

type UsersAndGroupsService struct {
	Drivers LdapBasedService
}

type ThumbnailSettings struct {
	TransferSecret      string `yaml:"transfer_secret"`
	WebdavAllowInsecure bool   `yaml:"webdav_allow_insecure"`
	Cs3AllowInsecure    bool   `yaml:"cs3_allow_insecure"`
}

type ThumbnailService struct {
	Thumbnail ThumbnailSettings
}

type Search struct {
	Events         Events
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type Audit struct {
	Events Events
}

type Sharing struct {
	Events Events
}

type StorageUsers struct {
	Events         Events
	MountID        string         `yaml:"mount_id"`
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type Gateway struct {
	StorageRegistry StorageRegistry `yaml:"storage_registry"`
}

type StorageRegistry struct {
	StorageUsersMountID string `yaml:"storage_users_mount_id"`
}

type Notifications struct {
	Notifications  struct{ Events Events } // The notifications config has a field called notifications
	ServiceAccount ServiceAccount          `yaml:"service_account"`
}

type Userlog struct {
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type AuthService struct {
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type Clientlog struct {
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

type WopiApp struct {
	Secret string `yaml:"secret"`
}

type App struct {
	Insecure bool `yaml:"insecure"`
}

type Collaboration struct {
	WopiApp WopiApp `yaml:"wopi"`
	App     App     `yaml:"app"`
}

type Nats struct {
	// The nats config has a field called nats
	Nats struct {
		TLSSkipVerifyClientCert bool `yaml:"tls_skip_verify_client_cert"`
	}
}

// Activitylog is the configuration for the activitylog service
type Activitylog struct {
	ServiceAccount ServiceAccount `yaml:"service_account"`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ServiceAccountID     string `yaml:"service_account_id"`
	ServiceAccountSecret string `yaml:"service_account_secret"`
}

// TODO: use the oCIS config struct instead of this custom struct
// We can't use it right now, because it would need  "omitempty" on
// all elements, in order to produce a slim config file with `ocis init`.
// We can't just add these "omitempty" tags, since we want to generate
// full example configuration files with that struct, too.
// Proposed solution to  get rid of this temporary solution:
// - use the oCIS config struct
// - set the needed values like below
// - marshal it to yaml
// - unmarshal it into yaml.Node
// - recurse through the nodes and delete empty / default ones
// - marshal it to yaml

// OcisConfig is the configuration for the oCIS services
type OcisConfig struct {
	TokenManager      TokenManager          `yaml:"token_manager"`
	MachineAuthAPIKey string                `yaml:"machine_auth_api_key"`
	SystemUserAPIKey  string                `yaml:"system_user_api_key"`
	TransferSecret    string                `yaml:"transfer_secret"`
	SystemUserID      string                `yaml:"system_user_id"`
	AdminUserID       string                `yaml:"admin_user_id"`
	Graph             GraphService          `yaml:"graph"`
	Idp               LdapBasedService      `yaml:"idp"`
	Idm               IdmService            `yaml:"idm"`
	Collaboration     Collaboration         `yaml:"collaboration"`
	Proxy             ProxyService          `yaml:"proxy"`
	Frontend          FrontendService       `yaml:"frontend"`
	AuthBasic         AuthbasicService      `yaml:"auth_basic"`
	AuthBearer        AuthbearerService     `yaml:"auth_bearer"`
	Users             UsersAndGroupsService `yaml:"users"`
	Groups            UsersAndGroupsService `yaml:"groups"`
	Ocdav             InsecureService       `yaml:"ocdav"`
	Ocm               OcmService            `yaml:"ocm"`
	Thumbnails        ThumbnailService      `yaml:"thumbnails"`
	Search            Search                `yaml:"search"`
	Audit             Audit                 `yaml:"audit"`
	Settings          SettingsService       `yaml:"settings"`
	Sharing           Sharing               `yaml:"sharing"`
	StorageUsers      StorageUsers          `yaml:"storage_users"`
	Notifications     Notifications         `yaml:"notifications"`
	Nats              Nats                  `yaml:"nats"`
	Gateway           Gateway               `yaml:"gateway"`
	Userlog           Userlog               `yaml:"userlog"`
	AuthService       AuthService           `yaml:"auth_service"`
	Clientlog         Clientlog             `yaml:"clientlog"`
	Activitylog       Activitylog           `yaml:"activitylog"`
}

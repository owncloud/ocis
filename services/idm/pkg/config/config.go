package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	IDM                Settings `yaml:"idm"`
	CreateDemoUsers    bool     `yaml:"create_demo_users" env:"IDM_CREATE_DEMO_USERS" desc:"Flag to enable or disable the creation of the demo users." introductionVersion:"pre5.0"`
	DemoUsersIssuerUrl string   `yaml:"demo_users_issuer_url" env:"OCIS_URL;OCIS_OIDC_ISSUER" desc:"The OIDC issuer URL to assign to the demo users." introductionVersion:"pre5.0"`

	ServiceUserPasswords ServiceUserPasswords `yaml:"service_user_passwords"`
	AdminUserID          string               `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID;IDM_ADMIN_USER_ID" desc:"ID of the user that should receive admin privileges. Consider that the UUID can be encoded in some LDAP deployment configurations like in .ldif files. These need to be decoded beforehand." introductionVersion:"pre5.0"`

	Context context.Context `yaml:"-"`
}

type Settings struct {
	LDAPSAddr    string `yaml:"ldaps_addr" env:"IDM_LDAPS_ADDR" desc:"Listen address for the LDAPS listener (ip-addr:port)." introductionVersion:"pre5.0"`
	Cert         string `yaml:"cert" env:"IDM_LDAPS_CERT" desc:"File name of the TLS server certificate for the LDAPS listener. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm." introductionVersion:"pre5.0"`
	Key          string `yaml:"key" env:"IDM_LDAPS_KEY" desc:"File name for the TLS certificate key for the server certificate. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm." introductionVersion:"pre5.0"`
	DatabasePath string `yaml:"database" env:"IDM_DATABASE_PATH" desc:"Full path to the IDM backend database. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/idm." introductionVersion:"pre5.0"`
}

type ServiceUserPasswords struct {
	OcisAdmin string `yaml:"admin_password" env:"IDM_ADMIN_PASSWORD" desc:"Password to set for the oCIS 'admin' user. Either cleartext or an argon2id hash." introductionVersion:"pre5.0"`
	Idm       string `yaml:"idm_password" env:"IDM_SVC_PASSWORD" desc:"Password to set for the 'idm' service user. Either cleartext or an argon2id hash." introductionVersion:"pre5.0"`
	Reva      string `yaml:"reva_password" env:"IDM_REVASVC_PASSWORD" desc:"Password to set for the 'reva' service user. Either cleartext or an argon2id hash." introductionVersion:"pre5.0"`
	Idp       string `yaml:"idp_password" env:"IDM_IDPSVC_PASSWORD" desc:"Password to set for the 'idp' service user. Either cleartext or an argon2id hash." introductionVersion:"pre5.0"`
}

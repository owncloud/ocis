package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	IDM             Settings `yaml:"idm"`
	CreateDemoUsers bool     `yaml:"create_demo_users" env:"IDM_CREATE_DEMO_USERS;ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"Flag to enabe/disable the creation of the demo users"`

	ServiceUserPasswords ServiceUserPasswords `yaml:"service_user_passwords"`

	Context context.Context `yaml:"-"`
}

type Settings struct {
	LDAPSAddr    string `yaml:"ldaps_addr" env:"IDM_LDAPS_ADDR" desc:"Listen address for the ldaps listener (ip-addr:port)"`
	Cert         string `yaml:"cert" env:"IDM_LDAPS_CERT" desc:"File name of the TLS server certificate for the ldaps listener"`
	Key          string `yaml:"key" env:"IDM_LDAPS_KEY" desc:"File name for the TLS certificate key for the server certificate"`
	DatabasePath string `yaml:"database" env:"IDM_DATABASE_PATH" desc:"Full path to the idm backend database"`
}

type ServiceUserPasswords struct {
	IdmAdmin string `yaml:"admin_password" env:"IDM_ADMIN_PASSWORD" desc:"Password to set for the \"idm\" service users. Either cleartext or an argon2id hash"`
	Reva     string `yaml:"reva_password" env:"IDM_REVASVC_PASSWORD" desc:"Password to set for the \"reva\" service users. Either cleartext or an argon2id hash"`
	Idp      string `yaml:"idp_password" env:"IDM_IDPSVC_PASSWORD" desc:"Password to set for the \"idp\" service users. Either cleartext or an argon2id hash"`
}

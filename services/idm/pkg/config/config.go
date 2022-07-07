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

	IDM             Settings `yaml:"idm"`
	CreateDemoUsers bool     `yaml:"create_demo_users" env:"IDM_CREATE_DEMO_USERS;ACCOUNTS_DEMO_USERS_AND_GROUPS" desc:"Flag to enable or disable the creation of the demo users."`

	ServiceUserPasswords ServiceUserPasswords `yaml:"service_user_passwords"`
	AdminUserID          string               `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID;IDM_ADMIN_USER_ID" desc:"ID of the user that should receive admin privileges."`

	Context context.Context `yaml:"-"`
}

type Settings struct {
	LDAPSAddr    string `yaml:"ldaps_addr" env:"IDM_LDAPS_ADDR" desc:"Listen address for the LDAPS listener (ip-addr:port)."`
	Cert         string `yaml:"cert" env:"IDM_LDAPS_CERT" desc:"File name of the TLS server certificate for the LDAPS listener."`
	Key          string `yaml:"key" env:"IDM_LDAPS_KEY" desc:"File name for the TLS certificate key for the server certificate."`
	DatabasePath string `yaml:"database" env:"IDM_DATABASE_PATH" desc:"Full path to the IDM backend database."`
}

type ServiceUserPasswords struct {
	OcisAdmin string `yaml:"admin_password" env:"IDM_ADMIN_PASSWORD" desc:"Password to set for the oCIS \"admin\" user. Either cleartext or an argon2id hash."`
	Idm       string `yaml:"idm_password" env:"IDM_SVC_PASSWORD" desc:"Password to set for the \"idm\" service user. Either cleartext or an argon2id hash."`
	Reva      string `yaml:"reva_password" env:"IDM_REVASVC_PASSWORD" desc:"Password to set for the \"reva\" service user. Either cleartext or an argon2id hash."`
	Idp       string `yaml:"idp_password" env:"IDM_IDPSVC_PASSWORD" desc:"Password to set for the \"idp\" service user. Either cleartext or an argon2id hash."`
}

package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Service: Service{
			Name: "idm",
		},
		CreateDemoUsers: true,
		ServiceUserPasswords: ServiceUserPasswords{
			IdmAdmin: "idm",
			Idp:      "idp",
			Reva:     "reva",
		},
		IDM: Settings{
			LDAPSAddr:    "127.0.0.1:9235",
			Cert:         path.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
			Key:          path.Join(defaults.BaseDataPath(), "idm", "ldap.key"),
			DatabasePath: path.Join(defaults.BaseDataPath(), "idm", "ocis.boltdb"),
		},
	}
}

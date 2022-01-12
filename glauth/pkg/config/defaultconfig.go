package config

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr: "127.0.0.1:9129",
		},
		Service: Service{
			Name: "glauth",
		},
		Ldap: Ldap{
			Enabled:   true,
			Addr:      "127.0.0.1:9125",
			Namespace: "com.owncloud.ldap",
		},
		Ldaps: Ldaps{
			Enabled:   true,
			Addr:      "127.0.0.1:9126",
			Namespace: "com.owncloud.ldaps",
			Cert:      path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
			Key:       path.Join(defaults.BaseDataPath(), "ldap", "ldap.key"),
		},
		Backend: Backend{
			Datastore:   "accounts",
			BaseDN:      "dc=ocis,dc=test",
			Insecure:    false,
			NameFormat:  "cn",
			GroupFormat: "ou",
			Servers:     nil,
			SSHKeyAttr:  "sshPublicKey",
			UseGraphAPI: true,
		},
		Fallback: FallbackBackend{
			Datastore:   "",
			BaseDN:      "dc=ocis,dc=test",
			Insecure:    false,
			NameFormat:  "cn",
			GroupFormat: "ou",
			Servers:     nil,
			SSHKeyAttr:  "sshPublicKey",
			UseGraphAPI: true,
		},
		RoleBundleUUID: "71881883-1768-46bd-a24d-a356a2afdf7f", // BundleUUIDRoleAdmin
	}
}

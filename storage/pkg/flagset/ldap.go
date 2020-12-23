package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// LDAPWithConfig applies LDAP cfg to the flagset
func LDAPWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "ldap-hostname",
			Value:       "localhost",
			Usage:       "LDAP hostname",
			EnvVars:     []string{"STORAGE_LDAP_HOSTNAME"},
			Destination: &cfg.Reva.LDAP.Hostname,
		},
		&cli.IntFlag{
			Name:        "ldap-port",
			Value:       9126,
			Usage:       "LDAP port",
			EnvVars:     []string{"STORAGE_LDAP_PORT"},
			Destination: &cfg.Reva.LDAP.Port,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       "dc=example,dc=org",
			Usage:       "LDAP basedn",
			EnvVars:     []string{"STORAGE_LDAP_BASE_DN"},
			Destination: &cfg.Reva.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-loginfilter",
			Value:       "(&(objectclass=posixAccount)(|(cn={{login}})(mail={{login}})))",
			Usage:       "LDAP login filter",
			EnvVars:     []string{"STORAGE_LDAP_LOGINFILTER"},
			Destination: &cfg.Reva.LDAP.LoginFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
			Usage:       "LDAP filter used when getting a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-attributefilter",
			Value:       "(&(objectclass=posixAccount)({{attr}}={{value}}))",
			Usage:       "LDAP filter used when searching for a user by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"STORAGE_LDAP_ATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.AttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-findfilter",
			Value:       "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
			Usage:       "LDAP filter used when searching for recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"STORAGE_LDAP_FINDFILTER"},
			Destination: &cfg.Reva.LDAP.FindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-groupfilter",
			// FIXME the storage implementation needs to use the memberof overlay to get the cn when it only has the uuid,
			// because the ldap schema either uses the dn or the member(of) attributes to establish membership
			Value:       "(&(objectclass=posixGroup)(ownclouduuid={{.OpaqueId}}*))", // This filter will never work
			Usage:       "LDAP filter used when getting the groups of a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       "cn=reva,ou=sysusers,dc=example,dc=org",
			Usage:       "LDAP bind dn",
			EnvVars:     []string{"STORAGE_LDAP_BIND_DN"},
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       "reva",
			Usage:       "LDAP bind password",
			EnvVars:     []string{"STORAGE_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-idp",
			Value:       "https://localhost:9200",
			Usage:       "Identity provider to use for users",
			EnvVars:     []string{"STORAGE_LDAP_IDP", "OCIS_URL"}, // STORAGE_LDAP_IDP takes precedence over OCIS_URL
			Destination: &cfg.Reva.LDAP.IDP,
		},
		// ldap dn is always the dn
		&cli.StringFlag{
			Name:        "ldap-schema-uid",
			Value:       "ownclouduuid",
			Usage:       "LDAP schema uid",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_UID"},
			Destination: &cfg.Reva.LDAP.Schema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-mail",
			Value:       "mail",
			Usage:       "LDAP schema mail",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.Schema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-displayName",
			Value:       "displayname",
			Usage:       "LDAP schema displayName",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.Schema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-cn",
			Value:       "cn",
			Usage:       "LDAP schema cn",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.Schema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-uidnumber",
			Value:       "uidnumber",
			Usage:       "LDAP schema uidnumber",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.Schema.UIDNumber,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-gidnumber",
			Value:       "gidnumber",
			Usage:       "LDAP schema gidnumber",
			EnvVars:     []string{"STORAGE_LDAP_SCHEMA_GIDNUMBER"},
			Destination: &cfg.Reva.LDAP.Schema.GIDNumber,
		},
	}
}

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

		// User specific filters

		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
			Usage:       "LDAP filter used when getting a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-userattributefilter",
			Value:       "(&(objectclass=posixAccount)({{attr}}={{value}}))",
			Usage:       "LDAP filter used when searching for a user by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"STORAGE_LDAP_USERATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.UserAttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-userfindfilter",
			Value:       "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
			Usage:       "LDAP filter used when searching for user recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"STORAGE_LDAP_USERFINDFILTER"},
			Destination: &cfg.Reva.LDAP.UserFindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-usergroupfilter",
			// FIXME the storage implementation needs to use the memberof overlay to get the cn when it only has the uuid,
			// because the ldap schema either uses the dn or the member(of) attributes to establish membership
			Value:       "(&(objectclass=posixGroup)(ownclouduuid={{.OpaqueId}}*))", // This filter will never work
			Usage:       "LDAP filter used when getting the groups of a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_USERGROUPFILTER"},
			Destination: &cfg.Reva.LDAP.UserGroupFilter,
		},

		// Group specific filters
		// These might not work at the moment. Need to be fixed

		&cli.StringFlag{
			Name:        "ldap-groupfilter",
			Value:       "(&(objectclass=posixGroup)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
			Usage:       "LDAP filter used when getting a group. The CS3 groupid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupattributefilter",
			Value:       "(&(objectclass=posixGroup)({{attr}}={{value}}))",
			Usage:       "LDAP filter used when searching for a group by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.GroupAttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupfindfilter",
			Value:       "(&(objectclass=posixGroup)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
			Usage:       "LDAP filter used when searching for group recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"STORAGE_LDAP_GROUPFINDFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-groupmemberfilter",
			// FIXME the storage implementation needs to use the members overlay to get the cn when it only has the uuid
			Value:       "(&(objectclass=posixAccount)(ownclouduuid={{.OpaqueId}}*))", // This filter will never work
			Usage:       "LDAP filter used when getting the members of a group. The CS3 groupid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPMEMBERFILTER"},
			Destination: &cfg.Reva.LDAP.GroupMemberFilter,
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

		// user schema

		&cli.StringFlag{
			Name:        "ldap-user-schema-uid",
			Value:       "ownclouduuid",
			Usage:       "LDAP user schema uid",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID"},
			Destination: &cfg.Reva.LDAP.UserSchema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-mail",
			Value:       "mail",
			Usage:       "LDAP user schema mail",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.UserSchema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-displayName",
			Value:       "displayname",
			Usage:       "LDAP user schema displayName",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.UserSchema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-cn",
			Value:       "cn",
			Usage:       "LDAP user schema cn",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.UserSchema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-uidnumber",
			Value:       "uidnumber",
			Usage:       "LDAP user schema uidnumber",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.UIDNumber,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-gidnumber",
			Value:       "gidnumber",
			Usage:       "LDAP user schema gidnumber",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.GIDNumber,
		},

		// group schema

		&cli.StringFlag{
			Name:        "ldap-group-schema-gid",
			Value:       "cn",
			Usage:       "LDAP group schema gid",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GID,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-mail",
			Value:       "mail",
			Usage:       "LDAP group schema mail",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-displayName",
			Value:       "cn",
			Usage:       "LDAP group schema displayName",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.GroupSchema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-cn",
			Value:       "cn",
			Usage:       "LDAP group schema cn",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.GroupSchema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-gidnumber",
			Value:       "gidnumber",
			Usage:       "LDAP group schema gidnumber",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GIDNumber,
		},
	}
}

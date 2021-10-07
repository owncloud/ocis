package flagset

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// LDAPWithConfig applies LDAP cfg to the flagset
func LDAPWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "ldap-hostname",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.Hostname, "localhost"),
			Usage:       "LDAP hostname",
			EnvVars:     []string{"STORAGE_LDAP_HOSTNAME"},
			Destination: &cfg.Reva.LDAP.Hostname,
		},
		&cli.IntFlag{
			Name:        "ldap-port",
			Value:       flags.OverrideDefaultInt(cfg.Reva.LDAP.Port, 9126),
			Usage:       "LDAP port",
			EnvVars:     []string{"STORAGE_LDAP_PORT"},
			Destination: &cfg.Reva.LDAP.Port,
		},
		&cli.StringFlag{
			Name:        "ldap-cacert",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.CACert, path.Join(defaults.BaseDataPath(), "ldap", "ldap.crt")),
			Usage:       "Path to a trusted Certificate file (in PEM format) for the LDAP Connection",
			EnvVars:     []string{"STORAGE_LDAP_CACERT"},
			Destination: &cfg.Reva.LDAP.CACert,
		},
		&cli.BoolFlag{
			Name:        "ldap-insecure",
			Value:       flags.OverrideDefaultBool(cfg.Reva.LDAP.Insecure, false),
			Usage:       "Disable TLS certificate and hostname validation",
			EnvVars:     []string{"STORAGE_LDAP_INSECURE"},
			Destination: &cfg.Reva.LDAP.Insecure,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.BaseDN, "dc=ocis,dc=test"),
			Usage:       "LDAP basedn",
			EnvVars:     []string{"STORAGE_LDAP_BASE_DN"},
			Destination: &cfg.Reva.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-loginfilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.LoginFilter, "(&(objectclass=posixAccount)(|(cn={{login}})(mail={{login}})))"),
			Usage:       "LDAP login filter",
			EnvVars:     []string{"STORAGE_LDAP_LOGINFILTER"},
			Destination: &cfg.Reva.LDAP.LoginFilter,
		},

		// User specific filters

		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserFilter, "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))"),
			Usage:       "LDAP filter used when getting a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-userattributefilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserAttributeFilter, "(&(objectclass=posixAccount)({{attr}}={{value}}))"),
			Usage:       "LDAP filter used when searching for a user by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"STORAGE_LDAP_USERATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.UserAttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-userfindfilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserFindFilter, "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))"),
			Usage:       "LDAP filter used when searching for user recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"STORAGE_LDAP_USERFINDFILTER"},
			Destination: &cfg.Reva.LDAP.UserFindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-usergroupfilter",
			// FIXME the storage implementation needs to use the memberof overlay to get the cn when it only has the uuid,
			// because the ldap schema either uses the dn or the member(of) attributes to establish membership
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserGroupFilter, "(&(objectclass=posixGroup)(ownclouduuid={{.OpaqueId}}*))"), // This filter will never work
			Usage:       "LDAP filter used when getting the groups of a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_USERGROUPFILTER"},
			Destination: &cfg.Reva.LDAP.UserGroupFilter,
		},

		// Group specific filters
		// These might not work at the moment. Need to be fixed

		&cli.StringFlag{
			Name:        "ldap-groupfilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupFilter, "(&(objectclass=posixGroup)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))"),
			Usage:       "LDAP filter used when getting a group. The CS3 groupid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupattributefilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupAttributeFilter, "(&(objectclass=posixGroup)({{attr}}={{value}}))"),
			Usage:       "LDAP filter used when searching for a group by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.GroupAttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupfindfilter",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupFindFilter, "(&(objectclass=posixGroup)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))"),
			Usage:       "LDAP filter used when searching for group recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"STORAGE_LDAP_GROUPFINDFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-groupmemberfilter",
			// FIXME the storage implementation needs to use the members overlay to get the cn when it only has the uuid
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupMemberFilter, "(&(objectclass=posixAccount)(ownclouduuid={{.OpaqueId}}*))"), // This filter will never work
			Usage:       "LDAP filter used when getting the members of a group. The CS3 groupid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"STORAGE_LDAP_GROUPMEMBERFILTER"},
			Destination: &cfg.Reva.LDAP.GroupMemberFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.BindDN, "cn=reva,ou=sysusers,dc=ocis,dc=test"),
			Usage:       "LDAP bind dn",
			EnvVars:     []string{"STORAGE_LDAP_BIND_DN"},
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.BindPassword, "reva"),
			Usage:       "LDAP bind password",
			EnvVars:     []string{"STORAGE_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-idp",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.IDP, "https://localhost:9200"),
			Usage:       "Identity provider to use for users",
			EnvVars:     []string{"STORAGE_LDAP_IDP", "OCIS_URL"}, // STORAGE_LDAP_IDP takes precedence over OCIS_URL
			Destination: &cfg.Reva.LDAP.IDP,
		},
		// ldap dn is always the dn

		// user schema

		&cli.StringFlag{
			Name:        "ldap-user-schema-uid",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.UID, "ownclouduuid"),
			Usage:       "LDAP user schema uid",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID"},
			Destination: &cfg.Reva.LDAP.UserSchema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-mail",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.Mail, "mail"),
			Usage:       "LDAP user schema mail",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.UserSchema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-displayName",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.DisplayName, "displayname"),
			Usage:       "LDAP user schema displayName",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.UserSchema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-cn",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.CN, "cn"),
			Usage:       "LDAP user schema cn",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.UserSchema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-uidnumber",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.UIDNumber, "uidnumber"),
			Usage:       "LDAP user schema uidnumber",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.UIDNumber,
		},
		&cli.StringFlag{
			Name:        "ldap-user-schema-gidnumber",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.UserSchema.GIDNumber, "gidnumber"),
			Usage:       "LDAP user schema gidnumber",
			EnvVars:     []string{"STORAGE_LDAP_USER_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.UserSchema.GIDNumber,
		},

		// group schema

		&cli.StringFlag{
			Name:        "ldap-group-schema-gid",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupSchema.GID, "cn"),
			Usage:       "LDAP group schema gid",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GID,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-mail",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupSchema.Mail, "mail"),
			Usage:       "LDAP group schema mail",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.GroupSchema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-displayName",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupSchema.DisplayName, "cn"),
			Usage:       "LDAP group schema displayName",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.GroupSchema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-cn",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupSchema.CN, "cn"),
			Usage:       "LDAP group schema cn",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.GroupSchema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-group-schema-gidnumber",
			Value:       flags.OverrideDefaultString(cfg.Reva.LDAP.GroupSchema.GIDNumber, "gidnumber"),
			Usage:       "LDAP group schema gidnumber",
			EnvVars:     []string{"STORAGE_LDAP_GROUP_SCHEMA_GID_NUMBER"},
			Destination: &cfg.Reva.LDAP.GroupSchema.GIDNumber,
		},
	}
}

package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// UsersWithConfig applies cfg to the root flagset
func UsersWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"REVA_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"REVA_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"REVA_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"REVA_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"REVA_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9145",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_SHARING_DEBUG_ADDR"},
			Destination: &cfg.Reva.Users.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"REVA_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"REVA_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"REVA_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},

		// REVA

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVars:     []string{"REVA_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},

		// LDAP

		&cli.StringFlag{
			Name:        "ldap-hostname",
			Value:       "localhost",
			Usage:       "LDAP hostname",
			EnvVars:     []string{"REVA_LDAP_HOSTNAME"},
			Destination: &cfg.Reva.LDAP.Hostname,
		},
		&cli.IntFlag{
			Name:        "ldap-port",
			Value:       9126,
			Usage:       "LDAP port",
			EnvVars:     []string{"REVA_LDAP_PORT"},
			Destination: &cfg.Reva.LDAP.Port,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       "dc=example,dc=org",
			Usage:       "LDAP basedn",
			EnvVars:     []string{"REVA_LDAP_BASE_DN"},
			Destination: &cfg.Reva.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       "(&(objectclass=posixAccount)(|(ownclouduuid={{.OpaqueId}})(cn={{.OpaqueId}})))",
			Usage:       "LDAP filter used when getting a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"REVA_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-attributefilter",
			Value:       "(&(objectclass=posixAccount)({{attr}}={{value}}))",
			Usage:       "LDAP filter used when searching for a user by claim/attribute. {{attr}} will be replaced with the attribute, {{value}} with the value.",
			EnvVars:     []string{"REVA_LDAP_ATTRIBUTEFILTER"},
			Destination: &cfg.Reva.LDAP.AttributeFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-findfilter",
			Value:       "(&(objectclass=posixAccount)(|(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)))",
			Usage:       "LDAP filter used when searching for recipients. {{query}} will be replaced with the search query",
			EnvVars:     []string{"REVA_LDAP_FINDFILTER"},
			Destination: &cfg.Reva.LDAP.FindFilter,
		},
		&cli.StringFlag{
			Name: "ldap-groupfilter",
			// FIXME the reva implementation needs to use the memberof overlay to get the cn when it only has the uuid,
			// because the ldap schema either uses the dn or the member(of) attributes to establish membership
			Value:       "(&(objectclass=posixGroup)(ownclouduuid={{.OpaqueId}}*))", // This filter will never work
			Usage:       "LDAP filter used when getting the groups of a user. The CS3 userid properties {{.OpaqueId}} and {{.Idp}} are available.",
			EnvVars:     []string{"REVA_LDAP_GROUPFILTER"},
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       "cn=reva,ou=sysusers,dc=example,dc=org",
			Usage:       "LDAP bind dn",
			EnvVars:     []string{"REVA_LDAP_BIND_DN"},
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       "reva",
			Usage:       "LDAP bind password",
			EnvVars:     []string{"REVA_LDAP_BIND_PASSWORD"},
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		&cli.StringFlag{
			Name:        "ldap-idp",
			Value:       "https://localhost:9200",
			Usage:       "Identity provider to use for users",
			EnvVars:     []string{"REVA_LDAP_IDP"},
			Destination: &cfg.Reva.LDAP.IDP,
		},
		// ldap dn is always the dn
		&cli.StringFlag{
			Name:        "ldap-schema-uid",
			Value:       "ownclouduuid",
			Usage:       "LDAP schema uid",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_UID"},
			Destination: &cfg.Reva.LDAP.Schema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-mail",
			Value:       "mail",
			Usage:       "LDAP schema mail",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_MAIL"},
			Destination: &cfg.Reva.LDAP.Schema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-displayName",
			Value:       "displayname",
			Usage:       "LDAP schema displayName",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_DISPLAYNAME"},
			Destination: &cfg.Reva.LDAP.Schema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-cn",
			Value:       "cn",
			Usage:       "LDAP schema cn",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_CN"},
			Destination: &cfg.Reva.LDAP.Schema.CN,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-uidnumber",
			Value:       "uidnumber",
			Usage:       "LDAP schema uidnumber",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_UID_NUMBER"},
			Destination: &cfg.Reva.LDAP.Schema.UIDNumber,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-gidnumber",
			Value:       "gidnumber",
			Usage:       "LDAP schema gidnumber",
			EnvVars:     []string{"REVA_LDAP_SCHEMA_GIDNUMBER"},
			Destination: &cfg.Reva.LDAP.Schema.GIDNumber,
		},
		&cli.StringFlag{
			Name:        "rest-client-id",
			Value:       "",
			Usage:       "User rest driver Client ID",
			EnvVars:     []string{"REVA_REST_CLIENT_ID"},
			Destination: &cfg.Reva.UserRest.ClientID,
		},
		&cli.StringFlag{
			Name:        "rest-client-secret",
			Value:       "",
			Usage:       "User rest driver Client Secret",
			EnvVars:     []string{"REVA_REST_CLIENT_SECRET"},
			Destination: &cfg.Reva.UserRest.ClientSecret,
		},
		&cli.StringFlag{
			Name:        "rest-redis-address",
			Value:       "localhost:6379",
			Usage:       "Address for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_ADDRESS"},
			Destination: &cfg.Reva.UserRest.RedisAddress,
		},
		&cli.StringFlag{
			Name:        "rest-redis-username",
			Value:       "",
			Usage:       "Username for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_USERNAME"},
			Destination: &cfg.Reva.UserRest.RedisUsername,
		},
		&cli.StringFlag{
			Name:        "rest-redis-password",
			Value:       "",
			Usage:       "Password for redis server",
			EnvVars:     []string{"REVA_REST_REDIS_PASSWORD"},
			Destination: &cfg.Reva.UserRest.RedisPassword,
		},
		&cli.IntFlag{
			Name:        "rest-user-groups-cache-expiration",
			Value:       5,
			Usage:       "Time in minutes for redis cache expiration.",
			EnvVars:     []string{"REVA_REST_CACHE_EXPIRATION"},
			Destination: &cfg.Reva.UserRest.UserGroupsCacheExpiration,
		},
		&cli.StringFlag{
			Name:        "rest-id-provider",
			Value:       "",
			Usage:       "The OIDC Provider",
			EnvVars:     []string{"REVA_REST_ID_PROVIDER"},
			Destination: &cfg.Reva.UserRest.IDProvider,
		},
		&cli.StringFlag{
			Name:        "rest-api-base-url",
			Value:       "",
			Usage:       "Base API Endpoint",
			EnvVars:     []string{"REVA_REST_API_BASE_URL"},
			Destination: &cfg.Reva.UserRest.APIBaseURL,
		},
		&cli.StringFlag{
			Name:        "rest-oidc-token-endpoint",
			Value:       "",
			Usage:       "Endpoint to generate token to access the API",
			EnvVars:     []string{"REVA_REST_OIDC_TOKEN_ENDPOINT"},
			Destination: &cfg.Reva.UserRest.OIDCTokenEndpoint,
		},
		&cli.StringFlag{
			Name:        "rest-target-api",
			Value:       "",
			Usage:       "The target application",
			EnvVars:     []string{"REVA_REST_TARGET_API"},
			Destination: &cfg.Reva.UserRest.TargetAPI,
		},

		// Services

		// Users

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_USERS_NETWORK"},
			Destination: &cfg.Reva.Users.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_USERS_PROTOCOL"},
			Destination: &cfg.Reva.Users.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9144",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_USERS_ADDR"},
			Destination: &cfg.Reva.Users.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("userprovider"), // TODO preferences
			Usage:   "--service userprovider [--service otherservice]",
			EnvVars: []string{"REVA_USERS_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "ldap",
			Usage:       "user driver: 'demo', 'json', 'ldap', or 'rest'",
			EnvVars:     []string{"REVA_USERS_DRIVER"},
			Destination: &cfg.Reva.Users.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVars:     []string{"REVA_USERS_JSON"},
			Destination: &cfg.Reva.Users.JSON,
		},
	}
}

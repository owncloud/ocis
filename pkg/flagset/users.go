package flagset

import (
	"github.com/micro/cli"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// UsersWithConfig applies cfg to the root flagset
func UsersWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{

		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVar:      "REVA_TRACING_ENABLED",
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVar:      "REVA_TRACING_TYPE",
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVar:      "REVA_TRACING_ENDPOINT",
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVar:      "REVA_TRACING_COLLECTOR",
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVar:      "REVA_TRACING_SERVICE",
			Destination: &cfg.Tracing.Service,
		},

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9145",
			Usage:       "Address to bind debug server",
			EnvVar:      "REVA_SHARING_DEBUG_ADDR",
			Destination: &cfg.Reva.Users.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVar:      "REVA_DEBUG_TOKEN",
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVar:      "REVA_DEBUG_PPROF",
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVar:      "REVA_DEBUG_ZPAGES",
			Destination: &cfg.Debug.Zpages,
		},

		// REVA

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVar:      "REVA_JWT_SECRET",
			Destination: &cfg.Reva.JWTSecret,
		},

		// LDAP

		&cli.StringFlag{
			Name:        "ldap-hostname",
			Value:       "localhost",
			Usage:       "LDAP hostname",
			EnvVar:      "REVA_LDAP_HOSTNAME",
			Destination: &cfg.Reva.LDAP.Hostname,
		},
		&cli.IntFlag{
			Name:        "ldap-port",
			Value:       389,
			Usage:       "LDAP port",
			EnvVar:      "REVA_LDAP_PORT",
			Destination: &cfg.Reva.LDAP.Port,
		},
		&cli.StringFlag{
			Name:        "ldap-base-dn",
			Value:       "dc=owncloud,dc=com",
			Usage:       "LDAP basedn",
			EnvVar:      "REVA_LDAP_BASE_DN",
			Destination: &cfg.Reva.LDAP.BaseDN,
		},
		&cli.StringFlag{
			Name:        "ldap-userfilter",
			Value:       "(objectclass=posixAccount)",
			Usage:       "LDAP userfilter",
			EnvVar:      "REVA_LDAP_USERFILTER",
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupfilter",
			Value:       "(objectclass=posixGroup)",
			Usage:       "LDAP groupfilter",
			EnvVar:      "REVA_LDAP_GROUPFILTER",
			Destination: &cfg.Reva.LDAP.GroupFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-dn",
			Value:       "cn=admin,dc=owncloud,dc=com",
			Usage:       "LDAP bind dn",
			EnvVar:      "REVA_LDAP_BIND_DN",
			Destination: &cfg.Reva.LDAP.BindDN,
		},
		&cli.StringFlag{
			Name:        "ldap-bind-password",
			Value:       "admin",
			Usage:       "LDAP bind password",
			EnvVar:      "REVA_LDAP_BIND_PASSWORD",
			Destination: &cfg.Reva.LDAP.BindPassword,
		},
		// ldap dn is always the dn
		&cli.StringFlag{
			Name:        "ldap-schema-uid",
			Value:       "uid",
			Usage:       "LDAP schema uid",
			EnvVar:      "REVA_LDAP_SCHEMA_UID",
			Destination: &cfg.Reva.LDAP.Schema.UID,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-mail",
			Value:       "mail",
			Usage:       "LDAP schema mail",
			EnvVar:      "REVA_LDAP_SCHEMA_MAIL",
			Destination: &cfg.Reva.LDAP.Schema.Mail,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-displayName",
			Value:       "displayName",
			Usage:       "LDAP schema displayName",
			EnvVar:      "REVA_LDAP_SCHEMA_DISPLAYNAME",
			Destination: &cfg.Reva.LDAP.Schema.DisplayName,
		},
		&cli.StringFlag{
			Name:        "ldap-schema-cn",
			Value:       "cn",
			Usage:       "LDAP schema cn",
			EnvVar:      "REVA_LDAP_SCHEMA_CN",
			Destination: &cfg.Reva.LDAP.Schema.CN,
		},

		// Services

		// Users

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVar:      "REVA_USERS_NETWORK",
			Destination: &cfg.Reva.Users.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVar:      "REVA_USERS_PROTOCOL",
			Destination: &cfg.Reva.Users.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9144",
			Usage:       "Address to bind reva service",
			EnvVar:      "REVA_USERS_ADDR",
			Destination: &cfg.Reva.Users.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVar:      "REVA_USERS_URL",
			Destination: &cfg.Reva.Users.URL,
		},
		&cli.StringFlag{
			Name:        "services",
			Value:       "userprovider", // TODO preferences
			Usage:       "comma separated list of services to include",
			EnvVar:      "REVA_USERS_SERVICES",
			Destination: &cfg.Reva.Users.Services,
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "demo",
			Usage:       "user driver: 'demo', 'json' or 'ldap'",
			EnvVar:      "REVA_USERS_DRIVER",
			Destination: &cfg.Reva.Users.Driver,
		},
		&cli.StringFlag{
			Name:        "json-config",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVar:      "REVA_USERS_JSON",
			Destination: &cfg.Reva.Users.JSON,
		},
	}
}

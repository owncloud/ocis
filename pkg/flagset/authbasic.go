package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// AuthBasicWithConfig applies cfg to the root flagset
func AuthBasicWithConfig(cfg *config.Config) []cli.Flag {
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
			Value:       "0.0.0.0:9147",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_AUTH_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBasic.DebugAddr,
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

		// Users

		&cli.StringFlag{
			Name:        "users-driver",
			Value:       "ldap",
			Usage:       "user driver: 'demo', 'json' or 'ldap'",
			EnvVars:     []string{"REVA_USERS_DRIVER"},
			Destination: &cfg.Reva.Users.Driver,
		},
		&cli.StringFlag{
			Name:        "users-json",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVars:     []string{"REVA_USERS_JSON"},
			Destination: &cfg.Reva.Users.JSON,
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
			Value:       "(&(objectclass=posixAccount)(cn=%s))",
			Usage:       "LDAP userfilter",
			EnvVars:     []string{"REVA_LDAP_USERFILTER"},
			Destination: &cfg.Reva.LDAP.UserFilter,
		},
		&cli.StringFlag{
			Name:        "ldap-groupfilter",
			Value:       "(&(objectclass=posixGroup)(cn=%s))",
			Usage:       "LDAP groupfilter",
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
			Value:       "uid",
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
			Value:       "sn",
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

		// Services

		// AuthBasic

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva auth-basic service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_AUTH_BASIC_NETWORK"},
			Destination: &cfg.Reva.AuthBasic.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_AUTH_BASIC_PROTOCOL"},
			Destination: &cfg.Reva.AuthBasic.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9146",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_AUTH_BASIC_ADDR"},
			Destination: &cfg.Reva.AuthBasic.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9146",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_AUTH_BASIC_URL"},
			Destination: &cfg.Reva.AuthBasic.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"),
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"REVA_AUTH_BASIC_SERVICES"},
		},
	}
}

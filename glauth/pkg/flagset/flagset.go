package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/glauth/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       false,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       false,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9129",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"GLAUTH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"GLAUTH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"GLAUTH_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"GLAUTH_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"GLAUTH_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"GLAUTH_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "glauth",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"GLAUTH_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9129",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"GLAUTH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"GLAUTH_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"GLAUTH_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"GLAUTH_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "role-bundle-id",
			Value:       "71881883-1768-46bd-a24d-a356a2afdf7f", // BundleUUIDRoleAdmin
			Usage:       "roleid used to make internal grpc requests",
			EnvVars:     []string{"GLAUTH_ROLE_BUNDLE_ID"},
			Destination: &cfg.RoleBundleUUID,
		},

		&cli.StringFlag{
			Name:        "ldap-addr",
			Value:       "0.0.0.0:9125",
			Usage:       "Address to bind ldap server",
			EnvVars:     []string{"GLAUTH_LDAP_ADDR"},
			Destination: &cfg.Ldap.Address,
		},
		&cli.BoolFlag{
			Name:        "ldap-enabled",
			Value:       true,
			Usage:       "Enable ldap server",
			EnvVars:     []string{"GLAUTH_LDAP_ENABLED"},
			Destination: &cfg.Ldap.Enabled,
		},

		&cli.StringFlag{
			Name:        "ldaps-addr",
			Value:       "0.0.0.0:9126",
			Usage:       "Address to bind ldap server",
			EnvVars:     []string{"GLAUTH_LDAPS_ADDR"},
			Destination: &cfg.Ldaps.Address,
		},
		&cli.BoolFlag{
			Name:        "ldaps-enabled",
			Value:       true,
			Usage:       "Enable ldap server",
			EnvVars:     []string{"GLAUTH_LDAPS_ENABLED"},
			Destination: &cfg.Ldaps.Enabled,
		},
		&cli.StringFlag{
			Name:        "ldaps-cert",
			Value:       "./ldap.crt",
			Usage:       "path to ldaps certificate in PEM format",
			EnvVars:     []string{"GLAUTH_LDAPS_CERT"},
			Destination: &cfg.Ldaps.Cert,
		},
		&cli.StringFlag{
			Name:        "ldaps-key",
			Value:       "./ldap.key",
			Usage:       "path to ldaps key in PEM format",
			EnvVars:     []string{"GLAUTH_LDAPS_KEY"},
			Destination: &cfg.Ldaps.Key,
		},

		// backend config

		&cli.StringFlag{
			Name:        "backend-basedn",
			Value:       "dc=example,dc=org",
			Usage:       "base distinguished name to expose",
			EnvVars:     []string{"GLAUTH_BACKEND_BASEDN"},
			Destination: &cfg.Backend.BaseDN,
		},
		&cli.StringFlag{
			Name:        "backend-name-format",
			Value:       "cn",
			Usage:       "name attribute for entries to expose. typically cn or uid",
			EnvVars:     []string{"GLAUTH_BACKEND_NAME_FORMAT"},
			Destination: &cfg.Backend.NameFormat,
		},
		&cli.StringFlag{
			Name:        "backend-group-format",
			Value:       "ou",
			Usage:       "name attribute for entries to expose. typically ou, cn or dc",
			EnvVars:     []string{"GLAUTH_BACKEND_GROUP_FORMAT"},
			Destination: &cfg.Backend.GroupFormat,
		},
		&cli.StringFlag{
			Name:        "backend-ssh-key-attr",
			Value:       "sshPublicKey",
			Usage:       "ssh key attribute for entries to expose",
			EnvVars:     []string{"GLAUTH_BACKEND_SSH_KEY_ATTR"},
			Destination: &cfg.Backend.SSHKeyAttr,
		},
		&cli.StringFlag{
			Name:  "backend-datastore",
			Value: "accounts",
			// TODO bring back config / flat file support
			Usage:       "datastore to use as the backend. one of accounts, ldap or owncloud",
			EnvVars:     []string{"GLAUTH_BACKEND_DATASTORE"},
			Destination: &cfg.Backend.Datastore,
		},
		&cli.BoolFlag{
			Name:        "backend-insecure",
			Value:       false,
			Usage:       "Allow insecure requests to the datastore",
			EnvVars:     []string{"GLAUTH_BACKEND_INSECURE"},
			Destination: &cfg.Backend.Insecure,
		},
		&cli.StringSliceFlag{
			Name:    "backend-server",
			Value:   cli.NewStringSlice("https://demo.owncloud.com/apps/graphapi/v1.0"),
			Usage:   `--backend-server http://internal1.example.com [--backend-server http://internal2.example.com]`,
			EnvVars: []string{"GLAUTH_BACKEND_SERVERS"},
		},
		&cli.BoolFlag{
			Name:        "backend-use-graphapi",
			Value:       true,
			Usage:       "use Graph API, only for owncloud datastore",
			EnvVars:     []string{"GLAUTH_BACKEND_USE_GRAPHAPI"},
			Destination: &cfg.Backend.UseGraphAPI,
		},

		// fallback config

		&cli.StringFlag{
			Name:        "fallback-basedn",
			Value:       "dc=example,dc=org",
			Usage:       "base distinguished name to expose",
			EnvVars:     []string{"GLAUTH_FALLBACK_BASEDN"},
			Destination: &cfg.Fallback.BaseDN,
		},
		&cli.StringFlag{
			Name:        "fallback-name-format",
			Value:       "cn",
			Usage:       "name attribute for entries to expose. typically cn or uid",
			EnvVars:     []string{"GLAUTH_FALLBACK_NAME_FORMAT"},
			Destination: &cfg.Fallback.NameFormat,
		},
		&cli.StringFlag{
			Name:        "fallback-group-format",
			Value:       "ou",
			Usage:       "name attribute for entries to expose. typically ou, cn or dc",
			EnvVars:     []string{"GLAUTH_FALLBACK_GROUP_FORMAT"},
			Destination: &cfg.Fallback.GroupFormat,
		},
		&cli.StringFlag{
			Name:        "fallback-ssh-key-attr",
			Value:       "sshPublicKey",
			Usage:       "ssh key attribute for entries to expose",
			EnvVars:     []string{"GLAUTH_FALLBACK_SSH_KEY_ATTR"},
			Destination: &cfg.Fallback.SSHKeyAttr,
		},
		&cli.StringFlag{
			Name:  "fallback-datastore",
			Value: "",
			// TODO bring back config / flat file support
			Usage:       "datastore to use as the fallback. one of accounts, ldap or owncloud",
			EnvVars:     []string{"GLAUTH_FALLBACK_DATASTORE"},
			Destination: &cfg.Fallback.Datastore,
		},
		&cli.BoolFlag{
			Name:        "fallback-insecure",
			Value:       false,
			Usage:       "Allow insecure requests to the datastore",
			EnvVars:     []string{"GLAUTH_FALLBACK_INSECURE"},
			Destination: &cfg.Fallback.Insecure,
		},
		&cli.StringSliceFlag{
			Name:    "fallback-server",
			Value:   cli.NewStringSlice("https://demo.owncloud.com/apps/graphapi/v1.0"),
			Usage:   `--fallback-server http://internal1.example.com [--fallback-server http://internal2.example.com]`,
			EnvVars: []string{"GLAUTH_FALLBACK_SERVERS"},
		},
		&cli.BoolFlag{
			Name:        "fallback-use-graphapi",
			Value:       true,
			Usage:       "use Graph API, only for owncloud datastore",
			EnvVars:     []string{"GLAUTH_FALLBACK_USE_GRAPHAPI"},
			Destination: &cfg.Fallback.UseGraphAPI,
		},
	}
}

package flagset

import (
	"path"

	"github.com/owncloud/ocis/glauth/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	pkgos "github.com/owncloud/ocis/ocis-pkg/os"
	"github.com/urfave/cli/v2"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"GLAUTH_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"GLAUTH_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"GLAUTH_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9129"),
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
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"GLAUTH_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       flags.OverrideDefaultString(cfg.File, ""),
			Usage:       "Path to config file",
			EnvVars:     []string{"GLAUTH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"GLAUTH_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"GLAUTH_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"GLAUTH_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"GLAUTH_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "glauth"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"GLAUTH_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9129"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"GLAUTH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
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
			Value:       flags.OverrideDefaultString(cfg.RoleBundleUUID, "71881883-1768-46bd-a24d-a356a2afdf7f"), // BundleUUIDRoleAdmin
			Usage:       "roleid used to make internal grpc requests",
			EnvVars:     []string{"GLAUTH_ROLE_BUNDLE_ID"},
			Destination: &cfg.RoleBundleUUID,
		},

		&cli.StringFlag{
			Name:        "ldap-addr",
			Value:       flags.OverrideDefaultString(cfg.Ldap.Addr, "0.0.0.0:9125"),
			Usage:       "Address to bind ldap server",
			EnvVars:     []string{"GLAUTH_LDAP_ADDR"},
			Destination: &cfg.Ldap.Addr,
		},
		&cli.BoolFlag{
			Name:        "ldap-enabled",
			Value:       flags.OverrideDefaultBool(cfg.Ldap.Enabled, true),
			Usage:       "Enable ldap server",
			EnvVars:     []string{"GLAUTH_LDAP_ENABLED"},
			Destination: &cfg.Ldap.Enabled,
		},

		&cli.StringFlag{
			Name:        "ldaps-addr",
			Value:       flags.OverrideDefaultString(cfg.Ldaps.Addr, "0.0.0.0:9126"),
			Usage:       "Address to bind ldap server",
			EnvVars:     []string{"GLAUTH_LDAPS_ADDR"},
			Destination: &cfg.Ldaps.Addr,
		},
		&cli.BoolFlag{
			Name:        "ldaps-enabled",
			Value:       flags.OverrideDefaultBool(cfg.Ldaps.Enabled, true),
			Usage:       "Enable ldap server",
			EnvVars:     []string{"GLAUTH_LDAPS_ENABLED"},
			Destination: &cfg.Ldaps.Enabled,
		},
		&cli.StringFlag{
			Name:        "ldaps-cert",
			Value:       flags.OverrideDefaultString(cfg.Ldaps.Cert, path.Join(pkgos.MustUserConfigDir("ocis", "ldap"), "ldap.crt")),
			Usage:       "path to ldaps certificate in PEM format",
			EnvVars:     []string{"GLAUTH_LDAPS_CERT"},
			Destination: &cfg.Ldaps.Cert,
		},
		&cli.StringFlag{
			Name:        "ldaps-key",
			Value:       flags.OverrideDefaultString(cfg.Ldaps.Key, path.Join(pkgos.MustUserConfigDir("ocis", "ldap"), "ldap.key")),
			Usage:       "path to ldaps key in PEM format",
			EnvVars:     []string{"GLAUTH_LDAPS_KEY"},
			Destination: &cfg.Ldaps.Key,
		},

		// backend config

		&cli.StringFlag{
			Name:        "backend-basedn",
			Value:       flags.OverrideDefaultString(cfg.Backend.BaseDN, "dc=ocis,dc=test"),
			Usage:       "base distinguished name to expose",
			EnvVars:     []string{"GLAUTH_BACKEND_BASEDN"},
			Destination: &cfg.Backend.BaseDN,
		},
		&cli.StringFlag{
			Name:        "backend-name-format",
			Value:       flags.OverrideDefaultString(cfg.Backend.NameFormat, "cn"),
			Usage:       "name attribute for entries to expose. typically cn or uid",
			EnvVars:     []string{"GLAUTH_BACKEND_NAME_FORMAT"},
			Destination: &cfg.Backend.NameFormat,
		},
		&cli.StringFlag{
			Name:        "backend-group-format",
			Value:       flags.OverrideDefaultString(cfg.Backend.GroupFormat, "ou"),
			Usage:       "name attribute for entries to expose. typically ou, cn or dc",
			EnvVars:     []string{"GLAUTH_BACKEND_GROUP_FORMAT"},
			Destination: &cfg.Backend.GroupFormat,
		},
		&cli.StringFlag{
			Name:        "backend-ssh-key-attr",
			Value:       flags.OverrideDefaultString(cfg.Backend.SSHKeyAttr, "sshPublicKey"),
			Usage:       "ssh key attribute for entries to expose",
			EnvVars:     []string{"GLAUTH_BACKEND_SSH_KEY_ATTR"},
			Destination: &cfg.Backend.SSHKeyAttr,
		},
		&cli.StringFlag{
			Name:  "backend-datastore",
			Value: flags.OverrideDefaultString(cfg.Backend.Datastore, "accounts"),
			// TODO bring back config / flat file support
			Usage:       "datastore to use as the backend. one of accounts, ldap or owncloud",
			EnvVars:     []string{"GLAUTH_BACKEND_DATASTORE"},
			Destination: &cfg.Backend.Datastore,
		},
		&cli.BoolFlag{
			Name:        "backend-insecure",
			Value:       flags.OverrideDefaultBool(cfg.Backend.Insecure, false),
			Usage:       "Allow insecure requests to the datastore",
			EnvVars:     []string{"GLAUTH_BACKEND_INSECURE"},
			Destination: &cfg.Backend.Insecure,
		},
		&cli.StringSliceFlag{
			Name:    "backend-server",
			Value:   cli.NewStringSlice(),
			Usage:   `--backend-server https://demo.owncloud.com/apps/graphapi/v1.0 [--backend-server "https://demo2.owncloud.com/apps/graphapi/v1.0"]`,
			EnvVars: []string{"GLAUTH_BACKEND_SERVERS"},
		},
		&cli.BoolFlag{
			Name:        "backend-use-graphapi",
			Value:       flags.OverrideDefaultBool(cfg.Backend.UseGraphAPI, true),
			Usage:       "use Graph API, only for owncloud datastore",
			EnvVars:     []string{"GLAUTH_BACKEND_USE_GRAPHAPI"},
			Destination: &cfg.Backend.UseGraphAPI,
		},

		// fallback config

		&cli.StringFlag{
			Name:        "fallback-basedn",
			Value:       flags.OverrideDefaultString(cfg.Fallback.BaseDN, "dc=ocis,dc=test"),
			Usage:       "base distinguished name to expose",
			EnvVars:     []string{"GLAUTH_FALLBACK_BASEDN"},
			Destination: &cfg.Fallback.BaseDN,
		},
		&cli.StringFlag{
			Name:        "fallback-name-format",
			Value:       flags.OverrideDefaultString(cfg.Fallback.NameFormat, "cn"),
			Usage:       "name attribute for entries to expose. typically cn or uid",
			EnvVars:     []string{"GLAUTH_FALLBACK_NAME_FORMAT"},
			Destination: &cfg.Fallback.NameFormat,
		},
		&cli.StringFlag{
			Name:        "fallback-group-format",
			Value:       flags.OverrideDefaultString(cfg.Fallback.GroupFormat, "ou"),
			Usage:       "name attribute for entries to expose. typically ou, cn or dc",
			EnvVars:     []string{"GLAUTH_FALLBACK_GROUP_FORMAT"},
			Destination: &cfg.Fallback.GroupFormat,
		},
		&cli.StringFlag{
			Name:        "fallback-ssh-key-attr",
			Value:       flags.OverrideDefaultString(cfg.Fallback.SSHKeyAttr, "sshPublicKey"),
			Usage:       "ssh key attribute for entries to expose",
			EnvVars:     []string{"GLAUTH_FALLBACK_SSH_KEY_ATTR"},
			Destination: &cfg.Fallback.SSHKeyAttr,
		},
		&cli.StringFlag{
			Name:  "fallback-datastore",
			Value: flags.OverrideDefaultString(cfg.Fallback.Datastore, ""),
			// TODO bring back config / flat file support
			Usage:       "datastore to use as the fallback. one of accounts, ldap or owncloud",
			EnvVars:     []string{"GLAUTH_FALLBACK_DATASTORE"},
			Destination: &cfg.Fallback.Datastore,
		},
		&cli.BoolFlag{
			Name:        "fallback-insecure",
			Value:       flags.OverrideDefaultBool(cfg.Fallback.Insecure, false),
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
			Value:       flags.OverrideDefaultBool(cfg.Fallback.UseGraphAPI, true),
			Usage:       "use Graph API, only for owncloud datastore",
			EnvVars:     []string{"GLAUTH_FALLBACK_USE_GRAPHAPI"},
			Destination: &cfg.Fallback.UseGraphAPI,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
		},
	}
}

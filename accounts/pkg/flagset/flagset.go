package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/accounts/pkg/config"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"ACCOUNTS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"ACCOUNTS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"ACCOUNTS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"ACCOUNTS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"ACCOUNTS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       cfg.Tracing.Type,
			Usage:       "Tracing backend type",
			EnvVars:     []string{"ACCOUNTS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       cfg.Tracing.Endpoint,
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"ACCOUNTS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       cfg.Tracing.Collector,
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"ACCOUNTS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       cfg.Tracing.Service,
			Usage:       "Service name for tracing",
			EnvVars:     []string{"ACCOUNTS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       cfg.HTTP.Namespace,
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"ACCOUNTS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       cfg.HTTP.Addr,
			Usage:       "Address to bind http server",
			EnvVars:     []string{"ACCOUNTS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       cfg.HTTP.Root,
			Usage:       "Root path of http server",
			EnvVars:     []string{"ACCOUNTS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.IntFlag{
			Name:        "http-cache-ttl",
			Value:       cfg.HTTP.CacheTTL,
			Usage:       "Set the static assets caching duration in seconds",
			EnvVars:     []string{"ACCOUNTS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       cfg.GRPC.Addr,
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"ACCOUNTS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.IntFlag{
			Name:        "accounts-hash-difficulty",
			Value:       cfg.Server.HashDifficulty,
			Usage:       "accounts password hash difficulty",
			EnvVars:     []string{"ACCOUNTS_HASH_DIFFICULTY"},
			Destination: &cfg.Server.HashDifficulty,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       cfg.Asset.Path,
			Usage:       "Path to custom assets",
			EnvVars:     []string{"ACCOUNTS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       cfg.TokenManager.JWTSecret,
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"ACCOUNTS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "storage-disk-path",
			Value:       cfg.Repo.Disk.Path,
			Usage:       "Path on the local disk, e.g. /var/tmp/ocis/accounts",
			EnvVars:     []string{"ACCOUNTS_STORAGE_DISK_PATH"},
			Destination: &cfg.Repo.Disk.Path,
		},
		&cli.StringFlag{
			Name:        "storage-cs3-provider-addr",
			Value:       cfg.Repo.CS3.ProviderAddr,
			Usage:       "bind address for the metadata storage provider",
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"},
			Destination: &cfg.Repo.CS3.ProviderAddr,
		},
		&cli.StringFlag{
			Name:        "storage-cs3-data-url",
			Value:       cfg.Repo.CS3.DataURL,
			Usage:       "http endpoint of the metadata storage",
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_URL"},
			Destination: &cfg.Repo.CS3.DataURL,
		},
		&cli.StringFlag{
			Name:        "storage-cs3-data-prefix",
			Value:       cfg.Repo.CS3.DataPrefix,
			Usage:       "path prefix for the http endpoint of the metadata storage, without leading slash",
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_PREFIX"},
			Destination: &cfg.Repo.CS3.DataPrefix,
		},
		&cli.StringFlag{
			Name:        "storage-cs3-jwt-secret",
			Value:       cfg.Repo.CS3.JWTSecret,
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.Repo.CS3.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "service-user-uuid",
			Value:       cfg.ServiceUser.UUID,
			Usage:       "uuid of the internal service user (required on EOS)",
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UUID"},
			Destination: &cfg.ServiceUser.UUID,
		},
		&cli.StringFlag{
			Name:        "service-user-username",
			Value:       cfg.ServiceUser.Username,
			Usage:       "username of the internal service user (required on EOS)",
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_USERNAME"},
			Destination: &cfg.ServiceUser.Username,
		},
		&cli.Int64Flag{
			Name:        "service-user-uid",
			Value:       cfg.ServiceUser.UID,
			Usage:       "uid of the internal service user (required on EOS)",
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UID"},
			Destination: &cfg.ServiceUser.UID,
		},
		&cli.Int64Flag{
			Name:        "service-user-gid",
			Value:       cfg.ServiceUser.GID,
			Usage:       "gid of the internal service user (required on EOS)",
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_GID"},
			Destination: &cfg.ServiceUser.GID,
		},
		&cli.Int64Flag{
			Name:        "uid-index-lower-bound",
			Value:       cfg.Index.UID.Lower,
			Usage:       "define a starting point for the account UID",
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.UID.Lower,
		},
		&cli.Int64Flag{
			Name:        "gid-index-lower-bound",
			Value:       cfg.Index.GID.Lower,
			Usage:       "define a starting point for the account GID",
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.GID.Lower,
		},
		&cli.Int64Flag{
			Name:        "uid-index-upper-bound",
			Value:       cfg.Index.UID.Upper,
			Usage:       "define an ending point for the account UID",
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.UID.Upper,
		},
		&cli.Int64Flag{
			Name:        "gid-index-upper-bound",
			Value:       cfg.Index.GID.Upper,
			Usage:       "define an ending point for the account GID",
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.GID.Upper,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}

// UpdateAccountWithConfig applies update command flags to cfg
func UpdateAccountWithConfig(cfg *config.Config, a *accounts.Account) []cli.Flag {
	if a.PasswordProfile == nil {
		a.PasswordProfile = &accounts.PasswordProfile{}
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.BoolFlag{
			Name:        "enabled",
			Usage:       "Enable the account",
			Destination: &a.AccountEnabled,
		},
		&cli.StringFlag{
			Name:        "displayname",
			Usage:       "Set the displayname for the account",
			Destination: &a.DisplayName,
		},
		&cli.StringFlag{
			Name:        "preferred-name",
			Usage:       "Set the preferred-name for the account",
			Destination: &a.PreferredName,
		},
		&cli.StringFlag{
			Name:        "on-premises-sam-account-name",
			Usage:       "Set the on-premises-sam-account-name",
			Destination: &a.OnPremisesSamAccountName,
		},
		&cli.Int64Flag{
			Name:        "uidnumber",
			Usage:       "Set the uidnumber for the account",
			Destination: &a.UidNumber,
		},
		&cli.Int64Flag{
			Name:        "gidnumber",
			Usage:       "Set the gidnumber for the account",
			Destination: &a.GidNumber,
		},
		&cli.StringFlag{
			Name:        "mail",
			Usage:       "Set the mail for the account",
			Destination: &a.Mail,
		},
		&cli.StringFlag{
			Name:        "description",
			Usage:       "Set the description for the account",
			Destination: &a.Description,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "Set the password for the account",
			Destination: &a.PasswordProfile.Password,
			// TODO read password from ENV?
		},
		&cli.StringSliceFlag{
			Name:  "password-policies",
			Usage: "Possible policies: DisableStrongPassword, DisablePasswordExpiration",
		},
		&cli.BoolFlag{
			Name:        "force-password-change",
			Usage:       "Force password change on next sign-in",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignIn,
		},
		&cli.BoolFlag{
			Name:        "force-password-change-mfa",
			Usage:       "Force password change on next sign-in with mfa",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignInWithMfa,
		},
	}
}

// AddAccountWithConfig applies create command flags to cfg
func AddAccountWithConfig(cfg *config.Config, a *accounts.Account) []cli.Flag {
	if a.PasswordProfile == nil {
		a.PasswordProfile = &accounts.PasswordProfile{}
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.BoolFlag{
			Name:        "enabled",
			Usage:       "Enable the account",
			Destination: &a.AccountEnabled,
		},
		&cli.StringFlag{
			Name:        "displayname",
			Usage:       "Set the displayname for the account",
			Destination: &a.DisplayName,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "Username will be written to preferred-name and on_premises_sam_account_name",
		},
		&cli.StringFlag{
			Name:        "preferred-name",
			Usage:       "Set the preferred-name for the account",
			Destination: &a.PreferredName,
		},
		&cli.StringFlag{
			Name:        "on-premises-sam-account-name",
			Usage:       "Set the on-premises-sam-account-name",
			Destination: &a.OnPremisesSamAccountName,
		},
		&cli.Int64Flag{
			Name:        "uidnumber",
			Usage:       "Set the uidnumber for the account",
			Destination: &a.UidNumber,
		},
		&cli.Int64Flag{
			Name:        "gidnumber",
			Usage:       "Set the gidnumber for the account",
			Destination: &a.GidNumber,
		},
		&cli.StringFlag{
			Name:        "mail",
			Usage:       "Set the mail for the account",
			Destination: &a.Mail,
		},
		&cli.StringFlag{
			Name:        "description",
			Usage:       "Set the description for the account",
			Destination: &a.Description,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "Set the password for the account",
			Destination: &a.PasswordProfile.Password,
			// TODO read password from ENV?
		},
		&cli.StringSliceFlag{
			Name:  "password-policies",
			Usage: "Possible policies: DisableStrongPassword, DisablePasswordExpiration",
		},
		&cli.BoolFlag{
			Name:        "force-password-change",
			Usage:       "Force password change on next sign-in",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignIn,
		},
		&cli.BoolFlag{
			Name:        "force-password-change-mfa",
			Usage:       "Force password change on next sign-in with mfa",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignInWithMfa,
		},
	}
}

// ListAccountsWithConfig applies list command flags to cfg
func ListAccountsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

// RemoveAccountWithConfig applies remove command flags to cfg
func RemoveAccountWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

// InspectAccountWithConfig applies inspect command flags to cfg
func InspectAccountWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       cfg.GRPC.Namespace,
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       cfg.Server.Name,
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

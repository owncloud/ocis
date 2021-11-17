package flagset

import (
	"github.com/owncloud/ocis/accounts/pkg/config"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/urfave/cli/v2"
)

// UpdateAccountWithConfig applies update command flags to cfg
func UpdateAccountWithConfig(cfg *config.Config, a *accounts.Account) []cli.Flag {
	if a.PasswordProfile == nil {
		a.PasswordProfile = &accounts.PasswordProfile{}
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "accounts"),
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
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "accounts"),
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
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "accounts"),
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
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "accounts"),
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
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "accounts"),
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

//go:build !simple
// +build !simple

package command

import (
	"github.com/owncloud/ocis/accounts/pkg/command"
	svcconfig "github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/flagset"
	"github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis/pkg/register"
	"github.com/owncloud/ocis/ocis/pkg/version"
	"github.com/urfave/cli/v2"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "accounts",
		Usage:    "Start accounts server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Accounts),
		Subcommands: []*cli.Command{
			command.ListAccounts(cfg.Accounts),
			command.AddAccount(cfg.Accounts),
			command.UpdateAccount(cfg.Accounts),
			command.RemoveAccount(cfg.Accounts),
			command.InspectAccount(cfg.Accounts),
			command.PrintVersion(cfg.Accounts),
		},
		Before: func(ctx *cli.Context) error {
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			origCmd := command.Server(configureAccounts(cfg))
			return handleOriginalAction(c, origCmd)
		},
	}
}

func configureAccounts(cfg *config.Config) *svcconfig.Config {
	cfg.Accounts.Log.Level = cfg.Log.Level
	cfg.Accounts.Log.Pretty = cfg.Log.Pretty
	cfg.Accounts.Log.Color = cfg.Log.Color
	cfg.Accounts.Server.Version = version.String

	if cfg.Tracing.Enabled {
		cfg.Accounts.Tracing.Enabled = cfg.Tracing.Enabled
		cfg.Accounts.Tracing.Type = cfg.Tracing.Type
		cfg.Accounts.Tracing.Endpoint = cfg.Tracing.Endpoint
		cfg.Accounts.Tracing.Collector = cfg.Tracing.Collector
	}

	if cfg.TokenManager.JWTSecret != "" {
		cfg.Accounts.TokenManager.JWTSecret = cfg.TokenManager.JWTSecret
		cfg.Accounts.Repo.CS3.JWTSecret = cfg.TokenManager.JWTSecret
	}

	return cfg.Accounts
}

func init() {
	register.AddCommand(AccountsCommand)
}

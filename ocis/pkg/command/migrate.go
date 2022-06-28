package command

import (
	"context"
	"fmt"
	"os"
	"sync"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	publicregistry "github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/manager/registry"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/config"
	sharingparser "github.com/owncloud/ocis/v2/services/sharing/pkg/config/parser"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

// Migrate is the entrypoint for the Migrate command.
func Migrate(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "migrate",
		Usage:    "migrate data from an existing to another instance",
		Category: "migration",
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				fmt.Printf("%v", err)
				return err
			}

			// Parse sharing config
			cfg.Sharing.Commons = cfg.Commons
			if err := sharingparser.ParseConfig(cfg.Sharing); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			return nil
		},
		Subcommands: []*cli.Command{
			MigrateShares(cfg),
			MigratePublicShares(cfg),
		},
	}
}

func init() {
	register.AddCommand(Migrate)
}

func MigrateShares(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "shares",
		Usage: "migrates shares from the previous to the new share manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "from",
				Value: "json",
				Usage: "Share manager to export the data from",
			},
			&cli.StringFlag{
				Name:  "to",
				Value: "cs3",
				Usage: "Share manager to import the data into",
			},
		},
		Action: func(c *cli.Context) error {
			log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
			ctx := log.WithContext(context.Background())
			rcfg := revaShareConfig(cfg.Sharing)
			oldDriver := c.String("from")
			newDriver := c.String("to")
			shareChan := make(chan *collaboration.Share)
			receivedShareChan := make(chan share.ReceivedShareWithUser)

			f, ok := registry.NewFuncs[oldDriver]
			if !ok {
				log.Error().Msg("Unknown share manager type '" + oldDriver + "'")
				os.Exit(1)
			}
			oldMgr, err := f(rcfg[oldDriver].(map[string]interface{}))
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate source share manager")
				os.Exit(1)
			}
			dumpMgr, ok := oldMgr.(share.DumpableManager)
			if !ok {
				log.Error().Msg("Share manager type '" + oldDriver + "' does not support dumping its shares.")
				os.Exit(1)
			}

			f, ok = registry.NewFuncs[newDriver]
			if !ok {
				log.Error().Msg("Unknown share manager type '" + newDriver + "'")
				os.Exit(1)
			}
			newMgr, err := f(rcfg[newDriver].(map[string]interface{}))
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate destination share manager")
				os.Exit(1)
			}
			loadMgr, ok := newMgr.(share.LoadableManager)
			if !ok {
				log.Error().Msg("Share manager type '" + newDriver + "' does not support loading a shares dump.")
				os.Exit(1)
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				log.Info().Msg("Migrating shares...")
				err = loadMgr.Load(ctx, shareChan, receivedShareChan)
				log.Info().Msg("Finished migrating shares.")
				if err != nil {
					log.Error().Err(err).Msg("Error while loading shares")
					os.Exit(1)
				}
				wg.Done()
			}()
			go func() {
				err = dumpMgr.Dump(ctx, shareChan, receivedShareChan)
				if err != nil {
					log.Error().Err(err).Msg("Error while dumping shares")
					os.Exit(1)
				}
				close(shareChan)
				close(receivedShareChan)
				wg.Done()
			}()
			wg.Wait()
			return nil
		},
	}
}

func MigratePublicShares(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "publicshares",
		Usage: "migrates public shares from the previous to the new public share manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "from",
				Value: "json",
				Usage: "Public share manager to export the data from",
			},
			&cli.StringFlag{
				Name:  "to",
				Value: "cs3",
				Usage: "Public share manager to import the data into",
			},
		},
		Action: func(c *cli.Context) error {
			log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
			ctx := log.WithContext(context.Background())

			rcfg := revaPublicShareConfig(cfg.Sharing)
			oldDriver := c.String("from")
			newDriver := c.String("to")
			shareChan := make(chan *publicshare.WithPassword)

			f, ok := publicregistry.NewFuncs[oldDriver]
			if !ok {
				log.Error().Msg("Unknown public share manager type '" + oldDriver + "'")
				os.Exit(1)
			}
			oldMgr, err := f(rcfg[oldDriver].(map[string]interface{}))
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate source public share manager")
				os.Exit(1)
			}
			dumpMgr, ok := oldMgr.(publicshare.DumpableManager)
			if !ok {
				log.Error().Msg("Public share manager type '" + oldDriver + "' does not support dumping its public shares.")
				os.Exit(1)
			}

			f, ok = publicregistry.NewFuncs[newDriver]
			if !ok {
				log.Error().Msg("Unknown public share manager type '" + newDriver + "'")
				os.Exit(1)
			}
			newMgr, err := f(rcfg[newDriver].(map[string]interface{}))
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate destination public share manager")
				os.Exit(1)
			}
			loadMgr, ok := newMgr.(publicshare.LoadableManager)
			if !ok {
				log.Error().Msg("Public share manager type '" + newDriver + "' does not support loading a public shares dump.")
				os.Exit(1)
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				log.Info().Msg("Migrating public shares...")
				err = loadMgr.Load(ctx, shareChan)
				log.Info().Msg("Finished migrating public shares.")
				if err != nil {
					log.Error().Err(err).Msg("Error while loading public shares")
					os.Exit(1)
				}
				wg.Done()
			}()
			go func() {
				err = dumpMgr.Dump(ctx, shareChan)
				if err != nil {
					log.Error().Err(err).Msg("Error while dumping public shares")
					os.Exit(1)
				}
				close(shareChan)
				wg.Done()
			}()
			wg.Wait()
			return nil
		},
	}
}

func revaShareConfig(cfg *sharing.Config) map[string]interface{} {
	return map[string]interface{}{
		"json": map[string]interface{}{
			"file":         cfg.UserSharingDrivers.JSON.File,
			"gateway_addr": cfg.Reva.Address,
		},
		"sql": map[string]interface{}{ // cernbox sql
			"db_username":                   cfg.UserSharingDrivers.SQL.DBUsername,
			"db_password":                   cfg.UserSharingDrivers.SQL.DBPassword,
			"db_host":                       cfg.UserSharingDrivers.SQL.DBHost,
			"db_port":                       cfg.UserSharingDrivers.SQL.DBPort,
			"db_name":                       cfg.UserSharingDrivers.SQL.DBName,
			"password_hash_cost":            cfg.UserSharingDrivers.SQL.PasswordHashCost,
			"enable_expired_shares_cleanup": cfg.UserSharingDrivers.SQL.EnableExpiredSharesCleanup,
			"janitor_run_interval":          cfg.UserSharingDrivers.SQL.JanitorRunInterval,
		},
		"owncloudsql": map[string]interface{}{
			"gateway_addr":     cfg.Reva.Address,
			"storage_mount_id": cfg.UserSharingDrivers.OwnCloudSQL.UserStorageMountID,
			"db_username":      cfg.UserSharingDrivers.OwnCloudSQL.DBUsername,
			"db_password":      cfg.UserSharingDrivers.OwnCloudSQL.DBPassword,
			"db_host":          cfg.UserSharingDrivers.OwnCloudSQL.DBHost,
			"db_port":          cfg.UserSharingDrivers.OwnCloudSQL.DBPort,
			"db_name":          cfg.UserSharingDrivers.OwnCloudSQL.DBName,
		},
		"cs3": map[string]interface{}{
			"gateway_addr":        cfg.UserSharingDrivers.CS3.ProviderAddr,
			"provider_addr":       cfg.UserSharingDrivers.CS3.ProviderAddr,
			"service_user_id":     cfg.UserSharingDrivers.CS3.SystemUserID,
			"service_user_idp":    cfg.UserSharingDrivers.CS3.SystemUserIDP,
			"machine_auth_apikey": cfg.UserSharingDrivers.CS3.SystemUserAPIKey,
		},
	}
}

func revaPublicShareConfig(cfg *sharing.Config) map[string]interface{} {
	return map[string]interface{}{
		"json": map[string]interface{}{
			"file":         cfg.PublicSharingDrivers.JSON.File,
			"gateway_addr": cfg.Reva.Address,
		},
		"sql": map[string]interface{}{
			"db_username":                   cfg.PublicSharingDrivers.SQL.DBUsername,
			"db_password":                   cfg.PublicSharingDrivers.SQL.DBPassword,
			"db_host":                       cfg.PublicSharingDrivers.SQL.DBHost,
			"db_port":                       cfg.PublicSharingDrivers.SQL.DBPort,
			"db_name":                       cfg.PublicSharingDrivers.SQL.DBName,
			"password_hash_cost":            cfg.PublicSharingDrivers.SQL.PasswordHashCost,
			"enable_expired_shares_cleanup": cfg.PublicSharingDrivers.SQL.EnableExpiredSharesCleanup,
			"janitor_run_interval":          cfg.PublicSharingDrivers.SQL.JanitorRunInterval,
		},
		"cs3": map[string]interface{}{
			"gateway_addr":        cfg.PublicSharingDrivers.CS3.ProviderAddr,
			"provider_addr":       cfg.PublicSharingDrivers.CS3.ProviderAddr,
			"service_user_id":     cfg.PublicSharingDrivers.CS3.SystemUserID,
			"service_user_idp":    cfg.PublicSharingDrivers.CS3.SystemUserIDP,
			"machine_auth_apikey": cfg.PublicSharingDrivers.CS3.SystemUserAPIKey,
		},
	}
}

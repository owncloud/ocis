package command

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/manager/registry"
	sharing "github.com/owncloud/ocis/v2/extensions/sharing/pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// Migrate is the entrypoint for the Migrate command.
func Migrate(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "migrate",
		Usage:    "migrate data from an existing instance to a new version",
		Category: "migration",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				s, _ := json.MarshalIndent(cfg, "", "	")
				fmt.Print(string(s))
				fmt.Printf("%v", err)
				return err
			}
			return nil
		},
		Subcommands: []*cli.Command{
			MigrateShares(cfg),
		},
	}
}

func init() {
	register.AddCommand(Migrate)
}

func MigrateShares(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "shares",
		Usage: "migrates shares from the previous to the new ",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			rcfg := revaConfig(cfg.Sharing)
			oldDriver := "json"
			newDriver := "cs3"
			shareChan := make(chan *collaboration.Share)
			receivedShareChan := make(chan share.ReceivedShareDump)

			f, ok := registry.NewFuncs[oldDriver]
			if !ok {
				fmt.Println("Unknown share manager type '" + oldDriver + "'")
				os.Exit(1)
			}
			oldMgr, err := f(rcfg[oldDriver].(map[string]interface{}))
			if err != nil {
				fmt.Println("failed to initiate source share manager", err)
				os.Exit(1)
			}
			if _, ok := oldMgr.(share.DumpableManager); !ok {
				fmt.Println("Share manager type '" + oldDriver + "' does not support migration.")
				os.Exit(1)
			}

			f, ok = registry.NewFuncs[newDriver]
			if !ok {
				fmt.Println("Unknown share manager type '" + oldDriver + "'")
				os.Exit(1)
			}
			newMgr, err := f(rcfg[newDriver].(map[string]interface{}))
			if err != nil {
				fmt.Println("failed to initiate source share manager", err)
				os.Exit(1)
			}
			if _, ok := newMgr.(share.LoadableManager); !ok {
				fmt.Println("Share manager type '" + newDriver + "' does not support migration.")
				os.Exit(1)
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				fmt.Println("Loading...")
				err = newMgr.(share.LoadableManager).Load(shareChan, receivedShareChan)
				fmt.Println("Finished loading...")
				if err != nil {
					fmt.Println("Error while loading shares", err)
					os.Exit(1)
				}
				wg.Done()
			}()
			go func() {
				err = oldMgr.(share.DumpableManager).Dump(shareChan, receivedShareChan)
				if err != nil {
					fmt.Println("Error while dumping shares", err)
					os.Exit(1)
				}
				wg.Done()
			}()
			wg.Wait()
			return nil
		},
	}
}

func revaConfig(cfg *sharing.Config) map[string]interface{} {
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

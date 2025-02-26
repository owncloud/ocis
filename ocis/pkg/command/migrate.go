package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	publicregistry "github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/manager/jsoncs3"
	"github.com/cs3org/reva/v2/pkg/share/manager/jsoncs3/providercache"
	"github.com/cs3org/reva/v2/pkg/share/manager/jsoncs3/shareid"
	"github.com/cs3org/reva/v2/pkg/share/manager/registry"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/timemanager"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/migrator"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	"github.com/mitchellh/mapstructure"
	tw "github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	oclog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	mregistry "github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/config"
	sharingparser "github.com/owncloud/ocis/v2/services/sharing/pkg/config/parser"
	"github.com/urfave/cli/v2"
)

// Migrate is the entrypoint for the Migrate command.
func Migrate(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "migrate",
		Usage:    "migrate data from an existing to another instance",
		Category: "migration",
		Subcommands: []*cli.Command{
			MigrateDecomposedfs(cfg),
			MigrateShares(cfg),
			MigratePublicShares(cfg),
			RebuildJSONCS3Indexes(cfg),
		},
	}
}

func init() {
	register.AddCommand(Migrate)
}

// RebuildJSONCS3Indexes rebuilds the share indexes from the shares json
func RebuildJSONCS3Indexes(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "rebuild-jsoncs3-indexes",
		Usage:       "rebuild the share indexes from the shares json",
		Subcommands: []*cli.Command{},
		Flags:       []cli.Flag{},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse sharing config
			cfg.Sharing.Commons = cfg.Commons
			return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
		},
		Action: func(c *cli.Context) error {
			log := logger()
			ctx := log.WithContext(context.Background())
			rcfg := revaShareConfig(cfg.Sharing)

			// Initialize registry to make service lookup work
			_ = mregistry.GetRegistry()

			// Get a jsoncs3 manager to operate its caches
			type config struct {
				GatewayAddr       string `mapstructure:"gateway_addr"`
				MaxConcurrency    int    `mapstructure:"max_concurrency"`
				ProviderAddr      string `mapstructure:"provider_addr"`
				ServiceUserID     string `mapstructure:"service_user_id"`
				ServiceUserIdp    string `mapstructure:"service_user_idp"`
				MachineAuthAPIKey string `mapstructure:"machine_auth_apikey"`
			}
			conf := &config{}
			if err := mapstructure.Decode(rcfg["jsoncs3"], conf); err != nil {
				err = errors.Wrap(err, "error creating a new manager")
				return err
			}
			s, err := metadata.NewCS3Storage(conf.GatewayAddr, conf.ProviderAddr, conf.ServiceUserID, conf.ServiceUserIdp, conf.MachineAuthAPIKey)
			if err != nil {
				return err
			}
			err = s.Init(ctx, "jsoncs3-share-manager-metadata")
			if err != nil {
				return err
			}
			gatewaySelector, err := pool.GatewaySelector(conf.GatewayAddr)
			if err != nil {
				return err
			}
			mgr, err := jsoncs3.New(s, gatewaySelector, 0, nil, 1)
			if err != nil {
				return err
			}

			// Rebuild indexes
			errorsOccured := false
			storages, err := s.ReadDir(ctx, "storages")
			if err != nil {
				return err
			}
			for iStorage, storage := range storages {
				fmt.Printf("Scanning storage %s (%d/%d)\n", storage, iStorage+1, len(storages))
				spaces, err := s.ReadDir(ctx, filepath.Join("storages", storage))
				if err != nil {
					fmt.Printf("failed! (%s)\n", err.Error())
					errorsOccured = true
					continue
				}

				for iSpace, space := range spaces {
					fmt.Printf("  Rebuilding space '%s' %d/%d...", strings.TrimSuffix(space, ".json"), iSpace+1, len(spaces))

					spaceBlob, err := s.SimpleDownload(ctx, filepath.Join("storages", storage, space))
					if err != nil {
						fmt.Printf(" failed! (%s)\n", err.Error())
						errorsOccured = true
						continue
					}
					shares := &providercache.Shares{}
					err = json.Unmarshal(spaceBlob, shares)
					if err != nil {
						fmt.Printf(" failed! (%s)\n", err.Error())
						errorsOccured = true
						continue
					}
					for _, share := range shares.Shares {
						err = mgr.Cache.Add(ctx, share.ResourceId.StorageId, share.ResourceId.SpaceId, share.Id.OpaqueId, share)
						if err != nil {
							fmt.Printf(" adding share '%s' to the cache failed! (%s)\n", share.Id.OpaqueId, err.Error())
							errorsOccured = true
						}
						err = mgr.CreatedCache.Add(ctx, share.Creator.OpaqueId, share.Id.OpaqueId)
						if err != nil {
							fmt.Printf(" adding share '%s' to the created cache failed! (%s)\n", share.Id.OpaqueId, err.Error())
							errorsOccured = true
						}

						spaceId := share.ResourceId.StorageId + shareid.IDDelimiter + share.ResourceId.SpaceId
						switch share.Grantee.Type {
						case provider.GranteeType_GRANTEE_TYPE_USER:
							userid := share.Grantee.GetUserId().GetOpaqueId()
							existingState, err := mgr.UserReceivedStates.Get(ctx, userid, spaceId, share.Id.OpaqueId)
							if err != nil {
								fmt.Printf(" retrieving current state of received share '%s' from the user cache failed! (%s)\n", share.Id.OpaqueId, err.Error())
								errorsOccured = true
							} else if existingState == nil {
								rs := &collaboration.ReceivedShare{
									Share: share,
									State: collaboration.ShareState_SHARE_STATE_PENDING,
								}
								err := mgr.UserReceivedStates.Add(ctx, userid, spaceId, rs)
								if err != nil {
									fmt.Printf(" adding share '%s' to the user cache failed! (%s)\n", share.Id.OpaqueId, err.Error())
									errorsOccured = true
								}
							}
						case provider.GranteeType_GRANTEE_TYPE_GROUP:
							groupid := share.Grantee.GetGroupId().GetOpaqueId()
							err := mgr.GroupReceivedCache.Add(ctx, groupid, spaceId)
							if err != nil {
								fmt.Printf(" adding share '%s' to the group cache failed! (%s)\n", share.Id.OpaqueId, err.Error())
								errorsOccured = true
							}
						}
					}
					fmt.Printf("  done\n")
				}
				fmt.Printf("done\n")
			}
			if errorsOccured {
				return errors.New("There were errors. Please review the logs or try again.")
			}

			return nil
		},
	}
}

// MigrateDecomposedfs is the entrypoint for the decomposedfs migrate command
func MigrateDecomposedfs(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "decomposedfs",
		Usage: "run a decomposedfs migration",
		Subcommands: []*cli.Command{
			ListDecomposedfsMigrations(cfg),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "direction",
				Aliases: []string{"d"},
				Value:   "migrate",
				Usage:   "direction of the migration to run ('migrate' or 'rollback')",
			},
			&cli.StringFlag{
				Name:    "migration",
				Aliases: []string{"m"},
				Value:   "",
				Usage:   "ID of the migration to run",
			},
			&cli.StringFlag{
				Name:     "root",
				Aliases:  []string{"r"},
				Required: true,
				Usage:    "Path to the root directory of the decomposedfs",
			},
		},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			log := logger()
			rootFlag := c.String("root")
			bod := lookup.DetectBackendOnDisk(rootFlag)
			backend := backend(rootFlag, bod)
			lu := lookup.New(backend, &options.Options{
				Root:            rootFlag,
				MetadataBackend: bod,
			}, &timemanager.Manager{})

			m := migrator.New(lu, log)

			err := m.RunMigration(c.String("migration"), c.String("direction") == "down")
			if err != nil {
				log.Error().Err(err).Msg("failed")
				return err
			}

			return nil
		},
	}
}

// ListDecomposedfsMigrations is the entrypoint for the decomposedfs list migrations command
func ListDecomposedfsMigrations(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list decomposedfs migrations",
		Action: func(c *cli.Context) error {
			rootFlag := c.String("root")
			bod := lookup.DetectBackendOnDisk(rootFlag)
			backend := backend(rootFlag, bod)
			lu := lookup.New(backend, &options.Options{
				Root:            rootFlag,
				MetadataBackend: bod,
			}, &timemanager.Manager{})

			m := migrator.New(lu, logger())
			migrationStates, err := m.Migrations()
			if err != nil {
				return err
			}

			migrations := []string{}
			for m := range migrationStates {
				migrations = append(migrations, m)
			}
			sort.Strings(migrations)

			table := tw.NewWriter(os.Stdout)
			table.SetHeader([]string{"Migration", "State", "Message"})
			table.SetAutoFormatHeaders(false)
			for _, migration := range migrations {
				table.Append([]string{migration, migrationStates[migration].State, migrationStates[migration].Message})
			}
			table.Render()

			return nil
		},
	}
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
				Value: "jsoncs3",
				Usage: "Share manager to import the data into",
			},
		},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse sharing config
			cfg.Sharing.Commons = cfg.Commons
			return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
		},
		Action: func(c *cli.Context) error {
			log := logger()
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
				Value: "jsoncs3",
				Usage: "Public share manager to import the data into",
			},
		},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse sharing config
			cfg.Sharing.Commons = cfg.Commons
			return configlog.ReturnError(sharingparser.ParseConfig(cfg.Sharing))
		},
		Action: func(c *cli.Context) error {
			log := logger()
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
		"jsoncs3": map[string]interface{}{
			"gateway_addr":        cfg.Reva.Address,
			"provider_addr":       cfg.UserSharingDrivers.JSONCS3.ProviderAddr,
			"service_user_id":     cfg.UserSharingDrivers.JSONCS3.SystemUserID,
			"service_user_idp":    cfg.UserSharingDrivers.JSONCS3.SystemUserIDP,
			"machine_auth_apikey": cfg.UserSharingDrivers.JSONCS3.SystemUserAPIKey,
		},
	}
}

func revaPublicShareConfig(cfg *sharing.Config) map[string]interface{} {
	return map[string]interface{}{
		"json": map[string]interface{}{
			"file":         cfg.PublicSharingDrivers.JSON.File,
			"gateway_addr": cfg.Reva.Address,
		},
		"jsoncs3": map[string]interface{}{
			"gateway_addr":        cfg.Reva.Address,
			"provider_addr":       cfg.PublicSharingDrivers.JSONCS3.ProviderAddr,
			"service_user_id":     cfg.PublicSharingDrivers.JSONCS3.SystemUserID,
			"service_user_idp":    cfg.PublicSharingDrivers.JSONCS3.SystemUserIDP,
			"machine_auth_apikey": cfg.PublicSharingDrivers.JSONCS3.SystemUserAPIKey,
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

func logger() *zerolog.Logger {
	log := oclog.NewLogger(
		oclog.Name("migrate"),
		oclog.Level("info"),
		oclog.Pretty(true),
		oclog.Color(true)).Logger
	return &log
}

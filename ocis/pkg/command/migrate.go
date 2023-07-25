package command

import (
	"context"
	"fmt"
	"os"
	"sync"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/publicshare"
	publicregistry "github.com/cs3org/reva/v2/pkg/publicshare/manager/registry"
	"github.com/cs3org/reva/v2/pkg/share"
	"github.com/cs3org/reva/v2/pkg/share/manager/registry"
	storageregistry "github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	oclog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/config"
	sharingparser "github.com/owncloud/ocis/v2/services/sharing/pkg/config/parser"
	storageparser "github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
	"github.com/urfave/cli/v2"
)

// Migrate is the entrypoint for the Migrate command.
func Migrate(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "migrate",
		Usage:    "migrate data from an existing to another instance",
		Category: "migration",
		Subcommands: []*cli.Command{
			MigrateSpace(cfg),
			MigrateShares(cfg),
			MigratePublicShares(cfg),
		},
	}
}

func init() {
	register.AddCommand(Migrate)
}

func MigrateSpace(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "space",
		Usage: "migrates a space from one storage provider to another using the storage interface",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "from",
				Value: "owncloudsql",
				Usage: "source type to read files from",
			},
			&cli.StringFlag{
				Name:  "to",
				Value: "ocis",
				Usage: "target type to write files to",
			},
			&cli.StringFlag{
				Name:  "space",
				Value: "",
				Usage: "space to migrate",
			},
			&cli.StringFlag{
				Name:  "username",
				Value: "",
				Usage: "username to migrate",
			},
			// TODO add continue-on-error option
		},
		Before: func(c *cli.Context) error {
			// Parse base config
			if err := parser.ParseConfig(cfg, true); err != nil {
				return configlog.ReturnError(err)
			}

			// Parse storage config
			cfg.StorageUsers.Commons = cfg.Commons

			// TODO hopefully we don't have to actually talk to other services
			cfg.StorageUsers.Commons.TokenManager.JWTSecret = "unused"
			cfg.StorageUsers.MountID = "unused"
			return configlog.ReturnError(storageparser.ParseConfig(cfg.StorageUsers))
		},
		Action: func(c *cli.Context) error {
			log := oclog.LoggerFromConfig("migrate", cfg.Log)
			ctx := log.WithContext(context.Background())

			rcfg := revaconfig.StorageProviderDrivers(cfg.StorageUsers)

			oldDriver := c.String("from")
			newDriver := c.String("to")
			space := c.String("space")
			username := c.String("username")
			u := &userv1beta1.User{
				Username: username, // needed for owncloudsql
			}
			ctx = revactx.ContextSetUser(ctx, u)

			f, ok := storageregistry.NewFuncs[oldDriver]
			if !ok {
				log.Error().Msg("Unknown source storage type '" + oldDriver + "'")
				os.Exit(1)
			}
			sourceFS, err := f(rcfg[oldDriver].(map[string]interface{}), nil)
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate source storage driver")
				os.Exit(1)
			}

			f, ok = storageregistry.NewFuncs[newDriver]
			if !ok {
				log.Error().Msg("Unknown target storage type '" + newDriver + "'")
				os.Exit(1)
			}

			// disable propagation during the import run
			//rcfg[newDriver]
			targetFS, err := f(rcfg[newDriver].(map[string]interface{}), nil)
			if err != nil {
				log.Error().Err(err).Msg("failed to initiate destination public share manager")
				os.Exit(1)
			}

			// 1. create space
			spaces, err := sourceFS.ListStorageSpaces(ctx, []*providerv1beta1.ListStorageSpacesRequest_Filter{{
				Type: providerv1beta1.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &providerv1beta1.ListStorageSpacesRequest_Filter_Id{
					Id: &providerv1beta1.StorageSpaceId{OpaqueId: space},
				},
			}}, true)
			if err != nil {
				return err
			}

			switch len(spaces) {
			case 0:
				return fmt.Errorf("unknown space %s", space)
			case 1:
				// cool
				res, err := targetFS.CreateStorageSpace(ctx, &providerv1beta1.CreateStorageSpaceRequest{Type: "personal"})
				if err != nil {
					return err
				}
				if res.Status.Code != cs3rpc.Code_CODE_OK {
					return fmt.Errorf("could not create space: (%d) %s", res.Status.Code, res.Status.Message)
				}
			default:
				return fmt.Errorf("expected on space with id %s, found %d", space, len(spaces))
			}

			// FIXME ... owncloudsql reads metadata from the database. This is probably the wrong approach to do an offline migration

			// 2. recreate tree

			children, err := sourceFS.ListFolder(ctx, &providerv1beta1.Reference{ResourceId: spaces[0].Root}, nil, nil) // TODO list all metadata keys with *?
			if err != nil {
				return err
			}

			for _, child := range children {
				childRef := &providerv1beta1.Reference{ResourceId: spaces[0].Root, Path: child.Name}
				switch child.Type {
				case providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER:
					err := targetFS.CreateDir(ctx, childRef)
					if err != nil {
						return err
					}
					err = targetFS.SetArbitraryMetadata(ctx, childRef, &providerv1beta1.ArbitraryMetadata{
						Metadata: map[string]string{
							"mtime": fmt.Sprintf("%d.%d", child.Mtime.Seconds, child.Mtime.Nanos),
							// TODO set etag? maybe not because we want to enforce updateing the fileid?
						},
					})
					if err != nil {
						return err
					}
					// TODO recursively descend

				case providerv1beta1.ResourceType_RESOURCE_TYPE_FILE:
					stream, err := sourceFS.Download(ctx, &providerv1beta1.Reference{ResourceId: child.Id})
					if err != nil {
						return err
					}
					_, err = targetFS.Upload(ctx, childRef, stream, nil)
					if err != nil {
						return err
					}
				default:
					// ignore
					continue
				}
			}

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
				Value: "cs3",
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
			log := oclog.LoggerFromConfig("migrate", cfg.Log)
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
			log := oclog.LoggerFromConfig("migrate", cfg.Log)
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
			"gateway_addr":        cfg.UserSharingDrivers.JSONCS3.ProviderAddr,
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

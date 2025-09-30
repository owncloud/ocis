package command

import (
	"context"
	"fmt"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	zlog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/urfave/cli/v2"
)

const (
	_spaceTypePersonal = "personal"
	_spaceStateTrashed = "trashed"
)

// Spaces wraps space-related sub-commands.
func Spaces(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "spaces",
		Usage: "manage spaces",
		Subcommands: []*cli.Command{
			PurgeDisabledSpaces(cfg),
		},
	}
}

// PurgeDisabledSpaces cli command purges disabled personal spaces.
//
// Note: The current implementation uses space modification time as a heuristic
// to determine if a space should be purged.
// The default value of the dry-run flag is true to avoid accidental purging.
func PurgeDisabledSpaces(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-disabled",
		Usage: "Purge disabled personal spaces that exceed the retention period",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Only show what would be purged without actually purging",
				Value:   true, // default must be true to avoid accidental purging!
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:     "space-id",
				Usage:    "Specific space ID to purge (omit to process all spaces)",
				Required: false,
				Aliases:  []string{"s"},
			},
			&cli.StringFlag{
				Name:     "retention-period",
				Usage:    "Retention period for disabled spaces (e.g., '336h', '14d')",
				Value:    "336h",
				Required: false,
				Aliases:  []string{"r"},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose logging",
				Value:   false,
				Aliases: []string{"v"},
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := cliLogger(c.Bool("verbose"))
			dryRun := c.Bool("dry-run")
			spaceID := c.String("space-id")
			retentionPeriod := c.String("retention-period")
			verbose := c.Bool("verbose")

			// Parse retention period
			duration, err := time.ParseDuration(retentionPeriod)
			if err != nil {
				return fmt.Errorf("invalid retention period format: %w", err)
			}

			log.Info().Msgf("Starting space purge process (dry-run: %v, retention: %v)", dryRun, retentionPeriod)

			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				return fmt.Errorf("error selecting gateway client: %w", err)
			}

			ctx, err := utils.GetServiceUserContext(cfg.ServiceAccount.ServiceAccountID, client, cfg.ServiceAccount.ServiceAccountSecret)
			if err != nil {
				return fmt.Errorf("could not get service user context: %w", err)
			}

			var spacesToProcess []*provider.StorageSpace
			if spaceID != "" {
				// Get specific space
				space, err := getSpaceByID(ctx, client, spaceID)
				if err != nil {
					return fmt.Errorf("failed to get space %s: %w", spaceID, err)
				}
				spacesToProcess = []*provider.StorageSpace{space}
			} else {
				// Get all personal spaces
				spaces, err := getAllPersonalSpaces(ctx, client)
				if err != nil {
					return fmt.Errorf("failed to list spaces: %w", err)
				}
				spacesToProcess = spaces
			}

			if verbose {
				log.Info().Msgf("Found %d spaces to process", len(spacesToProcess))
			}

			purgedCount := 0
			for _, space := range spacesToProcess {
				shouldPurge := shouldPurgeSpace(space, duration, log)
				if !shouldPurge {
					continue
				}

				mtime := utils.TSToTime(space.GetMtime()).Truncate(time.Second)
				if dryRun {
					log.Info().Msgf("Would purge space: %s (disabled at: %v)", space.Id.OpaqueId, mtime)
					purgedCount++
				} else {
					err := purgeSpace(ctx, client, space)
					if err != nil {
						log.Error().Err(err).Msgf("Failed to purge space %s (disabled at: %v)", space.Id.OpaqueId, mtime)
						continue
					}
					log.Info().Msgf("Purged space: %s (disabled at: %v)", space.Id.OpaqueId, mtime)
					purgedCount++
				}
			}

			if dryRun {
				log.Info().Msgf("Dry run completed. Would purge %d spaces", purgedCount)
			} else {
				log.Info().Msgf("Purge completed. Purged %d spaces", purgedCount)
			}

			return nil
		},
	}
}

// getSpaceByID retrieves a specific space by ID
func getSpaceByID(ctx context.Context, client gateway.GatewayAPIClient, spaceID string) (*provider.StorageSpace, error) {
	req := &provider.ListStorageSpacesRequest{
		Opaque: utils.AppendPlainToOpaque(nil, _spaceStateTrashed, _spaceStateTrashed),
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &provider.ListStorageSpacesRequest_Filter_Id{
					Id: &provider.StorageSpaceId{
						OpaqueId: spaceID,
					},
				},
			},
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: _spaceTypePersonal,
				},
			},
		},
	}

	resp, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("failed to list spaces: %s", resp.Status.Code)
	}

	if len(resp.StorageSpaces) == 0 {
		return nil, fmt.Errorf("space not found: %s", spaceID)
	}

	return resp.StorageSpaces[0], nil
}

// getAllPersonalSpaces retrieves all personal spaces
func getAllPersonalSpaces(ctx context.Context, client gateway.GatewayAPIClient) ([]*provider.StorageSpace, error) {
	req := &provider.ListStorageSpacesRequest{
		Opaque: utils.AppendPlainToOpaque(nil, _spaceStateTrashed, _spaceStateTrashed),
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: _spaceTypePersonal,
				},
			},
		},
	}

	resp, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("failed to list spaces: %s", resp.Status.Code)
	}

	return resp.StorageSpaces, nil
}

// shouldPurgeSpace determines if a space should be purged based on retention period.
func shouldPurgeSpace(space *provider.StorageSpace, retentionPeriod time.Duration, log zlog.Logger) bool {
	// Only process personal spaces
	if space.SpaceType != _spaceTypePersonal {
		return false
	}
	// Only process trashed spaces
	if !utils.ExistsInOpaque(space.GetOpaque(), _spaceStateTrashed) && utils.ReadPlainFromOpaque(space.GetOpaque(), _spaceStateTrashed) != _spaceStateTrashed {
		return false
	}

	// A disabled space have a DTime attribute set, but we can't directly access it then we use the space modification time
	disabledTime := utils.TSToTime(space.GetMtime())
	cutoffTime := disabledTime.Add(retentionPeriod)
	shouldPurge := time.Now().After(cutoffTime)

	if shouldPurge {
		log.Info().Msgf("Space %s last modified at %s, cutoff is %s - should purge",
			space.GetId().GetOpaqueId(), disabledTime.Truncate(time.Second), cutoffTime.Truncate(time.Second))
	} else {
		log.Debug().Msgf("Space %s last modified at %s, cutoff is %s - within retention period",
			space.GetId().GetOpaqueId(), disabledTime.Truncate(time.Second), cutoffTime.Truncate(time.Second))
	}

	return shouldPurge
}

// purgeSpace Purge the disabled space (permanent deletion).
func purgeSpace(ctx context.Context, client gateway.GatewayAPIClient, space *provider.StorageSpace) error {
	purgeFlag := utils.AppendPlainToOpaque(nil, "purge", "")
	resp, err := client.DeleteStorageSpace(ctx, &provider.DeleteStorageSpaceRequest{
		Opaque: purgeFlag,
		Id: &provider.StorageSpaceId{
			OpaqueId: space.Id.OpaqueId,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete storage space: %w", err)
	}

	if resp.Status.Code != rpc.Code_CODE_OK {
		return fmt.Errorf("delete storage space failed: %s", resp.Status.Code)
	}

	return nil
}

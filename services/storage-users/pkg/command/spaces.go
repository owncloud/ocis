package command

import (
	"context"
	"fmt"
	"os"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/olekukonko/tablewriter"
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
	_spaceTypeProject  = "project"

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

// PurgeDisabledSpaces cli command purges disabled spaces 'personal' or 'project'.
//
// Note: The current implementation uses space modification time as a heuristic
// to determine if a space should be purged.
// The default value of the dry-run flag is true to avoid accidental purging.
func PurgeDisabledSpaces(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge",
		Usage: "Purge disabled spaces 'personal' or 'project' that exceed the retention period",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "Only show what would be purged without actually purging",
				Value:   true, // default must be true to avoid accidental purging!
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "Type of spaces to purge (e.g., 'personal', 'project')",
				Value:    _spaceTypePersonal,
				Required: true,
				Aliases:  []string{"t"},
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
			spaceType := c.String("type")
			spaceID := c.String("space-id")
			retentionPeriod := c.String("retention-period")
			if spaceType != _spaceTypePersonal && spaceType != _spaceTypeProject {
				return fmt.Errorf("invalid space type: %s", spaceType)
			}

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
				space, err := getSpaceByID(ctx, client, spaceID, spaceType)
				if err != nil {
					return fmt.Errorf("failed to get space: %w", err)
				}
				if shouldPurgeSpace(space, duration, log, spaceType) {
					spacesToProcess = []*provider.StorageSpace{space}
				}
			} else {
				// Get all spaces
				spaces, err := getAllSpaces(ctx, client, spaceType)
				if err != nil {
					return fmt.Errorf("failed to list spaces: %w", err)
				}
				for _, space := range spaces {
					if shouldPurgeSpace(space, duration, log, spaceType) {
						spacesToProcess = append(spacesToProcess, space)
					}
				}
			}

			log.Info().Msgf("Found %d spaces to process", len(spacesToProcess))
			purgedCount := 0
			if len(spacesToProcess) > 0 {
				table := tablewriter.NewTable(os.Stdout)
				table.Header("id", "space type", "space name", "disabled at")
				for _, space := range spacesToProcess {

					mtime := utils.TSToTime(space.GetMtime()).Truncate(time.Second).UTC().Format(time.RFC3339)
					spaceId := space.GetId().GetOpaqueId()
					table.Append([]string{spaceId, space.GetSpaceType(), space.GetName(), mtime})

					if dryRun {
						log.Info().Msgf("Would purge space: %s (disabled at: %v)", spaceId, mtime)
						purgedCount++
					} else {
						err := purgeSpace(ctx, client, space)
						if err != nil {
							log.Error().Err(err).Msgf("Failed to purge space %s (disabled at: %v)", spaceId, mtime)
							continue
						}
						log.Info().Msgf("Purged space: %s (disabled at: %v)", spaceId, mtime)
						purgedCount++
					}
				}
				table.Render()
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
func getSpaceByID(ctx context.Context, client gateway.GatewayAPIClient, spaceID string, spaceType string) (*provider.StorageSpace, error) {
	req := &provider.ListStorageSpacesRequest{
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
					SpaceType: spaceType,
				},
			},
		},
	}

	resp, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("failed to list spaces: %s", resp.GetStatus().GetCode())
	}

	if len(resp.GetStorageSpaces()) == 0 {
		return nil, fmt.Errorf("space not found: %s", spaceID)
	}

	return resp.GetStorageSpaces()[0], nil
}

// getAllSpaces retrieves all spaces of a given type
func getAllSpaces(ctx context.Context, client gateway.GatewayAPIClient, spaceType string) ([]*provider.StorageSpace, error) {
	req := &provider.ListStorageSpacesRequest{
		Filters: []*provider.ListStorageSpacesRequest_Filter{
			{
				Type: provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &provider.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: spaceType,
				},
			},
		},
	}

	resp, err := client.ListStorageSpaces(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("failed to list spaces: %s", resp.GetStatus().GetCode())
	}

	return resp.GetStorageSpaces(), nil
}

// shouldPurgeSpace determines if a space should be purged based on retention period.
func shouldPurgeSpace(space *provider.StorageSpace, retentionPeriod time.Duration, log zlog.Logger, spaceType string) bool {
	// Only process spaces of the given type
	if space.GetSpaceType() != spaceType {
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
			OpaqueId: space.GetId().GetOpaqueId(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete storage space: %w", err)
	}
	if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return fmt.Errorf("delete storage space failed: %s", resp.GetStatus().GetCode())
	}

	return nil
}

package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/urfave/cli/v2"
)

// Consistency is the entry point for the consistency command
func Consistency(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "consistency",
		Usage: "Check consistency of uploads.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "id",
				DefaultText: "unset",
				Usage:       "filter sessions by upload session id",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging",
				Value: false,
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			f, ok := registry.NewFuncs[cfg.Driver]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown filesystem driver '%s'\n", cfg.Driver)
				os.Exit(1)
			}
			drivers := revaconfig.StorageProviderDrivers(cfg)
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}), nil, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadSessionLister)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support listing upload sessions\n", cfg.Driver)
				os.Exit(1)
			}

			filter := buildFilter(c)
			uploads, err := managingFS.ListUploadSessions(c.Context, filter)
			if err != nil {
				return err
			}

			// We need a lookup to read nodes, but we don't have easy access to the one created inside the driver.
			// We will recreate a minimal one for read-only access.
			// This assumes the driver is decomposedfs, which is the only one we really care about for this error.
			if cfg.Driver != "ocis" { // "ocis" is the driver name for decomposedfs in ocis config
				fmt.Printf("Warning: driver is '%s', consistency check is optimized for 'ocis' (decomposedfs).\n", cfg.Driver)
			}

			// Initialize a minimal lookup
			opts, err := options.New(drivers[cfg.Driver].(map[string]interface{}))
			if err != nil {
				return err
			}
			var backend metadata.Backend
			switch opts.MetadataBackend {
			case "xattrs":
				backend = metadata.NewXattrsBackend(opts.Root, opts.FileMetadataCache)
			case "messagepack":
				backend = metadata.NewMessagePackBackend(opts.Root, opts.FileMetadataCache)
			default:
				return fmt.Errorf("unknown metadata backend %s", opts.MetadataBackend)
			}
			lu := lookup.New(backend, opts, nil)

			verbose := c.Bool("verbose")

			fmt.Printf("Checking consistency for %d upload sessions...\n", len(uploads))

			for _, u := range uploads {
				ref := u.Reference()
				spaceID := ref.GetResourceId().GetSpaceId()
				nodeID := ref.GetResourceId().GetOpaqueId()
				if verbose {
					fmt.Printf("Checking upload %s (NodeID: %s, SpaceID: %s)\n", u.ID(), nodeID, spaceID)
				}
				if nodeID == "" {
					fmt.Printf("Error: nodeID is epmpty in upload %s\n", u.ID())
					continue
				}
				if spaceID == "" {
					fmt.Printf("Error: spaceID is epmpty in upload %s\n", u.ID())
					continue
				}

				err := checkNodeConsistency(c.Context, lu, spaceID, nodeID, verbose)
				if err != nil {
					fmt.Printf("Issue found in upload %s: %v\n", u.ID(), err)
				} else if verbose {
					fmt.Printf("Upload %s seems consistent.\n", u.ID())
				}
			}

			return nil
		},
	}
}

func checkNodeConsistency(ctx context.Context, lu *lookup.Lookup, spaceID, nodeID string, verbose bool) error {
	checkMetadataMsg := "check the node metadata"

	nodePath := lu.InternalPath(spaceID, nodeID)
	// Check if node exists
	_, err := os.Stat(nodePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Error: node %s does not exist on disk", nodeID)
		}
		return err
	}

	if _, err = os.Open(lu.MetadataBackend().MetadataPath(nodePath)); err != nil {
		return fmt.Errorf("Error: failed to read metafile for node %s. %s: %s", nodeID, checkMetadataMsg, lu.MetadataBackend().MetadataPath(nodePath))
	}

	attrs, err := lu.MetadataBackend().All(ctx, nodePath)
	if err != nil {
		return fmt.Errorf("Error: failed to read attributes for node %s: %w. %s", nodeID, err, checkMetadataMsg)
	}

	if len(attrs) == 0 {
		return fmt.Errorf("Error: node %s has no attributes. %s", nodeID, checkMetadataMsg)
	}

	parentID := string(attrs[prefixes.ParentidAttr])

	if verbose {
		name := string(attrs[prefixes.NameAttr])
		fmt.Printf("node %s (%s) has parent %s\n", nodeID, name, parentID)
	}

	if parentID == "" {
		if nodeID == spaceID {
			return nil
		}
		return fmt.Errorf("Error: node %s has no %s attribute. %s", nodeID, prefixes.ParentidAttr, checkMetadataMsg)
	}

	if parentID == node.RootID {
		return nil
	}

	// Traverse upwards
	currentID := parentID
	for {
		if currentID == spaceID || currentID == node.RootID {
			break
		}

		parentPath := lu.InternalPath(spaceID, currentID)
		_, err := os.Stat(parentPath)
		if err != nil {
			if os.IsNotExist(err) {
				// Parent is missing. Check trash bin.
				matches, globErr := filepath.Glob(parentPath + node.TrashIDDelimiter + "*")
				if globErr == nil && len(matches) > 0 {
					return fmt.Errorf("Info: node parent %s is in trash bin (found %d trashed versions, e.g. %s)", currentID, len(matches), filepath.Base(matches[0]))
				}

				return fmt.Errorf("Info: node parent %s is missing from disk", currentID)
			}
			return fmt.Errorf("Info: failed to stat parent %s: %w", currentID, err)
		}

		// Read parent's parent
		if _, err = os.Open(lu.MetadataBackend().MetadataPath(parentPath)); err != nil {
			return fmt.Errorf("Error: failed to read metafile for node %s. %s: %s", parentID, checkMetadataMsg, lu.MetadataBackend().MetadataPath(parentPath))
		}

		pAttrs, pErr := lu.MetadataBackend().All(ctx, parentPath)
		if pErr != nil {
			return fmt.Errorf("Error: failed to read attributes for parent %s: %w. %s", parentID, pErr, checkMetadataMsg)
		}

		if len(pAttrs) == 0 {
			return fmt.Errorf("Error: parent node %s has no attributes. %s", parentID, checkMetadataMsg)
		}

		nextParentID := string(pAttrs[prefixes.ParentidAttr])
		if verbose {
			pName := string(pAttrs[prefixes.NameAttr])
			fmt.Printf("parent %s (%s) has parent %s\n", parentID, pName, nextParentID)
		}

		if nextParentID == "" {
			return fmt.Errorf("Error: node parent %s has no parentID. %s", parentID, checkMetadataMsg)
		}
		currentID = nextParentID
	}

	return nil
}

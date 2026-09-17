package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/shamaton/msgpack/v2"
	"github.com/urfave/cli/v2"
)

// metadataExtension is the suffix of the messagepack metadata file of a node
const metadataExtension = ".mpk"

// RecalculateTreesize is the entry point for the recalculate-treesize command
func RecalculateTreesize(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name: "recalculate-treesize",
		Usage: "Recalculate the treesize of all directories in a space from the actual size of their children. " +
			"Use this to repair a space whose reported quota usage drifted from the data on disk",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "space-id",
				Usage:    "Space ID to recalculate (omit to process all spaces)",
				Required: false,
				Aliases:  []string{"s"},
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Only show which treesizes would be corrected without writing them",
				Value: true, // default must be true to avoid accidental writes!
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
			dryRun := c.Bool("dry-run")
			verbose := c.Bool("verbose")

			spaceIDs := []string{}
			if c.IsSet("space-id") {
				spaceIDs = append(spaceIDs, c.String("space-id"))
			} else {
				fmt.Println("Scanning all spaces for incorrect treesizes...")
				spaceIDs = globSpaceIDs(cfg)
			}

			if dryRun {
				fmt.Print("Dry run mode enabled, no treesize will be written\n\n")
			}

			var (
				table     *tablewriter.Table
				corrected int
			)
			table = tablewriter.NewTable(os.Stdout)
			table.Header("Space", "Node", "Stored", "Calculated", "Difference")

			for _, spaceID := range spaceIDs {
				if verbose {
					fmt.Printf("Recalculating treesizes for space: %s\n", spaceID)
				}
				root := filepath.Join(cfg.Drivers.OCIS.Root, "spaces", lookup.Pathify(spaceID, 1, 2), "nodes", lookup.Pathify(spaceID, 4, 2))
				_, n, err := recalculateNode(cfg, spaceID, root, dryRun, verbose, table)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error recalculating space %s: %s\n", spaceID, err)
					continue
				}
				corrected += n
			}

			if corrected > 0 {
				table.Render()
			}

			switch {
			case corrected == 0:
				fmt.Println("All treesizes are correct")
			case dryRun:
				fmt.Printf("Would correct %d treesizes. Run with --dry-run=false to apply\n", corrected)
			default:
				fmt.Printf("Corrected %d treesizes\n", corrected)
			}

			return nil
		},
	}
}

// recalculateNode calculates the size of the node at the given path. For directories it
// recurses into the children first, so the whole subtree is corrected bottom up, and
// compares the sum against the stored treesize. It returns the calculated size of the
// node and the number of corrected treesizes in the subtree.
func recalculateNode(cfg *config.Config, spaceID, path string, dryRun, verbose bool, table *tablewriter.Table) (uint64, int, error) {
	attribs, err := readAttributes(path)
	if err != nil {
		return 0, 0, err
	}

	// Files simply report their blobsize, there is nothing to recalculate. Only an
	// explicit file type is treated as a file, everything else is walked as a
	// directory: a space root for example carries no type attribute at all. This
	// matches how the propagator classifies children.
	if string(attribs[prefixes.TypeAttr]) == strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10) {
		size, err := strconv.ParseUint(string(attribs[prefixes.BlobsizeAttr]), 10, 64)
		if err != nil {
			// a node without a blobsize contributes nothing, but it is worth reporting
			if verbose {
				fmt.Fprintf(os.Stderr, "  Could not read blobsize of %s: %s\n", path, err)
			}
			return 0, 0, nil
		}
		return size, 0, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	var (
		calculated uint64
		corrected  int
	)
	for _, entry := range entries {
		child, err := filepath.EvalSymlinks(filepath.Join(path, entry.Name()))
		if err != nil {
			// dangling child entry, it contributes nothing
			if verbose {
				fmt.Fprintf(os.Stderr, "  Could not resolve child %s: %s\n", entry.Name(), err)
			}
			continue
		}
		size, n, err := recalculateNode(cfg, spaceID, child, dryRun, verbose, table)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "  Could not recalculate child %s: %s\n", child, err)
			}
			continue
		}
		calculated += size
		corrected += n
	}

	stored, storedErr := strconv.ParseUint(string(attribs[prefixes.TreesizeAttr]), 10, 64)
	if storedErr == nil && stored == calculated {
		return calculated, corrected, nil
	}

	storedStr := string(attribs[prefixes.TreesizeAttr])
	if storedStr == "" {
		storedStr = "unset"
	}
	diff := "n/a"
	if storedErr == nil {
		diff = fmt.Sprintf("%+d", int64(calculated)-int64(stored))
	}
	table.Append([]string{spaceID, nodeIDFromPath(path), storedStr, strconv.FormatUint(calculated, 10), diff})
	corrected++

	if !dryRun {
		if err := writeTreesize(path, attribs, calculated); err != nil {
			return calculated, corrected, fmt.Errorf("could not write treesize of %s: %w", path, err)
		}
	}

	return calculated, corrected, nil
}

// readAttributes reads the messagepack metadata of the node at the given path
func readAttributes(path string) (map[string][]byte, error) {
	b, err := os.ReadFile(path + metadataExtension)
	if err != nil {
		return nil, err
	}
	attribs := map[string][]byte{}
	if err := msgpack.Unmarshal(b, &attribs); err != nil {
		return nil, err
	}
	return attribs, nil
}

// writeTreesize sets the treesize attribute of the node at the given path
func writeTreesize(path string, attribs map[string][]byte, size uint64) error {
	attribs[prefixes.TreesizeAttr] = []byte(strconv.FormatUint(size, 10))
	b, err := msgpack.Marshal(attribs)
	if err != nil {
		return err
	}
	return os.WriteFile(path+metadataExtension, b, 0600)
}

// nodeIDFromPath rebuilds the node id from a pathified node path, i.e.
// .../nodes/5b/ba/1e/a7/-f185-4f31-8342-ed4b5743f096 -> 5bba1ea7-f185-4f31-8342-ed4b5743f096
func nodeIDFromPath(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	if len(parts) < 5 {
		return path
	}
	return strings.Join(parts[len(parts)-5:], "")
}

package command

import (
	"errors"
	"fmt"
	"io"
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
	"github.com/owncloud/reva/v2/pkg/storagespace"
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
				Usage:    "Space ID to recalculate, either <storage-id>$<space-id>!<opaque-id>, <storage-id>$<space-id> or the bare <space-id> (omit to process all spaces)",
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
				spaceID, err := parseTreesizeSpaceID(c.String("space-id"))
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
				spaceIDs = append(spaceIDs, spaceID)
			} else {
				fmt.Println("Scanning all spaces for incorrect treesizes...")
				var err error
				spaceIDs, err = listSpaceIDs(cfg)
				if err != nil {
					return cli.Exit(err.Error(), 1)
				}
			}

			return recalculateSpaces(cfg, spaceIDs, dryRun, verbose, os.Stdout, os.Stderr)
		},
	}
}

// parseTreesizeSpaceID extracts the bare space id from the id forms accepted by the
// --space-id flag: <storage-id>$<space-id>!<opaque-id>, <storage-id>$<space-id> and
// <space-id>.
func parseTreesizeSpaceID(id string) (string, error) {
	rid, err := storagespace.ParseID(id)
	if err != nil {
		return "", fmt.Errorf("invalid space id %q: %w", id, err)
	}
	spaceID := rid.GetSpaceId()
	// the id is used to build a path on disk, so only accept a plain, non empty id.
	// No legitimate space id contains "..", but pathifying one produces ".." path
	// elements, e.g. "ab......X" resolves to a sibling of the space.
	if spaceID == "" || spaceID == "." || strings.ContainsAny(spaceID, `/\!$`) || strings.Contains(spaceID, "..") {
		return "", fmt.Errorf("invalid space id %q", id)
	}
	return spaceID, nil
}

// listSpaceIDs returns the ids of all spaces below the storage root. Finding none is an
// error, it points to a wrong storage root or an unsupported storage driver.
func listSpaceIDs(cfg *config.Config) ([]string, error) {
	spaceIDs, err := globSpaceIDs(cfg)
	if err != nil {
		return nil, err
	}
	if len(spaceIDs) == 0 {
		return nil, fmt.Errorf("no spaces found in %s, check STORAGE_USERS_OCIS_ROOT and the storage driver", cfg.Drivers.OCIS.Root)
	}
	return spaceIDs, nil
}

// spaceRootPath returns the path of the root node of the space. It fails if the path
// would end up outside of the nodes directory of the space, as pathifying an id like
// "........x" or "ab......X" produces ".." path elements. parseTreesizeSpaceID already
// rejects such ids, this is defence in depth for ids found on disk.
func spaceRootPath(cfg *config.Config, spaceID string) (string, error) {
	spacesDir := filepath.Join(cfg.Drivers.OCIS.Root, "spaces")
	nodesDir := filepath.Join(spacesDir, lookup.Pathify(spaceID, 1, 2), "nodes")
	root := filepath.Join(nodesDir, lookup.Pathify(spaceID, 4, 2))
	if !isBelow(spacesDir, nodesDir) || !isBelow(nodesDir, root) {
		return "", fmt.Errorf("invalid space id %q", spaceID)
	}
	return root, nil
}

// isBelow reports whether path is strictly below base
func isBelow(base, path string) bool {
	rel, err := filepath.Rel(base, path)
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// recalculateSpaces recalculates the treesizes of the given spaces. A failing space does
// not stop the processing of the remaining ones, but makes the command fail.
func recalculateSpaces(cfg *config.Config, spaceIDs []string, dryRun, verbose bool, out, errOut io.Writer) error {
	if dryRun {
		fmt.Fprint(out, "Dry run mode enabled, no treesize will be written\n\n")
	}

	table := tablewriter.NewTable(out)
	table.Header("Space", "Node", "Stored", "Calculated", "Difference")

	var corrected, failed int
	for _, spaceID := range spaceIDs {
		if verbose {
			fmt.Fprintf(out, "Recalculating treesizes for space: %s\n", spaceID)
		}
		root, err := spaceRootPath(cfg, spaceID)
		if err != nil {
			fmt.Fprintf(errOut, "Error recalculating space %s: %s\n", spaceID, err)
			failed++
			continue
		}
		_, n, err := recalculateNode(cfg, spaceID, root, dryRun, verbose, table, errOut)
		// treesizes corrected before the error occurred count as well
		corrected += n
		if err != nil {
			fmt.Fprintf(errOut, "Error recalculating space %s: %s\n", spaceID, err)
			failed++
		}
	}

	if corrected > 0 {
		table.Render()
	}

	switch {
	case corrected > 0 && dryRun:
		fmt.Fprintf(out, "Would correct %d treesizes. Run with --dry-run=false to apply\n", corrected)
	case corrected > 0:
		fmt.Fprintf(out, "Corrected %d treesizes\n", corrected)
	case failed == 0:
		fmt.Fprintln(out, "All treesizes are correct")
	}

	if failed > 0 {
		return cli.Exit(fmt.Sprintf("%d of %d spaces could not be processed completely, their treesizes may still be incorrect", failed, len(spaceIDs)), 1)
	}
	return nil
}

// recalculateNode calculates the size of the node at the given path. For directories it
// recurses into the children first, so the whole subtree is corrected bottom up, and
// compares the sum against the stored treesize. It returns the calculated size of the
// node and the number of treesizes in the subtree that were corrected, or that would be
// in dry-run mode.
//
// Any error in the subtree fails it: no node above the failing one is written, as its
// sum would be incomplete. The error is printed once to errOut where it occurs, naming
// the exact path, the returned error only tells the caller that the subtree failed.
func recalculateNode(cfg *config.Config, spaceID, path string, dryRun, verbose bool, table *tablewriter.Table, errOut io.Writer) (uint64, int, error) {
	// resolve the start node as well, so a child link back to it is detected as a cycle
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	w := &treesizeWalker{spaceID: spaceID, dryRun: dryRun, verbose: verbose, table: table, errOut: errOut, chain: map[string]bool{}}
	return w.recalculate(path)
}

// treesizeWalker holds the state of the recalculation of one space
type treesizeWalker struct {
	spaceID string
	dryRun  bool
	verbose bool
	table   *tablewriter.Table
	errOut  io.Writer
	// chain holds the resolved paths of the nodes currently being recalculated, i.e. the
	// node and its ancestors, to detect child links that form a cycle
	chain map[string]bool
}

// nodeError is the error of the node at path. Its message only names the node, the
// cause has already been printed where it occurred.
type nodeError struct {
	path string
	err  error
}

func (e *nodeError) Error() string { return "could not recalculate " + e.path }
func (e *nodeError) Unwrap() error { return e.err }

// fail prints the error of the node at path and returns it
func (w *treesizeWalker) fail(path string, err error) error {
	fmt.Fprintf(w.errOut, "  Could not recalculate %s: %s\n", path, err)
	return &nodeError{path: path, err: err}
}

func (w *treesizeWalker) recalculate(path string) (uint64, int, error) {
	attribs, err := readAttributes(path)
	if err != nil {
		return 0, 0, w.fail(path, err)
	}

	// Files simply report their blobsize, there is nothing to recalculate. Only an
	// explicit file type is treated as a file, everything else is walked as a
	// directory: a space root for example carries no type attribute at all. This
	// matches how the propagator classifies children.
	if string(attribs[prefixes.TypeAttr]) == strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10) {
		size, err := strconv.ParseUint(string(attribs[prefixes.BlobsizeAttr]), 10, 64)
		if err != nil {
			// decomposedfs always sets the blobsize of a file, so this is corruption. Counting
			// it as 0 would write an undercounted treesize into the parents.
			return 0, 0, w.fail(path, fmt.Errorf("invalid blobsize: %w", err))
		}
		return size, 0, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return 0, 0, nil
		}
		return 0, 0, w.fail(path, err)
	}

	w.chain[path] = true
	defer delete(w.chain, path)

	var (
		calculated     uint64
		corrected      int
		failedChildren int
	)
	for _, entry := range entries {
		link := filepath.Join(path, entry.Name())
		child, err := filepath.EvalSymlinks(link)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			// dangling child entry, it contributes nothing
			if w.verbose {
				fmt.Fprintf(w.errOut, "  Skipping dangling child %s: %s\n", link, err)
			}
			continue
		case err != nil:
			// the remaining children are still processed, but the sum is incomplete
			_ = w.fail(link, err)
			failedChildren++
			continue
		case w.chain[child]:
			_ = w.fail(link, fmt.Errorf("child link forms a cycle, it resolves to its ancestor %s", child))
			failedChildren++
			continue
		}
		size, n, err := w.recalculate(child)
		// treesizes written below the failing node count as well
		corrected += n
		if err != nil {
			// already printed where it happened
			failedChildren++
			continue
		}
		calculated += size
	}

	if failedChildren > 0 {
		return 0, corrected, fmt.Errorf("%d children of %s could not be recalculated", failedChildren, path)
	}

	stored, storedErr := strconv.ParseUint(string(attribs[prefixes.TreesizeAttr]), 10, 64)
	if storedErr == nil && stored == calculated {
		return calculated, corrected, nil
	}

	if !w.dryRun {
		if err := writeTreesize(path, attribs, calculated); err != nil {
			// not counted and no table row, the treesize was not corrected
			return 0, corrected, w.fail(path, fmt.Errorf("could not write treesize: %w", err))
		}
	}

	storedStr := string(attribs[prefixes.TreesizeAttr])
	if storedStr == "" {
		storedStr = "unset"
	}
	diff := "n/a"
	if storedErr == nil {
		diff = fmt.Sprintf("%+d", int64(calculated)-int64(stored))
	}
	w.table.Append([]string{w.spaceID, nodeIDFromPath(path), storedStr, strconv.FormatUint(calculated, 10), diff})

	return calculated, corrected + 1, nil
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

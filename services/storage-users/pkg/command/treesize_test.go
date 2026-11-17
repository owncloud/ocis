package command

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/shamaton/msgpack/v2"
	"github.com/test-go/testify/require"
	"github.com/urfave/cli/v2"
)

// writeNode creates a node file with the given attributes on disk
func writeNode(t *testing.T, path string, attribs map[string][]byte) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
	b, err := msgpack.Marshal(attribs)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path+metadataExtension, b, 0600))
}

// writeFileNode creates a file node with the given blobsize. The node file itself
// has to exist as well, decomposedfs resolves children through it.
func writeFileNode(t *testing.T, path string, blobsize int) {
	t.Helper()
	writeNode(t, path, map[string][]byte{
		prefixes.TypeAttr:     []byte(strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10)),
		prefixes.BlobsizeAttr: []byte(strconv.Itoa(blobsize)),
	})
	require.NoError(t, os.WriteFile(path, []byte{}, 0600))
}

// writeDirNode creates a directory node with the given stored treesize
func writeDirNode(t *testing.T, path string, treesize int) {
	t.Helper()
	require.NoError(t, os.MkdirAll(path, 0700))
	writeNode(t, path, map[string][]byte{
		prefixes.TypeAttr:     []byte(strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_CONTAINER), 10)),
		prefixes.TreesizeAttr: []byte(strconv.Itoa(treesize)),
	})
}

// storedTreesize reads back the treesize attribute of a node
func storedTreesize(t *testing.T, path string) string {
	t.Helper()
	attribs, err := readAttributes(path)
	require.NoError(t, err)
	return string(attribs[prefixes.TreesizeAttr])
}

// linkChild links a child node into a parent directory, mirroring how decomposedfs
// references children by symlink
func linkChild(t *testing.T, parent, child, name string) {
	t.Helper()
	rel, err := filepath.Rel(parent, child)
	require.NoError(t, err)
	require.NoError(t, os.Symlink(rel, filepath.Join(parent, name)))
}

func TestRecalculateNode(t *testing.T) {
	t.Run("corrects a wrong treesize", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 999)
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 100)
		linkChild(t, dir, f, "file")

		size, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(100), size)
		require.Equal(t, 1, corrected)
		require.Equal(t, "100", storedTreesize(t, dir))
	})

	t.Run("leaves a correct treesize untouched", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 100)
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 100)
		linkChild(t, dir, f, "file")

		size, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(100), size)
		require.Equal(t, 0, corrected)
	})

	t.Run("does not write in dry-run mode", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 999)
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 100)
		linkChild(t, dir, f, "file")

		size, corrected, err := recalculateNode(nil, "spaceid", dir, true, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(100), size)
		require.Equal(t, 1, corrected)
		require.Equal(t, "999", storedTreesize(t, dir), "dry-run must not write the treesize")
	})

	t.Run("sums nested directories bottom up", func(t *testing.T) {
		root := t.TempDir()
		outer := filepath.Join(root, "outer")
		writeDirNode(t, outer, 0)
		inner := filepath.Join(root, "inner")
		writeDirNode(t, inner, 0)
		f1 := filepath.Join(root, "f1")
		writeFileNode(t, f1, 10)
		f2 := filepath.Join(root, "f2")
		writeFileNode(t, f2, 32)

		linkChild(t, inner, f2, "f2")
		linkChild(t, outer, f1, "f1")
		linkChild(t, outer, inner, "inner")

		size, corrected, err := recalculateNode(nil, "spaceid", outer, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(42), size)
		require.Equal(t, 2, corrected, "both the inner and the outer directory are corrected")
		require.Equal(t, "32", storedTreesize(t, inner))
		require.Equal(t, "42", storedTreesize(t, outer))
	})

	t.Run("treats a node without a type attribute as a directory", func(t *testing.T) {
		// space roots carry no type attribute
		root := t.TempDir()
		spaceRoot := filepath.Join(root, "space")
		require.NoError(t, os.MkdirAll(spaceRoot, 0700))
		writeNode(t, spaceRoot, map[string][]byte{
			prefixes.TreesizeAttr: []byte("999"),
		})
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 64)
		linkChild(t, spaceRoot, f, "file")

		size, corrected, err := recalculateNode(nil, "spaceid", spaceRoot, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(64), size)
		require.Equal(t, 1, corrected)
		require.Equal(t, "64", storedTreesize(t, spaceRoot))
	})

	t.Run("skips dangling child links", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 0)
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 20)
		linkChild(t, dir, f, "file")
		require.NoError(t, os.Symlink("../gone", filepath.Join(dir, "gone")))

		size, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, uint64(20), size, "the dangling entry contributes nothing")
	})

	t.Run("sets an unset treesize", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		require.NoError(t, os.MkdirAll(dir, 0700))
		writeNode(t, dir, map[string][]byte{
			prefixes.TypeAttr: []byte(strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_CONTAINER), 10)),
		})
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 7)
		linkChild(t, dir, f, "file")

		_, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout), io.Discard)
		require.NoError(t, err)
		require.Equal(t, 1, corrected)
		require.Equal(t, "7", storedTreesize(t, dir))
	})
}

func TestParseTreesizeSpaceID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"bare space id", "9d408cec-8f0a-4d33-8715-89df1217a10c", "9d408cec-8f0a-4d33-8715-89df1217a10c", false},
		{"storage id and space id", "storage-users-1$9d408cec-8f0a-4d33-8715-89df1217a10c", "9d408cec-8f0a-4d33-8715-89df1217a10c", false},
		{"project form with opaque id", "storage-users-1$9d408cec-8f0a-4d33-8715-89df1217a10c!9d408cec-8f0a-4d33-8715-89df1217a10c", "9d408cec-8f0a-4d33-8715-89df1217a10c", false},
		{"empty", "", "", true},
		{"storage id only", "storage-users-1$", "", true},
		{"opaque id only", "!abc", "", true},
		{"path traversal", "../../etc", "", true},
		{"path separator", "a/b", "", true},
		{"dot dot escaping the nodes dir", "ab......X", "", true},
		{"dot dot inside", "a..b", "", true},
		{"dot", ".", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTreesizeSpaceID(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

// writeSpace creates a space root with one file below the given decomposedfs root
func writeSpace(t *testing.T, storageRoot, spaceID string, treesize, blobsize int) string {
	t.Helper()
	root := filepath.Join(storageRoot, "spaces", lookup.Pathify(spaceID, 1, 2), "nodes", lookup.Pathify(spaceID, 4, 2))
	writeDirNode(t, root, treesize)
	f := filepath.Join(storageRoot, "file-"+spaceID)
	writeFileNode(t, f, blobsize)
	linkChild(t, root, f, "file")
	return root
}

func TestRecalculateSpaces(t *testing.T) {
	const good = "aaaaaaaa-0000-0000-0000-000000000001"
	const good2 = "cccccccc-0000-0000-0000-000000000003"
	const missing = "bbbbbbbb-0000-0000-0000-000000000002"

	t.Run("returns an error and no success message when a space is missing", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		var out, errOut bytes.Buffer

		err := recalculateSpaces(cfg, []string{missing}, true, false, &out, &errOut)
		require.Error(t, err)
		var exitErr cli.ExitCoder
		require.True(t, errors.As(err, &exitErr))
		require.NotEqual(t, 0, exitErr.ExitCode())
		require.NotContains(t, out.String(), "All treesizes are correct")
		require.Contains(t, errOut.String(), missing)
	})

	t.Run("processes remaining spaces and reports the failure count", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		root := writeSpace(t, cfg.Drivers.OCIS.Root, good, 999, 5)
		root2 := writeSpace(t, cfg.Drivers.OCIS.Root, good2, 999, 6)
		var out, errOut bytes.Buffer

		err := recalculateSpaces(cfg, []string{good, missing, good2}, false, false, &out, &errOut)
		require.Error(t, err)
		require.Contains(t, err.Error(), "1 of 3 spaces could not be processed")
		require.Equal(t, "5", storedTreesize(t, root))
		require.Equal(t, "6", storedTreesize(t, root2))
		require.NotContains(t, out.String(), "All treesizes are correct")
	})

	t.Run("succeeds when every space is fine", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		writeSpace(t, cfg.Drivers.OCIS.Root, good, 5, 5)
		var out, errOut bytes.Buffer

		require.NoError(t, recalculateSpaces(cfg, []string{good}, true, false, &out, &errOut))
		require.Contains(t, out.String(), "All treesizes are correct")
	})
}

func TestRecalculateNodeChildErrors(t *testing.T) {
	setup := func(t *testing.T) (dir, bad string) {
		root := t.TempDir()
		dir = filepath.Join(root, "dir")
		writeDirNode(t, dir, 999)
		ok := filepath.Join(root, "ok")
		writeFileNode(t, ok, 10)
		linkChild(t, dir, ok, "ok")
		bad = filepath.Join(root, "bad")
		writeFileNode(t, bad, 10)
		require.NoError(t, os.WriteFile(bad+metadataExtension, []byte("not msgpack \xff\xff"), 0600))
		linkChild(t, dir, bad, "bad")
		return dir, bad
	}

	t.Run("a failing child fails the subtree and the parent is not written", func(t *testing.T) {
		dir, bad := setup(t)
		var errOut bytes.Buffer

		_, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(io.Discard), &errOut)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, dir), "an incomplete sum must not be written")
		require.Contains(t, errOut.String(), bad)
	})

	t.Run("a failing child deep in the tree keeps every ancestor untouched", func(t *testing.T) {
		dir, _ := setup(t)
		top := filepath.Join(filepath.Dir(dir), "top")
		writeDirNode(t, top, 888)
		linkChild(t, top, dir, "dir")

		_, _, err := recalculateNode(nil, "spaceid", top, false, false, tablewriter.NewTable(io.Discard), io.Discard)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, dir))
		require.Equal(t, "888", storedTreesize(t, top))
	})
}

func TestRecalculateSpacesHardening(t *testing.T) {
	const good = "aaaaaaaa-0000-0000-0000-000000000001"

	t.Run("spaceRootPath rejects ids that escape the space nodes dir", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		for _, id := range []string{"........x", "..../..x", "......x", "..", "....", "../..", "ab......X", "."} {
			_, err := spaceRootPath(cfg, id)
			require.Error(t, err, id)
			require.Contains(t, err.Error(), "invalid space id", id)
		}
	})

	t.Run("spaceRootPath accepts a normal id", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		root, err := spaceRootPath(cfg, good)
		require.NoError(t, err)
		require.Equal(t, filepath.Join(cfg.Drivers.OCIS.Root, "spaces", "aa", "aaaaaa-0000-0000-0000-000000000001", "nodes", "aa", "aa", "aa", "aa", "-0000-0000-0000-000000000001"), root)
	})

	t.Run("an unreadable subtree fails the space without writing the root", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		root := writeSpace(t, cfg.Drivers.OCIS.Root, good, 999, 5)
		bad := filepath.Join(cfg.Drivers.OCIS.Root, "bad")
		writeFileNode(t, bad, 10)
		require.NoError(t, os.WriteFile(bad+metadataExtension, []byte("garbage \xff"), 0600))
		linkChild(t, root, bad, "bad")
		var out, errOut bytes.Buffer

		err := recalculateSpaces(cfg, []string{good}, false, false, &out, &errOut)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, root))
		require.NotContains(t, out.String(), "All treesizes are correct")
	})

	t.Run("prints the totals of the successful spaces when another fails", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		writeSpace(t, cfg.Drivers.OCIS.Root, good, 999, 5)
		var out, errOut bytes.Buffer

		err := recalculateSpaces(cfg, []string{good, "bbbbbbbb-0000-0000-0000-000000000002"}, false, false, &out, &errOut)
		require.Error(t, err)
		require.Contains(t, out.String(), "Corrected 1 treesizes")
		require.NotContains(t, err.Error(), "not verified")
	})
}

func TestListSpaceIDs(t *testing.T) {
	t.Run("fails when no space is found", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		_, err := listSpaceIDs(cfg)
		require.Error(t, err)
		require.Contains(t, err.Error(), cfg.Drivers.OCIS.Root)
	})

	t.Run("finds spaces", func(t *testing.T) {
		cfg := &config.Config{}
		cfg.Drivers.OCIS.Root = t.TempDir()
		writeSpace(t, cfg.Drivers.OCIS.Root, "aaaaaaaa-0000-0000-0000-000000000001", 0, 0)
		ids, err := listSpaceIDs(cfg)
		require.NoError(t, err)
		require.Equal(t, []string{"aaaaaaaa-0000-0000-0000-000000000001"}, ids)
	})
}

func TestRecalculateNodeInvalidBlobsize(t *testing.T) {
	for name, attribs := range map[string]map[string][]byte{
		"missing": {prefixes.TypeAttr: []byte(strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10))},
		"garbage": {
			prefixes.TypeAttr:     []byte(strconv.FormatUint(uint64(provider.ResourceType_RESOURCE_TYPE_FILE), 10)),
			prefixes.BlobsizeAttr: []byte("abc"),
		},
	} {
		t.Run(name+" blobsize fails the subtree and nothing is written", func(t *testing.T) {
			root := t.TempDir()
			dir := filepath.Join(root, "dir")
			writeDirNode(t, dir, 999)
			f := filepath.Join(root, "file")
			writeNode(t, f, attribs)
			require.NoError(t, os.WriteFile(f, []byte{}, 0600))
			linkChild(t, dir, f, "file")
			var errOut bytes.Buffer

			_, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(io.Discard), &errOut)
			require.Error(t, err)
			require.Equal(t, "999", storedTreesize(t, dir))
			require.Contains(t, errOut.String(), f, "reported even without --verbose")
		})
	}
}

// skipIfRoot skips tests that rely on permission errors, root bypasses them
func skipIfRoot(t *testing.T) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
}

func TestRecalculateNodeRound4(t *testing.T) {
	t.Run("an unresolvable child link other than a dangling one fails the subtree", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 999)
		f := filepath.Join(root, "file")
		writeFileNode(t, f, 20)
		linkChild(t, dir, f, "file")
		// a self referencing link fails with too many links, not with ENOENT
		loop := filepath.Join(dir, "loop")
		require.NoError(t, os.Symlink("loop", loop))
		var errOut bytes.Buffer

		_, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(io.Discard), &errOut)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, dir), "an incomplete sum must not be written")
		require.Contains(t, errOut.String(), "loop", "reported even without --verbose")
	})

	t.Run("a failed write is not counted and fails the subtree", func(t *testing.T) {
		skipIfRoot(t)
		root := t.TempDir()
		outer := filepath.Join(root, "outer")
		writeDirNode(t, outer, 0)
		inner := filepath.Join(root, "inner")
		writeDirNode(t, inner, 0)
		f := filepath.Join(root, "f")
		writeFileNode(t, f, 32)
		linkChild(t, inner, f, "f")
		linkChild(t, outer, inner, "inner")
		require.NoError(t, os.Chmod(outer+metadataExtension, 0400))
		var tableOut, errOut bytes.Buffer
		table := tablewriter.NewTable(&tableOut)

		_, corrected, err := recalculateNode(nil, "spaceid", outer, false, false, table, &errOut)
		require.Error(t, err)
		require.Equal(t, 1, corrected, "only the written inner treesize counts")
		require.Equal(t, "32", storedTreesize(t, inner))
		require.Equal(t, "0", storedTreesize(t, outer))
		require.NoError(t, table.Render())
		require.Contains(t, tableOut.String(), "inner")
		require.NotContains(t, tableOut.String(), "outer", "a node whose write failed has no table row")
		require.Contains(t, errOut.String(), outer)
	})

	t.Run("a child link back to an ancestor is a cycle and fails the subtree", func(t *testing.T) {
		root := t.TempDir()
		outer := filepath.Join(root, "outer")
		writeDirNode(t, outer, 999)
		inner := filepath.Join(root, "inner")
		writeDirNode(t, inner, 999)
		linkChild(t, outer, inner, "inner")
		linkChild(t, inner, outer, "outer")
		var errOut bytes.Buffer

		_, _, err := recalculateNode(nil, "spaceid", outer, false, false, tablewriter.NewTable(io.Discard), &errOut)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, outer))
		require.Equal(t, "999", storedTreesize(t, inner))
		require.Contains(t, errOut.String(), "cycle")
	})

	t.Run("a child link to the node itself is a cycle", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "dir")
		writeDirNode(t, dir, 999)
		linkChild(t, dir, dir, "self")
		var errOut bytes.Buffer

		_, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(io.Discard), &errOut)
		require.Error(t, err)
		require.Equal(t, "999", storedTreesize(t, dir))
		require.Contains(t, errOut.String(), "cycle")
	})

	t.Run("a deep failure is reported once at the leaf", func(t *testing.T) {
		root := t.TempDir()
		var parent string
		var levels []string
		for i := 0; i < 4; i++ {
			d := filepath.Join(root, "level"+strconv.Itoa(i))
			writeDirNode(t, d, 999)
			if parent != "" {
				linkChild(t, parent, d, "child")
			}
			levels = append(levels, d)
			parent = d
		}
		bad := filepath.Join(root, "bad")
		writeFileNode(t, bad, 10)
		require.NoError(t, os.WriteFile(bad+metadataExtension, []byte("garbage \xff"), 0600))
		linkChild(t, parent, bad, "bad")
		var errOut bytes.Buffer

		_, _, err := recalculateNode(nil, "spaceid", levels[0], false, false, tablewriter.NewTable(io.Discard), &errOut)
		require.Error(t, err)
		lines := strings.Split(strings.TrimSpace(errOut.String()), "\n")
		require.Len(t, lines, 1, errOut.String())
		require.Contains(t, lines[0], "bad")
		for _, l := range levels {
			require.Equal(t, "999", storedTreesize(t, l))
		}
	})
}

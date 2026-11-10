package command

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
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

		size, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout))
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

		size, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout))
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

		size, corrected, err := recalculateNode(nil, "spaceid", dir, true, false, tablewriter.NewTable(os.Stdout))
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

		size, corrected, err := recalculateNode(nil, "spaceid", outer, false, false, tablewriter.NewTable(os.Stdout))
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

		size, corrected, err := recalculateNode(nil, "spaceid", spaceRoot, false, false, tablewriter.NewTable(os.Stdout))
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

		size, _, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout))
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

		_, corrected, err := recalculateNode(nil, "spaceid", dir, false, false, tablewriter.NewTable(os.Stdout))
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
		require.Contains(t, err.Error(), "1 of 3 spaces failed")
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

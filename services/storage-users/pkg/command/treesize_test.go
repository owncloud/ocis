package command

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/olekukonko/tablewriter"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/shamaton/msgpack/v2"
	"github.com/test-go/testify/require"
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

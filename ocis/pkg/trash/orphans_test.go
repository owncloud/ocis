package trash_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis/pkg/trash"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/require"
)

const (
	_spaceID  = "ab" + "cdef"
	_nodeID   = "11223344-aaaa"
	_revision = _nodeID + ".REV.2026-11-10T08:43:15.41023Z"
	_trashed  = _nodeID + ".T.2026-11-10T08:49:21.088168Z"
	_minAge   = time.Hour
)

// old is older than _minAge, young is younger than _minAge
var (
	old   = time.Now().Add(-2 * _minAge)
	young = time.Now()
)

type fakeBlobstore struct {
	deleted []string
}

func (bs *fakeBlobstore) Delete(n *node.Node) error {
	bs.deleted = append(bs.deleted, n.SpaceID+"/"+n.BlobID)
	return nil
}

// nodesDir returns the directory that holds the files of _nodeID
func nodesDir(base string) string {
	return filepath.Join(base, "spaces", "ab", "cdef", "nodes", "11", "22", "33", "44")
}

func touch(t *testing.T, path string, mtime time.Time) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
	require.NoError(t, os.WriteFile(path, nil, 0600))
	require.NoError(t, os.Chtimes(path, mtime, mtime))
}

func writeMetadata(t *testing.T, path, blobID string, mtime time.Time) {
	t.Helper()
	b, err := msgpack.Marshal(map[string][]byte{"user.ocis.blobid": []byte(blobID)})
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0700))
	require.NoError(t, os.WriteFile(path, b, 0600))
	require.NoError(t, os.Chtimes(path, mtime, mtime))
}

// existing returns the names of the files that are left in the nodes directory
func existing(t *testing.T, base string) []string {
	t.Helper()
	entries, err := os.ReadDir(nodesDir(base))
	require.NoError(t, err)
	names := []string{}
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

func purge(t *testing.T, base string, bs trash.BlobDeleter, dryRun bool) string {
	t.Helper()
	return purgeWithOptions(t, base, trash.PurgeOrphansOptions{Blobstore: bs, DryRun: dryRun})
}

func purgeWithOptions(t *testing.T, base string, opts trash.PurgeOrphansOptions) string {
	t.Helper()
	out := &bytes.Buffer{}
	opts.MinAge = _minAge
	opts.Out = out
	require.NoError(t, trash.PurgeOrphanedFiles(base, opts))
	return out.String()
}

func TestPurgeOrphanedFilesRemovesLockFileOfMissingNode(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)
	touch(t, filepath.Join(dir, "other-node"), old)

	purge(t, base, &fakeBlobstore{}, false)

	require.ElementsMatch(t, []string{"other-node"}, existing(t, base))
}

func TestPurgeOrphanedFilesKeepsFilesOfExistingNode(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID), old)
	writeMetadata(t, filepath.Join(dir, _nodeID+".mpk"), "blob-1", old)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)
	touch(t, filepath.Join(dir, _revision+".mlock"), old)
	bs := &fakeBlobstore{}

	purge(t, base, bs, false)

	require.ElementsMatch(t, []string{
		_nodeID, _nodeID + ".mpk", _nodeID + ".mlock",
		_revision, _revision + ".mpk", _revision + ".mlock",
	}, existing(t, base))
	require.Empty(t, bs.deleted)
}

func TestPurgeOrphanedFilesKeepsLockFileOfExistingDirectory(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	require.NoError(t, os.MkdirAll(filepath.Join(dir, _nodeID), 0700))
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)

	purge(t, base, &fakeBlobstore{}, false)

	require.ElementsMatch(t, []string{_nodeID, _nodeID + ".mlock"}, existing(t, base))
}

func TestPurgeOrphanedFilesKeepsYoungLockFile(t *testing.T) {
	base := t.TempDir()
	// a new node gets its lock file before its node file
	touch(t, filepath.Join(nodesDir(base), _nodeID+".mlock"), young)

	purge(t, base, &fakeBlobstore{}, false)

	require.ElementsMatch(t, []string{_nodeID + ".mlock"}, existing(t, base))
}

func TestPurgeOrphanedFilesRemovesLockFilesOfTrashedNode(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _trashed), old)
	writeMetadata(t, filepath.Join(dir, _trashed+".mpk"), "blob-1", old)
	// the old version of Delete left the lock file at the old node path
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)
	// the lock file of a trashed node that was purged
	touch(t, filepath.Join(dir, _nodeID+".T.2026-01-01T00:00:00Z.mlock"), old)

	purge(t, base, &fakeBlobstore{}, false)

	require.ElementsMatch(t, []string{_trashed, _trashed + ".mpk"}, existing(t, base))
}

func TestPurgeOrphanedFilesRemovesRevisionOfMissingNode(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)
	touch(t, filepath.Join(dir, _revision+".mlock"), old)
	bs := &fakeBlobstore{}

	purge(t, base, bs, false)

	require.Empty(t, existing(t, base))
	require.Equal(t, []string{_spaceID + "/blob-2"}, bs.deleted)
}

func TestPurgeOrphanedFilesKeepsRevisionOfTrashedNode(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _trashed), old)
	writeMetadata(t, filepath.Join(dir, _trashed+".mpk"), "blob-1", old)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)
	bs := &fakeBlobstore{}

	purge(t, base, bs, false)

	require.ElementsMatch(t, []string{
		_trashed, _trashed + ".mpk", _revision, _revision + ".mpk",
	}, existing(t, base))
	require.Empty(t, bs.deleted)
}

func TestPurgeOrphanedFilesKeepsRevisionWithoutBlobstore(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)

	// without a blobstore, the blob ID must stay on the disk
	purge(t, base, nil, false)

	require.ElementsMatch(t, []string{_revision, _revision + ".mpk"}, existing(t, base))
}

func TestPurgeOrphanedFilesDryRunRemovesNothing(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)
	bs := &fakeBlobstore{}

	out := purge(t, base, bs, true)

	require.ElementsMatch(t, []string{_nodeID + ".mlock", _revision, _revision + ".mpk"}, existing(t, base))
	require.Empty(t, bs.deleted)
	require.Contains(t, out, "would remove 1 lock files and 1 revisions")
}

func TestPurgeOrphanedFilesPrintsOnlyTheSummaryByDefault(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)

	out := purge(t, base, &fakeBlobstore{}, false)

	require.NotContains(t, out, dir)
	require.Contains(t, out, "removed 1 lock files and 0 revisions")
}

func TestPurgeOrphanedFilesVerbosePrintsEachFile(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), old)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)

	out := purgeWithOptions(t, base, trash.PurgeOrphansOptions{Blobstore: &fakeBlobstore{}, DryRun: true, Verbose: true})

	require.Contains(t, out, "would remove "+filepath.Join(dir, _nodeID+".mlock"))
	require.Contains(t, out, "would delete blob blob-2 of space "+_spaceID)
	require.Contains(t, out, "would remove "+filepath.Join(dir, _revision))
	require.Contains(t, out, "would remove "+filepath.Join(dir, _revision+".mpk"))
}

func TestPurgeOrphanedFilesVerbosePrintsYoungFiles(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _nodeID+".mlock"), young)

	out := purgeWithOptions(t, base, trash.PurgeOrphansOptions{Blobstore: &fakeBlobstore{}, Verbose: true})

	require.Contains(t, out, "skip young file "+filepath.Join(dir, _nodeID+".mlock"))
}

func TestPurgeOrphanedFilesVerbosePrintsRevisionsWithoutBlobstore(t *testing.T) {
	base := t.TempDir()
	dir := nodesDir(base)
	touch(t, filepath.Join(dir, _revision), old)
	writeMetadata(t, filepath.Join(dir, _revision+".mpk"), "blob-2", old)

	quiet := purgeWithOptions(t, base, trash.PurgeOrphansOptions{})
	verbose := purgeWithOptions(t, base, trash.PurgeOrphansOptions{Verbose: true})

	require.NotContains(t, quiet, dir)
	require.Contains(t, verbose, "skip revision without blobstore "+filepath.Join(dir, _revision))
}

func TestPurgeOrphanedFilesFailsWithoutNodes(t *testing.T) {
	err := trash.PurgeOrphanedFiles(t.TempDir(), trash.PurgeOrphansOptions{Out: &bytes.Buffer{}})
	require.Error(t, err)
}

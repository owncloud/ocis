package revisions

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/shamaton/msgpack/v2"
	"github.com/test-go/testify/require"
)

// recordingBlobstore records the nodes passed to Delete so the test can assert
// the blob path inputs (SpaceID + BlobID).
type recordingBlobstore struct {
	deleted []*node.Node
}

func (r *recordingBlobstore) Delete(n *node.Node) error {
	r.deleted = append(r.deleted, n)
	return nil
}

// TestPurgeRevisionsDeletesBlobWithSpaceID guards against orphaning blobs: the
// blobstore derives the blob path from the node's SpaceID and BlobID, so
// PurgeRevisions must pass the SpaceID parsed from the revision path. Without
// it the blobstore targets the wrong path (no-op delete) while the revision
// metadata is still removed, leaving the blob orphaned.
func TestPurgeRevisionsDeletesBlobWithSpaceID(t *testing.T) {
	const (
		spaceID = "spaceid1"
		nodeID  = "nodeid1"
		blobID  = "blob-abc"
	)

	tmp := t.TempDir()
	dir := filepath.Join(tmp, "storage", "users", "spaces", spaceID, "nodes")
	require.NoError(t, os.MkdirAll(dir, 0o755))

	revPath := filepath.Join(dir, nodeID+".REV.2024-05-22T07:32:53.123456789Z.mpk")
	value, err := msgpack.Marshal(map[string][]byte{"user.ocis.blobid": []byte(blobID)})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(revPath, value, 0o644))

	nodes := make(chan string, 1)
	nodes <- revPath
	close(nodes)

	bs := &recordingBlobstore{}
	PurgeRevisions(nodes, bs, false, false)

	require.Len(t, bs.deleted, 1, "expected exactly one blob delete")
	require.Equal(t, blobID, bs.deleted[0].BlobID)
	require.Equal(t, spaceID, bs.deleted[0].SpaceID, "SpaceID must be passed to the blobstore, otherwise the blob is orphaned")
}

// TestPurgeRevisionsMixedWholeSecondAndFractional purges a decomposed tree that
// mixes whole-second and fractional revision triplets. Each triplet is the
// revision node plus its .mpk and .mlock companion: three files, one blob, and
// one revision. Dry-run keeps every file and does not call the blobstore.
func TestPurgeRevisionsMixedWholeSecondAndFractional(t *testing.T) {
	for _, d := range revisionDiscoverers() {
		t.Run(d.name, func(t *testing.T) {
			for _, dryRun := range []bool{true, false} {
				name := "delete"
				if dryRun {
					name = "dry-run"
				}
				t.Run(name, func(t *testing.T) {
					tree := newDecomposedRevisionTree(t)
					bs := &recordingBlobstore{}

					files, blobs, revs := PurgeRevisions(d.find(tree.root), bs, dryRun, false)

					require.Equal(t, tree.revisions*3, files, "files")
					require.Equal(t, tree.revisions, blobs, "blobs")
					require.Equal(t, tree.revisions, revs, "revisions")
					requirePresent(t, tree.preservedPaths)

					if dryRun {
						require.Empty(t, bs.deleted)
						requirePresent(t, tree.revisionPaths)
						return
					}

					require.Equal(t, tree.blobs, deletedBlobs(bs))
					requireAbsent(t, tree.revisionPaths)
				})
			}
		})
	}
}

// TestPurgeRevisionsRetainsInvalidWholeSecondMessagePack keeps a whole-second
// .mpk whose metadata is not MessagePack. A sibling whole-second triplet is
// purged in the same call, which shows the retained file was recognized as a
// revision and then left in place by the existing read error path. That path
// also skips the blob delete.
func TestPurgeRevisionsRetainsInvalidWholeSecondMessagePack(t *testing.T) {
	root := t.TempDir()
	nodePath := decomposedNodePath(root, revisionFixtureSpaceID, revisionFixtureNodeID)
	whole := time.Date(2024, 5, 22, 7, 32, 53, 0, time.UTC).UTC().Format(time.RFC3339Nano)
	ancient := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC).UTC().Format(time.RFC3339Nano)
	require.Equal(t, "2024-05-22T07:32:53Z", whole)
	require.Equal(t, "1601-01-01T00:00:00Z", ancient)

	validBase := nodePath + ".REV." + whole
	badMPK := nodePath + ".REV." + ancient + ".mpk"
	currentMPK := nodePath + ".mpk"

	writeFixtureFile(t, nodePath, nil)
	writeFixtureFile(t, currentMPK, revisionBlobMPK(t, "blob-current"))
	writeFixtureFile(t, validBase, nil)
	writeFixtureFile(t, validBase+".mpk", revisionBlobMPK(t, "blob-valid"))
	writeFixtureFile(t, validBase+".mlock", nil)
	writeFixtureFile(t, badMPK, []byte("not-msgpack"))

	nodes := make(chan string, 6)
	for _, p := range []string{
		badMPK,
		validBase,
		validBase + ".mpk",
		validBase + ".mlock",
		nodePath,
		currentMPK,
	} {
		nodes <- p
	}
	close(nodes)

	bs := &recordingBlobstore{}
	files, blobs, revs := PurgeRevisions(nodes, bs, false, false)

	require.Equal(t, 3, files, "files")
	require.Equal(t, 1, blobs, "blobs")
	require.Equal(t, 1, revs, "revisions")
	require.Equal(t, []purgedBlob{{
		spaceID: revisionFixtureSpaceID,
		blobID:  "blob-valid",
	}}, deletedBlobs(bs))
	requireAbsent(t, []string{validBase, validBase + ".mpk", validBase + ".mlock"})
	requirePresent(t, []string{badMPK, nodePath, currentMPK})
}

func deletedBlobs(bs *recordingBlobstore) []purgedBlob {
	got := make([]purgedBlob, 0, len(bs.deleted))
	for _, n := range bs.deleted {
		got = append(got, purgedBlob{spaceID: n.SpaceID, blobID: n.BlobID})
	}
	sort.Slice(got, func(i, j int) bool {
		if got[i].blobID != got[j].blobID {
			return got[i].blobID < got[j].blobID
		}
		return got[i].spaceID < got[j].spaceID
	})
	return got
}

func requirePresent(t *testing.T, paths []string) {
	t.Helper()
	for _, p := range paths {
		_, err := os.Stat(p)
		require.NoError(t, err, p)
	}
}

func requireAbsent(t *testing.T, paths []string) {
	t.Helper()
	for _, p := range paths {
		_, err := os.Stat(p)
		require.True(t, os.IsNotExist(err), "%s still present: %v", p, err)
	}
}

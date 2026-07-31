package revisions

import (
	"os"
	"path/filepath"
	"testing"

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

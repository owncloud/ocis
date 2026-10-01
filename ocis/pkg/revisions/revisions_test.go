package revisions

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/shamaton/msgpack/v2"
	"github.com/test-go/testify/require"
)

var (
	_basePath = "/spaces/8f/638374-6ea8-4f0d-80c4-66d9b49830a5/nodes/"
)

// func TestInit(t *testing.T) {
// initialize(10, 2)
// defer os.RemoveAll("test_temp")
// }

func TestGlob30(t *testing.T)  { test(t, 10, 2, glob) }
func TestGlob80(t *testing.T)  { test(t, 20, 3, glob) }
func TestGlob250(t *testing.T) { test(t, 50, 4, glob) }
func TestGlob600(t *testing.T) { test(t, 100, 5, glob) }

func TestWalk30(t *testing.T)  { test(t, 10, 2, walk) }
func TestWalk80(t *testing.T)  { test(t, 20, 3, walk) }
func TestWalk250(t *testing.T) { test(t, 50, 4, walk) }
func TestWalk600(t *testing.T) { test(t, 100, 5, walk) }

func TestList30(t *testing.T)  { test(t, 10, 2, list2) }
func TestList80(t *testing.T)  { test(t, 20, 3, list10) }
func TestList250(t *testing.T) { test(t, 50, 4, list20) }
func TestList600(t *testing.T) { test(t, 100, 5, list2) }

func TestGlobWorkers30(t *testing.T)  { test(t, 10, 2, globWorkersD1) }
func TestGlobWorkers80(t *testing.T)  { test(t, 20, 3, globWorkersD2) }
func TestGlobWorkers250(t *testing.T) { test(t, 50, 4, globWorkersD4) }
func TestGlobWorkers600(t *testing.T) { test(t, 100, 5, globWorkersD2) }

func BenchmarkGlob30(b *testing.B)        { benchmark(b, 10, 2, glob) }
func BenchmarkWalk30(b *testing.B)        { benchmark(b, 10, 2, walk) }
func BenchmarkList30(b *testing.B)        { benchmark(b, 10, 2, list2) }
func BenchmarkGlobWorkers30(b *testing.B) { benchmark(b, 10, 2, globWorkersD2) }

func BenchmarkGlob80(b *testing.B)        { benchmark(b, 20, 3, glob) }
func BenchmarkWalk80(b *testing.B)        { benchmark(b, 20, 3, walk) }
func BenchmarkList80(b *testing.B)        { benchmark(b, 20, 3, list2) }
func BenchmarkGlobWorkers80(b *testing.B) { benchmark(b, 20, 3, globWorkersD2) }

func BenchmarkGlob250(b *testing.B)        { benchmark(b, 50, 4, glob) }
func BenchmarkWalk250(b *testing.B)        { benchmark(b, 50, 4, walk) }
func BenchmarkList250(b *testing.B)        { benchmark(b, 50, 4, list2) }
func BenchmarkGlobWorkers250(b *testing.B) { benchmark(b, 50, 4, globWorkersD2) }

func BenchmarkGlobAT600(b *testing.B)          { benchmark(b, 100, 5, glob) }
func BenchmarkWalkAT600(b *testing.B)          { benchmark(b, 100, 5, walk) }
func BenchmarkList2AT600(b *testing.B)         { benchmark(b, 100, 5, list2) }
func BenchmarkList10AT600(b *testing.B)        { benchmark(b, 100, 5, list10) }
func BenchmarkList20AT600(b *testing.B)        { benchmark(b, 100, 5, list20) }
func BenchmarkGlobWorkersD1AT600(b *testing.B) { benchmark(b, 100, 5, globWorkersD1) }
func BenchmarkGlobWorkersD2AT600(b *testing.B) { benchmark(b, 100, 5, globWorkersD2) }
func BenchmarkGlobWorkersD4AT600(b *testing.B) { benchmark(b, 100, 5, globWorkersD4) }

func BenchmarkGlobAT22000(b *testing.B)          { benchmark(b, 2000, 10, glob) }
func BenchmarkWalkAT22000(b *testing.B)          { benchmark(b, 2000, 10, walk) }
func BenchmarkList2AT22000(b *testing.B)         { benchmark(b, 2000, 10, list2) }
func BenchmarkList10AT22000(b *testing.B)        { benchmark(b, 2000, 10, list10) }
func BenchmarkList20AT22000(b *testing.B)        { benchmark(b, 2000, 10, list20) }
func BenchmarkGlobWorkersD1AT22000(b *testing.B) { benchmark(b, 2000, 10, globWorkersD1) }
func BenchmarkGlobWorkersD2AT22000(b *testing.B) { benchmark(b, 2000, 10, globWorkersD2) }
func BenchmarkGlobWorkersD4AT22000(b *testing.B) { benchmark(b, 2000, 10, globWorkersD4) }

func BenchmarkGlob110000(b *testing.B)        { benchmark(b, 10000, 10, glob) }
func BenchmarkWalk110000(b *testing.B)        { benchmark(b, 10000, 10, walk) }
func BenchmarkList110000(b *testing.B)        { benchmark(b, 10000, 10, list2) }
func BenchmarkGlobWorkers110000(b *testing.B) { benchmark(b, 10000, 10, globWorkersD2) }

func benchmark(b *testing.B, numNodes int, numRevisions int, f func(string) <-chan string) {
	base := initialize(numNodes, numRevisions)
	defer os.RemoveAll(base)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := f(base)
		PurgeRevisions(ch, nil, false, false)
	}
	b.StopTimer()
}

func test(t *testing.T, numNodes int, numRevisions int, f func(string) <-chan string) {
	base := initialize(numNodes, numRevisions)
	defer os.RemoveAll(base)

	ch := f(base)
	_, _, revisions := PurgeRevisions(ch, nil, false, false)
	require.Equal(t, numNodes*numRevisions, revisions, "Deleted Revisions")
}

func glob(base string) <-chan string {
	return Glob(base + _basePath + "*/*/*/*/*")
}

func walk(base string) <-chan string {
	return Walk(base + _basePath)
}

func list2(base string) <-chan string {
	return List(base+_basePath, 2)
}

func list10(base string) <-chan string {
	return List(base+_basePath, 10)
}

func list20(base string) <-chan string {
	return List(base+_basePath, 20)
}

func globWorkersD1(base string) <-chan string {
	return GlobWorkers(base+_basePath, "*", "/*/*/*/*")
}

func globWorkersD2(base string) <-chan string {
	return GlobWorkers(base+_basePath, "*/*", "/*/*/*")
}

func globWorkersD4(base string) <-chan string {
	return GlobWorkers(base+_basePath, "*/*/*/*", "/*")
}

func initialize(numNodes int, numRevisions int) string {
	base := "test_temp_" + uuid.New().String()
	if err := os.Mkdir(base, os.ModePerm); err != nil {
		fmt.Println("Error creating test_temp directory", err)
		os.RemoveAll(base)
		os.Exit(1)
	}

	// create base path
	if err := os.MkdirAll(base+_basePath, fs.ModePerm); err != nil {
		fmt.Println("Error creating base path", err)
		os.RemoveAll(base)
		os.Exit(1)
	}

	for i := 0; i < numNodes; i++ {
		path := lookup.Pathify(uuid.New().String(), 4, 2)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(base+_basePath+dir, fs.ModePerm); err != nil {
			fmt.Println("Error creating test_temp directory", err)
			os.RemoveAll(base)
			os.Exit(1)
		}

		if _, err := os.Create(base + _basePath + path); err != nil {
			fmt.Println("Error creating file", err)
			os.RemoveAll(base)
			os.Exit(1)
		}
		for i := 0; i < numRevisions; i++ {
			os.Create(base + _basePath + path + ".REV.2024-05-22T07:32:53.89969" + strconv.Itoa(i) + "Z")
		}
	}
	return base
}

// revisionTimestampCase is a timestamp decomposedfs can write as a revision name.
// The stamp comes from time.Time.UTC().Format(time.RFC3339Nano), which omits the
// fractional part when the mtime falls on a whole second.
type revisionTimestampCase struct {
	name   string
	at     time.Time
	want   string
	blobID string
}

func revisionTimestampCases() []revisionTimestampCase {
	return []revisionTimestampCase{
		{
			name:   "contemporary whole second",
			at:     time.Date(2024, 5, 22, 7, 32, 53, 0, time.UTC),
			want:   "2024-05-22T07:32:53Z",
			blobID: "blob-contemporary",
		},
		{
			name:   "year 1601 whole second",
			at:     time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC),
			want:   "1601-01-01T00:00:00Z",
			blobID: "blob-1601",
		},
		{
			name:   "one fractional digit",
			at:     time.Date(2024, 6, 1, 12, 0, 0, 100000000, time.UTC),
			want:   "2024-06-01T12:00:00.1Z",
			blobID: "blob-frac-1",
		},
		{
			name:   "nine fractional digits",
			at:     time.Date(2024, 6, 1, 12, 0, 1, 123456789, time.UTC),
			want:   "2024-06-01T12:00:01.123456789Z",
			blobID: "blob-frac-9",
		},
	}
}

func TestRevisionTimestampsFollowWriterContract(t *testing.T) {
	for _, c := range revisionTimestampCases() {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, c.at.UTC().Format(time.RFC3339Nano))
		})
	}
}

const (
	revisionFixtureSpaceID = "8f638374-6ea8-4f0d-80c4-66d9b49830a5"
	revisionFixtureNodeID  = "aa111111-1111-4111-8111-111111111111"
	revisionFixtureOtherID = "bb222222-2222-4222-8222-222222222222"
)

// purgedBlob is the SpaceID and BlobID pair passed to the blobstore.
type purgedBlob struct {
	spaceID string
	blobID  string
}

// decomposedRevisionTree is a temporary decomposedfs tree:
//
//	<root>/spaces/<ss>/<space>/nodes/<n0>/<n1>/<n2>/<n3>/<node>
//	<root>/spaces/<ss>/<space>/nodes/<n0>/<n1>/<n2>/<n3>/<node>.REV.<RFC3339Nano>
//
// plus the .mpk and .mlock companion of every revision.
type decomposedRevisionTree struct {
	root           string
	revisions      int
	revisionPaths  []string
	preservedPaths []string
	blobs          []purgedBlob
}

type revisionDiscoverer struct {
	name string
	find func(root string) <-chan string
}

// revisionDiscoverers matches the mechanisms revisions purge can run, plus Walk.
func revisionDiscoverers() []revisionDiscoverer {
	return []revisionDiscoverer{
		{
			name: "glob",
			find: func(root string) <-chan string {
				return Glob(filepath.Join(root, "spaces", "*", "*", "nodes", "*", "*", "*", "*", "*"))
			},
		},
		{
			name: "globWorkers",
			find: func(root string) <-chan string {
				return GlobWorkers(filepath.Join(root, "spaces", "*", "*", "nodes"), "/*", "/*/*/*/*")
			},
		},
		{
			name: "walk",
			find: func(root string) <-chan string {
				return Walk(root)
			},
		},
		{
			name: "list",
			find: func(root string) <-chan string {
				return List(filepath.Join(root, "spaces"), 10)
			},
		},
	}
}

func decomposedNodePath(root, spaceID, nodeID string) string {
	return filepath.Join(root, "spaces", lookup.Pathify(spaceID, 1, 2), "nodes", lookup.Pathify(nodeID, 4, 2))
}

func revisionBlobMPK(t *testing.T, blobID string) []byte {
	t.Helper()
	b, err := msgpack.Marshal(map[string][]byte{"user.ocis.blobid": []byte(blobID)})
	require.NoError(t, err)
	return b
}

func writeFixtureFile(t *testing.T, path string, content []byte) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, content, 0o644))
}

// newDecomposedRevisionTree builds one space with a current node, an unrelated
// node, ordinary files, and one revision triplet per writer-contract timestamp.
func newDecomposedRevisionTree(t *testing.T) decomposedRevisionTree {
	t.Helper()

	root := t.TempDir()
	nodePath := decomposedNodePath(root, revisionFixtureSpaceID, revisionFixtureNodeID)
	otherPath := decomposedNodePath(root, revisionFixtureSpaceID, revisionFixtureOtherID)

	writeFixtureFile(t, nodePath, nil)
	writeFixtureFile(t, nodePath+".mpk", revisionBlobMPK(t, "blob-current"))
	writeFixtureFile(t, otherPath, nil)

	notes := filepath.Join(filepath.Dir(nodePath), "notes.txt")
	writeFixtureFile(t, notes, []byte("unrelated"))
	trashStamp := time.Date(2024, 5, 22, 7, 32, 53, 0, time.UTC).UTC().Format(time.RFC3339Nano)
	trash := nodePath + ".T." + trashStamp
	writeFixtureFile(t, trash, nil)
	readme := filepath.Join(root, "README.txt")
	writeFixtureFile(t, readme, []byte("unrelated"))
	ignored := filepath.Join(root, "spaces", "ignored.txt")
	writeFixtureFile(t, ignored, []byte("unrelated"))

	tree := decomposedRevisionTree{
		root: root,
		preservedPaths: []string{
			nodePath,
			nodePath + ".mpk",
			otherPath,
			notes,
			trash,
			readme,
			ignored,
		},
	}

	for _, c := range revisionTimestampCases() {
		stamp := c.at.UTC().Format(time.RFC3339Nano)
		base := nodePath + ".REV." + stamp
		writeFixtureFile(t, base, nil)
		writeFixtureFile(t, base+".mpk", revisionBlobMPK(t, c.blobID))
		writeFixtureFile(t, base+".mlock", nil)

		tree.revisions++
		tree.revisionPaths = append(tree.revisionPaths, base, base+".mpk", base+".mlock")
		tree.blobs = append(tree.blobs, purgedBlob{spaceID: revisionFixtureSpaceID, blobID: c.blobID})
	}

	sort.Strings(tree.revisionPaths)
	sort.Slice(tree.blobs, func(i, j int) bool {
		if tree.blobs[i].blobID != tree.blobs[j].blobID {
			return tree.blobs[i].blobID < tree.blobs[j].blobID
		}
		return tree.blobs[i].spaceID < tree.blobs[j].spaceID
	})
	return tree
}

// drainRevisionPaths reads a discovery channel to completion. GlobWorkers and
// List send from several goroutines, so callers compare the sorted result.
func drainRevisionPaths(ch <-chan string) []string {
	paths := make([]string, 0)
	for p := range ch {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

func TestDiscoverRevisionsWholeSecondAndFractional(t *testing.T) {
	tree := newDecomposedRevisionTree(t)

	for _, d := range revisionDiscoverers() {
		t.Run(d.name, func(t *testing.T) {
			got := drainRevisionPaths(d.find(tree.root))
			require.Equal(t, tree.revisionPaths, got)
		})
	}
}

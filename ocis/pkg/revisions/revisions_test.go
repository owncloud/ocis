package revisions

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/google/uuid"
	"github.com/test-go/testify/require"
)

var (
	_basePath = "test_temp/spaces/8f/638374-6ea8-4f0d-80c4-66d9b49830a5/nodes/"
)

// func TestInit(t *testing.T) {
// initialize(10, 2)
// defer os.RemoveAll("test_temp")
// }

func TestGlob30(t *testing.T)  { testGlob(t, 10, 2) }
func TestGlob80(t *testing.T)  { testGlob(t, 20, 3) }
func TestGlob250(t *testing.T) { testGlob(t, 50, 4) }
func TestGlob600(t *testing.T) { testGlob(t, 100, 5) }

func TestWalk30(t *testing.T)  { testWalk(t, 10, 2) }
func TestWalk80(t *testing.T)  { testWalk(t, 20, 3) }
func TestWalk250(t *testing.T) { testWalk(t, 50, 4) }
func TestWalk600(t *testing.T) { testWalk(t, 100, 5) }

func TestList30(t *testing.T)  { testList(t, 10, 2) }
func TestList80(t *testing.T)  { testList(t, 20, 3) }
func TestList250(t *testing.T) { testList(t, 50, 4) }
func TestList600(t *testing.T) { testList(t, 100, 5) }

func BenchmarkGlob30(b *testing.B) { benchmarkGlob(b, 10, 2) }
func BenchmarkWalk30(b *testing.B) { benchmarkWalk(b, 10, 2) }
func BenchmarkList30(b *testing.B) { benchmarkList(b, 10, 2) }

func BenchmarkGlob80(b *testing.B) { benchmarkGlob(b, 20, 3) }
func BenchmarkWalk80(b *testing.B) { benchmarkWalk(b, 20, 3) }
func BenchmarkList80(b *testing.B) { benchmarkList(b, 20, 3) }

func BenchmarkGlob250(b *testing.B) { benchmarkGlob(b, 50, 4) }
func BenchmarkWalk250(b *testing.B) { benchmarkWalk(b, 50, 4) }
func BenchmarkList250(b *testing.B) { benchmarkList(b, 50, 4) }

func BenchmarkGlob600(b *testing.B) { benchmarkGlob(b, 100, 5) }
func BenchmarkWalk600(b *testing.B) { benchmarkWalk(b, 100, 5) }
func BenchmarkList600(b *testing.B) { benchmarkList(b, 100, 5) }

func BenchmarkGlob11000(b *testing.B) { benchmarkGlob(b, 1000, 10) }
func BenchmarkWalk11000(b *testing.B) { benchmarkWalk(b, 1000, 10) }
func BenchmarkList11000(b *testing.B) { benchmarkList(b, 1000, 10) }

func BenchmarkGlob110000(b *testing.B) { benchmarkGlob(b, 10000, 10) }
func BenchmarkWalk110000(b *testing.B) { benchmarkWalk(b, 10000, 10) }
func BenchmarkList110000(b *testing.B) { benchmarkList(b, 10000, 10) }

func benchmarkGlob(b *testing.B, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	for i := 0; i < b.N; i++ {
		PurgeRevisionsGlob(_basePath+"*/*/*/*/*", nil, false, false)
	}
}

func benchmarkWalk(b *testing.B, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	for i := 0; i < b.N; i++ {
		PurgeRevisionsWalk(_basePath, nil, false, false)
	}
}

func benchmarkList(b *testing.B, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	for i := 0; i < b.N; i++ {
		PurgeRevisionsList(_basePath, nil, false, false)
	}
}

func testGlob(t *testing.T, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	_, _, revisions := PurgeRevisionsGlob(_basePath+"*/*/*/*/*", nil, false, false)
	require.Equal(t, numNodes*numRevisions, revisions, "Deleted Revisions")
}

func testWalk(t *testing.T, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	_, _, revisions := PurgeRevisionsWalk(_basePath, nil, false, false)
	require.Equal(t, numNodes*numRevisions, revisions, "Deleted Revisions")
}

func testList(t *testing.T, numNodes int, numRevisions int) {
	initialize(numNodes, numRevisions)
	defer os.RemoveAll("test_temp")

	_, _, revisions := PurgeRevisionsList(_basePath, nil, false, false)
	require.Equal(t, numNodes*numRevisions, revisions, "Deleted Revisions")
}

func initialize(numNodes int, numRevisions int) {
	// create base path
	if err := os.MkdirAll(_basePath, fs.ModePerm); err != nil {
		fmt.Println("Error creating test_temp directory", err)
		os.Exit(1)
	}

	for i := 0; i < numNodes; i++ {
		path := lookup.Pathify(uuid.New().String(), 4, 2)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(_basePath+dir, fs.ModePerm); err != nil {
			fmt.Println("Error creating test_temp directory", err)
			os.Exit(1)
		}

		if _, err := os.Create(_basePath + path); err != nil {
			fmt.Println("Error creating file", err)
			os.Exit(1)
		}
		for i := 0; i < numRevisions; i++ {
			os.Create(_basePath + path + ".REV.2024-05-22T07:32:53.89969" + strconv.Itoa(i) + "Z")
		}
	}
}

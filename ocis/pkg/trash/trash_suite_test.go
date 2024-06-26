package trash_test

import (
	"io/fs"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testhelper "github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/testhelpers"
	"github.com/owncloud/ocis/v2/ocis/pkg/trash"
)

func TestTrash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trash Suite")
}

// shared test helpers

func setupTestEnv(root string) *testhelper.TestEnv {
	Expect(os.MkdirAll(root, 0755)).ToNot(HaveOccurred())

	env, err := testhelper.NewTestEnv(map[string]interface{}{
		"root": root,
	})
	Expect(err).ToNot(HaveOccurred())

	return env
}

func abandonTrashedNodes(root, glob string) []trash.TrashDirs {
	dirs, err := fs.Glob(os.DirFS(root), glob)
	Expect(err).ToNot(HaveOccurred())

	var deletions []trash.TrashDirs
	for _, p := range dirs {
		td := trash.TrashDirs{
			// absolute path to the symlink
			LinkPath: path.Join(root, p),
		}

		t, err := os.Readlink(td.LinkPath)
		Expect(err).ToNot(HaveOccurred())

		// relative path to the node starting from the testRoot
		td.NodePath = path.Join(path.Dir(p), t)

		Expect(os.Remove(path.Join(root, td.NodePath))).ToNot(HaveOccurred())
		deletions = append(deletions, td)
	}

	return deletions
}

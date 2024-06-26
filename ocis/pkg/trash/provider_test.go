package trash_test

import (
	"os"
	"path"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/stretchr/testify/mock"

	testhelper "github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/testhelpers"
	"github.com/owncloud/ocis/v2/ocis/pkg/trash"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DataProvider", func() {
	var (
		env           *testhelper.TestEnv
		testRoot      string
		abandonedDirs []trash.TrashDirs
	)

	BeforeEach(func() {
		tr, err := os.MkdirTemp("", "ocis-test")
		Expect(err).ToNot(HaveOccurred())

		testRoot = tr
		env = setupTestEnv(path.Join(testRoot, "storage", "users"))
	})

	Context("ProduceData", func() {
		It("reports trashed directories even if the target node is abandoned", func() {
			env.Permissions.EXPECT().AssemblePermissions(mock.Anything, mock.Anything).Return(provider.ResourcePermissions{
				Stat:            true,
				CreateContainer: true,
				Delete:          true,
			}, nil)

			Expect(env.Fs.Delete(env.Ctx, &provider.Reference{
				ResourceId: env.SpaceRootRes,
				Path:       "/dir1/subdir1",
			})).ToNot(HaveOccurred())

			abandonedDirs = abandonTrashedNodes(testRoot, trash.TrashGlobPattern)
			Expect(abandonedDirs).To(HaveLen(1))

			dataProvider := trash.NewDataProvider(os.DirFS(testRoot), testRoot)
			Expect(dataProvider.ProduceData()).ToNot(HaveOccurred())

			e := <-dataProvider.Events
			// mixed paths
			// e.(trash.TrashDirs).LinkPath is absolute at the moment, is it necessary to make it relative?
			// e.(trash.TrashDirs).NodePath is relative at the moment, is it necessary to make it absolute?
			occurrences := []trash.TrashDirs{e.(trash.TrashDirs)}
			Expect(occurrences).To(Equal(abandonedDirs))
			Expect(occurrences).To(HaveLen(1))
		})
	})
})

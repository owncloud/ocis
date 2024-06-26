package trash_test

import (
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	testhelper "github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/testhelpers"
	"github.com/owncloud/ocis/v2/ocis/pkg/trash"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Trash", func() {
	var (
		env *testhelper.TestEnv
	)

	BeforeEach(func() {
		var err error
		env, err = testhelper.NewTestEnv(map[string]interface{}{
			"root":
		})
		Expect(err).ToNot(HaveOccurred())
	})

	Context("No empty trash directories", func() {
		When("a directory is removed", func() {
			JustBeforeEach(func() {
				env.Permissions.On("AssemblePermissions", mock.Anything, mock.Anything, mock.Anything).Return(provider.ResourcePermissions{
					Stat:            true,
					CreateContainer: true,
					Delete:          true,
				}, nil)

				err := env.Fs.Delete(env.Ctx, &provider.Reference{
					ResourceId: env.SpaceRootRes,
					Path:       "/dir1/subdir1",
				})
				Expect(err).ToNot(HaveOccurred())
			})
			It("does not find any trash dirs to remove", func() {
				err := trash.PurgeTrashOrphanedPaths(env.Root)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

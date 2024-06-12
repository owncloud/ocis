package assignmentscache_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cs3org/reva/v2/pkg/storage/utils/metadata"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/metadata/assignmentscache"
)

var _ = Describe("Assignmentscache", func() {
	var (
		c       assignmentscache.Cache
		storage metadata.Storage

		roleid     = "11111111-1111-1111-1111-111111111111"
		assignment = &settingsmsg.UserRoleAssignment{
			AccountUuid: "00000000-0000-0000-0000-000000000001",
			RoleId:      roleid,
			Id:          "00000001-0000-0000-0000-000000000000",
		}
		ctx    context.Context
		tmpdir string
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		tmpdir, err = os.MkdirTemp("", "assignmentscache-test")
		Expect(err).ToNot(HaveOccurred())

		err = os.MkdirAll(tmpdir, 0755)
		Expect(err).ToNot(HaveOccurred())

		storage, err = metadata.NewDiskStorage(tmpdir)
		Expect(err).ToNot(HaveOccurred())

		c = assignmentscache.New(storage, "basename", "assignments.json")
		Expect(c).ToNot(BeNil()) //nolint:all
	})

	AfterEach(func() {
		if tmpdir != "" {
			os.RemoveAll(tmpdir)
		}
	})

	Describe("Persist", func() {
		Context("with an existing entry", func() {
			BeforeEach(func() {
				Expect(c.Add(ctx, roleid, assignment)).To(Succeed())
			})

			It("updates the etag", func() {
				ra, _ := c.RoleAssignments.Load(roleid)
				oldEtag := ra.Etag
				Expect(oldEtag).ToNot(BeEmpty())

				Expect(c.Persist(ctx, roleid)).To(Succeed())

				ra, _ = c.RoleAssignments.Load(roleid)
				Expect(ra.Etag).ToNot(Equal(oldEtag))
			})
		})
	})

})

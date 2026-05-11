// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/client_test.go
package kiteworks_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks"
)

var _ = Describe("Client", func() {
	var (
		srv    *httptest.Server
		client *Client
	)

	BeforeEach(func() {
		srv = newMockServer()
		client = NewClient(srv.URL, "test-token", false)
	})

	AfterEach(func() {
		srv.Close()
	})

	Describe("GetTopFolders", func() {
		It("returns two top-level folders", func() {
			folders, err := client.GetTopFolders()
			Expect(err).ToNot(HaveOccurred())
			Expect(folders).To(HaveLen(2))
			Expect(folders[0].ID).To(Equal("f1"))
			Expect(folders[1].ID).To(Equal("f2"))
		})
	})

	Describe("GetFolder", func() {
		It("returns folder metadata", func() {
			fi, err := client.GetFolder("f1")
			Expect(err).ToNot(HaveOccurred())
			Expect(fi.ID).To(Equal("f1"))
			Expect(fi.Name).To(Equal("MyFiles"))
		})
	})

	Describe("ListFolder", func() {
		It("returns folder children", func() {
			dir, err := client.ListFolder("f1")
			Expect(err).ToNot(HaveOccurred())
			Expect(dir.Files).To(HaveLen(1))
			Expect(dir.Files[0].Name).To(Equal("hello.txt"))
		})
	})

	Describe("GetMe", func() {
		It("returns current user with quota", func() {
			u, err := client.GetMe()
			Expect(err).ToNot(HaveOccurred())
			Expect(u.ID).To(Equal("u1"))
			Expect(u.Quota.Allowed).To(Equal(int64(10737418240)))
		})
	})

	Describe("IsSharedWithUser", func() {
		It("returns false for owned folder", func() {
			folders, err := client.GetTopFolders()
			Expect(err).ToNot(HaveOccurred())
			Expect(folders[0].IsSharedWithUser()).To(BeFalse())
		})
		It("returns true for received share", func() {
			folders, err := client.GetTopFolders()
			Expect(err).ToNot(HaveOccurred())
			Expect(folders[1].IsSharedWithUser()).To(BeTrue())
		})
	})
})

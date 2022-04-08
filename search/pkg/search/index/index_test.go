package index_test

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/search/pkg/search"
	"github.com/owncloud/ocis/search/pkg/search/index"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Index", func() {
	var (
		i          *index.Index
		bleveIndex bleve.Index
		ref        *sprovider.Reference
		ri         *sprovider.ResourceInfo

		ctx context.Context
	)

	BeforeEach(func() {
		var err error
		bleveIndex, err = bleve.NewMemOnly(index.BuildMapping())
		Expect(err).ToNot(HaveOccurred())

		i, err = index.New(bleveIndex)
		Expect(err).ToNot(HaveOccurred())

		ref = &sprovider.Reference{
			ResourceId: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "rootopaqueid",
			},
		}
		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			Path: "foo.pdf",
			Size: 12345,
		}
	})

	Describe("New", func() {
		It("returns a new index instance", func() {
			i, err := index.New(bleveIndex)
			Expect(err).ToNot(HaveOccurred())
			Expect(i).ToNot(BeNil())
		})
	})

	Describe("NewPersisted", func() {
		It("returns a new index instance", func() {
			i, err := index.NewPersisted("")
			Expect(err).ToNot(HaveOccurred())
			Expect(i).ToNot(BeNil())
		})
	})

	Describe("Search", func() {
		It("finds files by prefix", func() {
			err := i.Add(ref, ri)
			Expect(err).ToNot(HaveOccurred())

			res, err := i.Search(ctx, &search.SearchIndexRequest{
				Reference: &sprovider.Reference{
					ResourceId: ref.ResourceId,
				},
				Query: "foo.pdf",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
			Expect(len(res.Matches)).To(Equal(1))
			Expect(res.Matches[0].Reference.ResourceId).To(Equal(ref.ResourceId))
			Expect(res.Matches[0].Info.Id).To(Equal(ri.Id))
			Expect(res.Matches[0].Info.Path).To(Equal(ri.Path))
			Expect(res.Matches[0].Info.Size).To(Equal(ri.Size))
		})

		PIt("finds files living deeper in the tree by prefix")
		PIt("finds directories by prefix")
		PIt("finds directories living deeper in the tree by prefix")
	})

	Describe("Scan", func() {
		PIt("adds the given resource recursively")
	})

	Describe("Index", func() {
		It("adds a resourceInfo to the index", func() {
			err := i.Add(ref, ri)
			Expect(err).ToNot(HaveOccurred())

			count, err := bleveIndex.DocCount()
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(uint64(1)))

			query := bleve.NewMatchQuery("foo.pdf")
			res, err := bleveIndex.Search(bleve.NewSearchRequest(query))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Hits.Len()).To(Equal(1))
		})

		It("updates an existing resource in the index", func() {
			err := i.Add(ref, ri)
			Expect(err).ToNot(HaveOccurred())
			count, _ := bleveIndex.DocCount()
			Expect(count).To(Equal(uint64(1)))

			err = i.Add(ref, ri)
			Expect(err).ToNot(HaveOccurred())
			count, _ = bleveIndex.DocCount()
			Expect(count).To(Equal(uint64(1)))
		})
	})

	Describe("Remove", func() {
		PIt("removes a resource from the index")
	})
})

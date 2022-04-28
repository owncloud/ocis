package index_test

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/extensions/search/pkg/search/index"
	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"

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
			Path: "./foo.pdf",
		}
		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			Path:     "foo.pdf",
			Size:     12345,
			MimeType: "application/pdf",
			Mtime:    &typesv1beta1.Timestamp{Seconds: 4000},
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
		Context("with a file in the root of the space", func() {
			BeforeEach(func() {
				err := i.Add(ref, ri)
				Expect(err).ToNot(HaveOccurred())
			})

			It("scopes the search to the specified space", func() {
				res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
					Ref: &searchmsg.Reference{
						ResourceId: &searchmsg.ResourceID{
							StorageId: "differentstorageid",
							OpaqueId:  "differentopaqueid",
						},
					},
					Query: "foo.pdf",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(0))
			})

			It("limits the search to the relevant fields", func() {
				res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
					Ref: &searchmsg.Reference{
						ResourceId: &searchmsg.ResourceID{
							StorageId: ref.ResourceId.StorageId,
							OpaqueId:  ref.ResourceId.OpaqueId,
						},
					},
					Query: "*" + ref.ResourceId.OpaqueId + "*",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(0))
			})

			It("returns all desired fields", func() {
				res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
					Ref: &searchmsg.Reference{
						ResourceId: &searchmsg.ResourceID{
							StorageId: ref.ResourceId.StorageId,
							OpaqueId:  ref.ResourceId.OpaqueId,
						},
					},
					Query: "foo.pdf",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(ref.ResourceId.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal(ref.Path))
				Expect(match.Entity.Id.OpaqueId).To(Equal(ri.Id.OpaqueId))
				Expect(match.Entity.Name).To(Equal(ri.Path))
				Expect(match.Entity.Size).To(Equal(ri.Size))
				Expect(match.Entity.MimeType).To(Equal(ri.MimeType))
				Expect(uint64(match.Entity.LastModifiedTime.AsTime().Unix())).To(Equal(ri.Mtime.Seconds))
			})

			It("finds files by name, prefix or substring match", func() {
				queries := []string{"foo.pdf", "foo*", "*oo.p*"}
				for _, query := range queries {
					res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
						Ref: &searchmsg.Reference{
							ResourceId: &searchmsg.ResourceID{
								StorageId: ref.ResourceId.StorageId,
								OpaqueId:  ref.ResourceId.OpaqueId,
							},
						},
						Query: query,
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(1), "query returned no result: "+query)
					Expect(res.Matches[0].Entity.Ref.ResourceId.OpaqueId).To(Equal(ref.ResourceId.OpaqueId))
					Expect(res.Matches[0].Entity.Ref.Path).To(Equal(ref.Path))
					Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal(ri.Id.OpaqueId))
					Expect(res.Matches[0].Entity.Name).To(Equal(ri.Path))
					Expect(res.Matches[0].Entity.Size).To(Equal(ri.Size))
				}
			})

			Context("and an additional file in a subdirectory", func() {
				var (
					nestedRef *sprovider.Reference
					nestedRI  *sprovider.ResourceInfo
				)

				BeforeEach(func() {
					nestedRef = &sprovider.Reference{
						ResourceId: &sprovider.ResourceId{
							StorageId: "storageid",
							OpaqueId:  "rootopaqueid",
						},
						Path: "./nested/nestedpdf.pdf",
					}
					nestedRI = &sprovider.ResourceInfo{
						Id: &sprovider.ResourceId{
							StorageId: "storageid",
							OpaqueId:  "nestedopaqueid",
						},
						Path: "nestedpdf.pdf",
						Size: 12345,
					}
					err := i.Add(nestedRef, nestedRI)
					Expect(err).ToNot(HaveOccurred())
				})

				It("finds files living deeper in the tree by filename, prefix or substring match", func() {
					queries := []string{"nestedpdf.pdf", "nested*", "*tedpdf.*"}
					for _, query := range queries {
						res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
							Ref: &searchmsg.Reference{
								ResourceId: &searchmsg.ResourceID{
									StorageId: ref.ResourceId.StorageId,
									OpaqueId:  ref.ResourceId.OpaqueId,
								},
							},
							Query: query,
						})
						Expect(err).ToNot(HaveOccurred())
						Expect(res).ToNot(BeNil())
						Expect(len(res.Matches)).To(Equal(1), "query returned no result: "+query)
					}
				})

				It("does not find the higher levels when limiting the searched directory", func() {
					res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
						Ref: &searchmsg.Reference{
							ResourceId: &searchmsg.ResourceID{
								StorageId: ref.ResourceId.StorageId,
								OpaqueId:  ref.ResourceId.OpaqueId,
							},
							Path: "./nested/",
						},
						Query: "foo.pdf",
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(0))
				})
			})
		})
	})

	Describe("Add", func() {
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
		It("removes a resource from the index", func() {
			err := i.Add(ref, ri)
			Expect(err).ToNot(HaveOccurred())
			count, _ := bleveIndex.DocCount()
			Expect(count).To(Equal(uint64(1)))

			err = i.Remove(ri.Id)
			Expect(err).ToNot(HaveOccurred())
			count, _ = bleveIndex.DocCount()
			Expect(count).To(Equal(uint64(0)))
		})
	})
})

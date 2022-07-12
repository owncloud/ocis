package index_test

import (
	"context"

	"github.com/blevesearch/bleve/v2"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/search/index"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Index", func() {
	var (
		i          *index.Index
		bleveIndex bleve.Index
		ctx        context.Context

		rootId = &sprovider.ResourceId{
			StorageId: "provider-1",
			SpaceId:   "spaceid",
			OpaqueId:  "rootopaqueid",
		}
		filename  string
		ref       *sprovider.Reference
		ri        *sprovider.ResourceInfo
		parentRef = &sprovider.Reference{
			ResourceId: rootId,
			Path:       "./my/sub d!r",
		}
		parentRi = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "provider-1",
				SpaceId:   "spaceid",
				OpaqueId:  "parentopaqueid",
			},
			Path:  "sub d!r",
			Size:  12345,
			Type:  sprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
			Mtime: &typesv1beta1.Timestamp{Seconds: 4000},
		}
		childRef = &sprovider.Reference{
			ResourceId: rootId,
			Path:       "./my/sub d!r/child.pdf",
		}
		childRi = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "provider-1",
				SpaceId:   "spaceid",
				OpaqueId:  "childopaqueid",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "provider-1",
				SpaceId:   "spaceid",
				OpaqueId:  "parentopaqueid",
			},
			Path:  "child.pdf",
			Size:  12345,
			Type:  sprovider.ResourceType_RESOURCE_TYPE_FILE,
			Mtime: &typesv1beta1.Timestamp{Seconds: 4000},
		}

		assertDocCount = func(rootId *sprovider.ResourceId, query string, expectedCount int) []*searchmsg.Match {
			res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
				Query: query,
				Ref: &searchmsg.Reference{
					ResourceId: &searchmsg.ResourceID{
						StorageId: "provider-1",
						SpaceId:   rootId.SpaceId,
						OpaqueId:  rootId.OpaqueId,
					},
				},
			})
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, len(res.Matches)).To(Equal(expectedCount), "query returned unexpected number of results: "+query)
			return res.Matches
		}
	)

	BeforeEach(func() {
		filename = "Foo.pdf"

		mapping, err := index.BuildMapping()
		Expect(err).ToNot(HaveOccurred())

		bleveIndex, err = bleve.NewMemOnly(mapping)
		Expect(err).ToNot(HaveOccurred())

		i, err = index.New(bleveIndex)
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		ref = &sprovider.Reference{
			ResourceId: rootId,
			Path:       "./" + filename,
		}
		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "provider-1",
				SpaceId:   "spaceid",
				OpaqueId:  "opaqueid",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "provider-1",
				SpaceId:   "spaceid",
				OpaqueId:  "someopaqueid",
			},
			Path:     filename,
			Size:     12345,
			Type:     sprovider.ResourceType_RESOURCE_TYPE_FILE,
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
		Context("by other fields than filename", func() {
			JustBeforeEach(func() {
				err := i.Add(ref, ri)
				Expect(err).ToNot(HaveOccurred())
			})

			It("finds files by size", func() {
				assertDocCount(ref.ResourceId, `Size:12345`, 1)
				assertDocCount(ref.ResourceId, `Size:>1000`, 1)
				assertDocCount(ref.ResourceId, `Size:<100000`, 1)

				assertDocCount(ref.ResourceId, `Size:12344`, 0)
				assertDocCount(ref.ResourceId, `Size:<1000`, 0)
				assertDocCount(ref.ResourceId, `Size:>100000`, 0)
			})
		})

		Context("by filename", func() {
			It("finds files with spaces in the filename", func() {
				ri.Path = "Foo oo.pdf"
				ref.Path = "./" + ri.Path
				err := i.Add(ref, ri)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(ref.ResourceId, `Name:foo\ o*`, 1)
			})

			It("finds files by digits in the filename", func() {
				ri.Path = "12345.pdf"
				ref.Path = "./" + ri.Path
				err := i.Add(ref, ri)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(ref.ResourceId, `Name:1234*`, 1)
			})

			Context("with a file in the root of the space", func() {
				JustBeforeEach(func() {
					err := i.Add(ref, ri)
					Expect(err).ToNot(HaveOccurred())
				})

				It("scopes the search to the specified space", func() {
					resourceId := &sprovider.ResourceId{
						StorageId: "provider-1",
						SpaceId:   "differentspaceid",
						OpaqueId:  "differentopaqueid",
					}
					assertDocCount(resourceId, `Name:foo.pdf`, 0)
				})

				It("limits the search to the specified fields", func() {
					assertDocCount(ref.ResourceId, "Name:*"+ref.ResourceId.OpaqueId+"*", 0)
				})

				It("returns all desired fields", func() {
					matches := assertDocCount(ref.ResourceId, "Name:foo.pdf", 1)
					match := matches[0]
					Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(ref.ResourceId.OpaqueId))
					Expect(match.Entity.Ref.Path).To(Equal(ref.Path))
					Expect(match.Entity.Id.OpaqueId).To(Equal(ri.Id.OpaqueId))
					Expect(match.Entity.Name).To(Equal(ri.Path))
					Expect(match.Entity.Size).To(Equal(ri.Size))
					Expect(match.Entity.Type).To(Equal(uint64(ri.Type)))
					Expect(match.Entity.MimeType).To(Equal(ri.MimeType))
					Expect(match.Entity.Deleted).To(BeFalse())
					Expect(uint64(match.Entity.LastModifiedTime.AsTime().Unix())).To(Equal(ri.Mtime.Seconds))
				})

				It("finds files by name, prefix or substring match", func() {
					queries := []string{"foo.pdf", "foo*", "*oo.p*"}
					for _, query := range queries {
						matches := assertDocCount(ref.ResourceId, query, 1)
						Expect(matches[0].Entity.Ref.ResourceId.OpaqueId).To(Equal(ref.ResourceId.OpaqueId))
						Expect(matches[0].Entity.Ref.Path).To(Equal(ref.Path))
						Expect(matches[0].Entity.Id.OpaqueId).To(Equal(ri.Id.OpaqueId))
						Expect(matches[0].Entity.Name).To(Equal(ri.Path))
						Expect(matches[0].Entity.Size).To(Equal(ri.Size))
					}
				})

				It("uses a lower-case index", func() {
					assertDocCount(ref.ResourceId, "Name:foo*", 1)
					assertDocCount(ref.ResourceId, "Name:Foo*", 0)
				})

				Context("and an additional file in a subdirectory", func() {
					var (
						nestedRef *sprovider.Reference
						nestedRI  *sprovider.ResourceInfo
					)

					BeforeEach(func() {
						nestedRef = &sprovider.Reference{
							ResourceId: &sprovider.ResourceId{
								StorageId: "provider-1",
								SpaceId:   "spaceid",
								OpaqueId:  "rootopaqueid",
							},
							Path: "./nested/nestedpdf.pdf",
						}
						nestedRI = &sprovider.ResourceInfo{
							Id: &sprovider.ResourceId{
								StorageId: "provider-1",
								SpaceId:   "spaceid",
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
							assertDocCount(ref.ResourceId, query, 1)
						}
					})

					It("does not find the higher levels when limiting the searched directory", func() {
						res, err := i.Search(ctx, &searchsvc.SearchIndexRequest{
							Ref: &searchmsg.Reference{
								ResourceId: &searchmsg.ResourceID{
									StorageId: ref.ResourceId.StorageId,
									SpaceId:   ref.ResourceId.SpaceId,
									OpaqueId:  ref.ResourceId.OpaqueId,
								},
								Path: "./nested/",
							},
							Query: "Name:foo.pdf",
						})
						Expect(err).ToNot(HaveOccurred())
						Expect(res).ToNot(BeNil())
						Expect(len(res.Matches)).To(Equal(0))
					})
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

	Describe("Delete", func() {
		It("marks a resource as deleted", func() {
			err := i.Add(parentRef, parentRi)
			Expect(err).ToNot(HaveOccurred())
			assertDocCount(rootId, `sub\ d!r`, 1)

			err = i.Delete(parentRi.Id)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d!r`, 0)
		})

		It("also marks child resources as deleted", func() {
			err := i.Add(parentRef, parentRi)
			Expect(err).ToNot(HaveOccurred())
			err = i.Add(childRef, childRi)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d\!r`, 1)
			assertDocCount(rootId, "child.pdf", 1)

			err = i.Delete(parentRi.Id)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d\!r`, 0)
			assertDocCount(rootId, "child.pdf", 0)
		})
	})

	Describe("Restore", func() {
		It("also marks child resources as restored", func() {
			err := i.Add(parentRef, parentRi)
			Expect(err).ToNot(HaveOccurred())
			err = i.Add(childRef, childRi)
			Expect(err).ToNot(HaveOccurred())
			err = i.Delete(parentRi.Id)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d!r`, 0)
			assertDocCount(rootId, "child.pdf", 0)

			err = i.Restore(parentRi.Id)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d!r`, 1)
			assertDocCount(rootId, "child.pdf", 1)
		})
	})

	Describe("Move", func() {
		It("moves the parent and its child resources", func() {
			err := i.Add(parentRef, parentRi)
			Expect(err).ToNot(HaveOccurred())
			err = i.Add(childRef, childRi)
			Expect(err).ToNot(HaveOccurred())

			parentRi.Path = "newname"
			err = i.Move(parentRi.Id, "./somewhere/else/newname")
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootId, `sub\ d!r`, 0)

			matches := assertDocCount(rootId, "Name:child.pdf", 1)
			Expect(matches[0].Entity.Ref.Path).To(Equal("./somewhere/else/newname/child.pdf"))
		})
	})
})

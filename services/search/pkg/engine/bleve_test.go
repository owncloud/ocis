package engine_test

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/storagespace"

	"github.com/blevesearch/bleve/v2"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
)

var _ = Describe("Bleve", func() {
	var (
		eng *engine.Bleve
		idx bleve.Index
		ctx context.Context

		doSearch = func(id string, query, path string) (*searchsvc.SearchIndexResponse, error) {
			rID, err := storagespace.ParseID(id)
			if err != nil {
				return nil, err
			}

			return eng.Search(ctx, &searchsvc.SearchIndexRequest{
				Query: query,
				Ref: &searchmsg.Reference{
					ResourceId: &searchmsg.ResourceID{
						StorageId: rID.StorageId,
						SpaceId:   rID.SpaceId,
						OpaqueId:  rID.OpaqueId,
					},
					Path: path,
				},
			})
		}

		assertDocCount = func(id string, query string, expectedCount int) []*searchmsg.Match {
			res, err := doSearch(id, query, "")

			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, len(res.Matches)).To(Equal(expectedCount), "query returned unexpected number of results: "+query)
			return res.Matches
		}

		rootResource   engine.Resource
		parentResource engine.Resource
		childResource  engine.Resource
	)

	BeforeEach(func() {
		mapping, err := engine.BuildBleveMapping()
		Expect(err).ToNot(HaveOccurred())

		idx, err = bleve.NewMemOnly(mapping)
		Expect(err).ToNot(HaveOccurred())

		eng = engine.NewBleveEngine(idx)
		Expect(err).ToNot(HaveOccurred())

		rootResource = engine.Resource{
			ID:       "1$2!2",
			RootID:   "1$2!2",
			Path:     ".",
			Document: content.Document{},
		}

		parentResource = engine.Resource{
			ID:       "1$2!3",
			ParentID: rootResource.ID,
			RootID:   rootResource.ID,
			Path:     "./parent d!r",
			Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_CONTAINER),
			Document: content.Document{Name: "parent d!r"},
		}

		childResource = engine.Resource{
			ID:       "1$2!4",
			ParentID: parentResource.ID,
			RootID:   rootResource.ID,
			Path:     "./parent d!r/child.pdf",
			Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_FILE),
			Document: content.Document{Name: "child.pdf"},
		}
	})

	Describe("New", func() {
		It("returns a new index instance", func() {
			b := engine.NewBleveEngine(idx)
			Expect(b).ToNot(BeNil())
		})
	})

	Describe("Search", func() {
		Context("by other fields than filename", func() {
			It("finds files by tags", func() {
				parentResource.Document.Tags = []string{"foo", "bar"}
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Tags:foo", 1)
				assertDocCount(rootResource.ID, "Tags:bar", 1)
				assertDocCount(rootResource.ID, "Tags:foo Tags:bar", 1)
				assertDocCount(rootResource.ID, "Tags:foo Tags:bar Tags:baz", 1)
				assertDocCount(rootResource.ID, "Tags:foo Tags:bar Tags:baz", 1)
				assertDocCount(rootResource.ID, "Tags:baz", 0)
			})

			It("finds files by size", func() {
				parentResource.Document.Size = 12345
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Size:12345", 1)
				assertDocCount(rootResource.ID, "Size:>1000", 1)
				assertDocCount(rootResource.ID, "Size:<100000", 1)
				assertDocCount(rootResource.ID, "Size:12344", 0)
				assertDocCount(rootResource.ID, "Size:<1000", 0)
				assertDocCount(rootResource.ID, "Size:>100000", 0)
			})
		})

		Context("by filename", func() {
			It("finds files with spaces in the filename", func() {
				parentResource.Document.Name = "Foo oo.pdf"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, `Name:foo\ o*`, 1)
			})

			It("finds files by digits in the filename", func() {
				parentResource.Document.Name = "12345.pdf"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Name:1234*", 1)
			})

			It("filters hidden files", func() {
				childResource.Hidden = true
				err := eng.Upsert(childResource.ID, childResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Hidden:T", 1)
				assertDocCount(rootResource.ID, "Hidden:F", 0)
			})

			Context("with a file in the root of the space", func() {
				It("scopes the search to the specified space", func() {
					parentResource.Document.Name = "foo.pdf"
					err := eng.Upsert(parentResource.ID, parentResource)
					Expect(err).ToNot(HaveOccurred())

					assertDocCount(rootResource.ID, "Name:foo.pdf", 1)
					assertDocCount("9$8!7", "Name:foo.pdf", 0)
				})
			})

			It("limits the search to the specified fields", func() {
				parentResource.Document.Name = "bar.pdf"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Name:bar.pdf", 1)
				assertDocCount(rootResource.ID, "Unknown:field", 0)
			})

			It("returns the total number of hits", func() {
				parentResource.Document.Name = "bar.pdf"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				res, err := doSearch(rootResource.ID, "Name:bar*", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(1)))
			})

			It("returns all desired fields", func() {
				parentResource.Document.Name = "bar.pdf"
				parentResource.Type = 3
				parentResource.MimeType = "application/pdf"

				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				matches := assertDocCount(rootResource.ID, fmt.Sprintf("Name:%s", parentResource.Name), 1)
				match := matches[0]
				Expect(match.Entity.Ref.Path).To(Equal(parentResource.Path))
				Expect(match.Entity.Name).To(Equal(parentResource.Name))
				Expect(match.Entity.Size).To(Equal(parentResource.Size))
				Expect(match.Entity.Type).To(Equal(parentResource.Type))
				Expect(match.Entity.MimeType).To(Equal(parentResource.MimeType))
				Expect(match.Entity.Deleted).To(BeFalse())
				Expect(match.Score > 0).To(BeTrue())
			})

			It("finds files by name, prefix or substring match", func() {
				parentResource.Document.Name = "foo.pdf"

				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				queries := []string{"foo.pdf", "foo*", "*oo.p*"}
				for _, query := range queries {
					err := eng.Upsert(parentResource.ID, parentResource)
					Expect(err).ToNot(HaveOccurred())

					assertDocCount(rootResource.ID, query, 1)
				}
			})

			It("does a case-insensitive search", func() {
				parentResource.Document.Name = "foo.pdf"

				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rootResource.ID, "Name:foo*", 1)
				assertDocCount(rootResource.ID, "Name:Foo*", 1)
			})

			Context("and an additional file in a subdirectory", func() {
				BeforeEach(func() {
					err := eng.Upsert(parentResource.ID, parentResource)
					Expect(err).ToNot(HaveOccurred())

					err = eng.Upsert(childResource.ID, childResource)
					Expect(err).ToNot(HaveOccurred())
				})

				It("finds files living deeper in the tree by filename, prefix or substring match", func() {
					queries := []string{"child.pdf", "child*", "*ld.*"}
					for _, query := range queries {
						assertDocCount(rootResource.ID, query, 1)
					}
				})
			})
		})

		Context("Highlights", func() {

			It("highlights only for content searches", func() {
				parentResource.Document.Name = "baz.pdf"
				parentResource.Document.Content = "foo bar baz"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				res, err := doSearch(rootResource.ID, "Name:baz*", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(1)))
				Expect(res.Matches[0].Entity.Highlights).To(Equal(""))
			})

			It("highlights search terms", func() {
				parentResource.Document.Name = "baz.pdf"
				parentResource.Document.Content = "foo bar baz"
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				res, err := doSearch(rootResource.ID, "Content:bar", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(1)))
				Expect(res.Matches[0].Entity.Highlights).To(Equal("foo <mark>bar</mark> baz"))
			})

		})

		Context("with a file in the root of the space and folder with a file. all of them have the same name", func() {
			BeforeEach(func() {
				parentResource := engine.Resource{
					ID:       "1$2!3",
					ParentID: rootResource.ID,
					RootID:   rootResource.ID,
					Path:     "./doc",
					Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_CONTAINER),
					Document: content.Document{Name: "doc"},
				}

				childResource := engine.Resource{
					ID:       "1$2!4",
					ParentID: parentResource.ID,
					RootID:   rootResource.ID,
					Path:     "./doc/doc.pdf",
					Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_FILE),
					Document: content.Document{Name: "doc.pdf"},
				}

				childResource2 := engine.Resource{
					ID:       "1$2!7",
					ParentID: parentResource.ID,
					RootID:   rootResource.ID,
					Path:     "./doc/file.pdf",
					Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_FILE),
					Document: content.Document{Name: "file.pdf"},
				}

				rootChildResource := engine.Resource{
					ID:       "1$2!5",
					ParentID: rootResource.ID,
					RootID:   rootResource.ID,
					Path:     "./doc.pdf",
					Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_FILE),
					Document: content.Document{Name: "doc.pdf"},
				}

				rootChildResource2 := engine.Resource{
					ID:       "1$2!6",
					ParentID: rootResource.ID,
					RootID:   rootResource.ID,
					Path:     "./file.pdf",
					Type:     uint64(sprovider.ResourceType_RESOURCE_TYPE_FILE),
					Document: content.Document{Name: "file.pdf"},
				}
				err := eng.Upsert(parentResource.ID, parentResource)
				Expect(err).ToNot(HaveOccurred())

				err = eng.Upsert(rootChildResource.ID, rootChildResource)
				Expect(err).ToNot(HaveOccurred())
				err = eng.Upsert(rootChildResource2.ID, rootChildResource2)
				Expect(err).ToNot(HaveOccurred())

				err = eng.Upsert(childResource.ID, childResource)
				Expect(err).ToNot(HaveOccurred())
				err = eng.Upsert(childResource2.ID, childResource2)
				Expect(err).ToNot(HaveOccurred())
			})
			It("search *doc* in a root", func() {
				res, err := doSearch(rootResource.ID, "Name:*doc*", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(3)))
			})
			It("search *doc* in a subfolder", func() {
				res, err := doSearch(rootResource.ID, "Name:*doc*", "./doc")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(2)))
			})
			It("search *file* in a root", func() {
				res, err := doSearch(rootResource.ID, "Name:*file*", "")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(2)))
			})
			It("search *file* in a subfolder", func() {
				res, err := doSearch(rootResource.ID, "Name:*file*", "./doc")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(1)))
			})
		})
	})

	Describe("Upsert", func() {
		It("adds a resourceInfo to the index", func() {
			err := eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			count, err := idx.DocCount()
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(uint64(1)))

			query := bleve.NewMatchQuery("child.pdf")
			res, err := idx.Search(bleve.NewSearchRequest(query))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Hits.Len()).To(Equal(1))
		})

		It("updates an existing resource in the index", func() {

			err := eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			countA, err := idx.DocCount()
			Expect(err).ToNot(HaveOccurred())
			Expect(countA).To(Equal(uint64(1)))

			err = eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			countB, err := idx.DocCount()
			Expect(err).ToNot(HaveOccurred())
			Expect(countB).To(Equal(uint64(1)))
		})
	})

	Describe("Delete", func() {
		It("marks a resource as deleted", func() {
			err := eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, "Name:*child*", 1)

			err = eng.Delete(childResource.ID)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, "Name:*child*", 0)
		})

		It("marks a child resources as deleted", func() {
			err := eng.Upsert(parentResource.ID, parentResource)
			Expect(err).ToNot(HaveOccurred())

			err = eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, parentResource.Document.Name, 1)
			assertDocCount(rootResource.ID, childResource.Document.Name, 1)

			err = eng.Delete(parentResource.ID)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, parentResource.Document.Name, 0)
			assertDocCount(rootResource.ID, childResource.Document.Name, 0)
		})
	})

	Describe("Restore", func() {
		It("also marks child resources as restored", func() {
			err := eng.Upsert(parentResource.ID, parentResource)
			Expect(err).ToNot(HaveOccurred())

			err = eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			err = eng.Delete(parentResource.ID)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, parentResource.Name, 0)
			assertDocCount(rootResource.ID, childResource.Name, 0)

			err = eng.Restore(parentResource.ID)
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, parentResource.Name, 1)
			assertDocCount(rootResource.ID, childResource.Name, 1)
		})
	})

	Describe("Move", func() {
		It("renames the parent and its child resources", func() {
			err := eng.Upsert(parentResource.ID, parentResource)
			Expect(err).ToNot(HaveOccurred())

			err = eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			parentResource.Path = "newname"
			err = eng.Move(parentResource.ID, parentResource.ParentID, "./my/newname")
			Expect(err).ToNot(HaveOccurred())

			assertDocCount(rootResource.ID, parentResource.Name, 0)

			matches := assertDocCount(rootResource.ID, "Name:child.pdf", 1)
			Expect(matches[0].Entity.ParentId.OpaqueId).To(Equal("3"))
			Expect(matches[0].Entity.Ref.Path).To(Equal("./my/newname/child.pdf"))
		})

		It("moves the parent and its child resources", func() {
			err := eng.Upsert(parentResource.ID, parentResource)
			Expect(err).ToNot(HaveOccurred())

			err = eng.Upsert(childResource.ID, childResource)
			Expect(err).ToNot(HaveOccurred())

			parentResource.Path = " "
			parentResource.ParentID = "1$2!somewhereopaqueid"

			err = eng.Move(parentResource.ID, parentResource.ParentID, "./somewhere/else/newname")
			Expect(err).ToNot(HaveOccurred())
			assertDocCount(rootResource.ID, `parent d!r`, 0)

			matches := assertDocCount(rootResource.ID, "Name:child.pdf", 1)
			Expect(matches[0].Entity.ParentId.OpaqueId).To(Equal("3"))
			Expect(matches[0].Entity.Ref.Path).To(Equal("./somewhere/else/newname/child.pdf"))

			matches = assertDocCount(rootResource.ID, `newname`, 1)
			Expect(matches[0].Entity.ParentId.OpaqueId).To(Equal("somewhereopaqueid"))
			Expect(matches[0].Entity.Ref.Path).To(Equal("./somewhere/else/newname"))

		})
	})
})

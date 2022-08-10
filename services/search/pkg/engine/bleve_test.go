package engine_test

import (
	"context"
	"fmt"
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

		createEntity = func(id sprovider.ResourceId, doc content.Document) engine.Resource {
			name := doc.Name

			if name == "" {
				name = "default.pdf"
			}

			return engine.Resource{
				ID:       fmt.Sprintf("%s$%s!%s", id.StorageId, id.SpaceId, id.OpaqueId),
				RootID:   fmt.Sprintf("%s$%s!%s", id.StorageId, id.SpaceId, id.OpaqueId),
				Path:     fmt.Sprintf("./%s", name),
				Document: doc,
			}
		}

		doSearch = func(id sprovider.ResourceId, query string) (*searchsvc.SearchIndexResponse, error) {
			return eng.Search(ctx, &searchsvc.SearchIndexRequest{
				Query: query,
				Ref: &searchmsg.Reference{
					ResourceId: &searchmsg.ResourceID{
						StorageId: id.StorageId,
						SpaceId:   id.SpaceId,
						OpaqueId:  id.OpaqueId,
					},
				},
			})
		}

		assertDocCount = func(id sprovider.ResourceId, query string, expectedCount int) []*searchmsg.Match {
			res, err := doSearch(id, query)

			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, len(res.Matches)).To(Equal(expectedCount), "query returned unexpected number of results: "+query)
			return res.Matches
		}
	)

	BeforeEach(func() {
		mapping, err := engine.BuildBleveMapping()
		Expect(err).ToNot(HaveOccurred())

		idx, err = bleve.NewMemOnly(mapping)
		Expect(err).ToNot(HaveOccurred())

		eng = engine.NewBleveEngine(idx)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("New", func() {
		It("returns a new index instance", func() {
			b := engine.NewBleveEngine(idx)
			Expect(b).ToNot(BeNil())
		})
	})

	Describe("Search", func() {
		Context("by other fields than filename", func() {
			It("finds files by size", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Size: 12345})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rid, `Size:12345`, 1)
				assertDocCount(rid, `Size:>1000`, 1)
				assertDocCount(rid, `Size:<100000`, 1)

				assertDocCount(rid, `Size:12344`, 0)
				assertDocCount(rid, `Size:<1000`, 0)
				assertDocCount(rid, `Size:>100000`, 0)
			})

			It("finds files by tags", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Tags: []string{"foo", "bar"}})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rid, `Tags:foo`, 1)
				assertDocCount(rid, `Tags:bar`, 1)
				assertDocCount(rid, `Tags:foo Tags:bar`, 1)
				assertDocCount(rid, `Tags:foo Tags:bar Tags:baz`, 1)
				assertDocCount(rid, `+Tags:foo +Tags:bar Tags:baz`, 1)
				assertDocCount(rid, `+Tags:foo +Tags:bar +Tags:baz`, 0)
				assertDocCount(rid, `Tags:baz`, 0)
			})
		})

		Context("by filename", func() {
			It("finds files with spaces in the filename", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "Foo oo.pdf"})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())
				assertDocCount(rid, `Name:foo\ o*`, 1)
			})

			It("finds files by digits in the filename", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "12345.pdf"})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rid, `Name:1234*`, 1)
			})

			Context("with a file in the root of the space", func() {
				It("scopes the search to the specified space", func() {
					rid := sprovider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					}
					r := createEntity(rid, content.Document{Name: "foo.pdf"})
					err := eng.Upsert(r.ID, r)
					Expect(err).ToNot(HaveOccurred())

					assertDocCount(rid, `Name:foo.pdf`, 1)
					assertDocCount(sprovider.ResourceId{
						StorageId: "9",
						SpaceId:   "8",
						OpaqueId:  "7",
					}, `Name:foo.pdf`, 0)
				})
			})

			It("limits the search to the specified fields", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "bar.pdf", Size: 789})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rid, `Name:bar.pdf`, 1)
				assertDocCount(rid, `Size:789`, 1)
				assertDocCount(rid, `Unknown:field`, 0)
			})

			It("returns the total number of hits", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "bar.pdf"})
				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				res, err := doSearch(rid, "Name:bar*")
				Expect(err).ToNot(HaveOccurred())
				Expect(res.TotalMatches).To(Equal(int32(1)))
			})

			It("returns all desired fields", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "bar.pdf"})
				r.Type = 3
				r.MimeType = "application/pdf"

				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				matches := assertDocCount(rid, fmt.Sprintf("Name:%s", r.Name), 1)
				match := matches[0]
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(rid.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal(r.Path))
				Expect(match.Entity.Name).To(Equal(r.Name))
				Expect(match.Entity.Size).To(Equal(r.Size))
				Expect(match.Entity.Type).To(Equal(r.Type))
				Expect(match.Entity.MimeType).To(Equal(r.MimeType))
				Expect(match.Entity.Deleted).To(BeFalse())
				Expect(match.Score > 0).To(BeTrue())
			})

			It("finds files by name, prefix or substring match", func() {
				queries := []string{"foo.pdf", "foo*", "*oo.p*"}
				for i, query := range queries {
					rid := sprovider.ResourceId{
						StorageId: string(rune(i + 1)),
						SpaceId:   string(rune(i + 2)),
						OpaqueId:  string(rune(i + 3)),
					}

					r := createEntity(rid, content.Document{Name: "foo.pdf"})
					r.Size = uint64(i + 250)
					err := eng.Upsert(r.ID, r)
					Expect(err).ToNot(HaveOccurred())

					matches := assertDocCount(rid, query, 1)
					Expect(matches[0].Entity.Ref.ResourceId.OpaqueId).To(Equal(rid.OpaqueId))
					Expect(matches[0].Entity.Ref.Path).To(Equal(r.Path))
					Expect(matches[0].Entity.Id.OpaqueId).To(Equal(rid.OpaqueId))
					Expect(matches[0].Entity.Name).To(Equal(r.Name))
					Expect(matches[0].Entity.Size).To(Equal(r.Size))
				}
			})

			It("uses a lower-case index", func() {
				rid := sprovider.ResourceId{
					StorageId: "1",
					SpaceId:   "2",
					OpaqueId:  "3",
				}
				r := createEntity(rid, content.Document{Name: "foo.pdf"})
				r.Type = 3
				r.MimeType = "application/pdf"

				err := eng.Upsert(r.ID, r)
				Expect(err).ToNot(HaveOccurred())

				assertDocCount(rid, "Name:foo*", 1)
				assertDocCount(rid, "Name:Foo*", 0)
			})

			Context("and an additional file in a subdirectory", func() {
				var (
					ridT sprovider.ResourceId
					entT engine.Resource
					ridD sprovider.ResourceId
					entD engine.Resource
				)

				BeforeEach(func() {
					ridT = sprovider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "3",
					}
					entT = createEntity(ridT, content.Document{Name: "top.pdf"})
					Expect(eng.Upsert(entT.ID, entT)).ToNot(HaveOccurred())

					ridD = sprovider.ResourceId{
						StorageId: "1",
						SpaceId:   "2",
						OpaqueId:  "4",
					}
					entD = createEntity(ridD, content.Document{Name: "deep.pdf"})
					entD.Path = "./nested/deep.pdf"
					Expect(eng.Upsert(entD.ID, entD)).ToNot(HaveOccurred())
				})

				It("finds files living deeper in the tree by filename, prefix or substring match", func() {
					queries := []string{"deep.pdf", "dee*", "*ep.*"}
					for _, query := range queries {
						assertDocCount(ridD, query, 1)
					}
				})

				It("does not find the higher levels when limiting the searched directory", func() {
					res, err := eng.Search(ctx, &searchsvc.SearchIndexRequest{
						Ref: &searchmsg.Reference{
							ResourceId: &searchmsg.ResourceID{
								StorageId: ridT.StorageId,
								SpaceId:   ridT.SpaceId,
								OpaqueId:  ridT.OpaqueId,
							},
							Path: "./nested/",
						},
						Query: "Name:top.pdf",
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(0))
				})
			})
		})
	})
})

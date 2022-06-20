package provider_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/v2/extensions/search/pkg/search/mocks"
	provider "github.com/owncloud/ocis/v2/extensions/search/pkg/search/provider"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
)

var _ = Describe("Searchprovider", func() {
	var (
		p           *provider.Provider
		gwClient    *cs3mocks.GatewayAPIClient
		indexClient *mocks.IndexClient

		ctx        context.Context
		eventsChan chan interface{}

		logger = log.NewLogger()
		user   = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
		otherUser = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "otheruser",
			},
		}
		personalSpace = &sprovider.StorageSpace{
			Opaque: &typesv1beta1.Opaque{
				Map: map[string]*typesv1beta1.OpaqueEntry{
					"path": {
						Decoder: "plain",
						Value:   []byte("/foo"),
					},
				},
			},
			Id:   &sprovider.StorageSpaceId{OpaqueId: "personalspace"},
			Root: &sprovider.ResourceId{StorageId: "storageid", OpaqueId: "storageid"},
			Name: "personalspace",
		}

		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				OpaqueId:  "opaqueid",
			},
			Path: "foo.pdf",
			Size: 12345,
		}
	)

	BeforeEach(func() {
		ctx = context.Background()
		eventsChan = make(chan interface{})
		gwClient = &cs3mocks.GatewayAPIClient{}
		indexClient = &mocks.IndexClient{}

		p = provider.New(gwClient, indexClient, "", eventsChan, logger)

		gwClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
			Status: status.NewOK(ctx),
			Token:  "authtoken",
		}, nil)
		gwClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
			Status: status.NewOK(context.Background()),
			Info:   ri,
		}, nil)
		indexClient.On("DocCount").Return(uint64(1), nil)
	})

	Describe("New", func() {
		It("returns a new instance", func() {
			p := provider.New(gwClient, indexClient, "", eventsChan, logger)
			Expect(p).ToNot(BeNil())
		})
	})

	Describe("IndexSpace", func() {
		It("walks the space and indexes all files", func() {
			gwClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil)

			res, err := p.IndexSpace(ctx, &searchsvc.IndexSpaceRequest{
				SpaceId: "storageid",
				UserId:  "user",
			})
			Expect(err).ToNot(HaveOccurred())
			Expect(res).ToNot(BeNil())
		})
	})

	Describe("Search", func() {
		It("fails when an empty query is given", func() {
			res, err := p.Search(ctx, &searchsvc.SearchRequest{
				Query: "",
			})
			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())
		})

		Context("with a personal space", func() {
			BeforeEach(func() {
				gwClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
					Matches: []*searchmsg.Match{
						{
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										OpaqueId:  personalSpace.Root.OpaqueId,
									},
									Path: "./path/to/Foo.pdf",
								},
								Id: &searchmsg.ResourceID{
									StorageId: personalSpace.Root.StorageId,
									OpaqueId:  "foo-id",
								},
								Name: "Foo.pdf",
							},
						},
					},
				}, nil)
			})

			It("lowercases the filename", func() {
				p.Search(ctx, &searchsvc.SearchRequest{
					Query: "Foo.pdf",
				})
				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "Name:*foo.pdf*"
				}))
			})

			It("does not mess with field-based searches", func() {
				p.Search(ctx, &searchsvc.SearchRequest{
					Query: "Size:<10",
				})
				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "Size:<10"
				}))
			})

			It("uppercases field names", func() {
				tests := []struct {
					Original string
					Expected string
				}{
					{Original: "size:<100", Expected: "Size:<100"},
				}
				for _, test := range tests {
					p.Search(ctx, &searchsvc.SearchRequest{
						Query: test.Original,
					})
					indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
						return req.Query == test.Expected
					}))
				}
			})

			It("escapes special characters", func() {
				p.Search(ctx, &searchsvc.SearchRequest{
					Query: "Foo oo.pdf",
				})
				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == `Name:*foo\ oo.pdf*`
				}))
			})

			It("searches the personal user space", func() {
				res, err := p.Search(ctx, &searchsvc.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Id.OpaqueId).To(Equal("foo-id"))
				Expect(match.Entity.Name).To(Equal("Foo.pdf"))
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(personalSpace.Root.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal("./path/to/Foo.pdf"))

				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "Name:*foo*" && req.Ref.ResourceId.OpaqueId == personalSpace.Root.OpaqueId && req.Ref.Path == ""
				}))
			})
		})

		Context("with received shares", func() {
			var (
				grantSpace      *sprovider.StorageSpace
				mountpointSpace *sprovider.StorageSpace
			)

			BeforeEach(func() {
				grantSpace = &sprovider.StorageSpace{
					SpaceType: "grant",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "otherspaceroot!otherspacegrant"},
					Root:      &sprovider.ResourceId{StorageId: "otherspaceroot", OpaqueId: "otherspacegrant"},
					Name:      "grantspace",
				}
				mountpointSpace = &sprovider.StorageSpace{
					SpaceType: "mountpoint",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "otherspaceroot!otherspacemountpoint"},
					Root:      &sprovider.ResourceId{StorageId: "otherspaceroot", OpaqueId: "otherspacemountpoint"},
					Name:      "mountpointspace",
					Opaque: &typesv1beta1.Opaque{
						Map: map[string]*typesv1beta1.OpaqueEntry{
							"grantStorageID": {Decoder: "plain", Value: []byte("otherspaceroot")},
							"grantOpaqueID":  {Decoder: "plain", Value: []byte("otherspacegrant")},
						},
					},
				}
				gwClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
					Status: status.NewOK(ctx),
					Path:   "/grant/path",
				}, nil)
			})

			It("searches the received spaces", func() {
				gwClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{grantSpace, mountpointSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
					Matches: []*searchmsg.Match{
						{
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: grantSpace.Root.StorageId,
										OpaqueId:  grantSpace.Root.OpaqueId,
									},
									Path: "./grant/path/to/Shared.pdf",
								},
								Id: &searchmsg.ResourceID{
									StorageId: grantSpace.Root.StorageId,
									OpaqueId:  "grant-shared-id",
								},
								Name: "Shared.pdf",
							},
						},
					},
				}, nil)

				res, err := p.Search(ctx, &searchsvc.SearchRequest{
					Query: "Foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Id.OpaqueId).To(Equal("grant-shared-id"))
				Expect(match.Entity.Name).To(Equal("Shared.pdf"))
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(mountpointSpace.Root.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal("./to/Shared.pdf"))

				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "Name:*foo*" && req.Ref.ResourceId.StorageId == grantSpace.Root.StorageId && req.Ref.Path == "./grant/path"
				}))
			})

			It("finds matches in both the personal space AND the grant", func() {
				gwClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, grantSpace, mountpointSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.StorageId == grantSpace.Root.StorageId
				})).Return(&searchsvc.SearchIndexResponse{
					Matches: []*searchmsg.Match{
						{
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: grantSpace.Root.StorageId,
										OpaqueId:  grantSpace.Root.OpaqueId,
									},
									Path: "./grant/path/to/Shared.pdf",
								},
								Id: &searchmsg.ResourceID{
									StorageId: grantSpace.Root.StorageId,
									OpaqueId:  "grant-shared-id",
								},
								Name: "Shared.pdf",
							},
						},
					},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.StorageId == personalSpace.Root.StorageId
				})).Return(&searchsvc.SearchIndexResponse{
					Matches: []*searchmsg.Match{
						{
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										OpaqueId:  personalSpace.Root.OpaqueId,
									},
									Path: "./path/to/Foo.pdf",
								},
								Id: &searchmsg.ResourceID{
									StorageId: personalSpace.Root.StorageId,
									OpaqueId:  "foo-id",
								},
								Name: "Foo.pdf",
							},
						},
					},
				}, nil)

				res, err := p.Search(ctx, &searchsvc.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(2))
				ids := []string{res.Matches[0].Entity.Id.OpaqueId, res.Matches[1].Entity.Id.OpaqueId}
				Expect(ids).To(ConsistOf("foo-id", "grant-shared-id"))
			})
		})
	})
})

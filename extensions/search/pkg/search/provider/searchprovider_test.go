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
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/extensions/search/pkg/search/mocks"
	provider "github.com/owncloud/ocis/extensions/search/pkg/search/provider"
	"github.com/owncloud/ocis/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
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
			Root: &sprovider.ResourceId{OpaqueId: "personalspaceroot"},
			Name: "personalspace",
		}

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

	Describe("events", func() {
		It("trigger an index update when a file has been uploaded", func() {
			called := false
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.FileUploaded{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}).Should(BeTrue())
		})

		It("removes an entry from the index when the file has been deleted", func() {
			called := false

			gwClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
				Status: status.NewNotFound(context.Background(), ""),
			}, nil)
			indexClient.On("Delete", mock.MatchedBy(func(id *sprovider.ResourceId) bool {
				return id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.ItemTrashed{
				Ref:       ref,
				ID:        ri.Id,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}).Should(BeTrue())
		})

		It("indexes items when they are being restored", func() {
			called := false
			indexClient.On("Restore", mock.MatchedBy(func(id *sprovider.ResourceId) bool {
				return id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.ItemRestored{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}).Should(BeTrue())
		})

		It("indexes items when a version has been restored", func() {
			called := false
			indexClient.On("Add", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.FileVersionRestored{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}).Should(BeTrue())
		})

		It("indexes items when they are being moved", func() {
			called := false
			indexClient.On("Move", mock.Anything, mock.MatchedBy(func(riToIndex *sprovider.ResourceInfo) bool {
				return riToIndex.Id.OpaqueId == ri.Id.OpaqueId
			})).Return(nil).Run(func(args mock.Arguments) {
				called = true
			})
			eventsChan <- events.ItemMoved{
				Ref:       ref,
				Executant: user.Id,
			}

			Eventually(func() bool {
				return called
			}).Should(BeTrue())
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
				gwClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					p := string(req.Opaque.Map["path"].Value)
					return p == "/"
				})).Return(&sprovider.ListStorageSpacesResponse{
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
					return req.Query == "foo" && req.Ref.ResourceId.OpaqueId == personalSpace.Root.OpaqueId && req.Ref.Path == ""
				}))
			})
		})

		Context("with received shares", func() {
			var (
				grantSpace *sprovider.StorageSpace
			)

			BeforeEach(func() {
				grantSpace = &sprovider.StorageSpace{
					SpaceType: "grant",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "otherspaceroot!otherspacegrant"},
					Root:      &sprovider.ResourceId{StorageId: "otherspaceroot", OpaqueId: "otherspacegrant"},
					Name:      "grantspace",
				}
				gwClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
					Status: status.NewOK(ctx),
					Path:   "/grant/path",
				}, nil)
			})

			It("searches the received spaces (grants)", func() {
				gwClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					p := string(req.Opaque.Map["path"].Value)
					return p == "/"
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{grantSpace},
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
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Id.OpaqueId).To(Equal("grant-shared-id"))
				Expect(match.Entity.Name).To(Equal("Shared.pdf"))
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(grantSpace.Root.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal("./to/Shared.pdf"))

				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "foo" && req.Ref.ResourceId.OpaqueId == grantSpace.Root.OpaqueId && req.Ref.Path == "./grant/path"
				}))
			})

			It("finds matches in both the personal space AND the grant", func() {
				gwClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					p := string(req.Opaque.Map["path"].Value)
					return p == "/"
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, grantSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.OpaqueId == grantSpace.Root.OpaqueId
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
					return req.Ref.ResourceId.OpaqueId == personalSpace.Root.OpaqueId
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

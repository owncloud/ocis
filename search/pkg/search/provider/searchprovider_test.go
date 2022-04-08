// Copyright 2018-2022 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package provider_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/owncloud/ocis/search/pkg/search"
	"github.com/owncloud/ocis/search/pkg/search/mocks"
	provider "github.com/owncloud/ocis/search/pkg/search/provider"
)

var _ = Describe("Searchprovider", func() {
	var (
		p           *provider.Provider
		gwClient    *cs3mocks.GatewayAPIClient
		indexClient *mocks.IndexClient

		ctx context.Context

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
	)

	BeforeEach(func() {
		ctx = context.Background()
		gwClient = &cs3mocks.GatewayAPIClient{}
		indexClient = &mocks.IndexClient{}

		p = provider.New(gwClient, indexClient)
	})

	Describe("New", func() {
		It("returns a new instance", func() {
			p := provider.New(gwClient, indexClient)
			Expect(p).ToNot(BeNil())
		})
	})

	Describe("Search", func() {
		It("fails when an empty query is given", func() {
			res, err := p.Search(ctx, &search.SearchRequest{
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
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&search.SearchIndexResult{
					Matches: []search.Match{
						{
							Reference: &sprovider.Reference{
								ResourceId: personalSpace.Root,
								Path:       "./path/to/Foo.pdf",
							},
							Info: &sprovider.ResourceInfo{
								Id: &sprovider.ResourceId{
									StorageId: personalSpace.Root.StorageId,
									OpaqueId:  "foo-id",
								},
								Path: "Foo.pdf",
							},
						},
					},
				}, nil)
			})

			It("searches the personal user space", func() {
				res, err := p.Search(ctx, &search.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Info.Id.OpaqueId).To(Equal("foo-id"))
				Expect(match.Info.Path).To(Equal("Foo.pdf"))
				Expect(match.Reference.ResourceId).To(Equal(personalSpace.Root))
				Expect(match.Reference.Path).To(Equal("./path/to/Foo.pdf"))

				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *search.SearchIndexRequest) bool {
					return req.Query == "foo" && req.Reference.ResourceId == personalSpace.Root && req.Reference.Path == ""
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
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&search.SearchIndexResult{
					Matches: []search.Match{
						search.Match{
							Reference: &sprovider.Reference{
								ResourceId: grantSpace.Root,
								Path:       "./grant/path/to/Foo.pdf",
							},
							Info: &sprovider.ResourceInfo{
								Id: &sprovider.ResourceId{
									StorageId: grantSpace.Root.StorageId,
									OpaqueId:  "grant-foo-id",
								},
								Path: "Foo.pdf",
							},
						},
					},
				}, nil)

				res, err := p.Search(ctx, &search.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Info.Id.OpaqueId).To(Equal("grant-foo-id"))
				Expect(match.Info.Path).To(Equal("Foo.pdf"))
				Expect(match.Reference.ResourceId).To(Equal(grantSpace.Root))
				Expect(match.Reference.Path).To(Equal("./to/Foo.pdf"))

				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *search.SearchIndexRequest) bool {
					return req.Query == "foo" && req.Reference.ResourceId == grantSpace.Root && req.Reference.Path == "./grant/path"
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
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *search.SearchIndexRequest) bool {
					return req.Reference.ResourceId == grantSpace.Root
				})).Return(&search.SearchIndexResult{
					Matches: []search.Match{
						search.Match{
							Reference: &sprovider.Reference{
								ResourceId: grantSpace.Root,
								Path:       "./grant/path/to/Foo.pdf",
							},
							Info: &sprovider.ResourceInfo{
								Id: &sprovider.ResourceId{
									StorageId: grantSpace.Root.StorageId,
									OpaqueId:  "grant-foo-id",
								},
								Path: "Foo.pdf",
							},
						},
					},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *search.SearchIndexRequest) bool {
					return req.Reference.ResourceId == personalSpace.Root
				})).Return(&search.SearchIndexResult{
					Matches: []search.Match{
						search.Match{
							Reference: &sprovider.Reference{
								ResourceId: personalSpace.Root,
								Path:       "./path/to/Foo.pdf",
							},
							Info: &sprovider.ResourceInfo{
								Id: &sprovider.ResourceId{
									StorageId: personalSpace.Root.StorageId,
									OpaqueId:  "foo-id",
								},
								Path: "Foo.pdf",
							},
						},
					},
				}, nil)

				res, err := p.Search(ctx, &search.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(2))
				ids := []string{res.Matches[0].Info.Id.OpaqueId, res.Matches[1].Info.Id.OpaqueId}
				Expect(ids).To(ConsistOf("foo-id", "grant-foo-id"))

			})
		})
	})
})

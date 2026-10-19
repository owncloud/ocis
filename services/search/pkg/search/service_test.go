package search_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaborationv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	sprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/search/pkg/config"
	"github.com/owncloud/ocis/v2/services/search/pkg/content"
	contentMocks "github.com/owncloud/ocis/v2/services/search/pkg/content/mocks"
	"github.com/owncloud/ocis/v2/services/search/pkg/engine"
	engineMocks "github.com/owncloud/ocis/v2/services/search/pkg/engine/mocks"
	"github.com/owncloud/ocis/v2/services/search/pkg/search"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/status"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var _ = Describe("Searchprovider", func() {
	var (
		s               search.Searcher
		extractor       *contentMocks.Extractor
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		indexClient     *engineMocks.Engine
		ctx             context.Context
		logger          = log.NewLogger()
		user            = &userv1beta1.User{
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
			Id:        &sprovider.StorageSpaceId{OpaqueId: "storageid$personalspace!personalspace"},
			Root:      &sprovider.ResourceId{StorageId: "storageid", SpaceId: "personalspace", OpaqueId: "personalspace"},
			Name:      "personalspace",
			SpaceType: "personal",
		}

		ri = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "opaqueid",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "parentopaqueid",
			},
			Path:  "./foo.pdf",
			Size:  12345,
			Mtime: &typesv1beta1.Timestamp{Seconds: 4000},
		}

		ri2 = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "opaqueid",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "parentopaqueid",
			},
			Type:  sprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
			Path:  "./.space",
			Size:  12345,
			Mtime: &typesv1beta1.Timestamp{Seconds: 4000},
		}

		ri2_1 = &sprovider.ResourceInfo{
			Id: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "opaqueid_1",
			},
			ParentId: &sprovider.ResourceId{
				StorageId: "storageid",
				SpaceId:   "spaceid",
				OpaqueId:  "opaqueid",
			},
			Type:  sprovider.ResourceType_RESOURCE_TYPE_FILE,
			Path:  "./.space/file.pdf",
			Size:  12345,
			Mtime: &typesv1beta1.Timestamp{Seconds: 4000},
		}
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		ctx = revactx.ContextSetUser(context.Background(), user)
		indexClient = &engineMocks.Engine{}
		extractor = &contentMocks.Extractor{}

		s = search.NewService(gatewaySelector, indexClient, extractor, logger, &config.Config{})

		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{
			Status: status.NewOK(ctx),
			Token:  "authtoken",
		}, nil)
		gatewayClient.On("GetPath", mock.Anything, mock.MatchedBy(func(req *sprovider.GetPathRequest) bool {
			return req.ResourceId.OpaqueId == ri.Id.OpaqueId
		})).Return(&sprovider.GetPathResponse{
			Status: status.NewOK(context.Background()),
			Path:   ri.Path,
		}, nil)
		gatewayClient.On("GetReceivedShare", mock.Anything, mock.Anything).Return(&collaborationv1beta1.GetReceivedShareResponse{
			Status: status.NewOK(ctx),
		}, nil)
		indexClient.On("DocCount").Return(uint64(1), nil)
		indexClient.On("Optimize", mock.Anything).Return(nil)
	})

	Describe("New", func() {
		It("returns a new instance", func() {
			s := search.NewService(gatewaySelector, indexClient, extractor, logger, &config.Config{})
			Expect(s).ToNot(BeNil())
		})
	})

	Describe("IndexSpace", func() {
		It("walks the space and indexes all files", func() {
			gatewayClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			indexClient.On("Upsert", mock.Anything, mock.MatchedBy(func(r engine.Resource) bool {
				return r.ID == "storageid$spaceid!opaqueid" && r.Path == "./foo.pdf"
			})).Return(nil)
			indexClient.On("Lookup", mock.Anything).Return(nil, engine.ErrResourceNotFound)
			gatewayClient.On("Stat", mock.Anything, mock.MatchedBy(func(sreq *sprovider.StatRequest) bool {
				return sreq.Ref.ResourceId.StorageId == "storageid" &&
					(sreq.Ref.ResourceId.OpaqueId == "spaceid" || sreq.Ref.ResourceId.OpaqueId == ri.Id.OpaqueId) &&
					sreq.Ref.ResourceId.SpaceId == "spaceid"
			})).Return(&sprovider.StatResponse{
				Status: status.NewOK(context.Background()),
				Info:   ri,
			}, nil)
			err := s.IndexSpace(context.Background(), &sprovider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"})
			Expect(err).ShouldNot(HaveOccurred())
			indexClient.AssertNumberOfCalls(GinkgoT(), "Upsert", 1)
		})

		It("walks the space skipping .space directory", func() {
			gatewayClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&sprovider.ListContainerResponse{
				Status: status.NewOK(context.Background()),
				Infos:  []*sprovider.ResourceInfo{ri2_1},
			}, nil)
			gatewayClient.On("Stat", mock.Anything, mock.MatchedBy(func(sreq *sprovider.StatRequest) bool {
				return sreq.Ref.ResourceId.StorageId == "storageid" &&
					sreq.Ref.ResourceId.OpaqueId == "spaceid" &&
					sreq.Ref.ResourceId.SpaceId == "spaceid"
			})).Return(&sprovider.StatResponse{
				Status: status.NewOK(context.Background()),
				Info:   ri2,
			}, nil)

			indexClient.AssertNotCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
			indexClient.AssertNotCalled(GinkgoT(), "Search", mock.Anything, mock.Anything)
			err := s.IndexSpace(context.Background(), &sprovider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"})
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("descends into already-indexed directories to index unvisited children", func() {
			// Simulate a large tree where all directories are already indexed
			// (mtime matches) but each contains one unindexed file — e.g. after
			// a Tika crash interrupted a previous reindex run.  The fix ensures
			// we never SkipDir for directories, so children are always visited.
			//
			// 10,000 directories: pessimistic-case benchmark confirming the walker
			// scales linearly even when every directory is already indexed.
			const numDirs = 10000

			gatewayClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			// The indexed mtime that the engine returns for directories.
			indexedMtime := time.Unix(5000, 0).UTC().Format(time.RFC3339Nano)

			// Build root: a container that lists numDirs child directories.
			rootInfo := &sprovider.ResourceInfo{
				Id:    &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "spaceid"},
				Type:  sprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Path:  ".",
				Mtime: &typesv1beta1.Timestamp{Seconds: 5000},
			}
			dirInfos := make([]*sprovider.ResourceInfo, numDirs)
			for i := 0; i < numDirs; i++ {
				dirInfos[i] = &sprovider.ResourceInfo{
					Id:    &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: fmt.Sprintf("dir-%d", i)},
					Type:  sprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
					Path:  fmt.Sprintf("dir-%d", i),
					Mtime: &typesv1beta1.Timestamp{Seconds: 5000},
				}
			}

			// Pre-build file infos for each directory (map for O(1) dispatch).
			fileInfoByDir := make(map[string]*sprovider.ResourceInfo, numDirs)
			for i := 0; i < numDirs; i++ {
				dirID := fmt.Sprintf("dir-%d", i)
				fileInfoByDir[dirID] = &sprovider.ResourceInfo{
					Id:       &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: fmt.Sprintf("file-%d", i)},
					ParentId: &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: dirID},
					Type:     sprovider.ResourceType_RESOURCE_TYPE_FILE,
					Path:     fmt.Sprintf("file-%d.pdf", i),
					Size:     1024,
					Mtime:    &typesv1beta1.Timestamp{Seconds: 5000},
				}
			}

			// Build ID-to-fileInfo and ID-to-path maps for O(1) Stat and GetPath dispatch.
			fileInfoByID := make(map[string]*sprovider.ResourceInfo, numDirs)
			filePathByID := make(map[string]string, numDirs)
			for dirID, fi := range fileInfoByDir {
				fileInfoByID[fi.Id.OpaqueId] = fi
				filePathByID[fi.Id.OpaqueId] = "/" + dirID + "/" + fi.Path
			}
			gatewayClient.On("GetPath", mock.Anything, mock.MatchedBy(func(req *sprovider.GetPathRequest) bool {
				return fileInfoByID[req.GetResourceId().GetOpaqueId()] != nil
			})).Return(
				func(_ context.Context, req *sprovider.GetPathRequest, _ ...grpc.CallOption) *sprovider.GetPathResponse {
					return &sprovider.GetPathResponse{Status: status.NewOK(context.Background()), Path: filePathByID[req.GetResourceId().GetOpaqueId()]}
				},
				func(_ context.Context, _ *sprovider.GetPathRequest, _ ...grpc.CallOption) error { return nil },
			)

			// Single Stat mock dispatching via map — avoids 10k separate matchers.
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(
				func(_ context.Context, sreq *sprovider.StatRequest, _ ...grpc.CallOption) *sprovider.StatResponse {
					if sreq.Ref.ResourceId != nil && sreq.Ref.ResourceId.OpaqueId == "spaceid" {
						return &sprovider.StatResponse{
							Status: status.NewOK(context.Background()),
							Info:   rootInfo,
						}
					}
					if fi, ok := fileInfoByID[sreq.Ref.GetResourceId().GetOpaqueId()]; ok {
						return &sprovider.StatResponse{
							Status: status.NewOK(context.Background()),
							Info:   fi,
						}
					}
					return &sprovider.StatResponse{Status: status.NewOK(context.Background()), Info: rootInfo}
				},
				func(_ context.Context, _ *sprovider.StatRequest, _ ...grpc.CallOption) error {
					return nil
				},
			)

			// Single ListContainer mock dispatching via map.
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(
				func(_ context.Context, req *sprovider.ListContainerRequest, _ ...grpc.CallOption) *sprovider.ListContainerResponse {
					opaqueID := req.Ref.ResourceId.OpaqueId
					if opaqueID == "spaceid" {
						return &sprovider.ListContainerResponse{
							Status: status.NewOK(context.Background()),
							Infos:  dirInfos,
						}
					}
					if fi, ok := fileInfoByDir[opaqueID]; ok {
						return &sprovider.ListContainerResponse{
							Status: status.NewOK(context.Background()),
							Infos:  []*sprovider.ResourceInfo{fi},
						}
					}
					return &sprovider.ListContainerResponse{
						Status: status.NewOK(context.Background()),
						Infos:  nil,
					}
				},
				func(_ context.Context, _ *sprovider.ListContainerRequest, _ ...grpc.CallOption) error {
					return nil
				},
			)

			// Lookup: root and directories are already indexed with matching mtime.
			// Use a map for O(1) lookup instead of linear scan (matters at 10k dirs).
			indexedIDs := make(map[string]bool, numDirs+1)
			indexedIDs["storageid$spaceid!spaceid"] = true
			for i := 0; i < numDirs; i++ {
				indexedIDs[fmt.Sprintf("storageid$spaceid!dir-%d", i)] = true
			}
			indexedResource := &engine.Resource{Document: content.Document{Mtime: indexedMtime}, Extracted: true}
			indexClient.On("Lookup", mock.MatchedBy(func(id string) bool {
				return indexedIDs[id] || len(id) > 18 // dir-N or file-N IDs
			})).Return(func(id string) (*engine.Resource, error) {
				if indexedIDs[id] {
					return indexedResource, nil
				}
				// Files are not indexed.
				return nil, engine.ErrResourceNotFound
			})

			// Track Upsert calls — each unindexed file should trigger one.
			var upsertCount atomic.Int32
			indexClient.On("Upsert", mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
				upsertCount.Add(1)
			}).Return(nil)

			err := s.IndexSpace(context.Background(), &sprovider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(int(upsertCount.Load())).To(Equal(numDirs), "every unindexed file inside an already-indexed directory must be visited")
		})

		It("skips files still in postprocessing", func() {
			gatewayClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)

			processingFile := &sprovider.ResourceInfo{
				Id: &sprovider.ResourceId{
					StorageId: "storageid",
					SpaceId:   "spaceid",
					OpaqueId:  "opaqueid",
				},
				ParentId: &sprovider.ResourceId{
					StorageId: "storageid",
					SpaceId:   "spaceid",
					OpaqueId:  "parentopaqueid",
				},
				Type:   sprovider.ResourceType_RESOURCE_TYPE_FILE,
				Path:   "./inflight.pdf",
				Size:   1024,
				Mtime:  &typesv1beta1.Timestamp{Seconds: 4000},
				Opaque: utils.AppendPlainToOpaque(nil, "status", "processing"),
			}

			gatewayClient.On("Stat", mock.Anything, mock.MatchedBy(func(sreq *sprovider.StatRequest) bool {
				return sreq.Ref.ResourceId.StorageId == "storageid" &&
					sreq.Ref.ResourceId.OpaqueId == "spaceid" &&
					sreq.Ref.ResourceId.SpaceId == "spaceid"
			})).Return(&sprovider.StatResponse{
				Status: status.NewOK(context.Background()),
				Info:   processingFile,
			}, nil)

			err := s.IndexSpace(context.Background(), &sprovider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"})
			Expect(err).ShouldNot(HaveOccurred())

			// The skip gate fires before any index lookup or extraction.
			indexClient.AssertNotCalled(GinkgoT(), "Lookup", mock.Anything)
			indexClient.AssertNotCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
			extractor.AssertNotCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
		})
	})

	Describe("IndexSpace abort threshold", func() {
		// A root container listing numFiles files. Stat by ID answers per the
		// missing/statBroken sets, Stat by path per staleByPath; Lookup finds only
		// the files in indexed, so every other file is attempted.
		const numFiles = 6
		var (
			rootInfo    *sprovider.ResourceInfo
			fileInfos   []*sprovider.ResourceInfo
			missing     map[string]bool // Stat answers NOT_FOUND
			statBroken  map[string]bool // Stat answers with an internal error
			staleByPath map[string]bool // Stat by path answers NOT_FOUND, Stat by ID works
			pathMissing map[string]bool // GetPath answers NOT_FOUND
			pathBroken  map[string]bool // GetPath answers with an internal error
			indexed     map[string]bool // already indexed and unchanged
			newService  = func(maxConsecutiveFailures int) *search.Service {
				cfg := &config.Config{}
				cfg.Extractor.MaxConsecutiveFailures = maxConsecutiveFailures
				s := search.NewService(gatewaySelector, indexClient, extractor, logger, cfg)
				s.SetExtractionRetryDelay(time.Millisecond)
				return s
			}
			spaceID = &sprovider.StorageSpaceId{OpaqueId: "storageid$spaceid!spaceid"}
		)

		BeforeEach(func() {
			missing = map[string]bool{}
			statBroken = map[string]bool{}
			staleByPath = map[string]bool{}
			pathMissing = map[string]bool{}
			pathBroken = map[string]bool{}
			indexed = map[string]bool{}
			gatewayClient.On("GetUserByClaim", mock.Anything, mock.Anything).Return(&userv1beta1.GetUserByClaimResponse{
				Status: status.NewOK(context.Background()),
				User:   user,
			}, nil)

			rootInfo = &sprovider.ResourceInfo{
				Id:    &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "spaceid"},
				Type:  sprovider.ResourceType_RESOURCE_TYPE_CONTAINER,
				Path:  ".",
				Mtime: &typesv1beta1.Timestamp{Seconds: 5000},
			}
			fileInfos = make([]*sprovider.ResourceInfo, numFiles)
			byID := map[string]*sprovider.ResourceInfo{}
			for i := range fileInfos {
				fileInfos[i] = &sprovider.ResourceInfo{
					Id:       &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: fmt.Sprintf("file-%d", i)},
					ParentId: &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "spaceid"},
					Type:     sprovider.ResourceType_RESOURCE_TYPE_FILE,
					Path:     fmt.Sprintf("file-%d.pdf", i),
					Size:     1024,
					Mtime:    &typesv1beta1.Timestamp{Seconds: 5000},
				}
				byID[fileInfos[i].Id.OpaqueId] = fileInfos[i]
			}

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(
				func(_ context.Context, sreq *sprovider.StatRequest, _ ...grpc.CallOption) *sprovider.StatResponse {
					id := sreq.GetRef().GetResourceId().GetOpaqueId()
					if p := sreq.GetRef().GetPath(); p != "" && p != "." {
						// By path, relative to the space root.
						for _, fi := range fileInfos {
							if "./"+fi.Path == p && !staleByPath[fi.Id.OpaqueId] {
								id = fi.Id.OpaqueId
								break
							}
							id = ""
						}
					}
					switch {
					case id == "spaceid":
						return &sprovider.StatResponse{Status: status.NewOK(context.Background()), Info: rootInfo}
					case statBroken[id]:
						return &sprovider.StatResponse{Status: status.NewInternal(context.Background(), "storage unavailable")}
					case missing[id] || byID[id] == nil:
						return &sprovider.StatResponse{Status: status.NewNotFound(context.Background(), id)}
					}
					return &sprovider.StatResponse{Status: status.NewOK(context.Background()), Info: byID[id]}
				},
				func(_ context.Context, _ *sprovider.StatRequest, _ ...grpc.CallOption) error { return nil },
			)
			gatewayClient.On("GetPath", mock.Anything, mock.MatchedBy(func(req *sprovider.GetPathRequest) bool {
				return byID[req.GetResourceId().GetOpaqueId()] != nil
			})).Return(
				func(_ context.Context, req *sprovider.GetPathRequest, _ ...grpc.CallOption) *sprovider.GetPathResponse {
					id := req.GetResourceId().GetOpaqueId()
					switch {
					case pathMissing[id]:
						return &sprovider.GetPathResponse{Status: status.NewNotFound(context.Background(), id)}
					case pathBroken[id]:
						return &sprovider.GetPathResponse{Status: status.NewInternal(context.Background(), "storage unavailable")}
					}
					return &sprovider.GetPathResponse{Status: status.NewOK(context.Background()), Path: "/" + byID[id].Path}
				},
				func(_ context.Context, _ *sprovider.GetPathRequest, _ ...grpc.CallOption) error { return nil },
			)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&sprovider.ListContainerResponse{
				Status: status.NewOK(context.Background()),
				Infos:  fileInfos,
			}, nil)
			indexClient.On("Lookup", mock.Anything).Return(
				func(id string) *engine.Resource {
					if indexed[id] {
						return &engine.Resource{ID: id, Extracted: true, Document: content.Document{Mtime: time.Unix(6000, 0).UTC().Format(time.RFC3339Nano)}}
					}
					return nil
				},
				func(id string) error {
					if indexed[id] {
						return nil
					}
					return engine.ErrResourceNotFound
				},
			)
			indexClient.On("Upsert", mock.Anything, mock.Anything).Return(nil)
		})

		It("does not count a file that no longer exists towards the abort threshold", func() {
			missing["file-0"] = true
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			err := newService(1).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred(), "a file deleted since it was listed is skipped, never counted")
			indexClient.AssertNumberOfCalls(GinkgoT(), "Upsert", numFiles-1)
		})

		It("keeps the failure count across a file that no longer exists", func() {
			// Only file-0 and file-2 fail extraction. file-1 is gone between them.
			// It neither counts nor resets the counter, so a threshold of 2
			// aborts at file-2, before the healthy file-3 could reset it.
			missing["file-1"] = true
			extractor.On("Extract", mock.Anything, mock.MatchedBy(func(ri *sprovider.ResourceInfo) bool {
				return ri.GetId().GetOpaqueId() == "file-0" || ri.GetId().GetOpaqueId() == "file-2"
			}), mock.Anything).Return(content.Document{}, errors.New("extractor down"))
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			err := newService(2).IndexSpace(context.Background(), spaceID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aborting index walk: 2 consecutive extraction failures"))
		})

		It("fails the walk when most files are not found, as in a storage provider outage", func() {
			// The gateway answers NOT_FOUND when it cannot reach the provider.
			for _, fi := range fileInfos {
				missing[fi.Id.OpaqueId] = true
			}

			err := newService(5).IndexSpace(context.Background(), spaceID)
			Expect(err).Should(HaveOccurred(), "an outage must not report the space as indexed")
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("index walk incomplete: %d of %d files were not found", numFiles, numFiles)))
			extractor.AssertNotCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
		})

		It("indexes a stale node that Stat cannot find by its path but can by its ID", func() {
			staleByPath["file-1"] = true
			var extracted []string
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				extracted = append(extracted, args.Get(1).(*sprovider.ResourceInfo).GetId().GetOpaqueId())
			}).Return(content.Document{}, nil)

			err := newService(5).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(extracted).To(ConsistOf("file-0", "file-1", "file-2", "file-3", "file-4", "file-5"))
			gatewayClient.AssertNumberOfCalls(GinkgoT(), "GetPath", 1) // only the stale node needs the ID fallback
		})

		It("does not count a file that is gone between the Stat and the GetPath", func() {
			staleByPath["file-1"] = true // forces the lookup by ID, which calls GetPath
			pathMissing["file-1"] = true
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			err := newService(1).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred(), "a vanished file is not an extraction failure")
			indexClient.AssertNumberOfCalls(GinkgoT(), "Upsert", numFiles-1)
		})

		It("counts a GetPath failure other than not-found", func() {
			staleByPath["file-1"] = true
			pathBroken["file-1"] = true
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			err := newService(1).IndexSpace(context.Background(), spaceID)
			Expect(err).Should(HaveOccurred(), "a storage failure must not look like a deleted file")
			Expect(err.Error()).To(ContainSubstring("aborting index walk: 1 consecutive extraction failures"))
		})

		It("resolves healthy files by path, without a GetPath call", func() {
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, nil)

			err := newService(5).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred())
			indexClient.AssertNumberOfCalls(GinkgoT(), "Upsert", numFiles)
			gatewayClient.AssertNotCalled(GinkgoT(), "GetPath", mock.Anything, mock.Anything)
		})

		It("aborts after the configured number of consecutive extraction failures", func() {
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, errors.New("extractor down"))

			err := newService(2).IndexSpace(context.Background(), spaceID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("aborting index walk: 2 consecutive extraction failures"))
			indexClient.AssertNotCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
		})

		It("never aborts when the threshold is 0", func() {
			var extractCalls atomic.Int32
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
				extractCalls.Add(1)
			}).Return(content.Document{}, errors.New("extractor down"))

			err := newService(0).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred(), "with the abort disabled the walk finishes despite failures")
			Expect(int(extractCalls.Load())).To(Equal(numFiles*6), "every file was attempted, each with all retries")
		})

		It("counts a Stat failure other than not-found towards the abort threshold", func() {
			for _, fi := range fileInfos {
				statBroken[fi.Id.OpaqueId] = true
			}

			err := newService(2).IndexSpace(context.Background(), spaceID)
			Expect(err).Should(HaveOccurred(), "a storage or auth error must not look like a deleted file")
			Expect(err.Error()).To(ContainSubstring("aborting index walk: 2 consecutive extraction failures"))
			extractor.AssertNotCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
		})

		It("stops with an error when the walk's context ends, instead of skipping the rest", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var extractCalls atomic.Int32
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
				extractCalls.Add(1)
				cancel() // e.g. the index command's request timeout fires during the walk
			}).Return(content.Document{}, nil)

			err := newService(0).IndexSpace(ctx, spaceID)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("index walk stopped"))
			Expect(int(extractCalls.Load())).To(Equal(1), "no file after the cancellation was attempted")
		})

		It("stops retrying a failed extraction as soon as the context ends", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var extractCalls atomic.Int32
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Run(func(_ mock.Arguments) {
				extractCalls.Add(1)
				cancel()
			}).Return(content.Document{}, errors.New("extractor down"))

			s := newService(0)
			s.SetExtractionRetryDelay(time.Hour)
			done := make(chan error, 1)
			go func() { done <- s.IndexSpace(ctx, spaceID) }()

			var err error
			Eventually(done, 5*time.Second).Should(Receive(&err), "the retry wait must end with the context, not after an hour")
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("index walk stopped"))
			Expect(int(extractCalls.Load())).To(Equal(1))
		})

		It("resets the counter on an already indexed, unchanged file", func() {
			// Failing files 0, 2 and 4 are each followed by an indexed one.
			for i := 1; i < numFiles; i += 2 {
				indexed[storagespace.FormatResourceID(fileInfos[i].Id)] = true
			}
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{}, errors.New("cannot extract"))

			err := newService(2).IndexSpace(context.Background(), spaceID)
			Expect(err).ShouldNot(HaveOccurred(), "failures separated by a healthy file are not consecutive")
		})
	})

	Describe("UpdateTags", func() {
		BeforeEach(func() {
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
				Status: status.NewOK(context.Background()),
				Info: &sprovider.ResourceInfo{
					Id:       ri.Id,
					ParentId: ri.ParentId,
					Path:     ri.Path,
					Name:     "foo.pdf",
					Size:     ri.Size,
					Mtime:    ri.Mtime,
					Type:     sprovider.ResourceType_RESOURCE_TYPE_FILE,
					ArbitraryMetadata: &sprovider.ArbitraryMetadata{
						Metadata: map[string]string{"tags": "important,work"},
					},
				},
			}, nil)
		})

		It("updates tags without content extraction when resource is indexed", func() {
			indexClient.On("Update", "storageid$spaceid!opaqueid", mock.Anything).Return(nil)

			s.UpdateTags(context.Background(), &sprovider.Reference{ResourceId: ri.Id})

			extractor.AssertNotCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
			indexClient.AssertCalled(GinkgoT(), "Update", "storageid$spaceid!opaqueid", mock.Anything)
		})

		It("falls back to full UpsertItem when resource is not yet indexed", func() {
			indexClient.On("Update", mock.Anything, mock.Anything).Return(engine.ErrResourceNotFound)
			extractor.On("Extract", mock.Anything, mock.Anything, mock.Anything).Return(content.Document{Name: "foo.pdf"}, nil)
			indexClient.On("Upsert", mock.Anything, mock.Anything).Return(nil)

			s.UpdateTags(context.Background(), &sprovider.Reference{ResourceId: ri.Id})

			extractor.AssertCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
		})

		It("logs an error and does not fall back on non-retriable errors", func() {
			indexClient.On("Update", mock.Anything, mock.Anything).Return(errors.New("connection refused"))

			s.UpdateTags(context.Background(), &sprovider.Reference{ResourceId: ri.Id})

			extractor.AssertNotCalled(GinkgoT(), "Extract", mock.Anything, mock.Anything, mock.Anything)
			indexClient.AssertNotCalled(GinkgoT(), "Upsert", mock.Anything, mock.Anything)
		})
	})

	Describe("UpdateTags when GetPath fails", func() {
		for name, getPathStatus := range map[string]func() *rpcv1beta1.Status{
			"with an internal error": func() *rpcv1beta1.Status { return status.NewInternal(context.Background(), "no path") },
			"with not found":         func() *rpcv1beta1.Status { return status.NewNotFound(context.Background(), "no path") },
		} {
			It("still updates the tags from the Stat result, GetPath answering "+name, func() {
				noPathID := &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "nopath"}
				gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
					Status: status.NewOK(context.Background()),
					Info: &sprovider.ResourceInfo{
						Id:   noPathID,
						Type: sprovider.ResourceType_RESOURCE_TYPE_FILE,
						ArbitraryMetadata: &sprovider.ArbitraryMetadata{
							Metadata: map[string]string{"tags": "important"},
						},
					},
				}, nil)
				gatewayClient.On("GetPath", mock.Anything, mock.MatchedBy(func(req *sprovider.GetPathRequest) bool {
					return req.GetResourceId().GetOpaqueId() == "nopath"
				})).Return(&sprovider.GetPathResponse{Status: getPathStatus()}, nil)
				indexClient.On("Update", "storageid$spaceid!nopath", mock.Anything).Return(nil)

				s.UpdateTags(context.Background(), &sprovider.Reference{ResourceId: noPathID})

				indexClient.AssertCalled(GinkgoT(), "Update", "storageid$spaceid!nopath", mock.Anything)
			})
		}
	})

	Describe("Search", func() {
		It("fails when an empty query is given", func() {
			res, err := s.Search(ctx, &searchsvc.SearchRequest{
				Query: "",
			})
			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())
		})

		Context("with a personal space", func() {
			BeforeEach(func() {
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										SpaceId:   personalSpace.Root.SpaceId,
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
				gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
					Status: status.NewOK(context.Background()),
					Info:   ri,
				}, nil)
			})

			It("does not mess with field-based searches", func() {
				_, err := s.Search(ctx, &searchsvc.SearchRequest{
					Query: "Size:<10",
				})
				Expect(err).ToNot(HaveOccurred())
				indexClient.AssertCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Query == "Size:<10"
				}))
			})

			It("searches the personal user space", func() {
				res, err := s.Search(ctx, &searchsvc.SearchRequest{
					Query: "foo",
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.TotalMatches).To(Equal(int32(1)))
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Id.OpaqueId).To(Equal("foo-id"))
				Expect(match.Entity.Name).To(Equal("Foo.pdf"))
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(personalSpace.Root.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal("./path/to/Foo.pdf"))
			})
		})

		Context("with a personal space with a filter", func() {
			BeforeEach(func() {
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace},
				}, nil)
				gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
					Status: status.NewOK(context.Background()),
					Info: &sprovider.ResourceInfo{
						Space: &sprovider.StorageSpace{Root: &sprovider.ResourceId{
							StorageId: "storageid",
							SpaceId:   "personalspace",
							OpaqueId:  "personalspace",
						}},
					},
				}, nil)
				gatewayClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
					Status: status.NewOK(ctx),
					Path:   "/path",
				}, nil)
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										SpaceId:   personalSpace.Root.SpaceId,
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
				res, err := s.Search(ctx, &searchsvc.SearchRequest{
					Query: "foo scope:storageid$personalspace!personalspace/path",
					Ref: &searchmsg.Reference{
						ResourceId: &searchmsg.ResourceID{
							StorageId: "storageid",
							SpaceId:   "personalspace",
							OpaqueId:  "personalspace",
						},
					},
				})

				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(res.TotalMatches).To(Equal(int32(1)))
				Expect(len(res.Matches)).To(Equal(1))
				match := res.Matches[0]
				Expect(match.Entity.Id.OpaqueId).To(Equal("foo-id"))
				Expect(match.Entity.Name).To(Equal("Foo.pdf"))
				Expect(match.Entity.Ref.ResourceId.OpaqueId).To(Equal(personalSpace.Root.OpaqueId))
				Expect(match.Entity.Ref.Path).To(Equal("./path/to/Foo.pdf"))
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
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacegrant"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacegrant"},
					Name:      "grantspace",
					RootInfo: &sprovider.ResourceInfo{
						Id: &sprovider.ResourceId{
							StorageId: "storageid",
							SpaceId:   "spaceid",
							OpaqueId:  "opaqueid",
						},
					},
				}
				mountpointSpace = &sprovider.StorageSpace{
					SpaceType: "mountpoint",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacemountpoint"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacemountpoint"},
					Name:      "mountpointspace",
					Opaque: &typesv1beta1.Opaque{
						Map: map[string]*typesv1beta1.OpaqueEntry{
							"grantStorageID": {Decoder: "plain", Value: []byte("storageproviderid")},
							"grantSpaceID":   {Decoder: "plain", Value: []byte("spaceid")},
							"grantOpaqueID":  {Decoder: "plain", Value: []byte("otherspacegrant")},
						},
					},
				}
				gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
					Status: status.NewOK(context.Background()),
					Info:   ri,
				}, nil)
				gatewayClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
					Status: status.NewOK(ctx),
					Path:   "/grant/path",
				}, nil)
			})

			It("searches the received spaces", func() {
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{grantSpace, mountpointSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.Anything).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: grantSpace.Root.StorageId,
										SpaceId:   grantSpace.Root.SpaceId,
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

				res, err := s.Search(ctx, &searchsvc.SearchRequest{
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
				Expect(match.Entity.RemoteItemId).To(Equal(&searchmsg.ResourceID{
					StorageId: grantSpace.RootInfo.Id.StorageId,
					SpaceId:   grantSpace.RootInfo.Id.SpaceId,
					OpaqueId:  grantSpace.RootInfo.Id.OpaqueId,
				}))
			})

			Context("when searching both spaces", func() {
				BeforeEach(func() {
					gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
						Status:        status.NewOK(ctx),
						StorageSpaces: []*sprovider.StorageSpace{personalSpace, grantSpace, mountpointSpace},
					}, nil)
					indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
						return req.Ref.ResourceId.OpaqueId == grantSpace.Root.SpaceId &&
							req.Ref.ResourceId.SpaceId == grantSpace.Root.SpaceId
					})).Return(&searchsvc.SearchIndexResponse{
						TotalMatches: 2,
						Matches: []*searchmsg.Match{
							{
								Score: 2,
								Entity: &searchmsg.Entity{
									Ref: &searchmsg.Reference{
										ResourceId: &searchmsg.ResourceID{
											StorageId: grantSpace.Root.StorageId,
											SpaceId:   grantSpace.Root.SpaceId,
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
							{
								Score: 0.01,
								Entity: &searchmsg.Entity{
									Ref: &searchmsg.Reference{
										ResourceId: &searchmsg.ResourceID{
											StorageId: grantSpace.Root.StorageId,
											SpaceId:   grantSpace.Root.SpaceId,
											OpaqueId:  grantSpace.Root.OpaqueId,
										},
										Path: "./grant/path/to/Irrelevant.pdf",
									},
									Id: &searchmsg.ResourceID{
										StorageId: grantSpace.Root.StorageId,
										OpaqueId:  "grant-irrelevant-id",
									},
									Name: "Irrelevant.pdf",
								},
							},
						},
					}, nil)
					indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
						return req.Ref.ResourceId.OpaqueId == personalSpace.Root.OpaqueId &&
							req.Ref.ResourceId.SpaceId == personalSpace.Root.SpaceId
					})).Return(&searchsvc.SearchIndexResponse{
						TotalMatches: 1,
						Matches: []*searchmsg.Match{
							{
								Score: 1,
								Entity: &searchmsg.Entity{
									Ref: &searchmsg.Reference{
										ResourceId: &searchmsg.ResourceID{
											StorageId: personalSpace.Root.StorageId,
											SpaceId:   personalSpace.Root.SpaceId,
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

				It("considers the search Ref parameter", func() {
					res, err := s.Search(ctx, &searchsvc.SearchRequest{
						Query: "foo",
						Ref: &searchmsg.Reference{
							ResourceId: &searchmsg.ResourceID{
								StorageId: "storageid",
								SpaceId:   "personalspace",
								OpaqueId:  "personalspace",
							},
						},
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(1))
					Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("foo-id"))
				})

				It("finds matches in both the personal space AND the grant", func() {
					res, err := s.Search(ctx, &searchsvc.SearchRequest{
						Query: "foo",
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(3))
					ids := []string{res.Matches[0].Entity.Id.OpaqueId, res.Matches[1].Entity.Id.OpaqueId, res.Matches[2].Entity.Id.OpaqueId}
					Expect(ids).To(ConsistOf("foo-id", "grant-shared-id", "grant-irrelevant-id"))
				})

				It("sorts and limits the combined results from all spaces", func() {
					res, err := s.Search(ctx, &searchsvc.SearchRequest{
						Query:    "foo",
						PageSize: 2,
					})
					Expect(err).ToNot(HaveOccurred())
					Expect(res).ToNot(BeNil())
					Expect(len(res.Matches)).To(Equal(2))
					ids := []string{res.Matches[0].Entity.Id.OpaqueId, res.Matches[1].Entity.Id.OpaqueId}
					Expect(ids).To(Equal([]string{"grant-shared-id", "foo-id"}))
				})
			})
		})

		Context("vault mode", func() {
			var vaultSpace *sprovider.StorageSpace

			BeforeEach(func() {
				vaultSpace = &sprovider.StorageSpace{
					SpaceType: "personal",
					Owner:     user,
					Id:        &sprovider.StorageSpaceId{OpaqueId: utils.VaultStorageProviderID + "$vaultspace!vaultspace"},
					Root:      &sprovider.ResourceId{StorageId: utils.VaultStorageProviderID, SpaceId: "vaultspace", OpaqueId: "vaultspace"},
					Name:      "vaultspace",
				}
				gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(&sprovider.StatResponse{
					Status: status.NewOK(context.Background()),
					Info:   ri,
				}, nil)
			})

			It("searches only the vault space when vault:true is in query AND context has MFA", func() {
				mfaCtx := revactx.SetMFA(ctx)
				// First call: general spaces (no storage_id opaque key) — returns only the personal space.
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					return utils.ReadPlainFromOpaque(req.Opaque, "storage_id") == ""
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(mfaCtx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace},
				}, nil)
				// Second call: vault-specific (storage_id set) — returns the vault space.
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					return utils.ReadPlainFromOpaque(req.Opaque, "storage_id") == utils.VaultStorageProviderID
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(mfaCtx),
					StorageSpaces: []*sprovider.StorageSpace{vaultSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == vaultSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: vaultSpace.Root.StorageId,
										SpaceId:   vaultSpace.Root.SpaceId,
										OpaqueId:  vaultSpace.Root.OpaqueId,
									},
									Path: "./vault-file.pdf",
								},
								Id: &searchmsg.ResourceID{
									StorageId: vaultSpace.Root.StorageId,
									OpaqueId:  "vault-file-id",
								},
								Name: "vault-file.pdf",
							},
						},
					},
				}, nil)

				res, err := s.Search(mfaCtx, &searchsvc.SearchRequest{Query: "secret vault:true"})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("vault-file-id"))
			})

			It("searches only regular spaces when MFA context is set but vault:true is absent", func() {
				mfaCtx := revactx.SetMFA(ctx)
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(mfaCtx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, vaultSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == personalSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										SpaceId:   personalSpace.Root.SpaceId,
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

				res, err := s.Search(mfaCtx, &searchsvc.SearchRequest{Query: "secret"})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("foo-id"))
			})

			It("searches only regular spaces when vault:true is in query but context has no MFA", func() {
				// ctx has no MFA — a non-MFA client cannot opt into vault mode
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, vaultSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == personalSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: personalSpace.Root.StorageId,
										SpaceId:   personalSpace.Root.SpaceId,
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

				res, err := s.Search(ctx, &searchsvc.SearchRequest{Query: "secret vault:true"})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("foo-id"))
			})

			It("regular search includes received shares but excludes vault spaces", func() {
				grantSpace := &sprovider.StorageSpace{
					SpaceType: "grant",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacegrant"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacegrant"},
					Name:      "grantspace",
					RootInfo: &sprovider.ResourceInfo{
						Id: &sprovider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
					},
				}
				mountpointSpace := &sprovider.StorageSpace{
					SpaceType: "mountpoint",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacemountpoint"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacemountpoint"},
					Name:      "mountpointspace",
					Opaque: &typesv1beta1.Opaque{
						Map: map[string]*typesv1beta1.OpaqueEntry{
							"grantStorageID": {Decoder: "plain", Value: []byte("storageproviderid")},
							"grantSpaceID":   {Decoder: "plain", Value: []byte("spaceid")},
							"grantOpaqueID":  {Decoder: "plain", Value: []byte("otherspacegrant")},
						},
					},
				}
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(ctx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, grantSpace, mountpointSpace, vaultSpace},
				}, nil)
				gatewayClient.On("GetPath", mock.Anything, mock.Anything).Return(&sprovider.GetPathResponse{
					Status: status.NewOK(ctx),
					Path:   "/grant/path",
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == grantSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: grantSpace.Root.StorageId,
										SpaceId:   grantSpace.Root.SpaceId,
										OpaqueId:  grantSpace.Root.OpaqueId,
									},
									Path: "./grant/path/to/Shared.pdf",
								},
								Id:   &searchmsg.ResourceID{StorageId: grantSpace.Root.StorageId, OpaqueId: "shared-id"},
								Name: "Shared.pdf",
							},
						},
					},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == personalSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{}, nil)

				res, err := s.Search(ctx, &searchsvc.SearchRequest{Query: "file"})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				// Result is the share, mapped to the mountpoint path.
				Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("shared-id"))
				Expect(res.Matches[0].Entity.Ref.ResourceId.OpaqueId).To(Equal(mountpointSpace.Root.OpaqueId))
				// Vault space was not searched.
				indexClient.AssertNotCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.StorageId == utils.VaultStorageProviderID
				}))
			})

			It("vault mode search excludes received shares from regular storage", func() {
				grantSpace := &sprovider.StorageSpace{
					SpaceType: "grant",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacegrant"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacegrant"},
					Name:      "grantspace",
				}
				mountpointSpace := &sprovider.StorageSpace{
					SpaceType: "mountpoint",
					Owner:     otherUser,
					Id:        &sprovider.StorageSpaceId{OpaqueId: "storageproviderid$spaceid!otherspacemountpoint"},
					Root:      &sprovider.ResourceId{StorageId: "storageproviderid", SpaceId: "spaceid", OpaqueId: "otherspacemountpoint"},
					Name:      "mountpointspace",
					Opaque: &typesv1beta1.Opaque{
						Map: map[string]*typesv1beta1.OpaqueEntry{
							"grantStorageID": {Decoder: "plain", Value: []byte("storageproviderid")},
							"grantSpaceID":   {Decoder: "plain", Value: []byte("spaceid")},
							"grantOpaqueID":  {Decoder: "plain", Value: []byte("otherspacegrant")},
						},
					},
				}
				mfaCtx := revactx.SetMFA(ctx)
				// First LSS call: general — returns personal + grant + mountpoint (no vault).
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					return utils.ReadPlainFromOpaque(req.Opaque, "storage_id") == ""
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(mfaCtx),
					StorageSpaces: []*sprovider.StorageSpace{personalSpace, grantSpace, mountpointSpace},
				}, nil)
				// Second LSS call: vault-specific — returns vault space.
				gatewayClient.On("ListStorageSpaces", mock.Anything, mock.MatchedBy(func(req *sprovider.ListStorageSpacesRequest) bool {
					return utils.ReadPlainFromOpaque(req.Opaque, "storage_id") == utils.VaultStorageProviderID
				})).Return(&sprovider.ListStorageSpacesResponse{
					Status:        status.NewOK(mfaCtx),
					StorageSpaces: []*sprovider.StorageSpace{vaultSpace},
				}, nil)
				indexClient.On("Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == vaultSpace.Root.SpaceId
				})).Return(&searchsvc.SearchIndexResponse{
					TotalMatches: 1,
					Matches: []*searchmsg.Match{
						{
							Score: 1,
							Entity: &searchmsg.Entity{
								Ref: &searchmsg.Reference{
									ResourceId: &searchmsg.ResourceID{
										StorageId: vaultSpace.Root.StorageId,
										SpaceId:   vaultSpace.Root.SpaceId,
										OpaqueId:  vaultSpace.Root.OpaqueId,
									},
									Path: "./vault-file.pdf",
								},
								Id:   &searchmsg.ResourceID{StorageId: vaultSpace.Root.StorageId, OpaqueId: "vault-file-id"},
								Name: "vault-file.pdf",
							},
						},
					},
				}, nil)

				res, err := s.Search(mfaCtx, &searchsvc.SearchRequest{Query: "file vault:true"})
				Expect(err).ToNot(HaveOccurred())
				Expect(res).ToNot(BeNil())
				Expect(len(res.Matches)).To(Equal(1))
				Expect(res.Matches[0].Entity.Id.OpaqueId).To(Equal("vault-file-id"))
				// Grant/mountpoint spaces were not searched.
				indexClient.AssertNotCalled(GinkgoT(), "Search", mock.Anything, mock.MatchedBy(func(req *searchsvc.SearchIndexRequest) bool {
					return req.Ref.ResourceId.SpaceId == grantSpace.Root.SpaceId
				}))
			})
		})
	})
})

var _ = DescribeTable("Parse Vault Mode",
	func(pattern, wantQuery string, wantVault bool) {
		gotQuery, gotVault := search.ParseVaultMode(pattern)
		Expect(gotQuery).To(Equal(wantQuery))
		Expect(gotVault).To(Equal(wantVault))
	},
	Entry("vault:true present",
		`file vault:true`,
		`file`,
		true,
	),
	Entry("vault:true at the beginning",
		`vault:true file`,
		`file`,
		true,
	),
	Entry("vault:true in the middle",
		`foo vault:true bar`,
		`foo  bar`,
		true,
	),
	Entry("vault:false is not vault mode",
		`+Name:*file* vault:false +Tags:&quot;foo&quot;`,
		`+Name:*file*  +Tags:&quot;foo&quot;`,
		false,
	),
	Entry("vault:false is not vault mode",
		`file vault:false`,
		`file`,
		false,
	),
	Entry("no vault token",
		`+Name:*file* +Tags:&quot;foo&quot;`,
		`+Name:*file* +Tags:&quot;foo&quot;`,
		false,
	),
)

var _ = DescribeTable("Parse Scope",
	func(pattern, wantSearch, wantScope string) {
		gotSearch, gotScope := search.ParseScope(pattern)
		Expect(gotSearch).To(Equal(wantSearch))
		Expect(gotScope).To(Equal(wantScope))
	},
	Entry("When scope is at the end of the line",
		`+Name:*file* +Tags:&quot;foo&quot; scope:<uuid>/folder/subfolder`,
		`+Name:*file* +Tags:&quot;foo&quot;`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When scope is at the end of the line 2",
		`+Name:*file* +Tags:&quot;foo&quot; scope:<uuid>/folder`,
		`+Name:*file* +Tags:&quot;foo&quot;`,
		`<uuid>/folder`,
	),
	Entry("When scope is at the end of the line 3",
		`file scope:<uuid>/folder/subfolder`,
		`file`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When scope is at the end of the line with a space",
		`+Name:*file* +Tags:&quot;foo&quot; scope: <uuid>/folder/subfolder`,
		`+Name:*file* +Tags:&quot;foo&quot;`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When scope is in the middle of the line",
		`+Name:*file* scope:<uuid>/folder/subfolder +Tags:&quot;foo&quot;`,
		`+Name:*file*  +Tags:&quot;foo&quot;`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When scope is at the end of the line",
		`scope:<uuid>/folder/subfolder +Name:*file*`,
		`+Name:*file*`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When scope is at the begging of the line",
		`scope:<uuid>/folder/subfolder file`,
		`file`,
		`<uuid>/folder/subfolder`,
	),
	Entry("When no scope",
		`+Name:*file* +Tags:&quot;foo&quot;`,
		`+Name:*file* +Tags:&quot;foo&quot;`,
		``,
	),
)

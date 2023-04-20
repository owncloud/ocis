// Copyright 2018-2021 CERN
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

package gateway

import (
	"context"
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/storage/cache"
	"github.com/cs3org/reva/v2/pkg/utils"
	"google.golang.org/grpc"
)

/*
   Cached Registry
*/

type cachedRegistryClient struct {
	c     registry.RegistryAPIClient
	cache cache.ProviderCache
}

func (c *cachedRegistryClient) ListStorageProviders(ctx context.Context, in *registry.ListStorageProvidersRequest, opts ...grpc.CallOption) (*registry.ListStorageProvidersResponse, error) {

	spaceID := sdk.DecodeOpaqueMap(in.Opaque)["space_id"]

	key := c.cache.GetKey(ctxpkg.ContextMustGetUser(ctx).GetId(), spaceID)
	if key != "" {
		s := &registry.ListStorageProvidersResponse{}
		if err := c.cache.PullFromCache(key, s); err == nil {
			return s, nil
		}
	}

	resp, err := c.c.ListStorageProviders(ctx, in, opts...)
	switch {
	case err != nil:
		return nil, err
	case resp.Status.Code != rpc.Code_CODE_OK:
		return resp, nil
	case spaceID == "":
		return resp, nil
	case spaceID == utils.ShareStorageSpaceID: // TODO do we need to compare providerid and spaceid separately?
		return resp, nil
	default:
		return resp, c.cache.PushToCache(key, resp)
	}
}

// not cached

func (c *cachedRegistryClient) GetStorageProviders(ctx context.Context, in *registry.GetStorageProvidersRequest, opts ...grpc.CallOption) (*registry.GetStorageProvidersResponse, error) {
	return c.c.GetStorageProviders(ctx, in, opts...)
}

func (c *cachedRegistryClient) GetHome(ctx context.Context, in *registry.GetHomeRequest, opts ...grpc.CallOption) (*registry.GetHomeResponse, error) {
	return c.c.GetHome(ctx, in, opts...)
}

/*
   Cached Storage Provider
*/

type cachedAPIClient struct {
	c                        provider.ProviderAPIClient
	statCache                cache.StatCache
	createHomeCache          cache.CreateHomeCache
	createPersonalSpaceCache cache.CreatePersonalSpaceCache
}

// Stat looks in cache first before forwarding to storage provider
func (c *cachedAPIClient) Stat(ctx context.Context, in *provider.StatRequest, opts ...grpc.CallOption) (*provider.StatResponse, error) {
	key := c.statCache.GetKey(ctxpkg.ContextMustGetUser(ctx).GetId(), in.GetRef(), in.GetArbitraryMetadataKeys(), in.GetFieldMask().GetPaths())
	if key != "" {
		s := &provider.StatResponse{}
		if err := c.statCache.PullFromCache(key, s); err == nil {
			return s, nil
		}
	}

	resp, err := c.c.Stat(ctx, in, opts...)
	switch {
	case err != nil:
		return nil, err
	case resp.Status.Code != rpc.Code_CODE_OK:
		return resp, nil
	case key == "":
		return resp, nil
	case strings.Contains(key, "sid:"+utils.ShareStorageProviderID):
		// We cannot cache shares at the moment:
		// we do not know when to invalidate them
		// FIXME: find a way to cache/invalidate them too
		return resp, nil
	case utils.ReadPlainFromOpaque(resp.GetInfo().GetOpaque(), "status") == "processing":
		return resp, nil
	default:
		return resp, c.statCache.PushToCache(key, resp)
	}
}

// CreateHome caches calls to CreateHome locally - anyways they only need to be called once per user
func (c *cachedAPIClient) CreateHome(ctx context.Context, in *provider.CreateHomeRequest, opts ...grpc.CallOption) (*provider.CreateHomeResponse, error) {
	key := c.createHomeCache.GetKey(ctxpkg.ContextMustGetUser(ctx).GetId())
	if key != "" {
		s := &provider.CreateHomeResponse{}
		if err := c.createHomeCache.PullFromCache(key, s); err == nil {
			return s, nil
		}
	}
	resp, err := c.c.CreateHome(ctx, in, opts...)
	switch {
	case err != nil:
		return nil, err
	case resp.Status.Code != rpc.Code_CODE_OK && resp.Status.Code != rpc.Code_CODE_ALREADY_EXISTS:
		return resp, nil
	case key == "":
		return resp, nil
	default:
		return resp, c.createHomeCache.PushToCache(key, resp)
	}
}

// methods below here are not cached, they just call the client directly

func (c *cachedAPIClient) AddGrant(ctx context.Context, in *provider.AddGrantRequest, opts ...grpc.CallOption) (*provider.AddGrantResponse, error) {
	return c.c.AddGrant(ctx, in, opts...)
}
func (c *cachedAPIClient) CreateContainer(ctx context.Context, in *provider.CreateContainerRequest, opts ...grpc.CallOption) (*provider.CreateContainerResponse, error) {
	return c.c.CreateContainer(ctx, in, opts...)
}
func (c *cachedAPIClient) Delete(ctx context.Context, in *provider.DeleteRequest, opts ...grpc.CallOption) (*provider.DeleteResponse, error) {
	return c.c.Delete(ctx, in, opts...)
}
func (c *cachedAPIClient) DenyGrant(ctx context.Context, in *provider.DenyGrantRequest, opts ...grpc.CallOption) (*provider.DenyGrantResponse, error) {
	return c.c.DenyGrant(ctx, in, opts...)
}
func (c *cachedAPIClient) GetPath(ctx context.Context, in *provider.GetPathRequest, opts ...grpc.CallOption) (*provider.GetPathResponse, error) {
	return c.c.GetPath(ctx, in, opts...)
}
func (c *cachedAPIClient) GetQuota(ctx context.Context, in *provider.GetQuotaRequest, opts ...grpc.CallOption) (*provider.GetQuotaResponse, error) {
	return c.c.GetQuota(ctx, in, opts...)
}
func (c *cachedAPIClient) InitiateFileDownload(ctx context.Context, in *provider.InitiateFileDownloadRequest, opts ...grpc.CallOption) (*provider.InitiateFileDownloadResponse, error) {
	return c.c.InitiateFileDownload(ctx, in, opts...)
}
func (c *cachedAPIClient) InitiateFileUpload(ctx context.Context, in *provider.InitiateFileUploadRequest, opts ...grpc.CallOption) (*provider.InitiateFileUploadResponse, error) {
	return c.c.InitiateFileUpload(ctx, in, opts...)
}
func (c *cachedAPIClient) ListGrants(ctx context.Context, in *provider.ListGrantsRequest, opts ...grpc.CallOption) (*provider.ListGrantsResponse, error) {
	return c.c.ListGrants(ctx, in, opts...)
}
func (c *cachedAPIClient) ListContainerStream(ctx context.Context, in *provider.ListContainerStreamRequest, opts ...grpc.CallOption) (provider.ProviderAPI_ListContainerStreamClient, error) {
	return c.c.ListContainerStream(ctx, in, opts...)
}
func (c *cachedAPIClient) ListContainer(ctx context.Context, in *provider.ListContainerRequest, opts ...grpc.CallOption) (*provider.ListContainerResponse, error) {
	return c.c.ListContainer(ctx, in, opts...)
}
func (c *cachedAPIClient) ListFileVersions(ctx context.Context, in *provider.ListFileVersionsRequest, opts ...grpc.CallOption) (*provider.ListFileVersionsResponse, error) {
	return c.c.ListFileVersions(ctx, in, opts...)
}
func (c *cachedAPIClient) ListRecycleStream(ctx context.Context, in *provider.ListRecycleStreamRequest, opts ...grpc.CallOption) (provider.ProviderAPI_ListRecycleStreamClient, error) {
	return c.c.ListRecycleStream(ctx, in, opts...)
}
func (c *cachedAPIClient) ListRecycle(ctx context.Context, in *provider.ListRecycleRequest, opts ...grpc.CallOption) (*provider.ListRecycleResponse, error) {
	return c.c.ListRecycle(ctx, in, opts...)
}
func (c *cachedAPIClient) Move(ctx context.Context, in *provider.MoveRequest, opts ...grpc.CallOption) (*provider.MoveResponse, error) {
	return c.c.Move(ctx, in, opts...)
}
func (c *cachedAPIClient) RemoveGrant(ctx context.Context, in *provider.RemoveGrantRequest, opts ...grpc.CallOption) (*provider.RemoveGrantResponse, error) {
	return c.c.RemoveGrant(ctx, in, opts...)
}
func (c *cachedAPIClient) PurgeRecycle(ctx context.Context, in *provider.PurgeRecycleRequest, opts ...grpc.CallOption) (*provider.PurgeRecycleResponse, error) {
	return c.c.PurgeRecycle(ctx, in, opts...)
}
func (c *cachedAPIClient) RestoreFileVersion(ctx context.Context, in *provider.RestoreFileVersionRequest, opts ...grpc.CallOption) (*provider.RestoreFileVersionResponse, error) {
	return c.c.RestoreFileVersion(ctx, in, opts...)
}
func (c *cachedAPIClient) RestoreRecycleItem(ctx context.Context, in *provider.RestoreRecycleItemRequest, opts ...grpc.CallOption) (*provider.RestoreRecycleItemResponse, error) {
	return c.c.RestoreRecycleItem(ctx, in, opts...)
}
func (c *cachedAPIClient) UpdateGrant(ctx context.Context, in *provider.UpdateGrantRequest, opts ...grpc.CallOption) (*provider.UpdateGrantResponse, error) {
	return c.c.UpdateGrant(ctx, in, opts...)
}
func (c *cachedAPIClient) CreateSymlink(ctx context.Context, in *provider.CreateSymlinkRequest, opts ...grpc.CallOption) (*provider.CreateSymlinkResponse, error) {
	return c.c.CreateSymlink(ctx, in, opts...)
}
func (c *cachedAPIClient) CreateReference(ctx context.Context, in *provider.CreateReferenceRequest, opts ...grpc.CallOption) (*provider.CreateReferenceResponse, error) {
	return c.c.CreateReference(ctx, in, opts...)
}
func (c *cachedAPIClient) SetArbitraryMetadata(ctx context.Context, in *provider.SetArbitraryMetadataRequest, opts ...grpc.CallOption) (*provider.SetArbitraryMetadataResponse, error) {
	return c.c.SetArbitraryMetadata(ctx, in, opts...)
}
func (c *cachedAPIClient) UnsetArbitraryMetadata(ctx context.Context, in *provider.UnsetArbitraryMetadataRequest, opts ...grpc.CallOption) (*provider.UnsetArbitraryMetadataResponse, error) {
	return c.c.UnsetArbitraryMetadata(ctx, in, opts...)
}
func (c *cachedAPIClient) SetLock(ctx context.Context, in *provider.SetLockRequest, opts ...grpc.CallOption) (*provider.SetLockResponse, error) {
	return c.c.SetLock(ctx, in, opts...)
}
func (c *cachedAPIClient) GetLock(ctx context.Context, in *provider.GetLockRequest, opts ...grpc.CallOption) (*provider.GetLockResponse, error) {
	return c.c.GetLock(ctx, in, opts...)
}
func (c *cachedAPIClient) RefreshLock(ctx context.Context, in *provider.RefreshLockRequest, opts ...grpc.CallOption) (*provider.RefreshLockResponse, error) {
	return c.c.RefreshLock(ctx, in, opts...)
}
func (c *cachedAPIClient) Unlock(ctx context.Context, in *provider.UnlockRequest, opts ...grpc.CallOption) (*provider.UnlockResponse, error) {
	return c.c.Unlock(ctx, in, opts...)
}
func (c *cachedAPIClient) GetHome(ctx context.Context, in *provider.GetHomeRequest, opts ...grpc.CallOption) (*provider.GetHomeResponse, error) {
	return c.c.GetHome(ctx, in, opts...)
}
func (c *cachedAPIClient) CreateStorageSpace(ctx context.Context, in *provider.CreateStorageSpaceRequest, opts ...grpc.CallOption) (*provider.CreateStorageSpaceResponse, error) {
	if in.Type == "personal" {
		key := c.createPersonalSpaceCache.GetKey(ctxpkg.ContextMustGetUser(ctx).GetId())
		if key != "" {
			s := &provider.CreateStorageSpaceResponse{}
			if err := c.createPersonalSpaceCache.PullFromCache(key, s); err == nil {
				return s, nil
			}
		}
		resp, err := c.c.CreateStorageSpace(ctx, in, opts...)
		switch {
		case err != nil:
			return nil, err
		case resp.Status.Code != rpc.Code_CODE_OK && resp.Status.Code != rpc.Code_CODE_ALREADY_EXISTS:
			return resp, nil
		case key == "":
			return resp, nil
		default:
			return resp, c.createPersonalSpaceCache.PushToCache(key, resp)
		}
	}
	return c.c.CreateStorageSpace(ctx, in, opts...)
}
func (c *cachedAPIClient) ListStorageSpaces(ctx context.Context, in *provider.ListStorageSpacesRequest, opts ...grpc.CallOption) (*provider.ListStorageSpacesResponse, error) {
	return c.c.ListStorageSpaces(ctx, in, opts...)
}
func (c *cachedAPIClient) UpdateStorageSpace(ctx context.Context, in *provider.UpdateStorageSpaceRequest, opts ...grpc.CallOption) (*provider.UpdateStorageSpaceResponse, error) {
	return c.c.UpdateStorageSpace(ctx, in, opts...)
}
func (c *cachedAPIClient) DeleteStorageSpace(ctx context.Context, in *provider.DeleteStorageSpaceRequest, opts ...grpc.CallOption) (*provider.DeleteStorageSpaceResponse, error) {
	return c.c.DeleteStorageSpace(ctx, in, opts...)
}

func (c *cachedAPIClient) TouchFile(ctx context.Context, in *provider.TouchFileRequest, opts ...grpc.CallOption) (*provider.TouchFileResponse, error) {
	return c.c.TouchFile(ctx, in, opts...)
}

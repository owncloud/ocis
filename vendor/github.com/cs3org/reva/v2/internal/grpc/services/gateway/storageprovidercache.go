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
	"encoding/json"
	"strings"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	sdk "github.com/cs3org/reva/v2/pkg/sdk/common"
	"github.com/cs3org/reva/v2/pkg/utils"
	"google.golang.org/grpc"
)

// available caches
const (
	stat = iota
	createhome
	listproviders
)

// allCaches is needed for initialization
var allCaches = []int{stat, createhome, listproviders}

// Caches holds all caches used by the gateway
type Caches []*ttlcache.Cache

// NewCaches initializes the caches. Optionally takes ttls for all caches.
// len(ttlsSeconds) is expected to be either 0, 1 or len(allCaches). Panics if not.
func NewCaches(ttlsSeconds ...int) Caches {
	numCaches := len(allCaches)

	ttls := make([]int, numCaches)
	switch len(ttlsSeconds) {
	case 0:
		// already done
	case 1:
		for i := 0; i < numCaches; i++ {
			ttls[i] = ttlsSeconds[0]
		}
	case numCaches:
		for i := 0; i < numCaches; i++ {
			ttls[i] = ttlsSeconds[i]
		}
	default:
		panic("caching misconfigured - pass 0, 1 or len(allCaches) arguments to NewCaches")
	}

	c := Caches{}
	for i := range allCaches {
		c = append(c, initCache(ttls[i]))
	}
	return c
}

// Close closes all caches - best to call it on teardown - ignores errors
func (c Caches) Close() {
	for _, cache := range c {
		cache.Close()
	}
}

// StorageProviderClient returns a (cached) client pointing to the storageprovider
func (c Caches) StorageProviderClient(p provider.ProviderAPIClient) provider.ProviderAPIClient {
	return &cachedAPIClient{
		c:      p,
		caches: c,
	}
}

// StorageRegistryClient returns a (cached) client pointing to the storageregistry
func (c Caches) StorageRegistryClient(p registry.RegistryAPIClient) registry.RegistryAPIClient {
	return &cachedRegistryClient{
		c:      p,
		caches: c,
	}
}

// RemoveStat removes a reference from the stat cache
func (c Caches) RemoveStat(user *userpb.User, res *provider.ResourceId) {
	uid := "uid:" + user.Id.OpaqueId
	sid := ""
	oid := ""
	if res != nil {
		sid = "sid:" + res.StorageId
		oid = "oid:" + res.OpaqueId
	}

	cache := c[stat]
	for _, key := range cache.GetKeys() {
		if strings.Contains(key, uid) {
			_ = cache.Remove(key)
			continue
		}

		if sid != "" && strings.Contains(key, sid) {
			_ = cache.Remove(key)
			continue
		}

		if oid != "" && strings.Contains(key, oid) {
			_ = cache.Remove(key)
			continue
		}
	}
}

// RemoveListStorageProviders removes a reference from the listproviders cache
func (c Caches) RemoveListStorageProviders(res *provider.ResourceId) {
	if res == nil {
		return
	}
	sid := res.StorageId

	cache := c[listproviders]
	for _, key := range cache.GetKeys() {
		if strings.Contains(key, sid) {
			_ = cache.Remove(key)
			continue
		}
	}
}

func initCache(ttlSeconds int) *ttlcache.Cache {
	cache := ttlcache.NewCache()
	_ = cache.SetTTL(time.Duration(ttlSeconds) * time.Second)
	cache.SkipTTLExtensionOnHit(true)
	return cache
}

func pullFromCache(cache *ttlcache.Cache, key string, dest interface{}) error {
	r, err := cache.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(r.([]byte), dest)
}

func pushToCache(cache *ttlcache.Cache, key string, src interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return cache.Set(key, b)
}

/*
   Cached Registry
*/

type cachedRegistryClient struct {
	c      registry.RegistryAPIClient
	caches Caches
}

func (c *cachedRegistryClient) ListStorageProviders(ctx context.Context, in *registry.ListStorageProvidersRequest, opts ...grpc.CallOption) (*registry.ListStorageProvidersResponse, error) {
	cache := c.caches[listproviders]

	user := ctxpkg.ContextMustGetUser(ctx)

	storageID := sdk.DecodeOpaqueMap(in.Opaque)["storage_id"]

	key := user.GetId().GetOpaqueId() + "!" + storageID
	if key != "!" {
		s := &registry.ListStorageProvidersResponse{}
		if err := pullFromCache(cache, key, s); err == nil {
			return s, nil
		}
	}

	resp, err := c.c.ListStorageProviders(ctx, in, opts...)
	switch {
	case err != nil:
		return nil, err
	case resp.Status.Code != rpc.Code_CODE_OK:
		return resp, nil
	case storageID == "":
		return resp, nil
	case storageID == utils.ShareStorageProviderID: // TODO do we need to compare providerid and spaceid separately?
		return resp, nil
	default:
		return resp, pushToCache(cache, key, resp)
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
	c      provider.ProviderAPIClient
	caches Caches
}

// generates a user specific key pointing to ref - used for statcache
// a key looks like: uid:1234-1233!sid:5678-5677!oid:9923-9934!path:/path/to/source
// as you see it adds "uid:"/"sid:"/"oid:" prefixes to the uuids so they can be differentiated
func statKey(user *userpb.User, ref *provider.Reference, metaDataKeys []string) string {
	if ref == nil || ref.ResourceId == nil || ref.ResourceId.StorageId == "" {
		return ""
	}

	key := "uid:" + user.Id.OpaqueId + "!sid:" + ref.ResourceId.StorageId + "!oid:" + ref.ResourceId.OpaqueId + "!path:" + ref.Path
	for _, k := range metaDataKeys {
		key += "!mdk:" + k
	}

	return key
}

// Stat looks in cache first before forwarding to storage provider
func (c *cachedAPIClient) Stat(ctx context.Context, in *provider.StatRequest, opts ...grpc.CallOption) (*provider.StatResponse, error) {
	cache := c.caches[stat]

	key := statKey(ctxpkg.ContextMustGetUser(ctx), in.Ref, in.ArbitraryMetadataKeys)
	if key != "" {
		s := &provider.StatResponse{}
		if err := pullFromCache(cache, key, s); err == nil {
			return s, nil
		}
	}
	resp, err := c.c.Stat(ctx, in, opts...)
	switch {
	case err != nil:
		return nil, err
	case resp.Status.Code != rpc.Code_CODE_OK && resp.Status.Code != rpc.Code_CODE_NOT_FOUND:
		return resp, nil
	case key == "":
		return resp, nil
	case strings.Contains(key, "sid:"+utils.ShareStorageProviderID):
		// We cannot cache shares at the moment:
		// we do not know when to invalidate them
		// FIXME: find a way to cache/invalidate them too
		return resp, nil
	default:
		return resp, pushToCache(cache, key, resp)
	}
}

// CreateHome caches calls to CreateHome locally - anyways they only need to be called once per user
func (c *cachedAPIClient) CreateHome(ctx context.Context, in *provider.CreateHomeRequest, opts ...grpc.CallOption) (*provider.CreateHomeResponse, error) {
	cache := c.caches[createhome]

	key := ctxpkg.ContextMustGetUser(ctx).Id.OpaqueId
	if key != "" {
		s := &provider.CreateHomeResponse{}
		if err := pullFromCache(cache, key, s); err == nil {
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
		return resp, pushToCache(cache, key, resp)
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

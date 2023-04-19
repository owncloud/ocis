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

package storage

import (
	"context"
	"io"
	"net/url"

	tusd "github.com/tus/tusd/pkg/handler"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	registry "github.com/cs3org/go-cs3apis/cs3/storage/registry/v1beta1"
)

// UploadFinishedFunc is a callback function used in storage drivers to indicate that an upload has finished
type UploadFinishedFunc func(spaceOwner, owner *userpb.UserId, ref *provider.Reference)

// FS is the interface to implement access to the storage.
type FS interface {
	GetHome(ctx context.Context) (string, error)
	CreateHome(ctx context.Context) error
	CreateDir(ctx context.Context, ref *provider.Reference) error
	TouchFile(ctx context.Context, ref *provider.Reference, markprocessing bool) error
	Delete(ctx context.Context, ref *provider.Reference) error
	Move(ctx context.Context, oldRef, newRef *provider.Reference) error
	GetMD(ctx context.Context, ref *provider.Reference, mdKeys, fieldMask []string) (*provider.ResourceInfo, error)
	ListFolder(ctx context.Context, ref *provider.Reference, mdKeys, fieldMask []string) ([]*provider.ResourceInfo, error)
	InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error)
	Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, uploadFunc UploadFinishedFunc) (provider.ResourceInfo, error)
	Download(ctx context.Context, ref *provider.Reference) (io.ReadCloser, error)
	ListRevisions(ctx context.Context, ref *provider.Reference) ([]*provider.FileVersion, error)
	DownloadRevision(ctx context.Context, ref *provider.Reference, key string) (io.ReadCloser, error)
	RestoreRevision(ctx context.Context, ref *provider.Reference, key string) error
	ListRecycle(ctx context.Context, ref *provider.Reference, key, relativePath string) ([]*provider.RecycleItem, error)
	RestoreRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string, restoreRef *provider.Reference) error
	PurgeRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string) error
	EmptyRecycle(ctx context.Context, ref *provider.Reference) error
	GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error)
	AddGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error
	DenyGrant(ctx context.Context, ref *provider.Reference, g *provider.Grantee) error
	RemoveGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error
	UpdateGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error
	ListGrants(ctx context.Context, ref *provider.Reference) ([]*provider.Grant, error)
	GetQuota(ctx context.Context, ref *provider.Reference) ( /*TotalBytes*/ uint64 /*UsedBytes*/, uint64 /*RemainingBytes*/, uint64, error)
	CreateReference(ctx context.Context, path string, targetURI *url.URL) error
	Shutdown(ctx context.Context) error
	SetArbitraryMetadata(ctx context.Context, ref *provider.Reference, md *provider.ArbitraryMetadata) error
	UnsetArbitraryMetadata(ctx context.Context, ref *provider.Reference, keys []string) error
	SetLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error
	GetLock(ctx context.Context, ref *provider.Reference) (*provider.Lock, error)
	RefreshLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock, existingLockID string) error
	Unlock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error
	// ListStorageSpaces lists the spaces in the storage.
	// The unrestricted parameter can be used to list other user's spaces when
	// the user has the necessary permissions.
	ListStorageSpaces(ctx context.Context, filter []*provider.ListStorageSpacesRequest_Filter, unrestricted bool) ([]*provider.StorageSpace, error)
	CreateStorageSpace(ctx context.Context, req *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error)
	UpdateStorageSpace(ctx context.Context, req *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error)
	DeleteStorageSpace(ctx context.Context, req *provider.DeleteStorageSpaceRequest) error
}

// UploadsManager defines the interface for FS implementations that allow for managing uploads
type UploadsManager interface {
	ListUploads() ([]tusd.FileInfo, error)
	PurgeExpiredUploads(chan<- tusd.FileInfo) error
}

// Registry is the interface that storage registries implement
// for discovering storage providers
type Registry interface {
	// GetProvider returns the Address of the storage provider that should be used for the given space.
	// Use it to determine where to create a new storage space.
	GetProvider(ctx context.Context, space *provider.StorageSpace) (*registry.ProviderInfo, error)
	// ListProviders returns the storage providers that match the given filter
	ListProviders(ctx context.Context, filters map[string]string) ([]*registry.ProviderInfo, error)
}

// PathWrapper is the interface to implement for path transformations
type PathWrapper interface {
	Unwrap(ctx context.Context, rp string) (string, error)
	Wrap(ctx context.Context, rp string) (string, error)
}

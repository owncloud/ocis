// Copyright 2018-2023 CERN
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

package ocm

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	ocmpb "github.com/cs3org/go-cs3apis/cs3/sharing/ocm/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typepb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/rhttp/router"
	"github.com/cs3org/reva/v2/pkg/sharedconf"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/cs3org/reva/v2/pkg/utils/cfg"
	"github.com/studio-b12/gowebdav"
)

func init() {
	registry.Register("ocmreceived", New)
}

type driver struct {
	c       *config
	gateway gateway.GatewayAPIClient
}

type config struct {
	GatewaySVC string `mapstructure:"gatewaysvc"`
}

func (c *config) ApplyDefaults() {
	c.GatewaySVC = sharedconf.GetGatewaySVC(c.GatewaySVC)
}

// New creates an OCM storage driver.
func New(m map[string]interface{}, _ events.Stream) (storage.FS, error) {
	var c config
	if err := cfg.Decode(m, &c); err != nil {
		return nil, err
	}

	gateway, err := pool.GetGatewayServiceClient(c.GatewaySVC)
	if err != nil {
		return nil, err
	}

	d := &driver{
		c:       &c,
		gateway: gateway,
	}

	return d, nil
}

func shareInfoFromPath(path string) (*ocmpb.ShareId, string) {
	// the path is of the type /share_id[/rel_path]
	shareID, rel := router.ShiftPath(path)
	return &ocmpb.ShareId{OpaqueId: shareID}, rel
}

func shareInfoFromReference(ref *provider.Reference) (*ocmpb.ShareId, string) {
	if ref.ResourceId == nil {
		return shareInfoFromPath(ref.Path)
	}

	return &ocmpb.ShareId{OpaqueId: ref.ResourceId.OpaqueId}, ref.Path
}

func (d *driver) getWebDAVFromShare(ctx context.Context, shareID *ocmpb.ShareId) (*ocmpb.ReceivedShare, string, string, error) {
	// TODO: we may want to cache the share
	res, err := d.gateway.GetReceivedOCMShare(ctx, &ocmpb.GetReceivedOCMShareRequest{
		Ref: &ocmpb.ShareReference{
			Spec: &ocmpb.ShareReference_Id{
				Id: shareID,
			},
		},
	})
	if err != nil {
		return nil, "", "", err
	}

	if res.Status.Code != rpc.Code_CODE_OK {
		if res.Status.Code == rpc.Code_CODE_NOT_FOUND {
			return nil, "", "", errtypes.NotFound("share not found")
		}
		return nil, "", "", errtypes.InternalError(res.Status.Message)
	}

	dav, ok := getWebDAVProtocol(res.Share.Protocols)
	if !ok {
		return nil, "", "", errtypes.NotFound("share does not contain a WebDAV endpoint")
	}

	return res.Share, dav.Uri, dav.SharedSecret, nil
}

func getWebDAVProtocol(protocols []*ocmpb.Protocol) (*ocmpb.WebDAVProtocol, bool) {
	for _, p := range protocols {
		if dav, ok := p.Term.(*ocmpb.Protocol_WebdavOptions); ok {
			return dav.WebdavOptions, true
		}
	}
	return nil, false
}

func (d *driver) webdavClient(ctx context.Context, ref *provider.Reference) (*gowebdav.Client, *ocmpb.ReceivedShare, string, error) {
	id, rel := shareInfoFromReference(ref)

	share, endpoint, secret, err := d.getWebDAVFromShare(ctx, id)
	if err != nil {
		return nil, nil, "", err
	}

	endpoint, err = url.PathUnescape(endpoint)
	if err != nil {
		return nil, nil, "", err
	}

	// FIXME: it's still not clear from the OCM APIs how to use the shared secret
	// will use as a token in the bearer authentication as this is the reva implementation
	c := gowebdav.NewClient(endpoint, "", "")
	c.SetHeader("Authorization", "Bearer "+secret)

	return c, share, rel, nil
}

func (d *driver) CreateDir(ctx context.Context, ref *provider.Reference) error {
	client, _, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return err
	}
	return client.MkdirAll(rel, 0)
}

func (d *driver) Delete(ctx context.Context, ref *provider.Reference) error {
	client, _, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return err
	}
	return client.RemoveAll(rel)
}

func (d *driver) TouchFile(ctx context.Context, ref *provider.Reference, markprocessing bool, mtime string) error {
	client, _, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return err
	}
	return client.Write(rel, []byte{}, 0)
}

func (d *driver) Move(ctx context.Context, oldRef, newRef *provider.Reference) error {
	client, _, relOld, err := d.webdavClient(ctx, oldRef)
	if err != nil {
		return err
	}
	_, relNew := shareInfoFromReference(newRef)

	return client.Rename(relOld, relNew, false)
}

func getPathFromShareIDAndRelPath(shareID *ocmpb.ShareId, relPath string) string {
	return filepath.Join("/", shareID.OpaqueId, relPath)
}

func convertStatToResourceInfo(ref *provider.Reference, f fs.FileInfo, share *ocmpb.ReceivedShare, relPath string) *provider.ResourceInfo {
	t := provider.ResourceType_RESOURCE_TYPE_FILE
	if f.IsDir() {
		t = provider.ResourceType_RESOURCE_TYPE_CONTAINER
	}

	var name string
	if share.ResourceType == provider.ResourceType_RESOURCE_TYPE_FILE {
		name = share.Name
	} else {
		name = f.Name()
	}

	webdav, _ := getWebDAVProtocol(share.Protocols)

	return &provider.ResourceInfo{
		Type:     t,
		Id:       ref.ResourceId,
		MimeType: mime.Detect(f.IsDir(), f.Name()),
		Path:     relPath,
		Name:     name,
		Size:     uint64(f.Size()),
		Mtime: &typepb.Timestamp{
			Seconds: uint64(f.ModTime().Unix()),
		},
		Owner:         share.Creator,
		PermissionSet: webdav.Permissions.Permissions,
		Checksum: &provider.ResourceChecksum{
			Type: provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_INVALID,
		},
	}
}

func (d *driver) GetMD(ctx context.Context, ref *provider.Reference, _ []string, _ []string) (*provider.ResourceInfo, error) {
	client, share, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return nil, err
	}

	info, err := client.Stat(rel)
	if err != nil {
		if gowebdav.IsErrNotFound(err) {
			return nil, errtypes.NotFound(ref.GetPath())
		}
		return nil, err
	}

	return convertStatToResourceInfo(ref, info, share, rel), nil
}

func (d *driver) ListFolder(ctx context.Context, ref *provider.Reference, _ []string, _ []string) ([]*provider.ResourceInfo, error) {
	client, share, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return nil, err
	}

	list, err := client.ReadDir(rel)
	if err != nil {
		return nil, err
	}

	res := make([]*provider.ResourceInfo, 0, len(list))
	for _, r := range list {
		res = append(res, convertStatToResourceInfo(ref, r, share, utils.MakeRelativePath(filepath.Join(rel, r.Name()))))
	}
	return res, nil
}

func (d *driver) InitiateUpload(ctx context.Context, ref *provider.Reference, _ int64, _ map[string]string) (map[string]string, error) {
	shareID, rel := shareInfoFromReference(ref)
	p := getPathFromShareIDAndRelPath(shareID, rel)

	return map[string]string{
		"simple": p,
	}, nil
}

func (d *driver) Upload(ctx context.Context, req storage.UploadRequest, _ storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	client, _, rel, err := d.webdavClient(ctx, req.Ref)
	if err != nil {
		return provider.ResourceInfo{}, err
	}
	client.SetInterceptor(func(method string, rq *http.Request) {
		// Set the content length on the request struct directly instead of the header.
		// The content-length header gets reset by the golang http library before
		// sendind out the request, resulting in chunked encoding to be used which
		// breaks the quota checks in ocdav.
		if method == "PUT" {
			rq.ContentLength = req.Length
		}
	})

	return provider.ResourceInfo{}, client.WriteStream(rel, req.Body, 0)
}

func (d *driver) Download(ctx context.Context, ref *provider.Reference) (io.ReadCloser, error) {
	client, _, rel, err := d.webdavClient(ctx, ref)
	if err != nil {
		return nil, err
	}

	return client.ReadStream(rel)
}

func (d *driver) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	shareID, rel := shareInfoFromReference(&provider.Reference{
		ResourceId: id,
	})
	return getPathFromShareIDAndRelPath(shareID, rel), nil
}

func (d *driver) Shutdown(ctx context.Context) error {
	return nil
}

func (d *driver) CreateHome(ctx context.Context) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) GetHome(ctx context.Context) (string, error) {
	return "", errtypes.NotSupported("operation not supported")
}

func (d *driver) ListRevisions(ctx context.Context, ref *provider.Reference) ([]*provider.FileVersion, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) DownloadRevision(ctx context.Context, ref *provider.Reference, key string) (io.ReadCloser, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) RestoreRevision(ctx context.Context, ref *provider.Reference, key string) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) ListRecycle(ctx context.Context, ref *provider.Reference, key, relativePath string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) RestoreRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string, restoreRef *provider.Reference) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) PurgeRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) EmptyRecycle(ctx context.Context, ref *provider.Reference) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) AddGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) DenyGrant(ctx context.Context, ref *provider.Reference, g *provider.Grantee) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) RemoveGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) UpdateGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) ListGrants(ctx context.Context, ref *provider.Reference) ([]*provider.Grant, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) GetQuota(ctx context.Context, ref *provider.Reference) ( /*TotalBytes*/ uint64 /*UsedBytes*/, uint64, uint64, error) {
	return 0, 0, 0, errtypes.NotSupported("operation not supported")
}

func (d *driver) CreateReference(ctx context.Context, path string, targetURI *url.URL) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) SetArbitraryMetadata(ctx context.Context, ref *provider.Reference, md *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) UnsetArbitraryMetadata(ctx context.Context, ref *provider.Reference, keys []string) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) SetLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) GetLock(ctx context.Context, ref *provider.Reference) (*provider.Lock, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) RefreshLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock, existingLockID string) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) Unlock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error {
	return errtypes.NotSupported("operation not supported")
}

func (d *driver) ListStorageSpaces(ctx context.Context, filters []*provider.ListStorageSpacesRequest_Filter, _ bool) ([]*provider.StorageSpace, error) {
	spaceTypes := map[string]struct{}{}
	var exists = struct{}{}
	appendTypes := []string{}
	for _, f := range filters {
		if f.Type == provider.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE {
			spaceType := f.GetSpaceType()
			if spaceType == "+mountpoint" {
				appendTypes = append(appendTypes, strings.TrimPrefix(spaceType, "+"))
				continue
			}
			spaceTypes[spaceType] = exists
		}
	}

	lrsRes, err := d.gateway.ListReceivedOCMShares(ctx, &ocmpb.ListReceivedOCMSharesRequest{})
	if err != nil {
		return nil, err
	}

	if len(spaceTypes) == 0 {
		spaceTypes["mountpoint"] = exists
	}
	for _, s := range appendTypes {
		spaceTypes[s] = exists
	}

	spaces := []*provider.StorageSpace{}
	for k := range spaceTypes {
		if k == "mountpoint" {
			for _, share := range lrsRes.Shares {
				root := &provider.ResourceId{
					StorageId: utils.PublicStorageProviderID,
					SpaceId:   share.Id.OpaqueId,
					OpaqueId:  share.Id.OpaqueId,
				}
				space := &provider.StorageSpace{
					Id: &provider.StorageSpaceId{
						OpaqueId: storagespace.FormatResourceID(*root),
					},
					SpaceType: "mountpoint",
					Owner: &userv1beta1.User{
						Id: share.Grantee.GetUserId(),
					},
					Root: root,
				}

				spaces = append(spaces, space)
			}
		}
	}

	return spaces, nil
}

func (d *driver) CreateStorageSpace(ctx context.Context, req *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) UpdateStorageSpace(ctx context.Context, req *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("operation not supported")
}

func (d *driver) DeleteStorageSpace(ctx context.Context, req *provider.DeleteStorageSpaceRequest) error {
	return errtypes.NotSupported("operation not supported")
}

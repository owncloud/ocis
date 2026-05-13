// vendor/github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kiteworks.go
package kiteworks

import (
	"context"
	"io"
	"math"
	"net/url"
	"strings"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
)

func init() {
	registry.Register("kiteworks", New)
}

// Config holds the driver configuration decoded from reva mapstructure
type Config struct {
	Endpoint  string `mapstructure:"endpoint"`
	Insecure  bool   `mapstructure:"insecure"`
	ChunkSize int64  `mapstructure:"chunk_size"`
}

// Driver implements storage.FS backed by the Kiteworks REST API
type Driver struct {
	cfg *Config
}

// New creates a new kiteworks storage driver
func New(m map[string]interface{}, _ events.Stream, _ *zerolog.Logger) (storage.FS, error) {
	cfg := &Config{}
	if err := mapstructure.Decode(m, cfg); err != nil {
		return nil, errors.Wrap(err, "kiteworks: error decoding config")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("kiteworks: 'endpoint' must be set")
	}
	if cfg.ChunkSize <= 0 {
		cfg.ChunkSize = defaultChunkSize
	}
	return &Driver{cfg: cfg}, nil
}

func (d *Driver) client(ctx context.Context) (*Client, error) {
	token, ok := ctxpkg.ContextGetToken(ctx)
	if !ok || token == "" {
		return nil, errtypes.PermissionDenied("kiteworks: no token in context")
	}
	return NewClient(d.cfg.Endpoint, token, d.cfg.Insecure), nil
}

// fileInfoToResourceInfo converts a Kiteworks FileInfo to a CS3 ResourceInfo
func fileInfoToResourceInfo(fi *FileInfo) *provider.ResourceInfo {
	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: "kiteworks",
			OpaqueId:  fi.ID,
		},
		Path: fi.Name,
		Type: provider.ResourceType_RESOURCE_TYPE_FILE,
		Size: uint64(fi.Size),
	}
	if fi.IsDir() {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_CONTAINER
	}
	if fi.Modified != nil {
		ri.Mtime = &typespb.Timestamp{
			Seconds: uint64(fi.Modified.Unix()),
		}
	}
	if fi.FingerPrints != nil {
		for _, fp := range fi.FingerPrints.FingerPrints {
			switch fp.Algo {
			case "sha1":
				ri.Checksum = &provider.ResourceChecksum{
					Type: provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_SHA1,
					Sum:  fp.Hash,
				}
			case "md5":
				ri.Checksum = &provider.ResourceChecksum{
					Type: provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_MD5,
					Sum:  fp.Hash,
				}
			case "adler32":
				ri.Checksum = &provider.ResourceChecksum{
					Type: provider.ResourceChecksumType_RESOURCE_CHECKSUM_TYPE_ADLER32,
					Sum:  fp.Hash,
				}
			}
		}
	}
	return ri
}

// spaceFromFileInfo converts a top-level Kiteworks folder to a CS3 StorageSpace
func spaceFromFileInfo(fi *FileInfo) *provider.StorageSpace {
	spaceType := "project"
	if fi.IsSharedWithUser() {
		spaceType = "mountpoint"
	}
	sp := &provider.StorageSpace{
		Id: &provider.StorageSpaceId{
			OpaqueId: fi.ID,
		},
		Root: &provider.ResourceId{
			StorageId: "kiteworks",
			OpaqueId:  fi.ID,
		},
		Name:      fi.Name,
		SpaceType: spaceType,
	}
	if fi.Modified != nil {
		sp.Mtime = &typespb.Timestamp{
			Seconds: uint64(fi.Modified.Unix()),
		}
	}
	return sp
}

// Shutdown implements storage.FS
func (d *Driver) Shutdown(_ context.Context) error { return nil }

// ListStorageSpaces implements storage.FS — returns each top-level Kiteworks folder as a space
func (d *Driver) ListStorageSpaces(ctx context.Context, _ []*provider.ListStorageSpacesRequest_Filter, _ bool) ([]*provider.StorageSpace, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	folders, err := c.GetTopFolders()
	if err != nil {
		return nil, err
	}
	spaces := make([]*provider.StorageSpace, 0, len(folders))
	for i := range folders {
		spaces = append(spaces, spaceFromFileInfo(&folders[i]))
	}
	return spaces, nil
}

// GetQuota implements storage.FS
func (d *Driver) GetQuota(ctx context.Context, _ *provider.Reference) (uint64, uint64, uint64, error) {
	c, err := d.client(ctx)
	if err != nil {
		return 0, 0, 0, err
	}
	u, err := c.GetMe()
	if err != nil {
		return 0, 0, 0, err
	}
	total := uint64(u.Quota.Allowed)
	used := uint64(u.Quota.Used)
	var remaining uint64
	if total > used {
		remaining = total - used
	}
	return total, used, remaining, nil
}

// GetMD implements storage.FS
func (d *Driver) GetMD(ctx context.Context, ref *provider.Reference, _, _ []string) (*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" && ref.GetPath() == "" {
		return nil, errtypes.NotFound("kiteworks: reference has no id or path")
	}
	if id == "" && ref.GetPath() != "" {
		// path-based lookup via search
		results, err := c.Search(ref.GetPath())
		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return nil, errtypes.NotFound(ref.GetPath())
		}
		return fileInfoToResourceInfo(&results[0]), nil
	}
	// Try folder first, fall back to file
	fi, err := c.GetFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return fileInfoToResourceInfo(fi), nil
}

// ListFolder implements storage.FS
func (d *Driver) ListFolder(ctx context.Context, ref *provider.Reference, _, _ []string) ([]*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" {
		return nil, errtypes.NotFound("kiteworks: reference has no id")
	}
	dir, err := c.ListFolder(id)
	if err != nil {
		return nil, err
	}
	var infos []*provider.ResourceInfo
	for i := range dir.Folders {
		infos = append(infos, fileInfoToResourceInfo(&dir.Folders[i]))
	}
	for i := range dir.Files {
		infos = append(infos, fileInfoToResourceInfo(&dir.Files[i]))
	}
	return infos, nil
}

// Download implements storage.FS
func (d *Driver) Download(ctx context.Context, ref *provider.Reference, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" {
		return nil, nil, errtypes.NotFound("kiteworks: reference has no id")
	}
	fi, err := c.GetFile(id)
	if err != nil {
		return nil, nil, err
	}
	ri := fileInfoToResourceInfo(fi)
	if openReaderFunc != nil && !openReaderFunc(ri) {
		return ri, nil, nil
	}
	resp, err := c.DownloadFile(id, "")
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, nil, &ClientError{StatusCode: resp.StatusCode, Body: body}
	}
	return ri, resp.Body, nil
}

// GetPathByID implements storage.FS
func (d *Driver) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	if id == nil || id.OpaqueId == "" {
		return "", errtypes.NotFound("kiteworks: missing resource id")
	}
	c, err := d.client(ctx)
	if err != nil {
		return "", err
	}
	// Try folder first, fall back to file
	fi, err := c.GetFolder(id.OpaqueId)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id.OpaqueId)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return fi.Path, nil
}

// --- Stubbed / not-supported operations ---

func (d *Driver) CreateStorageSpace(_ context.Context, _ *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: CreateStorageSpace")
}
func (d *Driver) UpdateStorageSpace(_ context.Context, _ *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: UpdateStorageSpace")
}
func (d *Driver) DeleteStorageSpace(_ context.Context, _ *provider.DeleteStorageSpaceRequest) error {
	return errtypes.NotSupported("kiteworks: DeleteStorageSpace")
}
func (d *Driver) CreateHome(_ context.Context) error {
	return errtypes.NotSupported("kiteworks: CreateHome")
}
func (d *Driver) GetHome(_ context.Context) (string, error) {
	return "", errtypes.NotSupported("kiteworks: GetHome")
}
func (d *Driver) CreateReference(_ context.Context, _ string, _ *url.URL) error {
	return errtypes.NotSupported("kiteworks: CreateReference")
}
func (d *Driver) CreateDir(ctx context.Context, ref *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	name := ref.GetPath()
	if parentID == "" || name == "" {
		return errtypes.NotFound("kiteworks: CreateDir requires a parent ID and name")
	}
	_, err = c.CreateFolder(parentID, name)
	return err
}
func (d *Driver) TouchFile(ctx context.Context, ref *provider.Reference, _ bool, _ string) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	name := ref.GetPath()
	if parentID == "" || name == "" {
		return errtypes.NotFound("kiteworks: TouchFile requires a parent ID and name")
	}
	_, err = uploadFile(c, parentID, name, 0, strings.NewReader(""), d.cfg.ChunkSize)
	return err
}
func (d *Driver) Delete(ctx context.Context, ref *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" {
		return errtypes.NotFound("kiteworks: reference has no id")
	}
	err = c.DeleteFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			return c.DeleteFile(id)
		}
		return err
	}
	return nil
}
func (d *Driver) Move(ctx context.Context, oldRef, newRef *provider.Reference) error {
	c, err := d.client(ctx)
	if err != nil {
		return err
	}
	sourceID := oldRef.GetResourceId().GetOpaqueId()
	if sourceID == "" {
		return errtypes.NotFound("kiteworks: source reference has no id")
	}
	destFolderID := newRef.GetResourceId().GetOpaqueId()
	if destFolderID == "" {
		// rename: try file first, fall back to folder
		err = c.RenameFile(sourceID, newRef.GetPath(), false)
		if err != nil {
			if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
				return c.RenameFolder(sourceID, newRef.GetPath())
			}
			return err
		}
		return nil
	}
	return c.MoveResource(sourceID, destFolderID, false)
}
func (d *Driver) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	parentID := ref.GetResourceId().GetOpaqueId()
	filename := metadata["filename"]
	if filename == "" {
		filename = ref.GetPath()
	}
	numChunks := 1
	if uploadLength > 0 {
		numChunks = int(math.Ceil(float64(uploadLength) / float64(d.cfg.ChunkSize)))
	}
	result, err := c.InitiateUpload(parentID, filename, uploadLength, numChunks)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"uploadID":  result.ID,
		"uploadURI": result.URI,
		"filename":  filename,
		"parentID":  parentID,
	}, nil
}
func (d *Driver) Upload(ctx context.Context, req storage.UploadRequest, uploadFunc storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	parentID := req.Ref.GetResourceId().GetOpaqueId()
	filename := req.Ref.GetPath()
	fi, err := uploadFile(c, parentID, filename, req.Length, req.Body, d.cfg.ChunkSize)
	if err != nil {
		return nil, err
	}
	ri := fileInfoToResourceInfo(fi)
	if uploadFunc != nil {
		u, ok := ctxpkg.ContextGetUser(ctx)
		if !ok {
			return nil, errtypes.PermissionDenied("kiteworks: no user in context")
		}
		uploadFunc(u.GetId(), u.GetId(), req.Ref)
	}
	return ri, nil
}
func (d *Driver) ListRevisions(_ context.Context, _ *provider.Reference) ([]*provider.FileVersion, error) {
	return nil, errtypes.NotSupported("kiteworks: ListRevisions")
}
func (d *Driver) DownloadRevision(_ context.Context, _ *provider.Reference, _ string, _ func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	return nil, nil, errtypes.NotSupported("kiteworks: DownloadRevision")
}
func (d *Driver) RestoreRevision(_ context.Context, _ *provider.Reference, _ string) error {
	return errtypes.NotSupported("kiteworks: RestoreRevision")
}
func (d *Driver) ListRecycle(_ context.Context, _ *provider.Reference, _, _ string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("kiteworks: ListRecycle")
}
func (d *Driver) RestoreRecycleItem(_ context.Context, _ *provider.Reference, _, _ string, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: RestoreRecycleItem")
}
func (d *Driver) PurgeRecycleItem(_ context.Context, _ *provider.Reference, _, _ string) error {
	return errtypes.NotSupported("kiteworks: PurgeRecycleItem")
}
func (d *Driver) EmptyRecycle(_ context.Context, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: EmptyRecycle")
}
func (d *Driver) AddGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: AddGrant — permission mapping not yet implemented")
}
func (d *Driver) DenyGrant(_ context.Context, _ *provider.Reference, _ *provider.Grantee) error {
	return errtypes.NotSupported("kiteworks: DenyGrant")
}
func (d *Driver) RemoveGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: RemoveGrant — permission mapping not yet implemented")
}
func (d *Driver) UpdateGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: UpdateGrant — permission mapping not yet implemented")
}
func (d *Driver) ListGrants(ctx context.Context, ref *provider.Reference) ([]*provider.Grant, error) {
	c, err := d.client(ctx)
	if err != nil {
		return nil, err
	}
	id := ref.GetResourceId().GetOpaqueId()
	if id == "" {
		return nil, errtypes.NotFound("kiteworks: reference has no id")
	}
	fi, err := c.GetFolder(id)
	if err != nil {
		if ce, ok := err.(*ClientError); ok && ce.StatusCode == 404 {
			fi, err = c.GetFile(id)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var grants []*provider.Grant
	for _, perm := range fi.Permissions {
		if !perm.Allowed {
			continue
		}
		grants = append(grants, &provider.Grant{
			Grantee: &provider.Grantee{
				Type: provider.GranteeType_GRANTEE_TYPE_USER,
				Id: &provider.Grantee_UserId{
					UserId: &userpb.UserId{OpaqueId: perm.Name},
				},
			},
			Permissions: &provider.ResourcePermissions{
				GetPath:              true,
				InitiateFileDownload: true,
				InitiateFileUpload:   true,
				ListContainer:        true,
				Stat:                 true,
			},
		})
	}
	return grants, nil
}
func (d *Driver) SetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("kiteworks: SetArbitraryMetadata")
}
func (d *Driver) UnsetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ []string) error {
	return errtypes.NotSupported("kiteworks: UnsetArbitraryMetadata")
}
func (d *Driver) GetLock(_ context.Context, _ *provider.Reference) (*provider.Lock, error) {
	return nil, errtypes.NotSupported("kiteworks: GetLock")
}
func (d *Driver) SetLock(_ context.Context, _ *provider.Reference, _ *provider.Lock) error {
	return errtypes.NotSupported("kiteworks: SetLock")
}
func (d *Driver) RefreshLock(_ context.Context, _ *provider.Reference, _ *provider.Lock, _ string) error {
	return errtypes.NotSupported("kiteworks: RefreshLock")
}
func (d *Driver) Unlock(_ context.Context, _ *provider.Reference, _ *provider.Lock) error {
	return errtypes.NotSupported("kiteworks: Unlock")
}

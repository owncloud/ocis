package kiteworks

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"

	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/kiteworks/kwlib"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/reva/v2/pkg/utils"
)

func init() {
	registry.Register("kiteworks", New)
}

// Config holds the driver configuration.
type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	APIToken string `mapstructure:"api_token"`
	Insecure bool   `mapstructure:"insecure"`
	MountID  string `mapstructure:"mount_id"`
}

// Driver implements storage.FS against a Kiteworks box (read-only).
type Driver struct {
	factory   *kwlib.APIClientFactory
	apiToken  string
	storageID string
	log       zerolog.Logger
}

// New returns a read-only Kiteworks storage driver.
func New(m map[string]interface{}, _ events.Stream, log *zerolog.Logger) (storage.FS, error) {
	c := &Config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	c.Endpoint = strings.TrimRight(c.Endpoint, "/")

	storageID := c.MountID
	if storageID == "" {
		storageID = "kiteworks"
	}

	l := zerolog.Nop()
	if log != nil {
		l = *log
	}

	return &Driver{
		factory:   kwlib.NewClientFactory(c.Endpoint, "reva-kiteworks/1.0", c.Insecure),
		apiToken:  c.APIToken,
		storageID: storageID,
		log:       l,
	}, nil
}

func (d *Driver) client(ctx context.Context) *kwlib.APIClient {
	token := d.apiToken
	if token == "" {
		token, _ = ctxpkg.ContextGetToken(ctx)
	}
	return d.factory.Build("", "", "", token, &d.log)
}

// toResourceInfo converts a kwlib.FileInfo to a CS3 ResourceInfo.
// spaceRootPath is the absolute KW path of the space root folder; it is stripped
// from fi.Path to produce a space-relative path with a leading "/".
func (d *Driver) toResourceInfo(fi *kwlib.FileInfo, spaceID, spaceRootPath string) *provider.ResourceInfo {
	relPath := strings.TrimPrefix(fi.Path, spaceRootPath)
	if !strings.HasPrefix(relPath, "/") {
		relPath = "/" + relPath
	}
	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: d.storageID,
			SpaceId:   spaceID,
			OpaqueId:  fi.ID,
		},
		Name:  fi.Name,
		Path:  relPath,
		Etag:  fi.ETag(),
		Mtime: utils.TimeToTS(fi.MTime()),
		PermissionSet: &provider.ResourcePermissions{
			Stat:                 true,
			GetPath:              true,
			InitiateFileDownload: fi.IsFile(),
			ListContainer:        fi.IsDir(),
		},
	}
	if fi.IsDir() {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_CONTAINER
		ri.MimeType = "httpd/unix-directory"
	} else {
		ri.Type = provider.ResourceType_RESOURCE_TYPE_FILE
		if fi.Size != nil {
			ri.Size = uint64(*fi.Size)
		}
	}
	return ri
}

// --- Read methods ---

func (d *Driver) Shutdown(_ context.Context) error { return nil }

func (d *Driver) ListStorageSpaces(ctx context.Context, _ []*provider.ListStorageSpacesRequest_Filter, _ bool) ([]*provider.StorageSpace, error) {
	c := d.client(ctx)
	dirs, err := c.GetTopFolders()
	if err != nil {
		return nil, err
	}

	spaces := make([]*provider.StorageSpace, 0, len(dirs.Data))
	for i := range dirs.Data {
		fi := &dirs.Data[i]
		spaces = append(spaces, &provider.StorageSpace{
			Id:        &provider.StorageSpaceId{OpaqueId: fi.ID},
			Name:      fi.Name,
			SpaceType: "project",
			Root: &provider.ResourceId{
				StorageId: d.storageID,
				SpaceId:   fi.ID,
				OpaqueId:  fi.ID,
			},
			RootInfo: d.toResourceInfo(fi, fi.ID, fi.Path),
			Mtime:    utils.TimeToTS(fi.MTime()),
			Opaque:   utils.AppendPlainToOpaque(nil, "spaceAlias", "project/"+fi.Name),
		})
	}
	return spaces, nil
}

// resolveRef walks a CS3 reference to the target KW node ID and space ID.
// If ref.Path is non-empty it resolves each component through ListFolderContents.
func (d *Driver) resolveRef(ctx context.Context, ref *provider.Reference) (nodeID, spaceID string, err error) {
	spaceID = ref.GetResourceId().GetSpaceId()
	nodeID = ref.GetResourceId().GetOpaqueId()
	if nodeID == "" {
		nodeID = spaceID
	}

	relPath := strings.Trim(strings.TrimPrefix(ref.GetPath(), "./"), "/.")
	if relPath == "" {
		return nodeID, spaceID, nil
	}

	c := d.client(ctx)
	for _, part := range strings.Split(relPath, "/") {
		if part == "" {
			continue
		}
		children, err := c.ListFolderContents(nodeID)
		if err != nil {
			return "", "", err
		}
		var found bool
		for i := range children {
			if children[i].Name == part {
				nodeID = children[i].ID
				found = true
				break
			}
		}
		if !found {
			return "", "", errtypes.NotFound(part)
		}
	}
	return nodeID, spaceID, nil
}

// spaceRootPath fetches the absolute KW path of the space root folder.
// Used by callers that need to convert absolute KW paths to space-relative paths.
func (d *Driver) spaceRootPath(ctx context.Context, spaceID string) (string, error) {
	root, err := d.client(ctx).GetFolderByID(spaceID)
	if err != nil {
		return "", err
	}
	return root.Path, nil
}

func (d *Driver) GetMD(ctx context.Context, ref *provider.Reference, _, _ []string) (*provider.ResourceInfo, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, err
	}
	return d.nodeMD(ctx, nodeID, spaceID, rootPath)
}

// nodeMD fetches metadata for a node, trying folder first then file.
func (d *Driver) nodeMD(ctx context.Context, nodeID, spaceID, spaceRootPath string) (*provider.ResourceInfo, error) {
	c := d.client(ctx)
	fi, err := c.GetFolderByID(nodeID)
	if err == nil {
		return d.toResourceInfo(fi, spaceID, spaceRootPath), nil
	}
	var ce *kwlib.ClientError
	if !errors.As(err, &ce) || ce.StatusCode != http.StatusNotFound {
		return nil, err
	}
	fi, err = c.GetFileByID(nodeID)
	if err != nil {
		var fce *kwlib.ClientError
		if errors.As(err, &fce) && fce.StatusCode == http.StatusNotFound {
			return nil, errtypes.NotFound(nodeID)
		}
		return nil, err
	}
	return d.toResourceInfo(fi, spaceID, spaceRootPath), nil
}

func (d *Driver) ListFolder(ctx context.Context, ref *provider.Reference, _, _ []string) ([]*provider.ResourceInfo, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, err
	}

	items, err := d.client(ctx).ListFolderContents(nodeID)
	if err != nil {
		return nil, err
	}

	infos := make([]*provider.ResourceInfo, 0, len(items))
	for i := range items {
		infos = append(infos, d.toResourceInfo(&items[i], spaceID, rootPath))
	}
	return infos, nil
}

func (d *Driver) Download(ctx context.Context, ref *provider.Reference, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	nodeID, spaceID, err := d.resolveRef(ctx, ref)
	if err != nil {
		return nil, nil, err
	}
	rootPath, err := d.spaceRootPath(ctx, spaceID)
	if err != nil {
		return nil, nil, err
	}

	ri, err := d.nodeMD(ctx, nodeID, spaceID, rootPath)
	if err != nil {
		return nil, nil, err
	}

	if !openReaderFunc(ri) {
		return ri, nil, nil
	}

	resp, err := d.client(ctx).GetFileContents(nodeID, "")
	if err != nil {
		return nil, nil, err
	}
	return ri, resp.Body, nil
}

func (d *Driver) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	rootPath, err := d.spaceRootPath(ctx, id.GetSpaceId())
	if err != nil {
		return "", err
	}
	ri, err := d.nodeMD(ctx, id.GetOpaqueId(), id.GetSpaceId(), rootPath)
	if err != nil {
		return "", err
	}
	return ri.Path, nil
}

func (d *Driver) ListGrants(_ context.Context, _ *provider.Reference) ([]*provider.Grant, error) {
	return []*provider.Grant{}, nil
}

func (d *Driver) GetQuota(ctx context.Context, _ *provider.Reference) (uint64, uint64, uint64, error) {
	q, err := d.client(ctx).GetQuotaInfo()
	if err != nil {
		// Non-fatal for read-only driver; return zero quota rather than failing.
		return 0, 0, 0, nil
	}
	total := uint64(q.FolderQuotaAllowed)
	used := uint64(q.FolderQuotaUsed)
	var remaining uint64
	if total > used {
		remaining = total - used
	}
	return total, used, remaining, nil
}

func (d *Driver) GetLock(_ context.Context, _ *provider.Reference) (*provider.Lock, error) {
	return nil, nil
}

func (d *Driver) ListRevisions(_ context.Context, _ *provider.Reference) ([]*provider.FileVersion, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) DownloadRevision(_ context.Context, _ *provider.Reference, _ string, _ func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	return nil, nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) ListRecycle(_ context.Context, _ *provider.Reference, _, _ string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

// --- Write methods: all return NotSupported ---

func (d *Driver) CreateReference(_ context.Context, _ string, _ *url.URL) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CreateDir(_ context.Context, _ *provider.Reference) (*storage.CreateDirResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) TouchFile(_ context.Context, _ *provider.Reference, _ bool, _ string) (*storage.TouchFileResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) Delete(_ context.Context, _ *provider.Reference) (*storage.DeleteResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) Move(_ context.Context, _, _ *provider.Reference) (*storage.MoveResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) InitiateUpload(_ context.Context, _ *provider.Reference, _ int64, _ map[string]string) (map[string]string, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) Upload(_ context.Context, _ storage.UploadRequest, _ storage.UploadFinishedFunc) (*provider.ResourceInfo, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) RestoreRevision(_ context.Context, _ *provider.Reference, _ string) (*storage.RestoreRevisionResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) RestoreRecycleItem(_ context.Context, _ *provider.Reference, _, _ string, _ *provider.Reference) (*storage.RestoreRecycleItemResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) PurgeRecycleItem(_ context.Context, _ *provider.Reference, _, _ string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) EmptyRecycle(_ context.Context, _ *provider.Reference) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) AddGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) DenyGrant(_ context.Context, _ *provider.Reference, _ *provider.Grantee) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) RemoveGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UpdateGrant(_ context.Context, _ *provider.Reference, _ *provider.Grant) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) SetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UnsetArbitraryMetadata(_ context.Context, _ *provider.Reference, _ []string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) SetLock(_ context.Context, _ *provider.Reference, _ *provider.Lock) (*storage.SetLockResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) RefreshLock(_ context.Context, _ *provider.Reference, _ *provider.Lock, _ string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) Unlock(_ context.Context, _ *provider.Reference, _ *provider.Lock) (*storage.UnlockResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CreateStorageSpace(_ context.Context, _ *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) UpdateStorageSpace(_ context.Context, _ *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) DeleteStorageSpace(_ context.Context, _ *provider.DeleteStorageSpaceRequest) (*storage.DeleteStorageSpaceResult, error) {
	return nil, errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CreateHome(_ context.Context) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) GetHome(_ context.Context) (string, error) {
	return "", errtypes.NotSupported("kiteworks: read-only driver")
}


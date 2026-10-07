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

const storageID = "kiteworks"

// Config holds the driver configuration.
type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Insecure bool   `mapstructure:"insecure"`
}

// Driver implements storage.FS against a Kiteworks box (read-only).
type Driver struct {
	factory *kwlib.APIClientFactory
	log     zerolog.Logger
}

// New returns a read-only Kiteworks storage driver.
func New(m map[string]interface{}, _ events.Stream, log *zerolog.Logger) (storage.FS, error) {
	c := &Config{}
	if err := mapstructure.Decode(m, c); err != nil {
		return nil, err
	}
	c.Endpoint = strings.TrimRight(c.Endpoint, "/")

	l := zerolog.Nop()
	if log != nil {
		l = *log
	}

	return &Driver{
		factory: kwlib.NewClientFactory(c.Endpoint, "reva-kiteworks/1.0", c.Insecure),
		log:     l,
	}, nil
}

func (d *Driver) client(ctx context.Context) *kwlib.APIClient {
	token, _ := ctxpkg.ContextGetToken(ctx)
	return d.factory.Build("", "", "", token, &d.log)
}

// toResourceInfo converts a kwlib.FileInfo to a CS3 ResourceInfo.
func (d *Driver) toResourceInfo(fi *kwlib.FileInfo, spaceID string) *provider.ResourceInfo {
	ri := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			StorageId: storageID,
			SpaceId:   spaceID,
			OpaqueId:  fi.ID,
		},
		Name:  fi.Name,
		Path:  fi.Path,
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
				StorageId: storageID,
				SpaceId:   fi.ID,
				OpaqueId:  fi.ID,
			},
			RootInfo: d.toResourceInfo(fi, fi.ID),
			Mtime:    utils.TimeToTS(fi.MTime()),
			Opaque:   utils.AppendPlainToOpaque(nil, "spaceAlias", "project/"+fi.Name),
		})
	}
	return spaces, nil
}

func resolveNodeID(ref *provider.Reference) (nodeID, spaceID string) {
	spaceID = ref.GetResourceId().GetSpaceId()
	nodeID = ref.GetResourceId().GetOpaqueId()
	if nodeID == "" {
		nodeID = spaceID
	}
	return
}

func (d *Driver) GetMD(ctx context.Context, ref *provider.Reference, _, _ []string) (*provider.ResourceInfo, error) {
	nodeID, spaceID := resolveNodeID(ref)
	return d.nodeMD(ctx, nodeID, spaceID)
}

// nodeMD fetches metadata for a node, trying folder first then file.
func (d *Driver) nodeMD(ctx context.Context, nodeID, spaceID string) (*provider.ResourceInfo, error) {
	c := d.client(ctx)
	fi, err := c.GetFolderByID(nodeID)
	if err == nil {
		return d.toResourceInfo(fi, spaceID), nil
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
	return d.toResourceInfo(fi, spaceID), nil
}

func (d *Driver) ListFolder(ctx context.Context, ref *provider.Reference, _, _ []string) ([]*provider.ResourceInfo, error) {
	nodeID, spaceID := resolveNodeID(ref)

	items, err := d.client(ctx).ListFolderContents(nodeID)
	if err != nil {
		return nil, err
	}

	infos := make([]*provider.ResourceInfo, 0, len(items))
	for i := range items {
		infos = append(infos, d.toResourceInfo(&items[i], spaceID))
	}
	return infos, nil
}

func (d *Driver) Download(ctx context.Context, ref *provider.Reference, openReaderFunc func(*provider.ResourceInfo) bool) (*provider.ResourceInfo, io.ReadCloser, error) {
	nodeID, spaceID := resolveNodeID(ref)

	ri, err := d.nodeMD(ctx, nodeID, spaceID)
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
	ri, err := d.nodeMD(ctx, id.GetOpaqueId(), id.GetSpaceId())
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
		return 0, 0, 0, err
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

func (d *Driver) MarkProcessing(_ context.Context, _ *provider.Reference, _ bool, _ string) error {
	return errtypes.NotSupported("kiteworks: read-only driver")
}

func (d *Driver) CommitUpload(_ context.Context, _ *provider.Reference, _ storage.UploadSource) (*provider.ResourceInfo, error) {
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

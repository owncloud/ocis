package cs3

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	idxerrs "github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
	metadatastorage "github.com/owncloud/ocis/ocis-pkg/metadata_storage"
	"google.golang.org/grpc/metadata"
)

func init() {
	registry.IndexConstructorRegistry["cs3"]["non_unique"] = NewNonUniqueIndexWithOptions
}

// NonUnique are fields for an index of type non_unique.
type NonUnique struct {
	caseInsensitive bool
	indexBy         string
	typeName        string
	filesDir        string
	indexBaseDir    string
	indexRootDir    string

	tokenManager    token.Manager
	storageProvider provider.ProviderAPIClient
	metadataStorage *metadatastorage.MetadataStorage

	cs3conf *Config
}

// NewNonUniqueIndexWithOptions instantiates a new NonUniqueIndex instance.
// /tmp/ocis/accounts/index.cs3/Pets/Bro*
// ├── Brown/
// │   └── rebef-123 -> /tmp/testfiles-395764020/pets/rebef-123
// ├── Green/
// │    ├── goefe-789 -> /tmp/testfiles-395764020/pets/goefe-789
// │    └── xadaf-189 -> /tmp/testfiles-395764020/pets/xadaf-189
// └── White/
//     └── wefwe-456 -> /tmp/testfiles-395764020/pets/wefwe-456
func NewNonUniqueIndexWithOptions(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	return &NonUnique{
		caseInsensitive: opts.CaseInsensitive,
		indexBy:         opts.IndexBy,
		typeName:        opts.TypeName,
		filesDir:        opts.FilesDir,
		indexBaseDir:    path.Join(opts.DataDir, "index.cs3"),
		indexRootDir:    path.Join(path.Join(opts.DataDir, "index.cs3"), strings.Join([]string{"non_unique", opts.TypeName, opts.IndexBy}, ".")),
		cs3conf: &Config{
			ProviderAddr: opts.ProviderAddr,
			JWTSecret:    opts.JWTSecret,
			ServiceUser:  opts.ServiceUser,
		},
	}
}

// Init initializes a non_unique index.
func (idx *NonUnique) Init() error {
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": idx.cs3conf.JWTSecret,
	})
	if err != nil {
		return err
	}
	idx.tokenManager = tokenManager

	client, err := pool.GetStorageProviderServiceClient(idx.cs3conf.ProviderAddr)
	if err != nil {
		return err
	}
	idx.storageProvider = client

	m, err := metadatastorage.NewMetadataStorage(idx.cs3conf.ProviderAddr)
	if err != nil {
		return err
	}
	idx.metadataStorage = &m

	if err := idx.makeDirIfNotExists(idx.indexBaseDir); err != nil {
		return err
	}

	if err := idx.makeDirIfNotExists(idx.indexRootDir); err != nil {
		return err
	}

	return nil
}

// Lookup exact lookup by value.
func (idx *NonUnique) Lookup(v string) ([]string, error) {
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	var matches = make([]string, 0)
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", idx.indexRootDir, v),
		},
	})

	if err != nil {
		return nil, err
	}

	for _, info := range res.Infos {
		matches = append(matches, path.Base(info.Path))
	}

	return matches, nil
}

// Add a new value to the index.
func (idx *NonUnique) Add(id, v string) (string, error) {
	if v == "" {
		return "", nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}

	newName := path.Join(idx.indexRootDir, v)
	if err := idx.makeDirIfNotExists(newName); err != nil {
		return "", err
	}

	if err := idx.createSymlink(id, path.Join(newName, id)); err != nil {
		if os.IsExist(err) {
			return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return "", err
	}

	return newName, nil
}

// Remove a value v from an index.
func (idx *NonUnique) Remove(id string, v string) error {
	if v == "" {
		return nil
	}
	if idx.caseInsensitive {
		v = strings.ToLower(v)
	}
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	deletePath := path.Join("/", idx.indexRootDir, v, id)
	resp, err := idx.storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Path: deletePath,
		},
	})

	if err != nil {
		return err
	}

	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		return &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
	}

	toStat := path.Join("/", idx.indexRootDir, v)
	lcResp, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: toStat,
		},
	})
	if err != nil {
		return err
	}

	if len(lcResp.Infos) == 0 {
		deletePath = path.Join("/", idx.indexRootDir, v)
		_, err := idx.storageProvider.Delete(ctx, &provider.DeleteRequest{
			Ref: &provider.Reference{
				Path: deletePath,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Update index from <oldV> to <newV>.
func (idx *NonUnique) Update(id, oldV, newV string) error {
	if idx.caseInsensitive {
		oldV = strings.ToLower(oldV)
		newV = strings.ToLower(newV)
	}

	if err := idx.Remove(id, oldV); err != nil {
		return err
	}

	if _, err := idx.Add(id, newV); err != nil {
		return err
	}

	return nil
}

// Search allows for glob search on the index.
func (idx *NonUnique) Search(pattern string) ([]string, error) {
	if idx.caseInsensitive {
		pattern = strings.ToLower(pattern)
	}

	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	foldersMatched := make([]string, 0)
	matches := make([]string, 0)
	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", idx.indexRootDir),
		},
	})

	if err != nil {
		return nil, err
	}

	for _, i := range res.Infos {
		if found, err := filepath.Match(pattern, path.Base(i.Path)); found {
			if err != nil {
				return nil, err
			}

			foldersMatched = append(foldersMatched, i.Path)
		}
	}

	for i := range foldersMatched {
		res, _ := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
			Ref: &provider.Reference{
				Path: foldersMatched[i],
			},
		})

		for _, info := range res.Infos {
			matches = append(matches, path.Base(info.Path))
		}
	}

	return matches, nil
}

// CaseInsensitive undocumented.
func (idx *NonUnique) CaseInsensitive() bool {
	return idx.caseInsensitive
}

// IndexBy undocumented.
func (idx *NonUnique) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *NonUnique) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
func (idx *NonUnique) FilesDir() string {
	return idx.filesDir
}

func (idx *NonUnique) makeDirIfNotExists(folder string) error {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}
	return storage.MakeDirIfNotExist(ctx, idx.storageProvider, folder)
}

func (idx *NonUnique) createSymlink(oldname, newname string) error {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	if _, err := idx.resolveSymlink(newname); err == nil {
		return os.ErrExist
	}

	err = idx.metadataStorage.SimpleUpload(ctx, newname, []byte(oldname))
	if err != nil {
		return err
	}
	return nil
}

func (idx *NonUnique) resolveSymlink(name string) (string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return "", err
	}

	b, err := idx.metadataStorage.SimpleDownload(ctx, name)
	if err != nil {
		if metadatastorage.IsNotFoundErr(err) {
			return "", os.ErrNotExist
		}
		return "", err
	}

	return string(b), err
}

func (idx *NonUnique) getAuthenticatedContext(ctx context.Context) (context.Context, error) {
	t, err := storage.AuthenticateCS3(ctx, idx.cs3conf.ServiceUser, idx.tokenManager)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)
	return ctx, nil
}

// Delete deletes the index folder from its storage.
func (idx *NonUnique) Delete() error {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	return deleteIndexRoot(ctx, idx.storageProvider, idx.indexRootDir)
}

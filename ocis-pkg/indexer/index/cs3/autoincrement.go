package cs3

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/owncloud/ocis/accounts/pkg/storage"

	idxerrs "github.com/owncloud/ocis/ocis-pkg/indexer/errors"

	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
)

// Autoincrement are fields for an index of type autoincrement.
type Autoincrement struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string

	tokenManager    token.Manager
	storageProvider provider.ProviderAPIClient
	dataProvider    dataProviderClient // Used to create and download data via http, bypassing reva upload protocol

	cs3conf *Config
	bound   *option.Bound
}

func init() {
	registry.IndexConstructorRegistry["cs3"]["autoincrement"] = NewAutoincrementIndex
}

// NewAutoincrementIndex instantiates a new AutoincrementIndex instance.
func NewAutoincrementIndex(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	u := &Autoincrement{
		indexBy:      opts.IndexBy,
		typeName:     opts.TypeName,
		filesDir:     opts.FilesDir,
		bound:        opts.Bound,
		indexBaseDir: path.Join(opts.DataDir, "index.cs3"),
		indexRootDir: path.Join(path.Join(opts.DataDir, "index.cs3"), strings.Join([]string{"autoincrement", opts.TypeName, opts.IndexBy}, ".")),
		cs3conf: &Config{
			ProviderAddr: opts.ProviderAddr,
			DataURL:      opts.DataURL,
			DataPrefix:   opts.DataPrefix,
			JWTSecret:    opts.JWTSecret,
			ServiceUser:  opts.ServiceUser,
		},
		dataProvider: dataProviderClient{
			baseURL: singleJoiningSlash(opts.DataURL, opts.DataPrefix),
			client: http.Client{
				Transport: http.DefaultTransport,
			},
		},
	}

	return u
}

// Init initializes an autoincrement index.
func (idx *Autoincrement) Init() error {
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

	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	if err := idx.makeDirIfNotExists(ctx, idx.indexBaseDir); err != nil {
		return err
	}

	if err := idx.makeDirIfNotExists(ctx, idx.indexRootDir); err != nil {
		return err
	}

	return nil
}

// Lookup exact lookup by value.
func (idx *Autoincrement) Lookup(v string) ([]string, error) {
	searchPath := path.Join(idx.indexRootDir, v)
	oldname, err := idx.resolveSymlink(searchPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return nil, err
	}

	return []string{oldname}, nil
}

// Add a new value to the index.
func (idx *Autoincrement) Add(id, v string) (string, error) {
	var newName string
	if v == "" {
		next, err := idx.next()
		if err != nil {
			return "", err
		}
		newName = path.Join(idx.indexRootDir, strconv.Itoa(next))
	} else {
		newName = path.Join(idx.indexRootDir, v)
	}
	if err := idx.createSymlink(id, newName); err != nil {
		if os.IsExist(err) {
			return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return "", err
	}

	return newName, nil
}

// Remove a value v from an index.
func (idx *Autoincrement) Remove(id string, v string) error {
	if v == "" {
		return nil
	}
	searchPath := path.Join(idx.indexRootDir, v)
	_, err := idx.resolveSymlink(searchPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return err
	}

	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	deletePath := path.Join(idx.indexRootDir, v)
	resp, err := idx.storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Path: deletePath,
		},
	})

	if err != nil {
		return err
	}

	// TODO Handle other error codes?
	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		return &idxerrs.NotFoundErr{}
	}

	return err
}

// Update index from <oldV> to <newV>.
func (idx *Autoincrement) Update(id, oldV, newV string) error {
	if err := idx.Remove(id, oldV); err != nil {
		return err
	}

	if _, err := idx.Add(id, newV); err != nil {
		return err
	}

	return nil
}

// Search allows for glob search on the index.
func (idx *Autoincrement) Search(pattern string) ([]string, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return nil, err
	}

	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: idx.indexRootDir,
		},
	})

	if err != nil {
		return nil, err
	}

	searchPath := idx.indexRootDir
	matches := make([]string, 0)
	for _, i := range res.GetInfos() {
		if found, err := filepath.Match(pattern, path.Base(i.Path)); found {
			if err != nil {
				return nil, err
			}

			oldPath, err := idx.resolveSymlink(path.Join(searchPath, path.Base(i.Path)))
			if err != nil {
				return nil, err
			}
			matches = append(matches, oldPath)
		}
	}

	return matches, nil
}

// CaseInsensitive undocumented.
func (idx *Autoincrement) CaseInsensitive() bool {
	return false
}

// IndexBy undocumented.
func (idx *Autoincrement) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *Autoincrement) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
func (idx *Autoincrement) FilesDir() string {
	return idx.filesDir
}

func (idx *Autoincrement) createSymlink(oldname, newname string) error {
	t, err := idx.authenticate(context.TODO())
	if err != nil {
		return err
	}

	if _, err := idx.resolveSymlink(newname); err == nil {
		return os.ErrExist
	}

	resp, err := idx.dataProvider.put(newname, strings.NewReader(oldname), t)
	if err != nil {
		return err
	}
	if err = resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

func (idx *Autoincrement) resolveSymlink(name string) (string, error) {
	t, err := idx.authenticate(context.TODO())
	if err != nil {
		return "", err
	}

	resp, err := idx.dataProvider.get(name, t)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", os.ErrNotExist
		}

		return "", fmt.Errorf("could not resolve symlink %s, got status %v", name, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err = resp.Body.Close(); err != nil {
		return "", err
	}
	return string(b), err
}

func (idx *Autoincrement) makeDirIfNotExists(ctx context.Context, folder string) error {
	return storage.MakeDirIfNotExist(ctx, idx.storageProvider, folder)
}

func (idx *Autoincrement) authenticate(ctx context.Context) (token string, err error) {
	return storage.AuthenticateCS3(ctx, idx.cs3conf.ServiceUser, idx.tokenManager)
}

func (idx *Autoincrement) next() (int, error) {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return -1, err
	}

	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: idx.indexRootDir,
		},
	})

	if err != nil {
		return -1, err
	}

	if len(res.GetInfos()) == 0 {
		return 0, nil
	}

	infos := res.GetInfos()
	sort.Slice(infos, func(i, j int) bool {
		a, _ := strconv.Atoi(path.Base(infos[i].Path))
		b, _ := strconv.Atoi(path.Base(infos[j].Path))
		return a < b
	})

	latest, err := strconv.Atoi(path.Base(infos[len(infos)-1].Path)) // would returning a string be a better interface?
	if err != nil {
		return -1, err
	}

	if int64(latest) < idx.bound.Lower {
		return int(idx.bound.Lower), nil
	}

	return latest + 1, nil
}

func (idx *Autoincrement) getAuthenticatedContext(ctx context.Context) (context.Context, error) {
	t, err := idx.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)
	return ctx, nil
}

// Delete deletes the index folder from its storage.
func (idx *Autoincrement) Delete() error {
	ctx, err := idx.getAuthenticatedContext(context.Background())
	if err != nil {
		return err
	}

	return deleteIndexRoot(ctx, idx.storageProvider, idx.indexRootDir)
}

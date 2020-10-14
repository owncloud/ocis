package cs3

import (
	"context"
	"fmt"
	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/accounts/pkg/indexer/index"
	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	"github.com/owncloud/ocis/accounts/pkg/indexer/registry"
)

type AutoincrementIndex struct {
	indexBy      string
	typeName     string
	filesDir     string
	indexBaseDir string
	indexRootDir string

	tokenManager    token.Manager
	storageProvider provider.ProviderAPIClient
	dataProvider    dataProviderClient // Used to create and download data via http, bypassing reva upload protocol

	cs3conf *Config
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

	u := &AutoincrementIndex{
		indexBy:      opts.IndexBy,
		typeName:     opts.TypeName,
		filesDir:     opts.FilesDir,
		indexBaseDir: path.Join(opts.DataDir, "index.cs3"),
		indexRootDir: path.Join(path.Join(opts.DataDir, "index.cs3"), strings.Join([]string{"autoincrement", opts.TypeName, opts.IndexBy}, ".")),
		cs3conf: &Config{
			ProviderAddr:    opts.ProviderAddr,
			DataURL:         opts.DataURL,
			DataPrefix:      opts.DataPrefix,
			JWTSecret:       opts.JWTSecret,
			ServiceUserName: "",
			ServiceUserUUID: "",
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

func (idx *AutoincrementIndex) Init() error {
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

	ctx := context.Background()
	tk, err := idx.authenticate(ctx)
	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, tk)

	if err := idx.makeDirIfNotExists(ctx, idx.indexBaseDir); err != nil {
		return err
	}

	if err := idx.makeDirIfNotExists(ctx, idx.indexRootDir); err != nil {
		return err
	}

	return nil
}

func (idx AutoincrementIndex) Lookup(v string) ([]string, error) {
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

func (idx AutoincrementIndex) Add(id, v string) (string, error) {
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

func (idx AutoincrementIndex) Remove(id string, v string) error {
	searchPath := path.Join(idx.indexRootDir, v)
	_, err := idx.resolveSymlink(searchPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return err
	}

	ctx := context.Background()
	t, err := idx.authenticate(ctx)
	if err != nil {
		return err
	}

	deletePath := path.Join("/meta", idx.indexRootDir, v)
	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	resp, err := idx.storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Spec: &provider.Reference_Path{Path: deletePath},
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

func (idx AutoincrementIndex) Update(id, oldV, newV string) error {
	if err := idx.Remove(id, oldV); err != nil {
		return err
	}

	if _, err := idx.Add(id, newV); err != nil {
		return err
	}

	return nil
}

func (idx AutoincrementIndex) Search(pattern string) ([]string, error) {
	ctx := context.Background()
	t, err := idx.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Spec: &provider.Reference_Path{Path: path.Join("/meta", idx.indexRootDir)},
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

func (idx AutoincrementIndex) IndexBy() string {
	return idx.indexBy
}

func (idx AutoincrementIndex) TypeName() string {
	return idx.typeName
}

func (idx AutoincrementIndex) FilesDir() string {
	return idx.filesDir
}

func (idx *AutoincrementIndex) createSymlink(oldname, newname string) error {
	t, err := idx.authenticate(context.TODO())
	if err != nil {
		return err
	}

	if _, err := idx.resolveSymlink(newname); err == nil {
		return os.ErrExist
	}

	_, err = idx.dataProvider.put(newname, strings.NewReader(oldname), t)
	if err != nil {
		return err
	}

	return nil

}

func (idx *AutoincrementIndex) resolveSymlink(name string) (string, error) {
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
	return string(b), err
}

func (idx *AutoincrementIndex) makeDirIfNotExists(ctx context.Context, folder string) error {
	var rootPathRef = &provider.Reference{
		Spec: &provider.Reference_Path{Path: fmt.Sprintf("/meta/%v", folder)},
	}

	resp, err := idx.storageProvider.Stat(ctx, &provider.StatRequest{
		Ref: rootPathRef,
	})

	if err != nil {
		return err
	}

	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		_, err := idx.storageProvider.CreateContainer(ctx, &provider.CreateContainerRequest{
			Ref: rootPathRef,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (idx *AutoincrementIndex) authenticate(ctx context.Context) (token string, err error) {
	u := &user.User{
		Id:     &user.UserId{},
		Groups: []string{},
	}
	if idx.cs3conf.ServiceUserName != "" {
		u.Id.OpaqueId = idx.cs3conf.ServiceUserUUID
	}
	return idx.tokenManager.MintToken(ctx, u)
}

func (idx AutoincrementIndex) next() (int, error) {
	ctx := context.Background()
	t, err := idx.authenticate(ctx)
	if err != nil {
		return -1, err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	res, err := idx.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Spec: &provider.Reference_Path{Path: path.Join("/meta", idx.indexRootDir)},
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
	return latest + 1, nil
}

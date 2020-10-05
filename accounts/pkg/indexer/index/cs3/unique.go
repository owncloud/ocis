package cs3

import (
	"context"
	"fmt"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	idxerrs "github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Unique struct {
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

type Config struct {
	ProviderAddr    string
	DataURL         string
	DataPrefix      string
	JWTSecret       string
	ServiceUserName string
	ServiceUserUUID string
}

// NewUniqueIndex instantiates a new UniqueIndex instance. Init() should be
// called afterward to ensure correct on-disk structure.
func NewUniqueIndex(typeName, indexBy, filesDir, indexBaseDir string, cfg *Config) Unique {
	return Unique{
		indexBy:      indexBy,
		typeName:     typeName,
		filesDir:     filesDir,
		indexBaseDir: indexBaseDir,
		indexRootDir: path.Join(indexBaseDir, strings.Join([]string{"unique", typeName, indexBy}, ".")),
		cs3conf:      cfg,
		dataProvider: dataProviderClient{
			client: http.Client{
				Transport: http.DefaultTransport,
			},
		},
	}
}

func (idx *Unique) Init() error {
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

// Add adds a value to the index, returns the path to the root-document
func (idx *Unique) Add(id, v string) (string, error) {
	newName := idx.indexURL(v)
	if err := idx.createSymlink(id, newName); err != nil {
		if os.IsExist(err) {
			return "", &idxerrs.AlreadyExistsErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return "", err
	}

	return newName, nil
}

func (idx *Unique) Lookup(v string) ([]string, error) {
	searchPath := singleJoiningSlash(idx.cs3conf.DataURL, path.Join(idx.cs3conf.DataPrefix, idx.indexRootDir, v))
	oldname, err := idx.resolveSymlink(searchPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = &idxerrs.NotFoundErr{TypeName: idx.typeName, Key: idx.indexBy, Value: v}
		}

		return nil, err
	}

	return []string{oldname}, nil
}

// 97d28b57
func (idx *Unique) Remove(id string, v string) error {
	searchPath := singleJoiningSlash(idx.cs3conf.DataURL, path.Join(idx.cs3conf.DataPrefix, idx.indexRootDir, v))
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

func (idx *Unique) Update(id, oldV, newV string) error {
	if err := idx.Remove(id, oldV); err != nil {
		return err
	}

	if _, err := idx.Add(id, newV); err != nil {
		return err
	}

	return nil
}

func (idx *Unique) Search(pattern string) ([]string, error) {
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

	searchPath := singleJoiningSlash(idx.cs3conf.DataURL, path.Join(idx.cs3conf.DataPrefix, idx.indexRootDir))

	matches := make([]string, 0)
	for _, i := range res.GetInfos() {
		if found, err := filepath.Match(pattern, path.Base(i.Path)); found {
			if err != nil {
				return nil, err
			}

			oldPath, err := idx.resolveSymlink(singleJoiningSlash(searchPath, path.Base(i.Path)))
			if err != nil {
				return nil, err
			}
			matches = append(matches, oldPath)
		}
	}

	return matches, nil

}

func (idx *Unique) IndexBy() string {
	return idx.indexBy
}

func (idx *Unique) TypeName() string {
	return idx.typeName
}

func (idx *Unique) FilesDir() string {
	return idx.filesDir
}

func (idx *Unique) createSymlink(oldname, newname string) error {
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

func (idx *Unique) resolveSymlink(name string) (string, error) {
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

func (idx *Unique) makeDirIfNotExists(ctx context.Context, folder string) error {
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

func (idx *Unique) indexURL(id string) string {
	return singleJoiningSlash(idx.cs3conf.DataURL, path.Join(idx.cs3conf.DataPrefix, idx.indexRootDir, id))
}

// TODO: this is copied from proxy. Find a better solution or move it to ocis-pkg
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (idx *Unique) authenticate(ctx context.Context) (token string, err error) {
	u := &user.User{
		Id:     &user.UserId{},
		Groups: []string{},
	}
	if idx.cs3conf.ServiceUserName != "" {
		u.Id.OpaqueId = idx.cs3conf.ServiceUserUUID
	}
	return idx.tokenManager.MintToken(ctx, u)
}

type dataProviderClient struct {
	client http.Client
}

func (d dataProviderClient) put(url string, body io.Reader, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", token)
	return d.client.Do(req)
}

func (d dataProviderClient) get(url string, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-access-token", token)
	return d.client.Do(req)
}

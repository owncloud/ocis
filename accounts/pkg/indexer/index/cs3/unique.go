package cs3

import (
	"context"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"io"
	"net/http"
	"path"
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

	return nil
}

func (idx *Unique) Add(id, v string) (string, error) {
	//oldName := path.Join(idx.filesDir, id)
	//newName := path.Join(idx.indexRootDir, v)

	panic("implement me")

}

func (idx *Unique) Lookup(v string) ([]string, error) {
	panic("implement me")
}

func (idx *Unique) Remove(id string, v string) error {
	panic("implement me")
}

func (idx *Unique) Update(id, oldV, newV string) error {
	panic("implement me")
}

func (idx *Unique) Search(pattern string) ([]string, error) {
	panic("implement me")
}

func (idx *Unique) IndexBy() string {
	panic("implement me")
}

func (idx *Unique) TypeName() string {
	panic("implement me")
}

func (idx *Unique) FilesDir() string {
	panic("implement me")
}

func (idx *Unique) fakeSymlink(oldname, newname string) {
	//idx.dataProvider.put()

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

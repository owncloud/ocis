package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"google.golang.org/grpc/metadata"
)

type CS3Repo struct {
	serviceID     string
	cfg           *config.Config
	tm            token.Manager
	storageClient provider.ProviderAPIClient
}

func NewCS3Repo(serviceID string, cfg *config.Config) (Repo, error) {
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": cfg.TokenManager.JWTSecret,
	})

	if err != nil {
		return nil, err
	}

	client, err := pool.GetStorageProviderServiceClient(cfg.Repo.CS3.ProviderAddr)
	if err != nil {
		return nil, err
	}

	return CS3Repo{
		serviceID:     serviceID,
		cfg:           cfg,
		tm:            tokenManager,
		storageClient: client,
	}, nil
}

func (r CS3Repo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	if err := r.makeRootDirIfNotExist(ctx, accountsFolder); err != nil {
		return err
	}

	var by []byte
	if by, err = json.Marshal(a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal account: %v", err.Error())
	}

	ureq, err := http.NewRequest("PUT", r.accountUrl(a.Id), bytes.NewReader(by))
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	if _, err := cl.Do(ureq); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	ureq, err := http.NewRequest("GET", r.accountUrl(id), nil)
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	resp, err := cl.Do(ureq)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) DeleteAccount(ctx context.Context, id string) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	_, err = r.storageClient.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Spec: &provider.Reference_Path{Path: fmt.Sprintf("/meta/accounts/%s", id)},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	if err := r.makeRootDirIfNotExist(ctx, groupsFolder); err != nil {
		return err
	}

	var by []byte
	if by, err = json.Marshal(g); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal account: %v", err.Error())
	}

	ureq, err := http.NewRequest("PUT", r.groupUrl(g.Id), bytes.NewReader(by))
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	if _, err := cl.Do(ureq); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)

	ureq, err := http.NewRequest("GET", r.groupUrl(id), nil)
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	resp, err := cl.Do(ureq)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(b, &g); err != nil {
		return err
	}

	return nil
}

func (r CS3Repo) DeleteGroup(ctx context.Context, id string) (err error) {
	t, err := r.authenticate(ctx)
	if err != nil {
		return err
	}

	ctx = metadata.AppendToOutgoingContext(ctx, token.TokenHeader, t)
	ureq, err := http.NewRequest("DELETE", r.groupUrl(id), nil)
	if err != nil {
		return err
	}

	ureq.Header.Add("x-access-token", t)
	cl := http.Client{
		Transport: http.DefaultTransport,
	}

	resp, err := cl.Do(ureq)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (r CS3Repo) authenticate(ctx context.Context) (token string, err error) {
	return r.tm.MintToken(ctx, &user.User{
		Id:     &user.UserId{},
		Groups: []string{},
	})
}

func (r CS3Repo) accountUrl(id string) string {
	return singleJoiningSlash(r.cfg.Repo.CS3.DriverURL, path.Join(r.cfg.Repo.CS3.DataPrefix, accountsFolder, id))
}

func (r CS3Repo) groupUrl(id string) string {
	return singleJoiningSlash(r.cfg.Repo.CS3.DriverURL, path.Join(r.cfg.Repo.CS3.DataPrefix, groupsFolder, id))
}

func (r CS3Repo) makeRootDirIfNotExist(ctx context.Context, folder string) error {
	var rootPathRef = &provider.Reference{
		Spec: &provider.Reference_Path{Path: fmt.Sprintf("/meta/%v", folder)},
	}

	resp, err := r.storageClient.Stat(ctx, &provider.StatRequest{
		Ref: rootPathRef,
	})

	if err != nil {
		return err
	}

	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		_, err := r.storageClient.CreateContainer(ctx, &provider.CreateContainerRequest{
			Ref: rootPathRef,
		})

		if err != nil {
			return err
		}
	}

	return nil
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

package storage

import (
	"context"
	"encoding/json"
	"path"
	"path/filepath"

	"github.com/cs3org/reva/pkg/auth/scope"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	v1beta11 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis/ocis-pkg/log"
	metadatastorage "github.com/owncloud/ocis/ocis-pkg/metadata_storage"
	"google.golang.org/grpc/metadata"
)

const (
	storageMountPath = "/meta"
)

// CS3Repo provides a cs3 implementation of the Repo interface
type CS3Repo struct {
	cfg             *config.Config
	tm              token.Manager
	storageProvider provider.ProviderAPIClient
	metadataStorage metadatastorage.MetadataStorage
}

// NewCS3Repo creates a new cs3 repo
func NewCS3Repo(cfg *config.Config) (Repo, error) {
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

	ms, err := metadatastorage.NewMetadataStorage(cfg.Repo.CS3.ProviderAddr)
	if err != nil {
		return nil, err
	}

	return CS3Repo{
		cfg:             cfg,
		tm:              tokenManager,
		storageProvider: client,
		metadataStorage: ms,
	}, nil
}

// WriteAccount writes an account via cs3 and modifies the provided account (e.g. with a generated id).
func (r CS3Repo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	if err := r.makeRootDirIfNotExist(ctx, accountsFolder); err != nil {
		return err
	}

	var by []byte
	if by, err = json.Marshal(a); err != nil {
		return err
	}

	err = r.metadataStorage.SimpleUpload(ctx, r.accountURL(a.Id), by)
	return err

}

// LoadAccount loads an account via cs3 by id and writes it to the provided account
func (r CS3Repo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	return r.loadAccount(ctx, id, a)
}

// LoadAccounts loads all the accounts from the cs3 api
func (r CS3Repo) LoadAccounts(ctx context.Context, a *[]*proto.Account) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	res, err := r.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", accountsFolder),
		},
	})
	if err != nil {
		return err
	}

	log := olog.NewLogger(olog.Pretty(r.cfg.Log.Pretty), olog.Color(r.cfg.Log.Color), olog.Level(r.cfg.Log.Level))
	for i := range res.Infos {
		acc := &proto.Account{}
		err := r.loadAccount(ctx, filepath.Base(res.Infos[i].Path), acc)
		if err != nil {
			log.Err(err).Msg("could not load account")
			continue
		}
		*a = append(*a, acc)
	}
	return nil
}

func (r CS3Repo) loadAccount(ctx context.Context, id string, a *proto.Account) error {
	account, err := r.metadataStorage.SimpleDownload(ctx, r.accountURL(id))
	if err != nil {
		if metadatastorage.IsNotFoundErr(err) {
			return &notFoundErr{"account", id}
		}
		return err
	}
	return json.Unmarshal(account, &a)
}

// DeleteAccount deletes an account via cs3 by id
func (r CS3Repo) DeleteAccount(ctx context.Context, id string) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	resp, err := r.storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", accountsFolder, id),
		},
	})

	if err != nil {
		return err
	}

	// TODO Handle other error codes?
	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		return &notFoundErr{"account", id}
	}

	return nil
}

// WriteGroup writes a group via cs3 and modifies the provided group (e.g. with a generated id).
func (r CS3Repo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	if err := r.makeRootDirIfNotExist(ctx, groupsFolder); err != nil {
		return err
	}

	var by []byte
	if by, err = json.Marshal(g); err != nil {
		return err
	}

	err = r.metadataStorage.SimpleUpload(ctx, r.groupURL(g.Id), by)
	return err
}

// LoadGroup loads a group via cs3 by id and writes it to the provided group
func (r CS3Repo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	return r.loadGroup(ctx, id, g)
}

// LoadGroups loads all the groups from the cs3 api
func (r CS3Repo) LoadGroups(ctx context.Context, g *[]*proto.Group) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	res, err := r.storageProvider.ListContainer(ctx, &provider.ListContainerRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", groupsFolder),
		},
	})
	if err != nil {
		return err
	}

	log := olog.NewLogger(olog.Pretty(r.cfg.Log.Pretty), olog.Color(r.cfg.Log.Color), olog.Level(r.cfg.Log.Level))
	for i := range res.Infos {
		grp := &proto.Group{}
		err := r.loadGroup(ctx, filepath.Base(res.Infos[i].Path), grp)
		if err != nil {
			log.Err(err).Msg("could not load account")
			continue
		}
		*g = append(*g, grp)
	}
	return nil
}

func (r CS3Repo) loadGroup(ctx context.Context, id string, g *proto.Group) error {
	group, err := r.metadataStorage.SimpleDownload(ctx, r.groupURL(id))
	if err != nil {
		if metadatastorage.IsNotFoundErr(err) {
			return &notFoundErr{"group", id}
		}
		return err
	}
	return json.Unmarshal(group, &g)
}

// DeleteGroup deletes a group via cs3 by id
func (r CS3Repo) DeleteGroup(ctx context.Context, id string) (err error) {
	ctx, err = r.getAuthenticatedContext(ctx)
	if err != nil {
		return err
	}

	resp, err := r.storageProvider.Delete(ctx, &provider.DeleteRequest{
		Ref: &provider.Reference{
			Path: path.Join("/", groupsFolder, id),
		},
	})

	if err != nil {
		return err
	}

	// TODO Handle other error codes?
	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		return &notFoundErr{"group", id}
	}

	return err
}

func (r CS3Repo) getAuthenticatedContext(ctx context.Context) (context.Context, error) {
	t, err := AuthenticateCS3(ctx, r.cfg.ServiceUser, r.tm)
	if err != nil {
		return nil, err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, t)
	return ctx, nil
}

// AuthenticateCS3 mints an auth token for communicating with cs3 storage based on a service user from config
func AuthenticateCS3(ctx context.Context, su config.ServiceUser, tm token.Manager) (token string, err error) {
	u := &user.User{
		Id: &user.UserId{
			OpaqueId: su.UUID,
			Type:     user.UserType_USER_TYPE_APPLICATION,
		},
		Groups:    []string{},
		UidNumber: su.UID,
		GidNumber: su.GID,
	}
	s, err := scope.AddOwnerScope(nil)
	if err != nil {
		return
	}
	return tm.MintToken(ctx, u, s)
}

func (r CS3Repo) accountURL(id string) string {
	return path.Join(accountsFolder, id)
}

func (r CS3Repo) groupURL(id string) string {
	return path.Join(groupsFolder, id)
}

func (r CS3Repo) makeRootDirIfNotExist(ctx context.Context, folder string) error {
	return MakeDirIfNotExist(ctx, r.storageProvider, folder)
}

// MakeDirIfNotExist will create a root node in the metadata storage. Requires an authenticated context.
func MakeDirIfNotExist(ctx context.Context, sp provider.ProviderAPIClient, folder string) error {
	var rootPathRef = &provider.Reference{
		Path: path.Join("/", folder),
	}

	resp, err := sp.Stat(ctx, &provider.StatRequest{
		Ref: rootPathRef,
	})

	if err != nil {
		return err
	}

	if resp.Status.Code == v1beta11.Code_CODE_NOT_FOUND {
		_, err := sp.CreateContainer(ctx, &provider.CreateContainerRequest{
			Ref: rootPathRef,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

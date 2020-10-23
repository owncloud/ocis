package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis/ocis-pkg/log"
)

var groupLock sync.Mutex

// DiskRepo provides A local filesystem implementation of the Repo interface
type DiskRepo struct {
	cfg *config.Config
	log olog.Logger
}

// NewDiskRepo creates A new disk repo
func NewDiskRepo(cfg *config.Config, log olog.Logger) DiskRepo {
	paths := []string{
		filepath.Join(cfg.Repo.Disk.Path, accountsFolder),
		filepath.Join(cfg.Repo.Disk.Path, groupsFolder),
	}
	for i := range paths {
		if _, err := os.Stat(paths[i]); err != nil {
			if os.IsNotExist(err) {
				if err = os.MkdirAll(paths[i], 0700); err != nil {
					log.Fatal().Err(err).Msgf("could not create data folder %v", paths[i])
				}
			}
		}
	}
	return DiskRepo{
		cfg: cfg,
		log: log,
	}
}

// WriteAccount to the local filesystem
func (r DiskRepo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	// leave only the group id
	r.deflateMemberOf(a)

	var bytes []byte
	if bytes, err = json.Marshal(a); err != nil {
		return err
	}

	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, a.Id)
	return ioutil.WriteFile(path, bytes, 0600)
}

// LoadAccount from the local filesystem
func (r DiskRepo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, id)
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		if os.IsNotExist(err) {
			err = &notFoundErr{"account", id}
		}
		return
	}

	return json.Unmarshal(data, a)
}

// LoadAccounts loads all the accounts from the local filesystem
func (r DiskRepo) LoadAccounts(ctx context.Context, a *[]*proto.Account) (err error) {
	root := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		acc := &proto.Account{}
		if err = r.LoadAccount(ctx, filepath.Base(path), acc); err != nil {
			r.log.Err(err).Msg("could not load account")
			return nil
		}
		*a = append(*a, acc)
		return nil
	})
}

// DeleteAccount from the local filesystem
func (r DiskRepo) DeleteAccount(ctx context.Context, id string) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, id)
	if err = os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			err = &notFoundErr{"account", id}
		}
	}

	return
}

// WriteGroup to the local filesystem
func (r DiskRepo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	// leave only the member id
	r.deflateMembers(g)

	var bytes []byte
	if bytes, err = json.Marshal(g); err != nil {
		return err
	}

	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, g.Id)

	groupLock.Lock()
	defer groupLock.Unlock()

	return ioutil.WriteFile(path, bytes, 0600)
}

// LoadGroup from the local filesystem
func (r DiskRepo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, id)

	groupLock.Lock()
	defer groupLock.Unlock()
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		if os.IsNotExist(err) {
			err = &notFoundErr{"group", id}
		}

		return
	}

	return json.Unmarshal(data, g)
}

// LoadGroups loads all the groups from the local filesystem
func (r DiskRepo) LoadGroups(ctx context.Context, g *[]*proto.Group) (err error) {
	root := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		grp := &proto.Group{}
		if err = r.LoadGroup(ctx, filepath.Base(path), grp); err != nil {
			r.log.Err(err).Msg("could not load group")
			return nil
		}
		*g = append(*g, grp)
		return nil
	})
}

// DeleteGroup from the local filesystem
func (r DiskRepo) DeleteGroup(ctx context.Context, id string) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, id)
	if err = os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			err = &notFoundErr{"account", id}
		}
	}

	return
}

// deflateMemberOf replaces the groups of A user with an instance that only contains the id
func (r DiskRepo) deflateMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	var deflated []*proto.Group
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != "" {
			deflated = append(deflated, &proto.Group{Id: a.MemberOf[i].Id})
		} else {
			// TODO fetch and use an id when group only has A name but no id
			r.log.Error().Str("id", a.Id).Interface("group", a.MemberOf[i]).Msg("resolving groups by name is not implemented yet")
		}
	}
	a.MemberOf = deflated
}

// deflateMembers replaces the users of A group with an instance that only contains the id
func (r DiskRepo) deflateMembers(g *proto.Group) {
	if g == nil {
		return
	}
	var deflated []*proto.Account
	for i := range g.Members {
		if g.Members[i].Id != "" {
			deflated = append(deflated, &proto.Account{Id: g.Members[i].Id})
		} else {
			// TODO fetch and use an id when group only has A name but no id
			r.log.Error().Str("id", g.Id).Interface("account", g.Members[i]).Msg("resolving members by name is not implemented yet")
		}
	}
	g.Members = deflated
}

package storage

import (
	"context"
	"encoding/json"
	"github.com/owncloud/ocis/accounts/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
	olog "github.com/owncloud/ocis/ocis-pkg/log"
)

var groupLock sync.Mutex

type DiskRepo struct {
	serviceID string
	cfg       *config.Config
	log       olog.Logger
}

func NewDiskRepo(serviceID string, cfg *config.Config, log olog.Logger) DiskRepo {
	return DiskRepo{
		serviceID: serviceID,
		cfg:       cfg,
		log:       log,
	}
}

// WriteAccount to the storage
func (r DiskRepo) WriteAccount(ctx context.Context, a *proto.Account) (err error) {
	// leave only the group id
	r.deflateMemberOf(a)

	var bytes []byte
	if bytes, err = json.Marshal(a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal account: %v", err.Error())
	}

	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, a.Id)

	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not write account: %v", err.Error())
	}
	return
}

// LoadAccount from the storage
func (r DiskRepo) LoadAccount(ctx context.Context, id string, a *proto.Account) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, id)

	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(r.serviceID, "could not read account: %v", err.Error())
	}

	if err = json.Unmarshal(data, a); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not unmarshal account: %v", err.Error())
	}
	return
}

// DeleteAccount from the storage
func (r DiskRepo) DeleteAccount(ctx context.Context, id string) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, accountsFolder, id)
	if err = os.Remove(path); err != nil {
		r.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove account")
		return merrors.InternalServerError(r.serviceID, "could not remove account: %v", err.Error())
	}

	return nil
}

// WriteGroup persists a given group to the storage
func (r DiskRepo) WriteGroup(ctx context.Context, g *proto.Group) (err error) {
	// leave only the member id
	r.deflateMembers(g)

	var bytes []byte
	if bytes, err = json.Marshal(g); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not marshal group: %v", err.Error())
	}

	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, g.Id)

	groupLock.Lock()
	defer groupLock.Unlock()
	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not write group: %v", err.Error())
	}
	return
}

// LoadGroup from the storage
func (r DiskRepo) LoadGroup(ctx context.Context, id string, g *proto.Group) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, id)

	groupLock.Lock()
	defer groupLock.Unlock()
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(r.serviceID, "could not read group: %v", err.Error())
	}

	if err = json.Unmarshal(data, g); err != nil {
		return merrors.InternalServerError(r.serviceID, "could not unmarshal group: %v", err.Error())
	}

	return
}

func (r DiskRepo) DeleteGroup(ctx context.Context, id string) (err error) {
	path := filepath.Join(r.cfg.Repo.Disk.Path, groupsFolder, id)
	if err = os.Remove(path); err != nil {
		r.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove group")
		return merrors.InternalServerError(r.serviceID, "could not remove group: %v", err.Error())
	}

	return nil
}

// deflateMemberOf replaces the groups of a user with an instance that only contains the id
func (r DiskRepo) deflateMemberOf(a *proto.Account) {
	if a == nil {
		return
	}
	var deflated []*proto.Group
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != "" {
			deflated = append(deflated, &proto.Group{Id: a.MemberOf[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			r.log.Error().Str("id", a.Id).Interface("group", a.MemberOf[i]).Msg("resolving groups by name is not implemented yet")
		}
	}
	a.MemberOf = deflated
}

// deflateMembers replaces the users of a group with an instance that only contains the id
func (r DiskRepo) deflateMembers(g *proto.Group) {
	if g == nil {
		return
	}
	var deflated []*proto.Account
	for i := range g.Members {
		if g.Members[i].Id != "" {
			deflated = append(deflated, &proto.Account{Id: g.Members[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			r.log.Error().Str("id", g.Id).Interface("account", g.Members[i]).Msg("resolving members by name is not implemented yet")
		}
	}
	g.Members = deflated
}

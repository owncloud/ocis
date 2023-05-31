// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

//go:build ceph
// +build ceph

package cephfs

import (
	"context"
	"fmt"
	"time"

	"github.com/ceph/go-ceph/cephfs/admin"
	rados2 "github.com/ceph/go-ceph/rados"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/pkg/errors"

	cephfs2 "github.com/ceph/go-ceph/cephfs"
	"github.com/dgraph-io/ristretto"
	"golang.org/x/sync/semaphore"
)

type cacheVal struct {
	perm  *cephfs2.UserPerm
	mount *cephfs2.MountInfo
}

//TODO: Add to cephfs obj

type connections struct {
	cache      *ristretto.Cache
	lock       *semaphore.Weighted
	ctx        context.Context
	userCache  *ristretto.Cache
	groupCache *ristretto.Cache
}

// TODO: make configurable/add to options
var usrLimit int64 = 1e4

func newCache() (c *connections, err error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     usrLimit,
		BufferItems: 64,
		OnEvict: func(item *ristretto.Item) {
			v := item.Value.(cacheVal)
			v.perm.Destroy()
			_ = v.mount.Unmount()
			_ = v.mount.Release()
		},
	})
	if err != nil {
		return
	}

	ucache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     10 * usrLimit,
		BufferItems: 64,
	})
	if err != nil {
		return
	}

	gcache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     10 * usrLimit,
		BufferItems: 64,
	})
	if err != nil {
		return
	}

	c = &connections{
		cache:      cache,
		lock:       semaphore.NewWeighted(usrLimit),
		ctx:        context.Background(),
		userCache:  ucache,
		groupCache: gcache,
	}

	return
}

func (c *connections) clearCache() {
	c.cache.Clear()
	c.cache.Close()
}

type adminConn struct {
	indexPoolName string
	subvolAdmin   *admin.FSAdmin
	adminMount    Mount
	radosConn     *rados2.Conn
	radosIO       *rados2.IOContext
}

func newAdminConn(poolName string) *adminConn {
	rados, err := rados2.NewConn()
	if err != nil {
		return nil
	}
	if err = rados.ReadDefaultConfigFile(); err != nil {
		return nil
	}

	if err = rados.Connect(); err != nil {
		return nil
	}

	pools, err := rados.ListPools()
	if err != nil {
		rados.Shutdown()
		return nil
	}

	var radosIO *rados2.IOContext
	if in(poolName, pools) {
		radosIO, err = rados.OpenIOContext(poolName)
		if err != nil {
			rados.Shutdown()
			return nil
		}
	} else {
		err = rados.MakePool(poolName)
		if err != nil {
			rados.Shutdown()
			return nil
		}
		radosIO, err = rados.OpenIOContext(poolName)
		if err != nil {
			rados.Shutdown()
			return nil
		}
	}

	mount, err := cephfs2.CreateFromRados(rados)
	if err != nil {
		rados.Shutdown()
		return nil
	}

	if err = mount.Mount(); err != nil {
		rados.Shutdown()
		destroyCephConn(mount, nil)
		return nil
	}

	return &adminConn{
		poolName,
		admin.NewFromConn(rados),
		mount,
		rados,
		radosIO,
	}
}

func newConn(user *User) *cacheVal {
	var perm *cephfs2.UserPerm
	mount, err := cephfs2.CreateMount()
	if err != nil {
		return destroyCephConn(mount, perm)
	}
	if err = mount.ReadDefaultConfigFile(); err != nil {
		return destroyCephConn(mount, perm)
	}
	if err = mount.Init(); err != nil {
		return destroyCephConn(mount, perm)
	}

	if user != nil { //nil creates admin conn
		perm = cephfs2.NewUserPerm(int(user.UidNumber), int(user.GidNumber), []int{})
		if err = mount.SetMountPerms(perm); err != nil {
			return destroyCephConn(mount, perm)
		}
	}

	if err = mount.MountWithRoot("/"); err != nil {
		return destroyCephConn(mount, perm)
	}

	if user != nil {
		if err = mount.ChangeDir(user.fs.conf.Root); err != nil {
			return destroyCephConn(mount, perm)
		}
	}

	return &cacheVal{
		perm:  perm,
		mount: mount,
	}
}

func (fs *cephfs) getUserByID(ctx context.Context, uid string) (*userpb.User, error) {
	if entity, found := fs.conn.userCache.Get(uid); found {
		return entity.(*userpb.User), nil
	}

	selector, err := pool.GatewaySelector(fs.conf.GatewaySvc)
	if err != nil {
		return nil, errors.Wrap(err, "error getting gateway selector")
	}
	client, err := selector.Next()
	if err != nil {
		return nil, errors.Wrap(err, "error selecting next gateway client")
	}
	getUserResp, err := client.GetUserByClaim(ctx, &userpb.GetUserByClaimRequest{
		Claim: "uid",
		Value: uid,
	})

	if err != nil {
		return nil, errors.Wrap(err, "cephfs: error getting user")
	}
	if getUserResp.Status.Code != rpc.Code_CODE_OK {
		return nil, errors.Wrap(err, "cephfs: grpc get user failed")
	}
	fs.conn.userCache.SetWithTTL(uid, getUserResp.User, 1, 24*time.Hour)
	fs.conn.userCache.SetWithTTL(getUserResp.User.Id.OpaqueId, getUserResp.User, 1, 24*time.Hour)

	return getUserResp.User, nil
}

func (fs *cephfs) getUserByOpaqueID(ctx context.Context, oid string) (*userpb.User, error) {
	if entity, found := fs.conn.userCache.Get(oid); found {
		return entity.(*userpb.User), nil
	}
	selector, err := pool.GatewaySelector(fs.conf.GatewaySvc)
	if err != nil {
		return nil, errors.Wrap(err, "error getting gateway selector")
	}
	client, err := selector.Next()
	if err != nil {
		return nil, errors.Wrap(err, "error selecting next gateway client")
	}
	getUserResp, err := client.GetUser(ctx, &userpb.GetUserRequest{
		UserId: &userpb.UserId{
			OpaqueId: oid,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cephfs: error getting user")
	}
	if getUserResp.Status.Code != rpc.Code_CODE_OK {
		return nil, errors.Wrap(err, "cephfs: grpc get user failed")
	}
	fs.conn.userCache.SetWithTTL(fmt.Sprint(getUserResp.User.UidNumber), getUserResp.User, 1, 24*time.Hour)
	fs.conn.userCache.SetWithTTL(oid, getUserResp.User, 1, 24*time.Hour)

	return getUserResp.User, nil
}

func (fs *cephfs) getGroupByID(ctx context.Context, gid string) (*grouppb.Group, error) {
	if entity, found := fs.conn.groupCache.Get(gid); found {
		return entity.(*grouppb.Group), nil
	}

	selector, err := pool.GatewaySelector(fs.conf.GatewaySvc)
	if err != nil {
		return nil, errors.Wrap(err, "error getting gateway selector")
	}
	client, err := selector.Next()
	if err != nil {
		return nil, errors.Wrap(err, "error selecting next gateway client")
	}
	getGroupResp, err := client.GetGroupByClaim(ctx, &grouppb.GetGroupByClaimRequest{
		Claim: "gid",
		Value: gid,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cephfs: error getting group")
	}
	if getGroupResp.Status.Code != rpc.Code_CODE_OK {
		return nil, errors.Wrap(err, "cephfs: grpc get group failed")
	}
	fs.conn.groupCache.SetWithTTL(gid, getGroupResp.Group, 1, 24*time.Hour)
	fs.conn.groupCache.SetWithTTL(getGroupResp.Group.Id.OpaqueId, getGroupResp.Group, 1, 24*time.Hour)

	return getGroupResp.Group, nil
}

func (fs *cephfs) getGroupByOpaqueID(ctx context.Context, oid string) (*grouppb.Group, error) {
	if entity, found := fs.conn.groupCache.Get(oid); found {
		return entity.(*grouppb.Group), nil
	}
	selector, err := pool.GatewaySelector(fs.conf.GatewaySvc)
	if err != nil {
		return nil, errors.Wrap(err, "error getting gateway selector")
	}
	client, err := selector.Next()
	if err != nil {
		return nil, errors.Wrap(err, "error selecting next gateway client")
	}
	getGroupResp, err := client.GetGroup(ctx, &grouppb.GetGroupRequest{
		GroupId: &grouppb.GroupId{
			OpaqueId: oid,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "cephfs: error getting group")
	}
	if getGroupResp.Status.Code != rpc.Code_CODE_OK {
		return nil, errors.Wrap(err, "cephfs: grpc get group failed")
	}
	fs.conn.userCache.SetWithTTL(fmt.Sprint(getGroupResp.Group.GidNumber), getGroupResp.Group, 1, 24*time.Hour)
	fs.conn.userCache.SetWithTTL(oid, getGroupResp.Group, 1, 24*time.Hour)

	return getGroupResp.Group, nil
}

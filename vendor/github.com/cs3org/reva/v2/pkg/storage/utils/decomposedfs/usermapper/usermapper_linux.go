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

package usermapper

import (
	"context"
	"fmt"
	"os/user"
	"runtime"
	"strconv"

	"golang.org/x/sys/unix"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
)

// UnixMapper is a user mapper that maps users to unix uids and gids
type UnixMapper struct {
	baseUid int
	baseGid int
}

// New returns a new user mapper
func NewUnixMapper() *UnixMapper {
	baseUid, _ := unix.SetfsuidRetUid(-1)
	baseGid, _ := unix.SetfsgidRetGid(-1)

	return &UnixMapper{
		baseUid: baseUid,
		baseGid: baseGid,
	}
}

// RunInUserScope runs the given function in the scope of the base user
func (um *UnixMapper) RunInBaseScope(f func() error) error {
	unscope, err := um.ScopeBase()
	if err != nil {
		return err
	}
	defer func() { _ = unscope() }()

	return f()
}

// ScopeBase returns to the base uid and gid returning a function that can be used to restore the previous scope
func (um *UnixMapper) ScopeBase() (func() error, error) {
	return um.ScopeUserByIds(-1, um.baseGid)
}

// ScopeUser returns to the base uid and gid returning a function that can be used to restore the previous scope
func (um *UnixMapper) ScopeUser(ctx context.Context) (func() error, error) {
	u := revactx.ContextMustGetUser(ctx)

	uid, gid, err := um.mapUser(u.Username)
	if err != nil {
		return nil, err
	}
	return um.ScopeUserByIds(uid, gid)
}

// ScopeUserByIds scopes the current user to the given uid and gid returning a function that can be used to restore the previous scope
func (um *UnixMapper) ScopeUserByIds(uid, gid int) (func() error, error) {
	runtime.LockOSThread() // Lock this Goroutine to the current OS thread

	var err error
	var prevUid int
	var prevGid int
	if uid >= 0 {
		prevUid, err = unix.SetfsuidRetUid(uid)
		if err != nil {
			return nil, err
		}
		if testUid, _ := unix.SetfsuidRetUid(-1); testUid != uid {
			return nil, fmt.Errorf("failed to setfsuid to %d", uid)
		}
	}
	if gid >= 0 {
		prevGid, err = unix.SetfsgidRetGid(gid)
		if err != nil {
			return nil, err
		}
		if testGid, _ := unix.SetfsgidRetGid(-1); testGid != gid {
			return nil, fmt.Errorf("failed to setfsgid to %d", gid)
		}
	}

	return func() error {
		if uid >= 0 {
			_ = unix.Setfsuid(prevUid)
		}
		if gid >= 0 {
			_ = unix.Setfsgid(prevGid)
		}
		runtime.UnlockOSThread()
		return nil
	}, nil
}

func (u *UnixMapper) mapUser(username string) (int, int, error) {
	userDetails, err := user.Lookup(username)
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(userDetails.Uid)
	if err != nil {
		return 0, 0, err
	}
	gid, err := strconv.Atoi(userDetails.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

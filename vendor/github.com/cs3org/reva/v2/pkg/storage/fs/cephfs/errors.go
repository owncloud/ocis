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

/*
 #include <string.h>
 #include <errno.h>
 #include <stdlib.h>
*/
import "C"
import (
	"fmt"

	"github.com/cs3org/reva/v2/pkg/errtypes"
)

func wrapErrorMsg(code C.int) string {
	return fmt.Sprintf("cephfs: ret=-%d, %s", code, C.GoString(C.strerror(code)))
}

var (
	errNotFound         = wrapErrorMsg(C.ENOENT)
	errFileExists       = wrapErrorMsg(C.EEXIST)
	errNoSpaceLeft      = wrapErrorMsg(C.ENOSPC)
	errIsADirectory     = wrapErrorMsg(C.EISDIR)
	errPermissionDenied = wrapErrorMsg(C.EACCES)
)

func getRevaError(err error) error {
	if err == nil {
		return nil
	}
	switch err.Error() {
	case errNotFound:
		return errtypes.NotFound("cephfs: dir entry not found")
	case errPermissionDenied:
		return errtypes.PermissionDenied("cephfs: permission denied")
	case errFileExists:
		return errtypes.AlreadyExists("cephfs: file already exists")
	case errNoSpaceLeft:
		return errtypes.InsufficientStorage("cephfs: no space left on device")
	default:
		return errtypes.InternalError(err.Error())
	}
}

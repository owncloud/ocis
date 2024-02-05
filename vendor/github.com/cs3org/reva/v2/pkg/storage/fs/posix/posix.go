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

package posix

import (
	"fmt"
	"path"

	microstore "go-micro.dev/v4/store"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/blobstore"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/fs/posix/tree"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/aspects"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/permissions"
	"github.com/cs3org/reva/v2/pkg/store"
)

func init() {
	registry.Register("posix", New)
}

// New returns an implementation to of the storage.FS interface that talk to
// a local filesystem.
func New(m map[string]interface{}, stream events.Stream) (storage.FS, error) {
	o, err := options.New(m)
	if err != nil {
		return nil, err
	}

	bs, err := blobstore.New(path.Join(o.Root))
	if err != nil {
		return nil, err
	}

	var lu *lookup.Lookup
	switch o.MetadataBackend {
	case "xattrs":
		lu = lookup.New(metadata.XattrsBackend{}, o)
	case "messagepack":
		lu = lookup.New(metadata.NewMessagePackBackend(o.Root, o.FileMetadataCache), o)
	default:
		return nil, fmt.Errorf("unknown metadata backend %s, only 'messagepack' or 'xattrs' (default) supported", o.MetadataBackend)
	}

	tp := tree.New(lu, bs, o, store.Create(
		store.Store(o.IDCache.Store),
		store.TTL(o.IDCache.TTL),
		store.Size(o.IDCache.Size),
		microstore.Nodes(o.IDCache.Nodes...),
		microstore.Database(o.IDCache.Database),
		microstore.Table(o.IDCache.Table),
		store.DisablePersistence(o.IDCache.DisablePersistence),
		store.Authentication(o.IDCache.AuthUsername, o.IDCache.AuthPassword),
	))

	permissionsSelector, err := pool.PermissionsSelector(o.PermissionsSVC, pool.WithTLSMode(o.PermTLSMode))
	if err != nil {
		return nil, err
	}

	p := permissions.NewPermissions(node.NewPermissions(lu), permissionsSelector)

	aspects := aspects.Aspects{
		Lookup:      lu,
		Tree:        tp,
		Permissions: p,
		EventStream: stream,
	}
	fs, err := decomposedfs.New(o, aspects)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

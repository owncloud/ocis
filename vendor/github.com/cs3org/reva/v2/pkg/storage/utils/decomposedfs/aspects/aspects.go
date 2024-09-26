// Copyright 2018-2024 CERN
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

package aspects

import (
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/permissions"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/trashbin"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/tree"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/usermapper"
)

// Aspects holds dependencies for handling aspects of the decomposedfs
type Aspects struct {
	Lookup            node.PathLookup
	Tree              node.Tree
	Blobstore         tree.Blobstore
	Trashbin          trashbin.Trashbin
	Permissions       permissions.Permissions
	EventStream       events.Stream
	DisableVersioning bool
	UserMapper        usermapper.Mapper
}

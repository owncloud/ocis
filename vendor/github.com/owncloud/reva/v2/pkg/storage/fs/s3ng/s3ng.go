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

package s3ng

import (
	"fmt"

	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/reva/v2/pkg/storage/fs/s3ng/blobstore"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs"
	"github.com/rs/zerolog"
)

func init() {
	registry.Register("s3ng", New)
}

// New returns an implementation to of the storage.FS interface that talk to
// a local filesystem.
func New(m map[string]interface{}, stream events.Stream, log *zerolog.Logger) (storage.FS, error) {
	o, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	if !o.S3ConfigComplete() {
		return nil, fmt.Errorf("S3 configuration incomplete")
	}

	defaultPutOptions := blobstore.Options{
		DisableContentSha256:  o.DisableContentSha256,
		DisableMultipart:      o.DisableMultipart,
		SendContentMd5:        o.SendContentMd5,
		ConcurrentStreamParts: o.ConcurrentStreamParts,
		NumThreads:            o.NumThreads,
		PartSize:              o.PartSize,
	}

	bs, err := blobstore.New(o.S3Endpoint, o.S3Region, o.S3Bucket, o.S3AccessKey, o.S3SecretKey, defaultPutOptions)
	if err != nil {
		return nil, err
	}

	return decomposedfs.NewDefault(m, bs, stream, log)
}

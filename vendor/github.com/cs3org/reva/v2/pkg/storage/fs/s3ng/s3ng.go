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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/cs3org/reva/v2/pkg/storage/fs/s3ng/blobstore"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs"
	"github.com/tus/tusd/pkg/s3store"
)

func init() {
	registry.Register("s3ng", New)
}

// New returns an implementation to of the storage.FS interface that talk to
// a local filesystem.
func New(m map[string]interface{}, stream events.Stream) (storage.FS, error) {
	o, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	if !o.S3ConfigComplete() {
		return nil, fmt.Errorf("S3 configuration incomplete")
	}

	bs, err := blobstore.New(o.S3Endpoint, o.S3Region, o.S3Bucket, o.S3AccessKey, o.S3SecretKey)
	if err != nil {
		return nil, err
	}

	s3Config := aws.NewConfig()
	s3Config.WithCredentials(credentials.NewStaticCredentials(o.S3AccessKey, o.S3SecretKey, "")).
		WithEndpoint(o.S3Endpoint).
		WithRegion(o.S3Region).
		WithS3ForcePathStyle(o.S3ForcePathStyle).
		WithDisableSSL(o.S3DisableSSL)

	tusDataStore := s3store.New(o.S3Bucket, s3.New(session.Must(session.NewSession()), s3Config))
	tusDataStore.ObjectPrefix = o.S3UploadObjectPrefix
	tusDataStore.MetadataObjectPrefix = o.S3UploadMetadataPrefix
	tusDataStore.TemporaryDirectory = o.S3UploadTemporaryDirectory

	return decomposedfs.NewDefault(m, bs, tusDataStore, stream)
}

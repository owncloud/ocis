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

package blobstore

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

// Blobstore provides an interface to an s3 compatible blobstore
type Blobstore struct {
	client *minio.Client

	bucket string
}

// New returns a new Blobstore
func New(endpoint, region, bucket, accessKey, secretKey string) (*Blobstore, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse s3 endpoint")
	}

	useSSL := u.Scheme != "http"
	client, err := minio.New(u.Host, &minio.Options{
		Region: region,
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup s3 client")
	}

	return &Blobstore{
		client: client,
		bucket: bucket,
	}, nil
}

// Upload stores some data in the blobstore under the given key
func (bs *Blobstore) Upload(node *node.Node, source string) error {
	reader, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "can not open source file to upload")
	}
	defer reader.Close()

	_, err1 := bs.client.PutObject(context.Background(), bs.bucket, bs.path(node), reader, node.Blobsize, minio.PutObjectOptions{ContentType: "application/octet-stream"})

	if err1 != nil {
		return errors.Wrapf(err, "could not store object '%s' into bucket '%s'", bs.path(node), bs.bucket)
	}
	return nil
}

// Download retrieves a blob from the blobstore for reading
func (bs *Blobstore) Download(node *node.Node) (io.ReadCloser, error) {
	reader, err := bs.client.GetObject(context.Background(), bs.bucket, bs.path(node), minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "could not download object '%s' from bucket '%s'", bs.path(node), bs.bucket)
	}

	stat, err := reader.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "blob path: %s", bs.path(node))
	}

	if node.Blobsize != stat.Size {
		return nil, fmt.Errorf("blob has unexpected size. %d bytes expected, got %d bytes", node.Blobsize, stat.Size)
	}

	return reader, nil
}

// Delete deletes a blob from the blobstore
func (bs *Blobstore) Delete(node *node.Node) error {
	err := bs.client.RemoveObject(context.Background(), bs.bucket, bs.path(node), minio.RemoveObjectOptions{})
	if err != nil {
		return errors.Wrapf(err, "could not delete object '%s' from bucket '%s'", bs.path(node), bs.bucket)
	}
	return nil
}

func (bs *Blobstore) path(node *node.Node) string {
	// https://aws.amazon.com/de/premiumsupport/knowledge-center/s3-prefix-nested-folders-difference/
	// Prefixes are used to partion a bucket. A prefix is everything except the filename.
	// For a file `BucketName/foo/bar/lorem.ipsum`, `BucketName/foo/bar/` is the prefix.
	// There are request limits per prefix, therefore you should have many prefixes.
	// There are no limits to prefixes per bucket, so in general it's better to have more then less.
	//
	// Since the spaceID is always the same for a space, we don't need to pathify that, because it would
	// not yield any performance gains
	return filepath.Clean(filepath.Join(node.SpaceID, lookup.Pathify(node.BlobID, 4, 2)))
}

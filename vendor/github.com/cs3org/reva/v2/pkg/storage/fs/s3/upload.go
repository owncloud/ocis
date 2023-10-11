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

package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/pkg/errors"
)

func (fs *s3FS) Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, req.Ref)
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "error resolving ref")
	}

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
		Body:   req.Body,
	}
	uploader := s3manager.NewUploaderWithClient(fs.client)
	result, err := uploader.Upload(upParams)

	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				return provider.ResourceInfo{}, errtypes.NotFound(fn)
			}
		}
		return provider.ResourceInfo{}, errors.Wrap(err, "s3fs: error creating object "+fn)
	}

	log.Debug().Interface("result", result) // todo cache etag?

	// return id, etag and mtime
	ri, err := fs.GetMD(ctx, req.Ref, []string{}, []string{"id", "etag", "mtime"})
	if err != nil {
		return provider.ResourceInfo{}, err
	}

	return *ri, nil
}

// InitiateUpload returns upload ids corresponding to different protocols it supports
func (fs *s3FS) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	return nil, errtypes.NotSupported("op not supported")
}

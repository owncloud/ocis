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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/mime"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func init() {
	registry.Register("s3", New)
}

type config struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	Prefix    string `mapstructure:"prefix"`
}

func parseConfig(m map[string]interface{}) (*config, error) {
	c := &config{}
	if err := mapstructure.Decode(m, c); err != nil {
		err = errors.Wrap(err, "error decoding conf")
		return nil, err
	}
	return c, nil
}

// New returns an implementation to of the storage.FS interface that talk to
// a s3 api.
func New(m map[string]interface{}) (storage.FS, error) {
	c, err := parseConfig(m)
	if err != nil {
		return nil, err
	}

	awsConfig := aws.NewConfig().
		WithHTTPClient(http.DefaultClient).
		WithMaxRetries(aws.UseServiceDefaultRetries).
		WithLogger(aws.NewDefaultLogger()).
		WithLogLevel(aws.LogOff).
		WithSleepDelay(time.Sleep).
		WithCredentials(credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, "")).
		WithEndpoint(c.Endpoint).
		WithS3ForcePathStyle(true).
		WithDisableSSL(true)

	if c.Region != "" {
		awsConfig.WithRegion(c.Region)
	} else {
		awsConfig.WithRegion("us-east-1")
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, errors.New("creating the S3 session")
	}

	s3Client := s3.New(sess)

	return &s3FS{client: s3Client, config: c}, nil
}

func (fs *s3FS) Shutdown(ctx context.Context) error {
	return nil
}

func (fs *s3FS) addRoot(p string) string {
	np := path.Join(fs.config.Prefix, p)
	return np
}

func (fs *s3FS) resolve(ctx context.Context, ref *provider.Reference) (string, error) {
	if strings.HasPrefix(ref.Path, "/") {
		return fs.addRoot(ref.GetPath()), nil
	}

	if ref.ResourceId != nil && ref.ResourceId.OpaqueId != "" {
		fn := path.Join("/", strings.TrimPrefix(ref.ResourceId.OpaqueId, "fileid-"))
		fn = fs.addRoot(fn)
		return fn, nil
	}

	// reference is invalid
	return "", fmt.Errorf("invalid reference %+v", ref)
}

func (fs *s3FS) removeRoot(np string) string {
	p := strings.TrimPrefix(np, fs.config.Prefix)
	if p == "" {
		p = "/"
	}
	return p
}

type s3FS struct {
	client *s3.S3
	config *config
}

// permissionSet returns the permission set for the current user
func (fs *s3FS) permissionSet(ctx context.Context) *provider.ResourcePermissions {
	// TODO fix permissions for share recipients by traversing reading acls up to the root? cache acls for the parent node and reuse it
	return &provider.ResourcePermissions{
		// owner has all permissions
		AddGrant:             true,
		CreateContainer:      true,
		Delete:               true,
		GetPath:              true,
		GetQuota:             true,
		InitiateFileDownload: true,
		InitiateFileUpload:   true,
		ListContainer:        true,
		ListFileVersions:     true,
		ListGrants:           true,
		ListRecycle:          true,
		Move:                 true,
		PurgeRecycle:         true,
		RemoveGrant:          true,
		RestoreFileVersion:   true,
		RestoreRecycleItem:   true,
		Stat:                 true,
		UpdateGrant:          true,
	}
}

func (fs *s3FS) normalizeObject(ctx context.Context, o *s3.Object, fn string) *provider.ResourceInfo {
	fn = fs.removeRoot(path.Join("/", fn))
	isDir := strings.HasSuffix(*o.Key, "/")
	md := &provider.ResourceInfo{
		Id: &provider.ResourceId{
			OpaqueId: "fileid-" + strings.TrimPrefix(fn, "/"),
		},
		Path:          fn,
		Type:          getResourceType(isDir),
		Etag:          *o.ETag,
		MimeType:      mime.Detect(isDir, fn),
		PermissionSet: fs.permissionSet(ctx),
		Size:          uint64(*o.Size),
		Mtime: &types.Timestamp{
			Seconds: uint64(o.LastModified.Unix()),
		},
	}
	appctx.GetLogger(ctx).Debug().
		Interface("object", o).
		Interface("metadata", md).
		Msg("normalized Object")
	return md
}

func getResourceType(isDir bool) provider.ResourceType {
	if isDir {
		return provider.ResourceType_RESOURCE_TYPE_CONTAINER
	}
	return provider.ResourceType_RESOURCE_TYPE_CONTAINER
}

func (fs *s3FS) normalizeHead(ctx context.Context, o *s3.HeadObjectOutput, fn string) *provider.ResourceInfo {
	fn = fs.removeRoot(path.Join("/", fn))
	isDir := strings.HasSuffix(fn, "/")
	md := &provider.ResourceInfo{
		Id:            &provider.ResourceId{OpaqueId: "fileid-" + strings.TrimPrefix(fn, "/")},
		Path:          fn,
		Type:          getResourceType(isDir),
		Etag:          *o.ETag,
		MimeType:      mime.Detect(isDir, fn),
		PermissionSet: fs.permissionSet(ctx),
		Size:          uint64(*o.ContentLength),
		Mtime: &types.Timestamp{
			Seconds: uint64(o.LastModified.Unix()),
		},
	}
	appctx.GetLogger(ctx).Debug().
		Interface("head", o).
		Interface("metadata", md).
		Msg("normalized Head")
	return md
}
func (fs *s3FS) normalizeCommonPrefix(ctx context.Context, p *s3.CommonPrefix) *provider.ResourceInfo {
	fn := fs.removeRoot(path.Join("/", *p.Prefix))
	md := &provider.ResourceInfo{
		Id:            &provider.ResourceId{OpaqueId: "fileid-" + strings.TrimPrefix(fn, "/")},
		Path:          fn,
		Type:          getResourceType(true),
		Etag:          "TODO(labkode)",
		MimeType:      mime.Detect(true, fn),
		PermissionSet: fs.permissionSet(ctx),
		Size:          0,
		Mtime: &types.Timestamp{
			Seconds: 0,
		},
	}
	appctx.GetLogger(ctx).Debug().
		Interface("prefix", p).
		Interface("metadata", md).
		Msg("normalized CommonPrefix")
	return md
}

// GetPathByID returns the path pointed by the file id
// In this implementation the file id is that path of the file without the first slash
// thus the file id always points to the filename
func (fs *s3FS) GetPathByID(ctx context.Context, id *provider.ResourceId) (string, error) {
	return path.Join("/", strings.TrimPrefix(id.OpaqueId, "fileid-")), nil
}

func (fs *s3FS) AddGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) DenyGrant(ctx context.Context, ref *provider.Reference, g *provider.Grantee) error {
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) ListGrants(ctx context.Context, ref *provider.Reference) ([]*provider.Grant, error) {
	return nil, errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) RemoveGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) UpdateGrant(ctx context.Context, ref *provider.Reference, g *provider.Grant) error {
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) GetQuota(ctx context.Context, ref *provider.Reference) (uint64, uint64, uint64, error) {
	return 0, 0, 0, nil
}

func (fs *s3FS) SetArbitraryMetadata(ctx context.Context, ref *provider.Reference, md *provider.ArbitraryMetadata) error {
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) UnsetArbitraryMetadata(ctx context.Context, ref *provider.Reference, keys []string) error {
	return errtypes.NotSupported("s3: operation not supported")
}

// GetLock returns an existing lock on the given reference
func (fs *s3FS) GetLock(ctx context.Context, ref *provider.Reference) (*provider.Lock, error) {
	return nil, errtypes.NotSupported("unimplemented")
}

// SetLock puts a lock on the given reference
func (fs *s3FS) SetLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error {
	return errtypes.NotSupported("unimplemented")
}

// RefreshLock refreshes an existing lock on the given reference
func (fs *s3FS) RefreshLock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error {
	return errtypes.NotSupported("unimplemented")
}

// Unlock removes an existing lock from the given reference
func (fs *s3FS) Unlock(ctx context.Context, ref *provider.Reference, lock *provider.Lock) error {
	return errtypes.NotSupported("unimplemented")
}

func (fs *s3FS) CreateReference(ctx context.Context, path string, targetURI *url.URL) error {
	// TODO(jfd):implement
	return errtypes.NotSupported("s3: operation not supported")
}

func (fs *s3FS) GetHome(ctx context.Context) (string, error) {
	return "", errtypes.NotSupported("eos: not supported")
}

func (fs *s3FS) CreateHome(ctx context.Context) error {
	return errtypes.NotSupported("s3fs: not supported")
}

func (fs *s3FS) CreateDir(ctx context.Context, ref *provider.Reference) error {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return nil
	}

	fn = fs.addRoot(fn) + "/" // append / to indicate folder // TODO only if fn does not end in /

	input := &s3.PutObjectInput{
		Bucket:        aws.String(fs.config.Bucket),
		Key:           aws.String(fn),
		ContentType:   aws.String("application/octet-stream"),
		ContentLength: aws.Int64(0),
	}

	result, err := fs.client.PutObject(input)
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				return errtypes.NotFound(ref.Path)
			}
		}
		// FIXME we also need already exists error, webdav expects 405 MethodNotAllowed
		return errors.Wrap(err, "s3fs: error creating dir "+ref.Path)
	}

	log.Debug().Interface("result", result) // todo cache etag?
	return nil
}

// TouchFile as defined in the storage.FS interface
func (fs *s3FS) TouchFile(ctx context.Context, ref *provider.Reference) error {
	return fmt.Errorf("unimplemented: TouchFile")
}

func (fs *s3FS) Delete(ctx context.Context, ref *provider.Reference) error {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "error resolving ref")
	}

	// first we need to find out if fn is a dir or a file

	_, err = fs.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
	})
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case s3.ErrCodeNoSuchKey:
				return errtypes.NotFound(fn)
			}
		}
		// it might be a directory, so we can batch delete the prefix + /
		iter := s3manager.NewDeleteListIterator(fs.client, &s3.ListObjectsInput{
			Bucket: aws.String(fs.config.Bucket),
			Prefix: aws.String(fn + "/"),
		})
		batcher := s3manager.NewBatchDeleteWithClient(fs.client)
		if err := batcher.Delete(aws.BackgroundContext(), iter); err != nil {
			return err
		}
		// ok, we are done
		return nil
	}

	// we found an object, let's get rid of it
	result, err := fs.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
	})
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case s3.ErrCodeNoSuchKey:
				return errtypes.NotFound(fn)
			}
		}
		return errors.Wrap(err, "s3fs: error deleting "+fn)
	}

	log.Debug().Interface("result", result)
	return nil
}

// CreateStorageSpace creates a storage space
func (fs *s3FS) CreateStorageSpace(ctx context.Context, req *provider.CreateStorageSpaceRequest) (*provider.CreateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("unimplemented: CreateStorageSpace")
}

func (fs *s3FS) moveObject(ctx context.Context, oldKey string, newKey string) error {

	// Copy
	// TODO double check CopyObject can deal with >5GB files.
	// Docs say we need to use multipart upload: https://docs.aws.amazon.com/AmazonS3/latest/API/RESTObjectCOPY.html
	_, err := fs.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(fs.config.Bucket),
		CopySource: aws.String("/" + fs.config.Bucket + oldKey),
		Key:        aws.String(newKey),
	})
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == s3.ErrCodeNoSuchBucket {
			return errtypes.NotFound(oldKey)
		}
		return err
	}
	// TODO cache etag and mtime?

	// Delete
	_, err = fs.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(oldKey),
	})
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeNoSuchBucket:
		case s3.ErrCodeNoSuchKey:
			return errtypes.NotFound(oldKey)
		}
		return err
	}
	return nil
}

func (fs *s3FS) Move(ctx context.Context, oldRef, newRef *provider.Reference) error {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, oldRef)
	if err != nil {
		return errors.Wrap(err, "error resolving ref")
	}

	newName, err := fs.resolve(ctx, newRef)
	if err != nil {
		return errors.Wrap(err, "error resolving ref")
	}

	// first we need to find out if fn is a dir or a file

	_, err = fs.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
	})
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case s3.ErrCodeNoSuchKey:
				return errtypes.NotFound(fn)
			}
		}

		// move directory
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(fs.config.Bucket),
			Prefix: aws.String(fn + "/"),
		}
		isTruncated := true

		for isTruncated {
			output, err := fs.client.ListObjectsV2(input)
			if err != nil {
				return errors.Wrap(err, "s3FS: error listing "+fn)
			}

			for _, o := range output.Contents {
				log.Debug().
					Interface("object", *o).
					Str("fn", fn).
					Msg("found Object")

				err := fs.moveObject(ctx, *o.Key, strings.Replace(*o.Key, fn+"/", newName+"/", 1))
				if err != nil {
					return err
				}
			}

			input.ContinuationToken = output.NextContinuationToken
			isTruncated = *output.IsTruncated
		}
		// ok, we are done
		return nil
	}

	// move single object
	err = fs.moveObject(ctx, fn, newName)
	if err != nil {
		return err
	}
	return nil
}

func (fs *s3FS) GetMD(ctx context.Context, ref *provider.Reference, mdKeys []string) (*provider.ResourceInfo, error) {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving ref")
	}

	// first try a head, works for files
	log.Debug().
		Str("fn", fn).
		Msg("trying HEAD")

	input := &s3.HeadObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
	}
	output, err := fs.client.HeadObject(input)
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case s3.ErrCodeNoSuchKey:
				return nil, errtypes.NotFound(fn)
			}
		}
		log.Debug().
			Str("fn", fn).
			Msg("trying to list prefix")
		// try by listing parent to find directory
		input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(fs.config.Bucket),
			Prefix:    aws.String(fn),
			Delimiter: aws.String("/"), // limit to a single directory
		}
		isTruncated := true

		for isTruncated {
			output, err := fs.client.ListObjectsV2(input)
			if err != nil {
				return nil, errors.Wrap(err, "s3FS: error listing "+fn)
			}

			for i := range output.CommonPrefixes {
				log.Debug().
					Interface("object", output.CommonPrefixes[i]).
					Str("fn", fn).
					Msg("found CommonPrefix")
				if *output.CommonPrefixes[i].Prefix == fn+"/" {
					return fs.normalizeCommonPrefix(ctx, output.CommonPrefixes[i]), nil
				}
			}

			input.ContinuationToken = output.NextContinuationToken
			isTruncated = *output.IsTruncated
		}
		return nil, errtypes.NotFound(fn)
	}

	return fs.normalizeHead(ctx, output, fn), nil
}

func (fs *s3FS) ListFolder(ctx context.Context, ref *provider.Reference, mdKeys []string) ([]*provider.ResourceInfo, error) {
	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving ref")
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(fs.config.Bucket),
		Prefix:    aws.String(fn + "/"),
		Delimiter: aws.String("/"), // limit to a single directory
	}
	isTruncated := true

	finfos := []*provider.ResourceInfo{}

	for isTruncated {
		output, err := fs.client.ListObjectsV2(input)
		if err != nil {
			return nil, errors.Wrap(err, "s3FS: error listing "+fn)
		}

		for i := range output.CommonPrefixes {
			finfos = append(finfos, fs.normalizeCommonPrefix(ctx, output.CommonPrefixes[i]))
		}

		for i := range output.Contents {
			finfos = append(finfos, fs.normalizeObject(ctx, output.Contents[i], *output.Contents[i].Key))
		}

		input.ContinuationToken = output.NextContinuationToken
		isTruncated = *output.IsTruncated
	}
	// TODO sort fileinfos?
	return finfos, nil
}

func (fs *s3FS) Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, _ storage.UploadFinishedFunc) error {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "error resolving ref")
	}

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
		Body:   r,
	}
	uploader := s3manager.NewUploaderWithClient(fs.client)
	result, err := uploader.Upload(upParams)

	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket {
				return errtypes.NotFound(fn)
			}
		}
		return errors.Wrap(err, "s3fs: error creating object "+fn)
	}

	log.Debug().Interface("result", result) // todo cache etag?
	return nil
}

func (fs *s3FS) Download(ctx context.Context, ref *provider.Reference) (io.ReadCloser, error) {
	log := appctx.GetLogger(ctx)

	fn, err := fs.resolve(ctx, ref)
	if err != nil {
		return nil, errors.Wrap(err, "error resolving ref")
	}

	// use GetObject instead of s3manager.Downloader:
	// the result.Body is a ReadCloser, which allows streaming
	// TODO double check we are not caching bytes in memory
	r, err := fs.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(fs.config.Bucket),
		Key:    aws.String(fn),
	})
	if err != nil {
		log.Error().Err(err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case s3.ErrCodeNoSuchKey:
				return nil, errtypes.NotFound(fn)
			}
		}
		return nil, errors.Wrap(err, "s3fs: error deleting "+fn)
	}
	return r.Body, nil
}

func (fs *s3FS) ListRevisions(ctx context.Context, ref *provider.Reference) ([]*provider.FileVersion, error) {
	return nil, errtypes.NotSupported("list revisions")
}

func (fs *s3FS) DownloadRevision(ctx context.Context, ref *provider.Reference, revisionKey string) (io.ReadCloser, error) {
	return nil, errtypes.NotSupported("download revision")
}

func (fs *s3FS) RestoreRevision(ctx context.Context, ref *provider.Reference, revisionKey string) error {
	return errtypes.NotSupported("restore revision")
}

func (fs *s3FS) PurgeRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string) error {
	return errtypes.NotSupported("purge recycle item")
}

func (fs *s3FS) EmptyRecycle(ctx context.Context, ref *provider.Reference) error {
	return errtypes.NotSupported("empty recycle")
}

func (fs *s3FS) ListRecycle(ctx context.Context, ref *provider.Reference, key, relativePath string) ([]*provider.RecycleItem, error) {
	return nil, errtypes.NotSupported("list recycle")
}

func (fs *s3FS) RestoreRecycleItem(ctx context.Context, ref *provider.Reference, key, relativePath string, restoreRef *provider.Reference) error {
	return errtypes.NotSupported("restore recycle")
}

func (fs *s3FS) ListStorageSpaces(ctx context.Context, filter []*provider.ListStorageSpacesRequest_Filter, unrestricted bool) ([]*provider.StorageSpace, error) {
	return nil, errtypes.NotSupported("list storage spaces")
}

// UpdateStorageSpace updates a storage space
func (fs *s3FS) UpdateStorageSpace(ctx context.Context, req *provider.UpdateStorageSpaceRequest) (*provider.UpdateStorageSpaceResponse, error) {
	return nil, errtypes.NotSupported("update storage space")
}

// DeleteStorageSpace deletes a storage space
func (fs *s3FS) DeleteStorageSpace(ctx context.Context, req *provider.DeleteStorageSpaceRequest) error {
	return errtypes.NotSupported("delete storage space")
}

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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	cephfs2 "github.com/ceph/go-ceph/cephfs"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctx2 "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/pkg/handler"
)

func (fs *cephfs) Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, uff storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	user := fs.makeUser(ctx)
	upload, err := fs.GetUpload(ctx, ref.GetPath())
	if err != nil {
		metadata := map[string]string{"sizedeferred": "true"}
		uploadIDs, err := fs.InitiateUpload(ctx, ref, 0, metadata)
		if err != nil {
			return provider.ResourceInfo{}, err
		}
		if upload, err = fs.GetUpload(ctx, uploadIDs["simple"]); err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "cephfs: error retrieving upload")
		}
	}

	uploadInfo := upload.(*fileUpload)

	p := uploadInfo.info.Storage["InternalDestination"]
	ok, err := IsChunked(p)
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "cephfs: error checking path")
	}
	if ok {
		var assembledFile string
		p, assembledFile, err = NewChunkHandler(ctx, fs).WriteChunk(p, r)
		if err != nil {
			return provider.ResourceInfo{}, err
		}
		if p == "" {
			if err = uploadInfo.Terminate(ctx); err != nil {
				return provider.ResourceInfo{}, errors.Wrap(err, "cephfs: error removing auxiliary files")
			}
			return provider.ResourceInfo{}, errtypes.PartialContent(ref.String())
		}
		uploadInfo.info.Storage["InternalDestination"] = p

		user.op(func(cv *cacheVal) {
			r, err = cv.mount.Open(assembledFile, os.O_RDONLY, 0)
		})
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "cephfs: error opening assembled file")
		}
		defer r.Close()
		defer user.op(func(cv *cacheVal) {
			_ = cv.mount.Unlink(assembledFile)
		})
	}
	ri := provider.ResourceInfo{
		// fill with at least fileid, mtime and etag
		Id: &provider.ResourceId{
			StorageId: uploadInfo.info.MetaData["providerID"],
			SpaceId:   uploadInfo.info.Storage["SpaceRoot"],
			OpaqueId:  uploadInfo.info.Storage["NodeId"],
		},
		Etag: uploadInfo.info.MetaData["etag"],
	}

	if mtime, err := utils.MTimeToTS(uploadInfo.info.MetaData["mtime"]); err == nil {
		ri.Mtime = &mtime
	}

	if _, err := uploadInfo.WriteChunk(ctx, 0, r); err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "cephfs: error writing to binary file")
	}

	return ri, uploadInfo.FinishUpload(ctx)
}

func (fs *cephfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	user := fs.makeUser(ctx)
	np, err := user.resolveRef(ref)
	if err != nil {
		return nil, errors.Wrap(err, "cephfs: error resolving reference")
	}

	info := tusd.FileInfo{
		MetaData: tusd.MetaData{
			"filename": filepath.Base(np),
			"dir":      filepath.Dir(np),
		},
		Size: uploadLength,
	}

	if metadata != nil {
		info.MetaData["providerID"] = metadata["providerID"]
		if metadata["mtime"] != "" {
			info.MetaData["mtime"] = metadata["mtime"]
		}
		if _, ok := metadata["sizedeferred"]; ok {
			info.SizeIsDeferred = true
		}
	}

	upload, err := fs.NewUpload(ctx, info)
	if err != nil {
		return nil, err
	}

	info, _ = upload.GetInfo(ctx)

	return map[string]string{
		"simple": info.ID,
		"tus":    info.ID,
	}, nil
}

// UseIn tells the tus upload middleware which extensions it supports.
func (fs *cephfs) UseIn(composer *tusd.StoreComposer) {
	composer.UseCore(fs)
	composer.UseTerminater(fs)
}

func (fs *cephfs) NewUpload(ctx context.Context, info tusd.FileInfo) (upload tusd.Upload, err error) {
	log := appctx.GetLogger(ctx)
	log.Debug().Interface("info", info).Msg("cephfs: NewUpload")

	user := fs.makeUser(ctx)

	fn := info.MetaData["filename"]
	if fn == "" {
		return nil, errors.New("cephfs: missing filename in metadata")
	}
	info.MetaData["filename"] = filepath.Clean(info.MetaData["filename"])

	dir := info.MetaData["dir"]
	if dir == "" {
		return nil, errors.New("cephfs: missing dir in metadata")
	}
	info.MetaData["dir"] = filepath.Clean(info.MetaData["dir"])

	np := filepath.Join(info.MetaData["dir"], info.MetaData["filename"])

	info.ID = uuid.New().String()

	binPath := fs.getUploadPath(info.ID)

	info.Storage = map[string]string{
		"Type":                "Cephfs",
		"BinPath":             binPath,
		"InternalDestination": np,

		"Idp":      user.Id.Idp,
		"UserId":   user.Id.OpaqueId,
		"UserName": user.Username,
		"UserType": utils.UserTypeToString(user.Id.Type),

		"LogLevel": log.GetLevel().String(),
	}

	// Create binary file with no content
	user.op(func(cv *cacheVal) {
		var f *cephfs2.File
		defer closeFile(f)
		f, err = cv.mount.Open(binPath, os.O_CREATE|os.O_WRONLY, filePermDefault)
		if err != nil {
			return
		}
	})
	//TODO: if we get two same upload ids, the second one can't upload at all
	if err != nil {
		return
	}

	upload = &fileUpload{
		info:     info,
		binPath:  binPath,
		infoPath: binPath + ".info",
		fs:       fs,
		ctx:      ctx,
	}

	if !info.SizeIsDeferred && info.Size == 0 {
		log.Debug().Interface("info", info).Msg("cephfs: finishing upload for empty file")
		// no need to create info file and finish directly
		err = upload.FinishUpload(ctx)

		return
	}

	// writeInfo creates the file by itself if necessary
	err = upload.(*fileUpload).writeInfo()

	return
}

func (fs *cephfs) getUploadPath(uploadID string) string {
	return filepath.Join(fs.conf.UploadFolder, uploadID)
}

// GetUpload returns the Upload for the given upload id
func (fs *cephfs) GetUpload(ctx context.Context, id string) (fup tusd.Upload, err error) {
	binPath := fs.getUploadPath(id)
	info := tusd.FileInfo{}
	if err != nil {
		return nil, errtypes.NotFound("bin path for upload " + id + " not found")
	}
	infoPath := binPath + ".info"

	var data bytes.Buffer
	f, err := fs.adminConn.adminMount.Open(infoPath, os.O_RDONLY, 0)
	if err != nil {
		return
	}
	_, err = io.Copy(&data, f)
	if err != nil {
		return
	}
	if err = json.Unmarshal(data.Bytes(), &info); err != nil {
		return
	}

	u := &userpb.User{
		Id: &userpb.UserId{
			Idp:      info.Storage["Idp"],
			OpaqueId: info.Storage["UserId"],
		},
		Username: info.Storage["UserName"],
	}
	ctx = ctx2.ContextSetUser(ctx, u)
	user := fs.makeUser(ctx)

	var stat Statx
	user.op(func(cv *cacheVal) {
		stat, err = cv.mount.Statx(binPath, cephfs2.StatxSize, 0)
	})
	if err != nil {
		return
	}
	info.Offset = int64(stat.Size)

	return &fileUpload{
		info:     info,
		binPath:  binPath,
		infoPath: infoPath,
		fs:       fs,
		ctx:      ctx,
	}, nil
}

type fileUpload struct {
	// info stores the current information about the upload
	info tusd.FileInfo
	// infoPath is the path to the .info file
	infoPath string
	// binPath is the path to the binary file (which has no extension)
	binPath string
	// only fs knows how to handle metadata and versions
	fs *cephfs
	// a context with a user
	ctx context.Context
}

// GetInfo returns the FileInfo
func (upload *fileUpload) GetInfo(ctx context.Context) (tusd.FileInfo, error) {
	return upload.info, nil
}

// GetReader returns an io.Reader for the upload
func (upload *fileUpload) GetReader(ctx context.Context) (file io.Reader, err error) {
	user := upload.fs.makeUser(upload.ctx)
	user.op(func(cv *cacheVal) {
		file, err = cv.mount.Open(upload.binPath, os.O_RDONLY, 0)
	})
	return
}

// WriteChunk writes the stream from the reader to the given offset of the upload
func (upload *fileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (n int64, err error) {
	var file io.WriteCloser
	user := upload.fs.makeUser(upload.ctx)
	user.op(func(cv *cacheVal) {
		file, err = cv.mount.Open(upload.binPath, os.O_WRONLY|os.O_APPEND, 0)
	})
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err = io.Copy(file, src)

	// If the HTTP PATCH request gets interrupted in the middle (e.g. because
	// the user wants to pause the upload), Go's net/http returns an io.ErrUnexpectedEOF.
	// However, for OwnCloudStore it's not important whether the stream has ended
	// on purpose or accidentally.
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return n, err
		}
	}

	upload.info.Offset += n
	err = upload.writeInfo()

	return n, err
}

// writeInfo updates the entire information. Everything will be overwritten.
func (upload *fileUpload) writeInfo() error {
	data, err := json.Marshal(upload.info)

	if err != nil {
		return err
	}
	user := upload.fs.makeUser(upload.ctx)
	user.op(func(cv *cacheVal) {
		var file io.WriteCloser
		if file, err = cv.mount.Open(upload.infoPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, filePermDefault); err != nil {
			return
		}
		defer file.Close()

		_, err = io.Copy(file, bytes.NewReader(data))
	})

	return err
}

// FinishUpload finishes an upload and moves the file to the internal destination
func (upload *fileUpload) FinishUpload(ctx context.Context) (err error) {

	np := upload.info.Storage["InternalDestination"]

	// TODO check etag with If-Match header
	// if destination exists
	// if _, err := os.Stat(np); err == nil {
	// the local storage does not store metadata
	// the fileid is based on the path, so no we do not need to copy it to the new file
	// the local storage does not track revisions
	// }

	// if destination exists
	// if _, err := os.Stat(np); err == nil {
	// create revision
	//	if err := upload.fs.archiveRevision(upload.ctx, np); err != nil {
	//		return err
	//	}
	// }

	user := upload.fs.makeUser(upload.ctx)
	log := appctx.GetLogger(ctx)

	user.op(func(cv *cacheVal) {
		err = cv.mount.Rename(upload.binPath, np)
	})
	if err != nil {
		return errors.Wrap(err, upload.binPath)
	}

	// only delete the upload if it was successfully written to the fs
	user.op(func(cv *cacheVal) {
		err = cv.mount.Unlink(upload.infoPath)
	})
	if err != nil {
		if err.Error() != errNotFound {
			log.Err(err).Interface("info", upload.info).Msg("cephfs: could not delete upload metadata")
		}
	}

	// TODO: set mtime if specified in metadata

	return
}

// To implement the termination extension as specified in https://tus.io/protocols/resumable-upload.html#termination
// - the storage needs to implement AsTerminatableUpload
// - the upload needs to implement Terminate

// AsTerminatableUpload returns a a TerminatableUpload
func (fs *cephfs) AsTerminatableUpload(upload tusd.Upload) tusd.TerminatableUpload {
	return upload.(*fileUpload)
}

// Terminate terminates the upload
func (upload *fileUpload) Terminate(ctx context.Context) (err error) {
	user := upload.fs.makeUser(upload.ctx)

	user.op(func(cv *cacheVal) {
		if err = cv.mount.Unlink(upload.infoPath); err != nil {
			return
		}
		err = cv.mount.Unlink(upload.binPath)
	})

	return
}

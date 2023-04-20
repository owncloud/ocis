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

package decomposedfs

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	tusd "github.com/tus/tusd/pkg/handler"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/upload"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

var _idRegexp = regexp.MustCompile(".*/([^/]+).info")

// Upload uploads data to the given resource
// TODO Upload (and InitiateUpload) needs a way to receive the expected checksum.
// Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, uff storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	up, err := fs.GetUpload(ctx, ref.GetPath())
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload")
	}

	uploadInfo := up.(*upload.Upload)

	p := uploadInfo.Info.Storage["NodeName"]
	if chunking.IsChunked(p) { // check chunking v1
		var assembledFile string
		p, assembledFile, err = fs.chunkHandler.WriteChunk(p, r)
		if err != nil {
			return provider.ResourceInfo{}, err
		}
		if p == "" {
			if err = uploadInfo.Terminate(ctx); err != nil {
				return provider.ResourceInfo{}, errors.Wrap(err, "ocfs: error removing auxiliary files")
			}
			return provider.ResourceInfo{}, errtypes.PartialContent(ref.String())
		}
		uploadInfo.Info.Storage["NodeName"] = p
		fd, err := os.Open(assembledFile)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)
		r = fd
	}

	if _, err := uploadInfo.WriteChunk(ctx, 0, r); err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing to binary file")
	}

	if err := uploadInfo.FinishUpload(ctx); err != nil {
		return provider.ResourceInfo{}, err
	}

	if uff != nil {
		info := uploadInfo.Info
		uploadRef := &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: info.MetaData["providerID"],
				SpaceId:   info.Storage["SpaceRoot"],
				OpaqueId:  info.Storage["SpaceRoot"],
			},
			Path: utils.MakeRelativePath(filepath.Join(info.MetaData["dir"], info.MetaData["filename"])),
		}
		owner, ok := ctxpkg.ContextGetUser(uploadInfo.Ctx)
		if !ok {
			return provider.ResourceInfo{}, errtypes.PreconditionFailed("error getting user from uploadinfo context")
		}
		spaceOwner := &userpb.UserId{
			OpaqueId: info.Storage["SpaceOwnerOrManager"],
		}
		uff(spaceOwner, owner.Id, uploadRef)
	}

	ri := provider.ResourceInfo{
		// fill with at least fileid, mtime and etag
		Id: &provider.ResourceId{
			StorageId: uploadInfo.Info.MetaData["providerID"],
			SpaceId:   uploadInfo.Info.Storage["SpaceRoot"],
			OpaqueId:  uploadInfo.Info.Storage["NodeId"],
		},
		Etag: uploadInfo.Info.MetaData["etag"],
	}

	if mtime, err := utils.MTimeToTS(uploadInfo.Info.MetaData["mtime"]); err == nil {
		ri.Mtime = &mtime
	}

	return ri, nil
}

// InitiateUpload returns upload ids corresponding to different protocols it supports
// TODO read optional content for small files in this request
// TODO InitiateUpload (and Upload) needs a way to receive the expected checksum. Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	log := appctx.GetLogger(ctx)

	n, err := fs.lu.NodeFromResource(ctx, ref)
	switch err.(type) {
	case nil:
		// ok
	case errtypes.IsNotFound:
		return nil, errtypes.PreconditionFailed(err.Error())
	default:
		return nil, err
	}

	// permissions are checked in NewUpload below

	relative, err := fs.lu.Path(ctx, n, node.NoCheck)
	if err != nil {
		return nil, err
	}

	lockID, _ := ctxpkg.ContextGetLockID(ctx)

	info := tusd.FileInfo{
		MetaData: tusd.MetaData{
			"filename": filepath.Base(relative),
			"dir":      filepath.Dir(relative),
			"lockid":   lockID,
		},
		Size: uploadLength,
		Storage: map[string]string{
			"SpaceRoot":           n.SpaceRoot.ID,
			"SpaceOwnerOrManager": n.SpaceOwnerOrManager(ctx).GetOpaqueId(),
		},
	}

	if metadata != nil {
		info.MetaData["providerID"] = metadata["providerID"]
		if mtime, ok := metadata["mtime"]; ok {
			info.MetaData["mtime"] = mtime
		}
		if expiration, ok := metadata["expires"]; ok {
			info.MetaData["expires"] = expiration
		}
		if _, ok := metadata["sizedeferred"]; ok {
			info.SizeIsDeferred = true
		}
		if checksum, ok := metadata["checksum"]; ok {
			parts := strings.SplitN(checksum, " ", 2)
			if len(parts) != 2 {
				return nil, errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
			}
			switch parts[0] {
			case "sha1", "md5", "adler32":
				info.MetaData["checksum"] = checksum
			default:
				return nil, errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
			}
		}
		if ifMatch, ok := metadata["if-match"]; ok {
			info.MetaData["if-match"] = ifMatch
		}
	}

	log.Debug().Interface("info", info).Interface("node", n).Interface("metadata", metadata).Msg("Decomposedfs: resolved filename")

	_, err = node.CheckQuota(n.SpaceRoot, n.Exists, uint64(n.Blobsize), uint64(info.Size))
	if err != nil {
		return nil, err
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
func (fs *Decomposedfs) UseIn(composer *tusd.StoreComposer) {
	composer.UseCore(fs)
	composer.UseTerminater(fs)
	composer.UseConcater(fs)
	composer.UseLengthDeferrer(fs)
}

// To implement the core tus.io protocol as specified in https://tus.io/protocols/resumable-upload.html#core-protocol
// - the storage needs to implement NewUpload and GetUpload
// - the upload needs to implement the tusd.Upload interface: WriteChunk, GetInfo, GetReader and FinishUpload

// NewUpload returns a new tus Upload instance
func (fs *Decomposedfs) NewUpload(ctx context.Context, info tusd.FileInfo) (tusd.Upload, error) {
	return upload.New(ctx, info, fs.lu, fs.tp, fs.p, fs.o.Root, fs.stream, fs.o.AsyncFileUploads, fs.o.Tokens)
}

// GetUpload returns the Upload for the given upload id
func (fs *Decomposedfs) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	return upload.Get(ctx, id, fs.lu, fs.tp, fs.o.Root, fs.stream, fs.o.AsyncFileUploads, fs.o.Tokens)
}

// ListUploads returns a list of all incomplete uploads
func (fs *Decomposedfs) ListUploads() ([]tusd.FileInfo, error) {
	return fs.uploadInfos(context.Background())
}

// PurgeExpiredUploads scans the fs for expired downloads and removes any leftovers
func (fs *Decomposedfs) PurgeExpiredUploads(purgedChan chan<- tusd.FileInfo) error {
	infos, err := fs.uploadInfos(context.Background())
	if err != nil {
		return err
	}

	for _, info := range infos {
		expires, err := strconv.Atoi(info.MetaData["expires"])
		if err != nil {
			continue
		}
		if int64(expires) < time.Now().Unix() {
			purgedChan <- info
			err = os.Remove(info.Storage["BinPath"])
			if err != nil {
				return err
			}
			err = os.Remove(filepath.Join(fs.o.Root, "uploads", info.ID+".info"))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AsTerminatableUpload returns a TerminatableUpload
// To implement the termination extension as specified in https://tus.io/protocols/resumable-upload.html#termination
// the storage needs to implement AsTerminatableUpload
func (fs *Decomposedfs) AsTerminatableUpload(up tusd.Upload) tusd.TerminatableUpload {
	return up.(*upload.Upload)
}

// AsLengthDeclarableUpload returns a LengthDeclarableUpload
// To implement the creation-defer-length extension as specified in https://tus.io/protocols/resumable-upload.html#creation
// the storage needs to implement AsLengthDeclarableUpload
func (fs *Decomposedfs) AsLengthDeclarableUpload(up tusd.Upload) tusd.LengthDeclarableUpload {
	return up.(*upload.Upload)
}

// AsConcatableUpload returns a ConcatableUpload
// To implement the concatenation extension as specified in https://tus.io/protocols/resumable-upload.html#concatenation
// the storage needs to implement AsConcatableUpload
func (fs *Decomposedfs) AsConcatableUpload(up tusd.Upload) tusd.ConcatableUpload {
	return up.(*upload.Upload)
}

func (fs *Decomposedfs) uploadInfos(ctx context.Context) ([]tusd.FileInfo, error) {
	infos := []tusd.FileInfo{}
	infoFiles, err := filepath.Glob(filepath.Join(fs.o.Root, "uploads", "*.info"))
	if err != nil {
		return nil, err
	}

	for _, info := range infoFiles {
		match := _idRegexp.FindStringSubmatch(info)
		if match == nil || len(match) < 2 {
			continue
		}
		up, err := fs.GetUpload(ctx, match[1])
		if err != nil {
			return nil, err
		}
		info, err := up.GetInfo(context.Background())
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}
	return infos, nil
}

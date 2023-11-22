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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	tusd "github.com/tus/tusd/pkg/handler"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/grpc/services/storageprovider"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/upload"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

var _idRegexp = regexp.MustCompile(".*/([^/]+).info")

// Upload uploads data to the given resource
// TODO Upload (and InitiateUpload) needs a way to receive the expected checksum.
// Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	up, err := fs.GetUpload(ctx, req.Ref.GetPath())
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload")
	}

	uploadInfo := up.(*upload.Upload)

	p := uploadInfo.Info.Storage["NodeName"]
	if chunking.IsChunked(p) { // check chunking v1
		var assembledFile string
		p, assembledFile, err = fs.chunkHandler.WriteChunk(p, req.Body)
		if err != nil {
			return provider.ResourceInfo{}, err
		}
		if p == "" {
			if err = uploadInfo.Terminate(ctx); err != nil {
				return provider.ResourceInfo{}, errors.Wrap(err, "ocfs: error removing auxiliary files")
			}
			return provider.ResourceInfo{}, errtypes.PartialContent(req.Ref.String())
		}
		uploadInfo.Info.Storage["NodeName"] = p
		fd, err := os.Open(assembledFile)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)
		req.Body = fd
	}

	if _, err := uploadInfo.WriteChunk(ctx, 0, req.Body); err != nil {
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
		executant, ok := ctxpkg.ContextGetUser(uploadInfo.Ctx)
		if !ok {
			return provider.ResourceInfo{}, errtypes.PreconditionFailed("error getting user from uploadinfo context")
		}
		spaceOwner := &userpb.UserId{
			OpaqueId: info.Storage["SpaceOwnerOrManager"],
		}
		uff(spaceOwner, executant.Id, uploadRef)
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
// It creates a node for new files to persist the fileid for the new child.
// TODO read optional content for small files in this request
// TODO InitiateUpload (and Upload) needs a way to receive the expected checksum. Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
// TODO needs a way to handle unknown filesize, currently uses the context
// FIXME headers is actually used to carry all kinds of headers
func (fs *Decomposedfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, headers map[string]string) (map[string]string, error) {

	n, err := fs.lu.NodeFromResource(ctx, ref)
	switch err.(type) {
	case nil:
		// ok
	case errtypes.IsNotFound:
		return nil, errtypes.PreconditionFailed(err.Error())
	default:
		return nil, err
	}

	sublog := appctx.GetLogger(ctx).With().Str("spaceid", n.SpaceID).Str("nodeid", n.ID).Int64("uploadLength", uploadLength).Interface("headers", headers).Logger()

	// permissions are checked in NewUpload below

	relative, err := fs.lu.Path(ctx, n, node.NoCheck)
	if err != nil {
		return nil, err
	}

	usr := ctxpkg.ContextMustGetUser(ctx)
	uploadMetadata := upload.Metadata{
		Filename:            n.Name,
		SpaceRoot:           n.SpaceRoot.ID,
		SpaceOwnerOrManager: n.SpaceOwnerOrManager(ctx).GetOpaqueId(),
		ProviderID:          headers["providerID"],
		MTime:               time.Now().UTC().Format(time.RFC3339Nano),
		NodeID:              n.ID,
		NodeParentID:        n.ParentID,
		ExecutantIdp:        usr.Id.Idp,
		ExecutantID:         usr.Id.OpaqueId,
		ExecutantType:       utils.UserTypeToString(usr.Id.Type),
		ExecutantUserName:   usr.Username,
		LogLevel:            sublog.GetLevel().String(),
	}

	tusMetadata := tusd.MetaData{}

	// checksum is sent as tus Upload-Checksum header and should not magically become a metadata property
	if checksum, ok := headers["checksum"]; ok {
		parts := strings.SplitN(checksum, " ", 2)
		if len(parts) != 2 {
			return nil, errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		switch parts[0] {
		case "sha1", "md5", "adler32":
			uploadMetadata.Checksum = checksum
		default:
			return nil, errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
	}

	// if mtime has been set via tus metadata, expose it as tus metadata
	if ocmtime, ok := headers["mtime"]; ok {
		if ocmtime != "null" {
			tusMetadata["mtime"] = ocmtime
			// overwrite mtime if requested
			mtime, err := utils.MTimeToTime(ocmtime)
			if err != nil {
				return nil, err
			}
			uploadMetadata.MTime = mtime.UTC().Format(time.RFC3339Nano)
		}
	}

	_, err = node.CheckQuota(ctx, n.SpaceRoot, n.Exists, uint64(n.Blobsize), uint64(uploadLength))
	if err != nil {
		return nil, err
	}

	// check permissions
	var (
		checkNode *node.Node
		path      string
	)
	if n.Exists {
		// check permissions of file to be overwritten
		checkNode = n
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}})
	} else {
		// check permissions of parent
		parent, perr := n.Parent(ctx)
		if perr != nil {
			return nil, errors.Wrap(perr, "Decomposedfs: error getting parent "+n.ParentID)
		}
		checkNode = parent
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}, Path: n.Name})
	}
	rp, err := fs.p.AssemblePermissions(ctx, checkNode) // context does not have a user?
	switch {
	case err != nil:
		return nil, err
	case !rp.InitiateFileUpload:
		return nil, errtypes.PermissionDenied(path)
	}

	// are we trying to overwrite a folder with a file?
	if n.Exists && n.IsDir(ctx) {
		return nil, errtypes.PreconditionFailed("resource is not a file")
	}

	// check lock
	// FIXME we cannot check the lock of a new file, because it would have to use the name ...
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	// treat 0 length uploads as deferred
	sizeIsDeferred := false
	if uploadLength == 0 {
		sizeIsDeferred = true
	}

	info := tusd.FileInfo{
		MetaData:       tusMetadata,
		Size:           uploadLength,
		SizeIsDeferred: sizeIsDeferred,
	}
	if lockID, ok := ctxpkg.ContextGetLockID(ctx); ok {
		uploadMetadata.LockID = lockID
	}
	uploadMetadata.Dir = filepath.Dir(relative)

	// rewrite filename for old chunking v1
	if chunking.IsChunked(n.Name) {
		uploadMetadata.Chunk = n.Name
		bi, err := chunking.GetChunkBLOBInfo(n.Name)
		if err != nil {
			return nil, err
		}
		n.Name = bi.Path
	}

	// TODO at this point we have no way to figure out the output or mode of the logger. we need that to reinitialize a logger in PreFinishResponseCallback
	// or better create a config option for the log level during PreFinishResponseCallback? might be easier for now

	// expires has been set by the storageprovider, do not expose as metadata. It is sent as a tus Upload-Expires header
	if expiration, ok := headers["expires"]; ok {
		if expiration != "null" { // TODO this is set by the storageprovider ... it cannot be set by cliensts, so it can never be the string 'null' ... or can it???
			uploadMetadata.Expires = expiration
		}
	}
	// only check preconditions if they are not empty
	// do not expose as metadata
	if headers["if-match"] != "" {
		uploadMetadata.HeaderIfMatch = headers["if-match"] // TODO drop?
	}
	if headers["if-none-match"] != "" {
		uploadMetadata.HeaderIfNoneMatch = headers["if-none-match"]
	}
	if headers["if-unmodified-since"] != "" {
		uploadMetadata.HeaderIfUnmodifiedSince = headers["if-unmodified-since"]
	}

	if uploadMetadata.HeaderIfNoneMatch == "*" && n.Exists {
		return nil, errtypes.Aborted(fmt.Sprintf("parent %s already has a child %s", n.ID, n.Name))
	}

	// create the upload
	u, err := fs.tusDataStore.NewUpload(ctx, info)
	if err != nil {
		return nil, err
	}

	info, err = u.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	uploadMetadata.ID = info.ID

	// keep track of upload
	err = upload.WriteMetadata(ctx, fs.lu, info.ID, uploadMetadata)
	if err != nil {
		return nil, err
	}

	sublog.Debug().Interface("info", info).Msg("Decomposedfs: initiated upload")

	return map[string]string{
		"simple": info.ID,
		"tus":    info.ID,
	}, nil
}

// GetDataStore returns the initialized Datastore
func (fs *Decomposedfs) GetDataStore() tusd.DataStore {
	return fs.tusDataStore
}

// PreFinishResponseCallback is called by the tus datatx, after all bytes have been transferred
func (fs *Decomposedfs) PreFinishResponseCallback(hook tusd.HookEvent) error {
	ctx := context.TODO()
	appctx.GetLogger(ctx).Debug().Interface("hook", hook).Msg("got PreFinishResponseCallback")
	ctx, span := tracer.Start(ctx, "PreFinishResponseCallback")
	defer span.End()

// NewUpload returns a new tus Upload instance
func (fs *Decomposedfs) NewUpload(ctx context.Context, info tusd.FileInfo) (tusd.Upload, error) {
	return upload.New(ctx, info, fs.lu, fs.tp, fs.p, fs.o.Root, fs.stream, fs.o.AsyncFileUploads, fs.o.Tokens)
}

// GetUpload returns the Upload for the given upload id
func (fs *Decomposedfs) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	return upload.Get(ctx, id, fs.lu, fs.tp, fs.o.Root, fs.stream, fs.o.AsyncFileUploads, fs.o.Tokens)
}

// ListUploadSessions returns the upload sessions for the given filter
func (fs *Decomposedfs) ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error) {
	var sessions []storage.UploadSession
	if filter.ID != nil && *filter.ID != "" {
		session, err := fs.getUploadSession(ctx, filepath.Join(fs.o.Root, "uploads", *filter.ID+".info"))
		if err != nil {
			return nil, err
		}
		sessions = []storage.UploadSession{session}
	} else {
		var err error
		sessions, err = fs.uploadSessions(ctx)
		if err != nil {
			return nil, err
		}
	}
	filteredSessions := []storage.UploadSession{}
	now := time.Now()
	for _, session := range sessions {
		if filter.Processing != nil && *filter.Processing != session.IsProcessing() {
			continue
		}
		if filter.Expired != nil {
			if *filter.Expired {
				if now.Before(session.Expires()) {
					continue
				}
			} else {
				if now.After(session.Expires()) {
					continue
				}
			}
		}
		filteredSessions = append(filteredSessions, session)
	}
	return filteredSessions, nil
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

func (fs *Decomposedfs) uploadSessions(ctx context.Context) ([]storage.UploadSession, error) {
	uploads := []storage.UploadSession{}
	infoFiles, err := filepath.Glob(filepath.Join(fs.o.Root, "uploads", "*.info"))
	if err != nil {
		return nil, err
	}

	for _, info := range infoFiles {
		progress, err := fs.getUploadSession(ctx, info)
		if err != nil {
			appctx.GetLogger(ctx).Error().Interface("path", info).Msg("Decomposedfs: could not getUploadSession")
			continue
		}

		uploads = append(uploads, progress)
	}
	return uploads, nil
}

func (fs *Decomposedfs) getUploadSession(ctx context.Context, path string) (storage.UploadSession, error) {
	match := _idRegexp.FindStringSubmatch(path)
	if match == nil || len(match) < 2 {
		return nil, fmt.Errorf("invalid upload path")
	}
	up, err := fs.GetUpload(ctx, match[1])
	if err != nil {
		return nil, err
	}
	info, err := up.GetInfo(context.Background())
	if err != nil {
		return nil, err
	}
	// upload processing state is stored in the node, for decomposedfs the NodeId is always set by InitiateUpload
	n, err := node.ReadNode(ctx, fs.lu, info.Storage["SpaceRoot"], info.Storage["NodeId"], true, nil, true)
	if err != nil {
		return nil, err
	}
	progress := upload.Progress{
		Path:       path,
		Info:       info,
		Processing: n.IsProcessing(ctx),
	}
	return progress, nil
}

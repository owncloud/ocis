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
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	"net/url"
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
	"golang.org/x/sync/errgroup"
)

var _idRegexp = regexp.MustCompile(".*/([^/]+).mpk")

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

	// if mtime has been set via the headers, expose it as tus metadata
	if ocmtime, ok := headers["mtime"]; ok && ocmtime != "null" {
		tusMetadata["mtime"] = ocmtime
		// overwrite mtime if requested
		mtime, err := utils.MTimeToTime(ocmtime)
		if err != nil {
			return nil, err
		}
		uploadMetadata.MTime = mtime.UTC().Format(time.RFC3339Nano)
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

	info := tusd.FileInfo{
		MetaData:       tusMetadata,
		Size:           uploadLength,
		SizeIsDeferred: uploadLength == 0, // treat 0 length uploads as deferred
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
	if expiration, ok := headers["expires"]; ok && expiration != "null" { // TODO this is set by the storageprovider ... it cannot be set by cliensts, so it can never be the string 'null' ... or can it???
		uploadMetadata.Expires, err = utils.MTimeToTime(expiration)
		if err != nil {
			return nil, err
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

	info := hook.Upload
	up, err := fs.tusDataStore.GetUpload(ctx, info.ID)
	if err != nil {
		return err
	}

	uploadMetadata, err := upload.ReadMetadata(ctx, fs.lu, info.ID)
	if err != nil {
		return err
	}

	// put lockID from upload back into context
	if uploadMetadata.LockID != "" {
		ctx = ctxpkg.ContextSetLockID(ctx, uploadMetadata.LockID)
	}

	// restore logger from file info
	log, err := logger.FromConfig(&logger.LogConf{
		Output: "stderr",
		Mode:   "json",
		Level:  uploadMetadata.LogLevel,
	})
	if err != nil {
		return err
	}

	ctx = appctx.WithLogger(ctx, log)

	// calculate the checksum of the written bytes
	// they will all be written to the metadata later, so we cannot omit any of them
	// TODO only calculate the checksum in sync that was requested to match, the rest could be async ... but the tests currently expect all to be present
	// TODO the hashes all implement BinaryMarshaler so we could try to persist the state for resumable upload. we would neet do keep track of the copied bytes ...

	log.Debug().Msg("calculating checksums")
	sha1h := sha1.New()
	md5h := md5.New()
	adler32h := adler32.New()
	{
		_, subspan := tracer.Start(ctx, "GetReader")
		reader, err := up.GetReader(ctx)
		subspan.End()
		if err != nil {
			// we can continue if no oc checksum header is set
			log.Info().Err(err).Interface("info", info).Msg("error getting Reader from upload")
		}
		if readCloser, ok := reader.(io.ReadCloser); ok {
			defer readCloser.Close()
		}

		r1 := io.TeeReader(reader, sha1h)
		r2 := io.TeeReader(r1, md5h)

		_, subspan = tracer.Start(ctx, "io.Copy")
		/*bytesCopied*/ _, err = io.Copy(adler32h, r2)
		subspan.End()
		if err != nil {
			log.Info().Err(err).Msg("error copying checksums")
		}
		/*
			if bytesCopied != info.Size {
				msg := fmt.Sprintf("mismatching upload length. expected %d, could only copy %d", info.Size, bytesCopied)
				log.Error().Interface("info", info).Msg(msg)
				return errtypes.InternalError(msg)
			}
		*/
	}

	// compare if they match the sent checksum
	// TODO the tus checksum extension would do this on every chunk, but I currently don't see an easy way to pass in the requested checksum. for now we do it in FinishUpload which is also called for chunked uploads
	if uploadMetadata.Checksum != "" {
		var err error
		parts := strings.SplitN(uploadMetadata.Checksum, " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		switch parts[0] {
		case "sha1":
			err = checkHash(parts[1], sha1h)
		case "md5":
			err = checkHash(parts[1], md5h)
		case "adler32":
			err = checkHash(parts[1], adler32h)
		default:
			err = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if err != nil {
			if tup, ok := up.(tusd.TerminatableUpload); ok {
				terr := tup.Terminate(ctx)
				if terr != nil {
					log.Error().Err(terr).Interface("info", info).Msg("failed to terminate upload")
				}
			}
			return err
		}
	}

	// update checksums
	uploadMetadata.ChecksumSHA1 = sha1h.Sum(nil)
	uploadMetadata.ChecksumMD5 = md5h.Sum(nil)
	uploadMetadata.ChecksumADLER32 = adler32h.Sum(nil)

	log.Debug().Str("id", info.ID).Msg("upload.UpdateMetadata")
	uploadMetadata, n, err := upload.UpdateMetadata(ctx, fs.lu, info.ID, info.Size, uploadMetadata)
	if err != nil {
		upload.Cleanup(ctx, fs.lu, n, info.ID, uploadMetadata.MTime, true)
		if tup, ok := up.(tusd.TerminatableUpload); ok {
			terr := tup.Terminate(ctx)
			if terr != nil {
				log.Error().Err(terr).Interface("info", info).Msg("failed to terminate upload")
			}
		}
		return err
	}

	if fs.stream != nil {
		user := &userpb.User{
			Id: &userpb.UserId{
				Type:     userpb.UserType(userpb.UserType_value[uploadMetadata.ExecutantType]),
				Idp:      uploadMetadata.ExecutantIdp,
				OpaqueId: uploadMetadata.ExecutantID,
			},
			Username: uploadMetadata.ExecutantUserName,
		}
		s, err := fs.downloadURL(ctx, info.ID)
		if err != nil {
			return err
		}

		log.Debug().Str("id", info.ID).Msg("events.Publish BytesReceived")
		if err := events.Publish(ctx, fs.stream, events.BytesReceived{
			UploadID:      info.ID,
			URL:           s,
			SpaceOwner:    n.SpaceOwnerOrManager(ctx),
			ExecutingUser: user,
			ResourceID:    &provider.ResourceId{SpaceId: n.SpaceID, OpaqueId: n.ID},
			Filename:      uploadMetadata.Filename, // TODO what and when do we publish chunking v2 names? Currently, this uses the chunk name.
			Filesize:      uint64(info.Size),
		}); err != nil {
			return err
		}
	}

	if n.Exists {
		// // copy metadata to a revision node
		log.Debug().Str("id", info.ID).Msg("copy metadata to a revision node")
		currentAttrs, err := n.Xattrs(ctx)
		if err != nil {
			return err
		}
		previousRevisionTime, err := n.GetMTime(ctx)
		if err != nil {
			return err
		}
		rm := upload.RevisionMetadata{
			MTime:           previousRevisionTime.UTC().Format(time.RFC3339Nano),
			BlobID:          n.BlobID,
			BlobSize:        n.Blobsize,
			ChecksumSHA1:    currentAttrs[prefixes.ChecksumPrefix+storageprovider.XSSHA1],
			ChecksumMD5:     currentAttrs[prefixes.ChecksumPrefix+storageprovider.XSMD5],
			ChecksumADLER32: currentAttrs[prefixes.ChecksumPrefix+storageprovider.XSAdler32],
		}
		revisionNode := n.RevisionNode(ctx, rm.MTime)

		rh, err := upload.CreateRevisionNode(ctx, fs.lu, revisionNode)
		if err != nil {
			return err
		}
		defer rh.Close()
		err = upload.WriteRevisionMetadataToNode(ctx, revisionNode, rm)
		if err != nil {
			return err
		}
	}

	sizeDiff := info.Size - n.Blobsize
	if !fs.o.AsyncFileUploads {
		// handle postprocessing synchronously
		log.Debug().Str("id", info.ID).Msg("upload.Finalize")
		err = upload.Finalize(ctx, fs.blobstore, uploadMetadata.MTime, info, n, uploadMetadata.BlobID) // moving or copying the blob only reads the blobid, no need to change the revision nodes nodeid

		log.Debug().Str("id", info.ID).Msg("upload.Cleanup")
		upload.Cleanup(ctx, fs.lu, n, info.ID, uploadMetadata.MTime, err != nil)
		if tup, ok := up.(tusd.TerminatableUpload); ok {
			log.Debug().Str("id", info.ID).Msg("tup.Terminate")
			terr := tup.Terminate(ctx)
			if terr != nil {
				log.Error().Err(terr).Interface("info", info).Msg("failed to terminate upload")
			}
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to upload")
			return err
		}
		log.Debug().Str("id", info.ID).Msg("upload.SetNodeToUpload")
		sizeDiff, err = upload.SetNodeToUpload(ctx, fs.lu, n, uploadMetadata)
		if err != nil {
			log.Error().Err(err).Msg("failed update Node to revision")
			return err
		}
	}
	log.Debug().Str("id", info.ID).Msg("fs.tp.Propagate")
	return fs.tp.Propagate(ctx, n, sizeDiff)
}

// URL returns a url to download an upload
func (fs *Decomposedfs) downloadURL(_ context.Context, id string) (string, error) {
	type transferClaims struct {
		jwt.StandardClaims
		Target string `json:"target"`
	}

	u, err := url.JoinPath(fs.o.Tokens.DownloadEndpoint, "tus/", id)
	if err != nil {
		return "", errors.Wrapf(err, "error joinging URL path")
	}
	ttl := time.Duration(fs.o.Tokens.TransferExpires) * time.Second
	claims := transferClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			Audience:  "reva",
			IssuedAt:  time.Now().Unix(),
		},
		Target: u,
	}

	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)

	tkn, err := t.SignedString([]byte(fs.o.Tokens.TransferSharedSecret))
	if err != nil {
		return "", errors.Wrapf(err, "error signing token with claims %+v", claims)
	}

	return url.JoinPath(fs.o.Tokens.DataGatewayEndpoint, tkn)
}

func checkHash(expected string, h hash.Hash) error {
	if expected != hex.EncodeToString(h.Sum(nil)) {
		return errtypes.ChecksumMismatch(fmt.Sprintf("invalid checksum: expected %s got %x", expected, h.Sum(nil)))
	}
	return nil
}

// Upload uploads data to the given resource
// is used by the simple datatx, after an InitiateUpload call
// TODO Upload (and InitiateUpload) needs a way to receive the expected checksum.
// Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) Upload(ctx context.Context, req storage.UploadRequest, uff storage.UploadFinishedFunc) (provider.ResourceInfo, error) {
	sublog := appctx.GetLogger(ctx).With().Str("path", req.Ref.Path).Int64("uploadLength", req.Length).Logger()
	up, err := fs.tusDataStore.GetUpload(ctx, req.Ref.GetPath())
	if err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: error retrieving upload")
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload")
	}

	uploadInfo, err := up.GetInfo(ctx)
	if err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: error retrieving upload info")
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload info")
	}

	uploadMetadata, err := upload.ReadMetadata(ctx, fs.lu, uploadInfo.ID)
	if err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: error retrieving upload metadata")
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error retrieving upload metadata")
	}

	if chunking.IsChunked(uploadMetadata.Chunk) { // check chunking v1, TODO, actually there is a 'OC-Chunked: 1' header, at least when the testsuite uses chunking v1
		var assembledFile, p string
		p, assembledFile, err = fs.chunkHandler.WriteChunk(uploadMetadata.Chunk, req.Body)
		if err != nil {
			sublog.Debug().Err(err).Msg("Decomposedfs: could not write chunk")
			return provider.ResourceInfo{}, err
		}
		if p == "" {
			sublog.Debug().Err(err).Str("chunk", uploadMetadata.Chunk).Msg("Decomposedfs: wrote chunk")
			return provider.ResourceInfo{}, errtypes.PartialContent(req.Ref.String())
		}
		uploadMetadata.Filename = p
		fd, err := os.Open(assembledFile)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)

		chunkStat, err := fd.Stat()
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: could not stat assembledFile for legacy chunking")
		}

		// fake a new upload with the correct size
		newInfo := tusd.FileInfo{
			Size:     chunkStat.Size(),
			MetaData: uploadInfo.MetaData,
		}
		nup, err := fs.tusDataStore.NewUpload(ctx, newInfo)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: could not create new tus upload for legacy chunking")
		}
		newInfo, err = nup.GetInfo(ctx)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: could not get info from upload")
		}
		uploadMetadata.ID = newInfo.ID
		uploadMetadata.BlobSize = newInfo.Size
		err = upload.WriteMetadata(ctx, fs.lu, newInfo.ID, uploadMetadata)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing upload metadata for legacy chunking")
		}

		_, err = nup.WriteChunk(ctx, 0, fd)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing to binary file for legacy chunking")
		}
		// use new upload and info
		up = nup
		uploadInfo, err = up.GetInfo(ctx)
		if err != nil {
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: could not get info for legacy chunking")
		}
	} else {
		// we need to call up.DeclareLength() before writing the chunk, but only if we actually got a length
		if req.Length > 0 {
			if ldx, ok := up.(tusd.LengthDeclarableUpload); ok {
				if err := ldx.DeclareLength(ctx, req.Length); err != nil {
					sublog.Debug().Err(err).Msg("Decomposedfs: error declaring length")
					return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error declaring length")
				}
			}
		}
		bytesWritten, err := up.WriteChunk(ctx, 0, req.Body)
		if err != nil {
			sublog.Debug().Err(err).Msg("Decomposedfs: error writing to binary file")
			return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error writing to binary file")
		}
		uploadInfo.Offset += bytesWritten
		if uploadInfo.SizeIsDeferred {
			// update the size and offset
			uploadInfo.SizeIsDeferred = false
			uploadInfo.Size = bytesWritten
		}
	}

	// This finishes the tus upload
	sublog.Debug().Msg("finishing upload")
	if err := up.FinishUpload(ctx); err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: error finishing upload")
		return provider.ResourceInfo{}, err
	}

	// we now need to handle to move/copy&delete to the target blobstore
	sublog.Debug().Msg("executing tusd prefinish callback")
	err = fs.PreFinishResponseCallback(tusd.HookEvent{Upload: uploadInfo})
	if err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: tusd callback failed")
		return provider.ResourceInfo{}, err
	}

	n, err := upload.ReadNode(ctx, fs.lu, uploadMetadata)
	if err != nil {
		sublog.Debug().Err(err).Msg("Decomposedfs: error reading node")
		return provider.ResourceInfo{}, err
	}

	if uff != nil {
		// TODO search needs to index the full path, so we return a reference relative to the space root.
		// but then the search has to walk the path. it might be more efficient if search called GetPath itself ... or we send the path as additional metadata in the event
		uploadRef := &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: uploadMetadata.ProviderID,
				SpaceId:   n.SpaceID,
				OpaqueId:  n.SpaceID,
			},
			Path: utils.MakeRelativePath(filepath.Join(uploadMetadata.Dir, uploadMetadata.Filename)),
		}
		executant, ok := ctxpkg.ContextGetUser(ctx)
		if !ok {
			return provider.ResourceInfo{}, errtypes.PreconditionFailed("error getting user from context")
		}

		sublog.Debug().Msg("calling upload finished func")
		uff(n.SpaceOwnerOrManager(ctx), executant.Id, uploadRef)
	}

	mtime, err := n.GetMTime(ctx)
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error getting mtime for '"+n.ID+"'")
	}
	etag, err := node.CalculateEtag(n, mtime)
	if err != nil {
		return provider.ResourceInfo{}, errors.Wrap(err, "Decomposedfs: error calculating etag '"+n.ID+"'")
	}
	ri := provider.ResourceInfo{
		// fill with at least fileid, mtime and etag
		Id: &provider.ResourceId{
			StorageId: uploadMetadata.ProviderID,
			SpaceId:   n.SpaceID,
			OpaqueId:  n.ID,
		},
		Etag: etag,
	}

	if mtime, err := utils.MTimeToTS(uploadInfo.MetaData["mtime"]); err == nil {
		ri.Mtime = &mtime
	}
	sublog.Debug().Msg("Decomposedfs: finished upload")

	return ri, nil
}

// ListUploadSessions returns the upload sessions for the given filter
func (fs *Decomposedfs) ListUploadSessions(ctx context.Context, filter storage.UploadSessionFilter) ([]storage.UploadSession, error) {
	var sessions []storage.UploadSession
	if filter.ID != nil && *filter.ID != "" {
		session, err := fs.getUploadSession(ctx, fs.lu.UploadPath("*"))
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

func (fs *Decomposedfs) uploadSessions(ctx context.Context) ([]storage.UploadSession, error) {
	uploads := []storage.UploadSession{}
	sessionFiles, err := filepath.Glob(fs.lu.UploadPath("*"))
	if err != nil {
		return nil, err
	}

	numWorkers := fs.o.MaxConcurrency
	if len(sessionFiles) < numWorkers {
		numWorkers = len(sessionFiles)
	}

	work := make(chan string, len(sessionFiles))
	results := make(chan storage.UploadSession, len(sessionFiles))

	g, ctx := errgroup.WithContext(ctx)

	// Distribute work
	g.Go(func() error {
		defer close(work)
		for _, itemPath := range sessionFiles {
			select {
			case work <- itemPath:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Spawn workers that'll concurrently work the queue
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			for path := range work {
				session, err := fs.getUploadSession(ctx, path)
				if err != nil {
					appctx.GetLogger(ctx).Error().Interface("path", path).Msg("Decomposedfs: could not getUploadSession")
					continue
				}

				select {
				case results <- session:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})
	}

	// Wait for things to settle down, then close results chan
	go func() {
		_ = g.Wait() // error is checked later
		close(results)
	}()

	// Collect results
	for ri := range results {
		uploads = append(uploads, ri)
	}
	return uploads, nil
}

func (fs *Decomposedfs) getUploadSession(ctx context.Context, path string) (storage.UploadSession, error) {
	match := _idRegexp.FindStringSubmatch(path)
	if match == nil || len(match) < 2 {
		return nil, fmt.Errorf("invalid upload path")
	}

	metadata, err := upload.ReadMetadata(ctx, fs.lu, match[1])
	if err != nil {
		return nil, err
	}
	// upload processing state is stored in the node, for decomposedfs the NodeId is always set by InitiateUpload
	var n *node.Node
	if metadata.NodeID == "" {
		// read parent first
		n, err = node.ReadNode(ctx, fs.lu, metadata.SpaceRoot, metadata.NodeParentID, true, nil, true)
		if err != nil {
			return nil, err
		}
		n, err = n.Child(ctx, metadata.Filename)
	} else {
		n, err = node.ReadNode(ctx, fs.lu, metadata.SpaceRoot, metadata.NodeID, true, nil, true)
	}
	if err != nil {
		return nil, err
	}
	tusUpload, err := fs.tusDataStore.GetUpload(ctx, metadata.ID)
	if err != nil {
		return nil, err
	}

	progress := upload.Progress{
		Upload:     tusUpload,
		Path:       path,
		Metadata:   metadata,
		Processing: n.IsProcessing(ctx),
	}

	return progress, nil
}

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
	"encoding/json"
	"fmt"
	"hash"
	"hash/adler32"
	"io"
	iofs "io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/xattrs"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tusd "github.com/tus/tusd/pkg/handler"
)

var defaultFilePerm = os.FileMode(0664)

// Upload uploads data to the given resource
// TODO Upload (and InitiateUpload) needs a way to receive the expected checksum.
// Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) Upload(ctx context.Context, ref *provider.Reference, r io.ReadCloser, uff storage.UploadFinishedFunc) error {
	upload, err := fs.GetUpload(ctx, ref.GetPath())
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error retrieving upload")
	}

	uploadInfo := upload.(*fileUpload)

	p := uploadInfo.info.Storage["NodeName"]
	if chunking.IsChunked(p) { // check chunking v1
		var assembledFile string
		p, assembledFile, err = fs.chunkHandler.WriteChunk(p, r)
		if err != nil {
			return err
		}
		if p == "" {
			if err = uploadInfo.Terminate(ctx); err != nil {
				return errors.Wrap(err, "ocfs: error removing auxiliary files")
			}
			return errtypes.PartialContent(ref.String())
		}
		uploadInfo.info.Storage["NodeName"] = p
		fd, err := os.Open(assembledFile)
		if err != nil {
			return errors.Wrap(err, "Decomposedfs: error opening assembled file")
		}
		defer fd.Close()
		defer os.RemoveAll(assembledFile)
		r = fd
	}

	if _, err := uploadInfo.WriteChunk(ctx, 0, r); err != nil {
		return errors.Wrap(err, "Decomposedfs: error writing to binary file")
	}

	if err := uploadInfo.FinishUpload(ctx); err != nil {
		return err
	}

	if uff != nil {
		info := uploadInfo.info
		uploadRef := &provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: storagespace.FormatStorageID(info.MetaData["providerID"], info.Storage["SpaceRoot"]),
				OpaqueId:  info.Storage["SpaceRoot"],
			},
			Path: utils.MakeRelativePath(filepath.Join(info.MetaData["dir"], info.MetaData["filename"])),
		}
		owner, ok := ctxpkg.ContextGetUser(uploadInfo.ctx)
		if !ok {
			return errtypes.PreconditionFailed("error getting user from uploadinfo context")
		}
		uff(owner.Id, uploadRef)
	}

	return nil
}

// InitiateUpload returns upload ids corresponding to different protocols it supports
// TODO read optional content for small files in this request
// TODO InitiateUpload (and Upload) needs a way to receive the expected checksum. Maybe in metadata as 'checksum' => 'sha1 aeosvp45w5xaeoe' = lowercase, space separated?
func (fs *Decomposedfs) InitiateUpload(ctx context.Context, ref *provider.Reference, uploadLength int64, metadata map[string]string) (map[string]string, error) {
	log := appctx.GetLogger(ctx)

	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		return nil, err
	}

	// permissions are checked in NewUpload below

	relative, err := fs.lu.Path(ctx, n)
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
			"SpaceRoot": n.SpaceRoot.ID,
		},
	}

	if metadata != nil {
		info.MetaData["providerID"] = metadata["providerID"]
		if mtime, ok := metadata["mtime"]; ok {
			info.MetaData["mtime"] = mtime
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
func (fs *Decomposedfs) NewUpload(ctx context.Context, info tusd.FileInfo) (upload tusd.Upload, err error) {

	log := appctx.GetLogger(ctx)
	log.Debug().Interface("info", info).Msg("Decomposedfs: NewUpload")

	if info.MetaData["filename"] == "" {
		return nil, errors.New("Decomposedfs: missing filename in metadata")
	}
	if info.MetaData["dir"] == "" {
		return nil, errors.New("Decomposedfs: missing dir in metadata")
	}

	n, err := fs.lu.NodeFromSpaceID(ctx, &provider.ResourceId{
		StorageId: info.Storage["SpaceRoot"],
	})
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error getting space root node")
	}

	n, err = fs.lookupNode(ctx, n, filepath.Join(info.MetaData["dir"], info.MetaData["filename"]))
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	log.Debug().Interface("info", info).Interface("node", n).Msg("Decomposedfs: resolved filename")

	// the parent owner will become the new owner
	p, perr := n.Parent()
	if perr != nil {
		return nil, errors.Wrap(perr, "Decomposedfs: error getting parent "+n.ParentID)
	}

	// check permissions
	var ok bool
	if n.Exists {
		// check permissions of file to be overwritten
		ok, err = fs.p.HasPermission(ctx, n, func(rp *provider.ResourcePermissions) bool {
			return rp.InitiateFileUpload
		})
	} else {
		// check permissions of parent
		ok, err = fs.p.HasPermission(ctx, p, func(rp *provider.ResourcePermissions) bool {
			return rp.InitiateFileUpload
		})
	}
	switch {
	case err != nil:
		return nil, errtypes.InternalError(err.Error())
	case !ok:
		return nil, errtypes.PermissionDenied(filepath.Join(n.ParentID, n.Name))
	}

	// if we are trying to overwriting a folder with a file
	if n.Exists && n.IsDir() {
		return nil, errtypes.PreconditionFailed("resource is not a file")
	}

	// check lock
	if info.MetaData["lockid"] != "" {
		ctx = ctxpkg.ContextSetLockID(ctx, info.MetaData["lockid"])
	}
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	info.ID = uuid.New().String()

	binPath, err := fs.getUploadPath(ctx, info.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error resolving upload path")
	}
	usr := ctxpkg.ContextMustGetUser(ctx)

	var spaceRoot string
	if info.Storage != nil {
		if spaceRoot, ok = info.Storage["SpaceRoot"]; !ok {
			spaceRoot = n.SpaceRoot.ID
		}
	} else {
		spaceRoot = n.SpaceRoot.ID
	}

	info.Storage = map[string]string{
		"Type":    "OCISStore",
		"BinPath": binPath,

		"NodeId":       n.ID,
		"NodeParentId": n.ParentID,
		"NodeName":     n.Name,
		"SpaceRoot":    spaceRoot,

		"Idp":      usr.Id.Idp,
		"UserId":   usr.Id.OpaqueId,
		"UserType": utils.UserTypeToString(usr.Id.Type),
		"UserName": usr.Username,

		"LogLevel": log.GetLevel().String(),
	}
	// Create binary file in the upload folder with no content
	log.Debug().Interface("info", info).Msg("Decomposedfs: built storage info")
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	u := &fileUpload{
		info:     info,
		binPath:  binPath,
		infoPath: filepath.Join(fs.o.Root, "uploads", info.ID+".info"),
		fs:       fs,
		ctx:      ctx,
	}

	// writeInfo creates the file by itself if necessary
	err = u.writeInfo()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (fs *Decomposedfs) getUploadPath(ctx context.Context, uploadID string) (string, error) {
	return filepath.Join(fs.o.Root, "uploads", uploadID), nil
}

// GetUpload returns the Upload for the given upload id
func (fs *Decomposedfs) GetUpload(ctx context.Context, id string) (tusd.Upload, error) {
	infoPath := filepath.Join(fs.o.Root, "uploads", id+".info")

	info := tusd.FileInfo{}
	data, err := ioutil.ReadFile(infoPath)
	if err != nil {
		if errors.Is(err, iofs.ErrNotExist) {
			// Interpret os.ErrNotExist as 404 Not Found
			err = tusd.ErrNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	stat, err := os.Stat(info.Storage["BinPath"])
	if err != nil {
		return nil, err
	}

	info.Offset = stat.Size()

	u := &userpb.User{
		Id: &userpb.UserId{
			Idp:      info.Storage["Idp"],
			OpaqueId: info.Storage["UserId"],
			Type:     utils.UserTypeMap(info.Storage["UserType"]),
		},
		Username: info.Storage["UserName"],
	}

	ctx = ctxpkg.ContextSetUser(ctx, u)
	// TODO configure the logger the same way ... store and add traceid in file info

	var opts []logger.Option
	opts = append(opts, logger.WithLevel(info.Storage["LogLevel"]))
	opts = append(opts, logger.WithWriter(os.Stderr, logger.ConsoleMode))
	l := logger.New(opts...)

	sub := l.With().Int("pid", os.Getpid()).Logger()

	ctx = appctx.WithLogger(ctx, &sub)

	return &fileUpload{
		info:     info,
		binPath:  info.Storage["BinPath"],
		infoPath: infoPath,
		fs:       fs,
		ctx:      ctx,
	}, nil
}

// lookupNode looks up nodes by path.
// This method can also handle lookups for paths which contain chunking information.
func (fs *Decomposedfs) lookupNode(ctx context.Context, spaceRoot *node.Node, path string) (*node.Node, error) {
	p := path
	isChunked := chunking.IsChunked(path)
	if isChunked {
		chunkInfo, err := chunking.GetChunkBLOBInfo(path)
		if err != nil {
			return nil, err
		}
		p = chunkInfo.Path
	}

	n, err := fs.lu.WalkPath(ctx, spaceRoot, p, true, func(ctx context.Context, n *node.Node) error { return nil })
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	if isChunked {
		n.Name = filepath.Base(path)
	}
	return n, nil
}

type fileUpload struct {
	// info stores the current information about the upload
	info tusd.FileInfo
	// infoPath is the path to the .info file
	infoPath string
	// binPath is the path to the binary file (which has no extension)
	binPath string
	// only fs knows how to handle metadata and versions
	fs *Decomposedfs
	// a context with a user
	// TODO add logger as well?
	ctx context.Context
}

// GetInfo returns the FileInfo
func (upload *fileUpload) GetInfo(ctx context.Context) (tusd.FileInfo, error) {
	return upload.info, nil
}

// WriteChunk writes the stream from the reader to the given offset of the upload
func (upload *fileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// calculate cheksum here? needed for the TUS checksum extension. https://tus.io/protocols/resumable-upload.html#checksum
	// TODO but how do we get the `Upload-Checksum`? WriteChunk() only has a context, offset and the reader ...
	// It is sent with the PATCH request, well or in the POST when the creation-with-upload extension is used
	// but the tus handler uses a context.Background() so we cannot really check the header and put it in the context ...
	n, err := io.Copy(file, src)

	// If the HTTP PATCH request gets interrupted in the middle (e.g. because
	// the user wants to pause the upload), Go's net/http returns an io.ErrUnexpectedEOF.
	// However, for the ocis driver it's not important whether the stream has ended
	// on purpose or accidentally.
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return n, err
		}
	}

	upload.info.Offset += n
	err = upload.writeInfo() // TODO info is written here ... we need to truncate in DiscardChunk

	return n, err
}

// GetReader returns an io.Reader for the upload
func (upload *fileUpload) GetReader(ctx context.Context) (io.Reader, error) {
	return os.Open(upload.binPath)
}

// writeInfo updates the entire information. Everything will be overwritten.
func (upload *fileUpload) writeInfo() error {
	data, err := json.Marshal(upload.info)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(upload.infoPath, data, defaultFilePerm)
}

// FinishUpload finishes an upload and moves the file to the internal destination
func (upload *fileUpload) FinishUpload(ctx context.Context) (err error) {

	// ensure cleanup
	defer upload.discardChunk()

	fi, err := os.Stat(upload.binPath)
	if err != nil {
		appctx.GetLogger(upload.ctx).Err(err).Msg("Decomposedfs: could not stat uploaded file")
		return
	}

	spaceID := upload.info.Storage["SpaceRoot"]
	n := node.New(
		spaceID,
		upload.info.Storage["NodeId"],
		upload.info.Storage["NodeParentId"],
		upload.info.Storage["NodeName"],
		fi.Size(),
		"",
		nil,
		upload.fs.lu,
	)
	n.SpaceRoot = node.New(spaceID, spaceID, "", "", 0, "", nil, upload.fs.lu)

	// check lock
	if upload.info.MetaData["lockid"] != "" {
		ctx = ctxpkg.ContextSetLockID(ctx, upload.info.MetaData["lockid"])
	}
	if err := n.CheckLock(ctx); err != nil {
		return err
	}

	var oldSize uint64
	if n.ID != "" {
		old, _ := node.ReadNode(ctx, upload.fs.lu, spaceID, n.ID, false)
		oldSize = uint64(old.Blobsize)
	}
	_, err = node.CheckQuota(n.SpaceRoot, n.ID != "", oldSize, uint64(fi.Size()))

	if err != nil {
		return err
	}

	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	targetPath := n.InternalPath()
	sublog := appctx.GetLogger(upload.ctx).
		With().
		Interface("info", upload.info).
		Str("binPath", upload.binPath).
		Str("targetPath", targetPath).
		Logger()

	// calculate the checksum of the written bytes
	// they will all be written to the metadata later, so we cannot omit any of them
	// TODO only calculate the checksum in sync that was requested to match, the rest could be async ... but the tests currently expect all to be present
	// TODO the hashes all implement BinaryMarshaler so we could try to persist the state for resumable upload. we would neet do keep track of the copied bytes ...
	sha1h := sha1.New()
	md5h := md5.New()
	adler32h := adler32.New()
	{
		f, err := os.Open(upload.binPath)
		if err != nil {
			sublog.Err(err).Msg("Decomposedfs: could not open file for checksumming")
			// we can continue if no oc checksum header is set
		}
		defer f.Close()

		r1 := io.TeeReader(f, sha1h)
		r2 := io.TeeReader(r1, md5h)

		if _, err := io.Copy(adler32h, r2); err != nil {
			sublog.Err(err).Msg("Decomposedfs: could not copy bytes for checksumming")
		}
	}
	// compare if they match the sent checksum
	// TODO the tus checksum extension would do this on every chunk, but I currently don't see an easy way to pass in the requested checksum. for now we do it in FinishUpload which is also called for chunked uploads
	if upload.info.MetaData["checksum"] != "" {
		parts := strings.SplitN(upload.info.MetaData["checksum"], " ", 2)
		if len(parts) != 2 {
			return errtypes.BadRequest("invalid checksum format. must be '[algorithm] [checksum]'")
		}
		switch parts[0] {
		case "sha1":
			err = upload.checkHash(parts[1], sha1h)
		case "md5":
			err = upload.checkHash(parts[1], md5h)
		case "adler32":
			err = upload.checkHash(parts[1], adler32h)
		default:
			err = errtypes.BadRequest("unsupported checksum algorithm: " + parts[0])
		}
		if err != nil {
			return err
		}
	}
	n.BlobID = upload.info.ID // This can be changed to a content hash in the future when reference counting for the blobs was added

	// defer writing the checksums until the node is in place

	// if target exists create new version
	versionsPath := ""
	if fi, err = os.Stat(targetPath); err == nil {
		// When the if-match header was set we need to check if the
		// etag still matches before finishing the upload.
		if ifMatch, ok := upload.info.MetaData["if-match"]; ok {
			var targetEtag string
			targetEtag, err = node.CalculateEtag(n.ID, fi.ModTime())
			if err != nil {
				return errtypes.InternalError(err.Error())
			}
			if ifMatch != targetEtag {
				return errtypes.Aborted("etag mismatch")
			}
		}

		// FIXME move versioning to blobs ... no need to copy all the metadata! well ... it does if we want to version metadata...
		// versions are stored alongside the actual file, so a rename can be efficient and does not cross storage / partition boundaries
		versionsPath = upload.fs.lu.InternalPath(spaceID, n.ID+node.RevisionIDDelimiter+fi.ModTime().UTC().Format(time.RFC3339Nano))

		// This move drops all metadata!!! We copy it below with CopyMetadata
		// FIXME the node must remain the same. otherwise we might restore share metadata
		if err = os.Rename(targetPath, versionsPath); err != nil {
			sublog.Err(err).
				Str("binPath", upload.binPath).
				Str("versionsPath", versionsPath).
				Msg("Decomposedfs: could not create version")
			return
		}

	}

	// upload the data to the blobstore
	file, err := os.Open(upload.binPath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = upload.fs.tp.WriteBlob(n, file)
	if err != nil {
		return errors.Wrap(err, "failed to upload file to blostore")
	}

	// now truncate the upload (the payload stays in the blobstore) and move it to the target path
	// TODO put uploads on the same underlying storage as the destination dir?
	// TODO trigger a workflow as the final rename might eg involve antivirus scanning
	if err = os.Truncate(upload.binPath, 0); err != nil {
		sublog.Err(err).
			Msg("Decomposedfs: could not truncate")
		return
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0700); err != nil {
		sublog.Warn().Err(err).Msg("Decomposedfs: could not create node dir, trying to write file anyway")
	}
	if err = os.Rename(upload.binPath, targetPath); err != nil {
		sublog.Error().Err(err).Msg("Decomposedfs: could not rename")
		return
	}
	if versionsPath != "" {
		// copy grant and arbitrary metadata
		// FIXME ... now restoring an older revision might bring back a grant that was removed!
		err = xattrs.CopyMetadata(versionsPath, targetPath, func(attributeName string) bool {
			return true
			// TODO determine all attributes that must be copied, currently we just copy all and overwrite changed properties
			/*
				return strings.HasPrefix(attributeName, xattrs.GrantPrefix) || // for grants
					strings.HasPrefix(attributeName, xattrs.MetadataPrefix) || // for arbitrary metadata
					strings.HasPrefix(attributeName, xattrs.FavPrefix) || // for favorites
					strings.HasPrefix(attributeName, xattrs.SpaceNameAttr) || // for a shared file
			*/
		})
		if err != nil {
			sublog.Info().Err(err).Msg("Decomposedfs: failed to copy xattrs")
		}
	}

	// now try write all checksums
	tryWritingChecksum(&sublog, n, "sha1", sha1h)
	tryWritingChecksum(&sublog, n, "md5", md5h)
	tryWritingChecksum(&sublog, n, "adler32", adler32h)

	// who will become the owner?  the owner of the parent actually ... not the currently logged in user
	err = n.WriteAllNodeMetadata()
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: could not write metadata")
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentInternalPath(), n.Name)
	var link string
	link, err = os.Readlink(childNameLink)
	if err == nil && link != "../"+n.ID {
		sublog.Err(err).
			Interface("node", n).
			Str("childNameLink", childNameLink).
			Str("link", link).
			Msg("Decomposedfs: child name link has wrong target id, repairing")

		if err = os.Remove(childNameLink); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not remove symlink child entry")
		}
	}
	if errors.Is(err, iofs.ErrNotExist) || link != "../"+n.ID {
		relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
		if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
			return errors.Wrap(err, "Decomposedfs: could not symlink child entry")
		}
	}

	// only delete the upload if it was successfully written to the storage
	if err = os.Remove(upload.infoPath); err != nil {
		if !errors.Is(err, iofs.ErrNotExist) {
			sublog.Err(err).Msg("Decomposedfs: could not delete upload info")
			return
		}
	}
	// use set arbitrary metadata?
	if upload.info.MetaData["mtime"] != "" {
		err := n.SetMtime(ctx, upload.info.MetaData["mtime"])
		if err != nil {
			sublog.Err(err).Interface("info", upload.info).Msg("Decomposedfs: could not set mtime metadata")
			return err
		}
	}

	n.Exists = true

	return upload.fs.tp.Propagate(upload.ctx, n)
}

func (upload *fileUpload) checkHash(expected string, h hash.Hash) error {
	if expected != hex.EncodeToString(h.Sum(nil)) {
		upload.discardChunk()
		return errtypes.ChecksumMismatch(fmt.Sprintf("invalid checksum: expected %s got %x", upload.info.MetaData["checksum"], h.Sum(nil)))
	}
	return nil
}
func tryWritingChecksum(log *zerolog.Logger, n *node.Node, algo string, h hash.Hash) {
	if err := n.SetChecksum(algo, h); err != nil {
		log.Err(err).
			Str("csType", algo).
			Bytes("hash", h.Sum(nil)).
			Msg("Decomposedfs: could not write checksum")
		// this is not critical, the bytes are there so we will continue
	}
}

func (upload *fileUpload) discardChunk() {
	if err := os.Remove(upload.binPath); err != nil {
		if !errors.Is(err, iofs.ErrNotExist) {
			appctx.GetLogger(upload.ctx).Err(err).Interface("info", upload.info).Str("binPath", upload.binPath).Interface("info", upload.info).Msg("Decomposedfs: could not discard chunk")
			return
		}
	}
	if err := os.Remove(upload.infoPath); err != nil {
		if !errors.Is(err, iofs.ErrNotExist) {
			appctx.GetLogger(upload.ctx).Err(err).Interface("info", upload.info).Str("infoPath", upload.infoPath).Interface("info", upload.info).Msg("Decomposedfs: could not discard chunk info")
			return
		}
	}
}

// To implement the termination extension as specified in https://tus.io/protocols/resumable-upload.html#termination
// - the storage needs to implement AsTerminatableUpload
// - the upload needs to implement Terminate

// AsTerminatableUpload returns a TerminatableUpload
func (fs *Decomposedfs) AsTerminatableUpload(upload tusd.Upload) tusd.TerminatableUpload {
	return upload.(*fileUpload)
}

// Terminate terminates the upload
func (upload *fileUpload) Terminate(ctx context.Context) error {
	if err := os.Remove(upload.infoPath); err != nil {
		if !errors.Is(err, iofs.ErrNotExist) {
			return err
		}
	}
	if err := os.Remove(upload.binPath); err != nil {
		if !errors.Is(err, iofs.ErrNotExist) {
			return err
		}
	}
	return nil
}

// To implement the creation-defer-length extension as specified in https://tus.io/protocols/resumable-upload.html#creation
// - the storage needs to implement AsLengthDeclarableUpload
// - the upload needs to implement DeclareLength

// AsLengthDeclarableUpload returns a LengthDeclarableUpload
func (fs *Decomposedfs) AsLengthDeclarableUpload(upload tusd.Upload) tusd.LengthDeclarableUpload {
	return upload.(*fileUpload)
}

// DeclareLength updates the upload length information
func (upload *fileUpload) DeclareLength(ctx context.Context, length int64) error {
	upload.info.Size = length
	upload.info.SizeIsDeferred = false
	return upload.writeInfo()
}

// To implement the concatenation extension as specified in https://tus.io/protocols/resumable-upload.html#concatenation
// - the storage needs to implement AsConcatableUpload
// - the upload needs to implement ConcatUploads

// AsConcatableUpload returns a ConcatableUpload
func (fs *Decomposedfs) AsConcatableUpload(upload tusd.Upload) tusd.ConcatableUpload {
	return upload.(*fileUpload)
}

// ConcatUploads concatenates multiple uploads
func (upload *fileUpload) ConcatUploads(ctx context.Context, uploads []tusd.Upload) (err error) {
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, partialUpload := range uploads {
		fileUpload := partialUpload.(*fileUpload)

		src, err := os.Open(fileUpload.binPath)
		if err != nil {
			return err
		}
		defer src.Close()

		if _, err := io.Copy(file, src); err != nil {
			return err
		}
	}

	return
}

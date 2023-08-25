// Copyright 2018-2022 CERN
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

package upload

import (
	"context"
	"encoding/json"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/logger"
	"github.com/cs3org/reva/v2/pkg/storage/utils/chunking"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/lockedfile"
	tusd "github.com/tus/tusd/pkg/handler"
)

var defaultFilePerm = os.FileMode(0664)

// PermissionsChecker defines an interface for checking permissions on a Node
type PermissionsChecker interface {
	AssemblePermissions(ctx context.Context, n *node.Node) (ap provider.ResourcePermissions, err error)
}

// New returns a new processing instance
func New(ctx context.Context, info tusd.FileInfo, lu *lookup.Lookup, tp Tree, p PermissionsChecker, fsRoot string, pub events.Publisher, async bool, tknopts options.TokenOptions) (upload *Upload, err error) {

	log := appctx.GetLogger(ctx)
	log.Debug().Interface("info", info).Msg("Decomposedfs: NewUpload")

	if info.MetaData["filename"] == "" {
		return nil, errors.New("Decomposedfs: missing filename in metadata")
	}
	if info.MetaData["dir"] == "" {
		return nil, errors.New("Decomposedfs: missing dir in metadata")
	}

	n, err := lu.NodeFromSpaceID(ctx, info.Storage["SpaceRoot"])
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error getting space root node")
	}

	n, err = lookupNode(ctx, n, filepath.Join(info.MetaData["dir"], info.MetaData["filename"]), lu)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	log.Debug().Interface("info", info).Interface("node", n).Msg("Decomposedfs: resolved filename")

	// the parent owner will become the new owner
	parent, perr := n.Parent(ctx)
	if perr != nil {
		return nil, errors.Wrap(perr, "Decomposedfs: error getting parent "+n.ParentID)
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
		checkNode = parent
		path, _ = storagespace.FormatReference(&provider.Reference{ResourceId: &provider.ResourceId{
			SpaceId:  checkNode.SpaceID,
			OpaqueId: checkNode.ID,
		}, Path: n.Name})
	}
	rp, err := p.AssemblePermissions(ctx, checkNode)
	switch {
	case err != nil:
		return nil, err
	case !rp.InitiateFileUpload:
		return nil, errtypes.PermissionDenied(path)
	}

	// are we trying to overwriting a folder with a file?
	if n.Exists && n.IsDir(ctx) {
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

	binPath := filepath.Join(fsRoot, "uploads", info.ID)
	usr := ctxpkg.ContextMustGetUser(ctx)

	var (
		spaceRoot string
		ok        bool
	)
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

		"NodeId":              n.ID,
		"NodeExists":          "true",
		"NodeParentId":        n.ParentID,
		"NodeName":            n.Name,
		"SpaceRoot":           spaceRoot,
		"SpaceOwnerOrManager": info.Storage["SpaceOwnerOrManager"],

		"Idp":      usr.Id.Idp,
		"UserId":   usr.Id.OpaqueId,
		"UserType": utils.UserTypeToString(usr.Id.Type),
		"UserName": usr.Username,

		"LogLevel": log.GetLevel().String(),
	}
	if !n.Exists {
		// fill future node info
		info.Storage["NodeId"] = uuid.New().String()
		info.Storage["NodeExists"] = "false"
	}
	if info.MetaData["if-none-match"] == "*" && info.Storage["NodeExists"] == "true" {
		return nil, errtypes.Aborted(fmt.Sprintf("parent %s already has a child %s", n.ID, n.Name))
	}
	// Create binary file in the upload folder with no content
	log.Debug().Interface("info", info).Msg("Decomposedfs: built storage info")
	file, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	u := buildUpload(ctx, info, binPath, filepath.Join(fsRoot, "uploads", info.ID+".info"), lu, tp, pub, async, tknopts)

	// writeInfo creates the file by itself if necessary
	err = u.writeInfo()
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Get returns the Upload for the given upload id
func Get(ctx context.Context, id string, lu *lookup.Lookup, tp Tree, fsRoot string, pub events.Publisher, async bool, tknopts options.TokenOptions) (*Upload, error) {
	infoPath := filepath.Join(fsRoot, "uploads", id+".info")

	info := tusd.FileInfo{}
	data, err := os.ReadFile(infoPath)
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

	up := buildUpload(ctx, info, info.Storage["BinPath"], infoPath, lu, tp, pub, async, tknopts)
	up.versionsPath = info.MetaData["versionsPath"]
	up.SizeDiff, _ = strconv.ParseInt(info.MetaData["sizeDiff"], 10, 64)
	return up, nil
}

// CreateNodeForUpload will create the target node for the Upload
func CreateNodeForUpload(upload *Upload, initAttrs node.Attributes) (*node.Node, error) {
	ctx, span := tracer.Start(upload.Ctx, "CreateNodeForUpload")
	defer span.End()
	_, subspan := tracer.Start(ctx, "os.Stat")
	fi, err := os.Stat(upload.binPath)
	subspan.End()
	if err != nil {
		return nil, err
	}

	fsize := fi.Size()
	spaceID := upload.Info.Storage["SpaceRoot"]
	n := node.New(
		spaceID,
		upload.Info.Storage["NodeId"],
		upload.Info.Storage["NodeParentId"],
		upload.Info.Storage["NodeName"],
		fsize,
		upload.Info.ID,
		provider.ResourceType_RESOURCE_TYPE_FILE,
		nil,
		upload.lu,
	)
	n.SpaceRoot, err = node.ReadNode(ctx, upload.lu, spaceID, spaceID, false, nil, false)
	if err != nil {
		return nil, err
	}

	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return nil, err
	}

	var f *lockedfile.File
	switch upload.Info.Storage["NodeExists"] {
	case "false":
		f, err = initNewNode(upload, n, uint64(fsize))
		if f != nil {
			appctx.GetLogger(upload.Ctx).Info().Str("lockfile", f.Name()).Interface("err", err).Msg("got lock file from initNewNode")
		}
	default:
		f, err = updateExistingNode(upload, n, spaceID, uint64(fsize))
		if f != nil {
			appctx.GetLogger(upload.Ctx).Info().Str("lockfile", f.Name()).Interface("err", err).Msg("got lock file from updateExistingNode")
		}
	}
	defer func() {
		if f == nil {
			return
		}
		if err := f.Close(); err != nil {
			appctx.GetLogger(upload.Ctx).Error().Err(err).Str("nodeid", n.ID).Str("parentid", n.ParentID).Msg("could not close lock")
		}
	}()
	if err != nil {
		return nil, err
	}

	mtime := time.Now()
	if upload.Info.MetaData["mtime"] != "" {
		// overwrite mtime if requested
		mtime, err = utils.MTimeToTime(upload.Info.MetaData["mtime"])
		if err != nil {
			return nil, err
		}
	}

	// overwrite technical information
	initAttrs.SetString(prefixes.MTimeAttr, mtime.UTC().Format(time.RFC3339Nano))
	initAttrs.SetInt64(prefixes.TypeAttr, int64(provider.ResourceType_RESOURCE_TYPE_FILE))
	initAttrs.SetString(prefixes.ParentidAttr, n.ParentID)
	initAttrs.SetString(prefixes.NameAttr, n.Name)
	initAttrs.SetString(prefixes.BlobIDAttr, n.BlobID)
	initAttrs.SetInt64(prefixes.BlobsizeAttr, n.Blobsize)
	initAttrs.SetString(prefixes.StatusPrefix, node.ProcessingStatus+upload.Info.ID)

	// update node metadata with new blobid etc
	err = n.SetXattrsWithContext(ctx, initAttrs, false)
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: could not write metadata")
	}

	// add etag to metadata
	upload.Info.MetaData["etag"], _ = node.CalculateEtag(n, mtime)

	// update nodeid for later
	upload.Info.Storage["NodeId"] = n.ID
	if err := upload.writeInfo(); err != nil {
		return nil, err
	}

	return n, nil
}

func initNewNode(upload *Upload, n *node.Node, fsize uint64) (*lockedfile.File, error) {
	// create folder structure (if needed)
	if err := os.MkdirAll(filepath.Dir(n.InternalPath()), 0700); err != nil {
		return nil, err
	}

	// create and write lock new node metadata
	f, err := lockedfile.OpenFile(upload.lu.MetadataBackend().LockfilePath(n.InternalPath()), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	// we also need to touch the actual node file here it stores the mtime of the resource
	h, err := os.OpenFile(n.InternalPath(), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return f, err
	}
	h.Close()

	if _, err := node.CheckQuota(upload.Ctx, n.SpaceRoot, false, 0, fsize); err != nil {
		return f, err
	}

	// link child name to parent if it is new
	childNameLink := filepath.Join(n.ParentPath(), n.Name)
	relativeNodePath := filepath.Join("../../../../../", lookup.Pathify(n.ID, 4, 2))
	log := appctx.GetLogger(upload.Ctx).With().Str("childNameLink", childNameLink).Str("relativeNodePath", relativeNodePath).Logger()
	log.Info().Msg("initNewNode: creating symlink")

	if err = os.Symlink(relativeNodePath, childNameLink); err != nil {
		log.Info().Err(err).Msg("initNewNode: symlink failed")
		if errors.Is(err, iofs.ErrExist) {
			log.Info().Err(err).Msg("initNewNode: symlink already exists")
			return f, errtypes.AlreadyExists(n.Name)
		}
		return f, errors.Wrap(err, "Decomposedfs: could not symlink child entry")
	}
	log.Info().Msg("initNewNode: symlink created")

	// on a new file the sizeDiff is the fileSize
	upload.SizeDiff = int64(fsize)
	upload.Info.MetaData["sizeDiff"] = strconv.Itoa(int(upload.SizeDiff))
	return f, nil
}

func updateExistingNode(upload *Upload, n *node.Node, spaceID string, fsize uint64) (*lockedfile.File, error) {
	targetPath := n.InternalPath()

	// write lock existing node before reading any metadata
	f, err := lockedfile.OpenFile(upload.lu.MetadataBackend().LockfilePath(targetPath), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	old, _ := node.ReadNode(upload.Ctx, upload.lu, spaceID, n.ID, false, nil, false)
	if _, err := node.CheckQuota(upload.Ctx, n.SpaceRoot, true, uint64(old.Blobsize), fsize); err != nil {
		return f, err
	}

	oldNodeMtime, err := old.GetMTime(upload.Ctx)
	if err != nil {
		return f, err
	}
	oldNodeEtag, err := node.CalculateEtag(old, oldNodeMtime)
	if err != nil {
		return f, err
	}

	// When the if-match header was set we need to check if the
	// etag still matches before finishing the upload.
	if ifMatch, ok := upload.Info.MetaData["if-match"]; ok {
		if ifMatch != oldNodeEtag {
			return f, errtypes.Aborted("etag mismatch")
		}
	}

	// When the if-none-match header was set we need to check if any of the
	// etags matches before finishing the upload.
	if ifNoneMatch, ok := upload.Info.MetaData["if-none-match"]; ok {
		if ifNoneMatch == "*" {
			return f, errtypes.Aborted("etag mismatch, resource exists")
		}
		for _, ifNoneMatchTag := range strings.Split(ifNoneMatch, ",") {
			if ifNoneMatchTag == oldNodeEtag {
				return f, errtypes.Aborted("etag mismatch")
			}
		}
	}

	// When the if-unmodified-since header was set we need to check if the
	// etag still matches before finishing the upload.
	if ifUnmodifiedSince, ok := upload.Info.MetaData["if-unmodified-since"]; ok {
		if err != nil {
			return f, errtypes.InternalError(fmt.Sprintf("failed to read mtime of node: %s", err))
		}
		ifUnmodifiedSince, err := time.Parse(time.RFC3339Nano, ifUnmodifiedSince)
		if err != nil {
			return f, errtypes.InternalError(fmt.Sprintf("failed to parse if-unmodified-since time: %s", err))
		}

		if oldNodeMtime.After(ifUnmodifiedSince) {
			return f, errtypes.Aborted("if-unmodified-since mismatch")
		}
	}

	upload.versionsPath = upload.lu.InternalPath(spaceID, n.ID+node.RevisionIDDelimiter+oldNodeMtime.UTC().Format(time.RFC3339Nano))
	upload.SizeDiff = int64(fsize) - old.Blobsize
	upload.Info.MetaData["versionsPath"] = upload.versionsPath
	upload.Info.MetaData["sizeDiff"] = strconv.Itoa(int(upload.SizeDiff))

	// create version node
	if _, err := os.Create(upload.versionsPath); err != nil {
		return f, err
	}

	// copy blob metadata to version node
	if err := upload.lu.CopyMetadataWithSourceLock(upload.Ctx, targetPath, upload.versionsPath, func(attributeName string, value []byte) (newValue []byte, copy bool) {
		return value, strings.HasPrefix(attributeName, prefixes.ChecksumPrefix) ||
			attributeName == prefixes.TypeAttr ||
			attributeName == prefixes.BlobIDAttr ||
			attributeName == prefixes.BlobsizeAttr ||
			attributeName == prefixes.MTimeAttr
	}, f, true); err != nil {
		return f, err
	}

	// keep mtime from previous version
	if err := os.Chtimes(upload.versionsPath, oldNodeMtime, oldNodeMtime); err != nil {
		return f, errtypes.InternalError(fmt.Sprintf("failed to change mtime of version node: %s", err))
	}

	return f, nil
}

// lookupNode looks up nodes by path.
// This method can also handle lookups for paths which contain chunking information.
func lookupNode(ctx context.Context, spaceRoot *node.Node, path string, lu *lookup.Lookup) (*node.Node, error) {
	p := path
	isChunked := chunking.IsChunked(path)
	if isChunked {
		chunkInfo, err := chunking.GetChunkBLOBInfo(path)
		if err != nil {
			return nil, err
		}
		p = chunkInfo.Path
	}

	n, err := lu.WalkPath(ctx, spaceRoot, p, true, func(ctx context.Context, n *node.Node) error { return nil })
	if err != nil {
		return nil, errors.Wrap(err, "Decomposedfs: error walking path")
	}

	if isChunked {
		n.Name = filepath.Base(path)
	}
	return n, nil
}

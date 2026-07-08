// Copyright 2018-2024 CERN
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
	"crypto/md5"  //nolint:gosec
	"crypto/sha1" //nolint:gosec
	"hash/adler32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/renameio/v2"
	"github.com/pkg/errors"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"github.com/owncloud/reva/v2/pkg/appctx"
	ctxpkg "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/errtypes"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/utils"
)

const defaultFilePerm = os.FileMode(0664)

// FileSession is the Session implementation for disk-backed uploads. While an
// upload is in progress, incoming bytes are staged in a .bin file and upload
// metadata (size, owner, checksums, etc.) is persisted in a .info file. Both
// survive process restarts, allowing TUS resumption.
//
// In scope: read/write the staged .bin file, persist/load upload metadata in the .info file.
// Out of scope: TUS protocol, checksums, event publishing, postprocessing — those live in coordinatedUpload and coordinator.
type FileSession struct {
	store *FileStore
	info  tusd.FileInfo
}

// Session is the driver-agnostic view of an upload session the Coordinator
// needs. Implementations must be pure state (CRUD): protocol orchestration
// belongs to coordinatedUpload or the coordinator itself.
type Session interface {
	storage.UploadSession

	// Data access — delegated to by coordinatedUpload for TUS reads/writes.
	GetInfo(ctx context.Context) (tusd.FileInfo, error)
	GetReader(ctx context.Context) (io.ReadCloser, error)
	WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error)

	// Internal coordinator plumbing.
	BinPath() string
	ProviderID() string
	SpaceID() string
	NodeID() string
	NodeExists() bool
	Dir() string
	URL(ctx context.Context) (string, error)
	SetScanData(result string, date time.Time)
	Checksums() storage.UploadChecksums
	SetChecksums(sha1, md5, adler32 []byte)
	Metadata() map[string]string
	Persist(ctx context.Context) error
	Cleanup(cleanBin, cleanInfo bool)
	Context(ctx context.Context) context.Context

	// Typed setters used by Coordinator.InitiateUpload to populate a new session
	// without knowing internal storage key names.
	SetStorageValue(key, value string)
	SetMetadata(key, value string)
	SetSize(size int64)
	SetSizeIsDeferred(value bool)
	SetExecutant(u *userpb.User)
	TouchBin() error
}

func (s *FileSession) GetInfo(_ context.Context) (tusd.FileInfo, error) {
	return s.info, nil
}

func (s *FileSession) GetReader(_ context.Context) (io.ReadCloser, error) {
	return os.Open(s.binPath())
}

func (s *FileSession) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := os.OpenFile(s.binPath(), os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	n, err := io.Copy(file, src)
	if err != nil && err != io.ErrUnexpectedEOF {
		return n, err
	}
	s.info.Offset += n
	return n, nil
}


// Purge removes all on-disk state for this session.
func (s *FileSession) Purge(ctx context.Context) {
	s.Cleanup(true, true)
}

// ScanData returns the AV scan result and scan date stored on the session.
func (s *FileSession) ScanData() (string, time.Time) {
	date := s.info.MetaData["scanDate"]
	if date == "" {
		return "", time.Time{}
	}
	d, _ := time.Parse(time.RFC3339, date)
	return s.info.MetaData["scanResult"], d
}

// ID returns the upload session ID.
func (s *FileSession) ID() string {
	return s.info.ID
}

// Filename returns the filename stored in the session.
func (s *FileSession) Filename() string {
	return s.info.Storage["NodeName"]
}

// Size returns the declared upload size.
func (s *FileSession) Size() int64 {
	return s.info.Size
}

// Offset returns the current upload offset.
func (s *FileSession) Offset() int64 {
	return s.info.Offset
}

// BinPath returns the path to the staged binary file.
func (s *FileSession) BinPath() string {
	return s.binPath()
}

// SpaceGid returns the numeric GID of the space owner, or "" if not set.
func (s *FileSession) SpaceGid() string {
	return s.info.Storage["SpaceGid"]
}

// ProviderID returns the storage provider ID stored in the session.
func (s *FileSession) ProviderID() string {
	return s.info.MetaData["providerID"]
}

// SpaceID returns the space (root) ID.
func (s *FileSession) SpaceID() string {
	return s.info.Storage["SpaceRoot"]
}

// NodeID returns the node ID for this upload.
func (s *FileSession) NodeID() string {
	return s.info.Storage["NodeId"]
}

// NodeExists returns whether the target node existed when the upload was initiated.
func (s *FileSession) NodeExists() bool {
	return s.info.Storage["NodeExists"] == "true"
}

// Dir returns the directory portion of the upload path.
func (s *FileSession) Dir() string {
	return s.info.Storage["Dir"]
}

// IsProcessing returns true if all bytes are received but postprocessing has not finished.
func (s *FileSession) IsProcessing() bool {
	return s.info.Size == s.info.Offset && s.info.MetaData["scanResult"] == ""
}

// SpaceOwner returns the space owner user ID.
func (s *FileSession) SpaceOwner() *userpb.UserId {
	return &userpb.UserId{
		OpaqueId: s.info.Storage["SpaceOwnerOrManager"],
		Idp:      s.info.Storage["SpaceOwnerIdp"],
		Type:     userpb.UserType(userpb.UserType_value[s.info.Storage["SpaceOwnerType"]]),
	}
}

// Executant returns the user ID of the user who initiated this upload.
func (s *FileSession) Executant() userpb.UserId {
	return userpb.UserId{
		Type:     utils.UserTypeMap(s.info.Storage["UserType"]),
		Idp:      s.info.Storage["Idp"],
		OpaqueId: s.info.Storage["UserId"],
	}
}

// Expires returns the upload expiry time.
func (s *FileSession) Expires() time.Time {
	var t time.Time
	if value, ok := s.info.MetaData["expires"]; ok {
		t, _ = utils.MTimeToTime(value)
	}
	return t
}

// Reference returns a CS3 reference for the resource being uploaded.
func (s *FileSession) Reference() provider.Reference {
	return provider.Reference{
		ResourceId: &provider.ResourceId{
			StorageId: s.info.MetaData["providerID"],
			SpaceId:   s.info.Storage["SpaceRoot"],
			OpaqueId:  s.info.Storage["NodeId"],
		},
	}
}

// Checksums returns the pre-computed checksums stored on the session.
func (s *FileSession) Checksums() storage.UploadChecksums {
	decode := func(key string) []byte {
		b, _ := hex.DecodeString(s.info.MetaData[key])
		return b
	}
	return storage.UploadChecksums{
		SHA1:    decode("checksumSHA1"),
		MD5:     decode("checksumMD5"),
		Adler32: decode("checksumAdler32"),
	}
}

// Metadata returns the upload metadata map passed to CommitUpload.
func (s *FileSession) Metadata() map[string]string {
	return map[string]string{
		"providerID":   s.info.MetaData["providerID"],
		"mtime":        s.info.MetaData["mtime"],
		"nodeExists":   s.info.Storage["NodeExists"],
		"versionsPath": s.info.MetaData["versionsPath"],
	}
}

// SetScanData stores AV scan results on the session.
func (s *FileSession) SetScanData(result string, date time.Time) {
	s.info.MetaData["scanResult"] = result
	s.info.MetaData["scanDate"] = date.Format(time.RFC3339)
}

// SetChecksums stores pre-computed checksums so CommitUpload can use them without
// re-reading the binary file.
func (s *FileSession) SetChecksums(sha1Sum, md5Sum, adler32Sum []byte) {
	s.info.MetaData["checksumSHA1"] = hex.EncodeToString(sha1Sum)
	s.info.MetaData["checksumMD5"] = hex.EncodeToString(md5Sum)
	s.info.MetaData["checksumAdler32"] = hex.EncodeToString(adler32Sum)
}

// SetMetadata sets a user-visible upload metadata field.
func (s *FileSession) SetMetadata(key, value string) {
	s.info.MetaData[key] = value
}

// SetStorageValue sets an internal storage field on the session.
func (s *FileSession) SetStorageValue(key, value string) {
	s.info.Storage[key] = value
}

// SetSize updates the declared upload size.
func (s *FileSession) SetSize(size int64) {
	s.info.Size = size
}

// SetSizeIsDeferred marks the upload size as not yet known.
func (s *FileSession) SetSizeIsDeferred(value bool) {
	s.info.SizeIsDeferred = value
}

// SetExecutant stores the identity of the user who initiated the upload.
func (s *FileSession) SetExecutant(u *userpb.User) {
	s.info.Storage["Idp"] = u.GetId().GetIdp()
	s.info.Storage["UserId"] = u.GetId().GetOpaqueId()
	s.info.Storage["UserType"] = utils.UserTypeToString(u.GetId().Type)
	s.info.Storage["UserName"] = u.GetUsername()
	s.info.Storage["UserDisplayName"] = u.GetDisplayName()
	b, _ := json.Marshal(u.GetOpaque())
	s.info.Storage["UserOpaque"] = string(b)
}

// TouchBin creates the empty staging file.
func (s *FileSession) TouchBin() error {
	f, err := os.OpenFile(s.binPath(), os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		return err
	}
	return f.Close()
}

// Persist writes the session metadata atomically to disk.
func (s *FileSession) Persist(ctx context.Context) error {
	infoPath := s.infoPath()
	if err := os.MkdirAll(filepath.Dir(infoPath), 0700); err != nil {
		return err
	}
	d, err := json.Marshal(s.info)
	if err != nil {
		return err
	}
	return renameio.WriteFile(infoPath, d, 0600)
}

// Cleanup removes the staged binary and/or info file.
// Node deletion and processing flag changes are the coordinator's responsibility.
func (s *FileSession) Cleanup(cleanBin, cleanInfo bool) {
	log := s.store.log
	if cleanBin {
		if err := os.Remove(s.binPath()); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Error().Str("path", s.binPath()).Err(err).Msg("filestore: removing staged binary failed")
		}
	}
	if cleanInfo {
		if err := os.Remove(s.infoPath()); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Error().Str("path", s.infoPath()).Err(err).Msg("filestore: removing session info failed")
		}
	}
}

// Context reconstructs a context carrying the user, logger, lock ID, and
// initiator ID that were recorded when the upload was initiated.
func (s *FileSession) Context(ctx context.Context) context.Context {
	sub := s.store.log.With().Int("pid", os.Getpid()).Logger()
	ctx = appctx.WithLogger(ctx, &sub)
	ctx = ctxpkg.ContextSetLockID(ctx, s.info.MetaData["lockid"])
	ctx = ctxpkg.ContextSetUser(ctx, s.executantUser())
	return ctxpkg.ContextSetInitiator(ctx, s.info.MetaData["initiatorid"])
}

// URL returns a signed JWT URL that the postprocessing service can use to
// download the staged binary via the data gateway.
func (s *FileSession) URL(_ context.Context) (string, error) {
	type transferClaims struct {
		jwt.RegisteredClaims
		Target string `json:"target"`
	}

	target := joinURLParts(s.store.opts.DownloadEndpoint, "tus/", s.info.ID)
	ttl := time.Duration(s.store.opts.TransferExpires) * time.Second
	claims := transferClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Audience:  jwt.ClaimStrings{"reva"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Target: target,
	}
	t := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	tkn, err := t.SignedString([]byte(s.store.opts.TransferSharedSecret))
	if err != nil {
		return "", errors.Wrapf(err, "filestore: error signing transfer token with claims %+v", claims)
	}
	return joinURLParts(s.store.opts.DataGatewayEndpoint, tkn), nil
}

func (s *FileSession) ToFileInfo() tusd.FileInfo {
	return s.info
}

func (s *FileSession) InitiatorID() string {
	return s.info.MetaData["initiatorid"]
}

func (s *FileSession) binPath() string {
	return filepath.Join(s.store.root, "uploads", s.info.ID)
}

func (s *FileSession) infoPath() string {
	return fileSessionPath(s.store.root, s.info.ID)
}

func (s *FileSession) executantUser() *userpb.User {
	var o *typespb.Opaque
	_ = json.Unmarshal([]byte(s.info.Storage["UserOpaque"]), &o)
	return &userpb.User{
		Id: &userpb.UserId{
			Type:     utils.UserTypeMap(s.info.Storage["UserType"]),
			Idp:      s.info.Storage["Idp"],
			OpaqueId: s.info.Storage["UserId"],
		},
		Username:    s.info.Storage["UserName"],
		DisplayName: s.info.Storage["UserDisplayName"],
		Opaque:      o,
	}
}

// fileSessionPath returns the path to the .info file for the given session ID.
func fileSessionPath(root, id string) string {
	return filepath.Join(root, "uploads", id+".info")
}

// calculateChecksums computes sha1, md5, and adler32 in a single pass over path.
func calculateChecksums(_ context.Context, path string) (hash.Hash, hash.Hash, hash.Hash32, error) {
	sha1h := sha1.New()   //nolint:gosec
	md5h := md5.New()    //nolint:gosec
	adler32h := adler32.New()

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r1 := io.TeeReader(f, sha1h)
	r2 := io.TeeReader(r1, md5h)
	if _, err = io.Copy(adler32h, r2); err != nil {
		return nil, nil, nil, err
	}
	return sha1h, md5h, adler32h, nil
}

func checkHash(expected string, h hash.Hash) error {
	got := hex.EncodeToString(h.Sum(nil))
	if expected != got {
		return errtypes.ChecksumMismatch(fmt.Sprintf("invalid checksum: expected %s got %s", expected, got))
	}
	return nil
}

// joinURLParts concatenates URL path segments, inserting "/" between them if needed.
func joinURLParts(parts ...string) string {
	var b strings.Builder
	for i, p := range parts {
		b.WriteString(p)
		if i < len(parts)-1 && !strings.HasSuffix(p, "/") {
			b.WriteByte('/')
		}
	}
	return b.String()
}

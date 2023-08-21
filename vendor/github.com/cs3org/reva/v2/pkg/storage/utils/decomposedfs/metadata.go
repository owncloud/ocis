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
	"path/filepath"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata/prefixes"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/pkg/errors"
)

// SetArbitraryMetadata sets the metadata on a resource
func (fs *Decomposedfs) SetArbitraryMetadata(ctx context.Context, ref *provider.Reference, md *provider.ArbitraryMetadata) (err error) {
	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error resolving ref")
	}
	sublog := appctx.GetLogger(ctx).With().Interface("node", n).Logger()

	if !n.Exists {
		err = errtypes.NotFound(filepath.Join(n.ParentID, n.Name))
		return err
	}

	rp, err := fs.p.AssemblePermissions(ctx, n)
	switch {
	case err != nil:
		return err
	case !rp.InitiateFileUpload: // TODO add explicit SetArbitraryMetadata grant to CS3 api, tracked in https://github.com/cs3org/cs3apis/issues/91
		f, _ := storagespace.FormatReference(ref)
		if rp.Stat {
			return errtypes.PermissionDenied(f)
		}
		return errtypes.NotFound(f)
	}

	// Set space owner in context
	storagespace.ContextSendSpaceOwnerID(ctx, n.SpaceOwnerOrManager(ctx))

	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return err
	}

	errs := []error{}
	// TODO should we really continue updating when an error occurs?
	if md.Metadata != nil {
		if val, ok := md.Metadata["mtime"]; ok {
			delete(md.Metadata, "mtime")
			if err := n.SetMtimeString(ctx, val); err != nil {
				errs = append(errs, errors.Wrap(err, "could not set mtime"))
			}
		}
		// TODO(jfd) special handling for atime?
		// TODO(jfd) allow setting birth time (btime)?
		// TODO(jfd) any other metadata that is interesting? fileid?
		// TODO unset when file is updated
		// TODO unset when folder is updated or add timestamp to etag?
		if val, ok := md.Metadata["etag"]; ok {
			delete(md.Metadata, "etag")
			if err := n.SetEtag(ctx, val); err != nil {
				errs = append(errs, errors.Wrap(err, "could not set etag"))
			}
		}
		if val, ok := md.Metadata[node.FavoriteKey]; ok {
			delete(md.Metadata, node.FavoriteKey)
			if u, ok := ctxpkg.ContextGetUser(ctx); ok {
				if uid := u.GetId(); uid != nil {
					if err := n.SetFavorite(ctx, uid, val); err != nil {
						sublog.Error().Err(err).
							Interface("user", u).
							Msg("could not set favorite flag")
						errs = append(errs, errors.Wrap(err, "could not set favorite flag"))
					}
				} else {
					sublog.Error().Interface("user", u).Msg("user has no id")
					errs = append(errs, errors.Wrap(errtypes.UserRequired("userrequired"), "user has no id"))
				}
			} else {
				sublog.Error().Interface("user", u).Msg("error getting user from ctx")
				errs = append(errs, errors.Wrap(errtypes.UserRequired("userrequired"), "error getting user from ctx"))
			}
		}
	}
	for k, v := range md.Metadata {
		attrName := prefixes.MetadataPrefix + k
		if err = n.SetXattrString(ctx, attrName, v); err != nil {
			errs = append(errs, errors.Wrap(err, "Decomposedfs: could not set metadata attribute "+attrName+" to "+k))
		}
	}

	switch len(errs) {
	case 0:
		return fs.tp.Propagate(ctx, n, 0)
	case 1:
		// TODO Propagate if anything changed
		return errs[0]
	default:
		// TODO Propagate if anything changed
		// TODO how to return multiple errors?
		return errors.New("multiple errors occurred, see log for details")
	}
}

// UnsetArbitraryMetadata unsets the metadata on the given resource
func (fs *Decomposedfs) UnsetArbitraryMetadata(ctx context.Context, ref *provider.Reference, keys []string) (err error) {
	n, err := fs.lu.NodeFromResource(ctx, ref)
	if err != nil {
		return errors.Wrap(err, "Decomposedfs: error resolving ref")
	}
	sublog := appctx.GetLogger(ctx).With().Interface("node", n).Logger()

	if !n.Exists {
		err = errtypes.NotFound(filepath.Join(n.ParentID, n.Name))
		return err
	}

	rp, err := fs.p.AssemblePermissions(ctx, n)
	switch {
	case err != nil:
		return err
	case !rp.InitiateFileUpload: // TODO use SetArbitraryMetadata grant to CS3 api, tracked in https://github.com/cs3org/cs3apis/issues/91
		f, _ := storagespace.FormatReference(ref)
		if rp.Stat {
			return errtypes.PermissionDenied(f)
		}
		return errtypes.NotFound(f)
	}

	// Set space owner in context
	storagespace.ContextSendSpaceOwnerID(ctx, n.SpaceOwnerOrManager(ctx))

	// check lock
	if err := n.CheckLock(ctx); err != nil {
		return err
	}

	errs := []error{}
	for _, k := range keys {
		switch k {
		case node.FavoriteKey:
			// the favorite flag is specific to the user, so we need to incorporate the userid
			u, ok := ctxpkg.ContextGetUser(ctx)
			if !ok {
				sublog.Error().
					Interface("user", u).
					Msg("error getting user from ctx")
				errs = append(errs, errors.Wrap(errtypes.UserRequired("userrequired"), "error getting user from ctx"))
				continue
			}
			var uid *userpb.UserId
			if uid = u.GetId(); uid == nil || uid.OpaqueId == "" {
				sublog.Error().
					Interface("user", u).
					Msg("user has no id")
				errs = append(errs, errors.Wrap(errtypes.UserRequired("userrequired"), "user has no id"))
				continue
			}
			fa := fmt.Sprintf("%s:%s:%s@%s", prefixes.FavPrefix, utils.UserTypeToString(uid.GetType()), uid.GetOpaqueId(), uid.GetIdp())
			if err := n.RemoveXattr(ctx, fa, true); err != nil {
				if metadata.IsAttrUnset(err) {
					continue // already gone, ignore
				}
				sublog.Error().Err(err).
					Interface("user", u).
					Str("key", fa).
					Msg("could not unset favorite flag")
				errs = append(errs, errors.Wrap(err, "could not unset favorite flag"))
			}
		default:
			if err = n.RemoveXattr(ctx, prefixes.MetadataPrefix+k, true); err != nil {
				if metadata.IsAttrUnset(err) {
					continue // already gone, ignore
				}
				sublog.Error().Err(err).
					Str("key", k).
					Msg("could not unset metadata")
				errs = append(errs, errors.Wrap(err, "could not unset metadata"))
			}
		}
	}
	switch len(errs) {
	case 0:
		return fs.tp.Propagate(ctx, n, 0)
	case 1:
		// TODO Propagate if anything changed
		return errs[0]
	default:
		// TODO Propagate if anything changed
		// TODO how to return multiple errors?
		return errors.New("multiple errors occurred, see log for details")
	}
}

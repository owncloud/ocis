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

package ocdav

import (
	"context"
	"fmt"
	"net/http"
	"path"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/errors"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/spacelookup"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func (s *svc) handlePathDelete(w http.ResponseWriter, r *http.Request, ns string) {
	fn := path.Join(ns, r.URL.Path)

	sublog := appctx.GetLogger(r.Context()).With().Str("path", fn).Logger()
	client, err := s.getClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	space, status, err := spacelookup.LookUpStorageSpaceForPath(r.Context(), client, fn)
	if err != nil {
		sublog.Error().Err(err).Msg("error sending a grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if status.Code != rpc.Code_CODE_OK {
		errors.HandleErrorStatus(&sublog, w, status)
		return
	}

	s.handleDelete(r.Context(), w, r, spacelookup.MakeRelativeReference(space, fn, false), sublog)
}

func (s *svc) handleDelete(ctx context.Context, w http.ResponseWriter, r *http.Request, ref *provider.Reference, log zerolog.Logger) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(ctx, "delete")
	defer span.End()

	req := &provider.DeleteRequest{Ref: ref}

	// FIXME the lock token is part of the application level protocol, it should be part of the DeleteRequest message not the opaque
	ih, ok := parseIfHeader(r.Header.Get(net.HeaderIf))
	if ok {
		if len(ih.lists) == 1 && len(ih.lists[0].conditions) == 1 {
			req.Opaque = utils.AppendPlainToOpaque(req.Opaque, "lockid", ih.lists[0].conditions[0].Token)
		}
	} else if r.Header.Get(net.HeaderIf) != "" {
		w.WriteHeader(http.StatusBadRequest)
		b, err := errors.Marshal(http.StatusBadRequest, "invalid if header", "")
		errors.HandleWebdavError(&log, w, b, err)
		return
	}

	client, err := s.getClient()
	if err != nil {
		log.Error().Err(err).Msg("error getting grpc client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := client.Delete(ctx, req)
	if err != nil {
		span.RecordError(err)
		log.Error().Err(err).Msg("error performing delete grpc request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch res.Status.Code {
	case rpc.Code_CODE_OK:
		w.WriteHeader(http.StatusNoContent)
	case rpc.Code_CODE_NOT_FOUND:
		w.WriteHeader(http.StatusNotFound)
		// TODO path might be empty or relative...
		m := fmt.Sprintf("Resource %v not found", ref.Path)
		b, err := errors.Marshal(http.StatusNotFound, m, "")
		errors.HandleWebdavError(&log, w, b, err)
	case rpc.Code_CODE_PERMISSION_DENIED:
		status := http.StatusForbidden
		if lockID := utils.ReadPlainFromOpaque(res.Opaque, "lockid"); lockID != "" {
			// http://www.webdav.org/specs/rfc4918.html#HEADER_Lock-Token says that the
			// Lock-Token value is a Coded-URL. We add angle brackets.
			w.Header().Set("Lock-Token", "<"+lockID+">")
			status = http.StatusLocked
		}
		// TODO path might be empty or relative...
		m := fmt.Sprintf("Permission denied to delete %v", ref.Path)
		// check if user has access to resource
		sRes, err := client.Stat(ctx, &provider.StatRequest{Ref: ref})
		if err != nil {
			span.RecordError(err)
			log.Error().Err(err).Msg("error performing stat grpc request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sRes.Status.Code != rpc.Code_CODE_OK {
			// return not found error so we dont leak existence of a file
			// TODO hide permission failed for users without access in every kind of request
			// TODO should this be done in the driver?
			status = http.StatusNotFound
			// TODO path might be empty or relative...
			m = fmt.Sprintf("%s not fount", ref.Path)
		}
		w.WriteHeader(status)
		b, err := errors.Marshal(status, m, "")
		errors.HandleWebdavError(&log, w, b, err)
	case rpc.Code_CODE_INTERNAL:
		if res.Status.Message == "can't delete mount path" {
			w.WriteHeader(http.StatusForbidden)
			b, err := errors.Marshal(http.StatusForbidden, res.Status.Message, "")
			errors.HandleWebdavError(&log, w, b, err)
		}
	default:
		status := status.HTTPStatusFromCode(res.Status.Code)
		w.WriteHeader(status)
		b, err := errors.Marshal(status, res.Status.Message, "")
		errors.HandleWebdavError(&log, w, b, err)
	}
}

func (s *svc) handleSpacesDelete(w http.ResponseWriter, r *http.Request, spaceID string) {
	ctx := r.Context()
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(ctx, "spaces_delete")
	defer span.End()

	sublog := appctx.GetLogger(ctx).With().Logger()

	ref, err := spacelookup.MakeStorageSpaceReference(spaceID, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// do not allow deleting spaces via dav endpoint - use graph endpoint instead
	// we get a relative reference coming from the space root
	// so if the path is "empty" we a referencing the space
	if ref.GetPath() == "." {
		sublog.Info().Msg("deleting spaces via dav is not allowed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.handleDelete(ctx, w, r, &ref, sublog)
}

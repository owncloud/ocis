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
	"errors"
	"fmt"
	"net/http"
	"path"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/spacelookup"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	rstatus "github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func (s *svc) handlePathMkcol(w http.ResponseWriter, r *http.Request, ns string) (status int, err error) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "mkcol")
	defer span.End()

	fn := path.Join(ns, r.URL.Path)
	for _, r := range nameRules {
		if !r.Test(fn) {
			return http.StatusBadRequest, fmt.Errorf("invalid name rule")
		}
	}
	sublog := appctx.GetLogger(ctx).With().Str("path", fn).Logger()

	// stat requested path to make sure it isn't existing yet
	// NOTE: It could be on another storage provider than the 'parent' of it
	sr, err := s.gwClient.Stat(ctx, &provider.StatRequest{
		Ref: &provider.Reference{
			Path: fn,
		},
	})
	switch {
	case err != nil:
		return http.StatusInternalServerError, err
	case sr.Status.Code == rpc.Code_CODE_OK:
		// https://www.rfc-editor.org/rfc/rfc4918#section-9.3.1:
		// 405 (Method Not Allowed) - MKCOL can only be executed on an unmapped URL.
		return http.StatusMethodNotAllowed, fmt.Errorf("The resource you tried to create already exists")
	case sr.Status.Code == rpc.Code_CODE_ABORTED:
		return http.StatusPreconditionFailed, errtypes.NewErrtypeFromStatus(sr.Status)
	case sr.Status.Code != rpc.Code_CODE_NOT_FOUND:
		return rstatus.HTTPStatusFromCode(sr.Status.Code), errtypes.NewErrtypeFromStatus(sr.Status)
	}

	parentPath := path.Dir(fn)

	space, rpcStatus, err := spacelookup.LookUpStorageSpaceForPath(ctx, s.gwClient, parentPath)
	switch {
	case err != nil:
		return http.StatusInternalServerError, err
	case rpcStatus.Code == rpc.Code_CODE_NOT_FOUND:
		// https://www.rfc-editor.org/rfc/rfc4918#section-9.3.1:
		// 409 (Conflict) - A collection cannot be made at the Request-URI until
		// one or more intermediate collections have been created.  The server
		// MUST NOT create those intermediate collections automatically.
		return http.StatusConflict, fmt.Errorf("intermediate collection does not exist")
	case rpcStatus.Code == rpc.Code_CODE_ABORTED:
		return http.StatusPreconditionFailed, errtypes.NewErrtypeFromStatus(rpcStatus)
	case rpcStatus.Code != rpc.Code_CODE_OK:
		return rstatus.HTTPStatusFromCode(rpcStatus.Code), errtypes.NewErrtypeFromStatus(rpcStatus)
	}

	return s.handleMkcol(ctx, w, r, spacelookup.MakeRelativeReference(space, parentPath, false), spacelookup.MakeRelativeReference(space, fn, false), sublog)
}

func (s *svc) handleSpacesMkCol(w http.ResponseWriter, r *http.Request, spaceID string) (status int, err error) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "spaces_mkcol")
	defer span.End()

	sublog := appctx.GetLogger(ctx).With().Str("path", r.URL.Path).Str("spaceid", spaceID).Str("handler", "mkcol").Logger()

	parentRef, err := spacelookup.MakeStorageSpaceReference(spaceID, path.Dir(r.URL.Path))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid space id")
	}
	childRef, _ := spacelookup.MakeStorageSpaceReference(spaceID, r.URL.Path)

	return s.handleMkcol(ctx, w, r, &parentRef, &childRef, sublog)
}

func (s *svc) handleMkcol(ctx context.Context, w http.ResponseWriter, r *http.Request, parentRef, childRef *provider.Reference, log zerolog.Logger) (status int, err error) {
	if r.Body != http.NoBody {
		// We currently do not support extended mkcol https://datatracker.ietf.org/doc/rfc5689/
		// TODO let clients send a body with properties to set on the new resource
		return http.StatusUnsupportedMediaType, fmt.Errorf("extended-mkcol not supported")
	}

	req := &provider.CreateContainerRequest{Ref: childRef}
	res, err := s.gwClient.CreateContainer(ctx, req)
	switch {
	case err != nil:
		return http.StatusInternalServerError, err
	case res.Status.Code == rpc.Code_CODE_OK:
		w.WriteHeader(http.StatusCreated)
		return 0, nil
	case res.Status.Code == rpc.Code_CODE_NOT_FOUND:
		// This should never happen because if the parent collection does not exist we should
		// get a Code_CODE_FAILED_PRECONDITION. We play stupid and return what the response gave us
		//lint:ignore ST1005 mimic the exact oc10 error message
		return http.StatusNotFound, errors.New("Resource not found")
	case res.Status.Code == rpc.Code_CODE_PERMISSION_DENIED:
		// check if user has access to parent
		sRes, err := s.gwClient.Stat(ctx, &provider.StatRequest{Ref: &provider.Reference{
			ResourceId: childRef.GetResourceId(),
			Path:       utils.MakeRelativePath(path.Dir(childRef.Path)),
		}})
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if sRes.Status.Code != rpc.Code_CODE_OK {
			// return not found error so we do not leak existence of a file
			// TODO hide permission failed for users without access in every kind of request
			// TODO should this be done in the driver?
			//lint:ignore ST1005 mimic the exact oc10 error message
			return http.StatusNotFound, errors.New("Resource not found")
		}
		return http.StatusForbidden, errors.New(sRes.Status.Message)
	case res.Status.Code == rpc.Code_CODE_ABORTED:
		return http.StatusPreconditionFailed, fmt.Errorf(res.Status.Message)
	case res.Status.Code == rpc.Code_CODE_FAILED_PRECONDITION:
		// https://www.rfc-editor.org/rfc/rfc4918#section-9.3.1:
		// 409 (Conflict) - A collection cannot be made at the Request-URI until
		// one or more intermediate collections have been created. The server
		// MUST NOT create those intermediate collections automatically.
		return http.StatusConflict, fmt.Errorf(res.Status.Message)
	case res.Status.Code == rpc.Code_CODE_ALREADY_EXISTS:
		// https://www.rfc-editor.org/rfc/rfc4918#section-9.3.1:
		// 405 (Method Not Allowed) - MKCOL can only be executed on an unmapped URL.
		//lint:ignore ST1005 mimic the exact oc10 error message
		return http.StatusMethodNotAllowed, errors.New("The resource you tried to create already exists")
	}
	return rstatus.HTTPStatusFromCode(res.Status.Code), errtypes.NewErrtypeFromStatus(res.Status)
}

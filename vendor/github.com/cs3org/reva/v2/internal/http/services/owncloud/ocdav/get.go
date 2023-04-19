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
	"io"
	"net/http"
	"path"
	"strconv"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/datagateway"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/errors"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/spacelookup"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
)

func (s *svc) handlePathGet(w http.ResponseWriter, r *http.Request, ns string) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "get")
	defer span.End()

	fn := path.Join(ns, r.URL.Path)

	sublog := appctx.GetLogger(ctx).With().Str("path", fn).Str("svc", "ocdav").Str("handler", "get").Logger()

	space, status, err := spacelookup.LookUpStorageSpaceForPath(ctx, s.gwClient, fn)
	if err != nil {
		sublog.Error().Err(err).Str("path", fn).Msg("failed to look up storage space")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if status.Code != rpc.Code_CODE_OK {
		errors.HandleErrorStatus(&sublog, w, status)
		return
	}

	s.handleGet(ctx, w, r, spacelookup.MakeRelativeReference(space, fn, false), "spaces", sublog)
}

func (s *svc) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request, ref *provider.Reference, dlProtocol string, log zerolog.Logger) {
	sReq := &provider.StatRequest{
		Ref: ref,
	}
	sRes, err := s.gwClient.Stat(ctx, sReq)
	if err != nil {
		log.Error().Err(err).Msg("error stat resource")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if sRes.Status.Code != rpc.Code_CODE_OK {
		errors.HandleErrorStatus(&log, w, sRes.Status)
		return
	}

	if sRes.Info.Type != provider.ResourceType_RESOURCE_TYPE_FILE {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(http.StatusOK)
		return
	}

	if status := utils.ReadPlainFromOpaque(sRes.GetInfo().GetOpaque(), "status"); status == "processing" {
		w.WriteHeader(http.StatusTooEarly)
		return
	}

	dReq := &provider.InitiateFileDownloadRequest{Ref: ref}
	dRes, err := s.gwClient.InitiateFileDownload(ctx, dReq)
	switch {
	case err != nil:
		log.Error().Err(err).Msg("error initiating file download")
		w.WriteHeader(http.StatusInternalServerError)
		return
	case dRes.Status.Code != rpc.Code_CODE_OK:
		errors.HandleErrorStatus(&log, w, dRes.Status)
		return
	}

	var ep, token string
	for _, p := range dRes.Protocols {
		if p.Protocol == dlProtocol {
			ep, token = p.DownloadEndpoint, p.Token
		}
	}

	httpReq, err := rhttp.NewRequest(ctx, http.MethodGet, ep, nil)
	if err != nil {
		log.Error().Err(err).Msg("error creating http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpReq.Header.Set(datagateway.TokenTransportHeader, token)

	if r.Header.Get(net.HeaderRange) != "" {
		httpReq.Header.Set(net.HeaderRange, r.Header.Get(net.HeaderRange))
	}

	httpClient := s.client

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Msg("error performing http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer httpRes.Body.Close()

	copyHeader(w.Header(), httpRes.Header)
	w.WriteHeader(httpRes.StatusCode)

	if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusPartialContent {
		// swallow the body and set content-length to 0 to prevent reverse proxies from trying to read from it
		w.Header().Set("Content-Length", "0")
		return
	}

	var c int64
	if c, err = io.Copy(w, httpRes.Body); err != nil {
		log.Error().Err(err).Msg("error finishing copying data to response")
	}
	if httpRes.Header.Get(net.HeaderContentLength) != "" {
		i, err := strconv.ParseInt(httpRes.Header.Get(net.HeaderContentLength), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("content-length", httpRes.Header.Get(net.HeaderContentLength)).Msg("invalid content length in datagateway response")
		}
		if i != c {
			log.Error().Int64("content-length", i).Int64("transferred-bytes", c).Msg("content length vs transferred bytes mismatch")
		}
	}
	// TODO we need to send the If-Match etag in the GET to the datagateway to prevent race conditions between stating and reading the file
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		for i := range values {
			dst.Add(key, values[i])
		}
	}
}

func (s *svc) handleSpacesGet(w http.ResponseWriter, r *http.Request, spaceID string) {
	ctx, span := s.tracerProvider.Tracer(tracerName).Start(r.Context(), "spaces_get")
	defer span.End()

	sublog := appctx.GetLogger(ctx).With().Str("path", r.URL.Path).Str("spaceid", spaceID).Str("handler", "get").Logger()

	ref, err := spacelookup.MakeStorageSpaceReference(spaceID, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.handleGet(ctx, w, r, &ref, "spaces", sublog)
}

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
	"encoding/json"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typespb "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/internal/http/services/datagateway"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/errors"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/net"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocdav/spacelookup"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
	tusd "github.com/tus/tusd/pkg/handler"
	"go.opentelemetry.io/otel/propagation"
)

// Propagator ensures the importer module uses the same trace propagation strategy.
var Propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

func (s *svc) handlePathTusPost(w http.ResponseWriter, r *http.Request, ns string) {
	ctx, span := appctx.GetTracerProvider(r.Context()).Tracer(tracerName).Start(r.Context(), "tus-post")
	defer span.End()

	// read filename from metadata
	meta := tusd.ParseMetadataHeader(r.Header.Get(net.HeaderUploadMetadata))
	if err := ValidateName(meta["filename"], s.nameValidators); err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	// append filename to current dir
	fn := path.Join(ns, r.URL.Path, meta["filename"])

	sublog := appctx.GetLogger(ctx).With().Str("path", fn).Logger()
	// check tus headers?

	ref := &provider.Reference{
		// FIXME ResourceId?
		Path: fn,
	}
	s.handleTusPost(ctx, w, r, meta, ref, sublog)
}

func (s *svc) handleSpacesTusPost(w http.ResponseWriter, r *http.Request, spaceID string) {
	ctx, span := appctx.GetTracerProvider(r.Context()).Tracer(tracerName).Start(r.Context(), "spaces-tus-post")
	defer span.End()

	// read filename from metadata
	meta := tusd.ParseMetadataHeader(r.Header.Get(net.HeaderUploadMetadata))
	if err := ValidateName(meta["filename"], s.nameValidators); err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}

	sublog := appctx.GetLogger(ctx).With().Str("spaceid", spaceID).Str("path", r.URL.Path).Logger()

	ref, err := spacelookup.MakeStorageSpaceReference(spaceID, path.Join(r.URL.Path, meta["filename"]))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.handleTusPost(ctx, w, r, meta, &ref, sublog)
}

func (s *svc) handleTusPost(ctx context.Context, w http.ResponseWriter, r *http.Request, meta map[string]string, ref *provider.Reference, log zerolog.Logger) {
	w.Header().Add(net.HeaderAccessControlAllowHeaders, strings.Join([]string{net.HeaderTusResumable, net.HeaderUploadLength, net.HeaderUploadMetadata, net.HeaderIfMatch}, ", "))
	w.Header().Add(net.HeaderAccessControlExposeHeaders, strings.Join([]string{net.HeaderTusResumable, net.HeaderUploadOffset, net.HeaderLocation}, ", "))
	w.Header().Set(net.HeaderTusExtension, "creation,creation-with-upload,checksum,expiration")

	w.Header().Set(net.HeaderTusResumable, "1.0.0")

	// Test if the version sent by the client is supported
	// GET methods are not checked since a browser may visit this URL and does
	// not include this header. This request is not part of the specification.
	if r.Header.Get(net.HeaderTusResumable) != "1.0.0" {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}
	if r.Header.Get(net.HeaderUploadLength) == "" {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}
	// r.Header.Get(net.HeaderOCChecksum)
	// TODO must be SHA1, ADLER32 or MD5 ... in capital letters????
	// curl -X PUT https://demo.owncloud.com/remote.php/webdav/testcs.bin -u demo:demo -d '123' -v -H 'OC-Checksum: SHA1:40bd001563085fc35165329ea1ff5c5ecbdbbeef'

	// TODO check Expect: 100-continue

	client, err := s.gatewaySelector.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sReq := &provider.StatRequest{
		Ref: ref,
	}
	sRes, err := client.Stat(ctx, sReq)
	if err != nil {
		log.Error().Err(err).Msg("error sending grpc stat request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sRes.Status.Code != rpc.Code_CODE_OK && sRes.Status.Code != rpc.Code_CODE_NOT_FOUND {
		errors.HandleErrorStatus(&log, w, sRes.Status)
		return
	}

	info := sRes.Info
	if info != nil && info.Type != provider.ResourceType_RESOURCE_TYPE_FILE {
		log.Warn().Msg("resource is not a file")
		w.WriteHeader(http.StatusConflict)
		return
	}

	if info != nil {
		clientETag := r.Header.Get(net.HeaderIfMatch)
		serverETag := info.Etag
		if clientETag != "" {
			if clientETag != serverETag {
				log.Warn().Str("client-etag", clientETag).Str("server-etag", serverETag).Msg("etags mismatch")
				w.WriteHeader(http.StatusPreconditionFailed)
				return
			}
		}
	}

	uploadLength, err := strconv.ParseInt(r.Header.Get(net.HeaderUploadLength), 10, 64)
	if err != nil {
		log.Debug().Err(err).Msg("wrong request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if uploadLength == 0 {
		tfRes, err := client.TouchFile(ctx, &provider.TouchFileRequest{
			Ref: ref,
		})
		if err != nil {
			log.Error().Err(err).Msg("error sending grpc stat request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if tfRes.Status.Code != rpc.Code_CODE_OK {
			log.Error().Interface("status", tfRes.Status).Msg("error touching file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	opaqueMap := map[string]*typespb.OpaqueEntry{
		net.HeaderUploadLength: {
			Decoder: "plain",
			Value:   []byte(r.Header.Get(net.HeaderUploadLength)),
		},
	}

	mtime := meta["mtime"]
	if mtime != "" {
		opaqueMap[net.HeaderOCMtime] = &typespb.OpaqueEntry{
			Decoder: "plain",
			Value:   []byte(mtime),
		}
	}

	// initiateUpload
	uReq := &provider.InitiateFileUploadRequest{
		Ref: ref,
		Opaque: &typespb.Opaque{
			Map: opaqueMap,
		},
	}

	uRes, err := client.InitiateFileUpload(ctx, uReq)
	if err != nil {
		log.Error().Err(err).Msg("error initiating file upload")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if uRes.Status.Code != rpc.Code_CODE_OK {
		if r.ProtoMajor == 1 {
			// drain body to avoid `connection closed` errors
			_, _ = io.Copy(io.Discard, r.Body)
		}
		if uRes.Status.Code == rpc.Code_CODE_NOT_FOUND {
			w.WriteHeader(http.StatusPreconditionFailed)
			return
		}
		errors.HandleErrorStatus(&log, w, uRes.Status)
		return
	}

	var ep, token string
	for _, p := range uRes.Protocols {
		if p.Protocol == "tus" {
			ep, token = p.UploadEndpoint, p.Token
		}
	}

	// TUS clients don't understand the reva transfer token. We need to append it to the upload endpoint.
	// The DataGateway has to take care of pulling it back into the request header upon request arrival.
	if token != "" {
		if !strings.HasSuffix(ep, "/") {
			ep += "/"
		}
		ep += token
	}

	w.Header().Set(net.HeaderLocation, ep)

	// for creation-with-upload extension forward bytes to dataprovider
	// TODO check this really streams
	if r.Header.Get(net.HeaderContentType) == "application/offset+octet-stream" {
		length, err := strconv.ParseInt(r.Header.Get(net.HeaderContentLength), 10, 64)
		if err != nil {
			log.Debug().Err(err).Msg("wrong request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var httpRes *http.Response

		// if we know the transfer secret we can directly talk to the dataprovider
		if s.c.TransferSharedSecret != "" {
			claims, err := datagateway.Verify(ctx, token, s.c.TransferSharedSecret)
			if err != nil {
				log.Error().Err(err).Msg("error verifying transfer token")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// directly send request to target
			ep = claims.Target
		}

		httpReq, err := rhttp.NewRequest(ctx, http.MethodPatch, ep, r.Body)
		if err != nil {
			log.Debug().Err(err).Msg("wrong request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		Propagator.Inject(ctx, propagation.HeaderCarrier(httpReq.Header))

		httpReq.Header.Set(net.HeaderContentType, r.Header.Get(net.HeaderContentType))
		httpReq.Header.Set(net.HeaderContentLength, r.Header.Get(net.HeaderContentLength))
		if r.Header.Get(net.HeaderUploadOffset) != "" {
			httpReq.Header.Set(net.HeaderUploadOffset, r.Header.Get(net.HeaderUploadOffset))
		} else {
			httpReq.Header.Set(net.HeaderUploadOffset, "0")
		}
		httpReq.Header.Set(net.HeaderTusResumable, r.Header.Get(net.HeaderTusResumable))

		httpRes, err = s.client.Do(httpReq)
		if err != nil {
			log.Error().Err(err).Msg("error doing PATCH request to data gateway")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer httpRes.Body.Close()

		w.Header().Set(net.HeaderUploadOffset, httpRes.Header.Get(net.HeaderUploadOffset))
		w.Header().Set(net.HeaderTusResumable, httpRes.Header.Get(net.HeaderTusResumable))
		w.Header().Set(net.HeaderTusUploadExpires, httpRes.Header.Get(net.HeaderTusUploadExpires))
		if httpRes.StatusCode != http.StatusNoContent {
			w.WriteHeader(httpRes.StatusCode)
			return
		}

		// check if upload was fully completed
		if length == 0 || httpRes.Header.Get(net.HeaderUploadOffset) == r.Header.Get(net.HeaderUploadLength) {
			// get uploaded file metadata

			if resid, err := storagespace.ParseID(httpRes.Header.Get(net.HeaderOCFileID)); err == nil {
				sReq.Ref = &provider.Reference{
					ResourceId: &resid,
				}
			}

			sRes, err := client.Stat(ctx, sReq)
			if err != nil {
				log.Error().Err(err).Msg("error sending grpc stat request")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if sRes.Status.Code != rpc.Code_CODE_OK && sRes.Status.Code != rpc.Code_CODE_NOT_FOUND {

				if sRes.Status.Code == rpc.Code_CODE_PERMISSION_DENIED {
					// the token expired during upload, so the stat failed
					// and we can't do anything about it.
					// the clients will handle this gracefully by doing a propfind on the file
					w.WriteHeader(http.StatusOK)
					return
				}

				errors.HandleErrorStatus(&log, w, sRes.Status)
				return
			}

			info := sRes.Info
			if info == nil {
				log.Error().Msg("No info found for uploaded file")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if httpRes != nil && httpRes.Header != nil && httpRes.Header.Get(net.HeaderOCMtime) != "" {
				// set the "accepted" value if returned in the upload response headers
				w.Header().Set(net.HeaderOCMtime, httpRes.Header.Get(net.HeaderOCMtime))
			}

			// get WebDav permissions for file
			isPublic := false
			if info.Opaque != nil && info.Opaque.Map != nil {
				if info.Opaque.Map["link-share"] != nil && info.Opaque.Map["link-share"].Decoder == "json" {
					ls := &link.PublicShare{}
					_ = json.Unmarshal(info.Opaque.Map["link-share"].Value, ls)
					isPublic = ls != nil
				}
			}
			isShared := !net.IsCurrentUserOwner(ctx, info.Owner)
			role := conversions.RoleFromResourcePermissions(info.PermissionSet, isPublic)
			permissions := role.WebDAVPermissions(
				info.Type == provider.ResourceType_RESOURCE_TYPE_CONTAINER,
				isShared,
				false,
				isPublic,
			)

			w.Header().Set(net.HeaderContentType, info.MimeType)
			w.Header().Set(net.HeaderOCFileID, storagespace.FormatResourceID(*info.Id))
			w.Header().Set(net.HeaderOCETag, info.Etag)
			w.Header().Set(net.HeaderETag, info.Etag)
			w.Header().Set(net.HeaderOCPermissions, permissions)

			t := utils.TSToTime(info.Mtime).UTC()
			lastModifiedString := t.Format(time.RFC1123Z)
			w.Header().Set(net.HeaderLastModified, lastModifiedString)
		}
	}

	w.WriteHeader(http.StatusCreated)
}

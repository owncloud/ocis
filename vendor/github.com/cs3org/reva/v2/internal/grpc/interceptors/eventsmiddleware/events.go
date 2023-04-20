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

package eventsmiddleware

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"

	"go-micro.dev/v4/util/log"
	"google.golang.org/grpc"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	v1beta12 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-micro/plugins/v4/events/natsjs"
)

const (
	defaultPriority = 200
)

func init() {
	rgrpc.RegisterUnaryInterceptor("eventsmiddleware", NewUnary)
}

// NewUnary returns a new unary interceptor that emits events when needed
// no lint because of the switch statement that should be extendable
//
//nolint:gocritic
func NewUnary(m map[string]interface{}) (grpc.UnaryServerInterceptor, int, error) {
	publisher, err := publisherFromConfig(m)
	if err != nil {
		return nil, 0, err
	}

	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Register a channel in the context to receive the space owner id from the handler(s) further down the stack
		var ownerID *user.UserId
		sendOwnerChan := make(chan *user.UserId)
		ctx = storagespace.ContextRegisterSendOwnerChan(ctx, sendOwnerChan)

		res, err := handler(ctx, req)
		if err != nil {
			return res, err
		}

		// Read the space owner id from the channel
		select {
		case ownerID = <-sendOwnerChan:
		default:
		}

		var executantID *user.UserId
		u, ok := revactx.ContextGetUser(ctx)
		if ok {
			executantID = u.Id
		}

		var ev interface{}
		switch v := res.(type) {
		case *collaboration.CreateShareResponse:
			if isSuccess(v) {
				ev = ShareCreated(v, executantID)
			}
		case *collaboration.RemoveShareResponse:
			if isSuccess(v) {
				ev = ShareRemoved(v, req.(*collaboration.RemoveShareRequest), executantID)
			}
		case *collaboration.UpdateShareResponse:
			if isSuccess(v) {
				ev = ShareUpdated(v, req.(*collaboration.UpdateShareRequest), executantID)
			}
		case *collaboration.UpdateReceivedShareResponse:
			if isSuccess(v) {
				ev = ReceivedShareUpdated(v, executantID)
			}
		case *link.CreatePublicShareResponse:
			if isSuccess(v) {
				ev = LinkCreated(v, executantID)
			}
		case *link.UpdatePublicShareResponse:
			if isSuccess(v) {
				ev = LinkUpdated(v, req.(*link.UpdatePublicShareRequest), executantID)
			}
		case *link.RemovePublicShareResponse:
			if isSuccess(v) {
				ev = LinkRemoved(v, req.(*link.RemovePublicShareRequest), executantID)
			}
		case *link.GetPublicShareByTokenResponse:
			if isSuccess(v) {
				ev = LinkAccessed(v, executantID)
			} else {
				ev = LinkAccessFailed(v, req.(*link.GetPublicShareByTokenRequest), executantID)
			}
		case *provider.AddGrantResponse:
			// TODO: update CS3 APIs
			// FIXME these should be part of the RemoveGrantRequest object
			// https://github.com/owncloud/ocis/issues/4312
			r := req.(*provider.AddGrantRequest)
			if isSuccess(v) && utils.ExistsInOpaque(r.Opaque, "spacegrant") {
				ev = SpaceShared(v, r, executantID)
			}
		case *provider.RemoveGrantResponse:
			r := req.(*provider.RemoveGrantRequest)
			if isSuccess(v) && utils.ExistsInOpaque(r.Opaque, "spacegrant") {
				ev = SpaceUnshared(v, req.(*provider.RemoveGrantRequest), executantID)
			}
		case *provider.CreateContainerResponse:
			if isSuccess(v) {
				ev = ContainerCreated(v, req.(*provider.CreateContainerRequest), ownerID, executantID)
			}
		case *provider.InitiateFileDownloadResponse:
			if isSuccess(v) {
				ev = FileDownloaded(v, req.(*provider.InitiateFileDownloadRequest), executantID)
			}
		case *provider.DeleteResponse:
			if isSuccess(v) {
				ev = ItemTrashed(v, req.(*provider.DeleteRequest), ownerID, executantID)
			}
		case *provider.MoveResponse:
			if isSuccess(v) {
				ev = ItemMoved(v, req.(*provider.MoveRequest), ownerID, executantID)
			}
		case *provider.PurgeRecycleResponse:
			if isSuccess(v) {
				ev = ItemPurged(v, req.(*provider.PurgeRecycleRequest), executantID)
			}
		case *provider.RestoreRecycleItemResponse:
			if isSuccess(v) {
				ev = ItemRestored(v, req.(*provider.RestoreRecycleItemRequest), ownerID, executantID)
			}
		case *provider.RestoreFileVersionResponse:
			if isSuccess(v) {
				ev = FileVersionRestored(v, req.(*provider.RestoreFileVersionRequest), ownerID, executantID)
			}
		case *provider.CreateStorageSpaceResponse:
			if isSuccess(v) && v.StorageSpace != nil { // TODO: Why are there CreateStorageSpaceResponses with nil StorageSpace?
				ev = SpaceCreated(v, executantID)
			}
		case *provider.UpdateStorageSpaceResponse:
			if isSuccess(v) {
				r := req.(*provider.UpdateStorageSpaceRequest)
				if r.StorageSpace.Name != "" {
					ev = SpaceRenamed(v, r, executantID)
				} else if utils.ExistsInOpaque(r.Opaque, "restore") {
					ev = SpaceEnabled(v, r, executantID)
				} else {
					ev = SpaceUpdated(v, r, executantID)
				}
			}
		case *provider.DeleteStorageSpaceResponse:
			if isSuccess(v) {
				r := req.(*provider.DeleteStorageSpaceRequest)
				if utils.ExistsInOpaque(r.Opaque, "purge") {
					ev = SpaceDeleted(v, r, executantID)
				} else {
					ev = SpaceDisabled(v, r, executantID)
				}
			}
		case *provider.TouchFileResponse:
			if isSuccess(v) {
				ev = FileTouched(v, req.(*provider.TouchFileRequest), ownerID, executantID)
			}
		}

		if ev != nil {
			if err := events.Publish(publisher, ev); err != nil {
				log.Error(err)
			}
		}

		return res, nil
	}
	return interceptor, defaultPriority, nil
}

// NewStream returns a new server stream interceptor
// that creates the application context.
func NewStream() grpc.StreamServerInterceptor {
	interceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO: Use ss.RecvMsg() and ss.SendMsg() to send events from a stream
		return handler(srv, ss)
	}
	return interceptor
}

// common interface to all responses
type su interface {
	GetStatus() *v1beta12.Status
}

func isSuccess(res su) bool {
	return res.GetStatus().Code == rpc.Code_CODE_OK
}

func publisherFromConfig(m map[string]interface{}) (events.Publisher, error) {
	typ := m["type"].(string)
	switch typ {
	default:
		return nil, fmt.Errorf("stream type '%s' not supported", typ)
	case "nats":
		address := m["address"].(string)
		cid := m["clusterID"].(string)

		enableTLS := m["enable-tls"].(bool)
		var tlsConf *tls.Config
		if enableTLS {
			skipVerify := m["tls-insecure"].(bool)
			var rootCAPool *x509.CertPool
			if val, ok := m["tls-root-ca-cert"]; ok {
				rootCACertPath := val.(string)
				if rootCACertPath != "" {
					f, err := os.Open(rootCACertPath)
					if err != nil {
						return nil, err
					}

					var certBytes bytes.Buffer
					if _, err := io.Copy(&certBytes, f); err != nil {
						return nil, err
					}

					rootCAPool = x509.NewCertPool()
					rootCAPool.AppendCertsFromPEM(certBytes.Bytes())
					skipVerify = false
				}
			}

			tlsConf = &tls.Config{
				InsecureSkipVerify: skipVerify,
				RootCAs:            rootCAPool,
			}
		}
		return stream.Nats(natsjs.TLSConfig(tlsConf), natsjs.Address(address), natsjs.ClusterID(cid))
	}
}

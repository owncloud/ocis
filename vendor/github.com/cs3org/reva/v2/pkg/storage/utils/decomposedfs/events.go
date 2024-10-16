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

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/cs3org/reva/v2/pkg/utils"
)

func (fs *Decomposedfs) publishEvent(ctx context.Context, evf func() (any, error)) {
	log := appctx.GetLogger(ctx)
	if fs.stream == nil {
		log.Error().Msg("Failed to publish event, stream is undefined")
		return
	}
	ev, err := evf()
	if err != nil || ev == nil {
		log.Error().Err(err).Msg("Failed to crete the event")
		return
	}
	if err := events.Publish(ctx, fs.stream, ev); err != nil {
		log.Error().Err(err).Msg("Failed to publish event")
	}
}

func (fs *Decomposedfs) moveEvent(ctx context.Context, oldRef, newRef *provider.Reference, oldNode, newNode *node.Node, orp, nrp *provider.ResourcePermissions) func() (any, error) {
	return func() (any, error) {
		executant, _ := revactx.ContextGetUser(ctx)
		ev := events.ItemMoved{
			SpaceOwner:        newNode.Owner(),
			Executant:         executant.GetId(),
			Ref:               newRef,
			OldReference:      oldRef,
			Timestamp:         utils.TSNow(),
			ImpersonatingUser: extractImpersonator(executant),
		}
		log := appctx.GetLogger(ctx)
		if nref, err := fs.refFromNode(ctx, newNode, newRef.GetResourceId().GetStorageId(), nrp); err == nil {
			ev.Ref = nref
		} else {
			log.Error().Err(err).Msg("Failed to get destination reference")
		}

		if oref, err := fs.refFromNode(ctx, oldNode, oldRef.GetResourceId().GetStorageId(), orp); err == nil {
			ev.OldReference = oref
		} else {
			log.Error().Err(err).Msg("Failed to get source reference")
		}

		return ev, nil
	}
}

func extractImpersonator(u *user.User) *user.User {
	var impersonator user.User
	if err := utils.ReadJSONFromOpaque(u.Opaque, "impersonating-user", &impersonator); err != nil {
		return nil
	}
	return &impersonator
}

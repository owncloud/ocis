package provider

import (
	"context"
	"sync"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/errtypes"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"google.golang.org/grpc/metadata"
)

// SpaceDebouncer debounces operations on spaces for a configurable amount of time
type SpaceDebouncer struct {
	after      time.Duration
	f          func(id *provider.StorageSpaceId, userID *user.UserId)
	pending    map[string]*time.Timer
	inProgress sync.Map

	mutex sync.Mutex
}

// NewSpaceDebouncer returns a new SpaceDebouncer instance
func NewSpaceDebouncer(d time.Duration, f func(id *provider.StorageSpaceId, userID *user.UserId)) *SpaceDebouncer {
	return &SpaceDebouncer{
		after:      d,
		f:          f,
		pending:    map[string]*time.Timer{},
		inProgress: sync.Map{},
	}
}

// Debounce restars the debounce timer for the given space
func (d *SpaceDebouncer) Debounce(id *provider.StorageSpaceId, userID *user.UserId) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if t := d.pending[id.OpaqueId]; t != nil {
		t.Stop()
	}

	d.pending[id.OpaqueId] = time.AfterFunc(d.after, func() {
		if _, ok := d.inProgress.Load(id.OpaqueId); ok {
			// Reschedule this run for when the previous run has finished
			d.mutex.Lock()
			d.pending[id.OpaqueId].Reset(d.after)
			d.mutex.Unlock()
			return
		}

		d.inProgress.Store(id.OpaqueId, true)
		defer d.inProgress.Delete(id.OpaqueId)
		d.f(id, userID)
	})
}

func (p *Provider) handleEvent(ev interface{}) {
	switch e := ev.(type) {
	case events.ItemTrashed:
		p.logger.Debug().Interface("event", ev).Msg("marking document as deleted")
		err := p.indexClient.Delete(e.ID)
		if err != nil {
			p.logger.Error().Err(err).Interface("Id", e.ID).Msg("failed to remove item from index")
		}
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.ItemRestored:
		p.logger.Debug().Interface("event", ev).Msg("marking document as restored")
		owner := &user.User{
			Id: e.Executant,
		}

		ownerCtx, err := p.getAuthContext(owner)
		if err != nil {
			return
		}
		statRes, err := p.statResource(ownerCtx, e.Ref, owner)
		if err != nil {
			p.logger.Error().Err(err).
				Str("storageid", e.Ref.GetResourceId().GetStorageId()).
				Str("spaceid", e.Ref.GetResourceId().GetSpaceId()).
				Str("opaqueid", e.Ref.GetResourceId().GetOpaqueId()).
				Str("path", e.Ref.GetPath()).
				Msg("failed to make stat call for the restored resource")
			return
		}

		switch statRes.Status.Code {
		case rpc.Code_CODE_OK:
			err = p.indexClient.Restore(statRes.Info.Id)
			if err != nil {
				p.logger.Error().Err(err).
					Str("storageid", e.Ref.GetResourceId().GetStorageId()).
					Str("spaceid", e.Ref.GetResourceId().GetSpaceId()).
					Str("opaqueid", e.Ref.GetResourceId().GetOpaqueId()).
					Str("path", e.Ref.GetPath()).
					Msg("failed to restore the changed resource in the index")
			}
		default:
			p.logger.Error().Interface("statRes", statRes).
				Str("storageid", e.Ref.GetResourceId().GetStorageId()).
				Str("spaceid", e.Ref.GetResourceId().GetSpaceId()).
				Str("opaqueid", e.Ref.GetResourceId().GetOpaqueId()).
				Str("path", e.Ref.GetPath()).
				Msg("failed to stat the restored resource")
		}
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.ItemMoved:
		p.logger.Debug().Interface("event", ev).Msg("resource has been moved, updating the document")
		owner := &user.User{
			Id: e.Executant,
		}

		ownerCtx, err := p.getAuthContext(owner)
		if err != nil {
			return
		}
		statRes, err := p.statResource(ownerCtx, e.Ref, owner)
		if err != nil {
			p.logger.Error().Err(err).Msg("failed to stat the moved resource")
			return
		}
		if statRes.Status.Code != rpc.Code_CODE_OK {
			p.logger.Error().Interface("statRes", statRes).Msg("failed to stat the moved resource")
			return
		}

		gpRes, err := p.getPath(ownerCtx, statRes.Info.Id, owner)
		if err != nil {
			p.logger.Error().Err(err).Interface("ref", e.Ref).Msg("failed to get path for moved resource")
			return
		}
		if gpRes.Status.Code != rpcv1beta1.Code_CODE_OK {
			p.logger.Error().Interface("status", gpRes.Status).Interface("ref", e.Ref).Msg("failed to get path for moved resource")
			return
		}

		err = p.indexClient.Move(statRes.GetInfo().GetId(), statRes.GetInfo().GetParentId(), gpRes.Path)
		if err != nil {
			p.logger.Error().Err(err).Msg("failed to move the changed resource in the index")
		}
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.ContainerCreated:
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.FileUploaded:
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.FileTouched:
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	case events.FileVersionRestored:
		p.reindexSpace(ev, e.Ref, e.Executant, e.SpaceOwner)
	default:
		// Not sure what to do here. Skip.
		return
	}
}

func (p *Provider) reindexSpace(ev interface{}, ref *provider.Reference, executant, owner *user.UserId) {
	p.logger.Debug().Interface("event", ev).Msg("resource has been changed, scheduling a space resync")

	spaceID := &provider.StorageSpaceId{
		OpaqueId: storagespace.FormatResourceID(provider.ResourceId{
			StorageId: ref.GetResourceId().GetStorageId(),
			SpaceId:   ref.GetResourceId().GetSpaceId(),
		}),
	}
	if owner != nil {
		p.indexSpaceDebouncer.Debounce(spaceID, owner)
	} else {
		p.indexSpaceDebouncer.Debounce(spaceID, executant)
	}
}

func (p *Provider) statResource(ctx context.Context, ref *provider.Reference, owner *user.User) (*provider.StatResponse, error) {
	return p.gwClient.Stat(ctx, &provider.StatRequest{Ref: ref})
}

func (p *Provider) getPath(ctx context.Context, id *provider.ResourceId, owner *user.User) (*provider.GetPathResponse, error) {
	return p.gwClient.GetPath(ctx, &provider.GetPathRequest{ResourceId: id})
}

func (p *Provider) getAuthContext(owner *user.User) (context.Context, error) {
	ownerCtx := ctxpkg.ContextSetUser(context.Background(), owner)
	authRes, err := p.gwClient.Authenticate(ownerCtx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + owner.GetId().GetOpaqueId(),
		ClientSecret: p.machineAuthAPIKey,
	})
	if err == nil && authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		err = errtypes.NewErrtypeFromStatus(authRes.Status)
	}
	if err != nil {
		p.logger.Error().Err(err).Interface("owner", owner).Interface("authRes", authRes).Msg("error using machine auth")
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ownerCtx, ctxpkg.TokenHeader, authRes.Token), nil
}

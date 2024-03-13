package command

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

var _registeredEvents = []events.Unmarshaller{
	events.ShareCreated{},
}

// ListenForEvents listens for events and acts accordingly
func ListenForEvents(ctx context.Context, cfg *config.Config, l log.Logger) error {
	bus, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Events))
	if err != nil {
		l.Error().Err(err).Msg("cannot connect to nats")
		return err
	}

	evChannel, err := events.Consume(bus, "frontend", _registeredEvents...)
	if err != nil {
		l.Error().Err(err).Msg("cannot consume from nats")
		return err
	}

	tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
	if err != nil {
		return err
	}

	gatewaySelector, err := pool.GatewaySelector(
		cfg.Reva.Address,
		pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
		pool.WithTLSMode(tm),
		pool.WithRegistry(registry.GetRegistry()),
	)
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway selector")
		return err
	}

	traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
	if err != nil {
		l.Error().Err(err).Msg("cannot initialize tracing")
		return err
	}

	grpcClient, err := grpc.NewClient(
		append(
			grpc.GetClientOptions(cfg.GRPCClientTLS),
			grpc.WithTraceProvider(traceProvider),
		)...,
	)
	if err != nil {
		l.Error().Err(err).Msg("cannot create grpc client")
		return err
	}

	valueService := settingssvc.NewValueService("com.owncloud.api.settings", grpcClient)

	for {
		select {
		case e := <-evChannel:
			switch ev := e.Event.(type) {
			default:
				l.Error().Interface("event", e).Msg("unhandled event")
			case events.ShareCreated:
				AutoAcceptShares(ev, cfg.AutoAcceptShares, l, gatewaySelector, valueService, cfg.ServiceAccount)
			}
		case <-ctx.Done():
			l.Info().Msg("context cancelled")
			return ctx.Err()
		}
	}
}

// AutoAcceptShares automatically accepts shares if configured by the admin or user
func AutoAcceptShares(ev events.ShareCreated, autoAcceptDefault bool, l log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], vs settingssvc.ValueService, cfg config.ServiceAccount) {
	gwc, err := gatewaySelector.Next()
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway client")
		return
	}
	ctx, err := utils.GetServiceUserContext(cfg.ServiceAccountID, gwc, cfg.ServiceAccountSecret)
	if err != nil {
		l.Error().Err(err).Msg("cannot impersonate user")
		return
	}

	uids, err := getUserIDs(ctx, gatewaySelector, ev.GranteeUserID, ev.GranteeGroupID)
	if err != nil {
		l.Error().Err(err).Msg("cannot get grantees")
		return
	}

	gwc, err = gatewaySelector.Next()
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway client")
		return
	}
	info, err := utils.GetResourceByID(ctx, ev.ItemID, gwc)
	if err != nil {
		l.Error().Err(err).Msg("error getting resource")
		return
	}

	for _, uid := range uids {
		if !autoAcceptShares(ctx, uid, autoAcceptDefault, vs) {
			continue
		}

		mountpoint, err := getMountpoint(ctx, l, ev.ItemID, uid, gatewaySelector, info)
		if err != nil {
			l.Error().Err(err).Msg("error getting mountpoint")
			continue

		}

		gwc, err := gatewaySelector.Next()
		if err != nil {
			l.Error().Err(err).Msg("cannot get gateway client")
			continue
		}
		resp, err := gwc.UpdateReceivedShare(ctx, updateShareRequest(ev.ShareID, uid, mountpoint))
		if err != nil {
			l.Error().Err(err).Msg("error sending grpc request")
			continue
		}

		if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
			l.Error().Interface("status", resp.GetStatus()).Str("userid", uid.GetOpaqueId()).Msg("unexpected status code while accepting share")
		}
	}

}

func getMountpoint(ctx context.Context, l log.Logger, itemid *provider.ResourceId, uid *user.UserId, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], info *provider.ResourceInfo) (string, error) {
	lrs, err := getSharesList(ctx, gatewaySelector, uid)
	if err != nil {
		return "", err
	}

	// we need to sort the received shares by mount point in order to make things easier to evaluate.
	base := path.Base(info.GetPath())
	mount := base
	mounts := make([]string, 0, len(lrs.Shares))
	var exists bool

	for _, s := range lrs.Shares {
		if s.State != collaboration.ShareState_SHARE_STATE_ACCEPTED {
			// we don't care about unaccepted shares
			continue
		}

		if utils.ResourceIDEqual(s.GetShare().GetResourceId(), itemid) {
			// a share to the same resource already exists and is mounted
			return s.GetMountPoint().GetPath(), nil
		}

		if s.GetMountPoint().GetPath() == mount {
			// does the shared resource still exist?
			gwc, err := gatewaySelector.Next()
			if err != nil {
				l.Error().Err(err).Msg("cannot get gateway client")
				continue
			}
			_, err = utils.GetResourceByID(ctx, s.GetShare().GetResourceId(), gwc)
			if err == nil {
				exists = true
			}
			// TODO we could delete shares here if the stat returns code NOT FOUND ... but listening for file deletes would be better
		}
		// collect all mount points
		mounts = append(mounts, s.GetMountPoint().GetPath())
	}

	// If the mount point really already exists, we need to insert a number into the filename
	if exists {
		// now we have a list of shares, we want to iterate over all of them and check for name collisions agents a mount points list
		for i := 1; i <= len(mounts)+1; i++ {
			ext := filepath.Ext(base)
			name := strings.TrimSuffix(base, ext)
			// be smart about .tar.(gz|bz) files
			if strings.HasSuffix(name, ".tar") {
				name = strings.TrimSuffix(name, ".tar")
				ext = ".tar" + ext
			}
			mount = name + " (" + strconv.Itoa(i) + ")" + ext
			if !slices.Contains(mounts, mount) {
				return mount, nil
			}
		}
	}
	return mount, nil
}

func getUserIDs(ctx context.Context, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], uid *user.UserId, gid *group.GroupId) ([]*user.UserId, error) {
	if uid != nil {
		return []*user.UserId{uid}, nil
	}

	gwc, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	res, err := gwc.GetGroup(ctx, &group.GetGroupRequest{GroupId: gid})
	if err != nil {
		return nil, err
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, errors.New("could not get group")
	}

	return res.GetGroup().GetMembers(), nil
}

func autoAcceptShares(ctx context.Context, uid *user.UserId, defaultValue bool, vs settingssvc.ValueService) bool {
	granteeCtx := metadata.Set(ctx, middleware.AccountID, uid.GetOpaqueId())
	if resp, err := vs.GetValueByUniqueIdentifiers(granteeCtx,
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: uid.GetOpaqueId(),
			SettingId:   defaults.SettingUUIDProfileAutoAcceptShares,
		},
	); err == nil {
		return resp.GetValue().GetValue().GetBoolValue()

	}
	return defaultValue
}

func updateShareRequest(shareID *collaboration.ShareId, uid *user.UserId, path string) *collaboration.UpdateReceivedShareRequest {
	return &collaboration.UpdateReceivedShareRequest{
		Opaque: utils.AppendJSONToOpaque(nil, "userid", uid),
		Share: &collaboration.ReceivedShare{
			Share: &collaboration.Share{
				Id: shareID,
			},
			MountPoint: &provider.Reference{
				Path: path,
			},
			State: collaboration.ShareState_SHARE_STATE_ACCEPTED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"state", "mount_point"}},
	}
}

// getSharesList gets the list of all shares for the given user.
func getSharesList(ctx context.Context, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], uid *user.UserId) (*collaboration.ListReceivedSharesResponse, error) {
	gwc, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}
	shares, err := gwc.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{
		Opaque: utils.AppendJSONToOpaque(nil, "userid", uid),
	})
	if err != nil {
		return nil, err
	}

	if shares.Status.Code != rpc.Code_CODE_OK {
		if shares.Status.Code == rpc.Code_CODE_NOT_FOUND {
			return nil, fmt.Errorf("not found")
		}
		return nil, fmt.Errorf(shares.GetStatus().GetMessage())
	}
	return shares, nil
}

package command

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"sort"
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
	bus, err := stream.NatsFromConfig(cfg.Service.Name, stream.NatsConfig(cfg.Events))
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

	gwc, err := gatewaySelector.Next()
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway client")
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
				AutoAcceptShares(ev, cfg.AutoAcceptShares, l, gwc, valueService, cfg.ServiceAccount)
			}
		case <-ctx.Done():
			l.Info().Msg("context cancelled")
			return ctx.Err()
		}
	}
}

// AutoAcceptShares automatically accepts shares if configured by the admin or user
func AutoAcceptShares(ev events.ShareCreated, autoAcceptDefault bool, l log.Logger, gwc gateway.GatewayAPIClient, vs settingssvc.ValueService, cfg config.ServiceAccount) {
	ctx, err := utils.GetServiceUserContext(cfg.ServiceAccountID, gwc, cfg.ServiceAccountSecret)
	if err != nil {
		l.Error().Err(err).Msg("cannot impersonate user")
		return
	}

	uids, err := getUserIDs(ctx, gwc, ev.GranteeUserID, ev.GranteeGroupID)
	if err != nil {
		l.Error().Err(err).Msg("cannot get granteess")
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

		mountpoint, err := getMountpoint(ctx, ev.ItemID, uid, gwc, info)
		if err != nil {
			l.Error().Err(err).Msg("error getting mountpoint")
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

func getMountpoint(ctx context.Context, itemid *provider.ResourceId, uid *user.UserId, gwc gateway.GatewayAPIClient, info *provider.ResourceInfo) (string, error) {
	lrs, err := getSharesList(ctx, gwc, uid)
	if err != nil {
		return "", err
	}

	// we need to sort the received shares by mount point in order to make things easier to evaluate.
	base := path.Base(info.GetPath())
	mount := base
	mountedShares := make([]*collaboration.ReceivedShare, 0, len(lrs.Shares))
	for _, s := range lrs.Shares {
		if s.State != collaboration.ShareState_SHARE_STATE_ACCEPTED {
			// we don't care about unaccepted shares
			continue
		}

		if utils.ResourceIDEqual(s.Share.ResourceId, itemid) {
			// a share to the same resource already exists and is mounted
			return s.MountPoint.Path, nil
		}

		mountedShares = append(mountedShares, s)
	}

	sort.Slice(mountedShares, func(i, j int) bool {
		return mountedShares[i].MountPoint.Path > mountedShares[j].MountPoint.Path
	})

	// now we have a list of shares, we want to iterate over all of them and check for name collisions
	for i, ms := range mountedShares {
		if ms.MountPoint.Path == mount {
			// does the shared resource still exist?
			_, err := utils.GetResourceByID(ctx, ms.Share.ResourceId, gwc)
			if err == nil {
				// The mount point really already exists, we need to insert a number into the filename
				ext := filepath.Ext(base)
				name := strings.TrimSuffix(base, ext)
				// be smart about .tar.(gz|bz) files
				if strings.HasSuffix(name, ".tar") {
					name = strings.TrimSuffix(name, ".tar")
					ext = ".tar" + ext
				}

				mount = fmt.Sprintf("%s (%s)%s", name, strconv.Itoa(i+1), ext)
			}
			// TODO we could delete shares here if the stat returns code NOT FOUND ... but listening for file deletes would be better
		}
	}
	return mount, nil
}

func getUserIDs(ctx context.Context, gwc gateway.GatewayAPIClient, uid *user.UserId, gid *group.GroupId) ([]*user.UserId, error) {
	if uid != nil {
		return []*user.UserId{uid}, nil
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
func getSharesList(ctx context.Context, client gateway.GatewayAPIClient, uid *user.UserId) (*collaboration.ListReceivedSharesResponse, error) {
	shares, err := client.ListReceivedShares(ctx, &collaboration.ListReceivedSharesRequest{
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

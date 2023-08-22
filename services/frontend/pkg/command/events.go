package command

import (
	"context"
	"errors"

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
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

var _registeredEvents = []events.Unmarshaller{
	events.ShareCreated{},
}

// ListenForEvents listens for events and acts accordingly
func ListenForEvents(cfg *config.Config, l log.Logger) {
	bus, err := stream.NatsFromConfig(cfg.Service.Name, stream.NatsConfig(cfg.Events))
	if err != nil {
		l.Error().Err(err).Msg("cannot connect to nats")
		return
	}

	evChannel, err := events.Consume(bus, "frontend", _registeredEvents...)
	if err != nil {
		l.Error().Err(err).Msg("cannot consume from nats")
	}

	tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
	if err != nil {
		return
	}

	gatewaySelector, err := pool.GatewaySelector(
		cfg.Reva.Address,
		pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
		pool.WithTLSMode(tm),
		pool.WithRegistry(registry.GetRegistry()),
	)
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway selector")
		return
	}

	gwc, err := gatewaySelector.Next()
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway client")
		return
	}

	traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
	if err != nil {
		l.Error().Err(err).Msg("cannot initialize tracing")
		return
	}

	grpcClient, err := grpc.NewClient(
		append(
			grpc.GetClientOptions(cfg.GRPCClientTLS),
			grpc.WithTraceProvider(traceProvider),
		)...,
	)
	if err != nil {
		l.Error().Err(err).Msg("cannot create grpc client")
		return
	}

	valueService := settingssvc.NewValueService("com.owncloud.api.settings", grpcClient)

	for e := range evChannel {
		switch ev := e.Event.(type) {
		default:
			l.Error().Interface("event", e).Msg("unhandled event")
		case events.ShareCreated:
			AutoAcceptShares(ev, cfg.AutoAcceptShares, l, gwc, valueService, cfg.ServiceAccount)
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

	for _, uid := range uids {
		if !autoAcceptShares(ctx, uid, autoAcceptDefault, vs) {
			continue
		}

		resp, err := gwc.UpdateReceivedShare(ctx, updateShareRequest(ev.ShareID, uid))
		if err != nil {
			l.Error().Err(err).Msg("error sending grpc request")
			continue
		}

		if resp.GetStatus().GetCode() != rpc.Code_CODE_OK {
			l.Error().Interface("status", resp.GetStatus()).Str("userid", uid.GetOpaqueId()).Msg("unexpected status code while accepting share")
		}
	}

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

func autoAcceptShares(ctx context.Context, u *user.UserId, defaultValue bool, vs settingssvc.ValueService) bool {
	granteeCtx := metadata.Set(ctx, middleware.AccountID, u.OpaqueId)
	if resp, err := vs.GetValueByUniqueIdentifiers(granteeCtx,
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: u.OpaqueId,
			SettingId:   defaults.SettingUUIDProfileAutoAcceptShares,
		},
	); err == nil {
		return resp.GetValue().GetValue().GetBoolValue()

	}
	return defaultValue
}

func updateShareRequest(shareID *collaboration.ShareId, uid *user.UserId) *collaboration.UpdateReceivedShareRequest {
	return &collaboration.UpdateReceivedShareRequest{
		Opaque: utils.AppendJSONToOpaque(nil, "userid", uid),
		Share: &collaboration.ReceivedShare{
			Share: &collaboration.Share{
				Id: shareID,
			},
			State: collaboration.ShareState_SHARE_STATE_ACCEPTED,
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"state"}},
	}
}

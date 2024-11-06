package command

import (
	"context"
	"errors"
	"sync"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/rs/zerolog"
	"go-micro.dev/v4/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/frontend/pkg/config"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
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

	traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
	if err != nil {
		l.Error().Err(err).Msg("cannot initialize tracing")
		return err
	}

	gatewaySelector, err := pool.GatewaySelector(
		cfg.Reva.Address,
		pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
		pool.WithTLSMode(tm),
		pool.WithRegistry(registry.GetRegistry()),
		pool.WithTracerProvider(traceProvider),
	)
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway selector")
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
				AutoAcceptShares(ev, cfg.AutoAcceptShares, l, gatewaySelector, valueService, cfg.ServiceAccount, cfg.MaxConcurrency)
			}
		case <-ctx.Done():
			l.Info().Msg("context cancelled")
			return nil
		}
	}
}

// AutoAcceptShares automatically accepts shares if configured by the admin or user
func AutoAcceptShares(ev events.ShareCreated, autoAcceptDefault bool, l log.Logger, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], vs settingssvc.ValueService, cfg config.ServiceAccount, maxConcurrency int) {
	gwc, err := gatewaySelector.Next()
	if err != nil {
		l.Error().Err(err).Msg("cannot get gateway client")
		return
	}
	ctx, err := utils.GetServiceUserContextWithContext(context.Background(), gwc, cfg.ServiceAccountID, cfg.ServiceAccountSecret)
	if err != nil {
		l.Error().Err(err).Msg("cannot impersonate user")
		return
	}

	userIDs, err := getUserIDs(ctx, gatewaySelector, ev.GranteeUserID, ev.GranteeGroupID)
	if err != nil {
		l.Error().Err(err).Msg("cannot get grantees")
		return
	}

	work := make(chan *user.UserId, len(userIDs))

	// Distribute work
	go func() {
		defer close(work)
		for _, userID := range userIDs {
			select {
			case work <- userID:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg := sync.WaitGroup{}
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range work {

				if !autoAcceptShares(ctx, userID, autoAcceptDefault, vs) {
					continue
				}

				gwc, err := gatewaySelector.Next()
				if err != nil {
					l.Error().Err(err).Msg("cannot get gateway client")
					continue
				}
				resp, err := gwc.UpdateReceivedShare(ctx, updateShareRequest(ev.ShareID, userID))
				if err != nil {
					l.Error().Err(err).Msg("error sending grpc request")
					continue
				}

				if code := resp.GetStatus().GetCode(); code != rpc.Code_CODE_OK {
					// log unexpected status codes if a share cannot be accepted...
					func() *zerolog.Event {
						switch code {
						// ... not found is not an error in the context of auto-accepting shares
						case rpc.Code_CODE_NOT_FOUND:
							return l.Debug()
						default:
							return l.Error()
						}
					}().Interface("status", resp.GetStatus()).Str("userid", userID.GetOpaqueId()).Msg("unexpected status code while accepting share")
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
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

package service

import (
	"context"
	"errors"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	micrometadata "go-micro.dev/v4/metadata"
)

type notificationFilter struct {
	log         log.Logger
	valueClient settingssvc.ValueService
}

func newNotificationFilter(l log.Logger, vc settingssvc.ValueService) *notificationFilter {
	return &notificationFilter{log: l, valueClient: vc}
}

// execute removes users who have disabled mail notifications for the event
func (nf notificationFilter) execute(ctx context.Context, users []*user.User, settingId string) []*user.User {
	var filteredUsers []*user.User

	for _, u := range users {
		userId := u.GetId().GetOpaqueId()
		enabled, err := getSetting(ctx, nf.valueClient, userId, settingId)
		if err != nil {
			nf.log.Error().Err(err).Str("userId", userId).Str("settingId", settingId).Msg("cannot get user event setting")
			filteredUsers = append(filteredUsers, u)
			continue
		}
		if enabled {
			filteredUsers = append(filteredUsers, u)
		}
	}

	return filteredUsers
}

func getSetting(ctx context.Context, vc settingssvc.ValueService, userId string, settingId string) (bool, error) {
	resp, err := vc.GetValueByUniqueIdentifiers(
		micrometadata.Set(ctx, middleware.AccountID, userId),
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: userId,
			SettingId:   settingId,
		},
	)

	if err != nil {
		return false, err
	}

	val := resp.GetValue().GetValue().GetCollectionValue().GetValues()
	for _, option := range val {
		if option.GetKey() == "mail" {
			return option.GetBoolValue(), nil
		}
	}
	return false, errors.New("no setting found")
}

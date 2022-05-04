package service

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"

	ssvc "github.com/owncloud/ocis/v2/extensions/settings/pkg/service/v0"
	olog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
)

const (
	// AccountManagementPermissionID is the hardcoded setting UUID for the account management permission
	AccountManagementPermissionID string = "8e587774-d929-4215-910b-a317b1e80f73"
	// AccountManagementPermissionName is the hardcoded setting name for the account management permission
	AccountManagementPermissionName string = "account-management"
	// GroupManagementPermissionID is the hardcoded setting UUID for the group management permission
	GroupManagementPermissionID string = "522adfbe-5908-45b4-b135-41979de73245"
	// GroupManagementPermissionName is the hardcoded setting name for the group management permission
	GroupManagementPermissionName string = "group-management"
	// SelfManagementPermissionID is the hardcoded setting UUID for the self management permission
	SelfManagementPermissionID string = "e03070e9-4362-4cc6-a872-1c7cb2eb2b8e"
	// SelfManagementPermissionName is the hardcoded setting name for the self management permission
	SelfManagementPermissionName string = "self-management"
)

// RegisterPermissions registers permissions for account management and group management with the settings service.
func RegisterPermissions(l *olog.Logger) {
	service := settingssvc.NewBundleService("com.owncloud.api.settings", grpc.DefaultClient)

	permissionRequests := generateAccountManagementPermissionsRequests()
	for i := range permissionRequests {
		res, err := service.AddSettingToBundle(context.Background(), &permissionRequests[i])
		bundleID := permissionRequests[i].BundleId
		if err != nil {
			l.Err(err).Str("bundle", bundleID).Str("setting", permissionRequests[i].Setting.Id).Msg("error adding permission to bundle")
		} else {
			l.Info().Str("bundle", bundleID).Str("setting", res.Setting.Id).Msg("successfully added permission to bundle")
		}
	}
}

func generateAccountManagementPermissionsRequests() []settingssvc.AddSettingToBundleRequest {
	return []settingssvc.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settingsmsg.Setting{
				Id:          AccountManagementPermissionID,
				Name:        AccountManagementPermissionName,
				DisplayName: "Account Management",
				Description: "This permission gives full access to everything that is related to account management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settingsmsg.Setting{
				Id:          GroupManagementPermissionID,
				Name:        GroupManagementPermissionName,
				DisplayName: "Group Management",
				Description: "This permission gives full access to everything that is related to group management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_GROUP,
					Id:   "all",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleUser,
			Setting: &settingsmsg.Setting{
				Id:          SelfManagementPermissionID,
				Name:        SelfManagementPermissionName,
				DisplayName: "Self Management",
				Description: "This permission gives access to self management.",
				Resource: &settingsmsg.Resource{
					Type: settingsmsg.Resource_TYPE_USER,
					Id:   "me",
				},
				Value: &settingsmsg.Setting_PermissionValue{
					PermissionValue: &settingsmsg.Permission{
						Operation:  settingsmsg.Permission_OPERATION_READWRITE,
						Constraint: settingsmsg.Permission_CONSTRAINT_OWN,
					},
				},
			},
		},
	}
}

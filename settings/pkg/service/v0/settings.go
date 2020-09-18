package svc

import settings "github.com/owncloud/ocis/settings/pkg/proto/v0"

const (
	// BundleUUIDRoleAdmin represents the admin role
	BundleUUIDRoleAdmin = "71881883-1768-46bd-a24d-a356a2afdf7f"

	// BundleUUIDRoleUser represents the user role.
	BundleUUIDRoleUser = "d7beeea8-8ff4-406b-8fb6-ab2dd81e6b11"

	// BundleUUIDRoleGuest represents the guest role.
	BundleUUIDRoleGuest = "38071a68-456a-4553-846a-fa67bf5596cc"
)

// generateBundlesDefaultRoles bootstraps the default roles.
func generateBundlesDefaultRoles() []*settings.Bundle {
	return []*settings.Bundle{
		generateBundleAdminRole(),
		generateBundleUserRole(),
		generateBundleGuestRole(),
	}
}

func generateBundleAdminRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleAdmin,
		Name:        "admin",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Admin",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}

func generateBundleUserRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleUser,
		Name:        "user",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "User",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}

func generateBundleGuestRole() *settings.Bundle {
	return &settings.Bundle{
		Id:          BundleUUIDRoleGuest,
		Name:        "guest",
		Type:        settings.Bundle_TYPE_ROLE,
		Extension:   "ocis-roles",
		DisplayName: "Guest",
		Resource: &settings.Resource{
			Type: settings.Resource_TYPE_SYSTEM,
		},
		Settings: []*settings.Setting{},
	}
}

package config

import (
	accounts "github.com/owncloud/ocis/v2/extensions/accounts/pkg/config/defaults"
	appProvider "github.com/owncloud/ocis/v2/extensions/app-provider/pkg/config/defaults"
	appRegistry "github.com/owncloud/ocis/v2/extensions/app-registry/pkg/config/defaults"
	audit "github.com/owncloud/ocis/v2/extensions/audit/pkg/config/defaults"
	authbasic "github.com/owncloud/ocis/v2/extensions/auth-basic/pkg/config/defaults"
	authbearer "github.com/owncloud/ocis/v2/extensions/auth-bearer/pkg/config/defaults"
	authmachine "github.com/owncloud/ocis/v2/extensions/auth-machine/pkg/config/defaults"
	frontend "github.com/owncloud/ocis/v2/extensions/frontend/pkg/config/defaults"
	gateway "github.com/owncloud/ocis/v2/extensions/gateway/pkg/config/defaults"
	graphExplorer "github.com/owncloud/ocis/v2/extensions/graph-explorer/pkg/config/defaults"
	graph "github.com/owncloud/ocis/v2/extensions/graph/pkg/config/defaults"
	groups "github.com/owncloud/ocis/v2/extensions/groups/pkg/config/defaults"
	idm "github.com/owncloud/ocis/v2/extensions/idm/pkg/config/defaults"
	idp "github.com/owncloud/ocis/v2/extensions/idp/pkg/config/defaults"
	nats "github.com/owncloud/ocis/v2/extensions/nats/pkg/config/defaults"
	notifications "github.com/owncloud/ocis/v2/extensions/notifications/pkg/config/defaults"
	ocdav "github.com/owncloud/ocis/v2/extensions/ocdav/pkg/config/defaults"
	ocs "github.com/owncloud/ocis/v2/extensions/ocs/pkg/config/defaults"
	proxy "github.com/owncloud/ocis/v2/extensions/proxy/pkg/config/defaults"
	search "github.com/owncloud/ocis/v2/extensions/search/pkg/config/defaults"
	settings "github.com/owncloud/ocis/v2/extensions/settings/pkg/config/defaults"
	sharing "github.com/owncloud/ocis/v2/extensions/sharing/pkg/config/defaults"
	storagepublic "github.com/owncloud/ocis/v2/extensions/storage-publiclink/pkg/config/defaults"
	storageshares "github.com/owncloud/ocis/v2/extensions/storage-shares/pkg/config/defaults"
	storageSystem "github.com/owncloud/ocis/v2/extensions/storage-system/pkg/config/defaults"
	storageusers "github.com/owncloud/ocis/v2/extensions/storage-users/pkg/config/defaults"
	store "github.com/owncloud/ocis/v2/extensions/store/pkg/config/defaults"
	thumbnails "github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/config/defaults"
	users "github.com/owncloud/ocis/v2/extensions/users/pkg/config/defaults"
	web "github.com/owncloud/ocis/v2/extensions/web/pkg/config/defaults"
	webdav "github.com/owncloud/ocis/v2/extensions/webdav/pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Runtime: Runtime{
			Port: "9250",
			Host: "localhost",
		},

		Accounts:          accounts.DefaultConfig(),
		AppProvider:       appProvider.DefaultConfig(),
		AppRegistry:       appRegistry.DefaultConfig(),
		Audit:             audit.DefaultConfig(),
		AuthBasic:         authbasic.DefaultConfig(),
		AuthBearer:        authbearer.DefaultConfig(),
		AuthMachine:       authmachine.DefaultConfig(),
		Frontend:          frontend.DefaultConfig(),
		Gateway:           gateway.DefaultConfig(),
		Graph:             graph.DefaultConfig(),
		GraphExplorer:     graphExplorer.DefaultConfig(),
		Groups:            groups.DefaultConfig(),
		IDM:               idm.DefaultConfig(),
		IDP:               idp.DefaultConfig(),
		Nats:              nats.DefaultConfig(),
		Notifications:     notifications.DefaultConfig(),
		OCDav:             ocdav.DefaultConfig(),
		OCS:               ocs.DefaultConfig(),
		Proxy:             proxy.DefaultConfig(),
		Search:            search.FullDefaultConfig(),
		Settings:          settings.DefaultConfig(),
		Sharing:           sharing.DefaultConfig(),
		StoragePublicLink: storagepublic.DefaultConfig(),
		StorageShares:     storageshares.DefaultConfig(),
		StorageSystem:     storageSystem.DefaultConfig(),
		StorageUsers:      storageusers.DefaultConfig(),
		Store:             store.DefaultConfig(),
		Thumbnails:        thumbnails.DefaultConfig(),
		Users:             users.DefaultConfig(),
		Web:               web.DefaultConfig(),
		WebDAV:            webdav.DefaultConfig(),
	}
}

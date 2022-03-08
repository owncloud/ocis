package config

import (
	accounts "github.com/owncloud/ocis/accounts/pkg/config/defaults"
	audit "github.com/owncloud/ocis/audit/pkg/config"
	glauth "github.com/owncloud/ocis/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/graph/pkg/config"
	idm "github.com/owncloud/ocis/idm/pkg/config"
	idp "github.com/owncloud/ocis/idp/pkg/config/defaults"
	nats "github.com/owncloud/ocis/nats/pkg/config"
	notifications "github.com/owncloud/ocis/notifications/pkg/config"
	ocs "github.com/owncloud/ocis/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/proxy/pkg/config"
	settings "github.com/owncloud/ocis/settings/pkg/config"
	storage "github.com/owncloud/ocis/storage/pkg/config"
	store "github.com/owncloud/ocis/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/thumbnails/pkg/config"
	web "github.com/owncloud/ocis/web/pkg/config"
	webdav "github.com/owncloud/ocis/webdav/pkg/config"
)

func DefaultConfig() *Config {
	return &Config{
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Runtime: Runtime{
			Port: "9250",
			Host: "localhost",
		},
		Audit:         audit.DefaultConfig(),
		Accounts:      accounts.DefaultConfig(),
		GLAuth:        glauth.DefaultConfig(),
		Graph:         graph.DefaultConfig(),
		IDP:           idp.DefaultConfig(),
		IDM:           idm.DefaultConfig(),
		Nats:          nats.DefaultConfig(),
		Notifications: notifications.DefaultConfig(),
		Proxy:         proxy.DefaultConfig(),
		GraphExplorer: graphExplorer.DefaultConfig(),
		OCS:           ocs.DefaultConfig(),
		Settings:      settings.DefaultConfig(),
		Web:           web.DefaultConfig(),
		Store:         store.DefaultConfig(),
		Thumbnails:    thumbnails.DefaultConfig(),
		WebDAV:        webdav.DefaultConfig(),
		Storage:       storage.DefaultConfig(),
	}
}

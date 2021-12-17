package config

import (
	accounts "github.com/owncloud/ocis/accounts/pkg/config"
	glauth "github.com/owncloud/ocis/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/graph/pkg/config"
	idp "github.com/owncloud/ocis/idp/pkg/config"
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
		Accounts:      accounts.DefaultConfig(),
		GLAuth:        glauth.DefaultConfig(),
		Graph:         graph.DefaultConfig(),
		IDP:           idp.DefaultConfig(),
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

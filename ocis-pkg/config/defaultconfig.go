package config

import (
	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/config/defaults"
	appprovider "github.com/owncloud/ocis/extensions/appprovider/pkg/config/defaults"
	audit "github.com/owncloud/ocis/extensions/audit/pkg/config/defaults"
	authbasic "github.com/owncloud/ocis/extensions/auth-basic/pkg/config/defaults"
	authbearer "github.com/owncloud/ocis/extensions/auth-bearer/pkg/config/defaults"
	authmachine "github.com/owncloud/ocis/extensions/auth-machine/pkg/config/defaults"
	frontend "github.com/owncloud/ocis/extensions/frontend/pkg/config/defaults"
	glauth "github.com/owncloud/ocis/extensions/glauth/pkg/config/defaults"
	graphExplorer "github.com/owncloud/ocis/extensions/graph-explorer/pkg/config/defaults"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/config/defaults"
	group "github.com/owncloud/ocis/extensions/group/pkg/config/defaults"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/config/defaults"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/config/defaults"
	nats "github.com/owncloud/ocis/extensions/nats/pkg/config/defaults"
	notifications "github.com/owncloud/ocis/extensions/notifications/pkg/config/defaults"
	ocdav "github.com/owncloud/ocis/extensions/ocdav/pkg/config/defaults"
	ocs "github.com/owncloud/ocis/extensions/ocs/pkg/config/defaults"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config/defaults"
	settings "github.com/owncloud/ocis/extensions/settings/pkg/config/defaults"
	sharing "github.com/owncloud/ocis/extensions/sharing/pkg/config/defaults"
	storagemetadata "github.com/owncloud/ocis/extensions/storage-metadata/pkg/config/defaults"
	storagepublic "github.com/owncloud/ocis/extensions/storage-publiclink/pkg/config/defaults"
	storageshares "github.com/owncloud/ocis/extensions/storage-shares/pkg/config/defaults"
	storageusers "github.com/owncloud/ocis/extensions/storage-users/pkg/config/defaults"
	storage "github.com/owncloud/ocis/extensions/storage/pkg/config/defaults"
	store "github.com/owncloud/ocis/extensions/store/pkg/config/defaults"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/config/defaults"
	user "github.com/owncloud/ocis/extensions/user/pkg/config/defaults"
	web "github.com/owncloud/ocis/extensions/web/pkg/config/defaults"
	webdav "github.com/owncloud/ocis/extensions/webdav/pkg/config/defaults"
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
		Audit:             audit.DefaultConfig(),
		Accounts:          accounts.DefaultConfig(),
		GLAuth:            glauth.DefaultConfig(),
		Graph:             graph.DefaultConfig(),
		IDP:               idp.DefaultConfig(),
		IDM:               idm.DefaultConfig(),
		Nats:              nats.DefaultConfig(),
		Notifications:     notifications.DefaultConfig(),
		Proxy:             proxy.DefaultConfig(),
		GraphExplorer:     graphExplorer.DefaultConfig(),
		OCS:               ocs.DefaultConfig(),
		Settings:          settings.DefaultConfig(),
		Web:               web.DefaultConfig(),
		Store:             store.DefaultConfig(),
		Thumbnails:        thumbnails.DefaultConfig(),
		WebDAV:            webdav.DefaultConfig(),
		Storage:           storage.DefaultConfig(),
		AuthBasic:         authbasic.FullDefaultConfig(),
		AuthBearer:        authbearer.FullDefaultConfig(),
		AuthMachine:       authmachine.FullDefaultConfig(),
		User:              user.FullDefaultConfig(),
		Group:             group.FullDefaultConfig(),
		Sharing:           sharing.FullDefaultConfig(),
		StorageMetadata:   storagemetadata.FullDefaultConfig(),
		StoragePublicLink: storagepublic.FullDefaultConfig(),
		StorageUsers:      storageusers.FullDefaultConfig(),
		StorageShares:     storageshares.FullDefaultConfig(),
		AppProvider:       appprovider.FullDefaultConfig(),
		Frontend:          frontend.FullDefaultConfig(),
		OCDav:             ocdav.FullDefaultConfig(),
	}
}

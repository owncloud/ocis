package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	activitylog "github.com/owncloud/ocis/v2/services/activitylog/pkg/config/defaults"
	antivirus "github.com/owncloud/ocis/v2/services/antivirus/pkg/config/defaults"
	appProvider "github.com/owncloud/ocis/v2/services/app-provider/pkg/config/defaults"
	appRegistry "github.com/owncloud/ocis/v2/services/app-registry/pkg/config/defaults"
	audit "github.com/owncloud/ocis/v2/services/audit/pkg/config/defaults"
	authapp "github.com/owncloud/ocis/v2/services/auth-app/pkg/config/defaults"
	authbasic "github.com/owncloud/ocis/v2/services/auth-basic/pkg/config/defaults"
	authbearer "github.com/owncloud/ocis/v2/services/auth-bearer/pkg/config/defaults"
	authmachine "github.com/owncloud/ocis/v2/services/auth-machine/pkg/config/defaults"
	authservice "github.com/owncloud/ocis/v2/services/auth-service/pkg/config/defaults"
	clientlog "github.com/owncloud/ocis/v2/services/clientlog/pkg/config/defaults"
	collaboration "github.com/owncloud/ocis/v2/services/collaboration/pkg/config/defaults"
	eventhistory "github.com/owncloud/ocis/v2/services/eventhistory/pkg/config/defaults"
	frontend "github.com/owncloud/ocis/v2/services/frontend/pkg/config/defaults"
	gateway "github.com/owncloud/ocis/v2/services/gateway/pkg/config/defaults"
	graph "github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	groups "github.com/owncloud/ocis/v2/services/groups/pkg/config/defaults"
	idm "github.com/owncloud/ocis/v2/services/idm/pkg/config/defaults"
	idp "github.com/owncloud/ocis/v2/services/idp/pkg/config/defaults"
	invitations "github.com/owncloud/ocis/v2/services/invitations/pkg/config/defaults"
	nats "github.com/owncloud/ocis/v2/services/nats/pkg/config/defaults"
	notifications "github.com/owncloud/ocis/v2/services/notifications/pkg/config/defaults"
	ocdav "github.com/owncloud/ocis/v2/services/ocdav/pkg/config/defaults"
	ocm "github.com/owncloud/ocis/v2/services/ocm/pkg/config/defaults"
	ocs "github.com/owncloud/ocis/v2/services/ocs/pkg/config/defaults"
	policies "github.com/owncloud/ocis/v2/services/policies/pkg/config/defaults"
	postprocessing "github.com/owncloud/ocis/v2/services/postprocessing/pkg/config/defaults"
	proxy "github.com/owncloud/ocis/v2/services/proxy/pkg/config/defaults"
	search "github.com/owncloud/ocis/v2/services/search/pkg/config/defaults"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/config/defaults"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/config/defaults"
	sse "github.com/owncloud/ocis/v2/services/sse/pkg/config/defaults"
	storagepublic "github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/config/defaults"
	storageshares "github.com/owncloud/ocis/v2/services/storage-shares/pkg/config/defaults"
	storageSystem "github.com/owncloud/ocis/v2/services/storage-system/pkg/config/defaults"
	storageusers "github.com/owncloud/ocis/v2/services/storage-users/pkg/config/defaults"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/config/defaults"
	userlog "github.com/owncloud/ocis/v2/services/userlog/pkg/config/defaults"
	users "github.com/owncloud/ocis/v2/services/users/pkg/config/defaults"
	web "github.com/owncloud/ocis/v2/services/web/pkg/config/defaults"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/config/defaults"
	webfinger "github.com/owncloud/ocis/v2/services/webfinger/pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		OcisURL: "https://localhost:9200",
		Runtime: Runtime{
			Port: "9250",
			Host: "localhost",
		},
		Reva: &shared.Reva{
			Address: "com.owncloud.api.gateway",
		},

		Activitylog:       activitylog.DefaultConfig(),
		Antivirus:         antivirus.DefaultConfig(),
		AppProvider:       appProvider.DefaultConfig(),
		AppRegistry:       appRegistry.DefaultConfig(),
		Audit:             audit.DefaultConfig(),
		AuthApp:           authapp.DefaultConfig(),
		AuthBasic:         authbasic.DefaultConfig(),
		AuthBearer:        authbearer.DefaultConfig(),
		AuthMachine:       authmachine.DefaultConfig(),
		AuthService:       authservice.DefaultConfig(),
		Clientlog:         clientlog.DefaultConfig(),
		Collaboration:     collaboration.DefaultConfig(),
		EventHistory:      eventhistory.DefaultConfig(),
		Frontend:          frontend.DefaultConfig(),
		Gateway:           gateway.DefaultConfig(),
		Graph:             graph.DefaultConfig(),
		Groups:            groups.DefaultConfig(),
		IDM:               idm.DefaultConfig(),
		IDP:               idp.DefaultConfig(),
		Invitations:       invitations.DefaultConfig(),
		Nats:              nats.DefaultConfig(),
		Notifications:     notifications.DefaultConfig(),
		OCDav:             ocdav.DefaultConfig(),
		OCM:               ocm.DefaultConfig(),
		OCS:               ocs.DefaultConfig(),
		Postprocessing:    postprocessing.DefaultConfig(),
		Policies:          policies.DefaultConfig(),
		Proxy:             proxy.DefaultConfig(),
		Search:            search.DefaultConfig(),
		Settings:          settings.DefaultConfig(),
		Sharing:           sharing.DefaultConfig(),
		SSE:               sse.DefaultConfig(),
		StoragePublicLink: storagepublic.DefaultConfig(),
		StorageShares:     storageshares.DefaultConfig(),
		StorageSystem:     storageSystem.DefaultConfig(),
		StorageUsers:      storageusers.DefaultConfig(),
		Thumbnails:        thumbnails.DefaultConfig(),
		Userlog:           userlog.DefaultConfig(),
		Users:             users.DefaultConfig(),
		Web:               web.DefaultConfig(),
		WebDAV:            webdav.DefaultConfig(),
		Webfinger:         webfinger.DefaultConfig(),
	}
}

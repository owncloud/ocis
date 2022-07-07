package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"

	appProvider "github.com/owncloud/ocis/v2/services/app-provider/pkg/config"
	appRegistry "github.com/owncloud/ocis/v2/services/app-registry/pkg/config"
	audit "github.com/owncloud/ocis/v2/services/audit/pkg/config"
	authbasic "github.com/owncloud/ocis/v2/services/auth-basic/pkg/config"
	authbearer "github.com/owncloud/ocis/v2/services/auth-bearer/pkg/config"
	authmachine "github.com/owncloud/ocis/v2/services/auth-machine/pkg/config"
	frontend "github.com/owncloud/ocis/v2/services/frontend/pkg/config"
	gateway "github.com/owncloud/ocis/v2/services/gateway/pkg/config"
	graphExplorer "github.com/owncloud/ocis/v2/services/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/v2/services/graph/pkg/config"
	groups "github.com/owncloud/ocis/v2/services/groups/pkg/config"
	idm "github.com/owncloud/ocis/v2/services/idm/pkg/config"
	idp "github.com/owncloud/ocis/v2/services/idp/pkg/config"
	nats "github.com/owncloud/ocis/v2/services/nats/pkg/config"
	notifications "github.com/owncloud/ocis/v2/services/notifications/pkg/config"
	ocdav "github.com/owncloud/ocis/v2/services/ocdav/pkg/config"
	ocs "github.com/owncloud/ocis/v2/services/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	search "github.com/owncloud/ocis/v2/services/search/pkg/config"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/config"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/config"
	storagepublic "github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/config"
	storageshares "github.com/owncloud/ocis/v2/services/storage-shares/pkg/config"
	storagesystem "github.com/owncloud/ocis/v2/services/storage-system/pkg/config"
	storageusers "github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	store "github.com/owncloud/ocis/v2/services/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
	users "github.com/owncloud/ocis/v2/services/users/pkg/config"
	web "github.com/owncloud/ocis/v2/services/web/pkg/config"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/config"
)

const (
	// SUPERVISED sets the runtime mode as supervised threads.
	SUPERVISED = iota

	// UNSUPERVISED sets the runtime mode as a single thread.
	UNSUPERVISED
)

type Mode int

// Runtime configures the oCIS runtime when running in supervised mode.
type Runtime struct {
	Port       string `yaml:"port" env:"OCIS_RUNTIME_PORT"`
	Host       string `yaml:"host" env:"OCIS_RUNTIME_HOST"`
	Extensions string `yaml:"services" env:"OCIS_RUN_EXTENSIONS;OCIS_RUN_SERVICES"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"shared"`

	Tracing *shared.Tracing `yaml:"tracing"`
	Log     *shared.Log     `yaml:"log"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `yaml:"ocis_url" desc:"URL, where oCIS is reachable for users."`

	Registry          string               `yaml:"registry"`
	TokenManager      *shared.TokenManager `yaml:"token_manager"`
	MachineAuthAPIKey string               `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`
	TransferSecret    string               `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET"`
	SystemUserID      string               `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID" desc:"ID of the oCIS storage-system system user. Admins need to set the ID for the storage-system system user in this config option which is then used to reference the user. Any reasonable long string is possible, preferably this would be an UUIDv4 format."`
	SystemUserAPIKey  string               `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY" desc:"API key for the storage-system system user."`
	AdminUserID       string               `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID" desc:"ID of a user, that should receive admin privileges."`
	Runtime           Runtime              `yaml:"runtime"`

	AppProvider       *appProvider.Config   `yaml:"app_provider"`
	AppRegistry       *appRegistry.Config   `yaml:"app_registry"`
	Audit             *audit.Config         `yaml:"audit"`
	AuthBasic         *authbasic.Config     `yaml:"auth_basic"`
	AuthBearer        *authbearer.Config    `yaml:"auth_bearer"`
	AuthMachine       *authmachine.Config   `yaml:"auth_machine"`
	Frontend          *frontend.Config      `yaml:"frontend"`
	Gateway           *gateway.Config       `yaml:"gateway"`
	Graph             *graph.Config         `yaml:"graph"`
	GraphExplorer     *graphExplorer.Config `yaml:"graph_explorer"`
	Groups            *groups.Config        `yaml:"groups"`
	IDM               *idm.Config           `yaml:"idm"`
	IDP               *idp.Config           `yaml:"idp"`
	Nats              *nats.Config          `yaml:"nats"`
	Notifications     *notifications.Config `yaml:"notifications"`
	OCDav             *ocdav.Config         `yaml:"ocdav"`
	OCS               *ocs.Config           `yaml:"ocs"`
	Proxy             *proxy.Config         `yaml:"proxy"`
	Settings          *settings.Config      `yaml:"settings"`
	Sharing           *sharing.Config       `yaml:"sharing"`
	StorageSystem     *storagesystem.Config `yaml:"storage_system"`
	StoragePublicLink *storagepublic.Config `yaml:"storage_public"`
	StorageShares     *storageshares.Config `yaml:"storage_shares"`
	StorageUsers      *storageusers.Config  `yaml:"storage_users"`
	Store             *store.Config         `yaml:"store"`
	Thumbnails        *thumbnails.Config    `yaml:"thumbnails"`
	Users             *users.Config         `yaml:"users"`
	Web               *web.Config           `yaml:"web"`
	WebDAV            *webdav.Config        `yaml:"webdav"`
	Search            *search.Config        `yaml:"search"`
}

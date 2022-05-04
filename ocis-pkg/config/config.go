package config

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"

	accounts "github.com/owncloud/ocis/v2/extensions/accounts/pkg/config"
	appProvider "github.com/owncloud/ocis/v2/extensions/app-provider/pkg/config"
	appRegistry "github.com/owncloud/ocis/v2/extensions/app-registry/pkg/config"
	audit "github.com/owncloud/ocis/v2/extensions/audit/pkg/config"
	authbasic "github.com/owncloud/ocis/v2/extensions/auth-basic/pkg/config"
	authbearer "github.com/owncloud/ocis/v2/extensions/auth-bearer/pkg/config"
	authmachine "github.com/owncloud/ocis/v2/extensions/auth-machine/pkg/config"
	frontend "github.com/owncloud/ocis/v2/extensions/frontend/pkg/config"
	gateway "github.com/owncloud/ocis/v2/extensions/gateway/pkg/config"
	glauth "github.com/owncloud/ocis/v2/extensions/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/v2/extensions/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/v2/extensions/graph/pkg/config"
	group "github.com/owncloud/ocis/v2/extensions/group/pkg/config"
	idm "github.com/owncloud/ocis/v2/extensions/idm/pkg/config"
	idp "github.com/owncloud/ocis/v2/extensions/idp/pkg/config"
	nats "github.com/owncloud/ocis/v2/extensions/nats/pkg/config"
	notifications "github.com/owncloud/ocis/v2/extensions/notifications/pkg/config"
	ocdav "github.com/owncloud/ocis/v2/extensions/ocdav/pkg/config"
	ocs "github.com/owncloud/ocis/v2/extensions/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/v2/extensions/proxy/pkg/config"
	search "github.com/owncloud/ocis/v2/extensions/search/pkg/config"
	settings "github.com/owncloud/ocis/v2/extensions/settings/pkg/config"
	sharing "github.com/owncloud/ocis/v2/extensions/sharing/pkg/config"
	storagepublic "github.com/owncloud/ocis/v2/extensions/storage-publiclink/pkg/config"
	storageshares "github.com/owncloud/ocis/v2/extensions/storage-shares/pkg/config"
	storagesystem "github.com/owncloud/ocis/v2/extensions/storage-system/pkg/config"
	storageusers "github.com/owncloud/ocis/v2/extensions/storage-users/pkg/config"
	store "github.com/owncloud/ocis/v2/extensions/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/v2/extensions/thumbnails/pkg/config"
	user "github.com/owncloud/ocis/v2/extensions/user/pkg/config"
	web "github.com/owncloud/ocis/v2/extensions/web/pkg/config"
	webdav "github.com/owncloud/ocis/v2/extensions/webdav/pkg/config"
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
	Extensions string `yaml:"extensions" env:"OCIS_RUN_EXTENSIONS"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"shared"`

	Tracing *shared.Tracing `yaml:"tracing"`
	Log     *shared.Log     `yaml:"log"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `yaml:"ocis_url"`

	Registry          string               `yaml:"registry"`
	TokenManager      *shared.TokenManager `yaml:"token_manager"`
	MachineAuthAPIKey string               `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
	TransferSecret    string               `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET"`
	SystemUserID      string               `yaml:"system_user_id" env:"OCIS_SYSTEM_USER_ID"`
	SystemUserAPIKey  string               `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY"`
	AdminUserID       string               `yaml:"admin_user_id" env:"OCIS_ADMIN_USER_ID"`
	Runtime           Runtime              `yaml:"runtime"`

	Accounts          *accounts.Config      `yaml:"accounts"`
	AppProvider       *appProvider.Config   `yaml:"app_provider"`
	AppRegistry       *appRegistry.Config   `yaml:"app_registry"`
	Audit             *audit.Config         `yaml:"audit"`
	AuthBasic         *authbasic.Config     `yaml:"auth_basic"`
	AuthBearer        *authbearer.Config    `yaml:"auth_bearer"`
	AuthMachine       *authmachine.Config   `yaml:"auth_machine"`
	Frontend          *frontend.Config      `yaml:"frontend"`
	Gateway           *gateway.Config       `yaml:"gateway"`
	GLAuth            *glauth.Config        `yaml:"glauth"`
	Graph             *graph.Config         `yaml:"graph"`
	GraphExplorer     *graphExplorer.Config `yaml:"graph_explorer"`
	Group             *group.Config         `yaml:"group"`
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
	User              *user.Config          `yaml:"user"`
	Web               *web.Config           `yaml:"web"`
	WebDAV            *webdav.Config        `yaml:"webdav"`
	Search            *search.Config        `yaml:"search"`
}

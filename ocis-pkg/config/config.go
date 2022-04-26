package config

import (
	"github.com/owncloud/ocis/ocis-pkg/shared"

	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/config"
	appprovider "github.com/owncloud/ocis/extensions/appprovider/pkg/config"
	audit "github.com/owncloud/ocis/extensions/audit/pkg/config"
	authbasic "github.com/owncloud/ocis/extensions/auth-basic/pkg/config"
	authbearer "github.com/owncloud/ocis/extensions/auth-bearer/pkg/config"
	authmachine "github.com/owncloud/ocis/extensions/auth-machine/pkg/config"
	frontend "github.com/owncloud/ocis/extensions/frontend/pkg/config"
	gateway "github.com/owncloud/ocis/extensions/gateway/pkg/config"
	glauth "github.com/owncloud/ocis/extensions/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/extensions/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/config"
	group "github.com/owncloud/ocis/extensions/group/pkg/config"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/config"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/config"
	nats "github.com/owncloud/ocis/extensions/nats/pkg/config"
	notifications "github.com/owncloud/ocis/extensions/notifications/pkg/config"
	ocdav "github.com/owncloud/ocis/extensions/ocdav/pkg/config"
	ocs "github.com/owncloud/ocis/extensions/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config"
	settings "github.com/owncloud/ocis/extensions/settings/pkg/config"
	sharing "github.com/owncloud/ocis/extensions/sharing/pkg/config"
	storagemetadata "github.com/owncloud/ocis/extensions/storage-metadata/pkg/config"
	storagepublic "github.com/owncloud/ocis/extensions/storage-publiclink/pkg/config"
	storageshares "github.com/owncloud/ocis/extensions/storage-shares/pkg/config"
	storageusers "github.com/owncloud/ocis/extensions/storage-users/pkg/config"
	store "github.com/owncloud/ocis/extensions/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	user "github.com/owncloud/ocis/extensions/user/pkg/config"
	web "github.com/owncloud/ocis/extensions/web/pkg/config"
	webdav "github.com/owncloud/ocis/extensions/webdav/pkg/config"
)

// TokenManager is the config for using the reva token manager
/*type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET"`
}*/

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
	*shared.Commons `yaml:"shared,omitempty"`

	Tracing *shared.Tracing `yaml:"tracing,omitempty"`
	Log     *shared.Log     `yaml:"log,omitempty"`

	Mode    Mode   `yaml:",omitempty"` // DEPRECATED
	File    string `yaml:",omitempty"`
	OcisURL string `yaml:"ocis_url,omitempty"`

	Registry          string               `yaml:"registry,omitempty"`
	TokenManager      *shared.TokenManager `yaml:"token_manager,omitempty"`
	MachineAuthAPIKey string               `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
	TransferSecret    string               `yaml:"transfer_secret,omitempty"`
	Runtime           Runtime              `yaml:"runtime,omitempty"`

	Audit             *audit.Config           `yaml:"audit,omitempty"`
	Accounts          *accounts.Config        `yaml:"accounts,omitempty"`
	GLAuth            *glauth.Config          `yaml:"glauth,omitempty"`
	Graph             *graph.Config           `yaml:"graph,omitempty"`
	GraphExplorer     *graphExplorer.Config   `yaml:"graph_explorer,omitempty"`
	IDP               *idp.Config             `yaml:"idp,omitempty"`
	IDM               *idm.Config             `yaml:"idm,omitempty"`
	Nats              *nats.Config            `yaml:"nats,omitempty"`
	Notifications     *notifications.Config   `yaml:"notifications,omitempty"`
	OCS               *ocs.Config             `yaml:"ocs,omitempty"`
	Web               *web.Config             `yaml:"web,omitempty"`
	Proxy             *proxy.Config           `yaml:"proxy,omitempty"`
	Settings          *settings.Config        `yaml:"settings,omitempty"`
	Gateway           *gateway.Config         `yaml:"gateway,omitempty"`
	Frontend          *frontend.Config        `yaml:"frontend,omitempty"`
	AuthBasic         *authbasic.Config       `yaml:"auth_basic,omitempty"`
	AuthBearer        *authbearer.Config      `yaml:"auth_bearer,omitempty"`
	AuthMachine       *authmachine.Config     `yaml:"auth_machine,omitempty"`
	User              *user.Config            `yaml:"user,omitempty"`
	Group             *group.Config           `yaml:"group,omitempty"`
	AppProvider       *appprovider.Config     `yaml:"app_provider,omitempty"`
	Sharing           *sharing.Config         `yaml:"sharing,omitempty"`
	StorageMetadata   *storagemetadata.Config `yaml:"storage_metadata,omitempty"`
	StoragePublicLink *storagepublic.Config   `yaml:"storage_public,omitempty"`
	StorageUsers      *storageusers.Config    `yaml:"storage_users,omitempty"`
	StorageShares     *storageshares.Config   `yaml:"storage_shares,omitempty"`
	OCDav             *ocdav.Config           `yaml:"ocdav,omitempty"`
	Store             *store.Config           `yaml:"store,omitempty"`
	Thumbnails        *thumbnails.Config      `yaml:"thumbnails,omitempty"`
	WebDAV            *webdav.Config          `yaml:"webdav,omitempty"`
}

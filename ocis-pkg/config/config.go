package config

import (
	"github.com/owncloud/ocis/ocis-pkg/shared"

	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/config"
	audit "github.com/owncloud/ocis/extensions/audit/pkg/config"
	glauth "github.com/owncloud/ocis/extensions/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/extensions/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/config"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/config"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/config"
	nats "github.com/owncloud/ocis/extensions/nats/pkg/config"
	notifications "github.com/owncloud/ocis/extensions/notifications/pkg/config"
	ocs "github.com/owncloud/ocis/extensions/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config"
	search "github.com/owncloud/ocis/extensions/search/pkg/config"
	settings "github.com/owncloud/ocis/extensions/settings/pkg/config"
	storage "github.com/owncloud/ocis/extensions/storage/pkg/config"
	store "github.com/owncloud/ocis/extensions/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	web "github.com/owncloud/ocis/extensions/web/pkg/config"
	webdav "github.com/owncloud/ocis/extensions/webdav/pkg/config"
)

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET"`
}

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

	Tracing shared.Tracing `yaml:"tracing"`
	Log     *shared.Log    `yaml:"log"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `yaml:"ocis_url"`

	Registry     string       `yaml:"registry"`
	TokenManager TokenManager `yaml:"token_manager"`
	Runtime      Runtime      `yaml:"runtime"`

	Audit         *audit.Config         `yaml:"audit"`
	Accounts      *accounts.Config      `yaml:"accounts"`
	GLAuth        *glauth.Config        `yaml:"glauth"`
	Graph         *graph.Config         `yaml:"graph"`
	GraphExplorer *graphExplorer.Config `yaml:"graph_explorer"`
	IDP           *idp.Config           `yaml:"idp"`
	IDM           *idm.Config           `yaml:"idm"`
	Nats          *nats.Config          `yaml:"nats"`
	Notifications *notifications.Config `yaml:"notifications"`
	OCS           *ocs.Config           `yaml:"ocs"`
	Web           *web.Config           `yaml:"web"`
	Proxy         *proxy.Config         `yaml:"proxy"`
	Search        *search.Config        `yaml:"search"`
	Settings      *settings.Config      `yaml:"settings"`
	Storage       *storage.Config       `yaml:"storage"`
	Store         *store.Config         `yaml:"store"`
	Thumbnails    *thumbnails.Config    `yaml:"thumbnails"`
	WebDAV        *webdav.Config        `yaml:"webdav"`
}

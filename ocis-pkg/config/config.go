package config

import (
	"github.com/owncloud/ocis/ocis-pkg/shared"

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

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET"`
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
	Port       string `ocisConfig:"port" env:"OCIS_RUNTIME_PORT"`
	Host       string `ocisConfig:"host" env:"OCIS_RUNTIME_HOST"`
	Extensions string `ocisConfig:"extensions" env:"OCIS_RUN_EXTENSIONS"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"shared"`

	Tracing shared.Tracing `ocisConfig:"tracing"`
	Log     *shared.Log    `ocisConfig:"log"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `ocisConfig:"ocis_url"`

	Registry     string       `ocisConfig:"registry"`
	TokenManager TokenManager `ocisConfig:"token_manager"`
	Runtime      Runtime      `ocisConfig:"runtime"`

	Accounts      *accounts.Config      `ocisConfig:"accounts"`
	GLAuth        *glauth.Config        `ocisConfig:"glauth"`
	Graph         *graph.Config         `ocisConfig:"graph"`
	GraphExplorer *graphExplorer.Config `ocisConfig:"graph_explorer"`
	IDP           *idp.Config           `ocisConfig:"idp"`
	OCS           *ocs.Config           `ocisConfig:"ocs"`
	Web           *web.Config           `ocisConfig:"web"`
	Proxy         *proxy.Config         `ocisConfig:"proxy"`
	Settings      *settings.Config      `ocisConfig:"settings"`
	Storage       *storage.Config       `ocisConfig:"storage"`
	Store         *store.Config         `ocisConfig:"store"`
	Thumbnails    *thumbnails.Config    `ocisConfig:"thumbnails"`
	WebDAV        *webdav.Config        `ocisConfig:"webdav"`
}

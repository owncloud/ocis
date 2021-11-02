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

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string `mapstructure:"addr"`
	Root string `mapstructure:"root"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string `mapstructure:"addr"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mapstructure:"jwt_secret"`
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
	Port       string `mapstructure:"port"`
	Host       string `mapstructure:"host"`
	Extensions string `mapstructure:"extensions"`
}

// Config combines all available configuration parts.
type Config struct {
	// Mode is mostly used whenever we need to run an extension. The technical debt this introduces is in regard of
	// what it does. Supervised (0) loads configuration from a unified config file because of known limitations of Viper; whereas
	// Unsupervised (1) MUST parse config from all known sources.
	Mode Mode
	File string

	Registry     string       `mapstructure:"registry"`
	Log          Log          `mapstructure:"log"`
	Debug        Debug        `mapstructure:"debug"`
	HTTP         HTTP         `mapstructure:"http"`
	GRPC         GRPC         `mapstructure:"grpc"`
	Tracing      Tracing      `mapstructure:"tracing"`
	TokenManager TokenManager `mapstructure:"token_manager"`
	Runtime      Runtime      `mapstructure:"runtime"`

	Accounts      *accounts.Config      `mapstructure:"accounts"`
	GLAuth        *glauth.Config        `mapstructure:"glauth"`
	Graph         *graph.Config         `mapstructure:"graph"`
	GraphExplorer *graphExplorer.Config `mapstructure:"graph_explorer"`
	IDP           *idp.Config           `mapstructure:"idp"`
	OCS           *ocs.Config           `mapstructure:"ocs"`
	Web           *web.Config           `mapstructure:"web"`
	Proxy         *proxy.Config         `mapstructure:"proxy"`
	Settings      *settings.Config      `mapstructure:"settings"`
	Storage       *storage.Config       `mapstructure:"storage"`
	Store         *store.Config         `mapstructure:"store"`
	Thumbnails    *thumbnails.Config    `mapstructure:"thumbnails"`
	WebDAV        *webdav.Config        `mapstructure:"webdav"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{
		Accounts:      accounts.DefaultConfig(),
		GLAuth:        glauth.New(),
		Graph:         graph.New(),
		GraphExplorer: graphExplorer.New(),
		IDP:           idp.New(),
		OCS:           ocs.New(),
		Web:           web.New(),
		Proxy:         proxy.DefaultConfig(),
		Settings:      settings.New(),
		Storage:       storage.New(),
		Store:         store.New(),
		Thumbnails:    thumbnails.New(),
		WebDAV:        webdav.New(),
	}
}

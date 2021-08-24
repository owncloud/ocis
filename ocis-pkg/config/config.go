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
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string
	Root string
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
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
	Port       string
	Host       string
	Extensions string
}

// Config combines all available configuration parts.
type Config struct {
	// Mode is mostly used whenever we need to run an extension. The technical debt this introduces is in regards of
	// what it does. Supervised (0) loads configuration from a unified config file because of known limitations of Viper; whereas
	// Unsupervised (1) MUST parse config from all known sources.
	Mode Mode
	File string

	Registry     string
	Log          Log
	Debug        Debug
	HTTP         HTTP
	GRPC         GRPC
	Tracing      Tracing
	TokenManager TokenManager
	Runtime      Runtime

	Accounts      *accounts.Config
	GLAuth        *glauth.Config
	Graph         *graph.Config
	GraphExplorer *graphExplorer.Config
	IDP           *idp.Config
	OCS           *ocs.Config
	Web           *web.Config
	Proxy         *proxy.Config
	Settings      *settings.Config
	Storage       *storage.Config
	Store         *store.Config
	Thumbnails    *thumbnails.Config
	WebDAV        *webdav.Config
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{
		Accounts:      accounts.New(),
		GLAuth:        glauth.New(),
		Graph:         graph.New(),
		GraphExplorer: graphExplorer.New(),
		IDP:           idp.New(),
		OCS:           ocs.New(),
		Web:           web.New(),
		Proxy:         proxy.New(),
		Settings:      settings.New(),
		Storage:       storage.New(),
		Store:         store.New(),
		Thumbnails:    thumbnails.New(),
		WebDAV:        webdav.New(),
	}
}

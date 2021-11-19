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

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string `ocisConfig:"addr"`
	Root string `ocisConfig:"root"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string `ocisConfig:"addr"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret"`
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
	Port       string `ocisConfig:"port"`
	Host       string `ocisConfig:"host"`
	Extensions string `ocisConfig:"extensions"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"shared"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `ocisConfig:"ocis_url"`

	Registry     string       `ocisConfig:"registry"`
	Log          shared.Log   `ocisConfig:"log"`
	Debug        Debug        `ocisConfig:"debug"`
	HTTP         HTTP         `ocisConfig:"http"`
	GRPC         GRPC         `ocisConfig:"grpc"`
	Tracing      Tracing      `ocisConfig:"tracing"`
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

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{
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

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9010",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr: "127.0.0.1:9000",
			Root: "/",
		},
		GRPC: GRPC{
			Addr: "127.0.0.1:9001",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "ocis",
		},
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

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}

// StructMappings binds a set of environment variables to a destination on cfg. Iterating over this set and editing the
// Destination value of a binding will alter the original value, as it is a pointer to its memory address. This lets
// us propagate changes easier.
func StructMappings(cfg *Config) []shared.EnvBinding {
	return structMappings(cfg)
}

func structMappings(cfg *Config) []shared.EnvBinding {
	return []shared.EnvBinding{
		{
			EnvVars:     []string{"OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		{
			EnvVars:     []string{"OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		{
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		{
			EnvVars:     []string{"OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"OCIS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"OCIS_RUNTIME_PORT"},
			Destination: &cfg.Runtime.Port,
		},
		{
			EnvVars:     []string{"OCIS_RUNTIME_HOST"},
			Destination: &cfg.Runtime.Host,
		},
		{
			EnvVars:     []string{"OCIS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		{
			EnvVars:     []string{"OCIS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		{
			EnvVars:     []string{"OCIS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		{
			EnvVars:     []string{"OCIS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		{
			EnvVars:     []string{"OCIS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"OCIS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"OCIS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		{
			EnvVars:     []string{"OCIS_RUN_EXTENSIONS"},
			Destination: &cfg.Runtime.Extensions,
		},
	}
}

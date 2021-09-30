package flagset

import (
	"encoding/json"

	registrypb "github.com/cs3org/go-cs3apis/cs3/app/registry/v1beta1"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

type mimeTypesCfg struct {
	cfg  map[string]config.MimeTypeConfig
	dest *map[string]config.MimeTypeConfig
}

// mimeTypes implements the Generic interface

// Set sets the json decoded value of the struct
func (m *mimeTypesCfg) Set(value string) error {
	// decode the (json) value as a map
	err := json.Unmarshal([]byte(value), &m.cfg)
	if err != nil {
		return err
	}
	*m.dest = m.cfg
	return nil
}

// String return a string representation (json encoded) of the struct
func (m *mimeTypesCfg) String() string {
	// encode the map into json
	b, err := json.Marshal(m.cfg)
	if err != nil {
		return ""
	}
	return string(b)
}

func overrideDefaultMimeTypesCfg(cfg map[string]config.MimeTypeConfig, default_ *mimeTypesCfg) *mimeTypesCfg {
	if len(cfg) == 0 {
		*default_.dest = default_.cfg // set destination
		return default_
	}
	return &mimeTypesCfg{
		cfg:  cfg,
		dest: default_.dest,
	}
}

type providersCfg struct {
	cfg  map[string]registrypb.ProviderInfo
	dest *map[string]registrypb.ProviderInfo
}

// mimeTypes implements the Generic interface

// Set sets the json decoded value of the struct
func (m *providersCfg) Set(value string) error {
	// decode the (json) value as a map
	err := json.Unmarshal([]byte(value), &m.cfg)
	if err != nil {
		return err
	}
	*m.dest = m.cfg
	return nil
}

// String return a string representation (json encoded) of the struct
func (m *providersCfg) String() string {
	// encode the map into json
	b, err := json.Marshal(m.cfg)
	if err != nil {
		return ""
	}
	return string(b)
}

func overrideDefaultProviderCfg(cfg map[string]registrypb.ProviderInfo, default_ *providersCfg) *providersCfg {
	if len(cfg) == 0 {
		*default_.dest = default_.cfg // set destination
		return default_
	}
	return &providersCfg{
		cfg:  cfg,
		dest: default_.dest,
	}
}

// AppProviderWithConfig applies cfg to the root flagset
func AppRegistryWithConfig(cfg *config.Config) []cli.Flag {

	flags := []cli.Flag{

		// AppRegistry

		&cli.GenericFlag{
			Name:    "mime-types",
			Value:   overrideDefaultMimeTypesCfg(cfg.Reva.AppRegistry.MimeTypes, &mimeTypesCfg{dest: &cfg.Reva.AppRegistry.MimeTypes}),
			Usage:   "Configuration for mime types",
			EnvVars: []string{"APP_REGISTRY_MIME_TYPES"},
		},

		&cli.GenericFlag{
			Name:    "providers",
			Value:   overrideDefaultProviderCfg(cfg.Reva.AppRegistry.Providers, &providersCfg{dest: &cfg.Reva.AppRegistry.Providers}),
			Usage:   "Configuration for providers",
			EnvVars: []string{"APP_REGISTRY_PROVIDERS"},
		},
	}

	return flags
}

package config

import (
	"path"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

var (
	// decoderConfigTagname sets the tag name to be used from the config structs
	// currently we only support "yaml" because we only support config loading
	// from yaml files and the yaml parser has no simple way to set a custom tag name to use
	decoderConfigTagName = "yaml"
)

// BindSourcesToStructs assigns any config value from a config file / env variable to struct `dst`. Its only purpose
// is to solely modify `dst`, not dealing with the config structs; and do so in a thread safe manner.
func BindSourcesToStructs(extension string, dst interface{}) (*gofig.Config, error) {
	cnf := gofig.NewWithOptions(extension)
	cnf.WithOptions(func(options *gofig.Options) {
		options.DecoderConfig.TagName = decoderConfigTagName
	})
	cnf.AddDriver(gooyaml.Driver)

	cfgFile := path.Join(defaults.BaseConfigPath(), extension+".yaml")
	_ = cnf.LoadFiles([]string{cfgFile}...)

	err := cnf.BindStruct("", &dst)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

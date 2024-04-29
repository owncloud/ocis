package config

import (
	"github.com/a8m/envsubst/parse"
	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"io/fs"
	"os"
	"path"
	"strings"
)

var (
	// decoderConfigTagName sets the tag name to be used from the config structs
	// currently we only support "yaml" because we only support config loading
	// from yaml files and the yaml parser has no simple way to set a custom tag name to use
	decoderConfigTagName = "yaml"
)

// BindSourcesToStructs assigns any config value from a config file / env variable to struct `dst`. Its only purpose
// is to solely modify `dst`, not dealing with the config structs; and do so in a thread safe manner.
func BindSourcesToStructs(service string, dst interface{}) (*gofig.Config, error) {
	fileSystem := os.DirFS("/")
	filePath := strings.TrimLeft(path.Join(defaults.BaseConfigPath(), service+".yaml"), "/")
	return bindSourcesToStructs(fileSystem, filePath, service, os.Environ(), dst)
}

func bindSourcesToStructs(fileSystem fs.FS, filePath, service string, env []string, dst interface{}) (*gofig.Config, error) {
	cnf := gofig.NewWithOptions(service)
	cnf.WithOptions(func(options *gofig.Options) {
		options.ParseEnv = false
		options.DecoderConfig.TagName = decoderConfigTagName
	})
	cnf.AddDriver(gooyaml.Driver)

	yamlContent, err := fs.ReadFile(fileSystem, filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cnf, nil
		}

		return nil, err
	}
	// replace env vars ..
	yamlContent, err = substitute(yamlContent, env)
	if err != nil {
		return nil, err
	}

	_ = cnf.LoadSources("yaml", yamlContent)

	err = cnf.BindStruct("", &dst)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

func substitute(b []byte, env []string) ([]byte, error) {
	s, err := parse.New("bytes", env,
		&parse.Restrictions{}).Parse(string(b))
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

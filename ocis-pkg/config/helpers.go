package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

var (
	// supportedExtensions is determined by gookit/config.
	// we only support the official yaml file ending (http://yaml.org/faq.html) to
	// mitigate the loading order problem.
	// It would raise this question: does yaml win over yml or vice versa!?
	supportedExtensions = []string{
		"yaml",
	}
	// decoderConfigTagname sets the tag name to be used from the config structs
	// currently we only support "yaml" because we only support config loading
	// from yaml files and the yaml parser has no simple way to set a custom tag name to use
	decoderConfigTagName = "yaml"
)

// configSources returns a slice with matched expected config files.
// It uses globbing to match a config file by name, and retrieve any supported extension supported by our drivers.
// It sanitizes the output depending on the list of drivers provided.
func configSources(filename string, drivers []string) []string {
	var sources []string

	locations := []string{
		defaults.BaseConfigPath(),
	}

	for i := range locations {
		dirFS := os.DirFS(locations[i])
		pattern := filename + ".*"
		matched, _ := fs.Glob(dirFS, pattern)
		if len(matched) > 0 {
			// prepend path to results
			for j := 0; j < len(matched); j++ {
				matched[j] = filepath.Join(locations[i], matched[j])
			}
		}
		sources = append(sources, matched...)
	}

	return sanitizeExtensions(sources, drivers, func(a, b string) bool {
		return strings.HasSuffix(filepath.Base(a), b)
	})
}

// sanitizeExtensions removes elements from "set" which extensions are not in "ext".
func sanitizeExtensions(set []string, ext []string, f func(a, b string) bool) []string {
	var r []string
	for i := 0; i < len(set); i++ {
		for j := 0; j < len(ext); j++ {
			if f(filepath.Base(set[i]), ext[j]) {
				r = append(r, set[i])
			}
		}
	}
	return r
}

// BindSourcesToStructs assigns any config value from a config file / env variable to struct `dst`. Its only purpose
// is to solely modify `dst`, not dealing with the config structs; and do so in a thread safe manner.
func BindSourcesToStructs(extension string, dst interface{}) (*gofig.Config, error) {
	sources := configSources(extension, supportedExtensions)
	cnf := gofig.NewWithOptions(extension)
	cnf.WithOptions(func(options *gofig.Options) {
		options.DecoderConfig.TagName = decoderConfigTagName
	})
	cnf.AddDriver(gooyaml.Driver)
	_ = cnf.LoadFiles(sources...)

	err := cnf.BindStruct("", &dst)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

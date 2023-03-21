package config

import (
	"crypto/tls"
	"fmt"
	"path"

	gofig "github.com/gookit/config/v2"
	gooyaml "github.com/gookit/config/v2/yaml"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

var (
	// decoderConfigTagname sets the tag name to be used from the config structs
	// currently we only support "yaml" because we only support config loading
	// from yaml files and the yaml parser has no simple way to set a custom tag name to use
	decoderConfigTagName = "yaml"
)

// BindSourcesToStructs assigns any config value from a config file / env variable to struct `dst`. Its only purpose
// is to solely modify `dst`, not dealing with the config structs; and do so in a thread safe manner.
func BindSourcesToStructs(service string, dst interface{}) (*gofig.Config, error) {
	cnf := gofig.NewWithOptions(service)
	cnf.WithOptions(func(options *gofig.Options) {
		options.DecoderConfig.TagName = decoderConfigTagName
	})
	cnf.AddDriver(gooyaml.Driver)

	cfgFile := path.Join(defaults.BaseConfigPath(), service+".yaml")
	_ = cnf.LoadFiles([]string{cfgFile}...)

	err := cnf.BindStruct("", &dst)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

// BuildTLSConfig returns a tls.Config struct for the given configuration.
// When tls is enabled it will try to load the given certificate or generate a self signed certificate
func BuildTLSConfig(l log.Logger, enabled bool, certPath, keyPath, address string) (*tls.Config, error) {
	if enabled {
		var cert tls.Certificate
		var err error
		if certPath != "" {
			// Generate a self-signing cert if no certificate is present
			if err := ociscrypto.GenCert(certPath, keyPath, l); err != nil {
				return nil, err
			}
			cert, err = tls.LoadX509KeyPair(certPath, keyPath)
			if err != nil {
				return nil, fmt.Errorf("error loading server certificate and key: %w", err)
			}
		} else {
			cert, err = ociscrypto.GenTempCertForAddr(address)
			if err != nil {
				return nil, err
			}
		}
		return &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
			//MinVersion:   tls.VersionTLS12,
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		}, nil
	}
	return nil, nil

}

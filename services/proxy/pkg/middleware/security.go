package middleware

import (
	"github.com/a8m/envsubst"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

// LoadCSPConfig loads CSP header configuration from a yaml file.
func LoadCSPConfig(proxyCfg *config.Config) (*config.CSP, error) {
	yamlContent, err := loadCSPYaml(proxyCfg)
	if err != nil {
		return nil, err
	}
	// replace env vars ..
	yamlContent, err = envsubst.Bytes(yamlContent)
	if err != nil {
		return nil, err
	}

	// read yaml
	cspConfig := config.CSP{}
	err = yaml.Unmarshal(yamlContent, &cspConfig)
	if err != nil {
		return nil, err
	}

	return &cspConfig, nil
}

func loadCSPYaml(proxyCfg *config.Config) ([]byte, error) {
	if proxyCfg.CSPConfigFileLocation == "" {
		return []byte(config.DefaultCSPConfig), nil
	}
	return os.ReadFile(proxyCfg.CSPConfigFileLocation)
}

// Security is a middleware to apply security relevant http headers like CSP.
func Security(cspConfig *config.CSP) func(h http.Handler) http.Handler {
	cspBuilder := cspbuilder.Builder{
		Directives: cspConfig.Directives,
	}

	secureMiddleware := secure.New(secure.Options{
		BrowserXssFilter:        true,
		ContentSecurityPolicy:   cspBuilder.MustBuild(),
		ContentTypeNosniff:      true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		FrameDeny:               true,
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		STSSeconds:              315360000,
		STSPreload:              true,
	})
	return func(next http.Handler) http.Handler {
		return secureMiddleware.Handler(next)
	}
}

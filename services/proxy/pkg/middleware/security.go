package middleware

import (
	"net/http"
	"os"

	gofig "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"
)

// LoadCSPConfig loads CSP header configuration from a yaml file.
func LoadCSPConfig(proxyCfg *config.Config) (*config.CSP, error) {
	yamlContent, err := loadCSPYaml(proxyCfg)
	if err != nil {
		return nil, err
	}
	return loadCSPConfig(yamlContent)
}

// LoadCSPConfig loads CSP header configuration from a yaml file.
func loadCSPConfig(yamlContent []byte) (*config.CSP, error) {
	// substitute env vars and load to struct
	gofig.WithOptions(gofig.ParseEnv)
	gofig.AddDriver(yaml.Driver)

	err := gofig.LoadSources("yaml", yamlContent)
	if err != nil {
		return nil, err
	}

	// read yaml
	cspConfig := config.CSP{}
	err = gofig.BindStruct("", &cspConfig)
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
		BrowserXssFilter:             true,
		ContentSecurityPolicy:        cspBuilder.MustBuild(),
		ContentTypeNosniff:           true,
		CustomFrameOptionsValue:      "SAMEORIGIN",
		FrameDeny:                    true,
		ReferrerPolicy:               "strict-origin-when-cross-origin",
		STSSeconds:                   315360000,
		STSPreload:                   true,
		PermittedCrossDomainPolicies: "none",
		RobotTag:                     "none",
	})
	return func(next http.Handler) http.Handler {
		return secureMiddleware.Handler(next)
	}
}

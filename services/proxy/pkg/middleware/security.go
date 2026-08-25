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
	serviceName := proxyCfg.Service.Name
	yamlContent, err := loadCSPYaml(proxyCfg)
	if err != nil {
		return nil, err
	}
	return loadCSPConfig(serviceName, yamlContent)
}

// LoadCSPConfig loads CSP header configuration from a yaml file.
func loadCSPConfig(svcName string, yamlContent []byte) (*config.CSP, error) {
	// The LoadCSPConfig (public) is expected to be called once, so
	// creating a new gookit config instance each time here is fine.
	// substitute env vars and load to struct
	cfg := gofig.NewWithOptions(svcName, gofig.ParseEnv)
	cfg.AddDriver(yaml.Driver)

	err := cfg.LoadSources("yaml", yamlContent)
	if err != nil {
		return nil, err
	}

	// read yaml
	cspConfig := config.CSP{}
	err = cfg.BindStruct("", &cspConfig)
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
func Security(cfg *config.Config, cspConfig *config.CSP) func(h http.Handler) http.Handler {
	cspBuilder := cspbuilder.Builder{
		Directives: cspConfig.Directives,
	}

	secureMiddleware := secure.New(secure.Options{
		ContentSecurityPolicy:        cspBuilder.MustBuild(),
		ContentTypeNosniff:           true,
		CustomFrameOptionsValue:      "SAMEORIGIN",
		CustomBrowserXssValue:        "0",
		FrameDeny:                    true,
		ReferrerPolicy:               "no-referrer",
		STSSeconds:                   315360000,
		ForceSTSHeader:               cfg.HTTP.ForceStrictTransportSecurity,
		STSIncludeSubdomains:         true,
		STSPreload:                   true,
		PermittedCrossDomainPolicies: "none",
		RobotTag:                     "none",
	})
	return func(next http.Handler) http.Handler {
		return secureMiddleware.Handler(next)
	}
}

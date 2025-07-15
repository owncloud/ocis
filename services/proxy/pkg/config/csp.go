package config

import (
	_ "embed"
)

// CSP defines CSP header directives
type CSP struct {
	Directives map[string][]string `yaml:"directives"`
}

//go:embed csp.yaml
var DefaultCSPConfig string

package config

// UnifiedRoles contains all settings related to unified roles.
type UnifiedRoles struct {
	AvailableRoles []string `yaml:"available_roles" env:"GRAPH_AVAILABLE_ROLES" desc:"A comma separated list of roles that are available for assignment." introductionVersion:"7.0.0"`
}

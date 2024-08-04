package config

// UnifiedRoles contains all settings related to unified roles.
type UnifiedRoles struct {
	AvailableRoles []string `yaml:"available_roles" env:"UNIFIED_ROLES_AVAILABLE_ROLES" desc:"A list of roles that are available for assignment." introductionVersion:"%%NEXT%%"`
}

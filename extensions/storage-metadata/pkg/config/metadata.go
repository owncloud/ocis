package config

func MetadataDrivers(cfg *Config) map[string]interface{} {
	return map[string]interface{}{
		"ocis": map[string]interface{}{
			"root":                cfg.Drivers.OCIS.Root,
			"user_layout":         cfg.Drivers.OCIS.UserLayout,
			"treetime_accounting": false,
			"treesize_accounting": false,
			"permissionssvc":      cfg.Drivers.OCIS.PermissionsEndpoint,
		},
	}
}

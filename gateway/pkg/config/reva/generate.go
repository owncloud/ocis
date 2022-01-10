package reva

import (
	"github.com/owncloud/ocis/gateway/pkg/config"
)

func Config(cfg *config.Config) (map[string]interface{}, error) {

	storageRegistryRules, err := rules(cfg)
	if err != nil {
		return nil, err
	}

	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			//"max_cpus":             cfg.Reva.Users.MaxCPUs,
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret": cfg.TokenManager.JWTSecret,
			"gatewaysvc": cfg.Reva.Address,
			//"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Network,
			"address": cfg.GRPC.Addr,
			"services": map[string]interface{}{
				"gateway": map[string]interface{}{
					"appregistrysvc":         cfg.ServiceMap.AppRegistryAddr,
					"authregistrysvc":        cfg.ServiceMap.AuthRegistryAddr,
					"groupprovidersvc":       cfg.ServiceMap.GroupProviderAddr,
					"ocmshareprovidersvc":    cfg.ServiceMap.OCMShareProviderAddr,
					"preferencessvc":         cfg.ServiceMap.PreferenceAddr,
					"publicshareprovidersvc": cfg.ServiceMap.PublicShareProviderAddr,
					"storageregistrysvc":     cfg.ServiceMap.StorageRegistryAddr,
					"userprovidersvc":        cfg.ServiceMap.UserProviderAddr,
					"usershareprovidersvc":   cfg.ServiceMap.UserShareProviderAddr,

					//"commit_share_to_storage_grant": cfg.Reva.Gateway.CommitShareToStorageGrant,
					//"commit_share_to_storage_ref":   cfg.Reva.Gateway.CommitShareToStorageRef,
					//"share_folder":                  cfg.Reva.Gateway.ShareFolder, // ShareFolder is the location where to create shares in the recipient's storage provider.

					// other
					//"disable_home_creation_on_login": cfg.Reva.Gateway.DisableHomeCreationOnLogin,
					//"datagateway":                    cfg.Reva.DataGateway.PublicURL,
					//"transfer_shared_secret":         cfg.Reva.TransferSecret,
					//"transfer_expires":               cfg.Reva.TransferExpires,
					//"home_mapping":                   cfg.Reva.Gateway.HomeMapping,
					//"etag_cache_ttl":                 cfg.Reva.Gateway.EtagCacheTTL,
				},
				"authregistry": map[string]interface{}{
					"driver": "static",
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"rules": map[string]interface{}{
								"basic":        cfg.ServiceMap.AuthBasicAddr,
								"bearer":       cfg.ServiceMap.AuthBearerAddr,
								"machine":      cfg.ServiceMap.AuthMachineAddr,
								"publicshares": cfg.ServiceMap.AuthPublicSharesAddr,
							},
						},
					},
				},
				//"appregistry": map[string]interface{}{
				//	"driver": "static",
				//	"drivers": map[string]interface{}{
				//		"static": map[string]interface{}{
				//			"mime_types": mimetypes(cfg, logger),
				//		},
				//	},
				//},
				"storageregistry": map[string]interface{}{
					"driver": cfg.StorageRegistry.Driver,
					"drivers": map[string]interface{}{
						"static": map[string]interface{}{
							"home_provider": cfg.StorageRegistry.HomeProvider,
							"rules":         storageRegistryRules,
						},
					},
				},
			},
		},
	}
	return rcfg, nil
}

func rules(cfg *config.Config) (map[string]map[string]interface{}, error) {

	// if a list of rules is given it overrides the generated rules from below
	//if len(cfg.Reva.StorageRegistry.Rules) > 0 {
	//	rules := map[string]map[string]interface{}{}
	//	for i := range cfg.Reva.StorageRegistry.Rules {
	//		parts := strings.SplitN(cfg.Reva.StorageRegistry.Rules[i], "=", 2)
	//		rules[parts[0]] = map[string]interface{}{"address": parts[1]}
	//	}
	//	return rules
	//}

	// check if the rules have to be read from a json file
	// if cfg.Reva.StorageRegistry.JSON != "" {
	// 	data, err := ioutil.ReadFile(cfg.Reva.StorageRegistry.JSON)
	// 	if err != nil {
	// 		logger.Error().Err(err).Msg("Failed to read storage registry rules from JSON file: " + cfg.Reva.StorageRegistry.JSON)
	// 		return nil
	// 	}
	// 	var rules map[string]map[string]interface{}
	// 	if err = json.Unmarshal(data, &rules); err != nil {
	// 		logger.Error().Err(err).Msg("Failed to unmarshal storage registry rules")
	// 		return nil
	// 	}
	// 	return rules
	// }

	// generate rules based on default config
	ret := map[string]map[string]interface{}{
		cfg.StorageRegistry.Storages.StorageHome.MountPath:        {"address": cfg.ServiceMap.StorageHomeAddr},
		cfg.StorageRegistry.Storages.StorageHome.AlternativeID:    {"address": cfg.ServiceMap.StorageHomeAddr},
		cfg.StorageRegistry.Storages.StorageUsers.MountPath:       {"address": cfg.ServiceMap.StorageUsersAddr},
		cfg.StorageRegistry.Storages.StorageUsers.MountID + ".*":  {"address": cfg.ServiceMap.StorageUsersAddr},
		cfg.StorageRegistry.Storages.StoragePublicShare.MountPath: {"address": cfg.ServiceMap.StoragePublicShare},
		cfg.StorageRegistry.Storages.StoragePublicShare.MountID:   {"address": cfg.ServiceMap.StoragePublicShare},
		// medatada storage not part of the global namespace
	}

	return ret, nil
}

// func mimetypes(cfg *config.Config, logger log.Logger) []map[string]interface{} {

// 	type mimeTypeConfig struct {
// 		MimeType      string `json:"mime_type" mapstructure:"mime_type"`
// 		Extension     string `json:"extension" mapstructure:"extension"`
// 		Name          string `json:"name" mapstructure:"name"`
// 		Description   string `json:"description" mapstructure:"description"`
// 		Icon          string `json:"icon" mapstructure:"icon"`
// 		DefaultApp    string `json:"default_app" mapstructure:"default_app"`
// 		AllowCreation bool   `json:"allow_creation" mapstructure:"allow_creation"`
// 	}
// 	var mimetypes []mimeTypeConfig
// 	var m []map[string]interface{}

// 	// load default app mimetypes from a json file
// 	if cfg.Reva.AppRegistry.MimetypesJSON != "" {
// 		data, err := ioutil.ReadFile(cfg.Reva.AppRegistry.MimetypesJSON)
// 		if err != nil {
// 			logger.Error().Err(err).Msg("Failed to read app registry mimetypes from JSON file: " + cfg.Reva.AppRegistry.MimetypesJSON)
// 			return nil
// 		}
// 		if err = json.Unmarshal(data, &mimetypes); err != nil {
// 			logger.Error().Err(err).Msg("Failed to unmarshal storage registry rules")
// 			return nil
// 		}
// 		if err := mapstructure.Decode(mimetypes, &m); err != nil {
// 			logger.Error().Err(err).Msg("Failed to decode defaultapp registry mimetypes to mapstructure")
// 			return nil
// 		}
// 		return m
// 	}

// 	logger.Info().Msg("No app registry mimetypes JSON file provided, loading default configuration")

// 	mimetypes = []mimeTypeConfig{
// 		{
// 			MimeType:    "application/pdf",
// 			Extension:   "pdf",
// 			Name:        "PDF",
// 			Description: "PDF document",
// 		},
// 		{
// 			MimeType:      "application/vnd.oasis.opendocument.text",
// 			Extension:     "odt",
// 			Name:          "OpenDocument",
// 			Description:   "OpenDocument text document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:      "application/vnd.oasis.opendocument.spreadsheet",
// 			Extension:     "ods",
// 			Name:          "OpenSpreadsheet",
// 			Description:   "OpenDocument spreadsheet document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:      "application/vnd.oasis.opendocument.presentation",
// 			Extension:     "odp",
// 			Name:          "OpenPresentation",
// 			Description:   "OpenDocument presentation document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
// 			Extension:     "docx",
// 			Name:          "Microsoft Word",
// 			Description:   "Microsoft Word document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
// 			Extension:     "xlsx",
// 			Name:          "Microsoft Excel",
// 			Description:   "Microsoft Excel document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
// 			Extension:     "pptx",
// 			Name:          "Microsoft PowerPoint",
// 			Description:   "Microsoft PowerPoint document",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:    "application/vnd.jupyter",
// 			Extension:   "ipynb",
// 			Name:        "Jupyter Notebook",
// 			Description: "Jupyter Notebook",
// 		},
// 		{
// 			MimeType:      "text/markdown",
// 			Extension:     "md",
// 			Name:          "Markdown file",
// 			Description:   "Markdown file",
// 			AllowCreation: true,
// 		},
// 		{
// 			MimeType:    "application/compressed-markdown",
// 			Extension:   "zmd",
// 			Name:        "Compressed markdown file",
// 			Description: "Compressed markdown file",
// 		},
// 	}

// 	if err := mapstructure.Decode(mimetypes, &m); err != nil {
// 		logger.Error().Err(err).Msg("Failed to decode defaultapp registry mimetypes to mapstructure")
// 		return nil
// 	}
// 	return m

// }

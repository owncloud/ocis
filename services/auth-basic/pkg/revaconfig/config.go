package revaconfig

import "github.com/owncloud/ocis/v2/services/auth-basic/pkg/config"

// AuthBasicConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func AuthBasicConfigFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"authprovider": map[string]interface{}{
					"auth_manager": cfg.AuthProvider,
					"auth_managers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.AuthProviders.JSON.File,
						},
						"ldap": ldapConfigFromString(cfg.AuthProviders.LDAP),
						"owncloudsql": map[string]interface{}{
							"dbusername":        cfg.AuthProviders.OwnCloudSQL.DBUsername,
							"dbpassword":        cfg.AuthProviders.OwnCloudSQL.DBPassword,
							"dbhost":            cfg.AuthProviders.OwnCloudSQL.DBHost,
							"dbport":            cfg.AuthProviders.OwnCloudSQL.DBPort,
							"dbname":            cfg.AuthProviders.OwnCloudSQL.DBName,
							"idp":               cfg.AuthProviders.OwnCloudSQL.IDP,
							"nobody":            cfg.AuthProviders.OwnCloudSQL.Nobody,
							"join_username":     cfg.AuthProviders.OwnCloudSQL.JoinUsername,
							"join_ownclouduuid": cfg.AuthProviders.OwnCloudSQL.JoinOwnCloudUUID,
						},
					},
				},
			},
		},
	}
	return rcfg
}

func ldapConfigFromString(cfg config.LDAPProvider) map[string]interface{} {
	return map[string]interface{}{
		"uri":               cfg.URI,
		"cacert":            cfg.CACert,
		"insecure":          cfg.Insecure,
		"bind_username":     cfg.BindDN,
		"bind_password":     cfg.BindPassword,
		"user_base_dn":      cfg.UserBaseDN,
		"group_base_dn":     cfg.GroupBaseDN,
		"user_filter":       cfg.UserFilter,
		"group_filter":      cfg.GroupFilter,
		"user_scope":        cfg.UserScope,
		"group_scope":       cfg.GroupScope,
		"user_objectclass":  cfg.UserObjectClass,
		"group_objectclass": cfg.GroupObjectClass,
		"login_attributes":  cfg.LoginAttributes,
		"idp":               cfg.IDP,
		"user_schema": map[string]interface{}{
			"id":              cfg.UserSchema.ID,
			"idIsOctetString": cfg.UserSchema.IDIsOctetString,
			"mail":            cfg.UserSchema.Mail,
			"displayName":     cfg.UserSchema.DisplayName,
			"userName":        cfg.UserSchema.Username,
		},
		"group_schema": map[string]interface{}{
			"id":              cfg.GroupSchema.ID,
			"idIsOctetString": cfg.GroupSchema.IDIsOctetString,
			"mail":            cfg.GroupSchema.Mail,
			"displayName":     cfg.GroupSchema.DisplayName,
			"groupName":       cfg.GroupSchema.Groupname,
			"member":          cfg.GroupSchema.Member,
		},
	}
}

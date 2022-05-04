package revaconfig

import (
	"github.com/owncloud/ocis/v2/extensions/group/pkg/config"
)

// GroupsConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func GroupsConfigFromStruct(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
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
				"groupprovider": map[string]interface{}{
					"driver": cfg.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"groups": cfg.Drivers.JSON.File,
						},
						"ldap": ldapConfigFromString(cfg.Drivers.LDAP),
						"rest": map[string]interface{}{
							"client_id":                      cfg.Drivers.REST.ClientID,
							"client_secret":                  cfg.Drivers.REST.ClientSecret,
							"redis_address":                  cfg.Drivers.REST.RedisAddr,
							"redis_username":                 cfg.Drivers.REST.RedisUsername,
							"redis_password":                 cfg.Drivers.REST.RedisPassword,
							"group_members_cache_expiration": cfg.GroupMembersCacheExpiration,
							"id_provider":                    cfg.Drivers.REST.IDProvider,
							"api_base_url":                   cfg.Drivers.REST.APIBaseURL,
							"oidc_token_endpoint":            cfg.Drivers.REST.OIDCTokenEndpoint,
							"target_api":                     cfg.Drivers.REST.TargetAPI,
						},
					},
				},
			},
		},
	}
}

func ldapConfigFromString(cfg config.LDAPDriver) map[string]interface{} {
	return map[string]interface{}{
		"uri":               cfg.URI,
		"cacert":            cfg.CACert,
		"insecure":          cfg.Insecure,
		"bind_username":     cfg.BindDN,
		"bind_password":     cfg.BindPassword,
		"user_base_dn":      cfg.UserBaseDN,
		"group_base_dn":     cfg.GroupBaseDN,
		"user_scope":        cfg.UserScope,
		"group_scope":       cfg.GroupScope,
		"user_filter":       cfg.UserFilter,
		"group_filter":      cfg.GroupFilter,
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

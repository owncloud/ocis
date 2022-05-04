package revaconfig

import (
	"github.com/owncloud/ocis/extensions/user/pkg/config"
)

// UsersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func UsersConfigFromStruct(cfg *config.Config) map[string]interface{} {
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
				"userprovider": map[string]interface{}{
					"driver": cfg.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.Drivers.JSON.File,
						},
						"ldap": ldapConfigFromString(cfg.Drivers.LDAP),
						"owncloudsql": map[string]interface{}{
							"dbusername":           cfg.Drivers.OwnCloudSQL.DBUsername,
							"dbpassword":           cfg.Drivers.OwnCloudSQL.DBPassword,
							"dbhost":               cfg.Drivers.OwnCloudSQL.DBHost,
							"dbport":               cfg.Drivers.OwnCloudSQL.DBPort,
							"dbname":               cfg.Drivers.OwnCloudSQL.DBName,
							"idp":                  cfg.Drivers.OwnCloudSQL.IDP,
							"nobody":               cfg.Drivers.OwnCloudSQL.Nobody,
							"join_username":        cfg.Drivers.OwnCloudSQL.JoinUsername,
							"join_ownclouduuid":    cfg.Drivers.OwnCloudSQL.JoinOwnCloudUUID,
							"enable_medial_search": cfg.Drivers.OwnCloudSQL.EnableMedialSearch,
						},
					},
				},
			},
		},
	}
	return rcfg
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

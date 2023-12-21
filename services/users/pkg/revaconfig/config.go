package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/users/pkg/config"
)

// UsersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func UsersConfigFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
			"grpc_client_options":       cfg.Reva.GetGRPCClientConfig(),
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			"tls_settings": map[string]interface{}{
				"enabled":     cfg.GRPC.TLS.Enabled,
				"certificate": cfg.GRPC.TLS.Cert,
				"key":         cfg.GRPC.TLS.Key,
			},
			// TODO build services dynamically
			"services": map[string]interface{}{
				"userprovider": map[string]interface{}{
					"driver": cfg.Driver,
					"drivers": map[string]interface{}{
						"json": map[string]interface{}{
							"users": cfg.Drivers.JSON.File,
						},
						"ldap": ldapConfigFromStruct(cfg),
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
			"interceptors": map[string]interface{}{
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "users",
				},
			},
		},
	}
	return rcfg
}

func ldapConfigFromStruct(cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"uri":                        cfg.Drivers.LDAP.URI,
		"cacert":                     cfg.Drivers.LDAP.CACert,
		"insecure":                   cfg.Drivers.LDAP.Insecure,
		"bind_username":              cfg.Drivers.LDAP.BindDN,
		"bind_password":              cfg.Drivers.LDAP.BindPassword,
		"user_base_dn":               cfg.Drivers.LDAP.UserBaseDN,
		"group_base_dn":              cfg.Drivers.LDAP.GroupBaseDN,
		"user_scope":                 cfg.Drivers.LDAP.UserScope,
		"group_scope":                cfg.Drivers.LDAP.GroupScope,
		"user_substring_filter_type": cfg.Drivers.LDAP.UserSubstringFilterType,
		"user_filter":                cfg.Drivers.LDAP.UserFilter,
		"group_filter":               cfg.Drivers.LDAP.GroupFilter,
		"user_objectclass":           cfg.Drivers.LDAP.UserObjectClass,
		"group_objectclass":          cfg.Drivers.LDAP.GroupObjectClass,
		"user_disable_mechanism":     cfg.Drivers.LDAP.DisableUserMechanism,
		"user_enabled_property":      cfg.Drivers.LDAP.UserSchema.Enabled,
		"user_type_property":         cfg.Drivers.LDAP.UserTypeAttribute,
		"group_local_disabled_dn":    cfg.Drivers.LDAP.LdapDisabledUsersGroupDN,
		"idp":                        cfg.Drivers.LDAP.IDP,
		"user_schema": map[string]interface{}{
			"id":              cfg.Drivers.LDAP.UserSchema.ID,
			"idIsOctetString": cfg.Drivers.LDAP.UserSchema.IDIsOctetString,
			"mail":            cfg.Drivers.LDAP.UserSchema.Mail,
			"displayName":     cfg.Drivers.LDAP.UserSchema.DisplayName,
			"userName":        cfg.Drivers.LDAP.UserSchema.Username,
		},
		"group_schema": map[string]interface{}{
			"id":              cfg.Drivers.LDAP.GroupSchema.ID,
			"idIsOctetString": cfg.Drivers.LDAP.GroupSchema.IDIsOctetString,
			"mail":            cfg.Drivers.LDAP.GroupSchema.Mail,
			"displayName":     cfg.Drivers.LDAP.GroupSchema.DisplayName,
			"groupName":       cfg.Drivers.LDAP.GroupSchema.Groupname,
			"member":          cfg.Drivers.LDAP.GroupSchema.Member,
		},
		"gateway_addr": cfg.Reva.Address,
	}
}

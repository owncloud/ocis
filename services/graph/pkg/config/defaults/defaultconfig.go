package defaults

import (
	"path"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.graph",
			Root:      "/graph",
		},
		Service: config.Service{
			Name: "graph",
		},
		Application: config.Application{
			DisplayName: "ownCloud Infinite Scale",
		},
		API: config.API{
			GroupMembersPatchLimit: 20,
			UsernameMatch:          "default",
		},
		Reva: shared.DefaultRevaConfig(),
		Spaces: config.Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
			// 30 minutes
			GroupsCacheTTL: 1800,
			// 30 minutes
			UsersCacheTTL: 1800,
		},
		Identity: config.Identity{
			Backend: "ldap",
			LDAP: config.LDAP{
				URI:                      "ldaps://localhost:9235",
				Insecure:                 false,
				CACert:                   path.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
				BindDN:                   "uid=libregraph,ou=sysusers,o=libregraph-idm",
				UseServerUUID:            false,
				UsePasswordModExOp:       true,
				WriteEnabled:             true,
				UserBaseDN:               "ou=users,o=libregraph-idm",
				UserSearchScope:          "sub",
				UserFilter:               "",
				UserObjectClass:          "inetOrgPerson",
				UserEmailAttribute:       "mail",
				UserDisplayNameAttribute: "displayName",
				UserNameAttribute:        "uid",
				// FIXME: switch this to some more widely available attribute by default
				//        ideally this needs to	be constant for the lifetime of a users
				UserIDAttribute:           "owncloudUUID",
				UserTypeAttribute:         "ownCloudUserType",
				UserEnabledAttribute:      "ownCloudUserEnabled",
				DisableUserMechanism:      "attribute",
				LdapDisabledUsersGroupDN:  "cn=DisabledUsersGroup,ou=groups,o=libregraph-idm",
				GroupBaseDN:               "ou=groups,o=libregraph-idm",
				GroupSearchScope:          "sub",
				GroupFilter:               "",
				GroupObjectClass:          "groupOfNames",
				GroupNameAttribute:        "cn",
				GroupIDAttribute:          "owncloudUUID",
				EducationResourcesEnabled: false,
			},
		},
		Events: config.Events{
			Endpoint:  "127.0.0.1:9233",
			Cluster:   "ocis-cluster",
			EnableTLS: false,
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for "envdecode".
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for "envdecode".
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.Tracing{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}

	if cfg.CacheStore == nil && cfg.Commons != nil && cfg.Commons.CacheStore != nil {
		cfg.CacheStore = &config.CacheStore{
			Type:    cfg.Commons.CacheStore.Type,
			Address: cfg.Commons.CacheStore.Address,
			Size:    cfg.Commons.CacheStore.Size,
		}
	} else if cfg.CacheStore == nil {
		cfg.CacheStore = &config.CacheStore{}
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}

	if cfg.Commons != nil {
		cfg.HTTP.TLS = cfg.Commons.HTTPServiceTLS
	}
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}

	// convert ttl to millisecond
	// the config is in seconds, therefore we need multiply it.
	cfg.Spaces.ExtendedSpacePropertiesCacheTTL = cfg.Spaces.ExtendedSpacePropertiesCacheTTL * int(time.Second)
	cfg.Spaces.GroupsCacheTTL = cfg.Spaces.GroupsCacheTTL * int(time.Second)
	cfg.Spaces.UsersCacheTTL = cfg.Spaces.UsersCacheTTL * int(time.Second)
}

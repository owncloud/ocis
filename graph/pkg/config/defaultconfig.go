package config

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.graph",
			Root:      "/graph",
		},
		Service: Service{
			Name: "graph",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
		},
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Spaces: Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
		},
		Identity: Identity{
			Backend: "cs3",
			LDAP: LDAP{
				URI:                      "ldap://localhost:9125",
				BindDN:                   "",
				BindPassword:             "",
				UserBaseDN:               "ou=users,dc=ocis,dc=test",
				UserSearchScope:          "sub",
				UserFilter:               "(objectClass=posixaccount)",
				UserEmailAttribute:       "mail",
				UserDisplayNameAttribute: "displayName",
				UserNameAttribute:        "uid",
				// FIXME: switch this to some more widely available attribute by default
				//        ideally this needs to	be constant for the lifetime of a users
				UserIDAttribute:    "ownclouduuid",
				GroupBaseDN:        "ou=groups,dc=ocis,dc=test",
				GroupSearchScope:   "sub",
				GroupFilter:        "(objectclass=groupOfNames)",
				GroupNameAttribute: "cn",
				GroupIDAttribute:   "cn",
			},
		},
	}
}

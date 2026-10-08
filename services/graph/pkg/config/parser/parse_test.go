package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
)

func TestValidateLDAPSettingsMasterIDRequiresUserFilter(t *testing.T) {
	cfg := &config.Config{
		Service: config.Service{Name: "graph"},
		Identity: config.Identity{
			LDAP: config.LDAP{
				BindPassword: "secret",
				UserFilter:   "",
			},
		},
		MultiInstance: config.MultiInstanceConfig{
			MasterID: "master-1",
		},
	}

	err := validateLDAPSettings(cfg)
	assert.Error(t, err)
}

func TestValidateLDAPSettingsMasterIDWithUserFilterPasses(t *testing.T) {
	cfg := &config.Config{
		Service: config.Service{Name: "graph"},
		Identity: config.Identity{
			LDAP: config.LDAP{
				BindPassword: "secret",
				UserFilter:   "(ownCloudMemberOf=instance-1)",
			},
		},
		MultiInstance: config.MultiInstanceConfig{
			MasterID: "master-1",
		},
	}

	err := validateLDAPSettings(cfg)
	assert.NoError(t, err)
}

func TestValidateLDAPSettingsNoMasterIDWithEmptyUserFilterPasses(t *testing.T) {
	cfg := &config.Config{
		Service: config.Service{Name: "graph"},
		Identity: config.Identity{
			LDAP: config.LDAP{
				BindPassword: "secret",
				UserFilter:   "",
			},
		},
		MultiInstance: config.MultiInstanceConfig{
			MasterID: "",
		},
	}

	err := validateLDAPSettings(cfg)
	assert.NoError(t, err)
}

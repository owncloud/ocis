package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/users/pkg/config"
)

func TestValidateMasterIDRequiresUserFilter(t *testing.T) {
	cfg := &config.Config{
		Service:      config.Service{Name: "users"},
		TokenManager: &config.TokenManager{JWTSecret: "secret"},
		Driver:       "ldap",
		Drivers: config.Drivers{
			LDAP: config.LDAPDriver{
				BindPassword: "secret",
				MasterID:     "master-1",
				UserFilter:   "",
			},
		},
	}

	err := Validate(cfg)
	assert.Error(t, err)
}

func TestValidateMasterIDWithUserFilterPasses(t *testing.T) {
	cfg := &config.Config{
		Service:      config.Service{Name: "users"},
		TokenManager: &config.TokenManager{JWTSecret: "secret"},
		Driver:       "ldap",
		Drivers: config.Drivers{
			LDAP: config.LDAPDriver{
				BindPassword: "secret",
				MasterID:     "master-1",
				UserFilter:   "(ownCloudMemberOf=instance-1)",
			},
		},
	}

	err := Validate(cfg)
	assert.NoError(t, err)
}

func TestValidateNoMasterIDWithEmptyUserFilterPasses(t *testing.T) {
	cfg := &config.Config{
		Service:      config.Service{Name: "users"},
		TokenManager: &config.TokenManager{JWTSecret: "secret"},
		Driver:       "ldap",
		Drivers: config.Drivers{
			LDAP: config.LDAPDriver{
				BindPassword: "secret",
				MasterID:     "",
				UserFilter:   "",
			},
		},
	}

	err := Validate(cfg)
	assert.NoError(t, err)
}

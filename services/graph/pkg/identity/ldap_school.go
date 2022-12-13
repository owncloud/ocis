package identity

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
)

type educationConfig struct {
	schoolBaseDN       string
	schoolFilter       string
	schoolObjectClass  string
	schoolScope        int
	schoolAttributeMap schoolAttributeMap
}

type schoolAttributeMap struct {
	displayName  string
	schoolNumber string
	id           string
}

func defaultEducationConfig() educationConfig {
	return educationConfig{
		schoolObjectClass:  "ocEducationSchool",
		schoolScope:        ldap.ScopeWholeSubtree,
		schoolAttributeMap: newSchoolAttributeMap(),
	}
}

func newEducationConfig(config config.LDAP) (educationConfig, error) {
	if config.EducationResourcesEnabled {
		var err error
		eduCfg := defaultEducationConfig()
		eduCfg.schoolBaseDN = config.EducationConfig.SchoolBaseDN
		if config.EducationConfig.SchoolSearchScope != "" {
			if eduCfg.schoolScope, err = stringToScope(config.EducationConfig.SchoolSearchScope); err != nil {
				return educationConfig{}, fmt.Errorf("error configuring school search scope: %w", err)
			}
		}
		return eduCfg, nil
	}
	return educationConfig{}, nil
}

func newSchoolAttributeMap() schoolAttributeMap {
	return schoolAttributeMap{
		displayName:  "ou",
		schoolNumber: "ocEducationSchoolNumber",
		id:           "owncloudUUID",
	}
}

// CreateSchool creates the supplied school in the identity backend.
func (i *LDAP) CreateSchool(ctx context.Context, school libregraph.EducationSchool) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// DeleteSchool deletes a given school, identified by id
func (i *LDAP) DeleteSchool(ctx context.Context, id string) error {
	return errNotImplemented
}

// GetSchool implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetSchool(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetSchools implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetSchools(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationSchool, error) {
	return nil, errNotImplemented
}

// GetSchoolMembers implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetSchoolMembers(ctx context.Context, id string) ([]*libregraph.User, error) {
	return nil, errNotImplemented
}

// AddMembersToSchool adds new members (reference by a slice of IDs) to supplied school in the identity backend.
func (i *LDAP) AddMembersToSchool(ctx context.Context, schoolID string, memberID []string) error {
	return errNotImplemented
}

// RemoveMemberFromSchool removes a single member (by ID) from a school
func (i *LDAP) RemoveMemberFromSchool(ctx context.Context, schoolID string, memberID string) error {
	return errNotImplemented
}

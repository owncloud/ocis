package identity

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofrs/uuid"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
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
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("CreateSchool")
	if !i.writeEnabled {
		return nil, errReadOnly
	}

	dn := fmt.Sprintf("%s=%s,%s",
		i.educationConfig.schoolAttributeMap.displayName,
		oldap.EscapeDNAttributeValue(school.GetDisplayName()),
		i.educationConfig.schoolBaseDN,
	)
	ar := ldap.NewAddRequest(dn, nil)
	ar.Attribute(i.educationConfig.schoolAttributeMap.displayName, []string{school.GetDisplayName()})
	ar.Attribute(i.educationConfig.schoolAttributeMap.schoolNumber, []string{school.GetSchoolNumber()})
	if !i.useServerUUID {
		ar.Attribute(i.educationConfig.schoolAttributeMap.id, []string{uuid.Must(uuid.NewV4()).String()})
	}
	objectClasses := []string{"organizationalUnit", i.educationConfig.schoolObjectClass, "top"}
	ar.Attribute("objectClass", objectClasses)

	if err := i.conn.Add(ar); err != nil {
		var lerr *ldap.Error
		logger.Debug().Err(err).Msg("error adding school")
		if errors.As(err, &lerr) {
			if lerr.ResultCode == ldap.LDAPResultEntryAlreadyExists {
				err = errorcode.New(errorcode.NameAlreadyExists, lerr.Error())
			}
		}
		return nil, err
	}

	// Read	back school from LDAP to get the generated UUID
	e, err := i.getSchoolByDN(ar.DN)
	if err != nil {
		return nil, err
	}
	return i.createSchoolModelFromLDAP(e), nil
}

// DeleteSchool deletes a given school, identified by id
func (i *LDAP) DeleteSchool(ctx context.Context, id string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("DeleteSchool")
	if !i.writeEnabled {
		return errReadOnly
	}
	e, err := i.getSchoolByID(id)
	if err != nil {
		return err
	}

	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}
	return nil
}

// GetSchool implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetSchool(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.EducationSchool, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetSchool")
	e, err := i.getSchoolByID(nameOrID)
	if err != nil {
		return nil, err
	}

	return i.createSchoolModelFromLDAP(e), nil
}

// GetSchools implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetSchools(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationSchool, error) {
	var filter string
	filter = fmt.Sprintf("(objectClass=%s)", i.educationConfig.schoolObjectClass)

	if i.educationConfig.schoolFilter != "" {
		filter = fmt.Sprintf("(&%s%s)", i.educationConfig.schoolFilter, filter)
	}

	searchRequest := ldap.NewSearchRequest(
		i.educationConfig.schoolBaseDN,
		i.educationConfig.schoolScope,
		ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{
			i.educationConfig.schoolAttributeMap.displayName,
			i.educationConfig.schoolAttributeMap.id,
			i.educationConfig.schoolAttributeMap.schoolNumber,
		},
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("GetSchools")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	schools := make([]*libregraph.EducationSchool, 0, len(res.Entries))
	for _, e := range res.Entries {
		school := i.createSchoolModelFromLDAP(e)
		// Skip invalid LDAP entries
		if school == nil {
			continue
		}
		schools = append(schools, school)
	}
	return schools, nil
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

func (i *LDAP) getSchoolByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.educationConfig.schoolAttributeMap.displayName,
		i.educationConfig.schoolAttributeMap.id,
		i.educationConfig.schoolAttributeMap.schoolNumber,
	}
	filter := fmt.Sprintf("(objectClass=%s)", i.educationConfig.schoolObjectClass)

	if i.educationConfig.schoolFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.educationConfig.schoolFilter)
	}
	return i.getEntryByDN(dn, attrs, filter)
}

func (i *LDAP) getSchoolByID(id string) (*ldap.Entry, error) {
	id = ldap.EscapeFilter(id)
	filter := fmt.Sprintf("(%s=%s)", i.educationConfig.schoolAttributeMap.id, id)
	return i.getSchoolByFilter(filter)
}

func (i *LDAP) getSchoolByFilter(filter string) (*ldap.Entry, error) {
	filter = fmt.Sprintf("(&%s(objectClass=%s)%s)",
		i.educationConfig.schoolFilter,
		i.educationConfig.schoolObjectClass,
		filter,
	)
	searchRequest := ldap.NewSearchRequest(
		i.educationConfig.schoolBaseDN,
		i.educationConfig.schoolScope,
		ldap.NeverDerefAliases, 1, 0, false,
		filter,
		[]string{
			i.educationConfig.schoolAttributeMap.displayName,
			i.educationConfig.schoolAttributeMap.id,
			i.educationConfig.schoolAttributeMap.schoolNumber,
		},
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("getSchoolByFilter")
	res, err := i.conn.Search(searchRequest)

	if err != nil {
		var errmsg string
		if lerr, ok := err.(*ldap.Error); ok {
			if lerr.ResultCode == ldap.LDAPResultSizeLimitExceeded {
				errmsg = fmt.Sprintf("too many results searching for school '%s'", filter)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).
					Str("schoolfilter", filter).Msg("too many results searching for school")
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	if len(res.Entries) == 0 {
		return nil, errNotFound
	}

	return res.Entries[0], nil
}

func (i *LDAP) createSchoolModelFromLDAP(e *ldap.Entry) *libregraph.EducationSchool {
	if e == nil {
		return nil
	}

	displayName := e.GetEqualFoldAttributeValue(i.educationConfig.schoolAttributeMap.displayName)
	id := e.GetEqualFoldAttributeValue(i.educationConfig.schoolAttributeMap.id)
	schoolNumber := e.GetEqualFoldAttributeValue(i.educationConfig.schoolAttributeMap.schoolNumber)

	if id != "" && displayName != "" && schoolNumber != "" {
		school := libregraph.NewEducationSchool()
		school.SetDisplayName(displayName)
		school.SetSchoolNumber(schoolNumber)
		school.SetId(id)
		return school
	}
	i.logger.Warn().Str("dn", e.DN).Str("id", id).Str("displayName", displayName).Str("schoolNumber", schoolNumber).Msg("Invalid School. Missing required attribute")
	return nil
}

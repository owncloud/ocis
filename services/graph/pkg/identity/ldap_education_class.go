package identity

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-ldap/ldap/v3"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

type educationClassAttributeMap struct {
	externalID     string
	classification string
}

func newEducationClassAttributeMap() educationClassAttributeMap {
	return educationClassAttributeMap{
		externalID:     "ocEducationExternalId",
		classification: "ocEducationClassType",
	}
}

// GetEducationClasses implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClasses(ctx context.Context, queryParam url.Values) ([]*libregraph.EducationClass, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetEducationClasses")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	var classFilter string
	if search != "" {
		search = ldap.EscapeFilter(search)
		classFilter = fmt.Sprintf(
			"(|(%s=%s)(%s=%s*)(%s=%s*))",
			i.educationConfig.classAttributeMap.externalID, search,
			i.groupAttributeMap.name, search,
			i.groupAttributeMap.id, search,
		)
	}
	classFilter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.educationConfig.classObjectClass, classFilter)

	classAttrs := i.getEducationClassAttrTypes(false)

	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, 0, 0, false,
		classFilter,
		classAttrs,
		nil,
	)
	logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("GetEducationClasses")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	classes := make([]*libregraph.EducationClass, 0, len(res.Entries))

	var c *libregraph.EducationClass
	for _, e := range res.Entries {
		if c = i.createEducationClassModelFromLDAP(e); c == nil {
			continue
		}
		classes = append(classes, c)
	}
	return classes, nil
}

// CreateEducationClass implements the EducationBackend interface for the LDAP backend.
// An EducationClass is mapped to an LDAP entry of the "groupOfNames" structural ObjectClass.
// With a few additional Attributes added on top via the "ocEducationClass" auxiallary ObjectClass.
func (i *LDAP) CreateEducationClass(ctx context.Context, class libregraph.EducationClass) (*libregraph.EducationClass, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("create educationClass")
	if !i.writeEnabled {
		return nil, errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	}
	ar, err := i.educationClassToAddRequest(class)
	if err != nil {
		return nil, err
	}

	if err := i.conn.Add(ar); err != nil {
		return nil, err
	}

	// Read	back group from LDAP to get the generated UUID
	e, err := i.getEducationClassByDN(ar.DN)
	if err != nil {
		return nil, err
	}
	return i.createEducationClassModelFromLDAP(e), nil
}

// GetEducationClass implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClass(ctx context.Context, id string, queryParam url.Values) (*libregraph.EducationClass, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetEducationClass")
	e, err := i.getEducationClassByID(id, false)
	if err != nil {
		return nil, err
	}
	var class *libregraph.EducationClass
	if class = i.createEducationClassModelFromLDAP(e); class == nil {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}
	return class, nil
}

// DeleteEducationClass implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) DeleteEducationClass(ctx context.Context, id string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("DeleteEducationClass")
	if !i.writeEnabled {
		return ErrReadOnly
	}
	e, err := i.getEducationClassByID(id, false)
	if err != nil {
		return err
	}

	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}

	// TODO update any users that are member of this school
	return nil
}

// GetEducationClassMembers implements the EducationBackend interface for the LDAP backend.
func (i *LDAP) GetEducationClassMembers(ctx context.Context, id string) ([]*libregraph.EducationUser, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetEducationClassMembers")
	e, err := i.getEducationClassByID(id, true)
	if err != nil {
		return nil, err
	}

	memberEntries, err := i.expandLDAPGroupMembers(ctx, e)
	result := make([]*libregraph.EducationUser, 0, len(memberEntries))
	if err != nil {
		return nil, err
	}
	for _, member := range memberEntries {
		if u := i.createEducationUserModelFromLDAP(member); u != nil {
			result = append(result, u)
		}
	}

	return result, nil
}

func (i *LDAP) educationClassToAddRequest(class libregraph.EducationClass) (*ldap.AddRequest, error) {
	plainGroup := i.educationClassToGroup(class)
	ldapAttrs, err := i.groupToLDAPAttrValues(*plainGroup)
	if err != nil {
		return nil, err
	}
	ldapAttrs, err = i.educationClassToLDAPAttrValues(class, ldapAttrs)
	if err != nil {
		return nil, err
	}

	ar := ldap.NewAddRequest(i.getEducationClassLDAPDN(class), nil)

	for attrType, values := range ldapAttrs {
		ar.Attribute(attrType, values)
	}
	return ar, nil
}

func (i *LDAP) educationClassToGroup(class libregraph.EducationClass) *libregraph.Group {
	group := libregraph.NewGroup()
	group.SetDisplayName(class.DisplayName)

	return group
}

func (i *LDAP) educationClassToLDAPAttrValues(class libregraph.EducationClass, attrs ldapAttributeValues) (ldapAttributeValues, error) {
	if externalID, ok := class.GetExternalIdOk(); ok {
		attrs[i.educationConfig.classAttributeMap.externalID] = []string{*externalID}
	}
	if classification, ok := class.GetClassificationOk(); ok {
		attrs[i.educationConfig.classAttributeMap.classification] = []string{*classification}
	}
	attrs["objectClass"] = append(attrs["objectClass"], i.educationConfig.classObjectClass)
	return attrs, nil
}

func (i *LDAP) getEducationClassAttrTypes(requestMembers bool) []string {
	attrs := []string{
		i.groupAttributeMap.name,
		i.groupAttributeMap.id,
		i.educationConfig.classAttributeMap.classification,
		i.educationConfig.classAttributeMap.externalID,
	}
	if requestMembers {
		attrs = append(attrs, i.groupAttributeMap.member)
	}
	return attrs
}

func (i *LDAP) getEducationClassByDN(dn string) (*ldap.Entry, error) {
	filter := fmt.Sprintf("(objectClass=%s)", i.educationConfig.classObjectClass)

	if i.groupFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.groupFilter)
	}

	return i.getEntryByDN(dn, i.getEducationClassAttrTypes(false), filter)
}

func (i *LDAP) createEducationClassModelFromLDAP(e *ldap.Entry) *libregraph.EducationClass {
	group := i.createGroupModelFromLDAP(e)
	return i.groupToEducationClass(*group, e)
}

func (i *LDAP) groupToEducationClass(group libregraph.Group, e *ldap.Entry) *libregraph.EducationClass {
	class := libregraph.NewEducationClass(group.GetDisplayName(), "")
	class.SetId(group.GetId())

	if e != nil {
		// Set the education User specific Attributes from the supplied LDAP Entry
		if externalID := e.GetEqualFoldAttributeValue(i.educationConfig.classAttributeMap.externalID); externalID != "" {
			class.SetExternalId(externalID)
		}
		if classification := e.GetEqualFoldAttributeValue(i.educationConfig.classAttributeMap.classification); classification != "" {
			class.SetClassification(classification)
		}
	}

	return class
}

func (i *LDAP) getEducationClassLDAPDN(class libregraph.EducationClass) string {
	return fmt.Sprintf("ocEducationExternalId=%s,%s", oldap.EscapeDNAttributeValue(class.GetExternalId()), i.groupBaseDN)
}

// getEducationClassByID looks up a class by id (and be ID or externalID)
func (i *LDAP) getEducationClassByID(id string, requestMembers bool) (*ldap.Entry, error) {
	id = ldap.EscapeFilter(id)
	filter := fmt.Sprintf("(|(%s=%s)(%s=%s))",
		i.groupAttributeMap.id, id,
		i.educationConfig.classAttributeMap.externalID, id)
	return i.getEducationClassByFilter(filter, requestMembers)
}

func (i *LDAP) getEducationClassByFilter(filter string, requestMembers bool) (*ldap.Entry, error) {
	filter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.educationConfig.classObjectClass, filter)
	return i.searchLDAPEntryByFilter(i.groupBaseDN, i.getEducationClassAttrTypes(requestMembers), filter)
}

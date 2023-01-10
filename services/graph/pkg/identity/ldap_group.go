package identity

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofrs/uuid"
	ldapdn "github.com/libregraph/idm/pkg/ldapdn"
	libregraph "github.com/owncloud/libre-graph-api-go"
	oldap "github.com/owncloud/ocis/v2/ocis-pkg/ldap"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"golang.org/x/exp/slices"
)

type groupAttributeMap struct {
	name         string
	id           string
	member       string
	memberSyntax string
}

// GetGroup implements the Backend Interface for the LDAP Backend
func (i *LDAP) GetGroup(ctx context.Context, nameOrID string, queryParam url.Values) (*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroup")
	e, err := i.getLDAPGroupByNameOrID(nameOrID, true)
	if err != nil {
		return nil, err
	}
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	var g *libregraph.Group
	if g = i.createGroupModelFromLDAP(e); g == nil {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}
	if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
		members, err := i.expandLDAPGroupMembers(ctx, e)
		if err != nil {
			return nil, err
		}
		if len(members) > 0 {
			m := make([]libregraph.User, 0, len(members))
			for _, ue := range members {
				if u := i.createUserModelFromLDAP(ue); u != nil {
					m = append(m, *u)
				}
			}
			g.Members = m
		}
	}
	return g, nil
}

// GetGroups implements the Backend Interface for the LDAP Backend
func (i *LDAP) GetGroups(ctx context.Context, queryParam url.Values) ([]*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroups")

	search := queryParam.Get("search")
	if search == "" {
		search = queryParam.Get("$search")
	}

	var expandMembers bool
	sel := strings.Split(queryParam.Get("$select"), ",")
	exp := strings.Split(queryParam.Get("$expand"), ",")
	if slices.Contains(sel, "members") || slices.Contains(exp, "members") {
		expandMembers = true
	}

	var groupFilter string
	if search != "" {
		search = ldap.EscapeFilter(search)
		groupFilter = fmt.Sprintf(
			"(|(%s=%s*)(%s=%s*))",
			i.groupAttributeMap.name, search,
			i.groupAttributeMap.id, search,
		)
	}
	groupFilter = fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.groupObjectClass, groupFilter)

	groupAttrs := []string{
		i.groupAttributeMap.name,
		i.groupAttributeMap.id,
	}
	if expandMembers {
		groupAttrs = append(groupAttrs, i.groupAttributeMap.member)
	}

	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, 0, 0, false,
		groupFilter,
		groupAttrs,
		nil,
	)
	logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("GetGroups")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		return nil, errorcode.New(errorcode.ItemNotFound, err.Error())
	}

	groups := make([]*libregraph.Group, 0, len(res.Entries))

	var g *libregraph.Group
	for _, e := range res.Entries {
		if g = i.createGroupModelFromLDAP(e); g == nil {
			continue
		}
		if expandMembers {
			members, err := i.expandLDAPGroupMembers(ctx, e)
			if err != nil {
				return nil, err
			}
			if len(members) > 0 {
				m := make([]libregraph.User, 0, len(members))
				for _, ue := range members {
					if u := i.createUserModelFromLDAP(ue); u != nil {
						m = append(m, *u)
					}
				}
				g.Members = m
			}
		}
		groups = append(groups, g)
	}
	return groups, nil
}

// GetGroupMembers implements the Backend Interface for the LDAP Backend
func (i *LDAP) GetGroupMembers(ctx context.Context, groupID string) ([]*libregraph.User, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("GetGroupMembers")
	e, err := i.getLDAPGroupByNameOrID(groupID, true)
	if err != nil {
		return nil, err
	}

	memberEntries, err := i.expandLDAPGroupMembers(ctx, e)
	result := make([]*libregraph.User, 0, len(memberEntries))
	if err != nil {
		return nil, err
	}
	for _, member := range memberEntries {
		if u := i.createUserModelFromLDAP(member); u != nil {
			result = append(result, u)
		}
	}

	return result, nil
}

// CreateGroup implements the Backend Interface for the LDAP Backend
// It is currently restricted to managing groups based on the "groupOfNames" ObjectClass.
// As "groupOfNames" requires a "member" Attribute to be present. Empty Groups (groups
// without a member) a represented by adding an empty DN as the single member.
func (i *LDAP) CreateGroup(ctx context.Context, group libregraph.Group) (*libregraph.Group, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("create group")
	if !i.writeEnabled {
		return nil, errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	}
	ar, err := i.groupToAddRequest(group)
	if err != nil {
		return nil, err
	}

	if err := i.conn.Add(ar); err != nil {
		return nil, err
	}

	// Read	back group from LDAP to get the generated UUID
	e, err := i.getGroupByDN(ar.DN)
	if err != nil {
		return nil, err
	}
	return i.createGroupModelFromLDAP(e), nil
}

// DeleteGroup implements the Backend Interface.
func (i *LDAP) DeleteGroup(ctx context.Context, id string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("DeleteGroup")
	if !i.writeEnabled {
		return errorcode.New(errorcode.NotAllowed, "server is configured read-only")
	}
	e, err := i.getLDAPGroupByID(id, false)
	if err != nil {
		return err
	}
	dr := ldap.DelRequest{DN: e.DN}
	if err = i.conn.Del(&dr); err != nil {
		return err
	}
	return nil
}

// AddMembersToGroup implements the Backend Interface for the LDAP backend.
// Currently it is limited to adding Users as Group members. Adding other groups
// as members is not yet implemented
func (i *LDAP) AddMembersToGroup(ctx context.Context, groupID string, memberIDs []string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("AddMembersToGroup")
	ge, err := i.getLDAPGroupByID(groupID, true)
	if err != nil {
		return err
	}

	mr := ldap.ModifyRequest{DN: ge.DN}
	// Handle empty groups (using the empty member attribute)
	current := ge.GetEqualFoldAttributeValues(i.groupAttributeMap.member)
	if len(current) == 1 && current[0] == "" {
		mr.Delete(i.groupAttributeMap.member, []string{""})
	}

	// Create a Set of current members for faster lookups
	currentSet := make(map[string]struct{}, len(current))
	for _, currentMember := range current {
		// We can ignore any empty member value here
		if currentMember == "" {
			continue
		}
		nCurrentMember, err := ldapdn.ParseNormalize(currentMember)
		if err != nil {
			// We couldn't parse the member value as a DN. Let's skip it, but log a warning
			logger.Warn().Str("memberDN", currentMember).Err(err).Msg("Couldn't parse DN")
			continue
		}
		currentSet[nCurrentMember] = struct{}{}
	}

	var newMemberDN []string
	for _, memberID := range memberIDs {
		me, err := i.getLDAPUserByID(memberID)
		if err != nil {
			return err
		}
		nDN, err := ldapdn.ParseNormalize(me.DN)
		if err != nil {
			logger.Error().Str("new member", me.DN).Err(err).Msg("Couldn't parse DN")
			return err
		}
		if _, present := currentSet[nDN]; !present {
			newMemberDN = append(newMemberDN, me.DN)
		} else {
			logger.Debug().Str("memberDN", me.DN).Msg("Member already present in group. Skipping")
		}
	}

	if len(newMemberDN) > 0 {
		mr.Add(i.groupAttributeMap.member, newMemberDN)

		if err := i.conn.Modify(&mr); err != nil {
			return err
		}
	}
	return nil
}

// RemoveMemberFromGroup implements the Backend Interface.
func (i *LDAP) RemoveMemberFromGroup(ctx context.Context, groupID string, memberID string) error {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("RemoveMemberFromGroup")
	ge, err := i.getLDAPGroupByID(groupID, true)
	if err != nil {
		logger.Debug().Str("backend", "ldap").Str("groupID", groupID).Msg("Error looking up group")
		return err
	}
	me, err := i.getLDAPUserByID(memberID)
	if err != nil {
		logger.Debug().Str("backend", "ldap").Str("memberID", memberID).Msg("Error looking up group member")
		return err
	}
	logger.Debug().Str("backend", "ldap").Str("groupdn", ge.DN).Str("member", me.DN).Msg("remove member")

	if mr, err := i.removeMemberFromGroupEntry(ge, me.DN); err == nil {
		return i.conn.Modify(mr)
	}
	return nil
}

func (i *LDAP) groupToAddRequest(group libregraph.Group) (*ldap.AddRequest, error) {
	ar := ldap.NewAddRequest(i.getGroupLDAPDN(group), nil)

	attrMap, err := i.groupToLDAPAttrValues(group)
	if err != nil {
		return nil, err
	}
	for attrType, values := range attrMap {
		ar.Attribute(attrType, values)
	}
	return ar, nil
}

func (i *LDAP) getGroupLDAPDN(group libregraph.Group) string {
	return fmt.Sprintf("cn=%s,%s", oldap.EscapeDNAttributeValue(group.GetDisplayName()), i.groupBaseDN)
}

func (i *LDAP) groupToLDAPAttrValues(group libregraph.Group) (map[string][]string, error) {
	attrs := map[string][]string{
		i.groupAttributeMap.name: {group.GetDisplayName()},
		"objectClass":            {"groupOfNames", "top"},
		// This is a crutch to allow groups without members for LDAP servers
		// that apply strict Schema checking. The RFCs define "member/uniqueMember"
		// as required attribute for groupOfNames/groupOfUniqueNames. So we
		// add an empty string (which is a valid DN) as the initial member.
		// It will be replaced once real members are added.
		// We might wanna use the newer, but not so broadly used "groupOfMembers"
		// objectclass (RFC2307bis-02) where "member" is optional.
		i.groupAttributeMap.member: {""},
	}

	if !i.useServerUUID {
		attrs["owncloudUUID"] = []string{uuid.Must(uuid.NewV4()).String()}
		attrs["objectClass"] = append(attrs["objectClass"], "owncloud")
	}
	return attrs, nil
}

func (i *LDAP) expandLDAPGroupMembers(ctx context.Context, e *ldap.Entry) ([]*ldap.Entry, error) {
	logger := i.logger.SubloggerWithRequestID(ctx)
	logger.Debug().Str("backend", "ldap").Msg("expandLDAPGroupMembers")
	result := []*ldap.Entry{}

	for _, memberDN := range e.GetEqualFoldAttributeValues(i.groupAttributeMap.member) {
		if memberDN == "" {
			continue
		}
		logger.Debug().Str("memberDN", memberDN).Msg("lookup")
		ue, err := i.getUserByDN(memberDN)
		if err != nil {
			// Ignore errors when reading a specific member fails, just log them and continue
			logger.Debug().Err(err).Str("member", memberDN).Msg("error reading group member")
			continue
		}
		result = append(result, ue)
	}

	return result, nil
}

func (i *LDAP) getLDAPGroupByID(id string, requestMembers bool) (*ldap.Entry, error) {
	id = ldap.EscapeFilter(id)
	filter := fmt.Sprintf("(%s=%s)", i.groupAttributeMap.id, id)
	return i.getLDAPGroupByFilter(filter, requestMembers)
}

func (i *LDAP) getLDAPGroupByNameOrID(nameOrID string, requestMembers bool) (*ldap.Entry, error) {
	nameOrID = ldap.EscapeFilter(nameOrID)
	filter := fmt.Sprintf("(|(%s=%s)(%s=%s))", i.groupAttributeMap.name, nameOrID, i.groupAttributeMap.id, nameOrID)
	return i.getLDAPGroupByFilter(filter, requestMembers)
}

func (i *LDAP) getLDAPGroupByFilter(filter string, requestMembers bool) (*ldap.Entry, error) {
	e, err := i.getLDAPGroupsByFilter(filter, requestMembers, true)
	if err != nil {
		return nil, err
	}
	if len(e) == 0 {
		return nil, errorcode.New(errorcode.ItemNotFound, "not found")
	}

	return e[0], nil
}

// Search for LDAP Groups matching the specified filter, if requestMembers is true the groupMemberShip
// attribute will be part of the result attributes. The LDAP filter is combined with the configured groupFilter
// resulting in a filter like "(&(LDAP.groupFilter)(objectClass=LDAP.groupObjectClass)(<filter_from_args>))"
func (i *LDAP) getLDAPGroupsByFilter(filter string, requestMembers, single bool) ([]*ldap.Entry, error) {
	attrs := []string{
		i.groupAttributeMap.name,
		i.groupAttributeMap.id,
	}

	if requestMembers {
		attrs = append(attrs, i.groupAttributeMap.member)
	}

	sizelimit := 0
	if single {
		sizelimit = 1
	}
	searchRequest := ldap.NewSearchRequest(
		i.groupBaseDN, i.groupScope, ldap.NeverDerefAliases, sizelimit, 0, false,
		fmt.Sprintf("(&%s(objectClass=%s)%s)", i.groupFilter, i.groupObjectClass, filter),
		attrs,
		nil,
	)
	i.logger.Debug().Str("backend", "ldap").
		Str("base", searchRequest.BaseDN).
		Str("filter", searchRequest.Filter).
		Int("scope", searchRequest.Scope).
		Int("sizelimit", searchRequest.SizeLimit).
		Interface("attributes", searchRequest.Attributes).
		Msg("getLDAPGroupsByFilter")
	res, err := i.conn.Search(searchRequest)
	if err != nil {
		var errmsg string
		if lerr, ok := err.(*ldap.Error); ok {
			if lerr.ResultCode == ldap.LDAPResultSizeLimitExceeded {
				errmsg = fmt.Sprintf("too many results searching for group '%s'", filter)
				i.logger.Debug().Str("backend", "ldap").Err(lerr).Msg(errmsg)
			}
		}
		return nil, errorcode.New(errorcode.ItemNotFound, errmsg)
	}
	return res.Entries, nil
}

// removeMemberFromGroupEntry creates an LDAP Modify request (not sending it)
// that would update the supplied entry to remove the specified member from the
// group
func (i *LDAP) removeMemberFromGroupEntry(group *ldap.Entry, memberDN string) (*ldap.ModifyRequest, error) {
	nOldMemberDN, err := ldapdn.ParseNormalize(memberDN)
	if err != nil {
		return nil, err
	}
	members := group.GetEqualFoldAttributeValues(i.groupAttributeMap.member)
	found := false
	for _, member := range members {
		if member == "" {
			continue
		}
		if nMember, err := ldapdn.ParseNormalize(member); err != nil {
			// We couldn't parse the member value as a DN. Let's keep it
			// as it is but log a warning
			i.logger.Warn().Str("memberDN", member).Err(err).Msg("Couldn't parse DN")
			continue
		} else {
			if nMember == nOldMemberDN {
				found = true
			}
		}
	}
	if !found {
		i.logger.Debug().Str("backend", "ldap").Str("groupdn", group.DN).Str("member", memberDN).
			Msg("The target is not a member of the group")
		return nil, ErrNotFound
	}

	mr := ldap.ModifyRequest{DN: group.DN}
	if len(members) == 1 {
		mr.Add(i.groupAttributeMap.member, []string{""})
	}
	mr.Delete(i.groupAttributeMap.member, []string{memberDN})
	return &mr, nil
}

func (i *LDAP) getGroupByDN(dn string) (*ldap.Entry, error) {
	attrs := []string{
		i.groupAttributeMap.id,
		i.groupAttributeMap.name,
	}
	filter := fmt.Sprintf("(objectClass=%s)", i.groupObjectClass)

	if i.groupFilter != "" {
		filter = fmt.Sprintf("(&%s(%s))", filter, i.groupFilter)
	}
	return i.getEntryByDN(dn, attrs, filter)
}

func (i *LDAP) getGroupsForUser(dn string) ([]*ldap.Entry, error) {
	groupFilter := fmt.Sprintf(
		"(%s=%s)",
		i.groupAttributeMap.member, dn,
	)
	userGroups, err := i.getLDAPGroupsByFilter(groupFilter, false, false)
	if err != nil {
		return nil, err
	}
	return userGroups, nil
}

func (i *LDAP) createGroupModelFromLDAP(e *ldap.Entry) *libregraph.Group {
	name := e.GetEqualFoldAttributeValue(i.groupAttributeMap.name)
	id := e.GetEqualFoldAttributeValue(i.groupAttributeMap.id)

	if id != "" && name != "" {
		return &libregraph.Group{
			DisplayName: &name,
			Id:          &id,
		}
	}
	i.logger.Warn().Str("dn", e.DN).Msg("Group is missing name or id")
	return nil
}

func (i *LDAP) groupsFromLDAPEntries(e []*ldap.Entry) []libregraph.Group {
	groups := make([]libregraph.Group, 0, len(e))
	for _, g := range e {
		if grp := i.createGroupModelFromLDAP(g); grp != nil {
			groups = append(groups, *grp)
		}
	}
	return groups
}

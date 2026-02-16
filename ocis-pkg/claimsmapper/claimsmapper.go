package claimsmapper

import (
	"regexp"
	"strings"
)

// ClaimsMapper is a configurable mapper to map oidc claims to ocis spaceIDs and roles
type ClaimsMapper struct {
	claimRegexp *regexp.Regexp
	roleMapping map[string]string
}

// NewClaimsMapper parses the config to create a new ClaimsMapper. It expects
//   - a regexp extracting the spaceID and the (unmapped) role from the claim.
//   - a roleMapping string of the form "oidcRole1:manager","oidcRole2:editor","oidcRole3:viewer"
//     unused roles can be omitted. Second part of the mapping must be a valid ocis role.
//     can be omitted if roles already match ocis roles
//
// Panics if regexp is not compilable
func NewClaimsMapper(reg string, roleMapping []string) ClaimsMapper {
	em := ClaimsMapper{
		claimRegexp: regexp.MustCompile(reg),
	}

	if len(roleMapping) == 0 {
		return em
	}

	em.roleMapping = make(map[string]string)
	for _, ms := range roleMapping {
		s := strings.Split(ms, ":")
		if len(s) != 2 {
			continue
		}
		em.roleMapping[s[0]] = s[1]
	}
	return em
}

// Exec extracts the spaceID and the role from a entitlement
func (em ClaimsMapper) Exec(e string) (match bool, spaceID string, role string) {
	s := em.claimRegexp.FindStringSubmatch(e)
	if len(s) != 3 {
		return
	}

	spaceID = s[1]
	if spaceID == "" {
		return
	}

	role = s[2]
	if em.roleMapping == nil {
		match = true
		return
	}

	role = em.roleMapping[role]
	if role != "" {
		match = true
		return
	}
	return false, "", ""
}

// Package rolemapping holds what the proxy and the graph service share to map values of a user to
// the name of a role: the proxy maps the values of the roles claim at login, the graph service the
// userClass of multi-instance guests who did not log in yet.
package rolemapping

import (
	"fmt"
	"regexp"
	"strings"
)

// ParsePair splits a 'claimvalue:rolename' pair of OCIS_MULTI_INSTANCE_GUEST_ROLE_MAPPING at its
// last colon, so the value may contain colons but the role name must not. Roles created through
// the settings service can't contain a colon, their names are restricted to letters, digits, '-'
// and '_'.
func ParsePair(pair string) (value string, roleName string, err error) {
	i := strings.LastIndex(pair, ":")
	// the env var parser trims list elements, the yaml parser does not
	if i >= 0 {
		value, roleName = strings.TrimSpace(pair[:i]), strings.TrimSpace(pair[i+1:])
	}
	if value == "" || roleName == "" {
		return "", "", fmt.Errorf("invalid guest role mapping '%s', expected 'claimvalue:rolename'", pair)
	}
	return value, roleName, nil
}

// Matcher matches values against the value of a mapping: an exact value, or a regular expression
// that must match a whole value.
type Matcher struct {
	value string
	// rx is nil if the value is no valid regular expression, then only the exact value matches
	rx *regexp.Regexp
}

// NewMatcher compiles the value of a mapping.
func NewMatcher(value string) Matcher {
	rx, err := regexp.Compile("^(?:" + value + ")$")
	if err != nil {
		rx = nil
	}
	return Matcher{value: value, rx: rx}
}

// Matches reports whether v matches the value of the mapping.
func (m Matcher) Matches(v string) bool {
	return v == m.value || (m.rx != nil && m.rx.MatchString(v))
}

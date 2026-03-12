package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// EnhanceFilterWithMasterID OR's the given LDAP filter with master-id match clauses.
// Returns filter unchanged when masterID is empty or both attributes are empty.
func EnhanceFilterWithMasterID(filter, masterID, memberAttr, guestAttr string) string {
	if masterID == "" {
		return filter
	}

	if memberAttr == "" && guestAttr == "" {
		return filter
	}

	var masterIDParts []string
	if memberAttr != "" {
		masterIDParts = append(masterIDParts,
			fmt.Sprintf("(%s=%s)", memberAttr, ldap.EscapeFilter(masterID)))
	}
	if guestAttr != "" {
		masterIDParts = append(masterIDParts,
			fmt.Sprintf("(%s=%s)", guestAttr, ldap.EscapeFilter(masterID)))
	}

	var masterIDFilter string
	if len(masterIDParts) == 1 {
		masterIDFilter = masterIDParts[0]
	} else if len(masterIDParts) > 1 {
		masterIDFilter = fmt.Sprintf("(|%s)", strings.Join(masterIDParts, ""))
	}

	if filter == "" {
		return masterIDFilter
	}

	return fmt.Sprintf("(|%s%s)", filter, masterIDFilter)
}

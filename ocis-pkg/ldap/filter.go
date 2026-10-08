package ldap

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// EnhanceFilterWithMasterID adds master-ID match clauses to the given LDAP filter with OR.
// The function returns the filter without a change in three cases.
// The first case is an empty masterID. The second case is two empty attributes.
// The third case is an empty filter.
// An empty filter means "no restriction".
// Do not OR an empty filter with a master-ID clause.
// This action would restrict an unrestricted search to the master ID only.
func EnhanceFilterWithMasterID(filter, masterID, memberAttr, guestAttr string) string {
	if masterID == "" {
		return filter
	}

	if memberAttr == "" && guestAttr == "" {
		return filter
	}

	if filter == "" {
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

	return fmt.Sprintf("(|%s%s)", filter, masterIDFilter)
}

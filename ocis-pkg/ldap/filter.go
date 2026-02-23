package ldap

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// EnhanceUserFilterFromEnv OR's the given LDAP filter with master-id match clauses
// read from OCIS_MULTI_INSTANCE_MASTER_ID, OCIS_LDAP_USER_MEMBER_ATTRIBUTE, and
// OCIS_LDAP_USER_GUEST_ATTRIBUTE. Returns filter unchanged when env vars are unset.
func EnhanceUserFilterFromEnv(filter string) string {
	return enhanceFilterWithMasterID(
		filter,
		os.Getenv("OCIS_MULTI_INSTANCE_MASTER_ID"),
		os.Getenv("OCIS_LDAP_USER_MEMBER_ATTRIBUTE"),
		os.Getenv("OCIS_LDAP_USER_GUEST_ATTRIBUTE"),
	)
}

func enhanceFilterWithMasterID(filter, masterID, memberAttr, guestAttr string) string {
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

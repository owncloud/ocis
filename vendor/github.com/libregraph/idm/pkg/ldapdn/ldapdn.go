package ldapdn

import (
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/cases"
)

// Normalize takes an ldap.DN struct and turns it into a "normalized" DN string
// by cases folding all RDN (attributetypes and values). Note: This currently
// handles all attributes as caseIgnoreStrings ignoring the Syntax the Attribute
// Type might have assigned.
func Normalize(dn *ldap.DN) string {
	var nDN string
	caseFold := cases.Fold()
	for r, rdn := range dn.RDNs {
		// FIXME to really normalize multivalued RDNs we'd need
		// to normalize the order of Attributes here as well
		for a, ava := range rdn.Attributes {
			if a > 0 {
				// This is a multivalued RDN.
				nDN += "+"
			} else if r > 0 {
				nDN += ","
			}
			nDN = nDN + caseFold.String(ava.Type) + "=" + caseFold.String(ava.Value)
		}
	}
	return nDN
}

// ParseNormalize normalizes the passed LDAP DN string by first parsing it (using ldap.ParseDN)
// and then casefolding all RDN using ldapdn.Normalize(). ParseNormalize will return an error
// when parsing the DN fails.
func ParseNormalize(dn string) (string, error) {
	parsed, err := ldap.ParseDN(dn)
	if err != nil {
		return "", err
	}
	return Normalize(parsed), nil
}

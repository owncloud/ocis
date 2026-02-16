package ldapdn

import (
	"bytes"
	"encoding/hex"

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
			nDN = nDN + caseFold.String(ava.Type) + "=" + encodeRDNValue(caseFold.String(ava.Value))
		}
	}
	return nDN
}

// encodeRDNValue applies the DN escaping rules (RFC4514) to the supplied
// string (the value part of an RDN). Returns the escaped string.
// Note: This function is taken from https://github.com/go-ldap/ldap/pull/104
func encodeRDNValue(rDNValue string) string {
	encodedBuf := bytes.Buffer{}

	escapeChar := func(c byte) {
		encodedBuf.WriteByte('\\')
		encodedBuf.WriteByte(c)
	}

	escapeHex := func(c byte) {
		encodedBuf.WriteByte('\\')
		encodedBuf.WriteString(hex.EncodeToString([]byte{c}))
	}

	for i := 0; i < len(rDNValue); i++ {
		char := rDNValue[i]
		if i == 0 && char == ' ' || char == '#' {
			// Special case leading space or number sign.
			escapeChar(char)
			continue
		}
		if i == len(rDNValue)-1 && char == ' ' {
			// Special case trailing space.
			escapeChar(char)
			continue
		}

		switch char {
		case '"', '+', ',', ';', '<', '>', '\\':
			// Each of these special characters must be escaped.
			escapeChar(char)
			continue
		}

		if char < ' ' || char > '~' {
			// All special character escapes are handled first
			// above. All bytes less than ASCII SPACE and all bytes
			// greater than ASCII TILDE must be hex-escaped.
			escapeHex(char)
			continue
		}

		// Any other character does not require escaping.
		encodedBuf.WriteByte(char)
	}

	return encodedBuf.String()
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

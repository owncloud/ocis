package ldapentry

import (
	"errors"
	"log"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/cases"

	"github.com/libregraph/idm/pkg/ldapdn"
)

func ApplyModify(old *ldap.Entry, mod *ldap.ModifyRequest) (newEntry *ldap.Entry, err error) {
	parsed, err := ldap.ParseDN(old.DN)
	if err != nil {
		return nil, err
	}
	nOldDN := ldapdn.Normalize(parsed)
	rdn := parsed.RDNs[0]
	nReqDN, err := ldapdn.ParseNormalize(mod.DN)
	if err != nil {
		return nil, err
	}
	// This shouldn't happen if we ge here (TM)
	if nOldDN != nReqDN {
		return nil, ldap.NewError(ldap.LDAPResultUnwillingToPerform, errors.New("DNs do not match"))
	}

	casefold := cases.Fold()
	newEntry = ldap.NewEntry(old.DN, map[string][]string{})
	newEntry.Attributes = old.Attributes
	for _, c := range mod.Changes {
		nType := casefold.String(c.Modification.Type)
		switch c.Operation {
		case ldap.AddAttribute:
			log.Printf("applying add for Attribute: %s", c.Modification.Type)
			newValues := entryApplyModAdd(newEntry.GetEqualFoldAttributeValues(nType), c.Modification.Vals)
			newEntry.Attributes = entryReplaceValues(newEntry.Attributes, c.Modification.Type, newValues)

		case ldap.ReplaceAttribute:
			log.Printf("applying replace for Attribute: %s", c.Modification.Type)
			// Modifies on RDN attributes need special care to make sure that the rdn Value is not removed
			for _, rdnAttr := range rdn.Attributes {
				if nType == casefold.String(rdnAttr.Type) {
					rdnPresent := false
					nRdnVal := casefold.String(rdnAttr.Value)
					for _, newVal := range c.Modification.Vals {
						if nRdnVal == casefold.String(newVal) {
							rdnPresent = true
							break
						}
					}
					if !rdnPresent {
						return nil, ldap.NewError(ldap.LDAPResultNotAllowedOnRDN, errors.New(""))
					}
				}
			}
			newEntry.Attributes = entryReplaceValues(newEntry.Attributes, c.Modification.Type, c.Modification.Vals)
		case ldap.DeleteAttribute:
			log.Printf("applying delete for Attribute: %s", c.Modification.Type)
			for _, rdnAttr := range rdn.Attributes {
				// Modifies on RDN attributes need special care
				if nType == casefold.String(rdnAttr.Type) {
					if len(c.Modification.Vals) == 0 {
						return nil, ldap.NewError(ldap.LDAPResultNotAllowedOnRDN, errors.New(""))
					}
					nRdnVal := casefold.String(rdnAttr.Value)
					for _, delVal := range c.Modification.Vals {
						if nRdnVal == casefold.String(delVal) {
							return nil, ldap.NewError(ldap.LDAPResultNotAllowedOnRDN, errors.New(""))
						}
					}
				}
			}
			newValues := entryApplyModDelete(old.GetEqualFoldAttributeValues(nType), c.Modification.Vals)
			newEntry.Attributes = entryReplaceValues(newEntry.Attributes, c.Modification.Type, newValues)
		}
	}
	return newEntry, nil
}

func entryReplaceValues(ea []*ldap.EntryAttribute, attrType string, newValues []string) (updatedAttrs []*ldap.EntryAttribute) {
	casefold := cases.Fold()
	nType := casefold.String(attrType)
	updated := false
	for _, attr := range ea {
		if casefold.String(attr.Name) == nType {
			updated = true
			if len(newValues) == 0 {
				continue
			}
			updatedAttrs = append(updatedAttrs, ldap.NewEntryAttribute(attr.Name, newValues))
		} else {
			updatedAttrs = append(updatedAttrs, attr)
		}
	}
	if !updated {
		if len(newValues) != 0 {
			updatedAttrs = append(updatedAttrs, ldap.NewEntryAttribute(attrType, newValues))
		}
	}
	return updatedAttrs
}

func entryApplyModAdd(curVals, addVals []string) (newVals []string) {
	newVals = curVals
	casefold := cases.Fold()
	for _, newVal := range addVals {
		present := false
		for _, val := range curVals {
			if casefold.String(newVal) == casefold.String(val) {
				present = true
				break
			}
		}
		if !present {
			newVals = append(newVals, newVal)
		}
	}
	return newVals
}

func entryApplyModDelete(curVals, delVals []string) (newVals []string) {
	casefold := cases.Fold()
	if len(delVals) == 0 {
		return []string{}
	}
	for _, curVal := range curVals {
		nCurVal := casefold.String(curVal)
		keep := true
		for _, del := range delVals {
			if nCurVal == casefold.String(del) {
				keep = false
				break
			}
		}
		if keep {
			newVals = append(newVals, curVal)
		}
	}
	return newVals
}

func EntryFromAddRequest(add *ldap.AddRequest) *ldap.Entry {
	attrs := map[string][]string{}

	for _, a := range add.Attributes {
		attrs[a.Type] = a.Vals
	}
	return ldap.NewEntry(add.DN, attrs)
}

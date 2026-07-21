package identity

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// attrsFromAddRequest flattens the attributes of an *ldap.AddRequest into the
// map[string][]string shape that ldap.NewEntry expects. It is used to synthesize
// the response entry from data already in hand, avoiding a read-after-write when
// oCIS generated the entry's ID itself (useServerUUID=false).
func attrsFromAddRequest(ar *ldap.AddRequest) map[string][]string {
	attrs := make(map[string][]string, len(ar.Attributes))
	for _, a := range ar.Attributes {
		attrs[a.Type] = a.Vals
	}
	return attrs
}

// applyModifyToEntry returns a copy of base with the changes from mr folded on
// (Replace / Add / Delete). It never mutates base. Attribute names are matched
// case-insensitively, consistent with ldap.Entry's GetEqualFold* accessors, so
// folding never produces a duplicate attribute entry that differs only in case.
//
// It is used to synthesize the response entry after an update from the pre-read
// entry plus the ModifyRequest being persisted, avoiding a read-after-write.
func applyModifyToEntry(base *ldap.Entry, mr *ldap.ModifyRequest) *ldap.Entry {
	// Deep-copy base into a name->values map so we never touch the original.
	attrs := make(map[string][]string, len(base.Attributes))
	// order preserves the original attribute order, with new attributes appended.
	order := make([]string, 0, len(base.Attributes))
	index := make(map[string]string, len(base.Attributes)) // fold-key -> stored name
	for _, a := range base.Attributes {
		vals := make([]string, len(a.Values))
		copy(vals, a.Values)
		attrs[a.Name] = vals
		order = append(order, a.Name)
		index[strings.ToLower(a.Name)] = a.Name
	}

	// name resolves the stored attribute name for attrType case-insensitively,
	// registering a new attribute (preserving the request's casing) if none matches.
	name := func(attrType string) string {
		if n, ok := index[strings.ToLower(attrType)]; ok {
			return n
		}
		index[strings.ToLower(attrType)] = attrType
		order = append(order, attrType)
		return attrType
	}

	for _, change := range mr.Changes {
		attrType := change.Modification.Type
		vals := change.Modification.Vals
		switch change.Operation {
		case ldap.ReplaceAttribute:
			n := name(attrType)
			cp := make([]string, len(vals))
			copy(cp, vals)
			attrs[n] = cp
		case ldap.AddAttribute:
			n := name(attrType)
			attrs[n] = append(attrs[n], vals...)
		case ldap.DeleteAttribute:
			if n, ok := index[strings.ToLower(attrType)]; ok {
				if len(vals) == 0 {
					// whole-attribute delete
					delete(attrs, n)
				} else {
					attrs[n] = removeValues(attrs[n], vals)
				}
			}
		}
	}

	result := ldap.NewEntry(base.DN, nil)
	result.Attributes = make([]*ldap.EntryAttribute, 0, len(order))
	for _, n := range order {
		v, ok := attrs[n]
		if !ok {
			// deleted attribute
			continue
		}
		result.Attributes = append(result.Attributes, ldap.NewEntryAttribute(n, v))
	}
	return result
}

// removeValues returns have with every element of remove stripped out.
func removeValues(have, remove []string) []string {
	if len(have) == 0 {
		return have
	}
	drop := make(map[string]struct{}, len(remove))
	for _, r := range remove {
		drop[r] = struct{}{}
	}
	kept := make([]string, 0, len(have))
	for _, v := range have {
		if _, ok := drop[v]; ok {
			continue
		}
		kept = append(kept, v)
	}
	return kept
}

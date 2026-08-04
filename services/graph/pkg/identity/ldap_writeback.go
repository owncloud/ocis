package identity

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// attrsFromAddRequest flattens an *ldap.AddRequest into the map ldap.NewEntry expects,
// synthesizing the response entry without a read-after-write (useServerUUID=false).
func attrsFromAddRequest(ar *ldap.AddRequest) map[string][]string {
	attrs := make(map[string][]string, len(ar.Attributes))
	for _, a := range ar.Attributes {
		vals := make([]string, len(a.Vals))
		copy(vals, a.Vals)
		attrs[a.Type] = vals
	}
	return attrs
}

// applyModifyToEntry synthesizes the post-update entry by folding mr's changes onto a
// copy of base (case-insensitive names), avoiding a read-after-write. base is unchanged.
func applyModifyToEntry(base *ldap.Entry, mr *ldap.ModifyRequest) *ldap.Entry {
	if base == nil {
		return nil
	}
	if mr == nil {
		mr = &ldap.ModifyRequest{DN: base.DN}
	}
	// Deep-copy base into a name->values map; base is never mutated.
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
					kept := removeValues(attrs[n], vals)
					if len(kept) == 0 {
						// the last value was removed; a real server drops the attribute
						delete(attrs, n)
					} else {
						attrs[n] = kept
					}
				}
			}
		}
	}

	dn := base.DN
	if mr.DN != "" {
		dn = mr.DN
	}
	result := ldap.NewEntry(dn, nil)
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

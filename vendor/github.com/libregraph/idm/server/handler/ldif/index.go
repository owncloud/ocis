/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

package ldif

import (
	"strconv"
	"strings"

	"github.com/armon/go-radix"
	"github.com/libregraph/idm/pkg/ldapserver"
	"github.com/spacewander/go-suffix-tree"
)

var indexAttributes = map[string]string{
	"entryCSN":     "eq",
	"entryUUID":    "eq",
	"objectClass":  "eq",
	"cn":           "pres,eq,sub",
	"gidNumber":    "eq",
	"mail":         "pres,eq,sub",
	"memberUid":    "eq",
	"ou":           "eq",
	"uid":          "eq",
	"uidNumber":    "eq",
	"uniqueMember": "eq",

	"sn":        "pres,eq,sub",
	"givenName": "pres,eq,sub",

	"mailAlternateAddress": "eq",

	// Additional indexes for attributes usually used with AD.
	"objectGUID":     "eq",
	"objectSID":      "eq",
	"otherMailbox":   "eq",
	"samAccountName": "eq",
}

func AddIndexAttribute(name string, indices string) {
	indexAttributes[name] = indices
}

func RemoveIndexAttribute(name string) {
	delete(indexAttributes, name)
}

type Index interface {
	Add(name, op string, values []string, entry *ldifEntry) bool
	Load(name, op string, params ...string) ([]*ldifEntry, bool)
}

type indexMap map[string][]*ldifEntry

func newIndexMap() indexMap {
	return make(indexMap)
}

func (im indexMap) Add(name, op string, values []string, entry *ldifEntry) bool {
	for _, value := range values {
		value = strings.ToLower(value)
		im[value] = append(im[value], entry)
	}
	return true
}

func (im indexMap) Load(name, op string, value ...string) ([]*ldifEntry, bool) {
	entries := im[strings.ToLower(value[0])]
	return entries, true
}

type indexSuffixTree struct {
	t *suffix.Tree
}

func newIndexSuffixTree() *indexSuffixTree {
	return &indexSuffixTree{
		t: suffix.NewTree(),
	}
}

func (ist indexSuffixTree) Add(name, op string, values []string, entry *ldifEntry) bool {
	for _, value := range values {
		sfx := []byte(value)
		var entries []*ldifEntry
		if v, ok := ist.t.Get(sfx); ok {
			entries = v.([]*ldifEntry)
		}
		entries = append(entries, entry)
		ist.t.Insert(sfx, entries)
	}
	return true
}

func (ist indexSuffixTree) Load(name, op string, value ...string) ([]*ldifEntry, bool) {
	var entries []*ldifEntry
	sfx := []byte(value[0])
	ist.t.WalkSuffix(sfx, func(key []byte, value interface{}) bool {
		entries = append(entries, value.([]*ldifEntry)...)
		return false
	})
	return entries, true
}

type indexRadixTree struct {
	t *radix.Tree
}

func newIndexRadixTree() *indexRadixTree {
	return &indexRadixTree{
		t: radix.New(),
	}
}

func (irt *indexRadixTree) Add(name, op string, values []string, entry *ldifEntry) bool {
	for _, value := range values {
		pfx := value
		var entries []*ldifEntry
		if v, ok := irt.t.Get(pfx); ok {
			entries = v.([]*ldifEntry)
		}
		entries = append(entries, entry)
		irt.t.Insert(pfx, entries)
	}
	return true
}

func (irt *indexRadixTree) Load(name, op string, value ...string) ([]*ldifEntry, bool) {
	var entries []*ldifEntry
	pfx := value[0]
	irt.t.WalkPrefix(pfx, func(key string, value interface{}) bool {
		entries = append(entries, value.([]*ldifEntry)...)
		return false
	})
	return entries, true
}

type indexSubTree struct {
	pres indexMap
	irt  *indexRadixTree
	ist  *indexSuffixTree
}

func newIndexSubTree() *indexSubTree {
	return &indexSubTree{
		pres: newIndexMap(),
		irt:  newIndexRadixTree(),
		ist:  newIndexSuffixTree(),
	}
}

func (idx *indexSubTree) Add(name, op string, values []string, entry *ldifEntry) bool {
	ok0 := idx.pres.Add(name, op, []string{""}, entry)
	ok1 := idx.irt.Add(name, op, values, entry)
	ok2 := idx.ist.Add(name, op, values, entry)
	return ok0 || ok1 || ok2
}

func (idx *indexSubTree) Load(name, op string, params ...string) ([]*ldifEntry, bool) {
	if len(params) != 2 {
		// Require one value and sub tag.
		return nil, false
	}
	tag, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		panic(err)
	}

	switch tag {
	case ldapserver.FilterSubstringsAny:
		// TODO(longsleep): Find a suitable way for full text search substring
		// matching, for example with Knuth-Morris-Pratt algorithm. Currently
		// we just do a presence match.
		return idx.pres.Load(name, op, "")
	case ldapserver.FilterSubstringsInitial:
		return idx.irt.Load(name, op, params[0])
	case ldapserver.FilterSubstringsFinal:
		return idx.ist.Load(name, op, params[0])
	default:
		return nil, false
	}
}

type indexMapRegister map[string]Index

func newIndexMapRegister() indexMapRegister {
	imr := make(indexMapRegister)
	for name, ops := range indexAttributes {
		for _, op := range strings.Split(ops, ",") {
			switch op {
			case "sub":
				imr[imr.getKey(name, op)] = newIndexSubTree()
			case "pres":
				imr[imr.getKey(name, op)] = newIndexMap()
			case "eq":
				imr[imr.getKey(name, op)] = newIndexMap()
			}
		}
	}
	return imr
}

func (imr indexMapRegister) getKey(name, op string) string {
	return strings.ToLower(name) + "," + op
}

func (imr indexMapRegister) Add(name, op string, values []string, entry *ldifEntry) bool {
	index, ok := imr[imr.getKey(name, op)]
	if !ok {
		// No matching index, refuse to add.
		return false
	}
	return index.Add(name, op, values, entry)
}

func (imr indexMapRegister) Load(name, op string, params ...string) ([]*ldifEntry, bool) {
	index, ok := imr[imr.getKey(name, op)]
	if !ok {
		// No such index.
		return nil, false
	}
	return index.Load(name, op, params...)
}

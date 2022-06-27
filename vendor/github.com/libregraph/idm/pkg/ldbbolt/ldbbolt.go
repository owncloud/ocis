/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2021 The LibreGraph Authors.
 */

// Package ldbbolt provides the lower-level Database functions for managing LDAP Entries
// in a	BoltDB database. Some implementation details:
//
// The database is currently separated in these three buckets
//
// - id2entry: This bucket contains the GOB encoded ldap.Entry instances keyed
//             by a unique 64bit ID
//
// - dn2id: This bucket is used as an index to lookup the ID of an entry by its DN. The DN
//          is used in an normalized (case-folded) form here.
//
// - id2children: This bucket uses the entry-ids as and index and the values contain a list
//                of the entry ids of its direct childdren
//
// Additional buckets will likely be added in the future to create efficient search indexes
package ldbbolt

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"

	"github.com/libregraph/idm/pkg/ldapdn"
	"github.com/libregraph/idm/pkg/ldapentry"
)

type LdbBolt struct {
	logger  logrus.FieldLogger
	db      *bolt.DB
	options *bolt.Options
	base    string
}

var (
	ErrEntryAlreadyExists = errors.New("entry already exists")
	ErrEntryNotFound      = errors.New("entry does not exist")
	ErrNonLeafEntry       = errors.New("entry is not a leaf entry")
)

func (bdb *LdbBolt) Configure(logger logrus.FieldLogger, baseDN, dbfile string, options *bolt.Options) error {
	bdb.logger = logger
	logger = logger.WithField("db", dbfile)
	logger.Debug("Open boltdb")
	db, err := bolt.Open(dbfile, 0o600, options)
	if err != nil {
		logger.WithError(err).Error("Error opening database")
		return err
	}
	bdb.db = db
	bdb.options = options
	bdb.base, _ = ldapdn.ParseNormalize(baseDN)
	return nil
}

// Initialize() opens the Database file and create the required buckets if they do not
// exist yet. After calling initialize the database is ready to process transactions
func (bdb *LdbBolt) Initialize() error {
	var err error
	logger := bdb.logger.WithField("db", bdb.db.Path())
	if bdb.options == nil || !bdb.options.ReadOnly {
		logger.Debug("Adding default buckets")
		err = bdb.db.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists([]byte("dn2id"))
			if err != nil {
				return fmt.Errorf("create bucket 'dn2id': %w", err)
			}
			_, err = tx.CreateBucketIfNotExists([]byte("id2children"))
			if err != nil {
				return fmt.Errorf("create bucket 'dn2id': %w", err)
			}
			_, err = tx.CreateBucketIfNotExists([]byte("id2entry"))
			if err != nil {
				return fmt.Errorf("create bucket 'id2entry': %w", err)
			}
			return nil
		})
		if err != nil {
			logger.WithError(err).Error("Error creating default buckets")
		}
	}
	return err
}

// Performs basic LDAP searches, using the dn2id and id2children buckets to generate
// a list of Result entries. Currently this does strip of the non-request attribute
// Neither does it support LDAP filters. For now we rely on the frontent (LDAPServer)
// to both.
func (bdb *LdbBolt) Search(base string, scope int) ([]*ldap.Entry, error) {
	entries := []*ldap.Entry{}
	nDN, err := ldapdn.ParseNormalize(base)
	if err != nil {
		return entries, err
	}

	err = bdb.db.View(func(tx *bolt.Tx) error {
		entryID := bdb.getIDByDN(tx, nDN)
		var entryIDs []uint64
		if entryID == 0 {
			return fmt.Errorf("not found")
		}
		switch scope {
		case ldap.ScopeBaseObject:
			entryIDs = append(entryIDs, entryID)
		case ldap.ScopeSingleLevel:
			entryIDs = bdb.getChildrenIDs(tx, entryID)
		case ldap.ScopeWholeSubtree:
			entryIDs = append(entryIDs, entryID)
			entryIDs = append(entryIDs, bdb.getSubtreeIDs(tx, entryID)...)
		}
		for _, id := range entryIDs {
			entry, err := bdb.getEntryByID(tx, id)
			if err != nil {
				return err
			}
			entries = append(entries, entry)
		}
		return nil
	})
	return entries, err
}

func idToBytes(id uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, id)
	return b
}

func (bdb *LdbBolt) getChildrenIDs(tx *bolt.Tx, parent uint64) []uint64 {
	id2Children := tx.Bucket([]byte("id2children"))
	children := id2Children.Get(idToBytes(parent))
	r := bytes.NewReader(children)
	ids := make([]uint64, len(children)/8)
	if err := binary.Read(r, binary.LittleEndian, &ids); err != nil {
		bdb.logger.Error(err)
	}
	// This logging it too verbose even for the "debug" level. Leaving
	// it here commented out as it can be helpful during development.
	// bdb.logger.WithFields(logrus.Fields{
	// 	"parentid": parent,
	// 	"children": ids,
	// }).Debug("getChildrenIDs")
	return ids
}

func (bdb *LdbBolt) getSubtreeIDs(tx *bolt.Tx, root uint64) []uint64 {
	var res []uint64
	children := bdb.getChildrenIDs(tx, root)
	res = append(res, children...)
	for _, child := range children {
		res = append(res, bdb.getSubtreeIDs(tx, child)...)
	}
	// This logging it too verbose even for the "debug" level. Leaving
	// it here commented out as it can be helpful during development.
	// bdb.logger.WithFields(logrus.Fields{
	// 	"rootid":  root,
	// 	"subtree": res,
	// }).Debug("getSubtreeIDs")
	return res
}

func (bdb *LdbBolt) EntryPut(e *ldap.Entry) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(e); err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	dn, _ := ldap.ParseDN(e.DN)
	parentDN := &ldap.DN{
		RDNs: dn.RDNs[1:],
	}
	nDN := ldapdn.Normalize(dn)

	if !strings.HasSuffix(nDN, bdb.base) {
		return fmt.Errorf("'%s' is not a descendant of '%s'", e.DN, bdb.base)
	}

	nParentDN := ldapdn.Normalize(parentDN)
	err := bdb.db.Update(func(tx *bolt.Tx) error {
		id2entry := tx.Bucket([]byte("id2entry"))
		id := bdb.getIDByDN(tx, nDN)
		if id != 0 {
			return ErrEntryAlreadyExists
		}
		var err error
		if id, err = id2entry.NextSequence(); err != nil {
			return err
		}

		if err := id2entry.Put(idToBytes(id), buf.Bytes()); err != nil {
			return err
		}
		if nDN != bdb.base {
			if err := bdb.addID2Children(tx, nParentDN, id); err != nil {
				return err
			}
		}
		dn2id := tx.Bucket([]byte("dn2id"))
		if err := dn2id.Put([]byte(nDN), idToBytes(id)); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (bdb *LdbBolt) EntryDelete(dn string) error {
	parsed, err := ldap.ParseDN(dn)
	if err != nil {
		return err
	}
	pparentDN := &ldap.DN{
		RDNs: parsed.RDNs[1:],
	}
	pdn := ldapdn.Normalize(pparentDN)

	ndn := ldapdn.Normalize(parsed)
	err = bdb.db.Update(func(tx *bolt.Tx) error {
		// Does this entry even exist?
		entryID := bdb.getIDByDN(tx, ndn)
		if entryID == 0 {
			return ErrEntryNotFound
		}

		// Refuse to delete if the entry has childs
		id2Children := tx.Bucket([]byte("id2children"))
		children := id2Children.Get(idToBytes(entryID))
		if len(children) != 0 {
			return ErrNonLeafEntry
		}

		// Update id2children bucket (remove entryid from parent)
		parentid := bdb.getIDByDN(tx, pdn)
		if parentid == 0 {
			return ErrEntryNotFound
		}
		children = id2Children.Get(idToBytes(parentid))
		r := bytes.NewReader(children)
		var newids []byte
		idBytes := make([]byte, 8)
		for _, err = io.ReadFull(r, idBytes); err == nil; _, err = io.ReadFull(r, idBytes) {
			if entryID != binary.LittleEndian.Uint64(idBytes) {
				newids = append(newids, idBytes...)
			}
		}
		if err = id2Children.Put(idToBytes(parentid), newids); err != nil {
			return fmt.Errorf("error updating id2Children index for %d: %w", parentid, err)
		}

		// Remove entry from dn2id bucket
		dn2id := tx.Bucket([]byte("dn2id"))
		err = dn2id.Delete([]byte(ndn))
		if err != nil {
			return err
		}
		id2entry := tx.Bucket([]byte("id2entry"))
		err = id2entry.Delete(idToBytes(entryID))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func (bdb *LdbBolt) EntryModify(req *ldap.ModifyRequest) error {
	ndn, err := ldapdn.ParseNormalize(req.DN)
	if err != nil {
		return err
	}
	err = bdb.db.Update(func(tx *bolt.Tx) error {
		oldEntry, id, innerErr := bdb.getEntryByDN(tx, ndn)
		if innerErr != nil {
			return innerErr
		}
		newEntry, innerErr := ldapentry.ApplyModify(oldEntry, req)
		if innerErr != nil {
			return innerErr
		}
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		if innerErr := enc.Encode(newEntry); innerErr != nil {
			return innerErr
		}
		id2entry := tx.Bucket([]byte("id2entry"))
		if innerErr := id2entry.Put(idToBytes(id), buf.Bytes()); innerErr != nil {
			return innerErr
		}

		return nil
	})
	return err
}

func (bdb *LdbBolt) addID2Children(tx *bolt.Tx, nParentDN string, newChildID uint64) error {
	bdb.logger.Debugf("AddID2Children '%s' id '%d'", nParentDN, newChildID)
	parentID := bdb.getIDByDN(tx, nParentDN)
	if parentID == 0 {
		return fmt.Errorf("parent not found '%s'", nParentDN)
	}

	bdb.logger.Debugf("Parent ID: %v", parentID)

	id2Children := tx.Bucket([]byte("id2children"))

	// FIXME add sanity check here if ID is already present
	children := id2Children.Get(idToBytes(parentID))
	children = append(children, idToBytes(newChildID)...)
	if err := id2Children.Put(idToBytes(parentID), children); err != nil {
		return fmt.Errorf("error updating id2Children index for %d: %w", parentID, err)
	}

	bdb.logger.Debugf("AddID2Children '%d' id '%v'", parentID, children)
	return nil
}

func (bdb *LdbBolt) getIDByDN(tx *bolt.Tx, nDN string) uint64 {
	dn2id := tx.Bucket([]byte("dn2id"))
	if dn2id == nil {
		bdb.logger.Debugf("Bucket 'dn2id' does not exist")
		return 0
	}
	id := dn2id.Get([]byte(nDN))
	if id == nil {
		bdb.logger.Debugf("DN: '%s' not found", nDN)
		return 0
	}
	return binary.LittleEndian.Uint64(id)
}

func (bdb *LdbBolt) getEntryByID(tx *bolt.Tx, id uint64) (entry *ldap.Entry, err error) {
	id2entry := tx.Bucket([]byte("id2entry"))
	entrybytes := id2entry.Get(idToBytes(id))
	buf := bytes.NewBuffer(entrybytes)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&entry); err != nil {
		return nil, fmt.Errorf("error decoding entry id: %d, %w", id, err)
	}
	return entry, nil
}

func (bdb *LdbBolt) getEntryByDN(tx *bolt.Tx, ndn string) (entry *ldap.Entry, id uint64, err error) {
	id = bdb.getIDByDN(tx, ndn)
	if id == 0 {
		return nil, id, ErrEntryNotFound
	}
	entry, err = bdb.getEntryByID(tx, id)
	return entry, id, err
}

func (bdb *LdbBolt) Close() {
	bdb.db.Close()
}

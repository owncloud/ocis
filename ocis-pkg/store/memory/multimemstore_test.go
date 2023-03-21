package memory

import (
	"context"
	"strconv"
	"testing"

	"go-micro.dev/v4/store"
)

func TestWriteReadTables(t *testing.T) {
	cache := NewMultiMemStore()

	record1 := &store.Record{
		Key:   "sameKey",
		Value: []byte("from record1"),
	}
	record2 := &store.Record{
		Key:   "sameKey",
		Value: []byte("from record2"),
	}

	_ = cache.Write(record1)
	_ = cache.Write(record2, store.WriteTo("DB02", "Table02"))

	records1, _ := cache.Read("sameKey")
	if len(records1) != 1 {
		t.Fatalf("Wrong number of records, expected 1, got %d", len(records1))
	}
	if records1[0].Key != "sameKey" {
		t.Errorf("Wrong key, expected \"sameKey\", got %s", records1[0].Key)
	}
	if string(records1[0].Value) != "from record1" {
		t.Errorf("Wrong value, expected \"from record1\", got %s", string(records1[0].Value))
	}

	records2, _ := cache.Read("sameKey", store.ReadFrom("DB02", "Table02"))
	if len(records2) != 1 {
		t.Fatalf("Wrong number of records, expected 1, got %d", len(records2))
	}
	if records2[0].Key != "sameKey" {
		t.Errorf("Wrong key, expected \"sameKey\", got %s", records2[0].Key)
	}
	if string(records2[0].Value) != "from record2" {
		t.Errorf("Wrong value, expected \"from record2\", got %s", string(records2[0].Value))
	}
}

func TestDeleteTables(t *testing.T) {
	cache := NewMultiMemStore()

	record1 := &store.Record{
		Key:   "sameKey",
		Value: []byte("from record1"),
	}
	record2 := &store.Record{
		Key:   "sameKey",
		Value: []byte("from record2"),
	}

	_ = cache.Write(record1)
	_ = cache.Write(record2, store.WriteTo("DB02", "Table02"))

	records1, _ := cache.Read("sameKey")
	if len(records1) != 1 {
		t.Fatalf("Wrong number of records, expected 1, got %d", len(records1))
	}
	if records1[0].Key != "sameKey" {
		t.Errorf("Wrong key, expected \"sameKey\", got %s", records1[0].Key)
	}
	if string(records1[0].Value) != "from record1" {
		t.Errorf("Wrong value, expected \"from record1\", got %s", string(records1[0].Value))
	}

	records2, _ := cache.Read("sameKey", store.ReadFrom("DB02", "Table02"))
	if len(records2) != 1 {
		t.Fatalf("Wrong number of records, expected 1, got %d", len(records2))
	}
	if records2[0].Key != "sameKey" {
		t.Errorf("Wrong key, expected \"sameKey\", got %s", records2[0].Key)
	}
	if string(records2[0].Value) != "from record2" {
		t.Errorf("Wrong value, expected \"from record2\", got %s", string(records2[0].Value))
	}

	_ = cache.Delete("sameKey")
	if _, err := cache.Read("sameKey"); err != store.ErrNotFound {
		t.Errorf("Key \"sameKey\" still exists after deletion")
	}

	records2, _ = cache.Read("sameKey", store.ReadFrom("DB02", "Table02"))
	if len(records2) != 1 {
		t.Fatalf("Wrong number of records, expected 1, got %d", len(records2))
	}
	if records2[0].Key != "sameKey" {
		t.Errorf("Wrong key, expected \"sameKey\", got %s", records2[0].Key)
	}
	if string(records2[0].Value) != "from record2" {
		t.Errorf("Wrong value, expected \"from record2\", got %s", string(records2[0].Value))
	}
}

func TestListTables(t *testing.T) {
	cache := NewMultiMemStore()

	record1 := &store.Record{
		Key:   "key001",
		Value: []byte("from record1"),
	}
	record2 := &store.Record{
		Key:   "key002",
		Value: []byte("from record2"),
	}

	_ = cache.Write(record1)
	_ = cache.Write(record2, store.WriteTo("DB02", "Table02"))

	keys, _ := cache.List(store.ListFrom("DB02", "Table02"))
	expectedKeys := []string{"key002"}
	if len(keys) != 1 {
		t.Fatalf("Wrong number of keys, expected 1, got %d", len(keys))
	}
	for index, key := range keys {
		if expectedKeys[index] != key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func TestWriteSizeLimit(t *testing.T) {
	cache := NewMultiMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 2,
				},
			),
		),
	)

	record := &store.Record{}
	for i := 0; i < 4; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		_ = cache.Write(record)
		_ = cache.Write(record, store.WriteTo("DB02", "Table02"))
	}

	keys1, _ := cache.List()
	expectedKeys1 := []string{"2", "3"}
	if len(keys1) != 2 {
		t.Fatalf("Wrong number of keys, expected 2, got %d", len(keys1))
	}
	for index, key := range keys1 {
		if expectedKeys1[index] != key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKeys1[index], key)
		}
	}

	keys2, _ := cache.List(store.ListFrom("DB02", "Table02"))
	expectedKeys2 := []string{"2", "3"}
	if len(keys2) != 2 {
		t.Fatalf("Wrong number of keys, expected 2, got %d", len(keys2))
	}
	for index, key := range keys2 {
		if expectedKeys2[index] != key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKeys2[index], key)
		}
	}
}

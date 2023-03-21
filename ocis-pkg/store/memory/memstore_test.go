package memory

import (
	"context"
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"sync/atomic"

	"go-micro.dev/v4/store"
)

func TestWriteAndRead(t *testing.T) {
	cache := NewMemStore()
	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzzz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"zozzz":      "vooouyenbdnya",
		"zzaz":       "viaooouyenbdnya",
		"mbzzaamb":   "viunya",
	}

	for key, value := range data {
		record := &store.Record{
			Key:   key,
			Value: []byte(value),
		}
		_ = cache.Write(record)
	}

	t.Run("Plain", func(t *testing.T) {
		readPlain(t, cache)
	})
	t.Run("Prefix", func(t *testing.T) {
		readPrefix(t, cache)
	})
	t.Run("Suffix", func(t *testing.T) {
		readSuffix(t, cache)
	})
	t.Run("PrefixSuffix", func(t *testing.T) {
		readPrefixSuffix(t, cache)
	})
	t.Run("PrefixLimitOffset", func(t *testing.T) {
		readPrefixLimitOffset(t, cache)
	})
	t.Run("SuffixLimitOffset", func(t *testing.T) {
		readSuffixLimitOffset(t, cache)
	})
	t.Run("PrefixSuffixLimitOffset", func(t *testing.T) {
		readPrefixSuffixLimitOffset(t, cache)
	})
}

func readPlain(t *testing.T, cache store.Store) {
	// expected data in the cache
	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzzz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"zozzz":      "vooouyenbdnya",
		"zzaz":       "viaooouyenbdnya",
		"mbzzaamb":   "viunya",
	}
	for key, value := range data {
		records, _ := cache.Read(key)
		if len(records) != 1 {
			t.Fatalf("Plain read for key %s returned %d records", key, len(records))
		}
		if key != records[0].Key {
			t.Errorf("Plain read for key %s returned got wrong key %s", key, records[0].Key)
		}
		v := string(records[0].Value)
		if value != v {
			t.Errorf("Plain read for key %s returned different value, expected %s, got %s", key, value, v)
		}
	}
}

func readPrefix(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "abaya", Value: "v329487"},
		{Key: "abayakjdkj", Value: "v989898"},
	}

	pref2 := []struct {
		Key   string
		Value string
	}{
		{Key: "zozzz", Value: "vooouyenbdnya"},
		{Key: "zzaz", Value: "viaooouyenbdnya"},
		{Key: "zzzz", Value: "viaooouyenbdnya"},
	}

	records, _ := cache.Read("abaya", store.ReadPrefix())
	if len(records) != 2 {
		t.Fatalf("Prefix read for \"abaya\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for prefix \"abaya\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for prefix \"abaya\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}

	records, _ = cache.Read("z", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Prefix read for \"z\" returned %d records, expected 3", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref2[index].Key != record.Key {
			t.Errorf("Unexpected key for prefix \"z\", index %d, expected %s, got %s", index, pref2[index].Key, record.Key)
		}
		if pref2[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for prefix \"z\", index %d, expected %s, got %s", index, pref2[index].Value, record.Value)
		}
	}
}

func readSuffix(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "abaaz", Value: "v398342"},
		{Key: "zzaz", Value: "viaooouyenbdnya"},
	}
	pref2 := []struct {
		Key   string
		Value string
	}{
		{Key: "abaaz", Value: "v398342"},
		{Key: "zzaz", Value: "viaooouyenbdnya"},
		{Key: "abazzz", Value: "v57869nbdnya"},
		{Key: "zozzz", Value: "vooouyenbdnya"},
		{Key: "zzzz", Value: "viaooouyenbdnya"},
	}

	records, _ := cache.Read("az", store.ReadSuffix())
	if len(records) != 2 {
		t.Fatalf("Suffix read for \"az\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for suffix \"az\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for suffix \"az\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}

	records, _ = cache.Read("z", store.ReadSuffix())
	if len(records) != 5 {
		t.Fatalf("Suffix read for \"z\" returned %d records, expected 5", len(records))
	}
	for index, record := range records {
		if pref2[index].Key != record.Key {
			t.Errorf("Unexpected key for suffix \"z\", index %d, expected %s, got %s", index, pref2[index].Key, record.Key)
		}
		if pref2[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for suffix \"z\", index %d, expected %s, got %s", index, pref2[index].Value, record.Value)
		}
	}
}

func readPrefixSuffix(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "zozzz", Value: "vooouyenbdnya"},
		{Key: "zzaz", Value: "viaooouyenbdnya"},
		{Key: "zzzz", Value: "viaooouyenbdnya"},
	}
	pref2 := []struct {
		Key   string
		Value string
	}{
		{Key: "mbmbmb", Value: "viuyenbdnya"},
		{Key: "mbzzaamb", Value: "viunya"},
	}

	records, _ := cache.Read("z", store.ReadPrefix(), store.ReadSuffix())
	if len(records) != 3 {
		t.Fatalf("Prefix-Suffix read for \"z\" returned %d records, expected 3", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for prefix-suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for prefix-suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}

	records, _ = cache.Read("mb", store.ReadPrefix(), store.ReadSuffix())
	if len(records) != 2 {
		t.Fatalf("Prefix-Suffix read for \"mb\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref2[index].Key != record.Key {
			t.Errorf("Unexpected key for prefix-suffix \"mb\", index %d, expected %s, got %s", index, pref2[index].Key, record.Key)
		}
		if pref2[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for prefix-suffix \"mb\", index %d, expected %s, got %s", index, pref2[index].Value, record.Value)
		}
	}
}

func readPrefixLimitOffset(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "abaaz", Value: "v398342"},
		{Key: "abaya", Value: "v329487"},
	}
	pref2 := []struct {
		Key   string
		Value string
	}{
		{Key: "abayakjdkj", Value: "v989898"},
		{Key: "abazzz", Value: "v57869nbdnya"},
	}

	records, _ := cache.Read("aba", store.ReadPrefix(), store.ReadLimit(2))
	if len(records) != 2 {
		t.Fatalf("Limit prefix read for \"aba\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for limit prefix \"aba\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for limit prefix \"aba\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}

	records, _ = cache.Read("aba", store.ReadPrefix(), store.ReadLimit(2), store.ReadOffset(2))
	if len(records) != 2 {
		t.Fatalf("Offset-limit prefix read for \"aba\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref2[index].Key != record.Key {
			t.Errorf("Unexpected key for offset-limit prefix \"aba\", index %d, expected %s, got %s", index, pref2[index].Key, record.Key)
		}
		if pref2[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for offset-limit prefix \"aba\", index %d, expected %s, got %s", index, pref2[index].Value, record.Value)
		}
	}
}

func readSuffixLimitOffset(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "abaaz", Value: "v398342"},
		{Key: "zzaz", Value: "viaooouyenbdnya"},
	}
	pref2 := []struct {
		Key   string
		Value string
	}{
		{Key: "abazzz", Value: "v57869nbdnya"},
		{Key: "zozzz", Value: "vooouyenbdnya"},
	}

	records, _ := cache.Read("z", store.ReadSuffix(), store.ReadLimit(2))
	if len(records) != 2 {
		t.Fatalf("Limit suffix read for \"z\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for limit suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for limit suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}

	records, _ = cache.Read("z", store.ReadSuffix(), store.ReadLimit(2), store.ReadOffset(2))
	if len(records) != 2 {
		t.Fatalf("Offset-limit suffix read for \"z\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref2[index].Key != record.Key {
			t.Errorf("Unexpected key for offset-limit suffix \"z\", index %d, expected %s, got %s", index, pref2[index].Key, record.Key)
		}
		if pref2[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for offset-limit suffix \"z\", index %d, expected %s, got %s", index, pref2[index].Value, record.Value)
		}
	}
}

func readPrefixSuffixLimitOffset(t *testing.T, cache store.Store) {
	pref1 := []struct {
		Key   string
		Value string
	}{
		{Key: "zzaz", Value: "viaooouyenbdnya"},
		{Key: "zzzz", Value: "viaooouyenbdnya"},
	}

	records, _ := cache.Read("z", store.ReadPrefix(), store.ReadSuffix(), store.ReadOffset(1), store.ReadLimit(2))
	if len(records) != 2 {
		t.Fatalf("Limit suffix read for \"z\" returned %d records, expected 2", len(records))
	}
	for index, record := range records {
		// it should be sorted alphabetically
		if pref1[index].Key != record.Key {
			t.Errorf("Unexpected key for limit suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Key, record.Key)
		}
		if pref1[index].Value != string(record.Value) {
			t.Errorf("Unexpected value for limit suffix \"z\", index %d, expected %s, got %s", index, pref1[index].Value, record.Value)
		}
	}
}

func TestWriteExpiryAndRead(t *testing.T) {
	cache := NewMemStore()

	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzaz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"mbzzaamb":   "viunya",
		"zozzz":      "vooouyenbdnya",
	}

	for key, value := range data {
		record := &store.Record{
			Key:    key,
			Value:  []byte(value),
			Expiry: time.Second * 1000,
		}
		_ = cache.Write(record)
	}

	records, _ := cache.Read("zzaz")
	if len(records) != 1 {
		t.Fatalf("Failed read for \"zzaz\" returned %d records, expected 1", len(records))
	}
	record := records[0]
	if record.Expiry < 999*time.Second || record.Expiry > 1000*time.Second {
		// The expiry will be adjusted on retrieval
		t.Errorf("Abnormal expiry range: expected %d-%d, got %d", 999*time.Second, 1000*time.Second, record.Expiry)
	}
}

func TestWriteExpiryWithExpiryAndRead(t *testing.T) {
	cache := NewMemStore()

	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzaz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"mbzzaamb":   "viunya",
		"zozzz":      "vooouyenbdnya",
	}

	for key, value := range data {
		record := &store.Record{
			Key:    key,
			Value:  []byte(value),
			Expiry: time.Second * 1000,
		}
		// write option will override the record data
		_ = cache.Write(record, store.WriteExpiry(time.Now().Add(time.Hour)))
	}

	records, _ := cache.Read("zzaz")
	if len(records) != 1 {
		t.Fatalf("Failed read for \"zzaz\" returned %d records, expected 1", len(records))
	}
	record := records[0]
	if record.Expiry < 3599*time.Second || record.Expiry > 3600*time.Second {
		// The expiry will be adjusted on retrieval
		t.Errorf("Abnormal expiry range: expected %d-%d, got %d", 3599*time.Second, 3600*time.Second, record.Expiry)
	}
}

func TestWriteExpiryWithTTLAndRead(t *testing.T) {
	cache := NewMemStore()

	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzaz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"mbzzaamb":   "viunya",
		"zozzz":      "vooouyenbdnya",
	}

	for key, value := range data {
		record := &store.Record{
			Key:    key,
			Value:  []byte(value),
			Expiry: time.Second * 1000,
		}
		// write option will override the record data, TTL takes precedence
		_ = cache.Write(record, store.WriteTTL(20*time.Second), store.WriteExpiry(time.Now().Add(time.Hour)))
	}

	records, _ := cache.Read("zzaz")
	if len(records) != 1 {
		t.Fatalf("Failed read for \"zzaz\" returned %d records, expected 1", len(records))
	}
	record := records[0]
	if record.Expiry < 19*time.Second || record.Expiry > 20*time.Second {
		// The expiry will be adjusted on retrieval
		t.Errorf("Abnormal expiry range: expected %d-%d, got %d", 19*time.Second, 20*time.Second, record.Expiry)
	}
}

func TestDelete(t *testing.T) {
	cache := NewMemStore()
	record := &store.Record{
		Key:   "record",
		Value: []byte("value for record"),
	}

	records, err := cache.Read("record")
	if err != store.ErrNotFound && len(records) > 0 {
		t.Fatal("Found key in cache but it shouldn't be there")
	}

	_ = cache.Write(record)
	records, err = cache.Read("record")
	if err != nil {
		t.Fatal("Key not found in cache after inserting it")
	}
	if len(records) != 1 {
		t.Fatal("Multiple keys found in cache after inserting it")
	}
	if records[0].Key != "record" && string(records[0].Value) != "value for record" {
		t.Fatal("Wrong record retrieved")
	}

	err = cache.Delete("record")
	if err != nil {
		t.Fatal("Error deleting the record")
	}

	records, err = cache.Read("record")
	if err != store.ErrNotFound && len(records) > 0 {
		t.Fatal("Found key in cache but it shouldn't be there")
	}
}

func TestList(t *testing.T) {
	cache := NewMemStore()
	data := map[string]string{
		"abaya":      "v329487",
		"abaaz":      "v398342",
		"abayakjdkj": "v989898",
		"zzzz":       "viaooouyenbdnya",
		"abazzz":     "v57869nbdnya",
		"mbmbmb":     "viuyenbdnya",
		"zozzz":      "vooouyenbdnya",
		"aboyo":      "v889487",
		"zzaaaz":     "v999487",
	}

	for key, value := range data {
		record := &store.Record{
			Key:   key,
			Value: []byte(value),
		}
		_ = cache.Write(record)
	}

	t.Run("Plain", func(t *testing.T) {
		listPlain(t, cache)
	})
	t.Run("Prefix", func(t *testing.T) {
		listPrefix(t, cache)
	})
	t.Run("Suffix", func(t *testing.T) {
		listSuffix(t, cache)
	})
	t.Run("PrefixSuffix", func(t *testing.T) {
		listPrefixSuffix(t, cache)
	})
	t.Run("LimitOffset", func(t *testing.T) {
		listLimitOffset(t, cache)
	})
	t.Run("PrefixLimitOffset", func(t *testing.T) {
		listPrefixLimitOffset(t, cache)
	})
	t.Run("SuffixLimitOffset", func(t *testing.T) {
		listSuffixLimitOffset(t, cache)
	})
	t.Run("PrefixSuffixLimitOffset", func(t *testing.T) {
		listPrefixSuffixLimitOffset(t, cache)
	})
}

func listPlain(t *testing.T, cache store.Store) {
	keys, _ := cache.List()
	expectedKeys := []string{"abaaz", "abaya", "abayakjdkj", "abazzz", "aboyo", "mbmbmb", "zozzz", "zzaaaz", "zzzz"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listPrefix(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListPrefix("aba"))
	expectedKeys := []string{"abaaz", "abaya", "abayakjdkj", "abazzz"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listSuffix(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListSuffix("z"))
	expectedKeys := []string{"zzaaaz", "abaaz", "abazzz", "zozzz", "zzzz"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listPrefixSuffix(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListPrefix("ab"), store.ListSuffix("z"))
	expectedKeys := []string{"abaaz", "abazzz"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listLimitOffset(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListLimit(3), store.ListOffset(2))
	expectedKeys := []string{"abayakjdkj", "abazzz", "aboyo"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listPrefixLimitOffset(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListPrefix("aba"), store.ListLimit(2), store.ListOffset(1))
	expectedKeys := []string{"abaya", "abayakjdkj"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listSuffixLimitOffset(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListSuffix("z"), store.ListLimit(2), store.ListOffset(1))
	expectedKeys := []string{"abaaz", "abazzz"}
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func listPrefixSuffixLimitOffset(t *testing.T, cache store.Store) {
	keys, _ := cache.List(store.ListPrefix("a"), store.ListSuffix("z"), store.ListLimit(2), store.ListOffset(1))
	expectedKeys := []string{"abazzz"} // only 2 available, and we skip the first one
	if len(keys) != len(expectedKeys) {
		t.Fatalf("Wrong number of keys, expected %d, got %d", len(expectedKeys), len(keys))
	}

	for index, key := range keys {
		if key != expectedKeys[index] {
			t.Errorf("Wrong key in the list in index %d, expected %s, got %s", index, expectedKeys[index], key)
		}
	}
}

func TestEvictWriteUpdate(t *testing.T) {
	cache := NewMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 3,
				},
			)),
	)

	for i := 0; i < 3; i++ {
		v := strconv.Itoa(i)
		record := &store.Record{
			Key:   v,
			Value: []byte(v),
		}
		_ = cache.Write(record)
	}

	// update first item
	updatedRecord := &store.Record{
		Key:   "0",
		Value: []byte("zero"),
	}
	_ = cache.Write(updatedRecord)

	// new record, to force eviction
	newRecord := &store.Record{
		Key:   "new",
		Value: []byte("newNew"),
	}
	_ = cache.Write(newRecord)

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Wrong number of record returned, expected 3, got %d", len(records))
	}

	expectedKV := []struct {
		Key   string
		Value string
	}{
		{Key: "0", Value: "zero"},
		{Key: "2", Value: "2"},
		{Key: "new", Value: "newNew"},
	}

	for index, record := range records {
		if record.Key != expectedKV[index].Key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKV[index].Key, record.Key)
		}
		if string(record.Value) != expectedKV[index].Value {
			t.Errorf("Wrong value  for index %d, expected %s, got %s", index, expectedKV[index].Value, string(record.Value))
		}
	}
}

func TestEvictRead(t *testing.T) {
	cache := NewMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 3,
				},
			)),
	)

	for i := 0; i < 3; i++ {
		v := strconv.Itoa(i)
		record := &store.Record{
			Key:   v,
			Value: []byte(v),
		}
		_ = cache.Write(record)
	}

	// Read first item
	_, _ = cache.Read("0")

	// new record, to force eviction
	newRecord := &store.Record{
		Key:   "new",
		Value: []byte("newNew"),
	}
	_ = cache.Write(newRecord)

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Wrong number of record returned, expected 3, got %d", len(records))
	}

	expectedKV := []struct {
		Key   string
		Value string
	}{
		{Key: "0", Value: "0"},
		{Key: "2", Value: "2"},
		{Key: "new", Value: "newNew"},
	}

	for index, record := range records {
		if record.Key != expectedKV[index].Key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKV[index].Key, record.Key)
		}
		if string(record.Value) != expectedKV[index].Value {
			t.Errorf("Wrong value  for index %d, expected %s, got %s", index, expectedKV[index].Value, string(record.Value))
		}
	}
}

func TestEvictReadPrefix(t *testing.T) {
	cache := NewMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 3,
				},
			)),
	)

	for i := 0; i < 3; i++ {
		v := strconv.Itoa(i)
		record := &store.Record{
			Key:   v,
			Value: []byte(v),
		}
		_ = cache.Write(record)
	}

	// Read prefix won't change evcition list
	_, _ = cache.Read("0", store.ReadPrefix())

	// new record, to force eviction
	newRecord := &store.Record{
		Key:   "new",
		Value: []byte("newNew"),
	}
	_ = cache.Write(newRecord)

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Wrong number of record returned, expected 3, got %d", len(records))
	}

	expectedKV := []struct {
		Key   string
		Value string
	}{
		{Key: "1", Value: "1"},
		{Key: "2", Value: "2"},
		{Key: "new", Value: "newNew"},
	}

	for index, record := range records {
		if record.Key != expectedKV[index].Key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKV[index].Key, record.Key)
		}
		if string(record.Value) != expectedKV[index].Value {
			t.Errorf("Wrong value  for index %d, expected %s, got %s", index, expectedKV[index].Value, string(record.Value))
		}
	}
}

func TestEvictReadSuffix(t *testing.T) {
	cache := NewMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 3,
				},
			)),
	)

	for i := 0; i < 3; i++ {
		v := strconv.Itoa(i)
		record := &store.Record{
			Key:   v,
			Value: []byte(v),
		}
		_ = cache.Write(record)
	}

	// Read suffix won't change evcition list
	_, _ = cache.Read("0", store.ReadSuffix())

	// new record, to force eviction
	newRecord := &store.Record{
		Key:   "new",
		Value: []byte("newNew"),
	}
	_ = cache.Write(newRecord)

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Wrong number of record returned, expected 3, got %d", len(records))
	}

	expectedKV := []struct {
		Key   string
		Value string
	}{
		{Key: "1", Value: "1"},
		{Key: "2", Value: "2"},
		{Key: "new", Value: "newNew"},
	}

	for index, record := range records {
		if record.Key != expectedKV[index].Key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKV[index].Key, record.Key)
		}
		if string(record.Value) != expectedKV[index].Value {
			t.Errorf("Wrong value  for index %d, expected %s, got %s", index, expectedKV[index].Value, string(record.Value))
		}
	}
}

func TestEvictList(t *testing.T) {
	cache := NewMemStore(
		store.WithContext(
			NewContext(
				context.Background(),
				map[string]interface{}{
					"maxCap": 3,
				},
			)),
	)

	for i := 0; i < 3; i++ {
		v := strconv.Itoa(i)
		record := &store.Record{
			Key:   v,
			Value: []byte(v),
		}
		_ = cache.Write(record)
	}

	// List won't change evcition list
	_, _ = cache.List()

	// new record, to force eviction
	newRecord := &store.Record{
		Key:   "new",
		Value: []byte("newNew"),
	}
	_ = cache.Write(newRecord)

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 3 {
		t.Fatalf("Wrong number of record returned, expected 3, got %d", len(records))
	}

	expectedKV := []struct {
		Key   string
		Value string
	}{
		{Key: "1", Value: "1"},
		{Key: "2", Value: "2"},
		{Key: "new", Value: "newNew"},
	}

	for index, record := range records {
		if record.Key != expectedKV[index].Key {
			t.Errorf("Wrong key for index %d, expected %s, got %s", index, expectedKV[index].Key, record.Key)
		}
		if string(record.Value) != expectedKV[index].Value {
			t.Errorf("Wrong value  for index %d, expected %s, got %s", index, expectedKV[index].Value, string(record.Value))
		}
	}
}

func TestExpireReadPrefix(t *testing.T) {
	cache := NewMemStore()

	record := &store.Record{}
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		if i%2 == 0 {
			record.Expiry = time.Duration(i) * time.Minute
		} else {
			record.Expiry = time.Duration(-i) * time.Minute
		}
		_ = cache.Write(record)
	}

	records, _ := cache.Read("", store.ReadPrefix())
	if len(records) != 10 {
		t.Fatalf("Wrong number of records, expected 10, got %d", len(records))
	}

	var expKeys []string
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expKeys = append(expKeys, strconv.Itoa(i))
		}
	}
	sort.Strings(expKeys)

	expKeyIndex := 0
	for _, record := range records {
		if record.Key != expKeys[expKeyIndex] {
			t.Fatalf("Wrong expected key, expected %s, got %s", expKeys[expKeyIndex], record.Key)
		}
		expKeyIndex++
	}
}

func TestExpireReadSuffix(t *testing.T) {
	cache := NewMemStore()

	record := &store.Record{}
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		if i%2 == 0 {
			record.Expiry = time.Duration(i) * time.Minute
		} else {
			record.Expiry = time.Duration(-i) * time.Minute
		}
		_ = cache.Write(record)
	}

	records, _ := cache.Read("", store.ReadSuffix())
	if len(records) != 10 {
		t.Fatalf("Wrong number of records, expected 10, got %d", len(records))
	}

	var expKeys []string
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expKeys = append(expKeys, strconv.Itoa(i))
		}
	}
	sort.Slice(expKeys, func(i, j int) bool {
		return reverseString(expKeys[i]) < reverseString(expKeys[j])
	})

	expKeyIndex := 0
	for _, record := range records {
		if record.Key != expKeys[expKeyIndex] {
			t.Fatalf("Wrong expected key, expected %s, got %s", expKeys[expKeyIndex], record.Key)
		}
		expKeyIndex++
	}
}

func TestExpireList(t *testing.T) {
	cache := NewMemStore()

	record := &store.Record{}
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		if i%2 == 0 {
			record.Expiry = time.Duration(i) * time.Minute
		} else {
			record.Expiry = time.Duration(-i) * time.Minute
		}
		_ = cache.Write(record)
	}

	keys, _ := cache.List()
	if len(keys) != 10 {
		t.Fatalf("Wrong number of records, expected 10, got %d", len(keys))
	}

	var expKeys []string
	for i := 0; i < 20; i++ {
		if i%2 == 0 {
			expKeys = append(expKeys, strconv.Itoa(i))
		}
	}
	sort.Strings(expKeys)

	expKeyIndex := 0
	for _, key := range keys {
		if key != expKeys[expKeyIndex] {
			t.Fatalf("Wrong expected key, expected %s, got %s", expKeys[expKeyIndex], key)
		}
		expKeyIndex++
	}
}

func TestExpireListPrefix(t *testing.T) {
	cache := NewMemStore()

	record := &store.Record{}
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		if i%2 == 0 {
			record.Expiry = time.Duration(i) * time.Minute
		} else {
			record.Expiry = time.Duration(-i) * time.Minute
		}
		_ = cache.Write(record)
	}

	keys, _ := cache.List(store.ListPrefix("1"))
	if len(keys) != 5 {
		t.Fatalf("Wrong number of records, expected 5, got %d", len(keys))
	}

	var expKeys []string
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		if i%2 == 0 && strings.HasPrefix(v, "1") {
			expKeys = append(expKeys, v)
		}
	}
	sort.Strings(expKeys)

	expKeyIndex := 0
	for _, key := range keys {
		if key != expKeys[expKeyIndex] {
			t.Fatalf("Wrong expected key, expected %s, got %s", expKeys[expKeyIndex], key)
		}
		expKeyIndex++
	}
}

func TestExpireListSuffix(t *testing.T) {
	cache := NewMemStore()

	record := &store.Record{}
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		record.Key = v
		record.Value = []byte(v)
		if i%2 == 0 {
			record.Expiry = time.Duration(i) * time.Minute
		} else {
			record.Expiry = time.Duration(-i) * time.Minute
		}
		_ = cache.Write(record)
	}

	keys, _ := cache.List(store.ListSuffix("8"))
	if len(keys) != 2 {
		t.Fatalf("Wrong number of records, expected 2, got %d", len(keys))
	}

	var expKeys []string
	for i := 0; i < 20; i++ {
		v := strconv.Itoa(i)
		if i%2 == 0 && strings.HasSuffix(v, "8") {
			expKeys = append(expKeys, v)
		}
	}
	sort.Slice(expKeys, func(i, j int) bool {
		return reverseString(expKeys[i]) < reverseString(expKeys[j])
	})

	expKeyIndex := 0
	for _, key := range keys {
		if key != expKeys[expKeyIndex] {
			t.Fatalf("Wrong expected key, expected %s, got %s", expKeys[expKeyIndex], key)
		}
		expKeyIndex++
	}
}

func TestConcurrentWrite(t *testing.T) {
	nThreads := []int{3, 10, 50}

	for _, threads := range nThreads {
		t.Run("T"+strconv.Itoa(threads), func(t *testing.T) {
			cache := NewMemStore(
				store.WithContext(
					NewContext(
						context.Background(),
						map[string]interface{}{
							"maxCap": 50000,
						},
					)),
			)

			var wg sync.WaitGroup
			var index int64

			wg.Add(threads)
			for i := 0; i < threads; i++ {
				go func(cache store.Store, ind *int64) {
					j := atomic.AddInt64(ind, 1) - 1
					for j < 100000 {
						v := strconv.FormatInt(j, 10)
						record := &store.Record{
							Key:   v,
							Value: []byte(v),
						}
						_ = cache.Write(record)
						j = atomic.AddInt64(ind, 1) - 1
					}
					wg.Done()
				}(cache, &index)
			}
			wg.Wait()

			records, _ := cache.Read("", store.ReadPrefix())
			if len(records) != 50000 {
				t.Fatalf("Wrong number of records, expected 50000, got %d", len(records))
			}
			for _, record := range records {
				if record.Key != string(record.Value) {
					t.Fatalf("Wrong record found, key %s, value %s", record.Key, string(record.Value))
				}
			}
		})
	}
}

func BenchmarkWrite(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// records will be copied, so it's safe to overwrite the previous record
				v := strconv.Itoa(i)
				record.Key = v
				record.Value = []byte(v)
				_ = cache.Write(record)
			}
		})
	}
}

func BenchmarkRead(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		for i := 0; i < size; i++ {
			// records will be copied, so it's safe to overwrite the previous record
			v := strconv.Itoa(i)
			record.Key = v
			record.Value = []byte(v)
			_ = cache.Write(record)
		}
		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v := strconv.Itoa(i)
				_, _ = cache.Read(v)
			}
		})
	}
}

func BenchmarkWriteMedKey(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	h := fnv.New128()
	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Reset()
				v := strconv.Itoa(i)
				bys := []byte(v)
				h.Write(bys)
				// records will be copied, so it's safe to overwrite the previous record
				record.Key = hex.EncodeToString(h.Sum(nil))
				record.Value = bys
				_ = cache.Write(record)
			}
		})
	}
}

func BenchmarkReadMedKey(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	h := fnv.New128()
	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		for i := 0; i < size; i++ {
			h.Reset()
			v := strconv.Itoa(i)
			bys := []byte(v)
			h.Write(bys)
			// records will be copied, so it's safe to overwrite the previous record
			record.Key = hex.EncodeToString(h.Sum(nil))
			record.Value = bys
			_ = cache.Write(record)
		}
		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Reset()
				v := strconv.Itoa(i)
				bys := []byte(v)
				h.Write(bys)
				_, _ = cache.Read(hex.EncodeToString(h.Sum(nil)))
			}
		})
	}
}

func BenchmarkReadMedKeyPrefix(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	h := fnv.New128()
	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		for i := 0; i < size; i++ {
			h.Reset()
			v := strconv.Itoa(i)
			bys := []byte(v)
			h.Write(bys)
			// records will be copied, so it's safe to overwrite the previous record
			record.Key = hex.EncodeToString(h.Sum(nil))
			record.Value = bys
			_ = cache.Write(record)
		}
		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Reset()
				v := strconv.Itoa(i)
				bys := []byte(v)
				h.Write(bys)
				_, _ = cache.Read(hex.EncodeToString(h.Sum(nil))[:10], store.ReadPrefix(), store.ReadLimit(50))
			}
		})
	}
}

func BenchmarkReadMedKeySuffix(b *testing.B) {
	cacheSizes := []int{512, 1024, 10000, 50000, 1000000}

	h := fnv.New128()
	for _, size := range cacheSizes {
		cache := NewMemStore(
			store.WithContext(
				NewContext(
					context.Background(),
					map[string]interface{}{
						"maxCap": size,
					},
				)),
		)
		record := &store.Record{}

		for i := 0; i < size; i++ {
			h.Reset()
			v := strconv.Itoa(i)
			bys := []byte(v)
			h.Write(bys)
			// records will be copied, so it's safe to overwrite the previous record
			record.Key = hex.EncodeToString(h.Sum(nil))
			record.Value = bys
			_ = cache.Write(record)
		}
		b.Run("CacheSize"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.Reset()
				v := strconv.Itoa(i)
				bys := []byte(v)
				h.Write(bys)
				_, _ = cache.Read(hex.EncodeToString(h.Sum(nil))[23:], store.ReadSuffix(), store.ReadLimit(50))
			}
		})
	}
}

func concurrentStoreBench(b *testing.B, threads int) {
	benchTest := map[string]int{
		"DefCap": 512,
		"LimCap": 3,
		"BigCap": 1000000,
	}
	for testname, size := range benchTest {
		b.Run(testname, func(b *testing.B) {
			cache := NewMemStore(
				store.WithContext(
					NewContext(
						context.Background(),
						map[string]interface{}{
							"maxCap": size,
						},
					)),
			)

			b.SetParallelism(threads)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				h := fnv.New128()
				record := &store.Record{}
				for pb.Next() {
					h.Reset()
					v := strconv.Itoa(rand.Int()) //nolint:gosec
					bys := []byte(v)
					h.Write(bys)
					// records will be copied, so it's safe to overwrite the previous record
					record.Key = hex.EncodeToString(h.Sum(nil))
					record.Value = bys
					_ = cache.Write(record)
				}
			})
		})
	}
}

func concurrentRetrieveBench(b *testing.B, threads int) {
	benchTest := map[string]int{
		"DefCap": 512,
		"LimCap": 3,
		"BigCap": 1000000,
	}
	for testname, size := range benchTest {
		b.Run(testname, func(b *testing.B) {
			cache := NewMemStore(
				store.WithContext(
					NewContext(
						context.Background(),
						map[string]interface{}{
							"maxCap": size,
						},
					)),
			)

			record := &store.Record{}
			for i := 0; i < size; i++ {
				v := strconv.Itoa(i)
				record.Key = v
				record.Value = []byte(v)
				_ = cache.Write(record)
			}

			b.SetParallelism(threads)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					v := strconv.Itoa(rand.Intn(size * 2)) //nolint:gosec
					_, _ = cache.Read(v)
				}
			})
		})
	}
}

func concurrentRemoveBench(b *testing.B, threads int) {
	benchTest := map[string]int{
		"DefCap": 512,
		"LimCap": 3,
		"BigCap": 1000000,
	}
	for testname, size := range benchTest {
		b.Run(testname, func(b *testing.B) {
			cache := NewMemStore(
				store.WithContext(
					NewContext(
						context.Background(),
						map[string]interface{}{
							"maxCap": size,
						},
					)),
			)

			record := &store.Record{}
			for i := 0; i < size; i++ {
				v := strconv.Itoa(i)
				record.Key = v
				record.Value = []byte(v)
				_ = cache.Write(record)
			}

			b.SetParallelism(threads)
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				record := &store.Record{}
				for pb.Next() {
					v := strconv.Itoa(rand.Intn(size * 2)) //nolint:gosec
					_ = cache.Delete(v)
					record.Key = v
					record.Value = []byte(v)
					_ = cache.Write(record)
				}
			})
		})
	}
}

func BenchmarkConcurrent(b *testing.B) {
	threads := []int{3, 10, 50}
	for _, nThreads := range threads {
		nt := strconv.Itoa(nThreads)
		b.Run("StoreT"+nt, func(b *testing.B) {
			concurrentStoreBench(b, nThreads)
		})
		b.Run("RetrieveT"+nt, func(b *testing.B) {
			concurrentRetrieveBench(b, nThreads)
		})
		b.Run("RemoveT"+nt, func(b *testing.B) {
			concurrentRemoveBench(b, nThreads)
		})
	}
}

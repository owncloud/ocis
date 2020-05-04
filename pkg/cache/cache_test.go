package cache

import (
	"testing"
)

// Prevents from invalid import cycle.
type AccountsCacheEntry struct {
	Email string
	UUID  string
}

func TestSet(t *testing.T) {
	c := NewCache(
		Size(256),
	)

	err := c.Set("accounts", "hello@foo.bar", AccountsCacheEntry{
		Email: "hello@foo.bar",
		UUID:  "9c31b040-59e2-4a2b-926b-334d9e3fbd05",
	})
	if err != nil {
		t.Error(err)
	}

	if c.Length("accounts") != 1 {
		t.Errorf("expected length 1 got `%v`", len(c.entries))
	}

	item, err := c.Get("accounts", "hello@foo.bar")
	if err != nil {
		t.Error(err)
	}

	if cachedEntry, ok := item.V.(AccountsCacheEntry); !ok {
		t.Errorf("invalid cached value type")
	} else {
		if cachedEntry.Email != "hello@foo.bar" {
			t.Errorf("invalid value. Expected `hello@foo.bar` got: `%v`", cachedEntry.Email)
		}
	}
}

func TestGet(t *testing.T) {
	svcCache := NewCache(
		Size(256),
	)

	err := svcCache.Set("accounts", "node", "0.0.0.0:1234")
	if err != nil {
		t.Error(err)
	}

	raw, err := svcCache.Get("accounts", "node")
	if err != nil {
		t.Error(err)
	}

	v, ok := raw.V.(string)
	if !ok {
		t.Errorf("invalid type on service node key")
	}

	if v != "0.0.0.0:1234" {
		t.Errorf("expected `0.0.0.0:1234` got `%v`", v)
	}
}

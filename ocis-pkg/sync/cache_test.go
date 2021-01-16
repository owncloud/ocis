package sync

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestCache_Get(t *testing.T) {
	size := 1024
	c := NewCache(size)

	for i := 0; i < size; i++ {
		c.Set(strconv.Itoa(i), i, time.Now().Add(10*time.Second))
	}

	for i := 0; i < size; i++ {
		assert.Equal(t, i, c.Get(strconv.Itoa(i)).V, "entry value is the same")
	}

	assert.Nil(t, c.Get("unknown"), "entry is nil if unknown")

	wait := 10 * time.Millisecond
	c.Set("expired", size, time.Now().Add(wait))
	time.Sleep(wait + 1)
	assert.Nil(t, c.Get(strconv.Itoa(size)), "entry is nil if it's expired")
}

func TestCache_Set(t *testing.T) {
	c := NewCache(1)

	c.Set("new", "new", time.Now().Add(10*time.Millisecond))
	assert.Equal(t, "new", c.Get("new").V, "new entries can be added")
	assert.Equal(t, 1, c.length, "adding new entries will increase the cache size")

	replacedExpiration := time.Now().Add(10 * time.Millisecond)
	c.Set("new", "updated", replacedExpiration)
	assert.Equal(t, "updated", c.Get("new").V, "entry values can be updated")
	assert.Equal(t, replacedExpiration, c.Get("new").expiration, "entry expiration can be updated")

	time.Sleep(11 * time.Millisecond)
	c.Set("eviction", "eviction", time.Now())
	assert.Equal(t, 1, c.length, "expired entries get removed")
}

func TestCache_Unset(t *testing.T) {
	c := NewCache(1)

	c.Set("new", "new", time.Now().Add(10*time.Millisecond))
	c.Unset("new")
	assert.Nil(t, c.Get("new"), "entries can be removed")
	assert.Equal(t, 0, c.length, "removing a entry decreases the cache size")
}

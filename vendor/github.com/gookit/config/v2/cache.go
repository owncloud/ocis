package config

// The read caches(strCache, sArrCache, sMapCache) are shared by all the getter methods
// and are filled on the read path. They are guarded by cacheMu, and are invalidated by
// ClearCaches(), which also bumps cacheGen.
//
// Rules:
//   - the lock order is hookMu -> lock -> cacheMu. the data write path holds the data lock
//     and then calls ClearCaches(), so a cache helper must NEVER take the data lock while
//     cacheMu is held, or it will deadlock. that's why they are small leaf methods.
//   - a value is only stored when cacheGen is not changed by a data write in the meantime,
//     so a value read before a data change can not be stored into the cache after it.

// cacheGeneration return the current cache generation.
func (c *Config) cacheGeneration() uint64 {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	return c.cacheGen
}

// cachedString get a string value from the cache.
func (c *Config) cachedString(key string) (value string, ok bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	value, ok = c.strCache[key]
	return
}

// cacheString store a string value to the cache by generation.
func (c *Config) cacheString(gen uint64, key, value string) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.cacheGen != gen {
		return
	}
	if c.strCache == nil {
		c.strCache = make(map[string]string)
	}
	c.strCache[key] = value
}

// cachedStringArr get a string array value from the cache.
func (c *Config) cachedStringArr(key string) (arr strArr, ok bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	arr, ok = c.sArrCache[key]
	return
}

// cacheStringArr store a string array value to the cache by generation.
func (c *Config) cacheStringArr(gen uint64, key string, arr strArr) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.cacheGen != gen {
		return
	}
	if c.sArrCache == nil {
		c.sArrCache = make(map[string]strArr)
	}
	c.sArrCache[key] = arr
}

// cachedStringMap get a string map value from the cache.
func (c *Config) cachedStringMap(key string) (mp strMap, ok bool) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	mp, ok = c.sMapCache[key]
	return
}

// cacheStringMap store a string map value to the cache by generation.
func (c *Config) cacheStringMap(gen uint64, key string, mp strMap) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	if c.cacheGen != gen {
		return
	}
	if c.sMapCache == nil {
		c.sMapCache = make(map[string]strMap)
	}
	c.sMapCache[key] = mp
}

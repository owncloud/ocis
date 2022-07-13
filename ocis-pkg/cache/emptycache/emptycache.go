package emptycache

// EmptyCache implements the Cache interface
type EmptyCache struct {
}

// This method won't do anything. There is no need to initialize anything.
// An empty map can be used to fill the required parameter.
// No error will be returned
func (c *EmptyCache) Initialize(params map[string]interface{}) error {
	return nil
}

// No value will be stored. No error will be returned
func (c *EmptyCache) Store(key string, value string, ttl int64) error {
	return nil
}

// Since no value will be stored, no value will be retrieved.
// The returned values will always be ("", false, nil)
func (c *EmptyCache) Retrieve(key string) (string, bool, error) {
	return "", false, nil
}

// Nothing stored, so nothing to remove. There won't be any error
func (c *EmptyCache) Remove(key string) error {
	return nil
}

func NewEmptyCache() *EmptyCache {
	return &EmptyCache{}
}

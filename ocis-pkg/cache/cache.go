package cache

// Define the methods required to implement a cache.
//
// This interface just focus on access methods, with the additional `Initialize`
// one. The only expectation here is that the `Initialize` method is called
// before any other one, just after the cache instance is create
//
// ```
// func NewCustomInstance() *Cache {
//   .....
// }
//
// myCache := NewCustomInstance()
// err1 := myCache.Initialize(iniParams)
// ......
// err2 := myCache.Store("myKey", "myValue")
// ```
// There are no big requirements for the implementations. A lot of things are
// expected to be handled internally by the implementation, and maybe
// allowing some configuration via the `Initialize` method.
// * Capacity of the cache
// * Default ttl
// * Eviction policies
// * Data management
// * Connectivity with external services
type Cache interface {
	// Initialize the cache. The parameters will be used to setup the cache
	// with the corresponding values. The setup can include capacity, default ttl,
	// eviction policy, connectivity parameters with external services, etc.
	// If the cache requires a prefix because it could be shared among multiple
	// services, it should be configured here.
	//
	// The specific parameters will depend in the implementation. Additional
	// parameters might be sent, but the implementation should ignore them and
	// not throw an error.
	//
	// Return an error if something went wrong
	Initialize(params map[string]interface{}) error

	// Store the value in the key. A ttl of 0 should be used to indicate
	// that the key shouldn't expire.
	// Unless the cache uses a different policy (it should be explicitly
	// documented), the new value will overwrite the previous one
	//
	// The cache shouldn't need to guarantee that the key is removed when
	// the ttl is reached, but it MUST guarantee that the value won't be
	// retrieved.
	// There are 2 main approaches to take if the key has reached its ttl:
	// * The key will stay in the cache and it will be removed when it's
	// retrieved. So the `Retrieve` method will return ("", false, nil) on
	// a expired key, and the key will be removed at that point.
	// * A service will be started to monitor the ttl of the keys. This
	// service must be handled completely internally to the cache implementation
	//
	// Return an error if the value couldn't be saved.
	Store(key string, value string, ttl int64) error

	// Get the value stored in the key. The method will also return whether
	// the key exists in the cache.
	//
	// For the return values, the possible combinations are:
	// (<any string>, true, nil) -> the stored value, including the empty string
	// ("", false, nil) -> the key doesn't exists
	// ("", false, error) -> if an error happened, for example, connectivity lost.
	//
	// Return the value stored or an empty string as the first return value,
	// whether the value exists or not in the cache as the second return value,
	// and any possible error as third return value.
	Retrieve(key string) (string, bool, error)

	// Explicitly remove the key from the cache. If the key doesn't exists, the
	// method should ignore it and not throw an error.
	//
	// An error will be returned if the key isn't removed
	Remove(key string) error
}

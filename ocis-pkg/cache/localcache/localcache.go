package localcache

import (
	"container/list"
	"sync"
	"time"
)

type cacheinfo struct {
	keyRef     string
	value      string
	validUntil time.Time
}

// Implements a thread-safe local cache with a LRU replacement mechanism.
// The local cache has a maximum capacity (configurable) which will be guaranteed.
// The implementation is focused on performance. In order to do so, it's important
// to know that the replacement mechanism **WON'T** take into account expired items.
// This means that a valid element might be removed from the cache even having
// expired elements in the cache.
//
// To clarify, this cache guarantees that **AT MOST** an element will be accessible
// during its ttl, but not later. However, it doesn't guarantee that such element
// will be present and accessible during its whole ttl.
//
// A linked list is used to keep track of the order of removal. New elements are
// added at the end of the list, so they'll be the last ones to be removed. All
// operations affect this linked list one way or another, usually moving the target
// element to the last position.
// Note that there is also a map which is used to guarantee fast access to the items.
// All of this is handled transparently.
//
// Note that there are additional operations not covered by the interface. Those
// operations are NOT intended to be used. Those additional operations don't
// guarantee thread-safety and they should be used only for testing.
//
// * Ensures the max capacity is respected, preventing using too much memory
// * Ensures constant time in all operations of the interface regardless of
// the number of elements.
// * Expired elements might be present, but they'll be removed before retrieving
// * A least-recently-used policy is used when the cache needs to free space. As said
// expired items won't be removed preferently unless explicitly accessed.
type LocalCache struct {
	data    map[string]*list.Element
	keyList *list.List
	mutex   sync.Mutex
	maxCap  int
}

// Creates a LocalCache instance. You must call the `Initialize` afterwards
func NewLocalCache() *LocalCache {
	return &LocalCache{}
}

// Initialize the LocalCache instance. Internal data structures will be created.
// You can provide optional parameters to configure some parts of the instance,
// in particular, the following options are supported for the LocalCache:
// * "capacity": (int) The maximum number of elements that this cache will hold.
// The maximum capacity will be enforced regardless of the ttls of the elements,
// which means that "valid" elements might be removed from the cache to hold new
// elements.
//
// These are the defaults:
// * "capacity": 512
//
// Note that this method is intended to be call just after the creation of
// the cache. It should be called only once, but multiple calls are allowed.
// Each additional call will reinitialize the cache with the new parameters
// and the data contained will be lost.
// This method isn't thread-safe and it should be call only once by the thread
// creating the cache.
func (c *LocalCache) Initialize(params map[string]interface{}) error {
	dataCapacity := 512

	capacity, capacityExists := params["capacity"]
	if capacityExists {
		if capacityInt, typeOk := capacity.(int); typeOk {
			dataCapacity = capacityInt
		}
	}

	c.maxCap = dataCapacity
	c.data = make(map[string]*list.Element, dataCapacity)
	c.keyList = list.New()
	return nil
}

// Stores the value under the key for a ttl duration.
// Updating an element will move the element to the last position in the linked
// list, causing the element to be removed the last (if needed).
// A ttl of 0 will cause the value to be stored until "1<<62" Unix time,
// this mean it will expire on 146138514283-06-19 07:45:04 +0000 UTC.
// In practice, the value won't expire.
func (c *LocalCache) Store(key string, value string, ttl int64) error {
	var validUntil time.Time

	if ttl == 0 {
		validUntil = time.Unix(1<<62, 0)
	} else {
		duration := time.Duration(ttl) * time.Second
		validUntil = time.Now().Add(duration)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.data[key]
	if !ok {
		// if key isn't present
		nKeyItems := c.keyList.Len()
		if nKeyItems >= c.maxCap {
			// free some space
			toBeRemoved := c.keyList.Front()
			if toBeRemoved != nil {
				toBeRemovedInfo := toBeRemoved.Value.(*cacheinfo)
				delete(c.data, toBeRemovedInfo.keyRef)
				c.keyList.Remove(toBeRemoved)
			}
		}
		info := &cacheinfo{
			keyRef:     key,
			value:      value,
			validUntil: validUntil,
		}
		newElement := c.keyList.PushBack(info)
		c.data[key] = newElement
	} else {
		// if key is present, just update the values
		info := elem.Value.(*cacheinfo)
		info.value = value
		info.validUntil = validUntil
		c.keyList.MoveToBack(elem)
	}
	return nil
}

// Retrieve the stored element from the cache.
// The following return values are expected:
// * <stringvalue>, true, nil -> the value stored under the requested key. Note
// that the value can still be an empty string if such value was stored. Use the
// true/false returned value (second one) to know if such value exists.
// * "", false, nil -> the key didn't exist.
//
// Errors aren't expected in this method.
func (c *LocalCache) Retrieve(key string) (string, bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.data[key]
	if !ok {
		return "", false, nil
	}
	info := elem.Value.(*cacheinfo)

	value := info.value
	exists := true
	// check ttl
	if info.validUntil.After(time.Now()) {
		c.keyList.MoveToBack(elem)
	} else {
		// expired item
		delete(c.data, key)
		c.keyList.Remove(elem)
		value = ""
		exists = false
	}
	return value, exists, nil
}

// Remove the target key. If the key doesn't exist, this method won't
// do anything.
// No error is expected in this method.
func (c *LocalCache) Remove(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.data[key]
	if ok {
		delete(c.data, key)
		c.keyList.Remove(elem)
	}
	return nil
}

//
// Additional methods not covered by the interface. These methods are intended
// to be used to debug or test
//

// Get the current maximum capacity configured in this cache
func (c *LocalCache) MaxCap() int {
	return c.maxCap
}

// Get the length of the underlying map
func (c *LocalCache) MapLen() int {
	return len(c.data)
}

// Get the length of the underlying list
func (c *LocalCache) ListLen() int {
	return c.keyList.Len()
}

// A callback to be used to traverse the stored elements. The elements
// will be traversed using the expected deletion order at a given time. This
// should guarantee a more or less predictable order. Note that this order
// will change based on the access and new addition of elements in the cache.
type LCCallback func(key string, value string)

// This function isn't part if the interface and it's intended to be used
// only for debugging.
// This will traverse the whole list of element and the cache will be locked
// while traversing. Operations over this cache can't be performed while the
// traversal is ongoing. Don't use this function in production code.
// The order will be the expected deletion order of items at a given time. This
// means that it will change based on the cache usage.
func (c *LocalCache) TraverseList(callback LCCallback) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	currentElement := c.keyList.Front()
	for currentElement != nil {
		info := currentElement.Value.(*cacheinfo)
		callback(info.keyRef, info.value)
		currentElement = currentElement.Next()
	}
}

package ordered_map

import (
	"fmt"
)

type OrderedMap struct {
	store  map[interface{}]interface{}
	mapper map[interface{}]*node
	root   *node
}

func NewOrderedMap() *OrderedMap {
	om := &OrderedMap{
		store:  make(map[interface{}]interface{}),
		mapper: make(map[interface{}]*node),
		root:   newRootNode(),
	}
	return om
}

func NewOrderedMapWithArgs(args []*KVPair) *OrderedMap {
	om := NewOrderedMap()
	om.update(args)
	return om
}

func (om *OrderedMap) update(args []*KVPair) {
	for _, pair := range args {
		om.Set(pair.Key, pair.Value)
	}
}

func (om *OrderedMap) Set(key interface{}, value interface{}) {
	if _, ok := om.store[key]; ok == false {
		root := om.root
		last := root.Prev
		last.Next = newNode(last, root, key)
		root.Prev = last.Next
		om.mapper[key] = last.Next
	}
	om.store[key] = value
}

func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	val, ok := om.store[key]
	return val, ok
}

func (om *OrderedMap) Delete(key interface{}) {
	_, ok := om.store[key]
	if ok {
		delete(om.store, key)
	}
	root, rootFound := om.mapper[key]
	if rootFound {
		prev := root.Prev
		next := root.Next
		prev.Next = next
		next.Prev = prev
		delete(om.mapper, key)
	}
}

func (om *OrderedMap) String() string {
	builder := make([]string, len(om.store))

	var index int = 0
	iter := om.IterFunc()
	for kv, ok := iter(); ok; kv, ok = iter() {
		val, _ := om.Get(kv.Key)
		builder[index] = fmt.Sprintf("%v:%v, ", kv.Key, val)
		index++
	}
	return fmt.Sprintf("OrderedMap%v", builder)
}

func (om *OrderedMap) Iter() <-chan *KVPair {
	println("Iter() method is deprecated!. Use IterFunc() instead.")
	return om.UnsafeIter()
}

/*
Beware, Iterator leaks goroutines if we do not fully traverse the map.
For most cases, `IterFunc()` should work as an iterator.
 */
func (om *OrderedMap) UnsafeIter() <-chan *KVPair {
	keys := make(chan *KVPair)
	go func() {
		defer close(keys)
		var curr *node
		root := om.root
		curr = root.Next
		for curr != root {
			v, _ := om.store[curr.Value]
			keys <- &KVPair{curr.Value, v}
			curr = curr.Next
		}
	}()
	return keys
}

func (om *OrderedMap) IterFunc() func() (*KVPair, bool) {
	var curr *node
	root := om.root
	curr = root.Next
	return func() (*KVPair, bool) {
		for curr != root {
			tmp := curr
			curr = curr.Next
			v, _ := om.store[tmp.Value]
			return &KVPair{tmp.Value, v}, true
		}
		return nil, false
	}
}

func (om *OrderedMap) Len() int {
	return len(om.store)
}


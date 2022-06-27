package ordered_map

import "fmt"

type KVPair struct {
	Key   interface{}
	Value interface{}
}

func (k *KVPair) String() string {
	return fmt.Sprintf("%v:%v", k.Key, k.Value)
}

func (kv1 *KVPair) Compare(kv2 *KVPair) bool {
	return kv1.Key == kv2.Key && kv1.Value == kv2.Value
}
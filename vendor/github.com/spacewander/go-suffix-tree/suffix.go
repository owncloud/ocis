package suffix

import (
	"bytes"
	"sort"
)

// Return
// the first index of the mismatch byte (from right to left, starts from 1)
// len(left)+1 if left byte sequence is shorter than right one
// 0 if two byte sequences are equal
// -len(right)-1 if left byte sequence is longer than right one
func suffixDiff(left, right []byte) int {
	leftLen := len(left)
	rightLen := len(right)
	minLen := leftLen
	if minLen > rightLen {
		minLen = rightLen
	}
	for i := 1; i <= minLen; i++ {
		if left[leftLen-i] != right[rightLen-i] {
			return i
		}
	}
	if leftLen < rightLen {
		return leftLen + 1
	} else if leftLen == rightLen {
		return 0
	}
	return -rightLen - 1
}

type _Edge struct {
	label []byte
	// Could be either Node or Leaf
	point interface{}
}

type _Leaf struct {
	// For LongestSuffix and so on. We choice to use more memory(24 bytes per node)
	// over appending keys each time.
	originKey []byte
	value     interface{}
}

type _Node struct {
	edges []*_Edge
}

func (node *_Node) insertEdge(edge *_Edge) {
	newEdgeLabelLen := len(edge.label)
	idx := sort.Search(len(node.edges), func(i int) bool {
		return newEdgeLabelLen < len(node.edges[i].label)
	})
	node.edges = append(node.edges, nil)
	copy(node.edges[idx+1:], node.edges[idx:])
	node.edges[idx] = edge
}

func (node *_Node) removeEdge(idx int) {
	copy(node.edges[idx:], node.edges[idx+1:])
	node.edges[len(node.edges)-1] = nil
	node.edges = node.edges[:len(node.edges)-1]
}

// Reorder edge which is not shorter than before
func (node *_Node) backwardEdge(idx int) {
	edge := node.edges[idx]
	edgeLabelLen := len(edge.label)
	edgesLen := len(node.edges)
	if idx == edgesLen-1 {
		// Still longest, no need to change
		return
	}
	// Get the first edge which's label is longer than this edge...
	i := sort.Search(edgesLen-idx-1, func(j int) bool {
		return edgeLabelLen < len(node.edges[j+idx+1].label)
	})
	// ... and insert before it. (Note that we just add `idx` instead of `idx+1`)
	i += idx
	copy(node.edges[idx:i], node.edges[idx+1:i+1])
	node.edges[i] = edge
}

// Reorder edge which is shorter than before
func (node *_Node) forwardEdge(idx int) {
	edge := node.edges[idx]
	edgeLabelLen := len(edge.label)
	i := sort.Search(idx, func(j int) bool {
		return edgeLabelLen < len(node.edges[j].label)
	})
	copy(node.edges[i+1:idx+1], node.edges[i:idx])
	node.edges[i] = edge
}

func (node *_Node) insert(originKey []byte, key []byte, value interface{}) (
	oldValue interface{}, ok bool) {

	start := 0
	if len(node.edges) > 0 && len(node.edges[0].label) == 0 {
		// handle empty label as a special case, so the rest of labels don't share
		// common suffix
		if len(key) == 0 {
			leaf, _ := node.edges[0].point.(*_Leaf)
			oldValue = leaf.value
			leaf.value = value
			return oldValue, true
		}
		start++
	}
	for i := start; i < len(node.edges); i++ {
		edge := node.edges[i]
		gap := suffixDiff(key, edge.label)
		if gap == 0 {
			// CASE 1: key == label
			switch point := edge.point.(type) {
			case *_Leaf:
				// Leaf hitted, replace old value
				oldValue = point.value
				point.value = value
				return oldValue, true
			case *_Node:
				// Node hitted, insert a leaf under this Node
				return point.insert(originKey, []byte{}, value)
			}
		} else if gap < 0 {
			// CASE 2: key > label
			gap = -gap
			label := key[:len(key)-gap+1]
			switch point := edge.point.(type) {
			case *_Leaf:
				// Before: Node - "label" -> Leaf(Value1)
				// After: Node - "label" - Node - "" -> Leaf(Value1)
				//							|- "s" -> Leaf(Value2)
				// Create new Node, move old Leaf under new Node, and then
				//	insert a new Leaf
				newNode := &_Node{
					edges: []*_Edge{
						{
							label: []byte{},
							point: point,
						},
						{
							label: label,
							point: &_Leaf{
								originKey: originKey,
								value:     value,
							},
						},
					},
				}
				edge.point = newNode
				return nil, true
			case *_Node:
				// Before: Node - "label" -> Node - "" -> Leaf(Value1)
				// After: Node - "label" - Node - "" -> Leaf(Value1)
				//							|- "s" -> Leaf(Value2)
				// Insert a new Leaf with extra data as label
				return point.insert(originKey, label, value)
			}
		} else if gap > 1 {
			// CASE 3: mismatch(key, label) after first letter or key < label
			// Before: Node - "labels" -> Node/Leaf(Value1)
			// After: Node - "label" - Node - "s" -> Node/Leaf(Value1)
			//						    |- "" -> Leaf(Value2)
			// Before: Node - "label" -> Node/Leaf(Value1)
			// After: Node - "lab" - Node - "el" -> Node/Leaf(Value1)
			//							|- "or" -> Leaf(Value2)
			newEdge := &_Edge{
				label: edge.label[:len(edge.label)-gap+1],
				point: edge.point,
			}
			keyEdge := &_Edge{
				label: key[:len(key)-gap+1],
				point: &_Leaf{
					originKey: originKey,
					value:     value,
				},
			}
			newNode := &_Node{
				edges: make([]*_Edge, 2),
			}
			if len(newEdge.label) < len(keyEdge.label) {
				newNode.edges[0], newNode.edges[1] = newEdge, keyEdge
			} else {
				newNode.edges[0], newNode.edges[1] = keyEdge, newEdge
			}
			edge.point = newNode
			edge.label = edge.label[len(edge.label)-gap+1:]
			node.forwardEdge(i)
			return nil, true
		}
		// CASE 4: totally mismatch
	}

	leaf := &_Leaf{
		originKey: originKey,
		value:     value,
	}
	edge := &_Edge{
		label: key,
		point: leaf,
	}
	node.insertEdge(edge)
	return nil, true
}

func (node *_Node) get(key []byte) (value interface{}, found bool) {
	edges := node.edges
	start := 0
	if len(edges[0].label) == 0 {
		// handle empty label as a special case, so the rest of labels don't share
		// common suffix
		if len(key) == 0 {
			leaf, _ := edges[0].point.(*_Leaf)
			return leaf.value, true
		}
		start++
	}

	keyLen := len(key)
	for i := start; i < len(edges); i++ {
		edge := edges[i]
		edgeLabelLen := len(edge.label)
		if keyLen > edgeLabelLen {
			if bytes.Equal(key[len(key)-len(edge.label):], edge.label) {
				subKey := key[:len(key)-len(edge.label)]
				switch point := edge.point.(type) {
				case *_Leaf:
					return nil, false
				case *_Node:
					return point.get(subKey)
				}
			}
		} else if keyLen == edgeLabelLen {
			if bytes.Equal(key, edge.label) {
				switch point := edge.point.(type) {
				case *_Leaf:
					return point.value, true
				case *_Node:
					return point.get([]byte{})
				}
			}
		} else {
			break
		}
	}

	return nil, false
}

func (node *_Node) longestSuffix(key []byte) (matchedKey []byte, value interface{}, found bool) {
	edges := node.edges
	start := 0
	if len(edges[0].label) == 0 {
		// handle empty label as a special case, so the rest of labels don't share
		// common suffix
		if len(key) == 0 {
			leaf, _ := edges[0].point.(*_Leaf)
			return leaf.originKey, leaf.value, true
		}
		start++
	}

	keyLen := len(key)
	for i := start; i < len(edges); i++ {
		edge := edges[i]
		edgeLabelLen := len(edge.label)
		if keyLen > edgeLabelLen {
			if bytes.Equal(key[len(key)-len(edge.label):], edge.label) {
				subKey := key[:len(key)-len(edge.label)]
				switch point := edge.point.(type) {
				case *_Leaf:
					return point.originKey, point.value, true
				case *_Node:
					matchedKey, value, found := point.longestSuffix(subKey)
					if found {
						return matchedKey, value, found
					}
				}
			}
		} else if keyLen == edgeLabelLen {
			if bytes.Equal(key, edge.label) {
				switch point := edge.point.(type) {
				case *_Leaf:
					return point.originKey, point.value, true
				case *_Node:
					matchedKey, value, found := point.longestSuffix([]byte{})
					if found {
						return matchedKey, value, found
					}
				}
			}
		} else {
			break
		}
	}

	if start == 1 {
		leaf, _ := edges[0].point.(*_Leaf)
		return leaf.originKey, leaf.value, true
	}

	return nil, nil, false
}

func (node *_Node) mergeChildNode(idx int, child *_Node) {
	if len(child.edges) == 1 {
		edge := node.edges[idx]
		edge.point = child.edges[0].point
		edge.label = append(child.edges[0].label, edge.label...)
		node.backwardEdge(idx)
	}
	// When child has only one edge, we will remove the child and merge its label,
	// So there is no case that child has no edge.
}

func (node *_Node) remove(key []byte) (value interface{}, found bool, childRemoved bool) {
	edges := node.edges
	start := 0
	if len(edges[0].label) == 0 {
		// handle empty label as a special case, so the rest of labels don't share
		// common suffix
		if len(key) == 0 {
			leaf, _ := edges[0].point.(*_Leaf)
			value = leaf.value
			node.removeEdge(0)
			return value, true, true
		}
		start++
	}

	keyLen := len(key)
	for i := start; i < len(edges); i++ {
		edge := edges[i]
		edgeLabelLen := len(edge.label)
		if keyLen > edgeLabelLen {
			if bytes.Equal(key[len(key)-len(edge.label):], edge.label) {
				key := key[:len(key)-len(edge.label)]
				switch point := edge.point.(type) {
				case *_Node:
					value, found, childRemoved = point.remove(key)
					if childRemoved {
						node.mergeChildNode(i, point)
					}
					return value, found, false
				}
			}
		} else if keyLen == edgeLabelLen {
			if bytes.Equal(key, edge.label) {
				switch point := edge.point.(type) {
				case *_Leaf:
					value = point.value
					node.removeEdge(i)
					return value, true, true
				case *_Node:
					value, found, childRemoved = point.remove([]byte{})
					if childRemoved {
						node.mergeChildNode(i, point)
					}
					return value, found, false
				}
			}
		} else {
			break
		}
	}

	return nil, false, false
}

// return either _Leaf or _Node as interface{}
func (node *_Node) getPointHasSuffix(key []byte) (interface{}, []byte, bool) {
	edges := node.edges
	keyLen := len(key)
	for i := len(edges) - 1; i >= 0; i-- {
		edge := edges[i]
		edgeLabelLen := len(edge.label)
		if keyLen > edgeLabelLen {
			if bytes.Equal(key[len(key)-len(edge.label):], edge.label) {
				subKey := key[:len(key)-len(edge.label)]
				switch point := edge.point.(type) {
				case *_Leaf:
					return nil, nil, false
				case *_Node:
					return point.getPointHasSuffix(subKey)
				}
			}
		} else {
			if bytes.HasSuffix(edge.label, key) {
				return edge.point, edge.label[:len(edge.label)-len(key)], true
			}
		}
	}
	return nil, nil, false
}

func (node *_Node) walk(suffix []byte, f func(key []byte, value interface{}) bool, stop *bool) {
	for _, edge := range node.edges {
		if *stop {
			return
		}
		switch point := edge.point.(type) {
		case *_Leaf:
			*stop = f(append(edge.label, suffix...), point.value)
		case *_Node:
			point.walk(append(edge.label, suffix...), f, stop)
		}
	}
}

func (node *_Node) walkNode(suffix [][]byte, f func(labels [][]byte, value interface{})) {
	f(append([][]byte{nil}, suffix...), nil)
	nodes := []*_Edge{}
	leaves := []*_Edge{}
	for _, edge := range node.edges {
		switch edge.point.(type) {
		case *_Leaf:
			leaves = append(leaves, edge)
		case *_Node:
			nodes = append(nodes, edge)
		}
	}
	for _, edge := range leaves {
		leaf, _ := edge.point.(*_Leaf)
		f(append([][]byte{edge.label}, suffix...), leaf.value)
	}
	for _, edge := range nodes {
		node, _ := edge.point.(*_Node)
		node.walkNode(append([][]byte{edge.label}, suffix...), f)
	}
}

// Tree represents a suffix tree.
type Tree struct {
	root      *_Node
	leavesNum int
}

// NewTree create a suffix tree for future usage.
func NewTree() *Tree {
	return &Tree{
		root: &_Node{
			edges: []*_Edge{},
		},
		leavesNum: 0,
	}
}

// Insert suffix tree with given key and value. Return the previous value and a boolean to
// indicate whether the insertion is successful.
func (tree *Tree) Insert(key []byte, value interface{}) (oldValue interface{}, ok bool) {
	if key == nil {
		return nil, false
	}
	oldValue, ok = tree.root.insert(key, key, value)
	if ok && oldValue == nil {
		tree.leavesNum++
	}
	return oldValue, ok
}

// Get returns the value of given key and a boolean to indicate
// whether the value is found.
func (tree *Tree) Get(key []byte) (value interface{}, found bool) {
	if key == nil || len(tree.root.edges) == 0 {
		return nil, false
	}
	return tree.root.get(key)
}

// LongestSuffix is mostly like Get.
// It returns the key which is the longest suffix of the given key,
// and the value referred by this key.
// Plus a boolean to indicate whether the key/value, is found.
func (tree *Tree) LongestSuffix(key []byte) (matchedKey []byte, value interface{}, found bool) {
	if key == nil || len(tree.root.edges) == 0 {
		return nil, nil, false
	}
	return tree.root.longestSuffix(key)
}

// Remove returns the value of given key and a boolean to indicate
// whethe the value is found. Then the value will be removed.
func (tree *Tree) Remove(key []byte) (oldValue interface{}, found bool) {
	if key == nil || len(tree.root.edges) == 0 {
		return nil, false
	}
	oldValue, found, _ = tree.root.remove(key)
	if found {
		tree.leavesNum--
	}
	return oldValue, found
}

// Len returns the number of keys.
func (tree *Tree) Len() int {
	return tree.leavesNum
}

// Walk through the tree, call function with key and value.
// Once the function returns true, it will stop walking.
// The travelling order is DFS, in the same suffix level the shortest key comes first.
func (tree *Tree) Walk(f func(key []byte, value interface{}) bool) {
	stop := false
	tree.root.walk([]byte{}, f, &stop)
}

// WalkSuffix travels through nodes which have given suffix, calls function with key and value.
// Once the function returns true, it will stop walking.
// The travelling order is DFS, in the same suffix level the shortest key comes first.
func (tree *Tree) WalkSuffix(suffix []byte, f func(key []byte, value interface{}) bool) {
	if len(tree.root.edges) != 0 {
		stop := false
		if len(suffix) == 0 {
			tree.root.walk([]byte{}, f, &stop)
		} else {
			startingPoint, extraLabel, found := tree.root.getPointHasSuffix(suffix)
			if found {
				switch point := startingPoint.(type) {
				case *_Leaf:
					f(point.originKey, point.value)
				case *_Node:
					point.walk(append(extraLabel, suffix...), f, &stop)
				}
			}
		}
	}
}

// This API is for testing/debug
func (tree *Tree) walkNode(f func(labels [][]byte, value interface{})) {
	tree.root.walkNode([][]byte{}, f)
}

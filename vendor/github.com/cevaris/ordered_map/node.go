package ordered_map

import (
	"fmt"
	"bytes"
)

type node struct {
	Prev  *node
	Next  *node
	Value interface{}
}

func newRootNode() *node {
	root := &node{}
	root.Prev = root
	root.Next = root
	return root
}

func newNode(prev *node, next *node, key interface{}) *node {
	return &node{Prev: prev, Next: next, Value: key}
}

func (n *node) Add(value string) {
	root := n
	last := root.Prev
	last.Next = newNode(last, n, value)
	root.Prev = last.Next
}

func (n *node) String() string {
	var buffer bytes.Buffer
	if n.Value == "" {
		// Need to sentinel
		var curr *node
		root := n
		curr = root.Next
		for curr != root {
			buffer.WriteString(fmt.Sprintf("%s, ", curr.Value))
			curr = curr.Next
		}
	} else {
		// Else, print pointer value
		buffer.WriteString(fmt.Sprintf("%p, ", &n))
	}
	return fmt.Sprintf("LinkList[%v]", buffer.String())
}

func (n *node) IterFunc() func() (string, bool) {
	var curr *node
	root := n
	curr = root.Next
	return func() (string, bool) {
		for curr != root {
			tmp := curr.Value.(string)
			curr = curr.Next
			return tmp, true
		}
		return "", false
	}
}

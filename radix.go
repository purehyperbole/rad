package rad

import (
	"fmt"
	"os"
)

// Radix tree
type Radix struct {
	root *Node
}

// New Radix Tree
func New() *Radix {
	return &Radix{
		root: &Node{
			edges: &[256]*Node{},
		},
	}
}

// Insert or update a keypair
func (r *Radix) Insert(key []byte, value interface{}) bool {
	parent, node, pos, dv := r.findInsertionPoint(key)

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Recovered in f", x)

			parent.print()
			node.print()
			fmt.Println(pos)
			fmt.Println(dv)
			os.Exit(1)
		}
	}()

	switch {
	case pos < len(key) && node == nil:
		return r.insertNode(key, value, parent, node, pos, dv)
	case pos == len(key)+dv && len(node.prefix) == 0, pos == len(key)+dv && len(node.prefix) > 1:
		return r.updateNode(key, value, parent, node, pos, dv)
	case (len(key) - (pos + dv)) > 0:
		return r.splitThreeWay(key, value, parent, node, pos, dv)
	case (len(key) - (pos + dv)) == 0:
		return r.splitTwoWay(key, value, parent, node, pos, dv)
	}

	return false
}

// MustInsert attempts to insert a value until it is successful
func (r *Radix) MustInsert(key []byte, value interface{}) {
	for !r.Insert(key, value) {
		fmt.Println("retrying: ", string(key))
	}
}

// Lookup a value by key
func (r *Radix) Lookup(key []byte) interface{} {
	var i, dv int

	n := r.root

	for n.next(key[i]) != nil {
		n = n.next(key[i])
		i++

		if len(n.prefix) > 0 {
			dv = divergence(n.prefix, key[i:])

			if len(n.prefix) > dv {
				break
			}

			i = i + dv
		}

		// if we've found the key, break the loop
		if i == len(key) {
			break
		}
	}

	if n == nil || len(key) > i {
		return nil
	}

	return n.value
}

func (r *Radix) findInsertionPoint(key []byte) (*Node, *Node, int, int) {
	var pos, dv int
	var node, parent *Node

	node = r.root

	for node.next(key[pos]) != nil {
		parent = node
		node = node.next(key[pos])
		pos++

		if len(node.prefix) > 0 {
			dv = divergence(node.prefix, key[pos:])

			if len(node.prefix) > dv {
				return parent, node, pos, dv
			}

			pos = pos + dv
		}

		// if we've found the key, return its parent node so it can be pointed to the new node
		if pos == len(key) {
			return parent, node, pos, 0
		}
	}

	return node, nil, pos, dv
}

func (r *Radix) insertNode(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	return parent.swapNext(key[pos], nil, &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  &[256]*Node{},
	})
}

func (r *Radix) updateNode(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	edgePos := pos - (len(node.prefix) + 1)

	fmt.Println(len(key))
	fmt.Println(edgePos)
	fmt.Println(string(node.prefix))

	return parent.swapNext(key[edgePos], node, &Node{
		prefix: node.prefix,
		value:  value,
		edges:  node.edges,
	})
}

func (r *Radix) splitTwoWay(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	var pfx []byte

	// fix issue where key is found, but is occupied by another node with prefix
	if len(key) > pos {
		pfx = key[pos : pos+dv]
	} else {
		pos--
	}

	n1 := &Node{
		prefix: pfx,
		value:  value,
		edges:  &[256]*Node{},
	}

	n2 := &Node{
		prefix: node.prefix[dv+1:],
		value:  node.value,
		edges:  node.edges,
	}

	n1.setNext(node.prefix[dv], n2)

	return parent.swapNext(key[pos], node, n1)
}

func (r *Radix) splitThreeWay(key []byte, value interface{}, parent, node *Node, pos, dv int) bool {
	n1 := &Node{
		prefix: node.prefix[:dv],
		edges:  &[256]*Node{},
	}

	n2 := &Node{
		prefix: node.prefix[dv+1:],
		value:  node.value,
		edges:  node.edges,
	}

	n3 := &Node{
		prefix: key[pos+dv+1:],
		value:  value,
		edges:  &[256]*Node{},
	}

	n1.setNext(node.prefix[dv], n2)
	n1.setNext(key[pos+dv], n3)

	return parent.swapNext(key[pos-1], node, n1)
}

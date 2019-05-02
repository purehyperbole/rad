package rad

import (
	"fmt"
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
func (r *Radix) Insert(key []byte, value interface{}) {
	parent, node, pos, dv := r.findInsertionPoint(key)

	switch {
	case pos == len(key):
		r.updateNode(key, value, parent, node, pos, dv)
	case pos < len(key) && node == nil:
		r.insertNode(key, value, parent, node, pos, dv)

	case (len(key) - (pos + dv)) > 0:
		fmt.Println("DEBUG")
		fmt.Println(len(node.prefix))
		fmt.Println(dv)
		fmt.Println(pos)

		r.splitThreeWay(key, value, parent, node, pos, dv)
	default:
		r.splitTwoWay(key, value, parent, node, pos, dv)
	}
}

// Lookup a value by key
func (r *Radix) Lookup(key []byte) interface{} {
	var i, dv int

	n := r.root

	for n.next(key[i]) != nil {
		n = n.next(key[i])
		if n == nil {
			break
		}
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
			return parent, node, pos, dv
		}
	}

	return node, nil, pos, dv
}

func (r *Radix) insertNode(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("insert", string(key))

	parent.setNext(key[pos], &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  &[256]*Node{},
	})
}

func (r *Radix) updateNode(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("update")

	parent.setNext(key[pos], &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  node.edges,
	})
}

func (r *Radix) splitTwoWay(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("split two way", string(key))

	fmt.Println("KEY:", string(key))
	fmt.Println(pos)
	fmt.Println(dv)

	parent.print()
	if node != nil {
		node.print()
	} else {
		fmt.Println("NODE IS NIL!")
	}

	fmt.Println("create node", string(key[pos]), "->", string(key[pos+1:pos+dv]))
	fmt.Println("create node", string(node.prefix[dv]), "->", string(node.prefix[dv+1:]))

	n1 := &Node{
		prefix: key[pos+1 : pos+dv],
		value:  value,
		edges:  &[256]*Node{},
	}

	n2 := &Node{
		prefix: node.prefix[dv+1:],
		value:  node.value,
		edges:  node.edges,
	}

	n1.setNext(node.prefix[dv], n2)

	parent.setNext(key[pos], n1)
}

func (r *Radix) splitThreeWay(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("split three way", string(key))

	/*
		fmt.Println("KEY:", string(key))
		fmt.Println(pos)
		fmt.Println(dv)

		fmt.Println(string(key[pos]))

		parent.print()
		if node != nil {
			node.print()
		} else {
			fmt.Println("NODE IS NIL!")
		}
	*/

	n1 := &Node{
		// prefix: key[pos:dv],
		edges: &[256]*Node{},
	}

	n2 := &Node{
		prefix: node.prefix[dv+pos:],
		value:  node.value,
		edges:  node.edges,
	}

	n3 := &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  &[256]*Node{},
	}

	n1.setNext(node.prefix[dv], n2)
	n1.setNext(key[pos], n3)

	parent.setNext(key[pos-1], n1)
}

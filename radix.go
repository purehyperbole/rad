package rad

import "fmt"

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
	case pos == len(key) && node == nil:
		r.updateNode(key, value, parent, node, pos, dv)
	case pos < len(key) && node == nil:
		r.insertNode(key, value, parent, node, pos, dv)
	case (len(key) - (pos + dv)) > 0:
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
	//fmt.Println("insert", string(key))

	parent.setNext(key[pos], &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  &[256]*Node{},
	})
}

func (r *Radix) updateNode(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("update", string(key))

	parent.setNext(key[pos], &Node{
		prefix: key[pos+1:],
		value:  value,
		edges:  node.edges,
	})
}

func (r *Radix) splitTwoWay(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	//fmt.Println("split two way", string(key))

	var pfx []byte

	// fix issue where key is found, but is occupied by another node with prefix
	if len(key) > pos+1 {
		pfx = key[pos+1 : pos+dv]
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

	parent.setNext(key[pos], n1)
}

func (r *Radix) splitThreeWay(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	//fmt.Println("split three way", string(key))

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

	parent.setNext(key[pos-1], n1)
}

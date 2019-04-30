package rad

import (
	"fmt"
	"strings"
)

// Radix tree
type Radix struct {
	root *Node
}

// New Radix Tree
func New() *Radix {
	return &Radix{
		root: &Node{},
	}
}

// Insert or update a keypair
func (r *Radix) Insert(key []byte, value interface{}) {
	parent, node, pos, dv := r.findInsertionPoint(key)

	switch {

	/*
		case len(node.prefix) > 0 && len(node.prefix) > dv:
			r.splitNode(key, parent, node, pos, dv)
		case pos == len(key):
			r.updateNode(key, parent, node, pos, dv)
		case n.HasPrefix() && dv == n.PrefixSize() || i < len(key):
			r.insertNode(key, parent, node, pos, db)
	*/

	case pos == len(key):
		r.updateNode(key, value, parent, node, pos, dv)
	case pos < len(key) && dv == 0:
		r.insertNode(key, value, parent, node, pos, dv)
	default:
		r.splitNode(key, value, parent, node, pos, dv)
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

// Graphviz returns a graphviz formatted string of all the nodes in the tree
// this should only be run on trees with relatively few nodes
func (r *Radix) Graphviz() string {
	var gvzc int

	gvoutput := []string{"digraph G {"}

	r.graphviz(&gvoutput, &gvzc, "[-1] ROOT", r.root)

	gvoutput = append(gvoutput, "}")

	return fmt.Sprint(strings.Join(gvoutput, "\n"))
}

func (r *Radix) graphviz(gvoutput *[]string, gvzc *int, previous string, n *Node) {
	for i, next := range n.edges {
		if next != nil {
			(*gvzc)++

			n = next

			(*gvoutput) = append((*gvoutput), fmt.Sprintf("  \"%s\" -> \"[%d] %s\" [label=\"%s\"]", previous, *gvzc, string(n.prefix), string(byte(i))))

			r.graphviz(gvoutput, gvzc, fmt.Sprintf("[%d] %s", *gvzc, string(n.prefix)), n)
		}
	}
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
	fmt.Println("insert")

	parent.setNext(key[pos], &Node{
		prefix: key[pos+1:],
		value:  value,
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

func (r *Radix) splitNode(key []byte, value interface{}, parent, node *Node, pos, dv int) {
	fmt.Println("split")

	parent.print()
	if node != nil {
		node.print()
	}

	fmt.Println("create node", string(key[pos]), "->", string(key[pos+1:pos+dv]))
	fmt.Println("create node", string(node.prefix[dv]), "->", string(node.prefix[dv+1:]))

	/*
		n1 := &Node{
			prefix:
		}
	*/

	fmt.Println(pos)
	fmt.Println(dv)
}
